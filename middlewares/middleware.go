package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"payment-service/clients"
	"payment-service/common/response"
	"payment-service/config"
	"payment-service/constants"
	errConstant "payment-service/constants/error"
	"strings"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HandlePanic is a middleware that recovers from panics and returns a 500 Internal Server Error response.
func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic:: %v", r)

				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errConstant.ErrInternalServerError.Error(),
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}

// RateLimiter is a middleware that applies rate limiting to incoming requests using the provided limiter.
func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement rate limiting logic here
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)

		if err != nil {
			c.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: errConstant.ErrTooManyRequests.Error(),
			})
			c.Abort()
		}

		c.Next()
	}
}

// extractBearerToken extracts the Bearer token from the Authorization header.
func extractBearerToken(token string) string {
	// Split the token into parts
	parts := strings.SplitN(token, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})

	c.Abort()
}

// validateAPIKey validates the API key from the request headers.
func validateAPIKey(c *gin.Context) error {
	apiKey := c.GetHeader(constants.XAPIKey)

	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XServiceName)
	signatureKey := config.Config.SignatureKey

	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)

	hash := sha256.New()

	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	// Validate the API key (this is just a placeholder, implement your own logic)
	if apiKey != resultHash {
		return errConstant.ErrUnauthorized
	}

	return nil
}

// contains checks if a role is present in the list of roles.
func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

// CheckRole checks if the user has one of the required roles.
func CheckRole(roles []string, client clients.IClientRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := client.GetUser().GetUserByToken(c.Request.Context())

		if err != nil {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		if !contains(roles, user.Role) {
			responseUnauthorized(c, errConstant.ErrForbidden.Error())
			return
		}

		c.Next()
	}
}

// Authenticate checks if the request is authenticated.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		token := c.GetHeader(constants.Authorization)

		if token == "" {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		err = validateAPIKey(c)

		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		c.Next()
	}
}

// AuthenticateWithoutToken checks if the request is authenticated without a token.
// This middleware is used for public endpoints.
func AuthenticateWithoutToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := validateAPIKey(c)

		if err != nil {
			responseUnauthorized(c, err.Error())
			return
		}

		c.Next()
	}
}
