package main

import (
	"github.com/HenryMarkle/gmserver/api"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/gin-gonic/gin"
)

func main() {
	defer db.DB.Close()
	server := gin.Default()

	server.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello")
	})

	{
		v1 := server.Group("/v1")

		v1.POST("/signin", api.SignIn)
		v1.POST("/signout", api.Signout)
	}

	server.Run()
}
