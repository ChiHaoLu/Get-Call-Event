package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

type request struct {
	From *string `json:"from"`
	To   string  `json:"to"`
	Data string  `json:"data"`
}
type traceConfig struct {
	Tracer string `json:"tracer"`
}
type traceResult struct {
	From    string `json:"from"`
	Gas     string `json:"gas"`
	GasUsed string `json:"gasUsed"`
	To      string `json:"to"`
	Input   string `json:"input"`
	Output  string `json:"output"`
	Value   string `json:"value"`
	Type    string `json:"type"`
}

/*
curl https://eth-sepolia.g.alchemy.com/v2/QmC1LbvkXMzQTCqaLT6j0wXp3L_rNJjO \
-X POST \
-H "Content-Type: application/json" \
--data '{"method":"debug_traceCall","params":[{"from":null,"to":"0x1A37E0A92f6F2E06088607B5D87DfeeB95A4BEC2","data":"0xc5fde5e5"}, "latest", {"tracer": "callTracer"}],"id":1,"jsonrpc":"2.0"}'
*/

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
	fxSig, data := ConstructTxData()
	from, signedTx := GetSignedTx(data)
	fmt.Println(from)
	fmt.Println(fxSig)
	fmt.Println(signedTx.Hash().Hex())

	// Try RPC CALL
	var result traceResult
	req := request{nil, os.Getenv("CONTRACT_ADDR"), fxSig}
	config := traceConfig{"callTracer"}
	if err := client.Call(&result, "debug_traceCall", req, "latest", config); err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
