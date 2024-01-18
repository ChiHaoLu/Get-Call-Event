package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ChiHaoLu/Get-Call-Event/tracer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

type request struct {
	From common.Address `json:"from"`
	To   string         `json:"to"`
	Data string         `json:"data"`
}

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
	fnSig, data := ConstructTxData()
	from, signedTx := GetSignedTx(data)
	fmt.Println(signedTx.Hash().Hex())

	// Try callTracer
	var callTracerResult tracer.CallTracerResult
	req := request{
		From: from,
		To:   os.Getenv("CONTRACT_ADDR"),
		Data: fnSig,
	}
	callTracerConfig := tracer.TraceConfig{
		Tracer: "callTracer",
	}
	if err := client.Call(&callTracerResult, "debug_traceCall", req, "latest", callTracerConfig); err != nil {
		log.Fatal(err)
	}
	fmt.Println(callTracerResult)

	// Try custom tracer
	var eventTracerResult tracer.EventTracerResult
	eventTracerConfig := tracer.TraceConfig{
		Tracer: tracer.Loaded.EventTracer,
	}
	fmt.Println(eventTracerConfig.Tracer)
	if err := client.Call(&eventTracerResult, "debug_traceCall", req, "latest", eventTracerConfig); err != nil {
		log.Fatal(err)
	}
	fmt.Println(eventTracerResult)
}
