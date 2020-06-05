package operator

import (
	"github.com/bahadrix/corpushub/repository"
	"io/ioutil"
	"os"
	"testing"
)

var operator *Operator

var repoDeployKey = []byte(`-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACD+GJTtajodrPNZmTq10ZR4VcoFxvOaNc00RJf2SiqRqAAAAKCXRchil0XI
YgAAAAtzc2gtZWQyNTUxOQAAACD+GJTtajodrPNZmTq10ZR4VcoFxvOaNc00RJf2SiqRqA
AAAED5KICQxVeoWmSI7we4WYArFyjfIKa57+xq+p31EI95n/4YlO1qOh2s81mZOrXRlHhV
ygXG85o1zTREl/ZKKpGoAAAAG2JhaGFkaXJAQmFoYWRpcnMtaUJhZy5sb2NhbAEC
-----END OPENSSH PRIVATE KEY-----
`)

var mockRepoOptions = &repository.RepoOptions{
	Title:      "Mock Repo",
	URL:        "git@github.com:bahadrix/git-mock-repo.git",
	Branch:     "master",
	PrivateKey: repoDeployKey,
}

func TestMain(m *testing.M) {

	workFolder, err := ioutil.TempDir("", "corpushub_test_*")
	if err != nil {
		panic(err)
	}

	operator, err = NewOperator(workFolder, nil)
	if err != nil {
		panic(err)
	}

	ec := m.Run()
	os.Exit(ec)
}

func TestOperator(t *testing.T) {

	err := operator.AddRepo(mockRepoOptions)
	if err != nil {
		t.Error(err)
	}

	err = operator.SyncRepo(mockRepoOptions.GetNormalizedURI())
	if err != nil {
		t.Error(err)
	}

	results, err := operator.Search("Lorem")

	if err != nil {
		t.Error(err)
	}

	_ = results

}