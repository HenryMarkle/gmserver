package api

import "github.com/gin-gonic/gin"

func GetAllComments(ctx *gin.Context)             {}
func GetAllCommentsOfManager(ctx *gin.Context)    {}
func GetAllCommentsOfSubscriber(ctx *gin.Context) {}
func CreateComment(ctx *gin.Context)              {}
func DeleteComment(ctx *gin.Context)              {}
