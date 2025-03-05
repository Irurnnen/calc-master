package calc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var operands = "+-*/"

type ASTNodeType int

const (
	NumberType ASTNodeType = iota
	OperationType
)

// Структура для узла AST
type ASTNode struct {
	Type     ASTNodeType // Тип узла: "number" или "operation"
	Value    float64     // Значение (для чисел)
	Operator string      // Оператор (для операций)
	Left     *ASTNode    // Левый потомок
	Right    *ASTNode    // Правый потомок
}

func NewNumberASTNode(value string) (*ASTNode, error) {
	if isNumber(value) {
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, ErrParseFloat
		}
		return &ASTNode{
			Type:  NumberType,
			Value: num,
		}, nil
	}
	return nil, ErrUnknown
}

// Структура для задачи
type Task struct {
	ID        int
	Arg1      float64
	Arg2      float64
	Operation string
}

// Структура для результата задачи
type TaskResult struct {
	ID     int
	Result float64
}

// Структура для выражения
type Expression struct {
	ID     int
	Status string
	Result float64
	AST    *ASTNode
}

// Вспомогательные функции
func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func IsOperator(token string) bool {
	return strings.Contains(operands, token)
}

func ParseNumber(node *ASTNode) float64 {
	if node.Type == NumberType {
		return node.Value
	}
	return 0
}

func InfixToRPN(expr string) (string, error) {
	var output []string
	var stack []string
	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}

	tokens := strings.Fields(expr)
	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if IsOperator(token) {
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return "", errors.New("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		} else {
			return "", fmt.Errorf("invalid token: %s", token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" || stack[len(stack)-1] == ")" {
			return "", errors.New("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return strings.Join(output, " "), nil
}

// Создает AST из RPN
func CreateASTFromRPN(rpn string) (*ASTNode, error) {
	stack := []*ASTNode{}
	tokens := strings.Fields(rpn)

	for _, token := range tokens {
		if isNumber(token) {
			node, err := NewNumberASTNode(token)
			if err != nil {
				return nil, err
			}
			stack = append(stack, node)
		} else if IsOperator(token) {
			if len(stack) < 2 {
				return nil, ErrExtraOperands
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, &ASTNode{
				Type:     OperationType,
				Operator: token,
				Left:     left,
				Right:    right,
			})
		} else {
			return nil, fmt.Errorf("invalid token: %s", token)
		}
	}

	if len(stack) != 1 {
		return nil, errors.New("invalid RPN expression")
	}

	return stack[0], nil
}

// Разделяет AST на задачи
func SplitASTIntoTasks(node *ASTNode) []Task {
	var tasks []Task
	var taskID int

	var traverse func(*ASTNode)
	traverse = func(n *ASTNode) {
		if n.Type == NumberType {
			return // Числа не требуют задач
		}

		// Рекурсивно обходим левую и правую части
		traverse(n.Left)
		traverse(n.Right)

		// Создаем задачу для текущей операции
		taskID++
		task := Task{
			ID:        taskID,
			Arg1:      ParseNumber(n.Left),  // Левый операнд
			Arg2:      ParseNumber(n.Right), // Правый операнд
			Operation: n.Operator,           // Операция
		}
		tasks = append(tasks, task)
	}

	// Начинаем обход с корня AST
	traverse(node)
	return tasks
}

// Обновляет AST с результатом задачи
func UpdateASTWithResult(node *ASTNode, taskID int, result float64) {
	if node == nil {
		return
	}

	if node.Type == OperationType {
		UpdateASTWithResult(node.Left, taskID, result)
		UpdateASTWithResult(node.Right, taskID, result)
	}
}
