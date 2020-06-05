package index

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/search/highlight/highlighter/html"
	"os"
)

type Index struct {
	// Bleve Index
	bindex bleve.Index
}

func createBIndex(indexPath string) (bleve.Index, error) {

	// URI field mapping
	mURI := bleve.NewTextFieldMapping()
	mURI.Analyzer = keyword.Name

	// Document mapping
	mDoc := bleve.NewDocumentMapping()
	mDoc.AddFieldMappingsAt("uri", mURI)

	// Index mapping
	mIndex := bleve.NewIndexMapping()
	mIndex.AddDocumentMapping("document", mDoc)

	return bleve.New(indexPath, mIndex)
}

// bsr2SearchResult converts Bleve searchresult to our Search Result for the sake of decoupling
func bsr2SearchResult(bsr *bleve.SearchResult) *SearchResult {

	searchResult := &SearchResult{
		TotalHits: bsr.Total,
		MaxScore:  bsr.MaxScore,
		Took:      bsr.Took,
		Hits:      make([]*SearchHit, bsr.Hits.Len()),
	}

	for i, hit := range bsr.Hits {
		searchResult.Hits[i] = &SearchHit{
			ID:        hit.ID,
			Score:     hit.Score,
			Fragments: hit.Fragments,
			Fields:    hit.Fields,
		}
	}

	return searchResult
}

// NewIndex returns new Index object. If index folder already exists at indexPath
// loads it instead of creating new one.
func NewIndex(indexPath string) (*Index, error) {

	_, err := os.Stat(indexPath)
	pathExists := !os.IsNotExist(err)

	var bindex bleve.Index

	if pathExists {
		bindex, err = bleve.Open(indexPath)

	} else {
		bindex, err = createBIndex(indexPath)
	}

	if err != nil {
		return nil, err
	}

	return &Index{bindex: bindex}, nil

}

func (index *Index) IndexDocument(doc *Document) error {
	err := doc.Validate()

	if err != nil {
		return err
	}

	return index.bindex.Index(doc.URI, doc)
}

func (index *Index) SearchContent(stringQuery string) (*SearchResult, error) {
	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(stringQuery))
	searchRequest.Fields = []string{"uri", "title", "last_modified"}
	searchRequest.Highlight = bleve.NewHighlightWithStyle(html.Name)
	results, err := index.bindex.Search(searchRequest)

	if err != nil {
		return nil, err
	}

	return bsr2SearchResult(results), nil
}

func (index *Index) GetURIsByPrefix(prefix string) ([]string, error) {
	query := bleve.NewPrefixQuery(prefix)
	query.SetField("uri")

	searchRequest := bleve.NewSearchRequest(query)
	results, err := index.bindex.Search(searchRequest)

	if err != nil {
		return nil, err
	}

	ids := make([]string, results.Hits.Len())

	for i, hit := range results.Hits {
		ids[i] = hit.ID
	}

	return ids, nil
}

func (index *Index) Close() error {
	return index.bindex.Close()
}

func (index *Index) DeleteDocument(docURI string) error {
	return index.bindex.Delete(docURI)
}