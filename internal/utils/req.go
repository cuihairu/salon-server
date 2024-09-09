package utils

type Paging struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

func NewPaging(page int, pageSize int) *Paging {
	return &Paging{
		Page:     page,
		PageSize: pageSize,
	}
}

func NewPagingWithDefault() *Paging {
	return NewPagingWithPage(1)
}

func NewPagingWithPage(page int) *Paging {
	return NewPaging(page, 10)
}

func (p *Paging) SetPage(page int) {
	p.Page = page
}

func (p *Paging) SetPageSize(pageSize int) {
	p.PageSize = pageSize
}

func (p *Paging) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Paging) Limit() int {
	return p.PageSize
}
