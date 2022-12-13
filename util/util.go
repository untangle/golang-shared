package util

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"strings"
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
