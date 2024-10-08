package middleware

import (
	"github.com/gin-gonic/gin"
)

func Init(service []interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Keys = make(map[string]interface{})
		ctx.Keys["user"] = service[0]
		ctx.Next()
	}
}

//func Init(service ...interface{}) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		ctx.Keys = make(map[string]interface{})
//		ctx.Keys["user"] = service[0]
//		ctx.Next()
//	}
//}
