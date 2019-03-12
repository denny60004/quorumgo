package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"strconv"
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

var enode1 = `5c3c98e3a28a87e73ab40468212de7ab6cf0e2afa77781295925f32369c00baf30f664e52f8d152c02b069d6daa1a61f477e3c1eca64403529dfbd0c31e09524`
var enode2 = `9b98a96a8ba080ff4c7863e5fdf3211a7082b612d5897ae4eed687eec391eb421c8ed7c572ca17f257441a0cb544a7c184244dfdf9a114f5251da3dac72e7585`
var enode3 = `a51690b44ab39fd83c42b5a7c087ba222970951f06655ebbba1625267fad105fd238c9f092e05b2293f526e748b2fa423b22d66296f770037393c26a9e5d3543`
var enode4 = `a68df7cd75e9ea490653bdba7c6868f979944578e59c9efd2aa62878822f16f46a49a13289f6392923053be1acb3a6ec8e2fc92cae59de859fd5892071fbfa88`

var cruxPub1 = `BULeR8JyUWhiuuCMU/HLA0Q5pzkYT+cHII3ZKBey3Bo=`
var cruxPub2 = `QfeDAys9MPDs2XHExtc84jKGHxZg/aj52DTh0vtA3Xc=`
var cruxPub3 = `1iTZde/ndBHvzhcl7V68x44Vx7pl8nwx9LqnM/AfJUg=`
var cruxPub4 = `oNspPPgszVUFw0qmGFfWwh1uxVUXgvBxleXORHj07g8=`

var ip1 = `10.5.0.11`
var ip2 = `10.5.0.12`
var ip3 = `10.5.0.13`
var ip4 = `10.5.0.14`

var raftPort = `50400`
var gethPort = `21000`

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

	enodeID := "enode://" + enode4 + "@" + ip4 + ":" + gethPort + "?discport=0&raftport=" + raftPort
	raftID, err := ec.RaftAddPeer(context.TODO(), enodeID)
	logErr(err)
	log.Println("New raft ID: " + strconv.Itoa(int(raftID)))

	log.Println(hex.EncodeToString(x))
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
