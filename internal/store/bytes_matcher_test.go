package store_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

type bytesMatcher string

func (m bytesMatcher) Matches(x interface{}) bool {
	actual, ok := x.([]byte)
	if !ok {
		return false
	}

	actual = append(actual, '\n')

	expectedFileName := filepath.Join("testdata", string(m))
	if update {
		if err := ioutil.WriteFile(expectedFileName, actual, 0600); err != nil {
			log.Printf("failed to write %q: %v", expectedFileName, err)
			return false
		}
	}

	expected, err := ioutil.ReadFile(expectedFileName)
	if err != nil {
		log.Printf("failed to read %q: %v", expectedFileName, err)
		return false
	}

	return string(expected) == string(actual)
}

func (m bytesMatcher) String() string {
	return fmt.Sprintf("file `%s`", string(m))
}
