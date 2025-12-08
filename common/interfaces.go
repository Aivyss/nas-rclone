package common

import "iter"

type CheckableCollection interface {
	Len() int
	IsEmpty() bool
}

type ImmutableMultiValueMap[K comparable, V any] interface {
	CheckableCollection
	Get(key K) ([]V, bool)
	Keys() iter.Seq[K]
	Entries() iter.Seq2[K, []V]
}

type MultiValueMap[K comparable, V any] interface {
	ImmutableMultiValueMap[K, V]
	Put(key K, value V)
	Sort(key K, orderBy func(values []V, i, j int) bool)
	SortAll(orderBy func(values []V, i, j int) bool)
}

type ImmutableSet[K comparable] interface {
	CheckableCollection
	Keys() iter.Seq[K]
	Contains(key K) bool
}

type Set[K comparable] interface {
	ImmutableSet[K]
	Put(key K)
	PutAll(keys ...K)
}
