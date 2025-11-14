package operator

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Expression 表达式接口
type Expression interface {
	Evaluate(ctx *ExecutionContext) (interface{}, error)
}

// 表达式类型
const (
	ExprTypeLiteral   = "literal"
	ExprTypeVariable  = "variable"
	ExprTypeBinary    = "binary"
	ExprTypeUnary     = "unary"
)

// BinaryOperator 二元操作符
const (
	OpEQ = "=="
	OpNE = "!="
	OpGT = ">"
	OpGE = ">="
	OpLT = "<"
	OpLE = "<="
	OpAnd = "&&"
	OpOr = "||"
)

// UnaryOperator 一元操作符
const (
	OpNot = "!"
)

// LiteralExpression 字面量表达式
type LiteralExpression struct {
	Value interface{}
}

func (e *LiteralExpression) Evaluate(ctx *ExecutionContext) (interface{}, error) {
	return e.Value, nil
}

// VariableExpression 变量表达式
type VariableExpression struct {
	Name string
}

func (e *VariableExpression) Evaluate(ctx *ExecutionContext) (interface{}, error) {
	val := ctx.GetVariable(e.Name)
	if val == nil {
		return nil, fmt.Errorf("变量未定义: %s", e.Name)
	}
	return val, nil
}

// BinaryExpression 二元表达式
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (e *BinaryExpression) Evaluate(ctx *ExecutionContext) (interface{}, error) {
	leftVal, err := e.Left.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	rightVal, err := e.Right.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	switch e.Operator {
	case OpEQ:
		return equal(leftVal, rightVal), nil
	case OpNE:
		return !equal(leftVal, rightVal), nil
	case OpGT:
		return greaterThan(leftVal, rightVal), nil
	case OpGE:
		return greaterThanOrEqual(leftVal, rightVal), nil
	case OpLT:
		return lessThan(leftVal, rightVal), nil
	case OpLE:
		return lessThanOrEqual(leftVal, rightVal), nil
	case OpAnd:
		return logicalAnd(leftVal, rightVal), nil
	case OpOr:
		return logicalOr(leftVal, rightVal), nil
	default:
		return nil, fmt.Errorf("不支持的操作符: %s", e.Operator)
	}
}

// UnaryExpression 一元表达式
type UnaryExpression struct {
	Operator string
	Operand  Expression
}

func (e *UnaryExpression) Evaluate(ctx *ExecutionContext) (interface{}, error) {
	val, err := e.Operand.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	switch e.Operator {
	case OpNot:
		return logicalNot(val), nil
	default:
		return nil, fmt.Errorf("不支持的一元操作符: %s", e.Operator)
	}
}

// Parser 表达式解析器
type Parser struct {
	tokens []Token
	pos    int
}

// Token 词法单元
type Token struct {
	Type  TokenType
	Value string
}

// TokenType 词法单元类型
type TokenType int

const (
	TokenNumber TokenType = iota
	TokenString
	TokenIdentifier
	TokenOperator
	TokenParenLeft
	TokenParenRight
	TokenEOF
)

// NewParser 创建新的解析器
func NewParser(expression string) *Parser {
	tokens := tokenize(expression)
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// Parse 解析表达式
func (p *Parser) Parse() (Expression, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("表达式为空")
	}
	return p.parseExpression()
}

// parseExpression 解析表达式
func (p *Parser) parseExpression() (Expression, error) {
	return p.parseLogicalOr()
}

// parseLogicalOr 解析逻辑或表达式
func (p *Parser) parseLogicalOr() (Expression, error) {
	left, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}

	for p.match(OpOr) {
		operator := p.previous().Value
		right, err := p.parseLogicalAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Operator: operator, Right: right}
	}

	return left, nil
}

// parseLogicalAnd 解析逻辑与表达式
func (p *Parser) parseLogicalAnd() (Expression, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.match(OpAnd) {
		operator := p.previous().Value
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{Left: left, Operator: operator, Right: right}
	}

	return left, nil
}

// parseComparison 解析比较表达式
func (p *Parser) parseComparison() (Expression, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	operators := []string{OpEQ, OpNE, OpGT, OpGE, OpLT, OpLE}
	for _, op := range operators {
		if p.match(op) {
			operator := p.previous().Value
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			left = &BinaryExpression{Left: left, Operator: operator, Right: right}
		}
	}

	return left, nil
}

// parseTerm 解析项
func (p *Parser) parseTerm() (Expression, error) {
	if p.match(OpNot) {
		operator := p.previous().Value
		operand, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{Operator: operator, Operand: operand}, nil
	}

	if p.matchType(TokenParenLeft) {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if !p.matchType(TokenParenRight) {
			return nil, fmt.Errorf("期望右括号")
		}
		return expr, nil
	}

	return p.parsePrimary()
}

// parsePrimary 解析主表达式
func (p *Parser) parsePrimary() (Expression, error) {
	if p.matchType(TokenNumber) {
		val := p.previous().Value
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}
		return &LiteralExpression{Value: num}, nil
	}

	if p.matchType(TokenString) {
		val := p.previous().Value
		return &LiteralExpression{Value: strings.Trim(val, "'\"")}, nil
	}

	if p.matchType(TokenIdentifier) {
		name := p.previous().Value
		return &VariableExpression{Name: name}, nil
	}

	return nil, fmt.Errorf("期望表达式")
}

// match 匹配特定操作符
func (p *Parser) match(operator string) bool {
	if p.peek().Type == TokenOperator && p.peek().Value == operator {
		p.advance()
		return true
	}
	return false
}

// matchType 匹配特定类型
func (p *Parser) matchType(tokenType TokenType) bool {
	if p.peek().Type == tokenType {
		p.advance()
		return true
	}
	return false
}

// peek 查看当前token
func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF, Value: ""}
	}
	return p.tokens[p.pos]
}

// advance 前进到下一个token
func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

// previous 获取前一个token
func (p *Parser) previous() Token {
	if p.pos > 0 {
		return p.tokens[p.pos-1]
	}
	return Token{Type: TokenEOF, Value: ""}
}

// tokenize 词法分析
func tokenize(expression string) []Token {
	var tokens []Token
	var current strings.Builder
	i := 0

	for i < len(expression) {
		ch := rune(expression[i])

		// 跳过空白字符
		if unicode.IsSpace(ch) {
			i++
			continue
		}

		// 处理数字
		if unicode.IsDigit(ch) {
			current.Reset()
			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || expression[i] == '.') {
				current.WriteByte(expression[i])
				i++
			}
			tokens = append(tokens, Token{Type: TokenNumber, Value: current.String()})
			continue
		}

		// 处理字符串
		if ch == '\'' || ch == '"' {
			current.Reset()
			quote := ch
			i++
			for i < len(expression) && rune(expression[i]) != quote {
				current.WriteByte(expression[i])
				i++
			}
			i++
			tokens = append(tokens, Token{Type: TokenString, Value: string(quote) + current.String() + string(quote)})
			continue
		}

		// 处理标识符
		if unicode.IsLetter(ch) || ch == '_' {
			current.Reset()
			for i < len(expression) && (unicode.IsLetter(rune(expression[i])) || unicode.IsDigit(rune(expression[i])) || expression[i] == '_') {
				current.WriteByte(expression[i])
				i++
			}
			tokens = append(tokens, Token{Type: TokenIdentifier, Value: current.String()})
			continue
		}

		// 处理操作符
		if strings.Contains("==!><=&&||", string(ch)) {
			current.Reset()
			current.WriteRune(ch)
			
			// 检查多字符操作符
			if i+1 < len(expression) {
				nextCh := rune(expression[i+1])
				if strings.Contains("=!<>&|", string(nextCh)) && strings.Contains("==!><=&&||", string(ch)+string(nextCh)) {
					current.WriteRune(nextCh)
					i++
				}
			}
			i++
			tokens = append(tokens, Token{Type: TokenOperator, Value: current.String()})
			continue
		}

		// 处理括号
		if ch == '(' {
			tokens = append(tokens, Token{Type: TokenParenLeft, Value: "("})
			i++
			continue
		}
		if ch == ')' {
			tokens = append(tokens, Token{Type: TokenParenRight, Value: ")"})
			i++
			continue
		}

		i++
	}

	return tokens
}

// 比较函数
func equal(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func greaterThan(a, b interface{}) bool {
	numA, okA := toNumber(a)
	numB, okB := toNumber(b)
	if okA && okB {
		return numA > numB
	}
	return false
}

func greaterThanOrEqual(a, b interface{}) bool {
	numA, okA := toNumber(a)
	numB, okB := toNumber(b)
	if okA && okB {
		return numA >= numB
	}
	return false
}

func lessThan(a, b interface{}) bool {
	numA, okA := toNumber(a)
	numB, okB := toNumber(b)
	if okA && okB {
		return numA < numB
	}
	return false
}

func lessThanOrEqual(a, b interface{}) bool {
	numA, okA := toNumber(a)
	numB, okB := toNumber(b)
	if okA && okB {
		return numA <= numB
	}
	return false
}

// 逻辑函数
func logicalAnd(a, b interface{}) bool {
	boolA := toBoolean(a)
	boolB := toBoolean(b)
	return boolA && boolB
}

func logicalOr(a, b interface{}) bool {
	boolA := toBoolean(a)
	boolB := toBoolean(b)
	return boolA || boolB
}

func logicalNot(a interface{}) bool {
	return !toBoolean(a)
}

// 类型转换函数
func toNumber(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case float64:
		return v, true
	case string:
		num, err := strconv.ParseFloat(v, 64)
		return num, err == nil
	default:
		return 0, false
	}
}

func toBoolean(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0
	case string:
		return len(v) > 0
	default:
		return false
	}
}

// EvaluateExpression 评估表达式
func EvaluateExpression(expr string, ctx *ExecutionContext) (interface{}, error) {
	parser := NewParser(expr)
	expression, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return expression.Evaluate(ctx)
}

// EvaluateBoolean 评估布尔表达式
func EvaluateBoolean(expr string, ctx *ExecutionContext) (bool, error) {
	result, err := EvaluateExpression(expr, ctx)
	if err != nil {
		return false, err
	}
	return toBoolean(result), nil
}