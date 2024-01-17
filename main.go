package main

import (
	"os"
    "fmt"
    "log"

    "github.com/joho/godotenv"
    "github.com/ethereum/go-ethereum/rpc"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("err loading: %v", err)
    }

	// Setup client
    client, err := rpc.DialHTTP(os.Getenv("RPC_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

	// Get signed trasanction
    data := ConstructTxData()
	signedTx := GetSignedTx(data)
    fmt.Println(signedTx.Hash().Hex())

	// Try eth_call
    type request struct {
        To   string `json:"to"`
        Data string `json:"data"`
    }
    var result string

    req := request{os.Getenv("CONTRACT_ADDR"), "0x8da5cb5b"}
    if err := client.Call(&result, "eth_call", req, "latest"); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s\n", result)

	// TODO: Try debug_xxx
}