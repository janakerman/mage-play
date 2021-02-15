package version_test

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateTestRepo() (*git.Repository, error) {
	rootDir := os.TempDir() + ".test/"
	err := os.RemoveAll(rootDir)
	if err != nil {
		return nil, err
	}

	repoDir := rootDir + "repo/"
	os.MkdirAll(repoDir, 0777)
	if err != nil {
		return nil, err
	}

	// Cache is needed to avoid nil panic.
	return git.Init(filesystem.NewStorage(osfs.New(repoDir+".git"), cache.NewObjectLRUDefault()), osfs.New(repoDir))
}

func CommitWithTag(repo *git.Repository, msg, tag string) (*plumbing.Reference, error) {
	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	firstHash, err := w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}

	if tag == "" {
		return nil, nil
	}

	fmt.Printf("Creating tag %s on commit %s\n", tag, firstHash.String())
	// return repo.CreateTag(tag, firstHash, nil)
	return repo.CreateTag(tag, firstHash, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "name",
			Email: "email",
			When:  time.Now(),
		},
		Message: "asdf",
	})
}

func GoRunMage(repo *git.Repository, target string) error {
	workTree, err := repo.Worktree()
	Must(err)
	workingDir := workTree.Filesystem.Root()

	cmd := exec.Command("go", "run", "../test/mage.go", "-d", "../test", "-w", workingDir, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Test_VersionDoesNothingIfTagExists(t *testing.T) {
	repo, err := CreateTestRepo()
	Must(err)

	_, err = CommitWithTag(repo, "Merge 1", "1.0.0")
	fmt.Println(reflect.TypeOf(err))
	Must(err)
	_, err = CommitWithTag(repo, "Merge 2", "")
	Must(err)

	err = GoRunMage(repo, "version")
	Must(err)
}
