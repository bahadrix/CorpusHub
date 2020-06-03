package store

import (
	"errors"
	"fmt"
	"github.com/bahadrix/corpushub/repository"
	bolt "go.etcd.io/bbolt"
)

var (
	// Bucket keys
	bkRepos   = []byte("REPOS")
	bkOptions = []byte("OPTIONS")
	bkMeta    = []byte("META")

	// Error Definitions
	ErrRepoNotFound = errors.New("repository not found")
)

// Store structure:
// REPOS: bucket
//     ↳ [REPO URI]: bucket
//         ↳ options: repository.RepoOptions
type RepoStore struct {
	db *bolt.DB
}

// NewRepoStore returns new Repository Store. If db path exists repo loaded from there otherwise
// new database file created.
func NewRepoStore(dbPath string) (*RepoStore, error) {

	db, err := bolt.Open(dbPath, 0600, nil)

	if err != nil {
		return nil, err
	}

	return &RepoStore{
		db: db,
	}, nil
}

// PutRepoOptions creates or updates (upserting) options by its own normalized repo uri
func (s *RepoStore) PutRepoOptions(options *repository.RepoOptions) error {

	return s.db.Update(func(tx *bolt.Tx) error {

		// Get root REPOS bucket
		bRepos, err := tx.CreateBucketIfNotExists(bkRepos)
		if err != nil {
			return err
		}

		// Get repo key
		repoKey := []byte(options.GetNormalizedURI())

		// Get bucket for repoKey
		bRepo, err := bRepos.CreateBucketIfNotExists(repoKey)

		if err != nil {
			return err
		}

		// Serialize options
		repoData, err := options.Serialize()

		if err != nil {
			return err
		}

		// Put data to repo bucket
		return bRepo.Put(bkOptions, repoData)

	})
}

// GetRepoOptions returns options at given repository URI. Note that URI is not URL. URIs are normalized repository URLS.
func (s *RepoStore) GetRepoOptions(repoURI string) (options *repository.RepoOptions, err error) {

	err = s.db.View(func(tx *bolt.Tx) error {

		// Get root REPOS bucket
		bRepos := tx.Bucket(bkRepos)
		if bRepos == nil {
			return nil
		}

		// Get repo bucket
		repoKey := []byte(repoURI)
		bRepo := bRepos.Bucket(repoKey)

		// Get serialized options
		data := bRepo.Get(bkOptions)

		if data == nil {
			return nil
		}

		// Deserialize options and write it from the closure
		options, err = repository.DeserializeRepoOptions(data)

		return err
	})

	return
}

// DeleteRepo removes the repo from the store. If repo not exists it returns ErrRepoNotFound error
func (s *RepoStore) DeleteRepo(repoURI string) error {

	return s.db.Update(func(tx *bolt.Tx) error {
		// Get root REPOS bucket
		bRepos := tx.Bucket(bkRepos)
		if bRepos == nil {
			return ErrRepoNotFound
		}

		// Get repo bucket
		repoKey := []byte(repoURI)
		bRepo := bRepos.Bucket(repoKey)

		// Delete the repo if it exists
		if bRepo != nil {
			return bRepos.DeleteBucket(repoKey)
		}

		return nil

	})

}

func (s *RepoStore) FindAll() (items []string, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		bRepos := tx.Bucket(bkRepos)
		if bRepos == nil {
			return nil
		}

		items = make([]string, 0)
		return bRepos.ForEach(func(k, v []byte) error {
			items = append(items, fmt.Sprintf("%s", k))
			return nil
		})
	})

	return
}

// Exists returns true if the repo exists
func (s *RepoStore) Exists(repoURI string) (exists bool) {

	_ = s.db.View(func(tx *bolt.Tx) error {
		// Get root REPOS bucket
		bRepos := tx.Bucket(bkRepos)
		if bRepos == nil {
			return nil
		}

		// Get repo bucket
		repoKey := []byte(repoURI)
		bRepo := bRepos.Bucket(repoKey)

		exists = bRepo != nil

		return nil

	})

	return

}

// PutMeta adds custom key value pair to repo. If repo not exists it returns ErrRepoNotFound error
func (s *RepoStore) PutMeta(repoURI string, key []byte, value []byte) error {

	return s.db.Update(func(tx *bolt.Tx) error {
		// Get root REPOS bucket
		bRepos := tx.Bucket(bkRepos)
		if bRepos == nil {
			return ErrRepoNotFound
		}

		// Get repo bucket
		repoKey := []byte(repoURI)
		bRepo := bRepos.Bucket(repoKey)

		// Repo must be exist
		if bRepo == nil {
			return ErrRepoNotFound
		}

		bMeta, err := bRepo.CreateBucketIfNotExists(bkMeta)
		if err != nil {
			return err
		}

		return bMeta.Put(key, value)

	})

}

// GetMeta returns value at given meta key. If repo not exists it returns ErrRepoNotFound error.
// If meta key not found returns nil error and value
func (s *RepoStore) GetMeta(repoURI string, key []byte) (value []byte, err error) {

	err = s.db.View(func(tx *bolt.Tx) error {

		// Get root REPOS bucket
		bRepos := tx.Bucket(bkRepos)
		if bRepos == nil {
			return ErrRepoNotFound
		}

		// Get repo bucket
		repoKey := []byte(repoURI)
		bRepo := bRepos.Bucket(repoKey)

		// Repo must be exist
		if bRepo == nil {
			return ErrRepoNotFound
		}

		// Get meta bucket
		bMeta := bRepo.Bucket(bkMeta)

		if bMeta == nil {
			return nil
		}

		// Get value
		value = bMeta.Get(key)

		return nil

	})

	return
}

// DeleteMeta removes meta key. If repo not exists it returns ErrRepoNotFound error. If key not found nil error returned.
func (s *RepoStore) DeleteMeta(repoURI string, key []byte) error {

	return s.db.Update(func(tx *bolt.Tx) error {
		// Get root REPOS bucket
		bRepos := tx.Bucket(bkRepos)
		if bRepos == nil {
			return ErrRepoNotFound
		}

		// Get repo bucket
		repoKey := []byte(repoURI)
		bRepo := bRepos.Bucket(repoKey)

		// Repo must be exist
		if bRepo == nil {
			return ErrRepoNotFound
		}

		bMeta := bRepo.Bucket(bkMeta)

		if bMeta == nil {
			return nil
		}

		return bMeta.Delete(key)

	})

}
