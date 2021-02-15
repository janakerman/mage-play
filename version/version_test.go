package version_test

import (
	"context"
	"github.com/janakerman/mage-play/version"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

func Must(t *testing.T, err error) {
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
}

func CreateTestRepo() (*git.Repository, error) {
	rootDir := os.TempDir() + ".test/"
	err := os.RemoveAll(rootDir)
	if err != nil {
		return nil,  err
	}

	repoDir := rootDir + "remote/"
	err = os.MkdirAll(repoDir, 0777)
	if err != nil {
		return nil, err
	}

	// Cache is needed to avoid nil panic.
	repo, err := git.Init(filesystem.NewStorage(osfs.New(repoDir+".git"), cache.NewObjectLRUDefault()), osfs.New(repoDir))
	if err != nil {
		return  nil, err
	}
	return repo, nil
}

func CloneTestRepo() (*git.Repository, error) {
	rootDir := os.TempDir() + ".test/"

	repoDir := rootDir + "repo/"
	err := os.MkdirAll(repoDir, 0777)
	if err != nil {
		return nil, err
	}

	// Cache is needed to avoid nil panic.
	repo, err := git.Clone(filesystem.NewStorage(osfs.New(repoDir+".git"), cache.NewObjectLRUDefault()), osfs.New(repoDir), &git.CloneOptions{
		URL:               rootDir + "remote/",
		RemoteName:        "origin",
	})
	if err != nil {
		return  nil, err
	}
	return repo, nil
}

func CommitWithTag(repo *git.Repository, msg, tag string) (*plumbing.Hash, error) {
	w, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	commitHash, err := w.Commit(msg, &git.CommitOptions{
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
		return &commitHash, nil
	}

	_, err = repo.CreateTag(tag, commitHash, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "name",
			Email: "email",
			When:  time.Now(),
		},
		Message: "asdf",
	})
	if err != nil {
		return nil, err
	}

	return &commitHash, nil
}

func GoRunMage(repo *git.Repository, target string) error {
	workTree, err := repo.Worktree()
	if err != nil {
		return err
	}
	workingDir := workTree.Filesystem.Root()

	cmd := exec.Command("go", "run", "../test/mage.go", "-d", "../test", "-w", workingDir, target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Test_BumpVersionViaMage(t *testing.T) {
	remote, err := CreateTestRepo()
	Must(t, err)
	_, err = CommitWithTag(remote, "Merge 1", "1.0.0")
	Must(t, err)
	latest, err := CommitWithTag(remote, "Merge 2", "")
	Must(t, err)
	repo, err := CloneTestRepo()
	Must(t, err)


	err = GoRunMage(repo, "bumpVersion")
	Must(t, err)


	ref, err := remote.Tag("1.0.1")
	Must(t, err)
	r, err := remote.TagObject(ref.Hash())
	Must(t, err)

	if r.Target != *latest {
		t.Error("expected tag not present")
	}
}

func Test_BumpVersion(t *testing.T) {
	remote, err := CreateTestRepo()
	Must(t, err)
	_, err = CommitWithTag(remote, "Merge 1", "1.0.0")
	Must(t, err)
	latest, err := CommitWithTag(remote, "Merge 2", "")
	Must(t, err)
	repo, err := CloneTestRepo()
	Must(t, err)


	err = version.BumpVersion(context.TODO(), repo)
	Must(t, err)


	ref, err := remote.Tag("1.0.1")
	Must(t, err)
	r, err := remote.TagObject(ref.Hash())
	Must(t, err)

	if r.Target != *latest {
		t.Error("expected tag not present")
	}
}
