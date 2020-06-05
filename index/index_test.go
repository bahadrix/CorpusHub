package index

import (
	"fmt"
	"github.com/bahadrix/corpushub/repository"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var repoDeployKey = []byte(`-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACD+GJTtajodrPNZmTq10ZR4VcoFxvOaNc00RJf2SiqRqAAAAKCXRchil0XI
YgAAAAtzc2gtZWQyNTUxOQAAACD+GJTtajodrPNZmTq10ZR4VcoFxvOaNc00RJf2SiqRqA
AAAED5KICQxVeoWmSI7we4WYArFyjfIKa57+xq+p31EI95n/4YlO1qOh2s81mZOrXRlHhV
ygXG85o1zTREl/ZKKpGoAAAAG2JhaGFkaXJAQmFoYWRpcnMtaUJhZy5sb2NhbAEC
-----END OPENSSH PRIVATE KEY-----
`)

var workFolder string
var repo *repository.Repo
var indexPath string
var mdCount = 0

func TestMain(m *testing.M) {

	var err error
	workFolder, err = ioutil.TempDir("", "corpushub_test_*")

	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(workFolder)

	indexPath = fmt.Sprintf("%s/testindex", workFolder)

	repoPath := fmt.Sprintf("%s/git-mock-repo", workFolder)

	var mockRepoOptions = &repository.RepoOptions{
		Title:      "Mock Repo",
		URL:        "git@github.com:bahadrix/git-mock-repo.git",
		Branch:     "master",
		PrivateKey: repoDeployKey,
	}

	//defer os.RemoveAll(workFolder)

	repo, err = repository.NewRepo(repoPath, mockRepoOptions)

	if err != nil {
		log.Fatal("Error on initializing repo", err)
	}

	_, err = repo.Sync()

	if err != nil {
		log.Fatal("Error on synchronizing repo", err)
	}

	readmeBytes, err := repo.ReadFile("/README.md")

	if err != nil {
		log.Fatal("Error on reading repo", err)
	}

	_ = readmeBytes

	m.Run()

}

func TestCreateIndex(t *testing.T) {

	index, err := NewIndex(indexPath)

	if err != nil {
		t.Error("Can't create index ", err)
	}

	infos, err := repo.GetFileInfos("/", true)

	if err != nil {
		t.Error("Can't get file infos", err)
	}

	for uri, info := range infos {

		extension := strings.ToLower(filepath.Ext(info.Name()))

		if extension != ".md" {
			continue
		}

		mdBytes, err := repo.ReadFile(uri)

		if err != nil {
			t.Error("Can not read file ", err)
		}

		doc := &Document{
			URI:          uri,
			Title:        info.Name(),
			Content:      string(mdBytes),
			LastModified: info.ModTime(),
		}

		err = index.IndexDocument(doc)

		if err != nil {
			t.Error("Can not index document", err)
		}

		mdCount++

	}

	err = index.Close()

	if err != nil {
		t.Log("Can not cloes the index ", err)
	}
}

func TestSearch(t *testing.T) {

	index, err := NewIndex(indexPath)

	if err != nil {
		t.Error("Can't load index ", err)
	}

	result, err := index.SearchContent("Lorem")

	if err != nil {
		t.Error("Can't do search ", err)
	}

	assert.Greater(t, result.TotalHits, uint64(0))

	err = index.Close()

	if err != nil {
		t.Log("Can not cloes the index ", err)
	}
}

func TestGetURIs(t *testing.T) {
	index, err := NewIndex(indexPath)

	repoIndices, err := index.GetURIsByPrefix(repo.GetURI())
	//repoIndices, err := index.GetURIsByPrefix("repo://")

	if err != nil {
		t.Log("Can not get uris by prefix ", err)
	}

	assert.Equal(t, mdCount, len(repoIndices))

	err = index.Close()

	// repo://github.com:bahadrix/git-mock-repo/master/testdocs/sub1/README.md
	// repo://github.com:bahadrix/git-mock-repo/master
	if err != nil {
		t.Log("Can not cloes the index ", err)
	}

}
