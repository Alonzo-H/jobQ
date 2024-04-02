package job

import (
	"fmt"
	"slices"
	"time"
)

type Job struct {
	Id          uint64    `json:"id"`
	Type        ttype     `json:"type"`
	Status      status    `json:"status"`
	ConcludedAt time.Time `json:"-"`
}

type ttype string

const (
	TimeCritical    ttype = "TIME_CRITICAL"
	NotTimeCritical ttype = "NOT_TIME_CRITICAL"
)

var validTypes = []ttype{
	TimeCritical,
	NotTimeCritical,
}

func (t ttype) IsValid() error {
	if !slices.Contains(validTypes, t) {
		return fmt.Errorf("not a valid type")
	}
	return nil
}

type status string

const (
	Queued     = "QUEUED"
	InProgress = "IN_PROGRESS"
	Concluded  = "CONCLUDED"
)
