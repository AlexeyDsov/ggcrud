package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type DbApplyFunc func(*gorm.DB) (*gorm.DB, error)
type WhereFuncFactory func(value string) DbApplyFunc
type OrderFuncFactory func(bool) DbApplyFunc

type FieldsMap struct {
	_map               map[string]string
	whereFuncFactories map[string]WhereFuncFactory
	orderFuncFactories map[string]OrderFuncFactory
}

func NewFieldMap() FieldsMap {
	return FieldsMap{
		_map:               map[string]string{},
		whereFuncFactories: map[string]WhereFuncFactory{},
		orderFuncFactories: map[string]OrderFuncFactory{},
	}
}

func (f FieldsMap) AddSimple(field string) FieldsMap {
	f._map[field] = field
	return f
}

func (f FieldsMap) AddPair(alias, field string) FieldsMap {
	f._map[alias] = field
	return f
}

func (f FieldsMap) AddWhereFunc(alias string, f1 WhereFuncFactory) FieldsMap {
	f.whereFuncFactories[alias] = f1
	return f
}

func (f FieldsMap) AddWhereLike(alias string, field string) FieldsMap {
	f.whereFuncFactories[alias] = func(value string) DbApplyFunc {
		return func(db *gorm.DB) (*gorm.DB, error) {
			return db.Where(fmt.Sprintf("%s like ?", field), value), nil
		}
	}
	return f
}

func (f FieldsMap) AddWhereILike(alias string, field string) FieldsMap {
	f.whereFuncFactories[alias] = func(value string) DbApplyFunc {
		return func(db *gorm.DB) (*gorm.DB, error) {
			return db.Where(fmt.Sprintf("%s ilike ?", field), strings.ToLower(value)), nil
		}
	}
	return f
}

func (f FieldsMap) AddOrderFunc(alias string, f1 OrderFuncFactory) FieldsMap {
	f.orderFuncFactories[alias] = f1
	return f
}

type ByQuery[T any] struct {
	db *gorm.DB
	withModel[T]
	fieldMaps FieldsMap
}

func (b *ByQuery[T]) ByQuery(c *gin.Context) {
	var where []DbApplyFunc

	for from, to := range b.fieldMaps._map {
		if value, exists := c.GetQuery(from); exists {
			where = append(where, eqWhere(to, value))
		}
	}
	for from, whereFuncFactory := range b.fieldMaps.whereFuncFactories {
		if value, exists := c.GetQuery(from); exists {
			where = append(where, whereFuncFactory(value))
		}
	}

	limit, _ := b.toInt(c.DefaultQuery("_limit", "100"), 100)
	offset, _ := b.toInt(c.DefaultQuery("_offset", "0"), 0)
	orderByFuncs, err := b.parseOrderBy(c)
	if err != nil {
		log.Printf("error in order by funcs parsing: %s", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	db, err := b.fillWhere(b.db, where)
	if err != nil {
		log.Printf("error in where: %s", err)
		c.JSON(http.StatusBadGateway, nil)
		return
	}
	countDb, err := b.fillWhere(b.db, where)
	if err != nil {
		log.Printf("error in where for count: %s", err)
		c.JSON(http.StatusBadGateway, nil)
		return
	}

	db = db.Limit(limit).Offset(offset)
	for _, orderByFunc := range orderByFuncs {
		db, err = orderByFunc(db)
		if err != nil {
			log.Printf("error applying order by func: %s", err)
			c.JSON(http.StatusBadGateway, nil)
			return
		}
	}

	models := b.modelSlice()
	if err := db.Find(&models).Error; err != nil {
		log.Printf("error in executing query: %s", err)
		c.JSON(http.StatusBadGateway, nil)
		return
	}

	var count int64
	if err := countDb.Model(b.modelInstance()).Count(&count).Error; err != nil {
		log.Printf("error in executing count query: %s", err)
		c.JSON(http.StatusBadGateway, nil)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"items": models, "count": count})
}

func (b *ByQuery[T]) toInt(val string, def int) (int, bool) {
	if i, err := strconv.Atoi(val); err == nil {
		return i, true
	} else {
		return def, false
	}
}

func (b *ByQuery[T]) parseOrderBy(c *gin.Context) ([]DbApplyFunc, error) {
	orderByStr := c.Query("_orderBy")

	hasIdOrder := false
	var orderBy []DbApplyFunc
	for _, orderFieldStr := range strings.Split(orderByStr, ",") {
		orderFieldStr = strings.TrimSpace(orderFieldStr)
		if orderFieldStr == "" {
			continue
		}

		isDesc := false
		if strings.HasPrefix(orderFieldStr, "-") {
			isDesc = true
			orderFieldStr = strings.TrimPrefix(orderFieldStr, "-")
		}

		if field, e := b.fieldMaps._map[orderFieldStr]; e {
			if field == "id" {
				hasIdOrder = true
			}

			orderBy = append(orderBy, orderByField(field, isDesc))
		} else if orderFunc, e := b.fieldMaps.orderFuncFactories[orderFieldStr]; e {
			orderBy = append(orderBy, orderFunc(isDesc))
		} else {
			return nil, fmt.Errorf("could not find order field/function: %s", orderFieldStr)
		}
	}

	if !hasIdOrder {
		orderBy = append(orderBy, orderByField("id", false))
	}

	return orderBy, nil
}

func (b *ByQuery[T]) fillWhere(db *gorm.DB, wheres []DbApplyFunc) (*gorm.DB, error) {
	var err error
	for _, where := range wheres {
		db, err = where(db)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func eqWhere(field, value string) DbApplyFunc {
	return func(db *gorm.DB) (*gorm.DB, error) {
		return db.Where(fmt.Sprintf("%s = ?", field), value), nil
	}
}

func orderByField(field string, isDesk bool) DbApplyFunc {
	return func(db *gorm.DB) (*gorm.DB, error) {
		if isDesk {
			field += " desc"
		}
		return db.Order(field), nil
	}
}
