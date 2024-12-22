package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"/Users/senya/Desktop/go-course-first-sprint/pkg/calculator"
)

type RequestBody struct {
	Expression string `json:"expression"`
}

type SuccessResponse struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func isValidExpression(expression string) bool {
	match, _ := regexp.MatchString(`^[0-9+\-*/(). ]+$`, expression)
	return match
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Method not allowed"})
		return
	}

	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	if !isValidExpression(requestBody.Expression) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Expression is not valid"})
		return
	}

	result, err := calc(requestBody.Expression)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{Result: result})
}

// Функция calc - заглушка, предполагается, что она реализована где-то в другом месте
func calc(expression string) (string, error) {
	answer, err := calculator.Calc(expression)
	return answer, err
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}