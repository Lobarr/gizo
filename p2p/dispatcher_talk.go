package p2p

import (
	"github.com/gizo-network/gizo/core"
	"github.com/gammazero/nexus/wamp"
)

func (d Dispatcher) BlockReq(ctx context.Context. args wamp.List, kwargs, details wamp.Dict) *client.InvokeResult {
	blockinfo, _ := d.GetBC().GetBlockInfo(args[0].(string))
	block := blockinfo.GetBlock()
	return &client.InvokeResult{Args: wamp.List{block}}
}

func (d Dispatcher) BlockReqMult(ctx context.Context. args wamp.List, kwargs, details wamp.Dict) *client.InvokeResult {
	var blocks []core.Block
	blockinfo, _ := d.GetBC().GetBlockInfo(args[0].(string))
	for _, hash := range args {
		blockinfo, _ := d.GetBC().GetBlockInfo(args[0].(string))
		block := blockinfo.GetBlock()
	}
}