package server

import (
	"github.com/gin-gonic/gin"
	"st-go/errors"
	"net/http"
	"strings"
)

const (
	JWTAuthorizationHeaderKey = "Authorization"
)

func SendErrorResponse(context *gin.Context, err error) {
	if error, ok := err.(errors.Error); ok {
		context.JSON(error.Status(), &Response{Error: error.Error()})
	} else {
		context.JSON(http.StatusOK, &Response{Error: err.Error()})
	}
	if context.IsAborted() == false {
		context.Abort()
	}
}

type JWTParser interface {
	Parse(*gin.Context, string) error
}

func Authorize(jwtParser JWTParser) gin.HandlerFunc {
	return func(context *gin.Context) {
		header := context.GetHeader(JWTAuthorizationHeaderKey)
		if header == "" {
			SendErrorResponse(context, errors.Unauthorized)
			return
		}
		components := strings.Split(header, " ")
		if len(components) != 2 {
			SendErrorResponse(context, errors.Unauthorized)
			return
		}
		tokenString := components[1]
		err := jwtParser.Parse(context, tokenString)
		if err != nil {
			SendErrorResponse(context, errors.Unauthorized)
			return
		}
		context.Next()
	}
}