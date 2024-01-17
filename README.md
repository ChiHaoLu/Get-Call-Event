# Get Call Event

Get event with ETH RPC in Golang when usgin estimate method, e.g. `eth_call`, `debug_getRawTransaction`...

## Environment
- Go: 1.21.0
- The testing contract has only deployed on Sepolia, make sure giving correct `RPC_URL` in `.env`.

## Run

```
$ go mod download
$ go run *.go
```

## Reference

- [Ethereum Development with Go](https://goethereumbook.org/en/)
- [Ethereum JSON-RPC Specification](https://ethereum.github.io/execution-apis/api-documentation/)
- [Alchemy: eth_call - Ethereum](https://docs.alchemy.com/reference/eth-call)