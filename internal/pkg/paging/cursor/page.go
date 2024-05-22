package paging_cursor

type PageableRequest struct {
	Size   int64
	Cursor string
	// Should be Direction field (next/prev page).
	// In this benchmark does not matter, only paging 1 -> 2 -> ... supported.
}

type Paged[T any] struct {
	Content        []T
	CursorNextPage string
	Size           int64
	TotalElements  int64
	TotalPages     int64
}

func EmptyPage[T any]() Paged[T] {
	return Paged[T]{}
}
