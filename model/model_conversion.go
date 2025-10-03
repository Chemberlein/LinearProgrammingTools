package model

import "fmt"

// ToCanonicalForm converts the linear program to the canonical form (<= constraints).
func (lp *LinearProgram) ToCanonicalForm() {
	if lp.State != Undefined {
		return
	}
	lp.EnsureMaximization()
	lp.EnsureNonNegativeRhs()
	lp.ConvertToLeConstraints()
	lp.State = Canonical
}

// ToSlackForm converts the linear program to the slack form (equality constraints).
func (lp *LinearProgram) ToSlackForm() {
	if lp.State == Slack {
		return
	}
	if lp.State == Undefined {
		lp.ToCanonicalForm()
	}
	lp.ConvertToEqualities()
	lp.State = Slack
}

func (lp *LinearProgram) EnsureMaximization() {
	if lp.Objective == MINIMIZE {
		lp.Objective = MAXIMIZE
		for i := range lp.ObjCoeff {
			lp.ObjCoeff[i] *= -1
		}
	}
}

func (lp *LinearProgram) EnsureNonNegativeRhs() {
	for i := range lp.Rhs {
		if lp.Rhs[i] < 0 {
			lp.Rhs[i] *= -1
			lp.ConstraintCoeff[i] = MultiplyRow(lp.ConstraintCoeff[i], -1)
			lp.Comparisons[i] = FlipComparison(lp.Comparisons[i])
		}
	}
}

func (lp *LinearProgram) ConvertToLeConstraints() {
	var newConstraintCoeff [][]float64
	var newRhs []float64
	var newComparisons []Comparison
	newNbConstraints := 0

	for i := 0; i < lp.NbConstraints; i++ {
		switch lp.Comparisons[i] {
		case LE:
			newConstraintCoeff = append(newConstraintCoeff, lp.ConstraintCoeff[i])
			newRhs = append(newRhs, lp.Rhs[i])
			newComparisons = append(newComparisons, LE)
			newNbConstraints++
		case LO:
			newConstraintCoeff = append(newConstraintCoeff, lp.ConstraintCoeff[i])
			newRhs = append(newRhs, lp.Rhs[i])
			newComparisons = append(newComparisons, LE)
			newNbConstraints++
		case BE:
			newConstraintCoeff = append(newConstraintCoeff, MultiplyRow(lp.ConstraintCoeff[i], -1))
			newRhs = append(newRhs, -lp.Rhs[i])
			newComparisons = append(newComparisons, LE)
			newNbConstraints++
		case BI:
			newConstraintCoeff = append(newConstraintCoeff, MultiplyRow(lp.ConstraintCoeff[i], -1))
			newRhs = append(newRhs, -lp.Rhs[i])
			newComparisons = append(newComparisons, LE)
			newNbConstraints++
		case EQ:
			// Add <= constraint
			newConstraintCoeff = append(newConstraintCoeff, lp.ConstraintCoeff[i])
			newRhs = append(newRhs, lp.Rhs[i])
			newComparisons = append(newComparisons, LE)
			newNbConstraints++

			// Add >= constraint, which is then converted to <=
			newConstraintCoeff = append(newConstraintCoeff, MultiplyRow(lp.ConstraintCoeff[i], -1))
			newRhs = append(newRhs, -lp.Rhs[i])
			newComparisons = append(newComparisons, LE)
			newNbConstraints++
		}
	}
	lp.ConstraintCoeff = newConstraintCoeff
	lp.Rhs = newRhs
	lp.Comparisons = newComparisons
	lp.NbConstraints = newNbConstraints
}

func (lp *LinearProgram) ConvertToEqualities() {
	for i := 0; i < lp.NbConstraints; i++ {
		switch lp.Comparisons[i] {
		case LE:
			lp.AddSlackVariable(i)
		case BI:
			lp.AddSurplusVariable(i)
		case LO:
			lp.AddSlackVariable(i)
		case BE:
			lp.AddSurplusVariable(i)
		}
	}
}

func (lp *LinearProgram) AddSlackVariable(constraintIndex int) {
	lp.NbVariables++
	newVarName := fmt.Sprintf("s%d", len(lp.SlackVariablesNames)+1)
	lp.SlackVariablesNames = append(lp.SlackVariablesNames, newVarName)
	lp.ObjCoeff = append(lp.ObjCoeff, 0)

	for i := range lp.ConstraintCoeff {
		newCol := 0.0
		if i == constraintIndex {
			newCol = 1.0
		}
		lp.ConstraintCoeff[i] = append(lp.ConstraintCoeff[i], newCol)
	}
	lp.Comparisons[constraintIndex] = EQ
}

func (lp *LinearProgram) AddSurplusVariable(constraintIndex int) {
	lp.NbVariables++
	newVarName := fmt.Sprintf("s%d", len(lp.SlackVariablesNames)+1)
	lp.SlackVariablesNames = append(lp.SlackVariablesNames, newVarName)
	lp.ObjCoeff = append(lp.ObjCoeff, 0)

	for i := range lp.ConstraintCoeff {
		newCol := 0.0
		if i == constraintIndex {
			newCol = -1.0
		}
		lp.ConstraintCoeff[i] = append(lp.ConstraintCoeff[i], newCol)
	}
	lp.Comparisons[constraintIndex] = EQ
}

func MultiplyRow(row []float64, scalar float64) []float64 {
	newRow := make([]float64, len(row))
	for i, val := range row {
		newRow[i] = val * scalar
	}
	return newRow
}

func FlipComparison(comp Comparison) Comparison {
	switch comp {
	case LE:
		return BE
	case BE:
		return LE
	case LO:
		return BI
	case BI:
		return LO
	default:
		return comp
	}
}
