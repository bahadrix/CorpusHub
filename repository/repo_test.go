package repository

import (
	"fmt"
	"github.com/bahadrix/corpushub/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

func TestMain(m *testing.M) {
	var err error
	workFolder, err = ioutil.TempDir("", "corpushub_test_*")

	if err != nil {
		panic(err)
	}

	//fmt.Printf("Test folder created at: %s\n", workFolder)
	defer func() {
		//fmt.Printf("Removing folder %s\n", workFolder)
		err = os.RemoveAll(workFolder)

		if err != nil {
			fmt.Printf("Error on removing folder: %s\n", workFolder)
		}

	}()


	m.Run()

}

func TestClone(t *testing.T) {

	for _, branchToTest := range []string{"master", "test"}{
		t.Run(branchToTest, func(t *testing.T) {
			var mockRepoOptions = &model.RepoOptions{
				Title:      "Mock Repo",
				URL:        "git@github.com:bahadrix/git-mock-repo.git",
				Branch:     branchToTest,
				PrivateKey: repoDeployKey,
			}

			repoPath := fmt.Sprintf("%s/git-mock-repo", workFolder)

			defer os.RemoveAll(workFolder)

			// Test Cloning
			repo, err := NewRepo(repoPath, mockRepoOptions)

			if err != nil {
				t.Error("Error on initializing repo", err)
			}

			alreadyUpToDate, err := repo.Sync()

			if err != nil {
				t.Error("Error on syncing repo", err)
			}

			assert.False(t, alreadyUpToDate)

			// Test re-sync
			alreadyUpToDate, err = repo.Sync()

			if err != nil {
				t.Error("Error on syncing repo", err)
			}

			assert.True(t, alreadyUpToDate)

			// Read file directly and check branch is correct
			branchBytes, err := ioutil.ReadFile(filepath.Join(repoPath, "branch"))

			if err != nil {
				t.Error("Can't open branchToTest file", err)
			} else {
				branchName := strings.ToLower(strings.TrimSpace(string(branchBytes)))
				assert.Equal(t, branchToTest, branchName, "Wrong branch")
			}

			// Read file via repo interface with fileURI and check branch is correct
			branchBytes, err = repo.ReadFile(fmt.Sprintf("repo://github.com:bahadrix/git-mock-repo/%s/branch", branchToTest))

			if err != nil {
				t.Error("Error on syncing repo", err)
			} else {
				branchName := strings.ToLower(strings.TrimSpace(string(branchBytes)))
				assert.Equal(t, branchToTest, branchName, "Wrong branch")
			}

		})
	}




}
