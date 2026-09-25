package calc

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/shopspring/decimal"
)

const divisionPrecision = 10

type tokenType int

const (
	tokenNumber tokenType = iota
	tokenPlus
	tokenMinus
	tokenStar
	tokenSlash
	tokenLParen
	tokenRParen
	tokenEOF
)

type token struct {
	typ tokenType
	val string
}

func tokenize(expr string) ([]token, error) {
	var tokens []token
	runes := []rune(expr)
	i := 0
	for i < len(runes) {
		c := runes[i]
		switch {
		case unicode.IsSpace(c):
			i++
		case c == '+':
			tokens = append(tokens, token{tokenPlus, "+"})
			i++
		case c == '-':
			tokens = append(tokens, token{tokenMinus, "-"})
			i++
		case c == '*':
			tokens = append(tokens, token{tokenStar, "*"})
			i++
		case c == '/':
			tokens = append(tokens, token{tokenSlash, "/"})
			i++
		case c == '(':
			tokens = append(tokens, token{tokenLParen, "("})
			i++
		case c == ')':
			tokens = append(tokens, token{tokenRParen, ")"})
			i++
		case unicode.IsDigit(c) || c == '.':
			start := i
			seenDot := false
			for i < len(runes) && (unicode.IsDigit(runes[i]) || runes[i] == '.') {
				if runes[i] == '.' {
					if seenDot {
						return nil, fmt.Errorf("invalid number near %q: unexpected extra '.'", string(runes[start:i+1]))
					}
					seenDot = true
				}
				i++
			}
			numStr := string(runes[start:i])
			if numStr == "." || strings.HasSuffix(numStr, ".") {
				return nil, fmt.Errorf("invalid number: %q", numStr)
			}
			tokens = append(tokens, token{tokenNumber, numStr})
		default:
			return nil, fmt.Errorf("invalid character: %q", c)
		}
	}
	tokens = append(tokens, token{tokenEOF, ""})
	return tokens, nil
}

type parser struct {
	tokens []token
	pos    int
}

func (p *parser) current() token {
	return p.tokens[p.pos]
}

func (p *parser) advance() token {
	t := p.tokens[p.pos]
	if p.pos < len(p.tokens)-1 {
		p.pos++
	}
	return t
}

// expression := term (('+' | '-') term)*
func (p *parser) parseExpression() (decimal.Decimal, error) {
	left, err := p.parseTerm()
	if err != nil {
		return decimal.Zero, err
	}
	for {
		switch p.current().typ {
		case tokenPlus:
			p.advance()
			right, err := p.parseTerm()
			if err != nil {
				return decimal.Zero, err
			}
			left = left.Add(right)
		case tokenMinus:
			p.advance()
			right, err := p.parseTerm()
			if err != nil {
				return decimal.Zero, err
			}
			left = left.Sub(right)
		default:
			return left, nil
		}
	}
}

// term := factor (('*' | '/') factor)*
func (p *parser) parseTerm() (decimal.Decimal, error) {
	left, err := p.parseFactor()
	if err != nil {
		return decimal.Zero, err
	}
	for {
		switch p.current().typ {
		case tokenStar:
			p.advance()
			right, err := p.parseFactor()
			if err != nil {
				return decimal.Zero, err
			}
			left = left.Mul(right)
		case tokenSlash:
			p.advance()
			right, err := p.parseFactor()
			if err != nil {
				return decimal.Zero, err
			}
			if right.IsZero() {
				return decimal.Zero, errors.New("division by zero")
			}
			left = left.DivRound(right, divisionPrecision)
		default:
			return left, nil
		}
	}
}

// factor := NUMBER | '(' expression ')' | ('-' | '+') factor
func (p *parser) parseFactor() (decimal.Decimal, error) {
	tok := p.current()
	switch tok.typ {
	case tokenPlus:
		p.advance()
		return p.parseFactor()
	case tokenMinus:
		p.advance()
		val, err := p.parseFactor()
		if err != nil {
			return decimal.Zero, err
		}
		return val.Neg(), nil
	case tokenNumber:
		p.advance()
		val, err := decimal.NewFromString(tok.val)
		if err != nil {
			return decimal.Zero, fmt.Errorf("invalid number: %q", tok.val)
		}
		return val, nil
	case tokenLParen:
		p.advance()
		val, err := p.parseExpression()
		if err != nil {
			return decimal.Zero, err
		}
		if p.current().typ != tokenRParen {
			return decimal.Zero, errors.New("missing closing parenthesis")
		}
		p.advance()
		return val, nil
	default:
		return decimal.Zero, errors.New("unexpected end of expression")
	}
}

// Evaluate parses and evaluates a basic arithmetic expression (+, -, *, /,
// parentheses, decimal numbers) using exact decimal arithmetic. Division is
// rounded to divisionPrecision decimal places. The result is returned as a
// canonical string with trailing fractional zeros trimmed.
func Evaluate(expr string) (string, error) {
	if strings.TrimSpace(expr) == "" {
		return "", errors.New("expression must not be empty")
	}

	tokens, err := tokenize(expr)
	if err != nil {
		return "", err
	}

	p := &parser{tokens: tokens}
	result, err := p.parseExpression()
	if err != nil {
		return "", err
	}
	if p.current().typ != tokenEOF {
		return "", fmt.Errorf("unexpected token: %q", p.current().val)
	}

	return trimTrailingZeros(result.String()), nil
}

func trimTrailingZeros(s string) string {
	if !strings.Contains(s, ".") {
		return s
	}
	s = strings.TrimRight(s, "0")
	return strings.TrimRight(s, ".")
}
