package dto

import "math"

type PageResponse[T any] struct {
	Content       []T     `json:"content"`
	TotalElements int     `json:"totalElements"`
	TotalPages    int     `json:"totalPages"`
	Number        int     `json:"number"`
	Size          int     `json:"size,omitempty"`
	Count         []Count `json:"count,omitempty"`
}

type Count struct {
	TotalElement int `json:"totalElement"`
}

func (p *PageResponse[any]) SetPage(size int, page int) {
	if len(p.Content) > 0 {
		p.TotalElements = p.Count[0].TotalElement
		p.TotalPages = int(math.Ceil(float64(p.TotalElements) / float64(size)))
		p.Number = page
	}
	p.Size = size
}
