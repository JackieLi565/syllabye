package util

import "github.com/oapi-codegen/nullable"

// DefaultNil returns a pointer to the value if the expression is true, otherwise nil.
func DefaultNil[T any](expression bool, value *T) *T {
	if expression {
		return value
	}
	return nil
}

// DefaultNullable returns a nullable value struct if the expression is true, otherwise a null nullable struct.
func DefaultNullable[T any](expression bool, value T) nullable.Nullable[T] {
	if expression {
		return nullable.NewNullableWithValue(value)
	}

	return nullable.NewNullNullable[T]()
}
