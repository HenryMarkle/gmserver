package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NonEmptyQueryInt64OrAbort(ctx *gin.Context, params ...string) []int64 {
	if ctx == nil || len(params) == 0 {
		return nil
	}

	errors := make([]string, 0, len(params))
	values := make([]int64, 0, len(params))

	for _, p := range params {
		value := ctx.Query(p)

		if value == "" {
			errors = append(errors, fmt.Sprintf("Missing query parameter: %s", p))
		} else {
			converted, convErr := strconv.ParseInt(value, 10, 64)
			if convErr != nil {
				errors = append(errors, fmt.Sprintf("Invalid query paramter: %s", p))
			} else if len(errors) == 0 {
				values = append(values, converted)
			}
		}
	}

	if len(errors) > 0 {
		ctx.JSON(http.StatusBadRequest, errors)
		return nil
	}

	return values
}
