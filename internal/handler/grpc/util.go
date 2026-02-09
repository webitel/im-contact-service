package grpc

const (
	DefaultPageSize = 10
)

func ParsePagination(page int32, perPage int32) (int32, int32) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	return page, perPage
}

// ResolvePaging returns a slice of items up to the specified size and a boolean indicating if there is a next page.
// If size is zero or negative, all items are returned and next is false.
func ResolvePaging[C any](size int, items []C) (result []C, next bool) {
	if size > 0 && len(items) > size {
		return items[0:size], true
	}
	return items, false
}
