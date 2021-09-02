package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {
	workspace := getWorkspace(os.Args[1])
	fmt.Println("workspace:", workspace)
	update(workspace)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getWorkspace(workspace string) string {
	if dirExists(workspace) {
		return workspace
	}

	env := os.Getenv(workspace)

	if dirExists(env) {
		return env
	}

	return ""
}

func dirExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func update(path string) {
	dirs, err := ioutil.ReadDir(path)
	check(err)
	numberOfRepos := len(dirs)
	var wg sync.WaitGroup
	fmt.Printf("Repo: %v\n", numberOfRepos)
	wg.Add(numberOfRepos)
	for _, dir := range dirs {
		if dir.IsDir() {
			fullPath := path + "/" + dir.Name() + "/.git"
			go func(p string) {
				defer wg.Done()
				cmd := exec.Command("git", "--git-dir", p, "remote", "update")
				var out bytes.Buffer
				var stderr bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &stderr
				err := cmd.Run()
				if err != nil {
					log.Fatalf("cmd failed with:\nstderror: %s\nrepo-path %s", stderr.String(), p)
				}
				fmt.Printf("Updated Repo: %s\n", p)
			}(fullPath)
		}
	}
	wg.Wait()
}
