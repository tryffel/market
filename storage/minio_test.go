package storage

import (
	"fmt"
	"testing"
)

func TestMinioFullPath(t *testing.T) {
	path := "5f7f981a-0c73-49f2-8f79-99282dea836d"
	full := fullPath(path)

	if full != fmt.Sprintf("5f/7f/%s", path) {
		t.Error("Got invalid full path for file")
	}
}
