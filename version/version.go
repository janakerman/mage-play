package version

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5/config"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func BumpVersion(ctx context.Context, repo *git.Repository) error {
	tags, err := latestAncestorTags(repo)
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		return fmt.Errorf("no tags found on branch")
	}

	_, version, err := latestTag(tags)
	if err != nil {
		return fmt.Errorf("error parsing tags: %w", err)
	}

	err = bumpHead(repo, version)
	if err != nil {
		return fmt.Errorf("failed bumping version: %w", err)
	}

	err = repo.PushContext(ctx, &git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"refs/tags/*:refs/tags/*"},
		Progress:   os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("failed pushing tags: %w", err)
	}

	return nil
}

func bumpHead(repo *git.Repository, latestVersion semver.Version) error {
	ref, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to fetch head: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return fmt.Errorf("failed to get commit object: %w", err)
	}

	if strings.Contains(commit.Message, "#bump-major") {
		latestVersion.Major++
	} else if strings.Contains(commit.Message, "#bump-minor") {
		latestVersion.Minor++
	} else {
		latestVersion.Patch++
	}

	ref, err = repo.CreateTag(latestVersion.String(), commit.Hash, &git.CreateTagOptions{
		Message: fmt.Sprintf("Release %s", latestVersion.String()),
		Tagger: &object.Signature{ // TODO: Take from somewhere else.
			Name:  "name",
			Email: "email",
			When:  time.Time{},
		},
	})
	if err != nil {
		return fmt.Errorf("failed creating tag: %w", err)
	}
	return nil
}

func latestTag(tags []*object.Tag) (*object.Tag, semver.Version, error) {
	type pair struct {
		semver semver.Version
		tag *object.Tag
	}

	var versions []pair
	for _, t := range tags {
		v, err := semver.Make(t.Name)
		if err != nil {
			return nil, semver.Version{}, fmt.Errorf("failed parsing version: %w", err)
		}
		versions = append(versions, pair{
			semver: v,
			tag:    t,
		})
	}

	sort.SliceStable(versions, func(i, j int) bool {
		return versions[i].semver.Compare(versions[j].semver) > 0
	})

	return versions[0].tag, versions[0].semver, nil
}

func latestAncestorTags(repo *git.Repository) ([]*object.Tag, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	hashToTags, err := hashToTags(repo)
	if err != nil {
		return nil, err
	}

	cIter, err := repo.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderDFS,
	})
	if err != nil {
		return nil, err
	}

	var tags []*object.Tag
	err = cIter.ForEach(func(c *object.Commit) error {
		commitTags, ok := hashToTags[c.Hash]
		if !ok || len(commitTags) == 0 {
			return nil
		}

		if tags == nil { // We only care about the latest tags aka first iteration.
			tags = commitTags
		}
		return nil // TODO: Return storer.ErrStop to break loop.
	})
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// hashToTags returns a map of target hash (commit hash) to tag.
func hashToTags(repo *git.Repository) (map[plumbing.Hash][]*object.Tag, error) {
	tags, err := repo.TagObjects()
	if err != nil {
		return nil, err
	}

	hashToTag := map[plumbing.Hash][]*object.Tag{}
	err = tags.ForEach(func(t *object.Tag) error {
		hashToTag[t.Target] = append(hashToTag[t.Target], t)
		return nil
	})
	return hashToTag, nil
}
