package clients

import (
	"payment-service/domain/dto"
	"time"

	errConstants "payment-service/constants/error/payment"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type MidtransClient struct {
	ServerKey    string
	IsProduction bool
}

type IMidtransClient interface {
	CreatePaymentLink(request *dto.PaymentRequest) (*MidtransData, error)
}

func NewMidtransClient(serverKey string, isProduction bool) *MidtransClient {
	return &MidtransClient{
		ServerKey:    serverKey,
		IsProduction: isProduction,
	}
}

func (c *MidtransClient) CreatePaymentLink(request *dto.PaymentRequest) (*MidtransData, error) {
	var (
		snapClient   snap.Client
		isProduction = midtrans.Sandbox
	)

	expiryDateTime := request.ExpiredAt
	currentTime := time.Now()
	duration := expiryDateTime.Sub(currentTime)

	if duration <= 0 {
		logrus.Errorf("Expired at is invalid")
		return nil, errConstants.ErrExpireAtInvalid
	}

	expiryUnit := "minute"
	expirDuration := int64(duration.Minutes())

	if duration.Hours() >= 1 {
		expiryUnit = "hour"
		expirDuration = int64(duration.Hours())
	} else if duration.Hours() >= 24 {
		expiryUnit = "day"
		expirDuration = int64(duration.Hours() / 24)
	}

	if c.IsProduction {
		isProduction = midtrans.Production
	}

	snapClient.New(c.ServerKey, isProduction)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:     request.OrderID,
			GrossAmt: int64(request.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:  request.CustomerDetail.Name,
			Email: request.CustomerDetail.Email,
			Phone: request.CustomerDetail.Phone,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    request.ItemDetail[0].ID,
				Name:  request.ItemDetail[0].Name,
				Price: int64(request.ItemDetail[0].Amount),
				Qty:   int32(request.ItemDetail[0].Quantity),
			},
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     expiryUnit,
			Duration: expirDuration,
		},
	}

	response, err := snapClient.CreateTransaction(req)

	if err != nil {
		logrus.Errorf("Error create transaction: %v", err)
		return nil, err
	}

	return &MidtransData{
		RedirectURL: response.RedirectURL,
		Token:       response.Token,
	}, nil
}
