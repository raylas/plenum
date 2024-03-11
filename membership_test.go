package plenum

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMembership(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	address := "127.0.0.1"
	port := 1111
	members := []string{"127.0.0.1:1111"}

	serfInstance, _, err := createMembership(logger, address, port, members)
	if err != nil {
		t.Fatal(err)
	}

	p := &Plenum{
		membership: serfInstance,
	}

	membership := p.Membership()

	assert.Equal(t, 1, len(membership))
	assert.Equal(t, "127.0.0.1:1111", membership[0])
}
