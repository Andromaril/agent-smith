// Package retry служит для возможности ретрая операций
package retry

import (
	"fmt"
	"log"
	"time"

	"github.com/andromaril/agent-smith/internal/errormetric"
)

// Retry выполняет ретрай определенной функции 3 раза, если после 3 ретраев заданная функция не прошла - выводит ошибку
func Retry(function func() error) error {

	tries := 0
	wait := 1
	var err error
	for {
		if tries >= 3 {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("fatal  %w", e)
		}
		err = function()
		if err != nil {
			e := errormetric.NewMetricError(err)
			log.Printf("fatal start operation %q", e.Error())
			time.Sleep(time.Duration(wait) * time.Second)
			wait += 2
			tries++
			continue
		}
		break
	}
	return nil
}
