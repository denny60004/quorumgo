package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	a "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/bsostech/quorumgo/quorumclient"
)

var address1 = common.HexToAddress("ed9d02e382b34818e88b88a309c7fe71e65f419d")
var address2 = common.HexToAddress("ca843569e3427144cead5e4d5999a3d0ccf92b8e")
var address3 = common.HexToAddress("0fbdc686b912d7722dc86510934589e0aaf3b55a")
var address4 = common.HexToAddress("9186eb3d20cbd1f5f992a950d808c4495153abd5")

var cruxPub1 = `BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=`
var cruxPub2 = `QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc=`
var cruxPub3 = `1iTZde/ndBHvzhcl7V68x44Vx7pl8nwx9LqnM/AfJUg=`
var cruxPub4 = `oNspPPgszVUFw0qmGFfWwh1uxVUXgvBxleXORHj07g8=`

var abi, _ = a.JSON(strings.NewReader(`[{"constant":true,"inputs":[],"name":"storedData","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"retVal","type":"uint256"}],"payable":false,"type":"function"},{"inputs":[{"name":"initVal","type":"uint256"}],"payable":false,"type":"constructor"}]`))
var bin, _ = hex.DecodeString(`6060604052341561000f57600080fd5b604051602080610149833981016040528080519060200190919050505b806000819055505b505b610104806100456000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632a1afcd914605157806360fe47b11460775780636d4ce63c146097575b600080fd5b3415605b57600080fd5b606160bd565b6040518082815260200191505060405180910390f35b3415608157600080fd5b6095600480803590602001909190505060c3565b005b341560a157600080fd5b60a760ce565b6040518082815260200191505060405180910390f35b60005481565b806000819055505b50565b6000805490505b905600a165627a7a72305820d5851baab720bba574474de3d09dbeaabc674a15f4dd93b974908476542c23f00029`)

func main() {
	ec, err := quorumclient.Dial("http://localhost:22001")
	logErr(err)

	args, _ := abi.Constructor.Inputs.Pack(big.NewInt(42))
	cData := append(bin, args...)

	gas, err := ec.EstimateGas(context.TODO(), ethereum.CallMsg{From: address1, Data: cData})
	logErr(err)
	privateFor := []string{cruxPub2}
	tx := quorumclient.NewSendTxArgs(address1, nil, hexutil.Big(*big.NewInt(int64(gas))), cData, privateFor)
	passphrase := ""
	txHash, err := ec.SignAndSendTransaction(context.TODO(), tx, passphrase)
	logErr(err)

	rx := make(chan *types.Receipt)
	go func() {
		txHashStr := txHash.Hex()
		for {
			log.Println("wating for tx " + txHashStr + " mined...")
			receipt, _ := ec.TransactionReceipt(context.TODO(), txHash)
			if receipt != nil {
				rx <- receipt
				break
			} else {
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	receipt := <-rx
	log.Println("contract address: " + receipt.ContractAddress.Hex())

	getData, _ := abi.Pack("get")
	x, err := ec.CallContract(context.TODO(), ethereum.CallMsg{From: address1, To: &receipt.ContractAddress, Data: getData}, nil)
	logErr(err)

	log.Println(hex.EncodeToString(x))
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
