package api

import "github.com/gin-gonic/gin"

func GetAllEvents(ctx *gin.Context)
func DidUserSeeEvent(ctx *gin.Context)
func MarkEventAsSeen(ctx *gin.Context)
func MarkAllEventsAsSeen(ctx *gin.Context)
