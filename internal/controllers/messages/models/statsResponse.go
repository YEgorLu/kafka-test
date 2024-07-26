package models

type StatsResponse struct {
	CountAll       int64 `json:"countAll"`
	CountProcessed int64 `json:"countProcessed"`
}
