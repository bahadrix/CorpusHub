package repository

import "encoding/json"

type RepoOptions struct {
	Title      string `json:"title,omitempty" binding:"required"`
	URL        string `json:"URL,omitempty" binding:"required"`
	Branch     string `json:"branch,omitempty" binding:"required"`
	PrivateKey []byte `json:"privateKey,omitempty" binding:"required"`
}

func (ro *RepoOptions) GetNormalizedURI() string {
	return normalizeGitURL(ro.URL)
}

func (ro *RepoOptions) Serialize() ([]byte, error) {
	return json.Marshal(ro)
}

func DeserializeRepoOptions(data []byte) (*RepoOptions, error) {
	var ro RepoOptions
	err := json.Unmarshal(data, &ro)
	return &ro, err
}
