//Package errormetric служит для обработки ошибок

package errormetric

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetricError(t *testing.T) {
	err1 := errors.New("test error")
	err := NewMetricError(err1)
	assert.Error(t, err)
}

func TestMetricError_Error(t *testing.T) {
	err1 := errors.New("test error")
	err := NewMetricError(err1)
	ex := err.Error()
	exerr := fmt.Sprintf("%v %v", time.Now().Format("2006/01/02 15:04:05"), err1)
	require.Equal(t, exerr, ex)
}
