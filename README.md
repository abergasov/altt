## Overview
init after clone
## Usage
```
# check tests
make test
# check linter
make lint
# up in docker
make dev_up
```

than simply open in browser
http://127.0.0.1:8012/eth/balance/0xDf8ac28156209F5cbc89cDec419dd3dB3D3E326a

or for track specific erc20 token
http://127.0.0.1:8012/eth/usdc/balance/0xDf8ac28156209F5cbc89cDec419dd3dB3D3E326a

metrics are available on, proxy metrics are with prefix `balancer_proxy`
http://127.0.0.1:8012/metrics

helt check and available endpoints
http://127.0.0.1:8012/health
http://127.0.0.1:8012/ready

as there are zero external deps, they always serve `ok`