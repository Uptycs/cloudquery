package compute

import (
	"os"
	"testing"

	"github.com/Uptycs/cloudquery/utilities"
)

func TestMain(m *testing.M) {
	utilities.ReadTableConfigurations("../../")
	os.Exit(m.Run())
}
