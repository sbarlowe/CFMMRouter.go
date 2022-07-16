package router

import (
	"fmt"
	"math"
	"sync"

	"amalfi/pkg/dex"
	lbfgsb "amalfi/pkg/optim/go-lbfgsb-0.1.6"

	"gonum.org/v1/gonum/mat"
)

func (r *Router) Route() {
	// Optimizer set up
	optimizer := lbfgsb.NewLbfgsb(r.NumTokens)
	_ = optimizer

	// Define bounds
	bounds := [][2]float64{{1 + 1e-8, 1 + 1e-8}, {math.Inf(1), math.Inf(1)}}
	optimizer.SetBounds(bounds) // why does fortran impl require one less bound than julia code

	// Define nu
	var rv mat.Dense
	rv.Mul(mat.NewDense(2, 1, []float64{1, 1}), mat.NewDense(1, 1, []float64{0.5})) // why is this the initial marginal price?
	r.Nu = rv.RawMatrix().Data
	// Define objective functions
	var objective lbfgsb.GeneralObjectiveFunction
	objective.Function = r.f
	objective.Gradient = r.gradf
	fmt.Println(r.Nu)
	minimum, exitStatus := optimizer.Minimize(objective, r.Nu)
	_ = minimum
	fmt.Println(exitStatus.Message)
}

func (r *Router) f(nu []float64) float64 {
	for i, dualVariableComponent := range nu {
		if dualVariableComponent != r.Nu[i] {
			r.findArbs()
		}
	}
	acc := 0.0
	var tempDelta *mat.VecDense
	var tempLambda *mat.VecDense
	var tempNu *mat.VecDense
	for _, cfmm := range r.CFMMs {
		switch cfmm.Type() {
		case dex.Prod:
			tempDelta = mat.NewVecDense(2, cfmm.(dex.ProductTwoCoin).Deltas)
			tempLambda = mat.NewVecDense(2, cfmm.(dex.ProductTwoCoin).Lambdas)
		case dex.Geom:

			tempDelta = mat.NewVecDense(2, cfmm.(dex.GeometricMeanTwoCoin).Deltas)
			tempLambda = mat.NewVecDense(2, cfmm.(dex.GeometricMeanTwoCoin).Lambdas)
		}
		mat.Dot(tempDelta, tempLambda)

		tempNu = mat.NewVecDense(2, r.Nu)
		acc += (mat.Dot(tempLambda, tempNu) - mat.Dot(tempDelta, tempNu))

	}

	return 1 // here is a place where I could use help
}
func (r *Router) gradf(nu []float64) []float64 {
	g := mat.NewVecDense(len(nu), make([]float64, len(nu)))
	for i, dualVariableComponent := range nu {
		if dualVariableComponent != r.Nu[i] {
			r.findArbs()
		}
	}

	// call gradient

	var tempDelta *mat.VecDense
	var tempLambda *mat.VecDense
	var diff *mat.VecDense
	for _, cfmm := range r.CFMMs {
		switch cfmm.Type() {
		case dex.Prod:
			tempDelta = mat.NewVecDense(2, cfmm.(dex.ProductTwoCoin).Deltas)
			tempLambda = mat.NewVecDense(2, cfmm.(dex.ProductTwoCoin).Lambdas)
		case dex.Geom:

			tempDelta = mat.NewVecDense(2, cfmm.(dex.GeometricMeanTwoCoin).Deltas)
			tempLambda = mat.NewVecDense(2, cfmm.(dex.GeometricMeanTwoCoin).Lambdas)
		}

		diff.SubVec(tempLambda, tempDelta)
		g.AddVec(g, diff)
	}
	return nil
}

func (r *Router) findArbs() {
	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(r.NumCFMMs)
	for _, cfmm := range r.CFMMs {
		go cfmm.FindArb(r.Nu) // @TODO understand what Nu is, how their julia code updates it, why does each CFMM in the jula code ostensibly get it's own nu
	}
	wg.Wait()

}
