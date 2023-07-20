package sherlock

import (
	"bytes"
	"io"
	"os"
)

// testFile is a simple helper that reads test files into a buffer
func testFile(filename string) (bytes.Buffer, error) {

	result := bytes.Buffer{}

	// Open the test file
	file, err := os.Open("./test-files/" + filename)

	if err != nil {
		return result, err
	}

	defer file.Close()

	// Read the test file into a buffer
	_, err = io.Copy(&result, file)

	if err != nil {
		return result, err
	}

	// Return result
	return result, nil
}
