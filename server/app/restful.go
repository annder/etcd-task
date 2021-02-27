package app

import "github.com/gin-gonic/gin"

func Ok (ctx *gin.Context,data interface{}, message string) {
	ctx.JSON(200,gin.H{
		"data":data,
		"message":message,
		"code":200,
	})
}

func Err(ctx *gin.Context,code int, message string){
	ctx.JSON(200,gin.H{
		"data":gin.H{},
		"message":message,
		"code":code,
	})
}
