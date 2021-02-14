package version

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/janakerman/mage-play/version"
)

// Version tags the current commit with the next version number based on the commit message.
func Version() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil
	}

	return version.NewVersion(repo)
}
