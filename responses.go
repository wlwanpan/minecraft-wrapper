package wrapper

import "time"

type DataGetResponse struct {
	HurtByTimestamp int       `json:"HurtByTimestamp"`
	SleepTimer      time.Time `json:"SleepTimer"`
	Pos             []float64 `json:"Pos"`
}

func strToDataGet(raw string) (*DataGetResponse, error) {
	return &DataGetResponse{}, nil
}

func logJsonParser() ([]byte, error) {

}
