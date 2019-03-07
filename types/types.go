package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// SendTxArgs defines type of first argument of method `eth_sendTransaction`
type SendTxArgs struct {
	From        common.Address  `json:"from"`
	To          *common.Address `json:"to"`
	Gas         *hexutil.Big    `json:"gas"`
	GasPrice    *hexutil.Big    `json:"gasPrice"`
	Value       *hexutil.Big    `json:"value"`
	Data        hexutil.Bytes   `json:"data"`
	Nonce       *hexutil.Uint64 `json:"nonce"`
	PrivateFrom string          `json:"privateFrom"`
	PrivateFor  []string        `json:"privateFor"`
}
