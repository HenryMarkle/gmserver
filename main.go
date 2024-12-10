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
				_ = customers.DELETE("/delist/:id", api.MarkCustomerAsDeleted)
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
				dash := auth.Group("/dashboard")

				_ = dash.GET("/home", api.GetHomeInfo)
				_ = dash.GET("/general", api.GetHomeGeneralInfo)
				_ = dash.PATCH("/general", api.UpdateHomeGeneralInfo)
				_ = dash.GET("/plan-paragraph", api.GetPlanParagrarph)
				_ = dash.PATCH("/plan-paragraph", api.UpdatePlanParagraph)
				_ = dash.GET("/plans", api.GetHomePlans)
				_ = dash.GET("/plan/:id", api.GetPlanByID)
				_ = dash.POST("/plan/new", api.CreatePlan)
				_ = dash.DELETE("/plan/:id", api.DeletePlanByID)
				_ = dash.PATCH("/plan", api.ReplacePlanByID)
				_ = dash.GET("/ads", api.GetAdsInfo)
				_ = dash.PATCH("/ads", api.UpdateAdsInfo)
				_ = dash.GET("/products", api.GetHomeProducts)
				_ = dash.GET("/product/:id", api.GetProductByID)
				_ = dash.POST("/product/new", api.CreateHomeProduct)
				_ = dash.DELETE("/product/:id", api.DeleteHomeProductByID)
				_ = dash.PATCH("/product", api.UpdateHomeProduct)
				_ = dash.GET("/product/categories", api.GetProductCategories)
				_ = dash.GET("/products-in-categories", api.GetCategoryProducts)
				_ = dash.GET("/products/category/:id", api.GetProductsOfCategory)
				_ = dash.POST("/product-category/new", api.CreateProductCategory)
				_ = dash.DELETE("/product-category/:id", api.DeleteProductCategoryByID)
				_ = dash.DELETE("/products-of-category/:id", api.DeleteProductsOfCategory)
				_ = dash.GET("/products-exists-in-category/:productId", api.ProductExistsUnderCategory)
				_ = dash.GET("/contacts", api.GetContacts)
				_ = dash.PATCH("/contacts", api.GetContacts)
			}

			{
				admin := auth.Group("/admin")
				admin.Use(api.AdminOnly())
			}
		}
	}

	server.Run()
}
