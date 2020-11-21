package main

import (
	"sort"
	"strings"
)

// base vectors for phrases and targets.
type baseVectors struct {
	BaseWords   []string
	BaseTargets []string
}

// vectorized phrase and values.
type vectorizedVals struct {
	phrase  []float64
	targets []float64
}

// genBaseVectors - generates standard base vectors to standarize targets and phrases vectorization.
func genBaseVectors(data *Data) baseVectors {

	wordsMap := make(map[string]int, 0)
	targetsMap := make(map[string]int, 0)

	baseVects := baseVectors{}
	baseWords := make([]string, 0)
	baseTargets := make([]string, 0)

	for _, phrase := range data.phrase {
		words := strings.Split(phrase, "_")
		for _, word := range words {
			wordsMap[word]++
		}
	}

	for word := range wordsMap {
		baseWords = append(baseWords, word)
	}

	for _, targets := range data.target {
		for _, target := range targets {
			targetsMap[target]++
		}
	}

	for target := range targetsMap {
		baseTargets = append(baseTargets, target)
	}

	sort.Strings(baseWords)
	sort.Strings(baseTargets)
	baseVects.BaseWords = baseWords
	baseVects.BaseTargets = baseTargets

	return baseVects
}

// vectorizer - vectorize a phrase and its targets.
func vectorizer(phrase string, targets []string, baseVects baseVectors) vectorizedVals {

	vectorized := vectorizedVals{
		phrase:  make([]float64, len(baseVects.BaseWords)),
		targets: make([]float64, len(baseVects.BaseTargets)),
	}

	for i := 0; i < len(baseVects.BaseTargets); i++ {
		vectorized.targets[i] = 0.03
	}
	words := strings.Split(phrase, "_")
	for _, word := range words {
		for idx, baseWord := range baseVects.BaseWords {
			if baseWord == word {
				vectorized.phrase[idx]++
				break
			}
		}
	}

	for _, target := range targets {
		for idx, baseTarget := range baseVects.BaseTargets {
			if baseTarget == target {
				vectorized.targets[idx] = 0.97
			}
		}
	}

	return vectorized
}
