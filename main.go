package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func main() {

	repos := make(map[string]*git.Repository)
	for _, arg := range os.Args[1:] {
		workspace := getWorkspace(arg)
		fmt.Printf("WORKSPACE(ARG) :: %v(%v)\n", workspace, arg)
		findRepos(workspace, repos)
	}

	for r := range repos {
		fmt.Printf("\nREPO :: %v", r)
	}

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

func getWorkspace(path string) string {

	fullPath, err := filepath.Abs(os.ExpandEnv(path))
	checkError(err)
	_, err = os.Stat(fullPath)
	checkError(err)

	return fullPath
}
