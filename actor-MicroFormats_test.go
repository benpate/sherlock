package sherlock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMicroformatsActor(t *testing.T) {
	testActor(t, "actor-MicroFormat-2.html")
}

func testActor(t *testing.T, filename string) {

	// Read the test file
	buffer, err := testFile(filename)
	require.Nil(t, err, "Error reading file: %s", filename)

	// Create an accumulator
	acc := actorAccumulator{
		url:  "TEST_FILE",
		body: buffer,
	}

	// Try to parse the file
	client := Client{}
	client.actor_MicroFormats(&acc)

	// spew.Dump(acc.result)
}
