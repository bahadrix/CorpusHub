package operator

import (
	"crypto/sha1"
	"fmt"
	"github.com/bahadrix/corpushub/index"
	"github.com/bahadrix/corpushub/repository"
	"github.com/bahadrix/corpushub/server/store"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Operator struct {
	store     *store.RepoStore
	index     *index.Index
	options   *Options
	reposPath string
	rc *RepoCache
}

type Options struct {

}

const ( // Meta Keys
	mkURIHash = "uriHash"
	mkLastUpdate = "lastUpdate"
)


func NewOperator(dataPath string, options *Options) (*Operator, error) {

	storePath := path.Join(dataPath, "store.db")
	indexPath := path.Join(dataPath, "index")

	s, err := store.NewRepoStore(storePath)

	if err != nil {
		return nil, err
	}

	idx, err := index.NewIndex(indexPath)

	if err != nil {
		return nil, err
	}

	log.Info("Using store at ", storePath)
	log.Info("Using index at ", indexPath)


	op := &Operator{
		store:     s,
		index:     idx,
		options:   options,
		reposPath: path.Join(dataPath, "repos"),
	}

	op.rc = NewRepoCache(op.generateRepo)

	return op, nil
}

// generateRepo mainly used by repo cache
func (op *Operator) generateRepo(repoURI string) (*repository.Repo, error) {
	options, err := op.store.GetRepoOptions(repoURI)

	if err != nil {
		return nil, err
	}
	if options == nil {
		return nil, store.ErrRepoNotFound
	}

	hash, err := op.store.GetMeta(repoURI, []byte(mkURIHash))
	if err != nil {
		return nil, err
	}

	hashText := fmt.Sprintf("%x", hash)
	repoPath := path.Join(op.reposPath, hashText)

	return repository.NewRepo(repoPath, options)
}

func (op *Operator) AddRepo(options *repository.RepoOptions) error {

	uri := options.GetNormalizedURI()

	h := sha1.New()
	h.Write([]byte(uri))
	uriHash := h.Sum(nil)

	// Add repo options
	err := op.store.PutRepo(options, map[string][]byte{
		mkURIHash: uriHash,
		mkLastUpdate: []byte(time.Now().UTC().Format(time.RFC3339)),
	})
	if err != nil {
		return err
	}

	// Clear the cache to update on next request
	op.rc.Delete(uri)

	return nil
}

func (op *Operator) IndexRepo(repo *repository.Repo) error {

	fileInfos, err := repo.GetFileInfos("/", true)
	if err != nil {
		return err
	}

	indexedFiles, err := op.index.GetURIsByPrefix(repo.GetURI())
	_ = indexedFiles

	repoFiles := map[string]os.FileInfo{}

	for fileURI, info := range fileInfos {
		ext := strings.ToLower(filepath.Ext(info.Name()))

		if ext != ".md" {
			continue
		}

		repoFiles[fileURI] = info
	}

	// Remove deleted files
	for _, previousFile := range indexedFiles {
		_, notRemoved := repoFiles[previousFile]
		if !notRemoved { // File is not exist any more
			err := op.index.DeleteDocument(previousFile)
			if err != nil {
				return err
			}
		}
	}

	// Index all files
	for fileURI, info := range repoFiles {

		fileContent, err := repo.ReadFile(fileURI)

		if err != nil {
			return err
		}

		doc := &index.Document{
			URI:          fileURI,
			Title:        info.Name(),
			Content:      fmt.Sprintf("%s", fileContent),
			LastModified: info.ModTime(),
		}

		err = op.index.IndexDocument(doc)

		if err != nil {
			return err
		}

	}

	return nil

}

// SyncRepo synchronizes repo and corresponding indexes. Repo must be added before otherwise error returned.
func (op *Operator) SyncRepo(repoURI string) error {

	repo, err := op.rc.Get(repoURI)

	if err != nil {
		return err
	}

	alreadyUpToDate, err := repo.Sync()
	if err != nil {
		return err
	}

	log.Info("Repo synced:", repoURI, " already up to date: ", alreadyUpToDate)

	return op.IndexRepo(repo)

}

func (op *Operator) Search(queryString string) (*index.SearchResult, error) {

	return op.index.SearchContent(queryString)

}

