package api

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/gin-gonic/gin"
)

func GetBlogByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: 'id': %v", convErr)
	}

	blog, queryErr := db.GetBlogByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("failed to get blog by id: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, blog)
}

func GetAllBlogs(ctx *gin.Context) {
	blogs, queryErr := db.GetAllBlogs(db.DB)
	if queryErr != nil {
		common.Logger.Printf("failed to get all blogs: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, blogs)
}

func UpdateBlogByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: 'id': %v", convErr)
	}

	title := ctx.PostForm("title")
	subtitle := ctx.PostForm("subtitle")
	description := ctx.PostForm("description")
	viewsStr := ctx.PostForm("views")

	views, viewsConvErr := strconv.ParseInt(viewsStr, 10, 32)
	if viewsConvErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid form field: 'views': not a number", viewsConvErr)
	}

	image, formErr := ctx.FormFile("image")
	if formErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid form data: %v", formErr)
		return
	}

	imageContent, imageErr := image.Open()
	if imageErr != nil {
		common.Logger.Printf("failed to open recieved image content: %v", imageErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer imageContent.Close()

	var buf bytes.Buffer
	if _, readErr := io.Copy(&buf, imageContent); readErr != nil {
		common.Logger.Printf("failed to read recieved image content: %v", readErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	imageBytes := buf.Bytes()

	queryErr := db.UpdateBlogByID(db.DB, db.Blog{
		ID:          id,
		Title:       title,
		Subtitle:    subtitle,
		Description: description,
		Image:       imageBytes,
		Views:       int(views),
	})

	if queryErr != nil {
		common.Logger.Printf("failed to update blog by id: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func CreateBlog(ctx *gin.Context) {

	title := ctx.PostForm("title")
	subtitle := ctx.PostForm("subtitle")
	description := ctx.PostForm("description")
	viewsStr := ctx.PostForm("views")

	views, viewsConvErr := strconv.ParseInt(viewsStr, 10, 32)
	if viewsConvErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid form field: 'views': not a number", viewsConvErr)
	}

	image, formErr := ctx.FormFile("image")
	if formErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid form data: %v", formErr)
		return
	}

	imageContent, imageErr := image.Open()
	if imageErr != nil {
		common.Logger.Printf("failed to open recieved image content: %v", imageErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer imageContent.Close()

	var buf bytes.Buffer
	if _, readErr := io.Copy(&buf, imageContent); readErr != nil {
		common.Logger.Printf("failed to read recieved image content: %v", readErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	imageBytes := buf.Bytes()

	id, queryErr := db.CreateBlog(db.DB, db.Blog{
		Title:       title,
		Subtitle:    subtitle,
		Description: description,
		Image:       imageBytes,
		Views:       int(views),
	})

	if queryErr != nil {
		common.Logger.Printf("failed to create blog: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func DeleteBlogByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: 'id': %v", convErr)
	}

	queryErr := db.DeleteBlogByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("failed to delete blog by id: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
