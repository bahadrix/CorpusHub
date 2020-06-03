package repository

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/bahadrix/corpushub/util"
	"github.com/go-git/go-git/v5"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	ggssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)

type RepoOptions struct {
	Title      string `json:"title,omitempty"`
	URL        string `json:"URL,omitempty"`
	Branch     string `json:"branch,omitempty"`
	PrivateKey []byte `json:"privateKey,omitempty"`
}

func (ro *RepoOptions) GetNormalizedURI() string {
	return normalizeGitURL(ro.URL)
}

func (ro *RepoOptions) Serialize() ([]byte, error) {
	return json.Marshal(ro)
}

func DeserializeRepoOptions(data []byte ) (*RepoOptions, error) {
	var ro RepoOptions
	err := json.Unmarshal(data, &ro)
	return &ro, err
}

type Repo struct {
	gitrepo   *git.Repository
	worktree  *git.Worktree
	auth      transport.AuthMethod
	localPath string
	options   *RepoOptions
	repoURI   string
}

func normalizeGitURL(url string) string {

	normalizedURL := strings.TrimRight(strings.TrimSpace(strings.ToLower(url)), "/")

	// strip scheme
	schemeSplit := strings.Split(url, "://")
	normalizedURL = schemeSplit[len(schemeSplit)-1]

	// strip user
	userSplit := strings.Split(normalizedURL, "@")
	normalizedURL = userSplit[len(userSplit)-1]

	// trim .git extension
	return strings.TrimSuffix(normalizedURL, ".git")

}

func NewRepo(localPath string, options *RepoOptions) (*Repo, error) {

	signer, err := ssh.ParsePrivateKey(options.PrivateKey)

	if err != nil {
		return nil, err
	}

	if options.Branch == "" {
		options.Branch = "master"
	}

	return &Repo{
		auth: &ggssh.PublicKeys{
			User:   "git",
			Signer: signer,
		},
		localPath: localPath,
		options:   options,
		repoURI:   fmt.Sprintf("repo://%s/%s", normalizeGitURL(options.URL), options.Branch),
	}, nil
}

func (r *Repo) GenerateFullPath(path string) string {
	return fmt.Sprintf("%s/%s", r.repoURI, strings.Trim(path, "/"))
}

func (r *Repo) GetURI() string {
	return r.repoURI
}

func (r *Repo) GetFileInfos(path string, recursive bool) (infoMap map[string]*os.FileInfo, err error) {

	if path == "" {
		path = "/"
	}

	infoMap = make(map[string]*os.FileInfo)

	directoryQueue := list.New()

	directoryQueue.PushBack(path)

	for directoryQueue.Len() > 0 {
		dirPath := fmt.Sprintf("%s", directoryQueue.Remove(directoryQueue.Front()))

		files, err := r.worktree.Filesystem.ReadDir(fmt.Sprintf("%s", dirPath))

		if err != nil {
			return nil, err
		}

		for _, info := range files {
			if info.Name() == ".git" {
				continue
			}
			fpath := filepath.Join(dirPath, info.Name())
			fileURI := r.GenerateFullPath(fpath)
			infoMap[fileURI] = &info

			if recursive && info.IsDir() {
				directoryQueue.PushBack(fpath)
			}

		}

	}

	return infoMap, err
}

func (r *Repo) Sync() (alreadyUpToDate bool, err error) {

	alreadyUpToDate = false

	isCloned := false

	if r.gitrepo == nil {
		var localDirectoryExists bool
		localDirectoryExists, err = util.FileExists(r.localPath)

		if err != nil {
			return
		}

		if localDirectoryExists {
			r.gitrepo, err = git.PlainOpen(r.localPath)
		} else {
			r.gitrepo, err = git.PlainClone(r.localPath, false, &git.CloneOptions{
				Auth:          r.auth,
				URL:           r.options.URL,
				ReferenceName: plumbing.NewBranchReferenceName(r.options.Branch),
			})
			isCloned = true
		}

		if err != nil {
			return
		}

		r.worktree, err = r.gitrepo.Worktree()

		if err != nil {
			return
		}

	}

	if !isCloned {
		err = r.worktree.Pull(&git.PullOptions{
			Auth:          r.auth,
			ReferenceName: plumbing.NewBranchReferenceName(r.options.Branch),
		})

		if err == git.NoErrAlreadyUpToDate {
			alreadyUpToDate = true
			err = nil
		}
	}

	return alreadyUpToDate, err

}

func (r *Repo) ReadFile(fileURI string) ([]byte, error) {
	filePath := strings.TrimPrefix(fileURI, r.repoURI)

	file, err := r.worktree.Filesystem.Open(filePath)

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	chunk := make([]byte, 1024)

	for {
		n, err := file.Read(chunk)

		buf.Write(chunk[:n])

		if err == io.EOF {
			break
		} else if err != nil {
			return buf.Bytes(), err
		}
	}

	return buf.Bytes(), nil

}
