package model

var (
	OperatorEq   = "eq"
	OperatorNot  = "not"
	OperatorLike = "like"
)

type Filter struct {
	Page        int           `json:"page"`
	PageSize    int           `json:"pageSize"`
	FilterField []FilterField `json:"filterFields"`
}

type FilterField struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}
