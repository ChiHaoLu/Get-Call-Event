package main

import (
	"context"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ConstructTxData() (string, []byte) {
	methodID := crypto.Keccak256([]byte("EmitWithWrite()"))[:4]
	// fmt.Println(hexutil.Encode(methodID)) // 0x9d460f0c

	var data []byte
	data = append(data, methodID...)
	// If we need to give input params to contract function, we can append here
	return hexutil.Encode(methodID), data
}

func GetSignedTx(data []byte) (common.Address, *types.Transaction) {
	// We use ethClient to get on-chain method, use rpcClient to realize rpc call in low level
	client, err := ethclient.Dial(os.Getenv("RPC_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Setup ECDSA signing
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Construct the transaction Object
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tx := types.NewTransaction(nonce, common.HexToAddress(os.Getenv("CONTRACT_ADDR")), value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	return address, signedTx
}
