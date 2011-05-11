
package main

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGITForLeaks(t *testing.T) {
	gitcmd := envOrDefault("git", "/usr/local/git/bin/git")
	repo := "/Users/petar/popalg.org-git"
	filename := "sparsification-by-spanners"
	for i := 0; i < 10000; i++ {
		_, _, err := GITGetCreateUpdateTime(gitcmd, repo, filename)
		if err != nil {
			t.Errorf("git: %s", err)
		}
		if i % 100 == 0 {
			mem, _ := runtime.MemProfile(nil, false)
			fmt.Printf("i= %d,    mem= %d             \r", i, mem)
		}
	}
	fmt.Printf("\n")
}
