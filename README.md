# CheLinearProgramming

A Go package for solving linear programming problems.

## Features

*   Parse linear programming problems from JSON files.
*   Solve problems using the simplex algorithm.
*   Convert problems to canonical and slack forms.

## Installation

```bash
go get github.com/Chemberlein/LinearProgrammingTools
```

## Getting Started

### Defining a Problem

To define a linear programming problem, you need to create a JSON object with the following structure:

```json
{
  "numberOfVariables": 3,
  "numberOfConstraints": 3,
  "objectiveFunction": {
    "objective": "minimize",
    "equasion": "4*x1 -5*x2 +3*x3"
  },
  "constraints": [
    "+5*x2 <= 200",
    "+4*x2 +3*x3 <= 430",
    "12*x1 +4*x2 +3*x3 <= 430"
  ]
}
```

- `numberOfVariables`: The number of variables in the problem.
- `numberOfConstraints`: The number of constraints.
- `objectiveFunction`: An object defining the function to be optimized.
    - `objective`: Either "minimize" or "maximize".
    - `equasion`: The objective function equation.
- `constraints`: An array of strings, where each string is a constraint.

### Solving the Problem and Getting the Solution

The following example shows how to parse a JSON string, solve the linear programming problem, and print the solution.

```go
package main

import (
	"fmt"
	"github.com/Chemberlein/LinearProgrammingTools/parser"
	"github.com/Chemberlein/LinearProgrammingTools/solver"
)

func main() {
	jsonData := `
	{
	  "numberOfVariables": 3,
	  "numberOfConstraints": 3,
	  "objectiveFunction": {
	    "objective": "minimize",
	    "equasion": "4*x1 -5*x2 +3*x3"
	  },
	  "constraints": [
	    "+5*x2 <= 200",
	    "+4*x2 +3*x3 <= 430",
	    "12*x1 +4*x2 +3*x3 <= 430"
	  ]
	}
	`

	// Parse the JSON data into a LinearProgram object
	lp, err := parser.Parse(jsonData)
	if err != nil {
		fmt.Printf("Failed to parse JSON: %v
", err)
		return
	}

	// Solve the linear programming problem
	err = solver.Solve(lp)
	if err != nil {
		fmt.Printf("Failed to solve: %v
", err)
		return
	}

	// Convert the solution to JSON format
	solutionJSON, err := parser.ConvertLPToJSON(lp)
	if err != nil {
		fmt.Printf("Failed to convert solution to JSON: %v
", err)
		return
	}

	// Print the solution
	fmt.Println(solutionJSON)
}
```

### Interpreting the Solution

The output will be a JSON object containing the solution to the problem. The solution will include the optimal value of the objective function and the values of the variables that achieve this optimal value.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
