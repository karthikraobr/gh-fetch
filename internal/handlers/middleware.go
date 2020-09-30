package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type HttpError struct {
	status     int
	innerError error
}

func (e *HttpError) Error() string {
	return e.innerError.Error()
}

func NewHttpError(status int, err error) *HttpError {
	return &HttpError{status: status, innerError: err}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		var httpErr *HttpError
		if errors.As(errors.Cause(err.Err), &httpErr) {
			c.JSON(httpErr.status, gin.H{
				"error": httpErr.Error(),
			})
			return
		}
	}
}
