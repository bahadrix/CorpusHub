package main

import (
	"fmt"
	"github.com/bahadrix/corpushub/repository"

	"github.com/bahadrix/corpushub/model"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
	"log"
	"os"
	"time"
)

var repoDeployKey = []byte(`-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACD+GJTtajodrPNZmTq10ZR4VcoFxvOaNc00RJf2SiqRqAAAAKCXRchil0XI
YgAAAAtzc2gtZWQyNTUxOQAAACD+GJTtajodrPNZmTq10ZR4VcoFxvOaNc00RJf2SiqRqA
AAAED5KICQxVeoWmSI7we4WYArFyjfIKa57+xq+p31EI95n/4YlO1qOh2s81mZOrXRlHhV
ygXG85o1zTREl/ZKKpGoAAAAG2JhaGFkaXJAQmFoYWRpcnMtaUJhZy5sb2NhbAEC
-----END OPENSSH PRIVATE KEY-----
`)

type Document struct {
	URI     string `json:"uri"`
	Content string `json:"content"`
	Date time.Time `json:"date"`
}

func (d *Document) Type() string {
	return "Document"
}

func main() {
	var mockRepoOptions = &model.RepoOptions{
		Title:      "Mock Repo",
		URL:        "git@github.com:bahadrix/git-mock-repo.git",
		Branch:     "master",
		PrivateKey: repoDeployKey,
	}

	repoPath := fmt.Sprintf("tmp/git-mock-repo")

	//defer os.RemoveAll(workFolder)

	repo, err := repository.NewRepo(repoPath, mockRepoOptions)

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

	doc := &Document{
		URI:     repo.GenerateFullPath("/README.md"),
		Content: fmt.Sprintf("%s", readmeBytes),
		Date: time.Now(),
	}

	// Remove previous index storage

	os.RemoveAll("tmp/example.bleve")


	// Create index
	docMapping := bleve.NewDocumentMapping()


	uriFieldMapping := bleve.NewTextFieldMapping()
	uriFieldMapping.Analyzer = keyword.Name

	docMapping.AddFieldMappingsAt("uri", uriFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("Document", docMapping)
	index, err := bleve.New("tmp/example.bleve", indexMapping)

	if err != nil {
		log.Fatal("Error on creating new index", err)
	}

	err = index.Index(doc.URI, doc)

	if err != nil {
		log.Fatal("Error on indexing", err)
	}

	// Query document


	//query := bleve.NewQueryStringQuery("repo")
	query := bleve.NewPrefixQuery("repo://")
	query.SetField("uri")
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"uri", "date"}
	search.Highlight = bleve.NewHighlightWithStyle(ansi.Name)
	results, err := index.Search(search)

	if err != nil {
		log.Fatal("Error on searching", err)
	}

	fmt.Println(results)


	for _, hit := range results.Hits {

		fmt.Printf("%v", hit.Fragments)
	}

	//fmt.Printf("%v", results.Hits[0].Fragments["Content"])



}
