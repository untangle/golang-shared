package util

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// ContainsString checks if a string is contained in an array of strings
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// StringArrayToDB converts a string array into a single string, using pipe as a delimiter
func StringArrayToDB(s []string) string {
	return strings.Join(s, "|")
}

// WaitGroupDoneOrTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func WaitGroupDoneOrTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}

// Helper function to randomize the order of an array
func RandomizeSlice[T any](slice []T) {
	n := len(slice)

	// Shuffle
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Pulls specified filenames out of a tar or tar.gz, depending on if isGzip is set.
// An error and nil are are returned if a given fileName cannot be found
// A map of the found fileNames to the found file data will be returned. To check for a
// files existence, strings.Contains() is used on the full path of each file. For example, /tmp/settings.json
// and /example/settings.json will both be found. However, only one of the two will be added to the returned map.
// If files of the same name are needed, more of their paths will have to be provided to differentiate them.
func ExtractFilesFromTar(b []byte, isGzip bool, fileName ...string) (map[string][]byte, error) {
	var fileReader io.Reader = bytes.NewReader(b)
	var err error
	if isGzip {
		fileReader, err = gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
	}

	foundFiles := make(map[string][]byte)

	tarReader := tar.NewReader(fileReader)
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if header.Typeflag == tar.TypeReg {
			for _, name := range fileName {
				if strings.Contains(header.Name, name) {
					var settingsData []byte

					settingsData, err := io.ReadAll(tarReader)
					if err != nil {
						return nil, err
					}

					foundFiles[name] = settingsData
				}
			}

		}
	}

	return foundFiles, nil
}
