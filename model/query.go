package model

type QueryReq struct {
	Label string `json:"label" form:"label"`
	Key   string `json:"key" form:"key"`
	Limit string `json:"limit" form:"limit"`
}
