package entity

type Defaultable[T any] interface {
	Default() *T
}
