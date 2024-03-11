package plenum

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
)

// Membership returns the current membership list
func (p *Plenum) Membership() []string {
	serfMembers := p.membership.Members()
	members := make([]string, 0, len(serfMembers))
	for _, member := range serfMembers {
		members = append(members, member.Name)
	}
	return members
}

// createMembership creates a new Serf instance and joins the cluster
func createMembership(logger *log.Logger, address string, port int, members []string) (*serf.Serf, chan serf.Event, error) {
	procedures := make(chan serf.Event, 16)

	memberlistConfig := memberlist.DefaultLANConfig()
	memberlistConfig.Logger = logger
	memberlistConfig.BindAddr = address
	memberlistConfig.BindPort = port

	serfConfig := serf.DefaultConfig()
	serfConfig.Logger = logger
	serfConfig.NodeName = fmt.Sprintf("%s:%d", address, port)
	serfConfig.MemberlistConfig = memberlistConfig
	serfConfig.TombstoneTimeout = 1 * time.Minute // TODO: Make this configurable
	serfConfig.EventCh = procedures

	s, err := serf.Create(serfConfig)
	if err != nil {
		return nil, nil, err
	}

	if len(members) > 0 {
		_, err = s.Join(members, false)
		if err != nil {
			log.Fatal(fmt.Errorf("error joining serf: %s", err.Error()))
		}
	}

	return s, procedures, nil
}
