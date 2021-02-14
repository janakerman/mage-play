//+build mage

package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Version() error {
	fmt.Println("version")

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println("Working dir: ", dir)

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil
	}

	ref, err := repo.Head()
	if err != nil {
		return err
	}

	cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return err
	}

	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c.String())
		return nil
	})
	if err != nil {
		return nil
	}

	return nil
}
