package peer

import (
	"github.com/protolambda/ask"
	"github.com/protolambda/rumor/control/actor/base"
	"github.com/protolambda/rumor/control/actor/peer/metadata"
	"github.com/protolambda/rumor/control/actor/peer/status"
)

type PeerCmd struct {
	*base.Base
	*status.PeerStatusState
	*metadata.PeerMetadataState
}

func (c *PeerCmd) Cmd(route string) (cmd interface{}, err error) {
	switch route {
	case "connect":
		cmd = &PeerConnectCmd{Base: c.Base}
	case "disconnect":
		cmd = &PeerDisconnectCmd{Base: c.Base}
	case "protect":
		cmd = &PeerProtectCmd{Base: c.Base}
	case "unprotect":
		cmd = &PeerUnprotectCmd{Base: c.Base}
	case "trim":
		cmd = &PeerTrimCmd{Base: c.Base}
	case "list":
		cmd = &PeerListCmd{Base: c.Base}
	case "addrs":
		cmd = &PeerAddrsCmd{Base: c.Base}
	case "status":
		cmd = &status.PeerStatusCmd{Base: c.Base, PeerStatusState: c.PeerStatusState}
	case "metadata":
		cmd = &metadata.PeerMetadataCmd{Base: c.Base, PeerMetadataState: c.PeerMetadataState}
	default:
		return nil, ask.UnrecognizedErr
	}
	return cmd, nil
}

func (c *PeerCmd) Routes() []string {
	return []string{"connect", "disconnect", "protect", "unprotect", "trim", "list", "addrs", "status", "metadata"}
}

func (c *PeerCmd) Help() string {
	return "Manage peers"
}
