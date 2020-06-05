package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bahadrix/corpushub/repository"
	"github.com/bahadrix/corpushub/server/operator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testHost = "1.1.1.1"
	testPort = 80
)

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

func doRequest(router *gin.Engine, request *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, request)
	return rec
}

func TestRest(t *testing.T) {

	workFolder, err := ioutil.TempDir("", "corpushub_test_*")
	if err != nil {
		panic(err)
	}

	op, err := operator.NewOperator(workFolder, nil)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.TestMode)
	router := setupRouter(testHost, testPort, op)

	// Test ping
	req, _ := http.NewRequest("GET", "/v1/ping", nil)
	resp := doRequest(router, req)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, fmt.Sprintf("PONGv%s %s %d", VERSION, testHost, testPort), resp.Body.String())

	mockRepoData, _ := json.Marshal(mockRepoOptions)

	// Test add repo
	req, _ = http.NewRequest("PUT", "/v1/repos", bytes.NewReader(mockRepoData))
	resp = doRequest(router, req)

	assert.Equal(t, 200, resp.Code)

	// Test manual repo sync
	req, _ = http.NewRequest("GET", fmt.Sprintf("/v1/repo/sync?uri=%s", mockRepoOptions.GetNormalizedURI()), nil)
	resp = doRequest(router, req)

	assert.Equal(t, 200, resp.Code, resp.Body.String())

}
