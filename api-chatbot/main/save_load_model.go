package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/golang-chatbot/api-chatbot/mlp"
)

// SaveTrainedNN - saves the trained neural network
func SaveTrainedNN(net mlp.Network, baseVects baseVectors) {
	h, err := os.Create("data/hweights.model")
	defer h.Close()
	if err == nil {
		net.HiddenWeights.MarshalBinaryTo(h)
	} else {
		panic(err)
	}

	o, err := os.Create("data/oweights.model")
	defer o.Close()
	if err == nil {
		net.OutputWeights.MarshalBinaryTo(o)
	} else {
		panic(err)
	}

	b, err := json.Marshal(baseVects)
	if err == nil {
		err = ioutil.WriteFile("data/base_vectors.json", b, 0644)
		if err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
}

// LoadTrainedNN - loads a trained neural network
func LoadTrainedNN(net *mlp.Network, baseVects *baseVectors) {
	h, err := os.Open("data/hweights.model")
	defer h.Close()
	if err == nil {
		net.HiddenWeights.Reset()
		net.HiddenWeights.UnmarshalBinaryFrom(h)
	} else {
		panic(err)
	}

	o, err := os.Open("data/oweights.model")
	defer o.Close()
	if err == nil {
		net.OutputWeights.Reset()
		net.OutputWeights.UnmarshalBinaryFrom(o)
	} else {
		panic(err)
	}

	file, err := ioutil.ReadFile("data/base_vectors.json")
	if err == nil {
		json.Unmarshal([]byte(file), &baseVects)
	} else {
		panic(err)
	}
}
