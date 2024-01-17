package main

import (
	"os"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/rpc"
)

func main() {
    client, err := rpc.DialHTTP(os.Getenv("RPC_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    type request struct {
        To   string `json:"to"`
        Data string `json:"data"`
    }

    var result string

    req := request{"0x1A37E0A92f6F2E06088607B5D87DfeeB95A4BEC2", "0x8da5cb5b"}
    if err := client.Call(&result, "eth_call", req, "latest"); err != nil {
        log.Fatal(err)
    }

    owner := common.HexToAddress(result)
    fmt.Printf("%s\n", owner.Hex())
}