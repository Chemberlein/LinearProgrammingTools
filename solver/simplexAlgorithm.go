package solver

import (
	"fmt"
	"math"
	"strings"
	"github.com/Chemberlein/LinearProgrammingTools/model"
)

// SimplexTable represents the simplex tableau.
type SimplexTable struct {
	data           [][]float64 // (constraints + objective row) x (variables + slacks + RHS)
	basicVariables []float64
}

// String returns a string representation of the simplex table.
func (table *SimplexTable) String() string {
	var builder strings.Builder
	builder.WriteString("Simplex Tableau:\n")

	// Get column widths

	colWidths := make([]int, len(table.data[0]))
	for _, row := range table.data {
		for j, val := range row {
			width := len(fmt.Sprintf("%.2f", val))
			if width > colWidths[j] {
				colWidths[j] = width
			}
		}
	}

	// Print header
	for j := 0; j < len(table.data[0]); j++ {
		builder.WriteString(fmt.Sprintf("%*s ", colWidths[j], fmt.Sprintf("x%d", j+1)))
	}

	builder.WriteString("\n")

	// Print data
	for _, row := range table.data {
		for j, val := range row {
			builder.WriteString(fmt.Sprintf("%*.*f ", colWidths[j], 2, val))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// InitializeTableau creates the initial simplex tableau from a standardized linear program.
func (table *SimplexTable) InitializeTableau(problem *model.LinearProgram) {
	if problem.State != model.Slack {
		problem.ToSlackForm()
	}

	m := problem.NbConstraints
	n := problem.NbVariables // n is now the total number of variables including slacks
	n_orig := n - m          // number of original variables

	numRows := m + 1
	numCols := n + 1 // n variables + 1 RHS column

	table.data = make([][]float64, numRows)
	for i := range table.data {
		table.data[i] = make([]float64, numCols)
	}
	table.basicVariables = make([]float64, m)

	// Fill constraint rows
	for i := 0; i < m; i++ {
		// Copy constraint coefficients
		for j := 0; j < n; j++ {
			table.data[i][j] = problem.ConstraintCoeff[i][j]
		}

		// Set RHS
		table.data[i][n] = problem.Rhs[i]

		// Initial basic variables are the slack variables
		table.basicVariables[i] = float64(n_orig + i)
	}

	// Fill objective row (bottom row)
	rowIndex := m
	for j := 0; j < n; j++ {
		table.data[rowIndex][j] = -problem.ObjCoeff[j]
	}

	// The rest of the objective row is 0
	table.data[rowIndex][n] = 0.0 // Initial objective value
}

// FindEnteringVariable finds the entering variable based on Bland's rule.
// It searches for the first negative coefficient (smallest index) in the objective row.
func (table *SimplexTable) FindEnteringVariable() int {
	objectiveRow := len(table.data) - 1
	epsilon := 1e-10

	// Search for the first negative coefficient (smallest index)
	for j := 0; j < len(table.data[objectiveRow])-1; j++ {
		coefficient := table.data[objectiveRow][j]
		if coefficient < -epsilon { // Significantly negative
			return j // Return the first negative coefficient's index
		}
	}

	return -1 // Optimal solution found
}

// FindLeavingVariable finds the leaving variable using the minimum ratio test.
// When ratios are tied, it selects the one with the smallest row index (Bland's rule).
func (table *SimplexTable) FindLeavingVariable(pivotCol int) int {
	numConstraintRows := len(table.data) - 1
	rhsCol := len(table.data[0]) - 1
	smallestRatio := math.Inf(1)
	pivotRow := -1
	epsilon := 1e-10

	for i := 0; i < numConstraintRows; i++ {
		pivotColValue := table.data[i][pivotCol]
		rhsValue := table.data[i][rhsCol]

		if pivotColValue > epsilon {
			ratio := rhsValue / pivotColValue

			if ratio >= -epsilon { // Non-negative ratio
				// Use epsilon for comparison to handle ties
				if ratio < smallestRatio-epsilon {
					smallestRatio = ratio
					pivotRow = i
				}
				// If tied, pivotRow keeps the smaller index (Bland's rule)
			}
		}
	}

	return pivotRow
}

// PerformPivot performs the pivot operation on the tableau.
// It modifies the tableau in-place.
func (table *SimplexTable) PerformPivot(pivotRow, pivotCol int) {
	numRows := len(table.data)
	numCols := len(table.data[0])
	pivotElement := table.data[pivotRow][pivotCol]

	// Step 1: Normalize pivot row
	for j := 0; j < numCols; j++ {
		table.data[pivotRow][j] /= pivotElement
	}

	// Step 2: Eliminate pivot column in all other rows
	for i := 0; i < numRows; i++ {
		if i != pivotRow {
			factor := table.data[i][pivotCol]
			for j := 0; j < numCols; j++ {
				table.data[i][j] -= factor * table.data[pivotRow][j]
			}
		}
	}
}

// IsInitiallyFeasible checks if the initial tableau is feasible.
func (table *SimplexTable) IsInitiallyFeasible() bool {
	rhsCol := len(table.data[0]) - 1
	numConstraintRows := len(table.data) - 1

	for i := 0; i < numConstraintRows; i++ {
		if table.data[i][rhsCol] < 0 {
			return false
		}
	}

	return true
}

// ExtractSolution extracts the solution from the final simplex tableau.
func (table *SimplexTable) ExtractSolution(problem *model.LinearProgram) []float64 {
	numOrigVars := problem.NbVariables - problem.NbConstraints
	solution := make([]float64, numOrigVars)
	rhsCol := len(table.data[0]) - 1

	// Initialize solution with zeros
	for i := 0; i < numOrigVars; i++ {
		solution[i] = 0
	}

	for i := 0; i < len(table.basicVariables); i++ {
		basicVarIndex := int(table.basicVariables[i])
		if basicVarIndex < numOrigVars {
			solution[basicVarIndex] = table.data[i][rhsCol]
		}
	}

	// The last element of the solution is the objective value
	objectiveRow := len(table.data) - 1
	objectiveValue := table.data[objectiveRow][rhsCol]

	solution = append(solution, objectiveValue)

	return solution
}

// Solve will find the values for the variables.
func Solve(lp *model.LinearProgram) error {
	originalObjective := lp.Objective
	lp.ToSlackForm()

	var table SimplexTable
	table.InitializeTableau(lp)

	if !table.IsInitiallyFeasible() {
		return fmt.Errorf("infeasible problem")
	}

	for {
		pivotCol := table.FindEnteringVariable()
		if pivotCol == -1 {
			lp.ObjVar = table.ExtractSolution(lp)
			if originalObjective == model.MINIMIZE {
				lp.ObjVar[len(lp.ObjVar)-1] *= -1
			}
			lp.State = model.Undefined
			return nil
		}

		pivotRow := table.FindLeavingVariable(pivotCol)
		if pivotRow == -1 {
			return fmt.Errorf("Unbounded")
		}

		table.PerformPivot(pivotRow, pivotCol)

		table.basicVariables[pivotRow] = float64(pivotCol)
	}

	return nil
}
