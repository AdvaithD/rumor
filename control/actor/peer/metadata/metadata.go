package metadata

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/protolambda/ask"
	"github.com/protolambda/rumor/control/actor/base"
	"github.com/protolambda/rumor/p2p/rpc/methods"
	"github.com/protolambda/rumor/p2p/rpc/reqresp"
	"github.com/protolambda/rumor/p2p/track"
)

type PeerMetadataState struct {
	Following bool
	Local     methods.MetaData
}

type PeerMetadataCmd struct {
	*base.Base
	*PeerMetadataState
	Store track.ExtendedPeerstore
}

func (c *PeerMetadataCmd) Help() string {
	return "Manage and track peer metadata"
}

func (c *PeerMetadataCmd) Cmd(route string) (cmd interface{}, err error) {
	switch route {
	case "ping":
		cmd = &PeerMetadataPingCmd{Base: c.Base, PeerMetadataState: c.PeerMetadataState, Store: c.Store}
	case "pong":
		cmd = &PeerMetadataPongCmd{Base: c.Base, PeerMetadataState: c.PeerMetadataState, Book: c.Store}
	case "get":
		cmd = &PeerMetadataGetCmd{Base: c.Base, PeerMetadataState: c.PeerMetadataState}
	case "set":
		cmd = &PeerMetadataSetCmd{Base: c.Base, PeerMetadataState: c.PeerMetadataState}
	case "req":
		cmd = &PeerMetadataReqCmd{Base: c.Base, PeerMetadataState: c.PeerMetadataState, Book: c.Store}
	case "poll":
		cmd = &PeerMetadataPollCmd{Base: c.Base, PeerMetadataState: c.PeerMetadataState, Store: c.Store}
	case "serve":
		cmd = &PeerMetadataServeCmd{Base: c.Base, PeerMetadataState: c.PeerMetadataState}
	case "follow":
		cmd = &PeerMetadataFollowCmd{Base: c.Base, PeerMetadataState: c.PeerMetadataState}
	default:
		return nil, ask.UnrecognizedErr
	}
	return cmd, nil
}

func (c *PeerMetadataCmd) Routes() []string {
	return []string{"ping", "pong", "get", "set", "req", "poll", "serve", "follow"}
}

func (c *PeerMetadataState) fetch(book track.MetadataBook, sFn reqresp.NewStreamFn, ctx context.Context, peerID peer.ID, comp reqresp.Compression) (
	resCode reqresp.ResponseCode, errMsg string, data *methods.MetaData, err error) {

	err = methods.MetaDataRPCv1.RunRequest(ctx, sFn, peerID, comp, reqresp.RequestSSZInput{Obj: nil}, 1,
		func() error {
			// TODO
			return nil
		},
		func(chunk reqresp.ChunkedResponseHandler) error {
			resCode = chunk.ResultCode()
			switch resCode {
			case reqresp.ServerErrCode, reqresp.InvalidReqCode:
				msg, err := chunk.ReadErrMsg()
				if err != nil {
					return err
				}
				errMsg = msg
			case reqresp.SuccessCode:
				var meta methods.MetaData
				if err := chunk.ReadObj(&meta); err != nil {
					return err
				}
				data = &meta
				book.RegisterMetadata(peerID, meta)
			}
			return nil
		})
	return
}

func (c *PeerMetadataState) ping(sFn reqresp.NewStreamFn, ctx context.Context, peerID peer.ID, comp reqresp.Compression) (
	resCode reqresp.ResponseCode, errMsg string, data methods.Pong, err error) {

	p := methods.Ping(c.Local.SeqNumber)
	err = methods.PingRPCv1.RunRequest(ctx, sFn, peerID, comp, reqresp.RequestSSZInput{Obj: &p}, 1,
		func() error {
			// TODO
			return nil
		},
		func(chunk reqresp.ChunkedResponseHandler) error {
			resCode = chunk.ResultCode()
			switch resCode {
			case reqresp.ServerErrCode, reqresp.InvalidReqCode:
				msg, err := chunk.ReadErrMsg()
				if err != nil {
					return err
				}
				errMsg = msg
			case reqresp.SuccessCode:
				var pong methods.Pong
				if err := chunk.ReadObj(&pong); err != nil {
					return err
				}
				data = pong
			}
			return nil
		})
	return
}
