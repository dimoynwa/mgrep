package main

import (
	"fmt"
	"mgrep/worker"
	"mgrep/worklist"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
)

func discoverDirs(wl *worklist.Worklist, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Reading problem : ", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			newPath := filepath.Join(path, entry.Name())
			discoverDirs(wl, newPath)
		} else {
			newPath := filepath.Join(path, entry.Name())
			wl.Add(worklist.NewJob(newPath))
		}
	}
}

var args struct {
	SearchTerm string `arg:"positional,required"`
	SearchDir  string `arg:"positional"`
}

func main() {
	arg.MustParse(&args)
	var workerWg sync.WaitGroup

	wl := worklist.New(100)
	results := make(chan worker.Result, 100)

	numWorkers := 10

	workerWg.Add(1)
	go func() {
		defer workerWg.Done()
		discoverDirs(&wl, args.SearchDir)
		wl.Finalize(numWorkers)
	}()

	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for {
				entry := wl.Next()
				if entry.Path != "" {
					res := worker.FindInFile(entry.Path, args.SearchTerm)
					if res != nil {
						for _, r := range res.Inner {
							results <- r
						}
					}
				} else {
					return
				}
			}
		}()
	}

	blockWorkersWg := make(chan struct{})
	go func() {
		workerWg.Wait()
		close(blockWorkersWg)
	}()

	var displayWg sync.WaitGroup
	displayWg.Add(1)

	go func() {
		for {
			select {
			case <-blockWorkersWg:
				if len(results) == 0 {
					displayWg.Done()
					return
				}
			case r := <-results:
				fmt.Println(r.Path, ":", r.LineNum, ":", r.Line)
			}
		}
	}()

	displayWg.Wait()
}
