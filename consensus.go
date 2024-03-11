package plenum

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

// Consensus returns the current leader of the plenum
func (p *Plenum) Consensus() string {
	leader := p.consensus.Leader()
	return string(leader)
}

// createConsensus creates a new Raft node and conditionally bootstraps the cluster
func createConsensus(logger io.Writer, dataPath, address string, port int, members []string) (*raft.Raft, chan bool, error) {
	raftId := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%d", address, port))))
	raftAddr := address + ":" + strconv.Itoa(port)
	raftNotifyCh := make(chan bool, 1)

	dataDir := filepath.Join(dataPath, raftId)

	if err := ensureDataPath(dataDir); err != nil {
		return nil, nil, err
	}

	raftDbPath := filepath.Join(dataDir, "raft.db")
	raftDb, err := raftboltdb.NewBoltStore(raftDbPath)
	if err != nil {
		return nil, nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStore(dataDir, 1, os.Stdout)
	if err != nil {
		return nil, nil, err
	}

	transport, err := raft.NewTCPTransport(raftAddr, nil, 3, 10*time.Second, os.Stdout)
	if err != nil {
		return nil, nil, err
	}

	raftConfig := raft.DefaultConfig()
	raftConfig.LogOutput = logger
	raftConfig.LocalID = raft.ServerID(raftAddr)
	raftConfig.NotifyCh = raftNotifyCh

	r, err := raft.NewRaft(raftConfig, &fsm{}, raftDb, raftDb, snapshotStore, transport)
	if err != nil {
		return nil, nil, err
	}

	bootstrapConfig := raft.Configuration{
		Servers: []raft.Server{
			{
				Suffrage: raft.Voter,
				ID:       raft.ServerID(raftAddr),
				Address:  raft.ServerAddress(raftAddr),
			},
		},
	}

	for _, member := range members {
		if member == raftAddr {
			continue
		}

		bootstrapConfig.Servers = append(bootstrapConfig.Servers, raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(member),
			Address:  raft.ServerAddress(member),
		})
	}

	f := r.BootstrapCluster(bootstrapConfig)
	if err := f.Error(); err != nil {
		return nil, nil, err
	}

	return r, raftNotifyCh, nil
}

func ensureDataPath(path string) error {
	// Flush the data directory, if it exists
	if err := os.RemoveAll(path + "/"); err != nil {
		return err
	}

	// Create the data directory, if it doesn't exist
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	return nil
}

// Required by the Raft FSM interface
type fsm struct {
}

func (f *fsm) Apply(*raft.Log) interface{} {
	return nil
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (f *fsm) Restore(io.ReadCloser) error {
	return nil
}
