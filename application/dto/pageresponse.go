package dto

import "math"

type PageResponse struct {
	Content      []any   `json:"content"`
	TotalElement int     `json:"totalElement"`
	TotalPages   int     `json:"totalPages"`
	Number       int     `json:"number"`
	Size         int     `json:"size,omitempty"`
	Count        []Count `json:"count,omitempty"`
}

type Count struct {
	TotalElement int `json:"totalElement,omitempty"`
}

func (p *PageResponse) SetPage(size int, page int) {
	if len(p.Content) > 0 {
		p.TotalElement = p.Count[0].TotalElement
		p.TotalPages = int(math.Ceil(float64(p.TotalElement) / float64(size)))
		p.Number = page
	}
	p.Size = size
}
