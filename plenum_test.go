package plenum

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvene(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	address := "127.0.0.1"
	port := 3333
	members := []string{"127.0.0.1:3333"}

	plenum, err := Convene(logger, address, port, members)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, plenum)
	assert.NotNil(t, plenum.logger)
	assert.NotNil(t, plenum.membership)
	assert.NotNil(t, plenum.consensus)
	assert.NotNil(t, plenum.Procedures)
	assert.NotNil(t, plenum.IsLeader)
	assert.NotNil(t, plenum.IsClosed)
}
