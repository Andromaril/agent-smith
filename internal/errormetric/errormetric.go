//Package errormetric служит для обработки ошибок
package errormetric

import (
	"fmt"
	"time"
)
// MetricError структура для обработки ошибок
type MetricError struct {
	Time time.Time
	Err  error
}

// Error для вывода ошибки в string-формате
func (te *MetricError) Error() string {
	return fmt.Sprintf("%v %v", te.Time.Format("2006/01/02 15:04:05"), te.Err)
}

// NewMetricError создает новую структуру MetricError
func NewMetricError(err error) error {
	return &MetricError{
		Time: time.Now(),
		Err:  err,
	}
}
