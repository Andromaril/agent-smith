package errormetric

import (
	"fmt"
	"time"
)

type MetricError struct {
	Time time.Time
	Err  error
}

func (te *MetricError) Error() string {
	return fmt.Sprintf("%v %v", te.Time.Format("2006/01/02 15:04:05"), te.Err)
}

func NewMetricError(err error) error {
	return &MetricError{
		Time: time.Now(),
		Err:  err,
	}
}
