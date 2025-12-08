package common

import (
	"iter"
	"sort"
)

type multiValueMapImpl[K comparable, V any] map[K][]V

func NewMultiValueMap[K comparable, V any]() MultiValueMap[K, V] {
	return multiValueMapImpl[K, V](map[K][]V{})
}

func (m multiValueMapImpl[K, V]) Put(key K, value V) {
	m[key] = append(m[key], value)
}

func (m multiValueMapImpl[K, V]) Get(key K) ([]V, bool) {
	values, ok := m[key]
	return values, ok
}

func (m multiValueMapImpl[K, V]) Sort(key K, orderBy func(values []V, i int, j int) bool) {
	list, ok := m.Get(key)
	if !ok {
		return
	}

	sort.Slice(list, func(i, j int) bool {
		return orderBy(list, i, j)
	})
	m[key] = list
}

func (m multiValueMapImpl[K, V]) SortAll(orderBy func(values []V, i int, j int) bool) {
	for key := range m {
		m.Sort(key, orderBy)
	}
}

func (m multiValueMapImpl[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range m {
			if !yield(key) {
				return
			}
		}
	}
}

func (m multiValueMapImpl[K, V]) Len() int {
	return len(m)
}

func (m multiValueMapImpl[K, V]) IsEmpty() bool {
	return len(m) == 0
}

func (m multiValueMapImpl[K, V]) Entries() iter.Seq2[K, []V] {
	return func(yield func(K, []V) bool) {
		for key, values := range m {
			if !yield(key, values) {
				return
			}
		}
	}
}

type setImpl[K comparable] map[K]*struct{}

func (s setImpl[K]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range s {
			if !yield(key) {
				return
			}
		}
	}
}

func (s setImpl[K]) Contains(key K) bool {
	_, ok := s[key]
	return ok
}

func (s setImpl[K]) Len() int {
	return len(s)
}

func (s setImpl[K]) IsEmpty() bool {
	return len(s) == 0
}

func (s setImpl[K]) Put(key K) {
	s[key] = &struct{}{}
}

func (s setImpl[K]) PutAll(keys ...K) {
	for _, key := range keys {
		s.Put(key)
	}
}

func NewSet[K comparable]() Set[K] {
	return setImpl[K](map[K]*struct{}{})
}

func NewImmutableSet[K comparable](keys ...K) ImmutableSet[K] {
	setObj := NewSet[K]()
	setObj.PutAll(keys...)

	return setObj
}
