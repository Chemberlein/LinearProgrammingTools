package parser

import (
	"io/ioutil"
	"testing"

	"github.com/Chemberlein/LinearProgrammingTools/model"
)

func TestParse(t *testing.T) {
	// Test case 1: Valid JSON
	jsonData, err := ioutil.ReadFile("tests/example.json")
	if err != nil {
		t.Fatalf("Failed to read example.json: %v", err)
	}

	lp, err := Parse(string(jsonData))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if lp.NbVariables != 3 {
		t.Errorf("Expected NbVariables to be 3, got %d", lp.NbVariables)
	}

	if lp.NbConstraints != 3 {
		t.Errorf("Expected NbConstraints to be 3, got %d", lp.NbConstraints)
	}

	if lp.Objective != model.MINIMIZE {
		t.Errorf("Expected Objective to be MINIMIZE, got %d", lp.Objective)
	}

	expectedObjCoeff := []float64{4, -5, 3}
	if !equalFloat64Slices(lp.ObjCoeff, expectedObjCoeff) {
		t.Errorf("Expected ObjCoeff to be %v, got %v", expectedObjCoeff, lp.ObjCoeff)
	}

	expectedConstraintCoeff := [][]float64{
		{0, 5, 0},
		{0, 4, 3},
		{12, 4, 3},
	}
	if !equalFloat64Matrices(lp.ConstraintCoeff, expectedConstraintCoeff) {
		t.Errorf("Expected ConstraintCoeff to be %v, got %v", expectedConstraintCoeff, lp.ConstraintCoeff)
	}

	expectedRhs := []float64{200, 430, 430}
	if !equalFloat64Slices(lp.Rhs, expectedRhs) {
		t.Errorf("Expected Rhs to be %v, got %v", expectedRhs, lp.Rhs)
	}

	expectedComparisons := []model.Comparison{model.LO, model.LO, model.LO}
	if !equalComparisonSlices(lp.Comparisons, expectedComparisons) {
		t.Errorf("Expected Comparisons to be %v, got %v", expectedComparisons, lp.Comparisons)
	}
}

func TestParse2(t *testing.T) {
	// Test case 2: Valid JSON
	jsonData, err := ioutil.ReadFile("tests/example2.json")
	if err != nil {
		t.Fatalf("Failed to read example2.json: %v", err)
	}

	lp, err := Parse(string(jsonData))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if lp.NbVariables != 2 {
		t.Errorf("Expected NbVariables to be 2, got %d", lp.NbVariables)
	}

	if lp.NbConstraints != 3 {
		t.Errorf("Expected NbConstraints to be 3, got %d", lp.NbConstraints)
	}

	if lp.Objective != model.MAXIMIZE {
		t.Errorf("Expected Objective to be MAXIMIZE, got %d", lp.Objective)
	}

	expectedObjCoeff := []float64{3, 5}
	if !equalFloat64Slices(lp.ObjCoeff, expectedObjCoeff) {
		t.Errorf("Expected ObjCoeff to be %v, got %v", expectedObjCoeff, lp.ObjCoeff)
	}

	expectedConstraintCoeff := [][]float64{
		{1, 0},
		{0, 2},
		{3, 2},
	}
	if !equalFloat64Matrices(lp.ConstraintCoeff, expectedConstraintCoeff) {
		t.Errorf("Expected ConstraintCoeff to be %v, got %v", expectedConstraintCoeff, lp.ConstraintCoeff)
	}

	expectedRhs := []float64{4, 12, 18}
	if !equalFloat64Slices(lp.Rhs, expectedRhs) {
		t.Errorf("Expected Rhs to be %v, got %v", expectedRhs, lp.Rhs)
	}

	expectedComparisons := []model.Comparison{model.LE, model.LE, model.LE}
	if !equalComparisonSlices(lp.Comparisons, expectedComparisons) {
		t.Errorf("Expected Comparisons to be %v, got %v", expectedComparisons, lp.Comparisons)
	}
}

func TestParse3(t *testing.T) {
	// Test case 3: Valid JSON
	jsonData, err := ioutil.ReadFile("tests/example3.json")
	if err != nil {
		t.Fatalf("Failed to read example3.json: %v", err)
	}

	lp, err := Parse(string(jsonData))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if lp.NbVariables != 2 {
		t.Errorf("Expected NbVariables to be 2, got %d", lp.NbVariables)
	}

	if lp.NbConstraints != 2 {
		t.Errorf("Expected NbConstraints to be 2, got %d", lp.NbConstraints)
	}

	if lp.Objective != model.MINIMIZE {
		t.Errorf("Expected Objective to be MINIMIZE, got %d", lp.Objective)
	}

	expectedObjCoeff := []float64{2, 3}
	if !equalFloat64Slices(lp.ObjCoeff, expectedObjCoeff) {
		t.Errorf("Expected ObjCoeff to be %v, got %v", expectedObjCoeff, lp.ObjCoeff)
	}

	expectedConstraintCoeff := [][]float64{
		{1, 1},
		{2, 1},
	}
	if !equalFloat64Matrices(lp.ConstraintCoeff, expectedConstraintCoeff) {
		t.Errorf("Expected ConstraintCoeff to be %v, got %v", expectedConstraintCoeff, lp.ConstraintCoeff)
	}

	expectedRhs := []float64{10, 20}
	if !equalFloat64Slices(lp.Rhs, expectedRhs) {
		t.Errorf("Expected Rhs to be %v, got %v", expectedRhs, lp.Rhs)
	}

	expectedComparisons := []model.Comparison{model.BE, model.LE}
	if !equalComparisonSlices(lp.Comparisons, expectedComparisons) {
		t.Errorf("Expected Comparisons to be %v, got %v", expectedComparisons, lp.Comparisons)
	}
}

func equalFloat64Slices(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalFloat64Matrices(a, b [][]float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !equalFloat64Slices(a[i], b[i]) {
			return false
		}
	}
	return true
}

func equalComparisonSlices(a, b []model.Comparison) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
