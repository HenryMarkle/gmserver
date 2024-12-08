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

		{
			auth := v1.Group("/auth")
			auth.Use(api.Auth())
			auth.POST("/signout", api.Signout)
			auth.PATCH("/changepassword", api.ChangePassword)

			{
				_ = auth.Group("/comments")
			}
			{
				_ = auth.Group("/customers")
			}
			{
				_ = auth.Group("/event")
			}
			{
				_ = auth.Group("/users")
			}
			{
				_ = auth.Group("/trainers")
			}
			{
				_ = auth.Group("/exercises")
			}
			{
				_ = auth.Group("/dashboard")
			}

			{
				admin := auth.Group("/admin")
				admin.Use(api.AdminOnly())
			}
		}
	}

	server.Run()
}
