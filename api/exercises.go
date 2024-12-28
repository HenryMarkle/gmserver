package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/HenryMarkle/gmserver/dto"
)

func GetAllSections(ctx *gin.Context) {
	sections, queryErr := db.GetAllExerciseSections(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get exercise sections: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, sections)
}

func GetSectionByName(ctx *gin.Context) {
	name := ctx.Params.ByName("name")

	if name == "" {
		ctx.String(http.StatusBadRequest, "Required query parameter: name")
		return
	}

	section, queryErr := db.GetExerciseSectionByNameWithExercises(db.DB, name)
	if queryErr != nil {
		common.Logger.Printf("Failed to get a section by name: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	exercises := make([]dto.Excercise_Res, 0, len(section.Excercises))
	for _, e := range section.Excercises {
		exercises = append(exercises, dto.Excercise_Res{
			Name:        e.Name,
			Description: e.Description,
			ID:          e.ID,
			CategoryID:  e.CategoryID,
		})
	}

	ctx.JSON(http.StatusOK, dto.ExcerciseCategory_Res{
		ID:        section.ID,
		Name:      section.Name,
		Exercises: exercises,
	})
}

func CreateSection(ctx *gin.Context) {
	name := ctx.Query("name")

	if name == "" {
		ctx.String(http.StatusBadRequest, "Required query parameter: name")
		return
	}

	id, queryErr := db.CreateExerciseSection(db.DB, name)
	if queryErr != nil {
		common.Logger.Printf("Failed to create a section: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func DeleteSection(ctx *gin.Context) {
	params := NonEmptyQueryInt64OrAbort(ctx, "id")
	if params == nil {
		return
	}

	id := params[0]

	queryErr := db.DeleteExerciseSectionByIDWithExercises(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete an exercise section by ID (id: %d): %v\n", id, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func UpdateSectionById(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: id")
		return
	}

	newName := ctx.Query("name")
	if newName == "" {
		ctx.String(http.StatusBadRequest, "Required query parameter: name")
	}

	queryErr := db.UpdateExerciseSectionByID(db.DB, db.ExcerciseCategory{
		ID:   id,
		Name: newName,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed to update exercise section by ID (id: %d): %v\n", id, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func DeleteSectionWithExercises(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invliad parameter: id")
		return
	}

	queryErr := db.DeleteExerciseSectionByIDWithExercises(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete a section with its exercises by ID (id: %d): %v\n", id, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func CountSectionExercises(ctx *gin.Context) {
	name := ctx.Params.ByName("name")
	if name == "" {
		ctx.String(http.StatusBadRequest, "Required query parameter: name")
		return
	}

	count, queryErr := db.CountExercisesOfExerciseSectionByName(db.DB, name)
	if queryErr != nil {
		common.Logger.Printf("Failed to count exercises of section '%s': %v\n", name, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, count)
}

func GetAllSectionsWithExcercises(ctx *gin.Context) {
	sections, queryErr := db.GetAllExerciseSectionsWithExercises(db.DB)
	if queryErr != nil {
		common.Logger.Printf("Failed to get sections with exercises: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	sectionsDto := make([]dto.ExcerciseCategory_Res, 0, len(sections))
	for _, s := range sections {
		exercisesDto := make([]dto.Excercise_Res, 0, len(s.Excercises))
		for _, e := range s.Excercises {
			exercisesDto = append(exercisesDto, dto.Excercise_Res{
				ID:          e.ID,
				Name:        e.Name,
				Description: e.Description,
				CategoryID:  e.CategoryID,
			})
		}

		sectionsDto = append(sectionsDto, dto.ExcerciseCategory_Res{
			ID:        s.ID,
			Name:      s.Name,
			Exercises: exercisesDto,
		})
	}

	ctx.JSON(http.StatusOK, sectionsDto)
}

func GetAllExcercises(ctx *gin.Context) {
	exercises, queryErr := db.GetAllExercises(db.DB)

	if queryErr != nil {
		common.Logger.Printf("Failed to get all exercises: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	exercisesDto := make([]dto.Excercise_Res, 0, len(exercises))
	for _, e := range exercises {
		exercisesDto = append(exercisesDto, dto.Excercise_Res{
			ID:          e.ID,
			Name:        e.Name,
			Description: e.Description,
			CategoryID:  e.CategoryID,
		})
	}

	ctx.JSON(http.StatusOK, exercisesDto)
}

func GetAllExcercisesOfSection(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: id")
		return
	}

	exercises, queryErr := db.GetAllExercisesOfSection(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to get exercises of section (id: %d): %v\n", id, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	exercisesDto := make([]dto.Excercise_Res, 0, len(exercises))
	for _, e := range exercises {
		exercisesDto = append(exercisesDto, dto.Excercise_Res{
			ID:          e.ID,
			Name:        e.Name,
			Description: e.Description,
			CategoryID:  e.CategoryID,
		})
	}

	ctx.JSON(http.StatusOK, exercisesDto)
}

func CreateExcercise(ctx *gin.Context) {
	data := dto.CreateExcercise_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	id, queryErr := db.CreateExercise(db.DB, db.Excercise{
		Name:        data.Name,
		Description: data.Description,
		CategoryID:  data.CategoryID,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed to create an exercise: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func DeleteExcercise(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.String(http.StatusBadRequest, "Required query parameter: name")
	}

	queryErr := db.DeleteExerciseByName(db.DB, name)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete an exercise '%s': %v\n", name, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func DeleteExcerciseById(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: id")
		return
	}

	queryErr := db.DeleteExerciseByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("Failed to delete an exercise by ID (id: %d): %v\n", id, queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func UpdateExcerciseById(ctx *gin.Context) {
	data := dto.UpdateExcercise_Req{}

	bindErr := ctx.ShouldBindJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	queryErr := db.UpdateExercise(db.DB, db.Excercise{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		CategoryID:  data.CategoryID,
	})
	if queryErr != nil {
		common.Logger.Printf("Failed to update exercise: %v\n", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
func UpdateExcerciseById2(ctx *gin.Context) {}
