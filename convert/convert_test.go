package convert

import (
	"flag"
	"os"
	"testing"
)

var updateGolden = flag.Bool("updateGolden", false, "Set to true to update the golden files")

func TestMain(m *testing.M) {
	flag.Parse()
	exitVal := m.Run()
	// do any teardown here if necessary
	os.Exit(exitVal)
}
