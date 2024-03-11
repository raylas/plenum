package plenum

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConsensus(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	address := "127.0.0.1"
	port := 2222
	members := []string{"127.0.0.1:2222"}

	raftInstance, _, err := createConsensus(logger.Writer(), os.TempDir(), address, port, members)
	if err != nil {
		t.Fatal(err)
	}

	p := &Plenum{
		consensus: raftInstance,
	}

	time.Sleep(2 * time.Second) // Wait for consensus to form

	consensus := p.Consensus()

	// assert.Equal(t, 1, len(consensus))
	assert.Equal(t, "127.0.0.1:2222", consensus)
}
