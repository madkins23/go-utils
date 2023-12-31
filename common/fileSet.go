package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type FilterFn func(key string, val bool) bool
type MappingFn func(key string, val bool) error

type FileSet struct {
	path string
	file *os.File
	set  map[string]bool
	lock sync.Mutex
}

var (
	ErrFileSetMapBreak = errors.New("Break map early")
)

func NewFileSet(path string) (*FileSet, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return nil, errors.Wrap(err, "checking for file")
		}

		// File did not previously exist, create it.
		file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			return nil, errors.Wrap(err, "creating file")
		}

		err = file.Close()
		if err != nil {
			return nil, errors.Wrap(err, "closing new file")
		}

		_, err = os.Stat(path)
		if err != nil {
			return nil, errors.Wrap(err, "checking new file")
		}
	}

	// Read items from file.
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "opening file")
	}

	fileSet := &FileSet{
		path: path,
		file: file,
		set:  make(map[string]bool),
		lock: sync.Mutex{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			if !strings.HasPrefix(line, "-") {
				fileSet.set[line] = true
			} else if fileSet.set[line] {
				delete(fileSet.set, line)
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "scanning file")
	}

	return fileSet, nil
}

func (fileSet *FileSet) Add(item string) error {
	fileSet.lock.Lock()
	defer fileSet.lock.Unlock()

	if _, ok := fileSet.set[item]; ok {
		return nil
	}

	var err error
	fileSet.set[item] = true
	if _, err = fileSet.file.WriteString(item + "\n"); err == nil {
		err = fileSet.file.Sync()
	}

	return err
}

func (fileSet *FileSet) Filtered(fn FilterFn) (*FileSet, error) {
	_ = fileSet.Close()
	fileSet.lock.Lock()
	defer fileSet.lock.Unlock()

	var err error

	oldPath := fileSet.path + "-old"
	olderPath := fileSet.path + "-older"

	if err = os.Rename(oldPath, olderPath); err != nil {
		fmt.Printf("*** Failed to rename %s to %s:\n***   %s\n", oldPath, olderPath, err.Error())
	}

	if err = os.Rename(fileSet.path, oldPath); err != nil {
		fmt.Printf("*** Failed to rename %s to %s:\n***   %s\n", fileSet.path, oldPath, err.Error())
	}

	// TODO: Nice to have cleanup code if the following fails, copy file back in place.

	newSet, err := NewFileSet(fileSet.path)
	if err != nil {
		return nil, errors.Wrap(err, "filtering file set")
	}

	for key, val := range fileSet.set {
		if fn(key, val) {
			if err = newSet.Add(key); err != nil {
				break
			}
		}
	}

	return newSet, err
}

/*
	// Clean up failed list.
	oldFailed := failed
	oldFailedPath := filepath.Join(statusPath, "old-failed.txt")

	if err = os.Rename(failedPath, oldFailedPath); err != nil {
		_, _ = fmt.Printf("*** Failed to rename failure file: %s\n", err.Error())
		return
	}

	if failed, err = utility.NewFileSet(filepath.Join(statusPath, "failed.txt")); err != nil {
		_, _ = fmt.Printf("*** Failed to create new failure file set: %s\n", err.Error())
		return
	}

	if err = oldFailed.Close(); err != nil {
		_, _ = fmt.Printf("*** Failed to close old failure file set: %s\n", err.Error())
		return
	}

	for item, isFailure := range oldFailed.Map() {
		if isFailure && !finished.Has(item) {
			if err = failed.Add(item); err != nil {
				_, _ = fmt.Printf("*** Failed to add failure %s to file set: %s\n", item, err.Error())
				return
			}
		}
	}

	if err = os.Rename(oldFailedPath, filepath.Join(statusPath, "older-failed.txt")); err != nil {
		_, _ = fmt.Printf("*** Failed to rename old failure file: %s\n", err.Error())
		return
	}

	waiter := sync.WaitGroup{}
	for i := 0; i < failedAtStartup; {
		failure := failed.Pop()
		if failure == "" {
			break
		}

		match := regexURL.FindStringSubmatch("'" + failure + "'")
		if match == nil {
			fmt.Printf("!!! Unable to match failure:\n!!!   `%s`\n", failure)
		} else if !finished.Has(match[1]) {
			waiter.Add(1)
			chanPages <- pageTask{
				urlWABAC: match[1],
				uriPage3: match[3],
				date:     match[2],
				wait:     &waiter,
			}
			i++
		}
	}
*/

func (fileSet *FileSet) Close() error {
	fileSet.lock.Lock()
	defer fileSet.lock.Unlock()
	return fileSet.file.Close()
}

func (fileSet *FileSet) Has(item string) bool {
	fileSet.lock.Lock()
	defer fileSet.lock.Unlock()
	return fileSet.set[item]
}

func (fileSet *FileSet) Map(fn MappingFn) error {
	fileSet.lock.Lock()
	defer fileSet.lock.Unlock()

	for key, val := range fileSet.set {
		if err := fn(key, val); err != nil {
			if err == ErrFileSetMapBreak {
				return nil // not an error
			}
			return err
		}
	}

	return nil
}

func (fileSet *FileSet) Pop() string {
	for item := range fileSet.set {
		if item != "" && fileSet.set[item] {
			delete(fileSet.set, item)
			return item
		}
	}

	return ""
}

func (fileSet *FileSet) Size() uint {
	fileSet.lock.Lock()
	defer fileSet.lock.Unlock()
	return uint(len(fileSet.set))
}
