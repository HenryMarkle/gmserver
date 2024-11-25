package api

import "github.com/gin-gonic/gin"

func addUser(ctx *gin.Context)
func isUserSignedIn(ctx *gin.Context)
func getCurrentUserId(ctx *gin.Context)
func getTotalSalaries(ctx *gin.Context)
func updateUser(ctx *gin.Context)
func changeUserName(ctx *gin.Context)
func deleteUser(ctx *gin.Context)
func deleteUserById(ctx *gin.Context)
func countUsers(ctx *gin.Context)
func getUserByEmail(ctx *gin.Context)
func getUsersLeftChartData(ctx *gin.Context)
func getUsersCreatedChartData(ctx *gin.Context)
func getUserById(ctx *gin.Context)
func getGymName(ctx *gin.Context)
func getCurrentUser(ctx *gin.Context)
func changeGymName(ctx *gin.Context)
func getAllUsers(ctx *gin.Context)
func getAllAnnouncments(ctx *gin.Context)
func createAnnouncement(ctx *gin.Context)
func markAsRead(ctx *gin.Context)
