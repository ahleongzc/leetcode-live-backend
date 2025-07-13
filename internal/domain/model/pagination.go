package model

type Pagination struct {
	Offset  uint `json:"offset"`
	Limit   uint `json:"limit"`
	Total   uint `json:"total"`
	HasNext bool `json:"has_next"`
	HasPrev bool `json:"has_prev"`
}

func NewPagination() *Pagination {
	return &Pagination{}
}

func (p *Pagination) SetOffset(offset uint) *Pagination {
	if p == nil {
		return nil
	}
	p.Offset = offset
	return p
}

func (p *Pagination) SetLimit(limit uint) *Pagination {
	if p == nil {
		return nil
	}
	p.Limit = limit
	return p
}

func (p *Pagination) SetTotal(total uint) *Pagination {
	if p == nil {
		return nil
	}
	p.Total = total
	return p
}

func (p *Pagination) SetHasNext(hasNext bool) *Pagination {
	if p == nil {
		return nil
	}
	p.HasNext = hasNext
	return p
}

func (p *Pagination) SetHasPrev(hasPrev bool) *Pagination {
	if p == nil {
		return nil
	}
	p.HasPrev = hasPrev
	return p
}
