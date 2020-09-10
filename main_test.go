package redis

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initMetrics()
	os.Exit(m.Run())
}
