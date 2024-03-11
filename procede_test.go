package plenum

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProcede(t *testing.T) {
	address := "127.0.0.1"
	port := 4444
	members := []string{}

	p, err := Convene(nil, address, port, members)
	if err != nil {
		t.Fatal(err)
	}

	go p.Procede()

	isLeader := false
	go func(p *Plenum) {
		for <-p.IsLeader {
			isLeader = true
		}
	}(p)

	time.Sleep(2 * time.Second) // Wait for consensus to form

	assert.True(t, isLeader)
}
