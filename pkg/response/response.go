package response

import "github.com/gin-gonic/gin"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Total     int `json:"total"`
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalPage int `json:"total_page"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, status int, message string, data interface{}, meta Meta) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

func Error(c *gin.Context, status int, message string, err string) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Error:   err,
	})
}
