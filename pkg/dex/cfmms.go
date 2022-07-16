package dex

import (
	"math"

	"github.com/ethereum/go-ethereum/common"
	"gonum.org/v1/gonum/mat"
)

const (
	Prod int = 1
	Geom     = 2
)

type CFMM interface {
	FindArb([]float64)
	arbDelta(float64, float64, float64, float64) float64
	arbLambda(float64, float64, float64, float64) float64
	Type() int
}

type ProductTwoCoin struct {
	TokensTraded []common.Address
	Fee          float64
	Deltas       []float64
	Lambdas      []float64
	Reserves     []float64

	// Matrix
	LocalToGlobal *mat.Dense
}

type GeometricMeanTwoCoin struct {
	Address      common.Address
	TokensTraded []common.Address
	Fee          float64
	Deltas       []float64
	Lambdas      []float64
	Reserves     []float64
	Weights      []float64

	// Matrix
	LocalToGlobal *mat.Dense
}

// Interface function definitions for ProductTwoCoin

func (market ProductTwoCoin) FindArb(nu []float64) {

	k := market.Reserves[0] * market.Reserves[1]
	market.Deltas[0] = market.arbDelta(k, nu[1]/nu[0], market.Reserves[0], 0)
	market.Deltas[1] = market.arbDelta(k, nu[0]/nu[1], market.Reserves[1], 0)

	market.Lambdas[0] = market.arbLambda(k, nu[0]/nu[1], market.Reserves[0], 0)
	market.Lambdas[1] = market.arbLambda(k, nu[1]/nu[0], market.Reserves[1], 0)

}

func (market ProductTwoCoin) arbDelta(k float64, m float64, r float64, ignore float64) float64 { // ignore is a quick hack to make the interfaces work -- you can ignore

	val := math.Sqrt(market.Fee*m*k) / market.Fee
	if val >= 0 {
		return val
	}
	return 0

}

func (market ProductTwoCoin) arbLambda(k float64, m float64, r float64, ignore float64) float64 { // ignore is a quick hack to make the interfaces work -- you can ignore

	val := r - math.Sqrt(k/(m*market.Fee))
	if val >= 0 {
		return val
	}
	return 0

}

func (market ProductTwoCoin) Type() int {
	return Prod
}

// Interface function definitions for GeometricMeanTwoCoin\

func (market GeometricMeanTwoCoin) FindArb(nu []float64) {
	eta := market.Weights[0] / market.Weights[1]

	market.Deltas[0] = market.arbDelta(eta, nu[1]/nu[0], market.Reserves[1], market.Reserves[0])
	market.Deltas[1] = market.arbDelta(1/eta, nu[0]/nu[1], market.Reserves[0], market.Reserves[1])

	market.Lambdas[0] = market.arbLambda(1/eta, nu[0]/nu[1], market.Reserves[0], market.Reserves[1])
	market.Lambdas[1] = market.arbLambda(eta, nu[1]/nu[0], market.Reserves[1], market.Reserves[0])

}

func (market GeometricMeanTwoCoin) arbDelta(eta float64, m float64, r0 float64, r1 float64) float64 {
	// 	@inline geom_arb_δ(m, r1, r2, η, γ) = max((γ*m*η*r1*r2^η)^(1/(η+1)) - r2, 0)/γ

	val := (math.Pow(market.Fee*m*eta*r0*math.Pow(r1, eta), 1/(eta+1)) - r1) / market.Fee
	if val >= 0 {
		return val
	}
	return 0

}

func (market GeometricMeanTwoCoin) arbLambda(eta float64, m float64, r0 float64, r1 float64) float64 {
	// @inline geom_arb_λ(m, r1, r2, η, γ) = max(r1 - ((r2*r1^(1/η))/(η*γ*m))^(η/(1+η)), 0)

	val := r0 - math.Pow(((math.Pow(r1*r0, (1/eta)))/(eta*market.Fee*m)), (eta/(1+eta)))
	if val >= 0 {
		return val
	}
	return 0

}

func (market GeometricMeanTwoCoin) Type() int {
	return Geom
}
