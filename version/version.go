package version

import (
	"fmt"
	"sort"

	"github.com/blang/semver/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func NewVersion(repo *git.Repository) error {
	tags, err := latestAncestorTags(repo)
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		return fmt.Errorf("no tags found on branch")
	}

	var versions []semver.Version
	for _, t := range tags {
		v, err := semver.Make(t.Name)
		if err != nil {
			return fmt.Errorf("failed parsing version: %w", err)
		}
		versions = append(versions, v)
	}

	sort.SliceStable(versions, func(i, j int) bool {
		return versions[i].Compare(versions[j]) > 0
	})

	latest := versions[0]
	fmt.Printf("Latest tag: %s", latest.String())

	return nil
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
