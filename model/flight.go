package model

import "time"

type FilghtData struct {
	FlightNo  string    `json:"flight_no"`
	FromCode  string    `json:"from_code"`
	ToCode    string    `json:"to_code"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type FlightModel struct{}
