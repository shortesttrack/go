package httpserver

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func QueryString(context *gin.Context, key string) *string {
	value := context.Query(key)
	if value == `` {
		return nil
	}
	return &value
}

func QueryBool(context *gin.Context, k string) *bool {
	v := context.Query(k)
	if v == `` {
		return nil
	}
	r, e := strconv.ParseBool(v)
	if e != nil {
		return nil
	}
	return &r
}

func QueryInt(context *gin.Context, k string, fallback int) int {
	v := context.Query(k)
	if v == `` {
		return fallback
	}
	i, e := strconv.Atoi(v)
	if e != nil {
		return fallback
	}
	return i
}