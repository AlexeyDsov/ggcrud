package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type DeleteById[T any] struct {
	db *gorm.DB
	withModel[T]
	unscoped bool
}

func (d *DeleteById[T]) Delete(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	db := d.db
	if d.unscoped {
		db = db.Unscoped()
	}

	if err := db.Delete(d.modelInstance(), id).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, nil)
}
