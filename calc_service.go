package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"unicode"
)

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

var (
	ErrInvalidCharacter       = errors.New("invalid character")
	ErrUnbalancedParentheses  = errors.New("unbalanced parentheses")
	ErrInvalidExpression      = errors.New("invalid expression")
)

func Calc(expression string) (float64, error) {
	rpn, err := toRPN(expression)
	if err != nil {
		return 0, err
	}
	result, err := evalRPN(rpn)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func toRPN(expression string) ([]string, error) {
	var output []string
	var opsStack []rune

	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		switch {
			case unicode.IsDigit(char) || char == '.':
				num := string(char)
				for i+1 < len(expression) && (unicode.IsDigit(rune(expression[i+1])) || rune(expression[i+1]) == '.') {
					i++
					num += string(expression[i])
				}
				output = append(output, num)
			case char == '(':
				opsStack = append(opsStack, char)
			case char == ')':
				found := false
				for len(opsStack) > 0 {
					top := opsStack[len(opsStack)-1]
					opsStack = opsStack[:len(opsStack)-1]
					if top == '(' {
						found = true
						break
					}
					output = append(output, string(top))
				}
				if !found {
					return nil, ErrUnbalancedParentheses
				}
			case char == '+' || char == '-' || char == '*' || char == '/':
				for len(opsStack) > 0 {
					top := opsStack[len(opsStack)-1]
					if top == '(' || precedence[char] > precedence[top] {
						break
					}
					output = append(output, string(top))
					opsStack = opsStack[:len(opsStack)-1]
				}
				opsStack = append(opsStack, char)
			default:
				return nil, fmt.Errorf("%w: %c", ErrInvalidCharacter, char)
		}
	}

	for len(opsStack) > 0 {
		top := opsStack[len(opsStack)-1]
		opsStack = opsStack[:len(opsStack)-1]
		if top == '(' || top == ')' {
			return nil, ErrUnbalancedParentheses
		}
		output = append(output, string(top))
	}

	return output, nil
}

func evalRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		if value, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, value)
		} else {
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
				case "+":
					result = a + b
				case "-":
					result = a - b
				case "*":
					result = a * b
				case "/":
					if b == 0 {
						return 0, errors.New("division by zero")
					}
					result = a / b
				default:
					return 0, fmt.Errorf("%w: %s", ErrInvalidCharacter, token)
			}

			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	result, err := Calc(req.Expression)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if errors.Is(err, ErrInvalidCharacter) || errors.Is(err, ErrUnbalancedParentheses) || errors.Is(err, ErrInvalidExpression) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(CalcResponse{Error: "Expression is not valid"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(CalcResponse{Error: "Internal server error"})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CalcResponse{Result: fmt.Sprintf("%f", result)})
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
