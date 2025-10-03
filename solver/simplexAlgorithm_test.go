package solver

import (
	"encoding/json"
	"testing"

	"github.com/Chemberlein/LinearProgrammingTools/model"
)

func TestSolve(t *testing.T) {
	lp := &model.LinearProgram{
		NbConstraints: 3,
		NbVariables:   2,
		VariableNames: []string{"x1", "x2"},
		Objective:     model.MAXIMIZE,
		ObjCoeff:      []float64{3, 5},
		Comparisons:   []model.Comparison{model.LE, model.LE, model.LE},
		ConstraintCoeff: [][]float64{
			{1, 0},
			{0, 2},
			{3, 2},
		},
		Rhs: []float64{4, 12, 18},
	}

	err := Solve(lp)
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}

	// Expected solution for this problem is x1=2, x2=6, objective=36
	expectedSolution := []float64{2, 6, 36}
	if !equalFloat64Slices(lp.ObjVar, expectedSolution, 1e-9) {
		t.Errorf("Expected solution to be %v, got %v", expectedSolution, lp.ObjVar)
	}
}

func TestGetSolutionJSON(t *testing.T) {
	lp := &model.LinearProgram{
		NbConstraints: 3,
		NbVariables:   2,
		VariableNames: []string{"x1", "x2"},
		Objective:     model.MAXIMIZE,
		ObjCoeff:      []float64{3, 5},
		Comparisons:   []model.Comparison{model.LE, model.LE, model.LE},
		ConstraintCoeff: [][]float64{
			{1, 0},
			{0, 2},
			{3, 2},
		},
		Rhs: []float64{4, 12, 18},
	}

	err := Solve(lp)
	if err != nil {
		t.Fatalf("Solve() error = %v", err)
	}

	jsonString, err := lp.GetSolutionJSON()
	if err != nil {
		t.Fatalf("GetSolutionJSON() error = %v", err)
	}

	var solution map[string]float64
	err = json.Unmarshal([]byte(jsonString), &solution)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	expectedSolution := map[string]float64{"x1": 2, "x2": 6, "objective": 36}
	if solution["x1"] != expectedSolution["x1"] || solution["x2"] != expectedSolution["x2"] || solution["objective"] != expectedSolution["objective"] {
		t.Errorf("Expected solution to be %v, got %v", expectedSolution, solution)
	}
}

func equalFloat64Slices(a, b []float64, tolerance float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if (a[i]-b[i]) > tolerance || (b[i]-a[i]) > tolerance {
			return false
		}
	}
	return true
}

func TestSolve_Unbounded(t *testing.T) {
	lp := &model.LinearProgram{
		NbConstraints: 2,
		NbVariables:   2,
		VariableNames: []string{"x1", "x2"},
		Objective:     model.MAXIMIZE,
		ObjCoeff:      []float64{1, 1},
		Comparisons:   []model.Comparison{model.LE, model.LE},
		ConstraintCoeff: [][]float64{
			{1, -1},
			{-1, 1},
		},
		Rhs: []float64{1, 1},
	}

	err := Solve(lp)
	if err == nil || err.Error() != "Unbounded" {
		t.Errorf("Expected error to be 'Unbounded', got %v", err)
	}
}

func TestSolve_Infeasible(t *testing.T) {
	lp := &model.LinearProgram{
		NbConstraints: 2,
		NbVariables:   2,
		VariableNames: []string{"x1", "x2"},
		Objective:     model.MAXIMIZE,
		ObjCoeff:      []float64{1, 1},
		Comparisons:   []model.Comparison{model.LE, model.BE},
		ConstraintCoeff: [][]float64{
			{1, 1},
			{1, 1},
		},
		Rhs: []float64{1, 2},
	}

	err := Solve(lp)
	if err == nil {
		t.Errorf("Expected error for infeasible problem, got nil")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}

func TestIsInitiallyFeasible(t *testing.T) {
	t.Run("Feasible", func(t *testing.T) {
		table := &SimplexTable{
			data: [][]float64{
				{1, 1, 1, 0, 10},
				{2, 1, 0, 1, 15},
				{-1, -2, 0, 0, 0},
			},
		}
		if !table.IsInitiallyFeasible() {
			t.Errorf("Expected tableau to be initially feasible")
		}
	})

	t.Run("Infeasible", func(t *testing.T) {
		table := &SimplexTable{
			data: [][]float64{
				{1, 1, 1, 0, -10},
				{2, 1, 0, 1, 15},
				{-1, -2, 0, 0, 0},
			},
		}
		if table.IsInitiallyFeasible() {
			t.Errorf("Expected tableau to be initially infeasible")
		}
	})
}
