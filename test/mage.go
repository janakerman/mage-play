// +build ignore

// This file provides a no install way to run the Mage tool for use in this packages tests.
package main

import (
	"os"

	"github.com/magefile/mage/mage"
)

func main() {
	os.Exit(mage.Main())
}
