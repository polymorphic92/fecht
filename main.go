package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/pterm/pterm"
)

func main() {

	repos := make(map[string]*git.Repository)
	for _, arg := range os.Args[1:] {
		workspace := getWorkspace(arg)
		findRepos(workspace, repos)
	}
	fetchRepos(repos)

}
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func findRepos(workspace string, repos map[string]*git.Repository) {
	err := filepath.Walk(workspace, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && strings.HasSuffix(path, ".git") {
			r, repoErr := git.PlainOpen(path)
			if repoErr == nil {
				repos[path] = r
			}
		}
		return nil
	})
	checkError(err)
}

func fetchRepos(repos map[string]*git.Repository) {
	opts := &git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	}
	pterm.EnableDebugMessages()
	numberOfRepos := len(repos)

	pterm.FgGreen.Println("Number of Repos:", numberOfRepos, "\n")
	var wg sync.WaitGroup
	wg.Add(numberOfRepos)
	for path, repo := range repos {
		go func(p string, r *git.Repository) {
			defer wg.Done()
			fetchErr := r.Fetch(opts)
			if fetchErr != nil {
				if fetchErr == git.NoErrAlreadyUpToDate {
					pterm.Info.Println("Repo:", p, fetchErr)
				} else {
					pterm.Error.Println("Eccountered error %v whiling fetching %v repo", fetchErr, p)
				}
			} else {
				pterm.Success.Println("Fetched: %v\n", p)
			}

		}(path, repo)
	}
	wg.Wait()
}

func getWorkspace(path string) string {

	fullPath, err := filepath.Abs(os.ExpandEnv(path))
	checkError(err)
	_, err = os.Stat(fullPath)
	checkError(err)

	return fullPath
}
