package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Chemberlein/LinearProgrammingTools/model"
)

// Regular expressions for parsing
var (
	varRegex   = regexp.MustCompile(`([a-zA-Z]+\d*)`)
	coeffRegex = regexp.MustCompile(`([+-]?)\s*(\d*\.?\d*)\s*\*?\s*([a-zA-Z]+\d*)`)
	compRegex  = regexp.MustCompile(`([<>=]+)`)
)

// JSONLinearProgram is the structure for parsing the JSON input.
type JSONLinearProgram struct {
	NumberOfVariables   int               `json:"numberOfVariables"`
	NumberOfConstraints int               `json:"numberOfConstraints"`
	ObjectiveFunction   ObjectiveFunction `json:"objectiveFunction"`
	Constraints         []string          `json:"constraints"`
}

// ObjectiveFunction is the structure for parsing the objective function from JSON.
type ObjectiveFunction struct {
	Objective string `json:"objective"`
	Equation  string `json:"equasion"`
}

// Parse takes a JSON string and returns a LinearProgram.
func Parse(jsonData string) (*model.LinearProgram, error) {
	jsonLP := &JSONLinearProgram{}
	err := json.Unmarshal([]byte(jsonData), jsonLP)
	if err != nil {
		return nil, err
	}

	lp := &model.LinearProgram{
		NbConstraints: jsonLP.NumberOfConstraints,
		NbVariables:   jsonLP.NumberOfVariables,
	}

	parseVariableNames(lp, jsonLP)
	varMap := buildVarMap(lp)

	err = parseObjectiveFunction(lp, jsonLP, varMap)
	if err != nil {
		return nil, err
	}

	err = parseConstraints(lp, jsonLP, varMap)
	if err != nil {
		return nil, err
	}

	return lp, nil
}

func parseVariableNames(lp *model.LinearProgram, jsonLP *JSONLinearProgram) {
	matches := varRegex.FindAllString(jsonLP.ObjectiveFunction.Equation, -1)
	lp.VariableNames = unique(matches)
}

func buildVarMap(lp *model.LinearProgram) map[string]int {
	varMap := make(map[string]int)
	for i, name := range lp.VariableNames {
		varMap[name] = i
	}
	return varMap
}

func parseObjectiveFunction(lp *model.LinearProgram, jsonLP *JSONLinearProgram, varMap map[string]int) error {
	if strings.ToLower(jsonLP.ObjectiveFunction.Objective) == "minimize" {
		lp.Objective = model.MINIMIZE
	} else {
		lp.Objective = model.MAXIMIZE
	}

	lp.ObjCoeff = make([]float64, lp.NbVariables)
	parseEquation(jsonLP.ObjectiveFunction.Equation, lp.ObjCoeff, varMap)
	return nil
}

func parseConstraints(lp *model.LinearProgram, jsonLP *JSONLinearProgram, varMap map[string]int) error {
	lp.ConstraintCoeff = make([][]float64, lp.NbConstraints)
	lp.Rhs = make([]float64, lp.NbConstraints)
	lp.Comparisons = make([]model.Comparison, lp.NbConstraints)

	for i, constr := range jsonLP.Constraints {
		lp.ConstraintCoeff[i] = make([]float64, lp.NbVariables)
		parts := compRegex.Split(constr, -1)
		compStr := compRegex.FindString(constr)

		err := setComparison(lp, i, compStr)
		if err != nil {
			return err
		}

		parseEquation(parts[0], lp.ConstraintCoeff[i], varMap)

		bVal, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return err
		}
		lp.Rhs[i] = bVal
	}
	return nil
}

func setComparison(lp *model.LinearProgram, index int, compStr string) error {
	switch compStr {
	case "<":
		lp.Comparisons[index] = model.LO
	case "<=":
		lp.Comparisons[index] = model.LE
	case ">":
		lp.Comparisons[index] = model.BI
	case ">=":
		lp.Comparisons[index] = model.BE
	case "=":
		lp.Comparisons[index] = model.EQ
	default:
		return fmt.Errorf("invalid comparison operator: %s", compStr)
	}
	return nil
}

func parseEquation(equation string, coeffs []float64, varMap map[string]int) {
	matches := coeffRegex.FindAllStringSubmatch(equation, -1)

	for _, match := range matches {
		sign := 1.0
		if match[1] == "-" {
			sign = -1.0
		}

		coeff := 1.0
		if match[2] != "" {
			coeff, _ = strconv.ParseFloat(match[2], 64)
		}

		varName := match[3]
		if idx, ok := varMap[varName]; ok {
			coeffs[idx] = sign * coeff
		}
	}
}

// unique returns a unique slice of strings
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// ConvertLPToJSON converts a LinearProgram back to a JSON string.
func ConvertLPToJSON(lp *model.LinearProgram) (string, error) {
	jsonLP := &JSONLinearProgram{
		NumberOfVariables:   lp.NbVariables,
		NumberOfConstraints: lp.NbConstraints,
	}

	// Convert objective function
	objObjective := "maximize"
	if lp.Objective == model.MINIMIZE {
		objObjective = "minimize"
	}
	jsonLP.ObjectiveFunction = ObjectiveFunction{
		Objective: objObjective,
		Equation:  equationToString(lp.ObjCoeff, lp.VariableNames, lp.SlackVariablesNames),
	}

	// Convert constraints
	jsonLP.Constraints = make([]string, lp.NbConstraints)
	for i := 0; i < lp.NbConstraints; i++ {
		compStr, err := comparisonToString(lp.Comparisons[i])
		if err != nil {
			return "", err
		}
		jsonLP.Constraints[i] = equationToString(lp.ConstraintCoeff[i], lp.VariableNames, lp.SlackVariablesNames) + " " + compStr + " " + strconv.FormatFloat(lp.Rhs[i], 'f', -1, 64)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err := enc.Encode(jsonLP)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func equationToString(coeffs []float64, varNames []string, slackVariablesNames []string) string {
	var parts []string
	allVarNames := append(varNames, slackVariablesNames...)
	for i, coeff := range coeffs {
		if coeff == 0 {
			continue
		}
		// Format the coefficient to remove trailing zeros
		coeffStr := strconv.FormatFloat(coeff, 'f', -1, 64)
		// Add a plus sign for positive coefficients, but not for the first term if it's positive.
		if len(parts) > 0 && coeff > 0 {
			coeffStr = "+" + coeffStr
		}
		if i < len(allVarNames) {
			parts = append(parts, coeffStr+"*"+allVarNames[i])
		}
	}
	return strings.Join(parts, " ")
}

func comparisonToString(comp model.Comparison) (string, error) {
	switch comp {
	case model.EQ:
		return "=", nil
	case model.LO:
		return "<", nil
	case model.LE:
		return "<=", nil
	case model.BI:
		return ">", nil
	case model.BE:
		return ">=", nil
	default:
		return "", fmt.Errorf("invalid comparison operator: %d", comp)
	}
}