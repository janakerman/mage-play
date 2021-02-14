package version_test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

	return git.Init(filesystem.NewStorage(osfs.New(repoDir+".git"), nil), osfs.New(repoDir))
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

	return repo.CreateTag(tag, firstHash, nil)
}

func RunMage(repo *git.Repository, target string) error {
	workTree, err := repo.Worktree()
	Must(err)
	workingDir := workTree.Filesystem.Root()

	cmd := exec.Command("go", "run", "./test/mage.go", "-w", workingDir, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Test_Version(t *testing.T) {
	repo, err := CreateTestRepo()
	Must(err)

	CommitWithTag(repo, "Merge 1", "v1.0.0")
	CommitWithTag(repo, "Merge 2", "")

	RunMage(repo, "version")
}
