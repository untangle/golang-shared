package settings

import (
	"encoding/json"
	"fmt"
	"io"
)

// PathUnmarshaller unmarshals objects in a JSON object given a path
// through the JSON input object to that desired target object.
// e.g. given {"x": {"y": "z"}} and the path "x", "y", we can extract
// the JSON for just "z".
type PathUnmarshaller struct {
	// decoder object that we use to get tokens from.
	decoder *json.Decoder

	// should we call UseNumber when we decode the object we find?
	useNumber bool

	// the path being searched for.
	searchedForPath []string
}

func (unm *PathUnmarshaller) formatError(msg string, params ...interface{}) error {
	paramsList := make([]interface{}, len(params), len(params)+1)
	copy(paramsList, params)
	paramsList = append([]interface{}{fmt.Sprintf("%#v", unm.searchedForPath)},
		params...)
	return fmt.Errorf(
		"path unmarshaller: encountered error: looking for path %v "+msg,
		paramsList...)
}

func (unm *PathUnmarshaller) getToken() (json.Token, error) {
	next, err := unm.decoder.Token()
	if err != nil && err == io.EOF {
		return nil, unm.formatError("can't find path, EOF")
	} else if err != nil {
		return nil, unm.formatError("path unmarshaller: json error: %w", err)
	}
	return next, nil

}

// ignore JSON in the token stream until we hit delim.
func (unm *PathUnmarshaller) ignoreUntil(delim json.Delim) error {
	for {
		next, err := unm.getToken()
		if err != nil {
			return err
		}
		switch tok := next.(type) {
		case json.Delim:
			switch tok {
			case delim:
				return nil
			case '[':
				if err := unm.ignoreUntil(']'); err != nil {
					return err
				}
			case '{':
				if err := unm.ignoreUntil('}'); err != nil {
					return err
				}
			default:
				return unm.formatError("bad JSON, got unexpected token: %T (%s)",
					tok,
					tok)
			}
		default:
			return unm.formatError("unexpected token type: %T (%s)", tok, tok)
		case nil, string, int, int64, float32,
			float64, bool, json.Number:
		}
	}
}

// ignore the entire next object in the stream.
func (unm *PathUnmarshaller) ignoreNextObject() error {
	next, err := unm.getToken()
	if err != nil {
		return err
	}
	switch tok := next.(type) {
	case json.Delim:
		switch tok {
		case '[':
			if err := unm.ignoreUntil(']'); err != nil {
				return err
			}
		case '{':
			if err := unm.ignoreUntil('}'); err != nil {
				return err
			}
		default:
			return unm.formatError("bad/unexpected JSON token: %T (%s)", tok, tok)
		}
	default:
		return unm.formatError("unexpected json token: %T (%s)", tok, tok)
	case nil, string, int, int64, float32,
		float64, bool, json.Number:
		return nil
	}
	return nil
}

const (
	completeMatch = iota // the path is completely matched, return now.
	partialMatch         // we so far match but we haven't gone all the way.
	noMatch              // This is the wrong path, stop going down it.
)

type matchResult int32

// pathMatchStatus returns the status of the current path -- fullPath
// is what we want, partialPath is what we have so far. See the consts
// above.
func pathMatchStatus(fullPath []string, partialPath []string) matchResult {
	if len(fullPath) < len(partialPath) {
		return noMatch
	}

	for i, elem := range partialPath {
		if elem != fullPath[i] {
			return noMatch
		}
	}
	if len(fullPath) == len(partialPath) {
		return completeMatch
	}
	return partialMatch

}

// processObject is called with a 'value' object from a parent
// object. It processes it based on type:
// 1. Array -- ignore (eat it up).
// 2. Object -- recurse into it.
// 3. other -- ignore (eat it up).
func (unm *PathUnmarshaller) processObject(
	currentPath []string,
	searchPath []string) (io.Reader, error) {
	for {
		next, err := unm.getToken()
		if err != nil {
			return nil, err
		}
		switch tok := next.(type) {
		case json.Delim:
			switch tok {
			case '{':
				if result, err := unm.searchObject(currentPath,
					searchPath); result != nil {
					return result, err
				} else if err != nil {
					return nil, err
				}
			case '[':
				if err := unm.ignoreUntil(']'); err != nil {
					return nil, err
				}
			default:
				return nil, unm.formatError(
					"unexpected json delim: %s", tok)
			}
		case string, int64, int, float64, float32:
			return nil, nil
		default:
			return nil, unm.formatError("unexpected json token: %s", tok)
		}
	}
}

// searchObject is given a dict object and looks through the keys in
// that object to see if they match the desired input path in
// searchPath. It makes a recursive call to itself if it can extend
// the matching path.
func (unm *PathUnmarshaller) searchObject(
	currentPath []string,
	searchPath []string) (io.Reader, error) {
	for {
		next, err := unm.getToken()
		if err != nil {
			return nil, err
		}
		switch tok := next.(type) {
		case string:
			newPath := make([]string, len(currentPath), len(currentPath)+1)
			copy(newPath, currentPath)
			newPath = append(newPath, tok)
			matchResult := pathMatchStatus(searchPath, newPath)
			switch matchResult {
			case completeMatch:
				return unm.decoder.Buffered(), nil
			case partialMatch:
				return unm.processObject(newPath, searchPath)
			case noMatch:
				if err := unm.ignoreNextObject(); err != nil {
					return nil, err
				}
			}
		case json.Delim:
			if tok == '}' {
				return nil, nil
			}
			return nil, unm.formatError("unexpected JSON delim: %s", tok)
		default:
			return nil, unm.formatError("unexpected JSON token: %T (%s)", tok, tok)
		}
	}
}

// The JSON tokenizer does not have : tokens. So when we find the key
// we look for, and get the stream for that point, there is still a
// ':' in the stream. So we eat it up.
func (unm *PathUnmarshaller) fastForwardToValue(reader io.Reader) error {
	bytes := []byte{0}
	for {
		n, err := reader.Read(bytes)
		if n != 1 || err != nil {
			return unm.formatError(
				"couldn't read expected object from JSON stream: %w",
				err)
		} else if bytes[0] == ':' {
			return nil
		}
	}
}

// UseNumber will tell the underlying decoder to UseNumber using the
// json.Decoder.UseNumber method, for the object returned via
// UnmarshalAtPath.
func (unm *PathUnmarshaller) UseNumber() {
	unm.useNumber = true
}

// NewPathUnmarshaller creates a new PathUnmarshaller that will read
// from reader. reader should point to valid JSON. The JSON object
// should be a dict object (curly braces), not any other kind of JSON
// object.
func NewPathUnmarshaller(reader io.Reader) *PathUnmarshaller {
	return &PathUnmarshaller{
		decoder: json.NewDecoder(reader),
	}
}

// UnmarshalAtPath unmarshalls the JSON object found at path into
// output. path is a list of keys into the object. For example, in the
// object {"one": {"two": 3}}, the path "one", "two" corresponds to
// the value 3.  If the value is not found or some other error occurs,
// we return an error, otherwise we unmarshal into output. Output can
// be anything json.Unmarshal accepts.
func (unm *PathUnmarshaller) UnmarshalAtPath(output interface{}, path ...string) error {
	unm.searchedForPath = path
	outputReader, err := unm.processObject([]string{}, path)
	if err != nil {
		return err
	} else if outputReader == nil {
		return unm.formatError("couldn't find path: %s", path)
	}
	err = unm.fastForwardToValue(outputReader)
	if err != nil {
		return err
	}
	newDecoder := json.NewDecoder(outputReader)
	if unm.useNumber {
		newDecoder.UseNumber()
	}
	return newDecoder.Decode(output)
}
