# Get Call Event

Survey how to get event with ETH RPC in Golang when usgin estimate method, e.g. `eth_call`, `debug_XXX`...

It means the transaction won't be mined, it even can be reverted!

## Environment
- Go: 1.21.0
- The testing contract has only deployed on Sepolia, make sure giving `RPC_URL` in correct network in `.env`.

## `debug_traceCall` Problems

- The Infura doesn't supoort `debug_traceCall`, and Alchemy only supports `debug_traceCall` for Growth, Scale, or Enterprise, not available on the Free tier. ([ref.](https://docs.alchemy.com/reference/debug-api-quickstart))

```curl
curl https://eth-sepolia.g.alchemy.com/v2/<your_api_key> \
-X POST \
-H "Content-Type: application/json" \
--data '{"method":"debug_traceCall","params":[{"from":null,"to":"0x1A37E0A92f6F2E06088607B5D87DfeeB95A4BEC2","data":"0xc5fde5e5"}, "latest", {"tracer": "callTracer"}],"id":1,"jsonrpc":"2.0"}'
```

- However, A bundler should be considered as an extension, or a side-car, to a full node. It requires full access to the node's debug API for custom JS tracing. Most RPC providers do not support this. For maximum performance, it is also recommended to run the bundler in the same machine as the node. ([ref.](https://docs.stackup.sh/docs/erc-4337-bundler-installation#running-the-bundler))
- In order to implement the full spec storage access rules and opcode banning, it must run against a GETH node, which supports `debug_traceCall` with javascript "tracer" Specifically, `hardhat node` and `ganache` do NOT support this API. You can still run the bundler with such nodes, but with `--unsafe` so it would skip these security checks([ref.](https://github.com/eth-infinitism/bundler/tree/main?tab=readme-ov-file#running-local-node))

Hence, after open the Geth on `http://localhost:8545`, you can:
```curl
curl http://localhost:8545 \
-X POST \
-H "Content-Type: application/json" \
--data '{"method":"debug_traceCall","params":[{"from":null,"to":"0x1A37E0A92f6F2E06088607B5D87DfeeB95A4BEC2","data":"0xc5fde5e5"}, "latest", {"tracer": "callTracer"}],"id":1,"jsonrpc":"2.0"}'
>
{"jsonrpc":"2.0","id":1,"result":{"type":"CALL","from":"0x0000000000000000000000000000000000000000","to":"0x1a37e0a92f6f2e06088607b5d87dfeeb95a4bec2","value":"0x0","gas":"0x2fa9e38","gasUsed":"0x0","input":"0xc5fde5e5","output":"0x"}}
```
Or custom tracer:
```curl
curl http://localhost:8545 \
-X POST \
-H "Content-Type: application/json" \
--data '{"jsonrpc":"2.0","method":"debug_traceCall","params":[{"from":null,"to":"0x1A37E0A92f6F2E06088607B5D87DfeeB95A4BEC2","data":"0xc5fde5e5"}, "latest", {"tracer": "{\"data\":[],\"fault\":function(log){},\"step\":function(log){var topicCount=(log.op.toString().match(/LOG(\\d)/)||[])[1];if(topicCount){var res={address:log.contract.getAddress(),data:log.memory.slice(parseInt(log.stack.peek(0)),parseInt(log.stack.peek(0))+parseInt(log.stack.peek(1))),};for(var i=0;i<topicCount;i++)res[`topic`+i.toString()]=log.stack.peek(i+2);this.data.push(res);}},\"result\":function(){return this.data;}}" }],"id":1}'
>
{"jsonrpc":"2.0","id":1,"result":[]}
```

## Run
Open the Geth with docker in sepolia:
```
docker run --rm -ti --name geth -p 8545:8545 ethereum/client-go:v1.10.26 \
  --miner.gaslimit 12000000 \
  --http --http.api personal,eth,net,web3,debug \
  --http.vhosts '*,localhost,host.docker.internal' --http.addr "0.0.0.0" \
  --ignore-legacy-receipts --allow-insecure-unlock --rpc.allow-unprotected-txs \
  --dev \
  --verbosity 2 \
  --nodiscover --maxpeers 0 --mine --miner.threads 1 \
  --networkid 11155111
```

Fetch the event:
```
$ go mod download
$ go run *.go
```

## Result

| ETH RPC_CALL | Usage | Result | Reason |
|---|---|---|---|
|`eth_call`| For read-only functions, it returns what the read-only function returns. For functions that change the state of the contract, it executes that transaction locally and returns any data returned by the function. ([ref.](https://docs.alchemy.com/reference/eth-call))| ❌ | The issue is that eth_call is specced to return only a binary blob and we can't add the logs in there without breaking the API. ([ref.](https://github.com/ethereum/go-ethereum/issues/20694#issuecomment-677457387)) |
|`debug_traceCall`| Lets you run an `eth_call` on top of a given block. The block can be specified either by hash or by number. It takes the same input object as a `eth_call`. ([ref.](https://github.com/ethereum/go-ethereum/pull/21338)) | ✅  | We can use custom tracer to filter the log message.|
|`debug_traceTransaction`| Attempts to run the transaction in the exact same manner as it was executed on the network. It will replay any transaction that may have been executed prior to this one before it and will then attempt to execute the transaction that corresponds to the given hash. ([ref.](https://docs.alchemy.com/reference/debug-tracetransaction))| ✅ | `debug_traceTransaction` can see the failed transaction's revert message. Hence if we want to estimate and see the event log, we should send this transaction first. ([ref.](https://github.com/ethereum/go-ethereum/issues/25967)) |

## Reference

- [debug Namespace](https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-debug#debugtracecall)
- [Ethereum JSON-RPC Specification](https://ethereum.github.io/execution-apis/api-documentation/)
- [Alchemy: debug_traceCall](https://docs.alchemy.com/reference/debug-tracecall)
- [Skaled returns error on valid debug_traceCall request #1748](https://github.com/skalenetwork/skaled/issues/1748)
- [debug_traceCall](https://hackmd.io/@rajivpoc/debug-tracecall?utm_source=preview-mode&utm_medium=rec)
- [Extracting emitted events (logs) from geth transaction trace (debug_traceCall)](https://stackoverflow.com/questions/72064656/extracting-emitted-events-logs-from-geth-transaction-trace-debug-tracecall)