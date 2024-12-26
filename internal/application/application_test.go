package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expectedStatus int
		expectedResult float64
		expectedError  string
	}{
		{"Сложение", `{"expression": "1+2"}`, http.StatusOK, 3, ""},
		{"Вычитание", `{"expression": "5-3"}`, http.StatusOK, 2, ""},
		{"Умножение", `{"expression": "2*3"}`, http.StatusOK, 6, ""},
		{"Деление", `{"expression": "6/2"}`, http.StatusOK, 3, ""},
		{"Скобки", `{"expression": "(2+3)*4"}`, http.StatusOK, 20, ""},
		{"Пробелы", `{"expression": "1 + 2 * 3"}`, http.StatusOK, 7, ""},
		{"Тест со скобками", `{"expression": "2*2*(5+5)"}`, http.StatusOK, 40, ""},
		{"Пустая строка", `{"expression": ""}`, http.StatusUnprocessableEntity, 0, "Expression is not valid"},
		{"Недопустимые символы", `{"expression": "1+a"}`, http.StatusUnprocessableEntity, 0, "Expression is not valid"},
		{"Некорректный JSON", `{"expression": "1+2"`, http.StatusBadRequest, 0, "Invalid request body"},
		{"Деление на ноль", `{"expression": "6/0"}`, http.StatusInternalServerError, 0, "Expression is not valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(tt.expression))
			if err != nil {
				t.Fatalf("Ошибка создания запроса: %v", err)
			}

			rr := httptest.NewRecorder()

			calculateHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Ожидался статус код %v, получен %v", tt.expectedStatus, rr.Code)
			}

			if tt.expectedError != "" {
				var errorResponse ErrorResponse
				err = json.Unmarshal(rr.Body.Bytes(), &errorResponse)
				if err != nil {
					t.Fatalf("Ошибка декодирования ответа: %v", err)
				}
				if errorResponse.Error != tt.expectedError {
					t.Errorf("Ожидалась ошибка '%v', получена '%v'", tt.expectedError, errorResponse.Error)
				}
			} else {
				var successResponse SuccessResponse
				err = json.Unmarshal(rr.Body.Bytes(), &successResponse)
				if err != nil {
					t.Fatalf("Ошибка декодирования ответа: %v", err)
				}
				if successResponse.Result != tt.expectedResult {
					t.Errorf("Ожидался результат %v, получен %v", tt.expectedResult, successResponse.Result)
				}
			}
		})
	}
}

func TestCalculateHandler_MethodNotAllowed(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/calculate", nil)
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}

	rr := httptest.NewRecorder()

	calculateHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Ожидался статус код %v, получен %v", http.StatusMethodNotAllowed, rr.Code)
	}

	var errorResponse ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}
	if errorResponse.Error != "Method not allowed" {
		t.Errorf("Ожидалась ошибка 'Method not allowed', получена '%v'", errorResponse.Error)
	}
}