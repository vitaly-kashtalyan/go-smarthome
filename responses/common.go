package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpOkData struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"OK"`
	Data    interface{} `json:"data" example:"interface{}"`
}

type HttpCreate struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"OK"`
}

type HttpError struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}

func ResOkData(ctx *gin.Context, data interface{}) {
	ok := HttpOkData{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, ok)
}

func ResOkPagination(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func ResCreate(ctx *gin.Context) {
	ok := HttpCreate{
		Status:  http.StatusCreated,
		Message: http.StatusText(http.StatusCreated),
	}
	ctx.JSON(http.StatusOK, ok)
}

func ResStatus(ctx *gin.Context, code int) {
	er := HttpError{
		Status:  code,
		Message: http.StatusText(code),
	}
	ctx.JSON(http.StatusOK, er)
}

func ResError(ctx *gin.Context, code int, message error) {
	er := HttpError{
		Status:  code,
		Message: message.Error(),
	}
	ctx.JSON(http.StatusOK, er)
}
