package quorumclient

import "github.com/bsostech/quorumgo"

// Verify that quorumclient.Client implements the quorumgo interfaces.
var (
	_ = quorumgo.TransactionSender(&Client{})
	// _ = quorumgo.IstanbulReader(&Client{})
	// _ = quorumgo.IstanbulSender(&Client{})
	// _ = quorumgo.RaftReader(&Client{})
	// _ = quorumgo.RaftSender(&Client{})
)
