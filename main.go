package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	ggssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)



func getAuth() transport.AuthMethod {
	keyStr := "-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW\nQyNTUxOQAAACCONI/sdeaq9oDHHl+nz8ihZEHqYFSd2fKJcgYQMor8KQAAAKDKl+sTypfr\nEwAAAAtzc2gtZWQyNTUxOQAAACCONI/sdeaq9oDHHl+nz8ihZEHqYFSd2fKJcgYQMor8KQ\nAAAEC8HYlDnh9E94tNa7fdqkPrMouW4galiHofFyN51W8Sko40j+x15qr2gMceX6fPyKFk\nQepgVJ3Z8olyBhAyivwpAAAAG2JhaGFkaXJAQmFoYWRpcnMtaUJhZy5sb2NhbAEC\n-----END OPENSSH PRIVATE KEY-----\n"

	signer,  err := ssh.ParsePrivateKey([]byte(keyStr))
	if err != nil {
		panic(err)
	}

	return &ggssh.PublicKeys{
		User:                  "git",
		Signer:                signer,
	}

}

func clone() {

	_, err := git.PlainClone("tmp/doctest", false, &git.CloneOptions{
		Auth: getAuth(),
		URL: "git@github.com:bahadrix/doctest.git",

	})

	if err != nil {
		panic(err)
	}
}

func main() {
	// git@gitlab.vlstats.com:6161/bahadir/copsgen.git


	r, err := git.PlainOpen("tmp/doctest")

	if err != nil {
		panic(err)
	}

	w, _ := r.Worktree()

	err = w.Pull(&git.PullOptions{
		Auth: getAuth(),
	})

	h, _ := r.Head()
	c, _ := r.CommitObject(h.Hash())

	stats, _ := c.Stats()

	_ = stats
	for _, stat := range stats {
		_ = stat
	}

	commits, _ := r.CommitObjects()
	_ = commits
	print("1")
}
