package handlers

import "gorm.io/gorm"

type Factory struct {
	db *gorm.DB
}

func NewFactory(db *gorm.DB) *Factory {
	return &Factory{db: db}
}

func GetByIds[T any](f *Factory) *ById[T] {
	return &ById[T]{f.db, withModel[T]{}}
}

func GetByFieldIds[T any](f *Factory, field string) *ByFieldIds[T] {
	return &ByFieldIds[T]{f.db, withModel[T]{}, field}
}

func GetByQuery[T any](f *Factory, fieldsMap FieldsMap) *ByQuery[T] {
	return &ByQuery[T]{f.db, withModel[T]{}, fieldsMap}
}

func FactoryUpdateByPut[T any](h *Factory, allowedFields AllowedFields) *UpdateByPut[T] {
	return NewUpdateByPut[T](h.db, allowedFields)
}

func FactoryCreateByPost[T any](h *Factory, allowedFields AllowedFields) *CreateByPost[T] {
	return NewCreateByPost[T](h.db, allowedFields)
}

func FactoryDeleteById[T any](f *Factory, unscoped bool) *DeleteById[T] {
	return &DeleteById[T]{f.db, withModel[T]{}, unscoped}
}
