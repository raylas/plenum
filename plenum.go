package plenum

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
)

type (
	Plenum struct {
		logger *log.Logger

		membership *serf.Serf
		consensus  *raft.Raft

		Procedures chan serf.Event // Exposed channel for handling Serf events

		IsLeader chan bool      // Exposed channel for determining leadership
		IsClosed chan os.Signal // Exposed channel for graceful shutdown
	}
)

// Convene creates a new membership and consensus group
func Convene(
	logger *log.Logger,
	address string,
	port int,
	members []string,
) (*Plenum, error) {
	pLogger, sLogger, rLogger := setLogging(logger)

	membership, procedures, err := createMembership(sLogger, address, port, members)
	if err != nil {
		return nil, err
	}

	consensus, isLeaderCh, err := createConsensus(rLogger, os.TempDir(), address, port+1, members)
	if err != nil {
		return nil, err
	}

	isClosed := make(chan os.Signal, 1)
	signal.Notify(isClosed, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-isClosed

		pLogger.Println("member exiting the plenum")
		membership.Leave()
		os.Exit(1)
	}()

	return &Plenum{
		logger:     pLogger,
		membership: membership,
		consensus:  consensus,
		Procedures: procedures,
		IsLeader:   isLeaderCh,
		IsClosed:   isClosed,
	}, nil
}

func setLogging(logger *log.Logger) (*log.Logger, *log.Logger, io.Writer) {
	// Plenum core logger
	var pLogger *log.Logger
	if logger != nil {
		pLogger = logger
	} else {
		pLogger = log.New(os.Stdout, "", log.LstdFlags)
	}

	// Serf and Raft loggers
	var sLogger *log.Logger
	var rLogger io.Writer
	_, ok := os.LookupEnv("PLENUM_DEBUG")
	if ok {
		sLogger = pLogger
		rLogger = pLogger.Writer()
	} else {
		sLogger = log.New(io.Discard, "", 0)
		rLogger = io.Discard
	}

	return pLogger, sLogger, rLogger
}
