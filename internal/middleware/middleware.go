package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"strings"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"testing-demo/internal/logging"
)

// This is very similar to `ginzap.Ginzap`, but we wanted  more fields than
// that recorded, so this middleware was created.
func RequestLogger(timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ginzap.Ginzap(logging.GetContextLogger(c).Desugar().Named("request"), timeFormat, utc)(c)
	}
}

// Injects a `*zap.SugaredLogger` into the context, complete with
// the `operation-id` key, as a convenience for the APIs that need a logger
func ContextLogger(logger *zap.SugaredLogger) gin.HandlerFunc {
	sugar := logger
	return func(c *gin.Context) {
		operationID := c.GetString("operationID")
		if operationID == "" {
			operationID, _ = newUUID()
		}
		s := sugar.With(zap.String("operation-id", operationID))
		c.Set("logger", s)
	}
}

func GetContextLogger(c context.Context) (logger *zap.SugaredLogger) {
	if val := c.Value("logger"); val != nil {
		logger, _ = val.(*zap.SugaredLogger)
	}
	return
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		errs := c.Errors
		if len(errs) > 0 {
			// We should only ever have one error in context by the time we get here...
			// We remove the error from the context, since we handle and log it here
			var err error
			err, c.Errors = errs[0].Err, c.Errors[1:]
			logger := GetContextLogger(c)
			if err != nil {
				logger.Error(err)
				if strings.Contains(err.Error(), "can't be parsed") {
					c.IndentedJSON(
						http.StatusBadRequest,
						err.Error(),
					)
				} else {
					c.IndentedJSON(
						http.StatusInternalServerError,
						err.Error(),
					)
				}
			}
		}
	}
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
