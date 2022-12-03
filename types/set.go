package types

type Set[T comparable] map[T]struct{}

func (set *Set[T]) Add(item T) {
	(*set)[item] = struct{}{}
}

func (set *Set[T]) Delete(item T) {
	delete(*set, item)
}

func (set *Set[T]) Has(item T) bool {
	_, ok := (*set)[item]

	return ok
}
