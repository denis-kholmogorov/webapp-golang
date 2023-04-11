package errorhandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if marshalErr, ok := err.(MarshalError); ok {
					fmt.Printf("MarshalError error: %v\n", marshalErr)
					c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf("Bad request"))

				} else if notFoundError, ok := err.(NotFoundError); ok {
					fmt.Printf("NotFoundError error: %v\n", notFoundError)
					c.AbortWithStatusJSON(http.StatusNotFound, fmt.Sprintf("Object not found"))
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				}
				return
			}
		}()
		c.Next()
	}
}

type MarshalError struct {
	Message string
}

func (e *MarshalError) Error() string {
	return fmt.Sprintf("MarshalError: %s", e.Message)
}

type ErrorResponse struct {
	Message string
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("ErrorResponse: %s", e.Message)
}

type DbError struct {
	Message string
}

func (e *DbError) Error() string {
	return fmt.Sprintf("DbError: %s", e.Message)
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Object not found %s", e.Message)
}
