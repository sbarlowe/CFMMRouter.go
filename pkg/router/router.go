package router

import (
	"amalfi/pkg/dex"

	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	Exchanges []common.Address
	Endpoint  string

	PrivateKey string

	QueryAddress    common.Address
	ContractAddress common.Address

	BaseToken common.Address
}

type Router struct {

	// dual variable
	Nu []float64

	// market information
	CFMMs     []dex.CFMM
	NumCFMMs  int
	Tokens    []common.Address
	NumTokens int
}

func NewRouter() (*Router, error) {

	// Define Tokens
	token1 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	token2 := common.HexToAddress("0x0000000000000000000000000000000000000002")
	TokensI := []common.Address{token1, token2}

	// Define Pools
	equalPool := dex.ProductTwoCoin{TokensTraded: TokensI, Fee: 1, Reserves: []float64{1e6, 1e6}}
	unequalSmallPool := dex.ProductTwoCoin{TokensTraded: TokensI, Fee: 1, Reserves: []float64{1e3, 2e3}}
	weightedPool := dex.GeometricMeanTwoCoin{TokensTraded: TokensI, Fee: 1, Weights: []float64{0.4, 0.6}, Reserves: []float64{1e4, 2e4}}

	// Define Various Market Information
	CFMMsI := []dex.CFMM{equalPool, unequalSmallPool, weightedPool}

	//create new bot object
	return &Router{
		CFMMs:     CFMMsI,
		NumCFMMs:  len(CFMMsI),
		Tokens:    TokensI,
		NumTokens: len(TokensI),
	}, nil
}
