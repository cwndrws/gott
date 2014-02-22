package gott

import (
	"testing"
)

func TestForReality(t *testing.T) {
	if true == false {
		t.Fatal("RUNNNNNNNN!!!!!")
	}
}
