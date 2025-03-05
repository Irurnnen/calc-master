package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/Irurnnen/calc-master/pkg/calc"
)

var (
	expressions      = make(map[int]*calc.Expression) // Хранилище выражений
	expressionsMutex = &sync.Mutex{}                  // Мьютекс для потокобезопасного доступа
	tasks            = make(map[int]*calc.Task)       // Хранилище задач
	tasksMutex       = &sync.Mutex{}                  // Мьютекс для задач
	currentID        = 0                              // Уникальный ID для выражений
)

// Обработчик для добавления выражения
func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	// Преобразуем выражение в RPN
	rpn, err := calc.InfixToRPN(req.Expression)
	if err != nil {
		http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
		return
	}

	// Создаем AST из RPN
	ast, err := calc.CreateASTFromRPN(rpn)
	if err != nil {
		http.Error(w, "Failed to create AST", http.StatusInternalServerError)
		return
	}

	// Создаем новое выражение
	expressionsMutex.Lock()
	currentID++
	expr := &calc.Expression{
		ID:     currentID,
		Status: "pending",
		Result: 0,
		AST:    ast,
	}
	expressions[currentID] = expr
	expressionsMutex.Unlock()

	// Разделяем AST на задачи
	tasks := calc.SplitASTIntoTasks(ast)

	// Сохраняем задачи
	tasksMutex.Lock()
	for _, task := range tasks {
		tasks[task.ID] = task
	}
	tasksMutex.Unlock()

	// Возвращаем ID выражения
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": currentID})
}

// Обработчик для получения списка выражений
func HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	expressionsMutex.Lock()
	defer expressionsMutex.Unlock()

	var exprs []*calc.Expression
	for _, expr := range expressions {
		exprs = append(exprs, expr)
	}

	json.NewEncoder(w).Encode(map[string][]*calc.Expression{"expressions": exprs})
}

// Обработчик для получения выражения по ID
func HandleGetExpressionByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	expressionsMutex.Lock()
	defer expressionsMutex.Unlock()

	expr, exists := expressions[id]
	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]*calc.Expression{"expression": expr})
}

// Обработчик для получения задачи
func HandleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tasksMutex.Lock()
		defer tasksMutex.Unlock()

		for _, task := range tasks {
			json.NewEncoder(w).Encode(map[string]*calc.Task{"task": task})
			return
		}

		http.Error(w, "No tasks available", http.StatusNotFound)
	} else if r.Method == http.MethodPost {
		var req struct {
			ID     int     `json:"id"`
			Result float64 `json:"result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
			return
		}

		tasksMutex.Lock()
		defer tasksMutex.Unlock()

		task, exists := tasks[req.ID]
		if !exists {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		// Обновляем результат задачи
		expressionsMutex.Lock()
		for _, expr := range expressions {
			if expr.AST != nil {
				calc.UpdateASTWithResult(expr.AST, task.ID, req.Result)
			}
		}
		expressionsMutex.Unlock()

		// Удаляем задачу
		delete(tasks, req.ID)

		w.WriteHeader(http.StatusOK)
	}
}
