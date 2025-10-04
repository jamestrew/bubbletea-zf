package files

import (
	"io/fs"
	"iter"
	"slices"
	"sync"

	"github.com/charlievieth/fastwalk"
)

type CSlice struct {
	data []string
	mu   sync.RWMutex
}

func (c *CSlice) Append(s string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = append(c.data, s)
}

func (c *CSlice) Seq() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range c.Seq2() {
			if !yield(v) {
				return
			}
		}
	}
}

func (c *CSlice) Seq2() iter.Seq2[int, string] {
	c.mu.RLock()
	items := make([]string, len(c.data))
	copy(items, c.data)
	c.mu.RUnlock()
	return func(yield func(int, string) bool) {
		for i, v := range items {
			if !yield(i, v) {
				return
			}
		}
	}
}

// GetWithChannel collects files using channel-based concurrency
// Preferred over GetWithCSlice for larger file sets
func GetWithChannel(rootDir string) ([]string, error) {
	pathsChan := make(chan string, 1000)
	paths := make([]string, 0, 1000)
	var wg sync.WaitGroup

	wg.Go(func() {
		for p := range pathsChan {
			paths = append(paths, p)
		}
	})

	config := fastwalk.Config{Follow: true}
	err := fastwalk.Walk(&config, rootDir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if path == "./.git" {
			return fs.SkipDir
		}
		if path == "." || de.IsDir() {
			return nil
		}

		pathsChan <- path

		return nil
	})

	close(pathsChan)
	wg.Wait()

	if err != nil {
		return nil, err
	}

	return paths, nil
}

// GetWithCSlice collects files using mutex-based concurrent slice
func GetWithCSlice(rootDir string) ([]string, error) {
	paths := &CSlice{}

	config := fastwalk.Config{Follow: true}
	err := fastwalk.Walk(&config, rootDir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if path == "./.git" {
			return fs.SkipDir
		}
		if path == "." || de.IsDir() {
			return nil
		}

		paths.Append(path)

		return nil
	})

	if err != nil {
		return nil, err
	}

	pathsList := slices.Collect(paths.Seq())
	return pathsList, nil
}
