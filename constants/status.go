package constants

type PaymentStatus int
type PaymentStatusString string

const (
	Initial    PaymentStatus = 0
	Pending    PaymentStatus = 100
	Settlement PaymentStatus = 200
	Expire     PaymentStatus = 300

	InitialString    PaymentStatusString = "Initial"
	PendingString    PaymentStatusString = "Pending"
	SettlementString PaymentStatusString = "Settlement"
	ExpireString     PaymentStatusString = "Expire"
)

var mapPaymentStatusStringToInt = map[PaymentStatusString]PaymentStatus{
	InitialString:    Initial,
	PendingString:    Pending,
	SettlementString: Settlement,
	ExpireString:     Expire,
}

var mapPaymentStatusIntToString = map[PaymentStatus]PaymentStatusString{
	Initial:    InitialString,
	Pending:    PendingString,
	Settlement: SettlementString,
	Expire:     ExpireString,
}

func (p PaymentStatusString) String() string {
	return string(p)
}

func (p PaymentStatus) Int() int {
	return int(p)
}

func (p PaymentStatus) GetStatusString() PaymentStatusString {
	return mapPaymentStatusIntToString[p]
}

func (p PaymentStatusString) GetStatusInt() PaymentStatus {
	return mapPaymentStatusStringToInt[p]
}
