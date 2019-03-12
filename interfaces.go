package quorumgo

import (
	"context"

	"github.com/bsostech/quorumgo/types"
	"github.com/ethereum/go-ethereum/common"
)

// TransactionSender wraps transaction sending.
type TransactionSender interface {
	// Implement method `eth_sendTransaction`
	SendTransaction(ctx context.Context, tx *types.SendTxArgs) (common.Hash, error)
	// Implement method `personal_signAndSendTransaction`
	SignAndSendTransaction(ctx context.Context, tx *types.SendTxArgs, passphrasse string) (common.Hash, error)
}

// IstanbulReader provides access to istanbul.
type IstanbulReader interface {
	// Implement method `istanbul_getSnapshot`
	// Implement method `istanbul_getSnapshotAtHash`
	// Implement method `istanbul_getValidators`
	// Implement method `istanbul_getValidatorsAtHash`
	// Implement method `istanbul_candidates`
}

// IstanbulSender wraps methods related to istanbul.
type IstanbulSender interface {
	// Implement method `istanbul_propose`
	// Implement method `istanbul_discard`
}

// RaftReader provides access to raft.
type RaftReader interface {
	// Implement method `raft_addPeer`.
	RaftAddPeer(enodeID string) (uint16, error)
	// Implement method `raft_removePeer`.
}

// RaftSender wraps methods related to raft.
type RaftSender interface {
	// Implement method `raft_role`.
	// Implement method `raft_leader`.
	// Implement method `raft_cluster`.
}
