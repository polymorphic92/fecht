package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func main() {

	repos := make(map[string]*git.Repository)
	for _, arg := range os.Args[1:] {
		workspace := getWorkspace(arg)
		fmt.Printf("WORKSPACE(ARG) :: %v(%v)\n", workspace, arg)
		findRepos(workspace, repos)
	}
	fmt.Printf("\n\n")
	fetchRepos(repos)

}
func checkError(e error) {
	if e != nil {
		panic(e) // TODO fail gracefully instead
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

	for path, repo := range repos {
		fmt.Printf("Fetching repo : %v\n", path)
		fetchErr := repo.Fetch(opts)
		if fetchErr != nil {
			fmt.Printf("\tissue fetching repo(%v) :: %v\n", path, fetchErr)
		}
	}
}

func getWorkspace(path string) string {

	fullPath, err := filepath.Abs(os.ExpandEnv(path))
	checkError(err)
	_, err = os.Stat(fullPath)
	checkError(err)

	return fullPath
}
