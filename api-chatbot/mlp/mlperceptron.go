package mlp

// reference: https://sausheong.github.io/posts/how-to-build-a-simple-artificial-neural-network-with-go/

import (
	"gonum.org/v1/gonum/mat"
)

// Network - Neural network structure.
type Network struct {
	inputs        int
	hiddens       int
	outputs       int
	HiddenWeights *mat.Dense
	OutputWeights *mat.Dense
	learningRate  float64
}

// CreateNetwork - Creates a new NN.
func CreateNetwork(input, hidden, output int, rate float64) (net Network) {
	net = Network{
		inputs:       input,
		hiddens:      hidden,
		outputs:      output,
		learningRate: rate,
	}

	net.HiddenWeights = mat.NewDense(net.hiddens, net.inputs, randomArray(net.hiddens*net.inputs, float64(net.inputs)))
	net.OutputWeights = mat.NewDense(net.outputs, net.hiddens, randomArray(net.outputs*net.hiddens, float64(net.hiddens)))

	return net
}

// Predict - Use the NN to make a prediction.
func (net Network) Predict(inputData []float64) mat.Matrix {
	// Forward propagation
	inputs := mat.NewDense(len(inputData), 1, inputData)
	hiddenInputs := dot(net.HiddenWeights, inputs)
	hiddenOutputs := apply(sigmoid, hiddenInputs)
	finalInputs := dot(net.OutputWeights, hiddenOutputs)
	finalOutputs := apply(sigmoid, finalInputs)
	return finalOutputs
}

// Train - train Neural Network.
func (net *Network) Train(inputData []float64, targetData []float64) {
	// forward propagation.
	inputs := mat.NewDense(len(inputData), 1, inputData)
	hiddenInputs := dot(net.HiddenWeights, inputs)
	hiddenOutputs := apply(sigmoid, hiddenInputs)
	finalInputs := dot(net.OutputWeights, hiddenOutputs)
	finalOutputs := apply(sigmoid, finalInputs)

	// find errors.
	targets := mat.NewDense(len(targetData), 1, targetData)
	outputErrors := subtract(targets, finalOutputs)
	hiddenErrors := dot(net.OutputWeights.T(), outputErrors)

	// backward propagation.
	net.OutputWeights = add(net.OutputWeights,
		scale(net.learningRate,
			dot(multiply(outputErrors, sigmoidPrime(finalOutputs)),
				hiddenOutputs.T()))).(*mat.Dense)

	net.HiddenWeights = add(net.HiddenWeights,
		scale(net.learningRate,
			dot(multiply(hiddenErrors, sigmoidPrime(hiddenOutputs)),
				inputs.T()))).(*mat.Dense)
}
