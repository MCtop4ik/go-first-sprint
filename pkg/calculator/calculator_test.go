package calculator

import (
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		err      bool
	}{
		{"Сложение", "1+2", 3, false},
		{"Вычитание", "5-3", 2, false},
		{"Умножение", "2*3", 6, false},
		{"Деление", "6/2", 3, false},
		{"Деление на ноль", "6/0", 0, true},
		{"Сложное выражение", "2+3*4", 14, false},
		{"Скобки", "(2+3)*4", 20, false},
		{"Дробные числа", "1.5+2.5", 4, false},
		{"Пробелы", "1 + 2 * 3", 7, false},
		{"Пробелы2", "(1 + 2) * 3", 9, false},
		{"Пустая строка", "", 0, true},
		{"Недостаточно символов", "1+", 0, true},
		{"Некорректный символ", "1+a", 0, true},
		{"Некорректное выражение", "+1+2", 0, true},
		{"Некорректное выражение", "1+2-", 0, true},
		{"Тест со скобками", "2*2*(5+5)", 40, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.err {
						t.Errorf("Произошла паника: %v", r)
					}
				}
			}()

			result, err := Calc(tt.input)
			if tt.err {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но получен результат: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Ошибка: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Ожидалось %v, получено %v", tt.expected, result)
				}
			}
		})
	}
}