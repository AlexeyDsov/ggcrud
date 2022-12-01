package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type ById[T any] struct {
	db *gorm.DB
	withModel[T]
}

func (h *ById[T]) ById(c *gin.Context) {
	var r struct {
		Id int `form:"id"`
	}

	err := c.BindUri(&r)
	if err != nil {
		panic(err)
	}

	model := h.modelInstancePointer()
	if err := h.db.Find(&model, r.Id).Error; err != nil {
		panic(err)
	}

	status := http.StatusOK
	if model == nil {
		status = http.StatusNotFound
	}

	c.JSON(status, model)
}

func (h *ById[T]) ByIds(c *gin.Context) {
	var ids []int
	err := c.BindJSON(&ids)

	if err != nil {
		panic(err)
	}

	model := h.modelSlice()
	if err := h.db.Where("id in ?", ids).Find(&model).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, model)
}
