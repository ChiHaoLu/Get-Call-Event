# Get Call Event

Survey how to get event with ETH RPC in Golang when usgin estimate method, e.g. `eth_call`, `debug_XXX`...

It means the transaction won't be mined, it even can be reverted!

## Environment
- Go: 1.21.0
- The testing contract has only deployed on Sepolia, make sure giving `RPC_URL` in correct network in `.env`.
- The Infura doesn't supoort, and Alchemy only supports `debug_traceCall` for Growth, Scale, or Enterprise, not available on the Free tier. ([ref.](https://docs.alchemy.com/reference/debug-api-quickstart))

## Run

```
$ go mod download
$ go run *.go
```

## Result

| ETH RPC_CALL | Usage | Result | Reason |
|---|---|---|---|
|`eth_call`| For read-only functions, it returns what the read-only function returns. For functions that change the state of the contract, it executes that transaction locally and returns any data returned by the function. ([ref.](https://docs.alchemy.com/reference/eth-call))| ❌ | The issue is that eth_call is specced to return only a binary blob and we can't add the logs in there without breaking the API. ([ref.](https://github.com/ethereum/go-ethereum/issues/20694#issuecomment-677457387)) |
|`debug_traceCall`| Lets you run an `eth_call` on top of a given block. The block can be specified either by hash or by number. It takes the same input object as a `eth_call`.
It returns the same output as debug_traceTransaction. A tracer can be specified as a third argument, similar to `debug_traceTransaction`.([ref.](https://github.com/ethereum/go-ethereum/pull/21338)) | ⚠️||

## Reference

- [Ethereum Development with Go](https://goethereumbook.org/en/)
- [Ethereum JSON-RPC Specification](https://ethereum.github.io/execution-apis/api-documentation/)
- [Alchemy: debug_traceCall](https://docs.alchemy.com/reference/debug-tracecall)
- [Skaled returns error on valid debug_traceCall request #1748](https://github.com/skalenetwork/skaled/issues/1748)