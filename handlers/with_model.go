package handlers

type withModel[T any] struct {
}

func (w *withModel[T]) modelInstance() T {
	var model T
	return model
}

func (w *withModel[T]) modelInstancePointer() *T {
	var model T
	return &model
}

func (w *withModel[T]) modelSlice() []T {
	var models []T
	return models
}
