package calculator

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type Expression interface {
	toS() string
	inspect() interface{}
	isReducible() bool
	reduce() Expression
}

type Number struct {
	value float64
}

func (n Number) toS() string {
	return strconv.FormatFloat(n.value, 'f', 2, 64)
}

func (n Number) inspect() interface{} {
	return n.value
}

func (n Number) isReducible() bool {
	return false
}

func (n Number) reduce() Expression {
	return n
}

type Add struct {
	left, right Expression
}

func (a Add) toS() string {
	return "(" + a.left.toS() + " + " + a.right.toS() + ")"
}

func (a Add) inspect() interface{} {
	return a
}

func (a Add) isReducible() bool {
	return true
}

func (a Add) reduce() Expression {
	if a.left.isReducible() {
		return Add{a.left.reduce(), a.right}
	} else if a.right.isReducible() {
		return Add{a.left, a.right.reduce()}
	} else {
		return Number{a.left.inspect().(float64) + a.right.inspect().(float64)}
	}
}

type Subtract struct {
	left, right Expression
}

func (s Subtract) toS() string {
	return "(" + s.left.toS() + " - " + s.right.toS() + ")"
}

func (s Subtract) inspect() interface{} {
	return s
}

func (s Subtract) isReducible() bool {
	return true
}

func (s Subtract) reduce() Expression {
	if s.left.isReducible() {
		return Subtract{s.left.reduce(), s.right}
	} else if s.right.isReducible() {
		return Subtract{s.left, s.right.reduce()}
	} else {
		return Number{s.left.inspect().(float64) - s.right.inspect().(float64)}
	}
}

type Multiply struct {
	left, right Expression
}

func (m Multiply) toS() string {
	return "(" + m.left.toS() + " * " + m.right.toS() + ")"
}

func (m Multiply) inspect() interface{} {
	return m
}

func (m Multiply) isReducible() bool {
	return true
}

func (m Multiply) reduce() Expression {
	if m.left.isReducible() {
		return Multiply{m.left.reduce(), m.right}
	} else if m.right.isReducible() {
		return Multiply{m.left, m.right.reduce()}
	} else {
		return Number{m.left.inspect().(float64) * m.right.inspect().(float64)}
	}
}

type Divide struct {
	left, right Expression
}

func (d Divide) toS() string {
	return "(" + d.left.toS() + " / " + d.right.toS() + ")"
}

func (d Divide) inspect() interface{} {
	return d
}

func (d Divide) isReducible() bool {
	return true
}

func (d Divide) reduce() Expression {
	if d.left.isReducible() {
		return Divide{d.left.reduce(), d.right}
	} else if d.right.isReducible() {
		return Divide{d.left, d.right.reduce()}
	} else {
		rightValue := d.right.inspect().(float64)
		if rightValue == 0 {
			panic("division by zero")
		}
		return Number{d.left.inspect().(float64) / rightValue}
	}
}

type Machine struct {
	expression Expression
}

func (m *Machine) step() {
	m.expression = m.expression.reduce()
}

func (m *Machine) run() {
	for m.expression.isReducible() {
		m.step()
	}
}

type Parser struct {
	tokens []string
	pos    int
}

func (p *Parser) parseExpression() (Expression, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) && (p.tokens[p.pos] == "+" || p.tokens[p.pos] == "-") {
		op := p.tokens[p.pos]
		p.pos++
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		if op == "+" {
			left = Add{left, right}
		} else {
			left = Subtract{left, right}
		}
	}
	return left, nil
}

func (p *Parser) parseTerm() (Expression, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.tokens) && (p.tokens[p.pos] == "*" || p.tokens[p.pos] == "/") {
		op := p.tokens[p.pos]
		p.pos++
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		if op == "*" {
			left = Multiply{left, right}
		} else {
			left = Divide{left, right}
		}
	}
	return left, nil
}

func (p *Parser) parseFactor() (Expression, error) {
	if p.pos >= len(p.tokens) {
		return nil, errors.New("unexpected end of expression")
	}

	if p.tokens[p.pos] == "(" {
		p.pos++
		exp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		p.pos++
		return exp, nil
	}

	val, err := strconv.ParseFloat(p.tokens[p.pos], 64)
	if err != nil {
		return nil, errors.New("invalid number: " + p.tokens[p.pos])
	}
	p.pos++
	return Number{value: val}, nil
}

func tokenize(input string) []string {
	var tokens []string
	var sb strings.Builder

	for _, r := range input {
		if unicode.IsSpace(r) {
			continue
		}
		if r == '(' || r == ')' || r == '+' || r == '-' || r == '*' || r == '/' {
			if sb.Len() > 0 {
				tokens = append(tokens, sb.String())
				sb.Reset()
			}
			tokens = append(tokens, string(r))
		} else {
			sb.WriteRune(r)
		}
	}

	if sb.Len() > 0 {
		tokens = append(tokens, sb.String())
	}

	return tokens
}

func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	parser := Parser{tokens: tokens, pos: 0}
	input, err := parser.parseExpression()
	if err != nil {
		return 0, err
	}

	machine := Machine{input}
	machine.run()

	return machine.expression.inspect().(float64), nil
}