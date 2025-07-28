package calculator

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	ErrUnclosedParentheses = errors.New("unclosed parentheses")
	ErrNotEnoughOperands   = errors.New("not enough operands")
	ErrUnknownOperator     = errors.New("unknown operator")
	ErrInvalidToken        = errors.New("invalid token")
	ErrEmptyExpression     = errors.New("empty expression")
	ErrDivisionByZero      = errors.New("division by zero")
	ErrZeroBase            = errors.New("zero to a non-positive exponent")
	ErrNegativeBase        = errors.New("negative base to a non-integer exponent")
)

func Calc(expression string) (float64, error) {
	parts, err := splitTokens(expression)
	if err != nil {
		return 0, err
	}

	return evaluateExpression(parts)
}

func splitTokens(expression string) ([]string, error) {
	var (
		tokens     []string
		numBuilder strings.Builder
	)

	for idx, ch := range expression {
		switch ch {
		case ' ', '\t':
			continue
		case '+', '*', '/', '^', '(', ')':
			if numBuilder.Len() > 0 {
				tokens = append(tokens, numBuilder.String())
				numBuilder.Reset()
			}
			tokens = append(tokens, string(ch))
		case '-':
			if numBuilder.Len() > 0 {
				tokens = append(tokens, numBuilder.String())
				numBuilder.Reset()
			}
			if idx == 0 || expression[idx-1] == '(' || isOperator(string(expression[idx-1])) {
				numBuilder.WriteRune(ch)
			} else {
				tokens = append(tokens, string(ch))
			}
		case 'e':
			numBuilder.WriteString(strconv.FormatFloat(math.E, 'f', -1, 64))
		case 'p':
			numBuilder.WriteString(strconv.FormatFloat(math.Pi, 'f', -1, 64))
		default:
			numBuilder.WriteRune(ch)
		}
	}

	if numBuilder.Len() > 0 {
		tokens = append(tokens, numBuilder.String())
	}

	if err := verifyTokens(tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func isOperator(symbol string) bool {
	return symbol == "+" || symbol == "-" || symbol == "*" || symbol == "/" || symbol == "^"
}

func verifyTokens(tokens []string) error {
	if len(tokens) == 0 {
		return ErrEmptyExpression
	}

	openBrackets := 0
	for _, tkn := range tokens {
		if tkn == "(" {
			openBrackets++
		} else if tkn == ")" {
			openBrackets--
		}
	}

	if openBrackets != 0 {
		return ErrUnclosedParentheses
	}

	return nil
}

func evaluateExpression(tokens []string) (float64, error) {
	numStack := make([]float64, 0, len(tokens))
	opStack := make([]string, 0, len(tokens))

	priority := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
		"^": 3,
	}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if token == "(" {
			opStack = append(opStack, token)
		} else if token == ")" {
			for len(opStack) > 0 && opStack[len(opStack)-1] != "(" {
				if err := processOperator(&numStack, &opStack); err != nil {
					return 0, err
				}
			}
			opStack = opStack[:len(opStack)-1]
		} else if _, isOp := priority[token]; isOp {
			for len(opStack) > 0 && priority[opStack[len(opStack)-1]] >= priority[token] {
				if err := processOperator(&numStack, &opStack); err != nil {
					return 0, err
				}
			}
			opStack = append(opStack, token)
		} else {
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("%w: %s", ErrInvalidToken, token)
			}
			numStack = append(numStack, value)
		}
	}

	for len(opStack) > 0 {
		if err := processOperator(&numStack, &opStack); err != nil {
			return 0, err
		}
	}

	return numStack[0], nil
}

func processOperator(numStack *[]float64, opStack *[]string) error {
	if len(*numStack) < 2 {
		if len(*numStack) > 0 && (*opStack)[len(*opStack)-1] == "-" {
			(*numStack)[len(*numStack)-1] *= -1
			*opStack = (*opStack)[:len(*opStack)-1]
			return nil
		}

		return fmt.Errorf("%w for %s", ErrNotEnoughOperands, (*opStack)[len(*opStack)-1])
	}

	num2 := (*numStack)[len(*numStack)-1]
	num1 := (*numStack)[len(*numStack)-2]
	*numStack = (*numStack)[:len(*numStack)-2]

	op := (*opStack)[len(*opStack)-1]
	*opStack = (*opStack)[:len(*opStack)-1]

	var res float64
	switch op {
	case "+":
		res = num1 + num2
	case "-":
		res = num1 - num2
	case "*":
		res = num1 * num2
	case "/":
		if num2 == 0 {
			return fmt.Errorf("%w: %f/%f", ErrDivisionByZero, num1, num2)
		}

		res = num1 / num2
	case "^":
		if num1 == 0 && num2 <= 0 {
			return fmt.Errorf("%w: %f^%f", ErrZeroBase, num1, num2)
		}
		if num1 < 0 && math.Trunc(num2) != num2 {
			return fmt.Errorf("%w: %f^%f", ErrNegativeBase, num1, num2)
		}

		res = math.Pow(num1, num2)
	default:
		return fmt.Errorf("%w: %s", ErrUnknownOperator, op)
	}

	*numStack = append(*numStack, res)

	return nil
}
