package net

type Range[T any] interface {
	End() T
	Start() T
}
