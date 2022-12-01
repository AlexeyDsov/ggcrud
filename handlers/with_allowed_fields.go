package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type AllowedFields map[string][]string

func NewAllowedFields() AllowedFields {
	return AllowedFields{}
}

func (a AllowedFields) AddId(field string) AllowedFields {
	return a.Add(field, field)
}

func (a AllowedFields) Add(from string, to ...string) AllowedFields {
	a[from] = to
	return a
}

func (a AllowedFields) Map(fields []string) []string {
	var mapped []string
	for _, field := range fields {
		if toFields, found := a[field]; found {
			mapped = append(mapped, toFields...)
		}
	}

	return mapped
}

type withAllowedFields[T any] struct {
	withModel[T]

	allowedFields AllowedFields
}

func (u *withAllowedFields[T]) checkAndGetAllowedFields(c *gin.Context) ([]string, []byte, error) {
	var fields map[string]interface{}
	if err := c.BindJSON(&fields); err != nil {
		return nil, nil, err
	}

	var nonAllowedFields []string
	var requestFields []string
	for field, _ := range fields {
		if _, found := u.allowedFields[field]; !found {
			nonAllowedFields = append(nonAllowedFields, field)
		} else {
			requestFields = append(requestFields, field)
		}
	}

	if len(nonAllowedFields) > 0 {
		return nil, nil, fmt.Errorf("not allowed fields: %s", strings.Join(nonAllowedFields, ", "))
	}

	jsonString, err := json.Marshal(fields)
	return u.allowedFields.Map(requestFields), jsonString, err
}
