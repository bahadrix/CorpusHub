package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptions(t *testing.T) {

	opts := &RepoOptions{
		Title:      "Mock Repo2",
		URL:        "git@github.com:bahadrix/git-mock-repo.git",
		Branch:     "development",
		PrivateKey: []byte("le key"),
	}

	assert.Equal(t, "github.com:bahadrix/git-mock-repo", opts.GetNormalizedURI())

	data, err := opts.Serialize()
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, data)

	opts2, err := DeserializeRepoOptions(data)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, opts, opts2)




}