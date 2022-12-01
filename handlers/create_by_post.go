package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type CreateByPost[T any] struct {
	db *gorm.DB
	withAllowedFields[T]
}

func NewCreateByPost[T any](db *gorm.DB, allowedFields AllowedFields) *CreateByPost[T] {
	return &CreateByPost[T]{
		db:                db,
		withAllowedFields: withAllowedFields[T]{withModel: withModel[T]{}, allowedFields: allowedFields},
	}
}

func (u *CreateByPost[T]) Create(c *gin.Context) {
	dbFields, jsonString, err := u.checkAndGetAllowedFields(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	model := u.modelInstance()

	if err := json.Unmarshal(jsonString, &model); err != nil {
		panic(err)
	}

	if err := u.db.Model(&model).Select(dbFields[0], dbFields[1:]).Create(&model).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, model)
}
