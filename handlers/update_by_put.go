package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type UpdateByPut[T any] struct {
	db *gorm.DB
	withAllowedFields[T]
}

func NewUpdateByPut[T any](db *gorm.DB, allowedFields AllowedFields) *UpdateByPut[T] {
	return &UpdateByPut[T]{
		db:                db,
		withAllowedFields: withAllowedFields[T]{withModel: withModel[T]{}, allowedFields: allowedFields},
	}
}

func (u *UpdateByPut[T]) UpdateByStrField(param string, search func(*gorm.DB, string, *T) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.Param(param)
		if value == "" {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("empty param %s", param))
			return
		}

		dbFields, jsonString, err := u.checkAndGetAllowedFields(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		model := u.modelInstance()
		if err := search(u.db, value, &model); err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, nil)
			return
		} else if err != nil {
			c.JSON(http.StatusBadGateway, nil)
			log.Printf("error happened getting model: %s", err)
			return
		}

		if err := json.Unmarshal(jsonString, &model); err != nil {
			c.JSON(http.StatusBadGateway, nil)
			log.Printf("error happened unmarshalling to model: %s", err)
			return
		}

		if err := u.db.Model(&model).Select(dbFields[0], dbFields[1:]).Updates(model).Error; err != nil {
			c.JSON(http.StatusBadGateway, nil)
			log.Printf("error happened updating model: %s", err)
			return
		}

		c.JSON(http.StatusOK, model)
	}
}

func (u *UpdateByPut[T]) Update(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	dbFields, jsonString, err := u.checkAndGetAllowedFields(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	model := u.modelInstance()
	if err := u.db.Limit(1).Find(&model, id).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return
	} else if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(jsonString, &model); err != nil {
		panic(err)
	}

	if err := u.db.Model(&model).Select(dbFields[0], dbFields[1:]).Updates(model).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, model)
}
