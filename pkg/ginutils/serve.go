package ginutil

import "github.com/gin-gonic/gin"

func ServeFileAttachment(path, name string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.File(path)
	}
}
