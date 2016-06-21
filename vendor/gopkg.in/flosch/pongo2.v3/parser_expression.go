package pongo2

import (
	"bytes"
	"fmt"
	"math"
)

type Expression struct {
	// TODO: Add location token?
	expr1    IEvaluator
	expr2    IEvaluator
	op_token *Token
}

type relationalExpression struct {
	// TODO: Add location token?
	expr1    IEvaluator
	expr2    IEvaluator
	op_token *Token
}

type simpleExpression struct {
	negate        bool
	negative_sign bool
	term1         IEvaluator
	term2         IEvaluator
	op_token      *Token
}

type term struct {
	// TODO: Add location token?
	factor1  IEvaluator
	factor2  IEvaluator
	op_token *Token
}

type power struct {
	// TODO: Add location token?
	power1 IEvaluator
	power2 IEvaluator
}

func (expr *Expression) FilterApplied(name string) bool {
	return expr.expr1.FilterApplied(name) && (expr.expr2 == nil ||
		(expr.expr2 != nil && expr.expr2.FilterApplied(name)))
}

func (expr *relationalExpression) FilterApplied(name string) bool {
	return expr.expr1.FilterApplied(name) && (expr.expr2 == nil ||
		(expr.expr2 != nil && expr.expr2.FilterApplied(name)))
}

func (expr *simpleExpression) FilterApplied(name string) bool {
	return expr.term1.FilterApplied(name) && (expr.term2 == nil ||
		(expr.term2 != nil && expr.term2.FilterApplied(name)))
}

func (t *term) FilterApplied(name string) bool {
	return t.factor1.FilterApplied(name) && (t.factor2 == nil ||
		(t.factor2 != nil && t.factor2.FilterApplied(name)))
}

func (p *power) FilterApplied(name string) bool {
	return p.power1.FilterApplied(name) && (p.power2 == nil ||
		(p.power2 != nil && p.power2.FilterApplied(name)))
}

func (expr *Expression) GetPositionToken() *Token {
	return expr.expr1.GetPositionToken()
}

func (expr *relationalExpression) GetPositionToken() *Token {
	return expr.expr1.GetPositionToken()
}

func (expr *simpleExpression) GetPositionToken() *Token {
	return expr.term1.GetPositionToken()
}

func (expr *term) GetPositionToken() *Token {
	return expr.factor1.GetPositionToken()
}

func (expr *power) GetPositionToken() *Token {
	return expr.power1.GetPositionToken()
}

func (expr *Expression) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	buffer.WriteString(value.String())
	return nil
}

func (expr *relationalExpression) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	buffer.WriteString(value.String())
	return nil
}

func (expr *simpleExpression) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	buffer.WriteString(value.String())
	return nil
}

func (expr *term) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	buffer.WriteString(value.String())
	return nil
}

func (expr *power) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	buffer.WriteString(value.String())
	return nil
}

func (expr *Expression) Evaluate(ctx *ExecutionContext) (*Value, *Error) {
	v1, err := expr.expr1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if expr.expr2 != nil {
		v2, err := expr.expr2.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		switch expr.op_token.Val {
		case "and", "&&":
			return AsValue(v1.IsTrue() && v2.IsTrue()), nil
		case "or", "||":
			return AsValue(v1.IsTrue() || v2.IsTrue()), nil
		default:
			panic(fmt.Sprintf("unimplemented: %s", expr.op_token.Val))
		}
	} else {
		return v1, nil
	}
}

func (expr *relationalExpression) Evaluate(ctx *ExecutionContext) (*Value, *Error) {
	v1, err := expr.expr1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if expr.expr2 != nil {
		v2, err := expr.expr2.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		switch expr.op_token.Val {
		case "<=":
			if v1.IsFloat() || v2.IsFloat() {
				return AsValue(v1.Float() <= v2.Float()), nil
			} else {
				return AsValue(v1.Integer() <= v2.Integer()), nil
			}
		case ">=":
			if v1.IsFloat() || v2.IsFloat() {
				return AsValue(v1.Float() >= v2.Float()), nil
			} else {
				return AsValue(v1.Integer() >= v2.Integer()), nil
			}
		case "==":
			return AsValue(v1.EqualValueTo(v2)), nil
		case ">":
			if v1.IsFloat() || v2.IsFloat() {
				return AsValue(v1.Float() > v2.Float()), nil
			} else {
				return AsValue(v1.Integer() > v2.Integer()), nil
			}
		case "<":
			if v1.IsFloat() || v2.IsFloat() {
				return AsValue(v1.Float() < v2.Float()), nil
			} else {
				return AsValue(v1.Integer() < v2.Integer()), nil
			}
		case "!=", "<>":
			return AsValue(!v1.EqualValueTo(v2)), nil
		case "in":
			return AsValue(v2.Contains(v1)), nil
		default:
			panic(fmt.Sprintf("unimplemented: %s", expr.op_token.Val))
		}
	} else {
		return v1, nil
	}
}

func (expr *simpleExpression) Evaluate(ctx *ExecutionContext) (*Value, *Error) {
	t1, err := expr.term1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	result := t1

	if expr.negate {
		result = result.Negate()
	}

	if expr.negative_sign {
		if result.IsNumber() {
			switch {
			case result.IsFloat():
				result = AsValue(-1 * result.Float())
			case result.IsInteger():
				result = AsValue(-1 * result.Integer())
			default:
				panic("not possible")
			}
		} else {
			return nil, ctx.Error("Negative sign on a non-number expression", expr.GetPositionToken())
		}
	}

	if expr.term2 != nil {
		t2, err := expr.term2.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		switch expr.op_token.Val {
		case "+":
			if result.IsFloat() || t2.IsFloat() {
				// Result will be a float
				return AsValue(result.Float() + t2.Float()), nil
			} else {
				// Result will be an integer
				return AsValue(result.Integer() + t2.Integer()), nil
			}
		case "-":
			if result.IsFloat() || t2.IsFloat() {
				// Result will be a float
				return AsValue(result.Float() - t2.Float()), nil
			} else {
				// Result will be an integer
				return AsValue(result.Integer() - t2.Integer()), nil
			}
		default:
			panic("unimplemented")
		}
	}

	return result, nil
}

func (t *term) Evaluate(ctx *ExecutionContext) (*Value, *Error) {
	f1, err := t.factor1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if t.factor2 != nil {
		f2, err := t.factor2.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		switch t.op_token.Val {
		case "*":
			if f1.IsFloat() || f2.IsFloat() {
				// Result will be float
				return AsValue(f1.Float() * f2.Float()), nil
			}
			// Result will be int
			return AsValue(f1.Integer() * f2.Integer()), nil
		case "/":
			if f1.IsFloat() || f2.IsFloat() {
				// Result will be float
				return AsValue(f1.Float() / f2.Float()), nil
			}
			// Result will be int
			return AsValue(f1.Integer() / f2.Integer()), nil
		case "%":
			// Result will be int
			return AsValue(f1.Integer() % f2.Integer()), nil
		default:
			panic("unimplemented")
		}
	} else {
		return f1, nil
	}
}

func (pw *power) Evaluate(ctx *ExecutionContext) (*Value, *Error) {
	p1, err := pw.power1.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if pw.power2 != nil {
		p2, err := pw.power2.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		return AsValue(math.Pow(p1.Float(), p2.Float())), nil
	} else {
		return p1, nil
	}
}

func (p *Parser) parseFactor() (IEvaluator, *Error) {
	if p.Match(TokenSymbol, "(") != nil {
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		if p.Match(TokenSymbol, ")") == nil {
			return nil, p.Error("Closing bracket expected after expression", nil)
		}
		return expr, nil
	}

	return p.parseVariableOrLiteralWithFilter()
}

func (p *Parser) parsePower() (IEvaluator, *Error) {
	pw := new(power)

	power1, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	pw.power1 = power1

	if p.Match(TokenSymbol, "^") != nil {
		power2, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		pw.power2 = power2
	}

	if pw.power2 == nil {
		// Shortcut for faster evaluation
		return pw.power1, nil
	}

	return pw, nil
}

func (p *Parser) parseTerm() (IEvaluator, *Error) {
	return_term := new(term)

	factor1, err := p.parsePower()
	if err != nil {
		return nil, err
	}
	return_term.factor1 = factor1

	for p.PeekOne(TokenSymbol, "*", "/", "%") != nil {
		if return_term.op_token != nil {
			// Create new sub-term
			return_term = &term{
				factor1: return_term,
			}
		}

		op := p.Current()
		p.Consume()

		factor2, err := p.parsePower()
		if err != nil {
			return nil, err
		}

		return_term.op_token = op
		return_term.factor2 = factor2
	}

	if return_term.op_token == nil {
		// Shortcut for faster evaluation
		return return_term.factor1, nil
	}

	return return_term, nil
}

func (p *Parser) parseSimpleExpression() (IEvaluator, *Error) {
	expr := new(simpleExpression)

	if sign := p.MatchOne(TokenSymbol, "+", "-"); sign != nil {
		if sign.Val == "-" {
			expr.negative_sign = true
		}
	}

	if p.Match(TokenSymbol, "!") != nil || p.Match(TokenKeyword, "not") != nil {
		expr.negate = true
	}

	term1, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	expr.term1 = term1

	for p.PeekOne(TokenSymbol, "+", "-") != nil {
		if expr.op_token != nil {
			// New sub expr
			expr = &simpleExpression{
				term1: expr,
			}
		}

		op := p.Current()
		p.Consume()

		term2, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		expr.term2 = term2
		expr.op_token = op
	}

	if expr.negate == false && expr.negative_sign == false && expr.term2 == nil {
		// Shortcut for faster evaluation
		return expr.term1, nil
	}

	return expr, nil
}

func (p *Parser) parseRelationalExpression() (IEvaluator, *Error) {
	expr1, err := p.parseSimpleExpression()
	if err != nil {
		return nil, err
	}

	expr := &relationalExpression{
		expr1: expr1,
	}

	if t := p.MatchOne(TokenSymbol, "==", "<=", ">=", "!=", "<>", ">", "<"); t != nil {
		expr2, err := p.parseRelationalExpression()
		if err != nil {
			return nil, err
		}
		expr.op_token = t
		expr.expr2 = expr2
	} else if t := p.MatchOne(TokenKeyword, "in"); t != nil {
		expr2, err := p.parseSimpleExpression()
		if err != nil {
			return nil, err
		}
		expr.op_token = t
		expr.expr2 = expr2
	}

	if expr.expr2 == nil {
		// Shortcut for faster evaluation
		return expr.expr1, nil
	}

	return expr, nil
}

func (p *Parser) ParseExpression() (IEvaluator, *Error) {
	rexpr1, err := p.parseRelationalExpression()
	if err != nil {
		return nil, err
	}

	exp := &Expression{
		expr1: rexpr1,
	}

	if p.PeekOne(TokenSymbol, "&&", "||") != nil || p.PeekOne(TokenKeyword, "and", "or") != nil {
		op := p.Current()
		p.Consume()
		expr2, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		exp.expr2 = expr2
		exp.op_token = op
	}

	if exp.expr2 == nil {
		// Shortcut for faster evaluation
		return exp.expr1, nil
	}

	return exp, nil
}
