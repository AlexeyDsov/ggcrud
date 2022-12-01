package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type ByFieldIds[T any] struct {
	db *gorm.DB
	withModel[T]
	field string
}

func (h *ByFieldIds[T]) ByFieldIds(c *gin.Context) {
	var ids []int
	err := c.BindJSON(&ids)

	if err != nil {
		panic(err)
	}

	models := h.modelSlice()
	if err := h.db.Where(fmt.Sprintf("%s in ?", h.field), ids).Find(&models).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, models)
}

func (h *ByFieldIds[T]) ByFieldStringIds(c *gin.Context) {
	var ids []string
	err := c.BindJSON(&ids)

	if err != nil {
		panic(err)
	}

	models := h.modelSlice()
	if err := h.db.Where(fmt.Sprintf("%s in ?", h.field), ids).Find(&models).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, models)
}
