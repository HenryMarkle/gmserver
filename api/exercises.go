package api

import "github.com/gin-gonic/gin"

func GetAllSections(ctx *gin.Context)               {}
func GetSectionByName(ctx *gin.Context)             {}
func CreateSection(ctx *gin.Context)                {}
func DeleteSection(ctx *gin.Context)                {}
func UpdateSectionById(ctx *gin.Context)            {}
func DeleteSectionWithExercises(ctx *gin.Context)   {}
func CountSectionExercises(ctx *gin.Context)        {}
func GetAllSectionsWithExcercises(ctx *gin.Context) {}
func GetAllExcercises(ctx *gin.Context)             {}
func GetAllExcercisesOfSection(ctx *gin.Context)    {}
func CreateExcercise(ctx *gin.Context)              {}
func DeleteExcercise(ctx *gin.Context)              {}
func DeleteExcerciseById(ctx *gin.Context)          {}
func UpdateExcerciseById(ctx *gin.Context)          {}
func UpdateExcerciseById2(ctx *gin.Context)         {}
