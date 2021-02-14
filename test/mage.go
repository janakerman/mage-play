// +build ignore

// This file supports a no install run of Mage for use in tests.
package main

import (
	"os"

	"github.com/magefile/mage/mage"
)

func main() {
	os.Exit(mage.Main())
}
