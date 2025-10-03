package model

import (
	"testing"
)

func TestToCanonicalForm(t *testing.T) {
	lp := &LinearProgram{
		NbConstraints: 2,
		NbVariables:   2,
		VariableNames: []string{"x", "y"},
		Objective:     MINIMIZE,
		ObjCoeff:      []float64{1, 2},
		Comparisons:   []Comparison{BE, LE},
		ConstraintCoeff: [][]float64{{1, 1}, {2, 1}},
		Rhs:           []float64{10, 15},
	}

	lp.ToCanonicalForm()

	if lp.Objective != MAXIMIZE {
		t.Errorf("Expected objective to be MAXIMIZE, but got %d", lp.Objective)
	}

	if lp.ObjCoeff[0] != -1 || lp.ObjCoeff[1] != -2 {
		t.Errorf("Expected ObjCoeff to be [-1, -2], but got %v", lp.ObjCoeff)
	}

	if lp.Comparisons[0] != LE || lp.Comparisons[1] != LE {
		t.Errorf("Expected all comparisons to be LE, but got %v", lp.Comparisons)
	}

	if lp.ConstraintCoeff[0][0] != -1 || lp.ConstraintCoeff[0][1] != -1 || lp.Rhs[0] != -10 {
		t.Errorf("Expected first constraint to be flipped, but got ConstraintCoeff: %v, Rhs: %v", lp.ConstraintCoeff[0], lp.Rhs[0])
	}
}

func TestToSlackForm(t *testing.T) {
	lp := &LinearProgram{
		NbConstraints: 2,
		NbVariables:   2,
		VariableNames: []string{"x", "y"},
		Objective:     MAXIMIZE,
		ObjCoeff:      []float64{1, 2},
		Comparisons:   []Comparison{LE, LE},
		ConstraintCoeff: [][]float64{{1, 1}, {2, 1}},
		Rhs:           []float64{10, 15},
	}

	lp.ToSlackForm()

	if lp.State != Slack {
		t.Errorf("Expected state to be Slack, but got %d", lp.State)
	}

	if lp.NbVariables != 4 {
		t.Errorf("Expected NbVariables to be 4, but got %d", lp.NbVariables)
	}

	if len(lp.SlackVariablesNames) != 2 {
		t.Errorf("Expected 2 slack variables, but got %d", len(lp.SlackVariablesNames))
	}

	if lp.Comparisons[0] != EQ || lp.Comparisons[1] != EQ {
		t.Errorf("Expected all comparisons to be EQ, but got %v", lp.Comparisons)
	}

	if lp.ConstraintCoeff[0][2] != 1 || lp.ConstraintCoeff[1][3] != 1 {
		t.Errorf("Expected slack variables to be added correctly, but got ConstraintCoeff: %v", lp.ConstraintCoeff)
	}
}
