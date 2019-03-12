package quorumclient

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/bsostech/quorumgo/types"
)

// Client extends ethclient.Client's and rpc.Client's methods.
type Client struct {
	*ethclient.Client
	rc *rpc.Client
}

// NewClient creates a quorumclient.Client.
func NewClient(rc *rpc.Client) *Client {
	return &Client{
		Client: ethclient.NewClient(rc),
		rc:     rc,
	}
}

// Dial creates a quorumclient.Client connected to given RPC url.
func Dial(url string) (*Client, error) {
	rc, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return NewClient(rc), nil
}

// NewSendTxArgs creates a SendTxArgs.
func NewSendTxArgs(from common.Address, to *common.Address, gas hexutil.Big, data hexutil.Bytes, privateFor []string) *types.SendTxArgs {
	return &types.SendTxArgs{
		From:       from,
		To:         to,
		Gas:        &gas,
		Data:       data,
		PrivateFor: privateFor,
	}
}

// SignAndSendTransaction request an locked account to sign a transaction and submit.
func (ec *Client) SignAndSendTransaction(ctx context.Context, tx *types.SendTxArgs, passphrasse string) (common.Hash, error) {
	var txHash common.Hash
	err := ec.rc.CallContext(ctx, &txHash, "personal_signAndSendTransaction", tx, passphrasse)
	return txHash, err
}

// SendTransaction request an unlocked account to sign a transaction and submit.
func (ec *Client) SendTransaction(ctx context.Context, tx *types.SendTxArgs) (common.Hash, error) {
	var txHash common.Hash
	err := ec.rc.CallContext(ctx, &txHash, "eth_sendTransaction", tx)
	return txHash, err
}

// RaftAddPeer should contain enode, IP(real world), raft port
func (ec *Client) RaftAddPeer(ctx context.Context, enodeID string) (uint16, error) {
	var raftID uint16
	err := ec.rc.CallContext(ctx, &raftID, "raft_addPeer", enodeID)
	return raftID, err
}
