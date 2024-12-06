package main

type Optional[T any] struct {
	Defined bool
	value   T
}

func (o Optional[T]) IsDefined() bool {
	return o.Defined
}

func (o Optional[T]) GetValue() (T, bool) {
	return o.value, o.Defined
}

func (o *Optional[T]) SetValue(value T) {
	o.Defined = true
	o.value = value
}

func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{Defined: true, value: value}
}

func NewOptionalEmpty[T any]() Optional[T] {
	return Optional[T]{Defined: false}
}
