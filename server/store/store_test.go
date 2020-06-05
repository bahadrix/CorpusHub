package store

import (
	"github.com/bahadrix/corpushub/repository"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"testing"
)

var repoStore *RepoStore

func TestMain(m *testing.M) {

	workFolder, err := ioutil.TempDir("", "corpushub_test_*")

	if err != nil {
		panic(err)
	}

	repoStore, err = NewRepoStore(path.Join(workFolder, "repostore.db"))

	if err != nil {
		panic(err)
	}

	m.Run()
}


func TestOptions(t *testing.T) {

	// Add sample repo
	err := repoStore.PutRepo(&repository.RepoOptions{
		Title:      "Mock Repo",
		URL:        "git@github.com:bahadrix/git-mock-repo.git",
		Branch:     "master",
		PrivateKey: []byte("le key"),
	}, nil)

	if err != nil {
		t.Error(err)
	}

	// Get all URIs
	uris, err := repoStore.FindAll()

	if err != nil {
		t.Error(err)
	}

	assert.ElementsMatch(t, uris, []string{"github.com:bahadrix/git-mock-repo"})

	// Same repo with different details
	testOpts2 := &repository.RepoOptions{
		Title:      "Mock Repo2",
		URL:        "git@github.com:bahadrix/git-mock-repo.git",
		Branch:     "development",
		PrivateKey: []byte("le key"),
	}

	// Add this also to store
	err = repoStore.PutRepo(testOpts2, nil)
	if err != nil {
		t.Error(err)
	}

	// Get all URIs
	uris, err = repoStore.FindAll()
	if err != nil {
		t.Error(err)
	}

	// We must have still only one repo
	assert.ElementsMatch(t, uris, []string{"github.com:bahadrix/git-mock-repo"})

	// with same details
	opts, err := repoStore.GetRepoOptions("github.com:bahadrix/git-mock-repo")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, testOpts2, opts)

	// Add another repo
	err = repoStore.PutRepo(&repository.RepoOptions{
		Title:      "Mock Repo",
		URL:        "git@github.com:bahadrix/git-mock-repo2.git",
		Branch:     "master",
		PrivateKey: []byte("le key"),
	}, nil)

	if err != nil {
		t.Error(err)
	}

	// Get all URIs
	uris, err = repoStore.FindAll()
	if err != nil {
		t.Error(err)
	}

	assert.ElementsMatch(t, uris, []string{"github.com:bahadrix/git-mock-repo", "github.com:bahadrix/git-mock-repo2"})

	// Delete second repo
	err = repoStore.DeleteRepo("github.com:bahadrix/git-mock-repo2")
	if err != nil {
		t.Error(err)
	}

	// Get all URIs
	uris, err = repoStore.FindAll()
	if err != nil {
		t.Error(err)
	}

	assert.ElementsMatch(t, uris, []string{"github.com:bahadrix/git-mock-repo"})

	assert.False(t, repoStore.Exists("github.com:bahadrix/git-mock-repo2"))
	assert.True(t, repoStore.Exists("github.com:bahadrix/git-mock-repo"))

	// Meta testing

	// Try add meta to non existing repo
	err = repoStore.PutMeta("github.com:bahadrix/git-mock-repo2", []byte("Test meta key"), []byte("Test meta value"))
	assert.Equal(t, ErrRepoNotFound, err)

	// Add meta to existed repo
	err = repoStore.PutMeta("github.com:bahadrix/git-mock-repo", []byte("Test meta key"), []byte("Test meta value"))
	if err != nil {
		t.Error(err)
	}

	// Try get meta from non existing repo
	value, err := repoStore.GetMeta("github.com:bahadrix/git-mock-repo2", []byte("Test meta key"))
	assert.Equal(t, ErrRepoNotFound, err)
	assert.Nil(t, value)


	// Get meta
	value, err = repoStore.GetMeta("github.com:bahadrix/git-mock-repo", []byte("Test meta key"))
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t,  []byte("Test meta value"), value)

	// Delete meta
	err = repoStore.DeleteMeta("github.com:bahadrix/git-mock-repo", []byte("Test meta key"))
	if err != nil {
		t.Error(err)
	}

	// Try to get non existing meta from existing repo
	value, err = repoStore.GetMeta("github.com:bahadrix/git-mock-repo", []byte("Test meta key"))
	assert.Nil(t, value)
	assert.Nil(t, err)

	// Add metamap at repo creation
	metaMap := map[string][]byte{
		"maKey": []byte("Ma Value"),
		"maKey2": []byte("Ma Value2"),
	}
	err = repoStore.PutRepo(opts, metaMap)
	if err != nil {
		t.Error(err)
	}

	for k, v := range metaMap {
		value, err = repoStore.GetMeta("github.com:bahadrix/git-mock-repo", []byte(k))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t,  v, value)
	}




}
