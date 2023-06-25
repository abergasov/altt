## Overview
init after clone
### Usage
```
# check tests
make test
# check linter
make lint
# up in docker
make dev_up
```

than simply open in browser
http://127.0.0.1:8000/eth/balance/0xDf8ac28156209F5cbc89cDec419dd3dB3D3E326a

or for track specific erc20 token
http://127.0.0.1:8000/eth/usdc/balance/0xDf8ac28156209F5cbc89cDec419dd3dB3D3E326a

metrics are available on, proxy metrics are with prefix `balancer_proxy`
http://127.0.0.1:8000/metrics

helt check and available endpoints
http://127.0.0.1:8000/health
http://127.0.0.1:8000/ready

as there are zero external deps, they always serve `ok`

### Implementation details
main logic is in `internal/service/web3/balancer` package.

for rotate rpc nodes, there is `internal/service/rpc` with allow use multiple free rpc nodes and avoid limitation.

solution can be improved by caching known addresses and track changes from new transaction.
