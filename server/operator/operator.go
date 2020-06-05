package operator

import (
	"crypto/sha1"
	"github.com/bahadrix/corpushub/index"
	"github.com/bahadrix/corpushub/repository"
	"github.com/bahadrix/corpushub/server/store"
	log "github.com/sirupsen/logrus"
	"path"
	"time"
)



type Operator struct {
	store *store.RepoStore
	index *index.Index
	options *Options
	repoPath string
}

type Options struct {

	// Disables synchronization cycle. Normally it starts automatically
	DisableSync bool

	// Interval between repo sync polls
	SyncInterval time.Duration

}

var ( // Meta Keys
	mkUriHash = []byte("uriHash")
)

func NewOperator(dataPath string, options *Options) (*Operator, error) {

	storePath := path.Join(dataPath, "store.db")
	indexPath := path.Join(dataPath, "index")

	if options == nil {
		options = &Options{
			SyncInterval: 30 * time.Second,
		}
	}

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



	return &Operator{
		store:   s,
		index:   idx,
		options: options,
		repoPath: path.Join(dataPath, "repos"),
	}, nil

}

func (op *Operator) StartSyncCycle() {
	go op.syncCycle()
}

func (op *Operator) AddRepo(options *repository.RepoOptions) error {

	//updateTime := time.Now()

	// Add repo options
	err := op.store.PutRepoOptions(options)
	if err != nil {
		return err
	}

	// Add uri hash
	uri := options.GetNormalizedURI()
	h := sha1.New()
	h.Write([]byte(uri))
	uriHash := h.Sum(nil)
	return op.store.PutMeta(uri, mkUriHash, uriHash)

}

// SyncRepo synchronizes repo and corresponding indexes. Repo must be added before otherwise error returned.
func (op *Operator) SyncRepo(repoURI string) error {

	return nil

}

func (op *Operator) syncCycle() {

	op.store.FindAll()

}