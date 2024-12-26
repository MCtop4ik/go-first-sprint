package calculator

import (
	"fmt"
	"strconv"
	"strings"
)

func stringToFloat64(str string) float64 {
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return res
}

func isSign(value rune) bool {
	return value == '+' || value == '-' || value == '*' || value == '/'
}

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	if len(expression) < 3 {
		return 0, fmt.Errorf("Недостаточно символов для вычисления")
	}

	var res float64
	var b string
	var c rune = 0
	var resflag bool = false
	var isc int
	var countc int = 0

	for _, value := range expression {
		if isSign(value) {
			countc++
		}
	}

	if isSign(rune(expression[0])) || isSign(rune(expression[len(expression)-1])) {
		return 0, fmt.Errorf("Некорректное выражение")
	}

	for i, value := range expression {
		if value == '(' {
			isc = i
		}
		if value == ')' {
			calc, err := Calc(expression[isc+1 : i])
			if err != nil {
				return 0, fmt.Errorf("Некорректное выражение")
			}
			calcstr := strconv.FormatFloat(calc, 'f', -1, 64)
			i2 := i
			i -= len(expression[isc:i+1]) - len(calcstr)
			expression = strings.Replace(expression, expression[isc:i2+1], calcstr, 1)
		}
	}

	if countc > 1 {
		for i := 1; i < len(expression); i++ {
			value := rune(expression[i])

			if value == '*' || value == '/' {
				var imin int = i - 1
				if imin != 0 {
					for !isSign(rune(expression[imin])) && imin > 0 {
						imin--
					}
					imin++
				}
				var imax int = i + 1
				if imax == len(expression) {
					imax--
				} else {
					for !isSign(rune(expression[imax])) && imax < len(expression)-1 {
						imax++
					}
				}
				if imax == len(expression)-1 {
					imax++
				}
				calc, err := Calc(expression[imin:imax])
				if err != nil {
					return 0, fmt.Errorf("Некорректное выражение")
				}
				calcstr := strconv.FormatFloat(calc, 'f', -1, 64)
				i -= len(expression[isc:i+1]) - len(calcstr) - 1
				expression = strings.Replace(expression, expression[imin:imax], calcstr, 1)
			}
			if value == '+' || value == '-' || value == '*' || value == '/' {
				c = value
			}
		}
	}

	for _, value := range expression + "s" {
		switch {
		case value == ' ':
			continue
		case (value >= '0' && value <= '9') || value == '.':
			b += string(value)
		case isSign(value) || value == 's':
			if resflag {
				switch c {
				case '+':
					res += stringToFloat64(b)
				case '-':
					res -= stringToFloat64(b)
				case '*':
					res *= stringToFloat64(b)
				case '/':
					divisor := stringToFloat64(b)
					if divisor == 0 {
						return 0, fmt.Errorf("Деление на ноль невозможно")
					}
					res /= divisor
				}
			} else {
				resflag = true
				res = stringToFloat64(b)
			}
			b = strings.ReplaceAll(b, b, "")
			c = value

		case value == 's':
		default:
			return 0, fmt.Errorf("Некорректный ввод")
		}
	}
	return res, nil
}