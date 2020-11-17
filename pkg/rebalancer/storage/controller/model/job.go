package model

import "time"

// Job -
type Job struct {
	Name      string
	Namespace string
	Active    int32
	StartTime time.Time
}
