package repository

import "encoding/json"

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

