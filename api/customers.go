package api

import "github.com/gin-gonic/gin"

func GetTotalIncome(ctx *gin.Context)
func CountCustomers(ctx *gin.Context)
func CountCustomersEndingIn(ctx *gin.Context)
func CountCustomersExpiring(ctx *gin.Context)
func CreateCustomer(ctx *gin.Context)
func GetAllCustomers(ctx *gin.Context)
func GetCustomerByID(ctx *gin.Context)
func DeleteCustomerByID(ctx *gin.Context)
func MarkCustomerAsDeleted(ctx *gin.Context)
func UpdateCustomerByID(ctx *gin.Context)
