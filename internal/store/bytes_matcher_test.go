package store_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
)

type bytesMatcher string

func (m bytesMatcher) Matches(x interface{}) bool {
	actualBytes, ok := x.([]byte)
	if !ok {
		return false
	}

	var actual interface{}
	if err := json.Unmarshal(actualBytes, &actual); err != nil {
		log.Printf("failed unmarshall actual bytes %q: %v", string(actualBytes), err)
		return false
	}

	expectedFileName := filepath.Join("testdata", string(m))
	if update {
		expectedBytes, err := json.MarshalIndent(&actual, "", "    ")
		if err != nil {
			log.Printf("failed marshall actual %+v: %v", actual, err)
			return false
		}

		expectedBytes = append(expectedBytes, '\n')
		if err := ioutil.WriteFile(expectedFileName, expectedBytes, 0600); err != nil {
			log.Printf("failed to write %q: %v", expectedFileName, err)
			return false
		}

		return true
	}

	expectedBytes, err := ioutil.ReadFile(expectedFileName)
	if err != nil {
		log.Printf("failed to read %q: %v", expectedFileName, err)
		return false
	}

	var expected interface{}
	if err := json.Unmarshal(expectedBytes, &expected); err != nil {
		log.Printf("failed unmarshall expected bytes %q: %v", string(expectedBytes), err)
		return false
	}

	if !reflect.DeepEqual(expected, actual) {
		log.Printf("not equal %q\nExpected: %+v\nActual: %+v\n", expectedFileName, expected, actual)
		return false
	}

	return true
}

func (m bytesMatcher) String() string {
	return fmt.Sprintf("file `%s`", string(m))
}
