package calc

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Irurnnen/calc-master/pkg/increment"
)

const disallowedSymbolsRegular = `[^0-9\.+\-*\/()^\s]`
const spacesRegular = `\s`

var operands = "+-*/^"
var priority = map[string]int{"+": 1, "-": 1, "*": 2, "/": 2, "^": 3}

type ASTNodeType int

const (
	TypeNumber ASTNodeType = iota
	TypeOperation
)

type Operation string

const (
	OperationPlus     = "+"
	OperationMinus    = "-"
	OperationMultiply = "*"
	OperationDivision = "/"
)

type ASTNode struct {
	Type     ASTNodeType
	Value    float64
	Operator Operation
	Left     *ASTNode
	Right    *ASTNode
}

func NewNumberASTNode(token string) (*ASTNode, error) {
	num, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return nil, ErrParseFloat
	}
	return &ASTNode{
		Type:  TypeNumber,
		Value: num,
		Left:  &ASTNode{},
		Right: &ASTNode{},
	}, nil
}

type Task struct {
	ID        int
	Arg1      float64
	Arg2      float64
	Operation Operation
}

type TaskResult struct {
	ID     int
	Result float64
}

func ParseExpression(expression string) (*ASTNode, error) {
	// Checking validity of expression
	if err := ValidateExpression(expression); err != nil {
		return nil, err
	}

	// Tokenize expression
	tokens := tokenize(expression)

	// Validate Tokens
	if err := ValidateTokens(tokens); err != nil {
		return nil, err
	}
	// TODO: go to AST
	// Change to postfix
	postfixTokens := ToPostfix(tokens)

	// Calculate the expression
	result, err := parseFromRPNtoAST(postfixTokens)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ValidateExpression(expression string) error {
	// Check disallowed symbols
	re := regexp.MustCompile(disallowedSymbolsRegular)
	if re.MatchString(expression) {
		return ErrExtraCharacters
	}

	// Check correction of brackets
	var bracketBalance int
	for _, v := range expression {
		if v == '(' {
			bracketBalance++
		} else if v == ')' {
			bracketBalance--
			if bracketBalance < 0 {
				return ErrWrongBracketOrder
			}
		}
	}
	if bracketBalance != 0 {
		return ErrUnpairedBracket
	}

	return nil
}

func RemoveSpaces(expression string) string {
	re := regexp.MustCompile(spacesRegular)
	return re.ReplaceAllString(expression, "")
}

func tokenize(expression string) []string {
	var tokens []string
	var number string

	// Delete all space in expression
	expression = RemoveSpaces(expression)

	for _, character := range expression {
		if IsNumber(string(character)) {
			number += string(character)
			continue
		}
		if number != "" {
			tokens = append(tokens, number)
			number = ""
		}
		tokens = append(tokens, string(character))
	}
	if len(number) != 0 {
		tokens = append(tokens, number)
	}
	return tokens
}

// ValidateTokens checks tokens for several errors: ErrEmptyExpression,
// ErrMultipleOperands, ErrMultipleNumbers, ErrExtraOperands
func ValidateTokens(tokens []string) error {
	// Check exists of expression
	if len(tokens) == 0 {
		return ErrEmptyExpression
	}
	// Check multiple operators or multiple numbers
	for i := 1; i < len(tokens); i++ {
		if IsOperand(tokens[i-1]) && IsOperand(tokens[i]) {
			return ErrMultipleOperands
		}
		if IsNumber(tokens[i-1]) && IsNumber(tokens[i]) {
			return ErrMultipleNumbers
		}
	}

	// Check operands at the beginning and end
	if IsOperand(tokens[0]) || IsOperand(tokens[len(tokens)-1]) {
		return ErrExtraOperands
	}

	return nil
}

// IsNumber returns the true if token is a number otherwise false
func IsNumber(token string) bool {
	for _, v := range token {
		if v != '.' && (v < '0' || v > '9') {
			return false
		}

	}
	return true
}

// IsOperand returns the true if token is an operand otherwise false
func IsOperand(token string) bool {
	return strings.Contains(operands, token) && len(token) == 1
}

// To Postfix changes the order of tokens to reverse Polish notation
func ToPostfix(tokens []string) []string {
	var stack []string
	var output []string

	for _, token := range tokens {
		if IsNumber(token) {
			output = append(output, token)
			continue
		}
		switch token {
		case "(":
			stack = append(stack, token)
		case ")":
			for len(stack) != 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) != 0 {
				stack = stack[:len(stack)-1]
			}
		default:
			for len(stack) != 0 && stack[len(stack)-1] != "(" && priority[token] <= priority[stack[len(stack)-1]] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}
	for len(stack) != 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return output
}

// parseFromRPNtoAST solves tokens in Reverse Polish notation. This function return float64
func parseFromRPNtoAST(tokens []string) (*ASTNode, error) {
	stack := []*ASTNode{}
	for _, token := range tokens {
		// If token is number
		if IsNumber(token) {
			node, err := NewNumberASTNode(token)
			if err != nil {
				return nil, err
			}
			stack = append(stack, node)
			// stack = append(stack, num)
			continue
		}
		// If token is operand
		if len(stack) < 2 {
			return nil, ErrExtraOperands
		}

		// Extract left and right Nodes
		left, right := stack[len(stack)-2], stack[len(stack)-1]
		stack = stack[:len(stack)-2]

		// Create new Node
		node := &ASTNode{
			Type:     TypeNumber,
			Value:    0,
			Operator: Operation(token),
			Left:     left,
			Right:    right,
		}
		stack = append(stack, node)
	}

	// Check size of stack
	if len(stack) != 1 {
		return nil, nil
	}

	return stack[0], nil
}

func splitASTIntoTasks(node *ASTNode) ([]Task, map[int]*ASTNode) {
	tasks := []Task{}
	nodeToTaskID := make(map[int]*ASTNode)

	var traverse func(*ASTNode) int
	traverse = func(n *ASTNode) int {
		// Numbers doesn't need to solve
		if n.Type == TypeNumber {
			return -1
		}

		leftID := traverse(n.Left)
		rightID := traverse(n.Right)

		// Create task
		increment.GlobalIncrement.Add()
		task := Task{
			ID:        increment.GlobalIncrement.Get(),
			Arg1:      parseNumber(n.Left),
			Arg2:      parseNumber(n.Right),
			Operation: n.Operator,
		}
		tasks = append(tasks, task)
		nodeToTaskID[increment.GlobalIncrement.Get()] = n

		return increment.GlobalIncrement.Get()
	}
}

func parseNumber(node *ASTNode) float64 {
	if node.Type == TypeNumber {
		return node.Value
	}
	return 0
}
