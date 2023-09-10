package sherlock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMicroformatsActor1(t *testing.T) {
	testActor(t, "actor-microformat-1.html")
}

func TestMicroformatsActor2(t *testing.T) {
	testActor(t, "actor-microformat-2.html")
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

	t.Log(acc.result)
}
