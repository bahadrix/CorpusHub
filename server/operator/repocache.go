package operator

import (
	"github.com/bahadrix/corpushub/repository"
	"sync"
)

type RepoCache struct {
	cache map[string]*repository.Repo
	mutex sync.RWMutex
	generator func(repoURI string) (*repository.Repo, error)
}

func NewRepoCache(generator func(repoURI string) (*repository.Repo, error) ) *RepoCache {
	return &RepoCache{
		cache: map[string]*repository.Repo{},
		mutex:     sync.RWMutex{},
		generator: generator,
	}
}


func (c *RepoCache) Get(repoURI string) (*repository.Repo, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	repo, isCached := c.cache[repoURI]
	var err error

	if !isCached {
		repo, err = c.generator(repoURI)
		if err != nil {
			return nil, err
		}

		c.cache[repoURI] = repo
	}

	return repo, nil

}

// Delete removes entry from cache. Returns true if repo is found and removed.
func (c *RepoCache) Delete(repoURI string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, isCached := c.cache[repoURI]

	if isCached {
		delete(c.cache, repoURI)
	}
	return isCached
}