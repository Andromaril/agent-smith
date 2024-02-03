package utils

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestParseURL(t *testing.T) {
	r1, _:= url.Parse("/update/gauge/test/1")
	r2, _:= url.Parse("/update/gauge/test/2")
	r3, _:= url.Parse("/update/gauge/")
	r4, _:= url.Parse("/update/gauge/test/rt")

	tests := []struct {
		name string
		have *url.URL
		want []string
	}{
		{name: "test1", have: r1, want: []string{"", "update", "gauge", "test", "1"}},
		{name: "test_2", have: r2, want: []string{"", "update", "gauge", "test", "2"}},
		{name: "test_3", have: r3, want: []string{"", "update", "gauge", ""}},
		{name: "test_4", have: r4, want: []string{"", "update", "gauge", "test", "rt"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ParseURL(tt.have))
		})
	}
}
