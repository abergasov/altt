package entities

import (
	"altt/internal/utils"
	"errors"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Chain uint64

type Token string

const (
	ETH     Token = "ETH"
	MATIC   Token = "MATIC"
	USDC    Token = "USDC"
	BTC_b   Token = "BTC_b" // nolint:stylecheck
	USDT    Token = "USDT"
	DAI     Token = "DAI"
	FTM     Token = "FTM"
	AVAX    Token = "AVAX"
	BNB     Token = "BNB"
	AgEUR   Token = "AgEUR"
	LZAgEUR Token = "LZAgEUR"
	One     Token = "ONE"
	STG     Token = "STG"
	CELO    Token = "CELO"
	XDAI    Token = "XDAI"

	ChainEthereum  Chain = 1
	ChainOptimism  Chain = 10
	ChainPolygon   Chain = 137
	ChainFantom    Chain = 250
	ChainArbitrum  Chain = 42161
	ChainAvalanche Chain = 43114
	ChainBNB       Chain = 56
	ChainCoreDAO   Chain = 1116
	ChainHarmony   Chain = 1666600000
	ChainGnosis    Chain = 100
	ChainCelo      Chain = 42220
)

func (c Chain) String() string {
	names := map[Chain]string{
		ChainEthereum:  "eth",
		ChainOptimism:  "optimism",
		ChainPolygon:   "polygon",
		ChainFantom:    "fantom",
		ChainArbitrum:  "arbitrum",
		ChainAvalanche: "avalanche",
		ChainBNB:       "binance sc",
		ChainCoreDAO:   "coredao",
		ChainHarmony:   "harmony",
		ChainGnosis:    "gnosis",
		ChainCelo:      "celo",
	}
	if name, ok := names[c]; ok {
		return name
	}
	return "unknown"
}

func ChainFromString(src string) (Chain, error) {
	names := map[string]Chain{
		"eth":        ChainEthereum,
		"optimism":   ChainOptimism,
		"polygon":    ChainPolygon,
		"fantom":     ChainFantom,
		"arbitrum":   ChainArbitrum,
		"avalanche":  ChainAvalanche,
		"binance sc": ChainBNB,
		"coredao":    ChainCoreDAO,
		"harmony":    ChainHarmony,
		"gnosis":     ChainGnosis,
		"celo":       ChainCelo,
	}
	if chain, ok := names[strings.ToLower(src)]; ok {
		return chain, nil
	}
	return 0, ErrUnknownChain
}

func (c Chain) ChainIDBI() *big.Int {
	return big.NewInt(int64(c))
}

var (
	KnownCoins      = []Token{ETH, USDC, USDT, DAI, FTM, MATIC, AVAX, BTC_b, One, STG, CELO, XDAI, AgEUR}
	ErrUnknownToken = errors.New("unknown coin")
	ErrUnknownChain = errors.New("unknown chain")
)

var Tokens = map[Token]map[Chain]common.Address{
	STG: {
		ChainEthereum:  common.HexToAddress("0xAf5191B0De278C7286d6C7CC6ab6BB8A73bA2Cd6"),
		ChainArbitrum:  common.HexToAddress("0x6694340fc020c5E6B96567843da2df01b2CE1eb6"),
		ChainPolygon:   common.HexToAddress("0x2F6F07CDcf3588944Bf4C42aC74ff24bF56e7590"),
		ChainAvalanche: common.HexToAddress("0x2F6F07CDcf3588944Bf4C42aC74ff24bF56e7590"),
		ChainOptimism:  common.HexToAddress("0x296F55F8Fb28E498B858d0BcDA06D955B2Cb3f97"),
		ChainBNB:       common.HexToAddress("0xB0D502E938ed5f4df2E681fE6E419ff29631d62b"),
		ChainFantom:    common.HexToAddress("0x2F6F07CDcf3588944Bf4C42aC74ff24bF56e7590"),
	},
	ETH: {
		ChainEthereum: common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"),
	},
	AgEUR: {
		ChainArbitrum:  common.HexToAddress("0xFA5Ed56A203466CbBC2430a43c66b9D8723528E7"),
		ChainPolygon:   common.HexToAddress("0xE0B52e49357Fd4DAf2c15e02058DCE6BC0057db4"),
		ChainAvalanche: common.HexToAddress("0xAEC8318a9a59bAEb39861d10ff6C7f7bf1F96C57"),
		ChainGnosis:    common.HexToAddress("0x4b1E2c2762667331Bc91648052F646d1b0d35984"),
		ChainCelo:      common.HexToAddress("0xC16B81Af351BA9e64C1a069E3Ab18c244A1E3049"),
	},
	LZAgEUR: {
		ChainGnosis:   common.HexToAddress("0xFA5Ed56A203466CbBC2430a43c66b9D8723528E7"),
		ChainCelo:     common.HexToAddress("0xf1dDcACA7D17f8030Ab2eb54f2D9811365EFe123"),
		ChainOptimism: common.HexToAddress("0x840b25c87B626a259CA5AC32124fA752F0230a72"),
	},
	BTC_b: {
		ChainPolygon:   common.HexToAddress("0x2297aebd383787a160dd0d9f71508148769342e3"),
		ChainArbitrum:  common.HexToAddress("0x2297aEbD383787A160DD0d9F71508148769342E3"),
		ChainOptimism:  common.HexToAddress("0x2297aEbD383787A160DD0d9F71508148769342E3"),
		ChainAvalanche: common.HexToAddress("0x152b9d0FdC40C096757F570A51E494bd4b943E50"),
	},
	USDC: {
		ChainEthereum:  common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
		ChainPolygon:   common.HexToAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174"),
		ChainArbitrum:  common.HexToAddress("0xFF970A61A04b1cA14834A43f5dE4533eBDDB5CC8"),
		ChainOptimism:  common.HexToAddress("0x7F5c764cBc14f9669B88837ca1490cCa17c31607"),
		ChainFantom:    common.HexToAddress("0x04068DA6C83AFCFA0e13ba15A6696662335D5B75"),
		ChainAvalanche: common.HexToAddress("0xB97EF9Ef8734C71904D8002F8b6Bc66Dd9c48a6E"),
		ChainCoreDAO:   common.HexToAddress("0xa4151B2B3e269645181dCcF2D426cE75fcbDeca9"),
		ChainHarmony:   common.HexToAddress("0x9b5fae311A4A4b9d838f301C9c27b55d19BAa4Fb"),
	},
	DAI: {
		ChainEthereum:  common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"),
		ChainPolygon:   common.HexToAddress("0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063"),
		ChainArbitrum:  common.HexToAddress("0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1"),
		ChainOptimism:  common.HexToAddress("0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1"),
		ChainFantom:    common.HexToAddress("0x8D11eC38a3EB5E956B052f67Da8Bdc9bef8Abf3E"),
		ChainAvalanche: common.HexToAddress("0xd586E7F844cEa2F87f50152665BCbc2C279D8d70"),
	},
	USDT: {
		ChainEthereum:  common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
		ChainPolygon:   common.HexToAddress("0xc2132D05D31c914a87C6611C10748AEb04B58e8F"),
		ChainArbitrum:  common.HexToAddress("0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9"),
		ChainOptimism:  common.HexToAddress("0x94b008aA00579c1307B0EF2c499aD98a8ce58e58"),
		ChainFantom:    common.HexToAddress("0x1B27A9dE6a775F98aaA5B90B62a4e2A0B84DbDd9"),
		ChainAvalanche: common.HexToAddress("0x9702230A8Ea53601f5cD2dc00fDBc13d4dF4A8c7"),
	},
}

func GetTokenAddress(chainID Chain, token Token) (common.Address, error) {
	if _, ok := Tokens[token]; !ok {
		return common.Address{}, ErrUnknownToken
	}
	if _, ok := Tokens[token][chainID]; !ok {
		return common.Address{}, ErrUnknownChain
	}
	return Tokens[token][chainID], nil
}

func GetChain(chainID *big.Int) Chain {
	data := map[uint64]Chain{
		1:          ChainEthereum,
		10:         ChainOptimism,
		56:         ChainBNB,
		137:        ChainPolygon,
		250:        ChainFantom,
		42161:      ChainArbitrum,
		43114:      ChainAvalanche,
		1116:       ChainCoreDAO,
		1666600000: ChainHarmony,
		100:        ChainGnosis,
		42220:      ChainCelo,
	}
	if chain, ok := data[chainID.Uint64()]; ok {
		return chain
	}
	return 0
}

func TokenFromString(data string) (Token, error) {
	dataMap := map[string]Token{
		"ETH":   ETH,
		"USDC":  USDC,
		"DAI":   DAI,
		"USDT":  USDT,
		"MATIC": MATIC,
		"FTM":   FTM,
		"AVAX":  AVAX,
		"BNB":   BNB,
		"BTC.b": BTC_b,
		"BTC":   BTC_b,
		"ONE":   One,
		"STG":   STG,
		"CELO":  CELO,
		"XDAI":  XDAI,
	}
	if coin, ok := dataMap[strings.ToUpper(data)]; ok {
		return coin, nil
	}
	return "", ErrUnknownToken
}

func MapChainToFuel(chain Chain) Token {
	switch chain {
	case ChainArbitrum, ChainOptimism, ChainEthereum:
		return ETH
	case ChainPolygon:
		return MATIC
	case ChainFantom:
		return FTM
	case ChainAvalanche:
		return AVAX
	case ChainHarmony:
		return One
	case ChainBNB:
		return BNB
	case ChainCelo:
		return CELO
	case ChainGnosis:
		return XDAI
	}
	log.Fatal("no fuel for chain", chain)
	return ETH
}

func GetTokenDecimals(token Token) int {
	switch token {
	case MATIC, ETH, FTM, DAI, AVAX, BNB, AgEUR, One, STG, CELO, XDAI:
		return 18
	case BTC_b:
		return 8
	case USDC, USDT:
		return 6
	}
	log.Fatal("no decimals for token", token)
	return 0
}

func CoinFromWEI(token Token, wei *big.Int) string {
	decimals := GetTokenDecimals(token)
	return utils.CustomFromWei(wei, decimals)
}
