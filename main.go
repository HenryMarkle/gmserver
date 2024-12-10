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
				comments := auth.Group("/comments")

				_ = comments.GET("/all", api.GetAllComments)
				_ = comments.GET("/user/:id", api.GetAllCommentsOfManager)
				_ = comments.GET("/sub/:id", api.GetAllCommentsOfSubscriber)
				_ = comments.POST("/new", api.CreateComment)
				_ = comments.DELETE("/:id", api.DeleteComment)
			}
			{
				customers := auth.Group("/customers")

				_ = customers.GET("/total-income", api.GetTotalIncome)
				_ = customers.GET("/count-customers", api.CountCustomers)
				_ = customers.PUT("/count-ending", api.CountCustomersEndingIn)
				_ = customers.GET("/count-expired", api.CountCustomersExpiring)
				_ = customers.POST("/new", api.CreateCustomer)
				_ = customers.GET("/all", api.GetAllComments)
				_ = customers.GET("/:id", api.GetCustomerByID)
				_ = customers.DELETE("/:id", api.DeleteCustomerByID)
				_ = customers.DELETE("delist/:id", api.MarkCustomerAsDeleted)
				_ = customers.PATCH("/", api.UpdateCustomerByID)
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
