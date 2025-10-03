package model

import (
	"encoding/json"
	"fmt"
)

// Enums for Objective and Comparison
type Objectiv int

const (
	MINIMIZE Objectiv = iota
	MAXIMIZE
)

type Comparison int

const (
	EQ Comparison = iota
	LO
	LE
	BI
	BE
)

type LPState int

const (
	Undefined LPState = iota
	Canonical
	Slack
)

// LinearProgram represents a linear programming problem.
type LinearProgram struct {
	NbConstraints       int
	NbVariables         int
	VariableNames       []string
	SlackVariablesNames []string
	Objective           Objectiv
	ObjVar              []float64
	ObjCoeff            []float64
	Comparisons         []Comparison
	ConstraintCoeff     [][]float64
	Rhs                 []float64
	State               LPState
}

// GetSolutionJSON returns the solution of the linear program in JSON format.
func (lp *LinearProgram) GetSolutionJSON() (string, error) {
	if lp.ObjVar == nil {
		return "", fmt.Errorf("solution not available")
	}

	solution := make(map[string]float64)
	for i, name := range lp.VariableNames {
		solution[name] = lp.ObjVar[i]
	}
	solution["objective"] = lp.ObjVar[len(lp.ObjVar)-1]

	jsonBytes, err := json.Marshal(solution)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}