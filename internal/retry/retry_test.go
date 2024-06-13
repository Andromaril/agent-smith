// Package retry служит для возможности ретрая операций

package retry

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mulfunc(i int) error {
	l := i * 2
	fmt.Print(l)
	return nil
}

func TestRetry(t *testing.T) {
	operation := func() error {
		err := mulfunc(2)
		return err
	}
	err := Retry(operation)
	assert.NoError(t, err)

}
