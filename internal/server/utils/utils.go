package utils

import (
	"net/url"
	"strings"
)

func ParseURL(url *url.URL) []string {
	spath := strings.Split(url.Path, "/")
	return spath
}
