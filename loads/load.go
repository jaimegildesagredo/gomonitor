package loads

import "time"

type Load struct {
	Values    LoadValues
	CreatedAt time.Time
}

type LoadValues []float32
