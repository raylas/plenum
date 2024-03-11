package plenum

import (
	"strconv"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
)

// Procede starts the core consensus loop
func (p *Plenum) Procede() {
	done := false
	time.Sleep(1 * time.Second) // Accomodates initial leader channel eventing

	for !done {
		select {
		case <-p.IsClosed:
			done = true

		case <-p.IsLeader:
			p.logger.Println("member is leader")

		case ev := <-p.Procedures:
			leader := p.consensus.VerifyLeader()

			if memberEvent, ok := ev.(serf.MemberEvent); ok {
				for _, member := range memberEvent.Members {
					changedPeer := member.Addr.String() + ":" + strconv.Itoa(int(member.Port+1))

					switch memberEvent.EventType() {
					case serf.EventMemberJoin:
						if leader.Error() == nil {
							f := p.consensus.AddVoter(raft.ServerID(changedPeer), raft.ServerAddress(changedPeer), 0, 0)
							if err := f.Error(); err != nil {
								p.logger.Fatalf("error adding voter to consensus: %s", err)
							}
						}
						p.logger.Printf("member joined: %s", changedPeer)

					case serf.EventMemberLeave, serf.EventMemberFailed:
						if leader.Error() == nil {
							f := p.consensus.RemoveServer(raft.ServerID(changedPeer), 0, 0)
							if err := f.Error(); err != nil {
								p.logger.Fatalf("error removing voter from consensus: %s", err)
							}
						}
						p.logger.Printf("member left: %s", changedPeer)

					case serf.EventMemberReap:
						p.logger.Printf("member reaped: %s", changedPeer)
					}
				}
			}
		}
	}
}
