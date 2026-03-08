package dtos

type PagedMeta struct {
	Total    int64 `json:"total"`
	Page     int32 `json:"page"`
	PageSize int32 `json:"pageSize"`
}

type PagedResponse[T any] struct {
	Data []T       `json:"data"`
	Meta PagedMeta `json:"meta"`
}
