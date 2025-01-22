package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/HenryMarkle/gmserver/dto"
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

	ctx.JSON(http.StatusOK, &dto.GetBlog_Res{
		ID:          blog.ID,
		Title:       blog.Title,
		Subtitle:    blog.Subtitle,
		Description: blog.Description,
		Views:       blog.Views,
	})
}

func GetAllBlogs(ctx *gin.Context) {
	blogs, queryErr := db.GetAllBlogs(db.DB)
	if queryErr != nil {
		common.Logger.Printf("failed to get all blogs: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	dtoBlogs := make([]dto.GetBlog_Res, 0, len(blogs))

	for _, i := range blogs {
		dtoBlogs = append(dtoBlogs, dto.GetBlog_Res{
			ID:          i.ID,
			Title:       i.Title,
			Subtitle:    i.Subtitle,
			Description: i.Description,
			Views:       i.Views,
		})
	}

	ctx.JSON(http.StatusOK, dtoBlogs)
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
	var imageBytes []byte
	if formErr == nil {
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
		imageBytes = buf.Bytes()
	}

	queryErr := db.UpdateBlogByID(db.DB, db.Blog{
		ID:          id,
		Title:       title,
		Subtitle:    subtitle,
		Description: description,
		Image:       imageBytes,
		Views:       int(views),
		ImageType:   filepath.Ext(image.Filename)[1:],
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
	var imageBytes []byte
	if formErr == nil {

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
		imageBytes = buf.Bytes()
	}

	id, queryErr := db.CreateBlog(db.DB, db.Blog{
		Title:       title,
		Subtitle:    subtitle,
		Description: description,
		Image:       imageBytes,
		Views:       int(views),
		ImageType:   fmt.Sprintf("image/%s", filepath.Ext(image.Filename)[1:]),
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

func GetBlogImageByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: 'id': %v", convErr)
	}

	blog, queryErr := db.GetBlogByID(db.DB, id)
	if queryErr != nil {
		common.Logger.Printf("failed to get a blog by ID: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fmt.Printf("TYPE: %s", blog.ImageType)

	ctx.Data(http.StatusOK, blog.ImageType, blog.Image)
}
