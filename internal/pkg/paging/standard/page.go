package paging_standard

type PageableRequest struct {
	Number int64
	Size   int64
}

type Paged[T any] struct {
	Content       []T
	Number        int64
	Size          int64
	TotalElements int64
	TotalPages    int64
}

func EmptyPage[T any]() Paged[T] {
	return Paged[T]{}
}
