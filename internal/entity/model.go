package entity

type Id struct {
	ID string `json:"id"`
}

type OrderBy struct {
	Column string `json:"column"`
	Order  string `json:"order"`
}

type Filter struct {
	Column string `json:"column"`
	Type   string `json:"type"` // eq, ne, gt, gte, lt, lte, search
	Value  string `json:"value"`
}

type GetListFilter struct {
	Page    int       `json:"offset"`
	Limit   int       `json:"limit"`
	Filters []Filter  `json:"filters"`
	OrderBy []OrderBy `json:"order_by"`
}

type UpdateFieldItem struct {
	Column string `json:"column"`
	Value  string `json:"value"`
}

type UpdateFieldRequest struct {
	Filter []Filter          `json:"filter"`
	Items  []UpdateFieldItem `json:"items"`
}

type RowsEffected struct {
	RowsEffected int `json:"rows_effected"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
