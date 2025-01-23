package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/HenryMarkle/gmserver/common"
	"github.com/HenryMarkle/gmserver/db"
	"github.com/HenryMarkle/gmserver/dto"
	"github.com/gin-gonic/gin"
)

func findBlogImage(id int64) (string, bool, error) {
	entries, lookErr := os.ReadDir(filepath.Join(common.StoragePath, "blogs"))

	if lookErr != nil {
		return "", false, lookErr
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()

		if strings.HasPrefix(name, fmt.Sprintf("%d", id)) {
			return filepath.Join(common.StoragePath, "blogs", name), true, nil
		}
	}

	return "", false, nil
}

func deleteBlogImage(id int64) error {
	entries, lookErr := os.ReadDir(filepath.Join(common.StoragePath, "blogs"))

	if lookErr != nil {
		return lookErr
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()

		if strings.HasPrefix(name, fmt.Sprintf("%d", id)) {
			os.Remove(filepath.Join(common.StoragePath, "blogs", name))
			return nil
		}
	}

	return nil
}

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

	data := dto.UpdateBlog_Req{}
	bindErr := ctx.ShouldBindBodyWithJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	queryErr := db.UpdateBlogByID(db.DB, db.Blog{
		ID:          data.ID,
		Title:       data.Title,
		Subtitle:    data.Subtitle,
		Description: data.Description,
		Views:       data.Views,
	})

	if queryErr != nil {
		common.Logger.Printf("failed to update blog by id: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func CreateBlog(ctx *gin.Context) {

	data := dto.CreateBlog_Req{}
	bindErr := ctx.ShouldBindBodyWithJSON(&data)
	if bindErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid data: %v", bindErr)
		return
	}

	id, queryErr := db.CreateBlog(db.DB, db.Blog{
		Title:       data.Title,
		Subtitle:    data.Subtitle,
		Description: data.Description,
		Views:       data.Views,
	})

	if queryErr != nil {
		common.Logger.Printf("failed to create blog: %v", queryErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func UploadBlogImage(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: 'id': %v", convErr)
	}

	image, formErr := ctx.FormFile("image")
	var imageExt string
	if formErr != nil {
		ctx.String(http.StatusBadRequest, "Required file: 'image'.")
		return
	}

	if image.Size > 10<<20 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, "Image size is too large (max is 10 MB).")
		return
	}

	imageExt = filepath.Ext(image.Filename)

	if imageExt != ".png" && imageExt != ".jpg" && imageExt != ".jpeg" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, "Only files with extensions '.png', '.jpg' and '.jpeg' are allowed.")
		return
	}

	deleteErr := deleteBlogImage(id)
	if deleteErr != nil {
		common.Logger.Printf("failed to delete previous blog image (id: %d): %v", id, deleteErr)
	}

	imagePath := filepath.Join(common.StoragePath, "blogs", fmt.Sprintf("%d%s", id, imageExt))

	uploadErr := ctx.SaveUploadedFile(image, imagePath)
	if uploadErr != nil {
		common.Logger.Printf("failed to upload blog image (id: %d): %v", id, uploadErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
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

	deleteErr := deleteBlogImage(id)
	if deleteErr != nil {
		common.Logger.Printf("failed to delete previous blog image (id: %d): %v", id, deleteErr)
	}

	ctx.Status(http.StatusOK)
}

func GetBlogImageByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")

	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		ctx.String(http.StatusBadRequest, "Invalid parameter: 'id': %v", convErr)
	}

	path, found, lookErr := findBlogImage(id)

	if lookErr != nil {
		common.Logger.Printf("failed to find blog image (id: %d): %v", id, lookErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !found {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		common.Logger.Printf("failed to read blog image file (id: %d): %v", id, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Data(http.StatusOK, fmt.Sprintf("image/%s", filepath.Ext(path)[1:]), bytes)
}
