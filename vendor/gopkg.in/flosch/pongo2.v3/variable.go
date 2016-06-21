package pongo2

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	varTypeInt = iota
	varTypeIdent
)

type variablePart struct {
	typ int
	s   string
	i   int

	is_function_call bool
	calling_args     []functionCallArgument // needed for a function call, represents all argument nodes (INode supports nested function calls)
}

type functionCallArgument interface {
	Evaluate(*ExecutionContext) (*Value, *Error)
}

// TODO: Add location tokens
type stringResolver struct {
	location_token *Token
	val            string
}

type intResolver struct {
	location_token *Token
	val            int
}

type floatResolver struct {
	location_token *Token
	val            float64
}

type boolResolver struct {
	location_token *Token
	val            bool
}

type variableResolver struct {
	location_token *Token

	parts []*variablePart
}

type nodeFilteredVariable struct {
	location_token *Token

	resolver    IEvaluator
	filterChain []*filterCall
}

type nodeVariable struct {
	location_token *Token
	expr           IEvaluator
}

func (expr *nodeFilteredVariable) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	buffer.WriteString(value.String())
	return nil
}

func (expr *variableResolver) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	buffer.WriteString(value.String())
	return nil
}

func (expr *stringResolver) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	buffer.WriteString(value.String())
	return nil
}

func (expr *intResolver) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	buffer.WriteString(value.String())
	return nil
}

func (expr *floatResolver) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	buffer.WriteString(value.String())
	return nil
}

func (expr *boolResolver) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	buffer.WriteString(value.String())
	return nil
}

func (v *nodeFilteredVariable) GetPositionToken() *Token {
	return v.location_token
}

func (v *variableResolver) GetPositionToken() *Token {
	return v.location_token
}

func (v *stringResolver) GetPositionToken() *Token {
	return v.location_token
}

func (v *intResolver) GetPositionToken() *Token {
	return v.location_token
}

func (v *floatResolver) GetPositionToken() *Token {
	return v.location_token
}

func (v *boolResolver) GetPositionToken() *Token {
	return v.location_token
}

func (s *stringResolver) Evaluate(ctx *ExecutionContext) (*Value, *Error) {
	return AsValue(s.val), nil
}

func (i *intResolver) Evaluate(ctx *ExecutionContext) (*Value, *Error) {
	return AsValue(i.val), nil
}

func (f *floatResolver) Evaluate(ctx *ExecutionContext) (*Value, *Error) {
	return AsValue(f.val), nil
}

func (b *boolResolver) Evaluate(ctx *ExecutionContext) (*Value, *Error) {
	return AsValue(b.val), nil
}

func (s *stringResolver) FilterApplied(name string) bool {
	return false
}

func (i *intResolver) FilterApplied(name string) bool {
	return false
}

func (f *floatResolver) FilterApplied(name string) bool {
	return false
}

func (b *boolResolver) FilterApplied(name string) bool {
	return false
}

func (nv *nodeVariable) FilterApplied(name string) bool {
	return nv.expr.FilterApplied(name)
}

func (nv *nodeVariable) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	value, err := nv.expr.Evaluate(ctx)
	if err != nil {
		return err
	}

	if !nv.expr.FilterApplied("safe") && !value.safe && value.IsString() && ctx.Autoescape {
		// apply escape filter
		value, err = filters["escape"](value, nil)
		if err != nil {
			return err
		}
	}

	buffer.WriteString(value.String())
	return nil
}

func (vr *variableResolver) FilterApplied(name string) bool {
	return false
}

func (vr *variableResolver) String() string {
	parts := make([]string, 0, len(vr.parts))
	for _, p := range vr.parts {
		switch p.typ {
		case varTypeInt:
			parts = append(parts, strconv.Itoa(p.i))
		case varTypeIdent:
			parts = append(parts, p.s)
		default:
			panic("unimplemented")
		}
	}
	return strings.Join(parts, ".")
}

func (vr *variableResolver) resolve(ctx *ExecutionContext) (*Value, error) {
	var current reflect.Value
	var is_safe bool

	for idx, part := range vr.parts {
		if idx == 0 {
			// We're looking up the first part of the variable.
			// First we're having a look in our private
			// context (e. g. information provided by tags, like the forloop)
			val, in_private := ctx.Private[vr.parts[0].s]
			if !in_private {
				// Nothing found? Then have a final lookup in the public context
				val = ctx.Public[vr.parts[0].s]
			}
			current = reflect.ValueOf(val) // Get the initial value
		} else {
			// Next parts, resolve it from current

			// Before resolving the pointer, let's see if we have a method to call
			// Problem with resolving the pointer is we're changing the receiver
			is_func := false
			if part.typ == varTypeIdent {
				func_value := current.MethodByName(part.s)
				if func_value.IsValid() {
					current = func_value
					is_func = true
				}
			}

			if !is_func {
				// If current a pointer, resolve it
				if current.Kind() == reflect.Ptr {
					current = current.Elem()
					if !current.IsValid() {
						// Value is not valid (anymore)
						return AsValue(nil), nil
					}
				}

				// Look up which part must be called now
				switch part.typ {
				case varTypeInt:
					// Calling an index is only possible for:
					// * slices/arrays/strings
					switch current.Kind() {
					case reflect.String, reflect.Array, reflect.Slice:
						current = current.Index(part.i)
					default:
						return nil, fmt.Errorf("Can't access an index on type %s (variable %s)",
							current.Kind().String(), vr.String())
					}
				case varTypeIdent:
					// debugging:
					// fmt.Printf("now = %s (kind: %s)\n", part.s, current.Kind().String())

					// Calling a field or key
					switch current.Kind() {
					case reflect.Struct:
						current = current.FieldByName(part.s)
					case reflect.Map:
						current = current.MapIndex(reflect.ValueOf(part.s))
					default:
						return nil, fmt.Errorf("Can't access a field by name on type %s (variable %s)",
							current.Kind().String(), vr.String())
					}
				default:
					panic("unimplemented")
				}
			}
		}

		if !current.IsValid() {
			// Value is not valid (anymore)
			return AsValue(nil), nil
		}

		// If current is a reflect.ValueOf(pongo2.Value), then unpack it
		// Happens in function calls (as a return value) or by injecting
		// into the execution context (e.g. in a for-loop)
		if current.Type() == reflect.TypeOf(&Value{}) {
			tmp_value := current.Interface().(*Value)
			current = tmp_value.val
			is_safe = tmp_value.safe
		}

		// Check whether this is an interface and resolve it where required
		if current.Kind() == reflect.Interface {
			current = reflect.ValueOf(current.Interface())
		}

		// Check if the part is a function call
		if part.is_function_call || current.Kind() == reflect.Func {
			// Check for callable
			if current.Kind() != reflect.Func {
				return nil, fmt.Errorf("'%s' is not a function (it is %s).", vr.String(), current.Kind().String())
			}

			// Check for correct function syntax and types
			// func(*Value, ...) *Value
			t := current.Type()

			// Input arguments
			if len(part.calling_args) != t.NumIn() && !(len(part.calling_args) >= t.NumIn()-1 && t.IsVariadic()) {
				return nil,
					fmt.Errorf("Function input argument count (%d) of '%s' must be equal to the calling argument count (%d).",
						t.NumIn(), vr.String(), len(part.calling_args))
			}

			// Output arguments
			if t.NumOut() != 1 {
				return nil, fmt.Errorf("'%s' must have exactly 1 output argument.", vr.String())
			}

			// Evaluate all parameters
			parameters := make([]reflect.Value, 0)

			num_args := t.NumIn()
			is_variadic := t.IsVariadic()
			var fn_arg reflect.Type

			for idx, arg := range part.calling_args {
				pv, err := arg.Evaluate(ctx)
				if err != nil {
					return nil, err
				}

				if is_variadic {
					if idx >= t.NumIn()-1 {
						fn_arg = t.In(num_args - 1).Elem()
					} else {
						fn_arg = t.In(idx)
					}
				} else {
					fn_arg = t.In(idx)
				}

				if fn_arg != reflect.TypeOf(new(Value)) {
					// Function's argument is not a *pongo2.Value, then we have to check whether input argument is of the same type as the function's argument
					if !is_variadic {
						if fn_arg != reflect.TypeOf(pv.Interface()) && fn_arg.Kind() != reflect.Interface {
							return nil, fmt.Errorf("Function input argument %d of '%s' must be of type %s or *pongo2.Value (not %T).",
								idx, vr.String(), fn_arg.String(), pv.Interface())
						} else {
							// Function's argument has another type, using the interface-value
							parameters = append(parameters, reflect.ValueOf(pv.Interface()))
						}
					} else {
						if fn_arg != reflect.TypeOf(pv.Interface()) && fn_arg.Kind() != reflect.Interface {
							return nil, fmt.Errorf("Function variadic input argument of '%s' must be of type %s or *pongo2.Value (not %T).",
								vr.String(), fn_arg.String(), pv.Interface())
						} else {
							// Function's argument has another type, using the interface-value
							parameters = append(parameters, reflect.ValueOf(pv.Interface()))
						}
					}
				} else {
					// Function's argument is a *pongo2.Value
					parameters = append(parameters, reflect.ValueOf(pv))
				}
			}

			// Call it and get first return parameter back
			rv := current.Call(parameters)[0]

			if rv.Type() != reflect.TypeOf(new(Value)) {
				current = reflect.ValueOf(rv.Interface())
			} else {
				// Return the function call value
				current = rv.Interface().(*Value).val
				is_safe = rv.Interface().(*Value).safe
			}
		}
	}

	if !current.IsValid() {
		// Value is not valid (e. g. NIL value)
		return AsValue(nil), nil
	}

	return &Value{val: current, safe: is_safe}, nil
}

func (vr *variableResolver) Evaluate(ctx *ExecutionContext) (*Value, *Error) {
	value, err := vr.resolve(ctx)
	if err != nil {
		return AsValue(nil), ctx.Error(err.Error(), vr.location_token)
	}
	return value, nil
}

func (v *nodeFilteredVariable) FilterApplied(name string) bool {
	for _, filter := range v.filterChain {
		if filter.name == name {
			return true
		}
	}
	return false
}

func (v *nodeFilteredVariable) Evaluate(ctx *ExecutionContext) (*Value, *Error) {
	value, err := v.resolver.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	for _, filter := range v.filterChain {
		value, err = filter.Execute(value, ctx)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

// IDENT | IDENT.(IDENT|NUMBER)...
func (p *Parser) parseVariableOrLiteral() (IEvaluator, *Error) {
	t := p.Current()

	if t == nil {
		return nil, p.Error("Unexpected EOF, expected a number, string, keyword or identifier.", p.last_token)
	}

	// Is first part a number or a string, there's nothing to resolve (because there's only to return the value then)
	switch t.Typ {
	case TokenNumber:
		p.Consume()

		// One exception to the rule that we don't have float64 literals is at the beginning
		// of an expression (or a variable name). Since we know we started with an integer
		// which can't obviously be a variable name, we can check whether the first number
		// is followed by dot (and then a number again). If so we're converting it to a float64.

		if p.Match(TokenSymbol, ".") != nil {
			// float64
			t2 := p.MatchType(TokenNumber)
			if t2 == nil {
				return nil, p.Error("Expected a number after the '.'.", nil)
			}
			f, err := strconv.ParseFloat(fmt.Sprintf("%s.%s", t.Val, t2.Val), 64)
			if err != nil {
				return nil, p.Error(err.Error(), t)
			}
			fr := &floatResolver{
				location_token: t,
				val:            f,
			}
			return fr, nil
		} else {
			i, err := strconv.Atoi(t.Val)
			if err != nil {
				return nil, p.Error(err.Error(), t)
			}
			nr := &intResolver{
				location_token: t,
				val:            i,
			}
			return nr, nil
		}
	case TokenString:
		p.Consume()
		sr := &stringResolver{
			location_token: t,
			val:            t.Val,
		}
		return sr, nil
	case TokenKeyword:
		p.Consume()
		switch t.Val {
		case "true":
			br := &boolResolver{
				location_token: t,
				val:            true,
			}
			return br, nil
		case "false":
			br := &boolResolver{
				location_token: t,
				val:            false,
			}
			return br, nil
		default:
			return nil, p.Error("This keyword is not allowed here.", nil)
		}
	}

	resolver := &variableResolver{
		location_token: t,
	}

	// First part of a variable MUST be an identifier
	if t.Typ != TokenIdentifier {
		return nil, p.Error("Expected either a number, string, keyword or identifier.", t)
	}

	resolver.parts = append(resolver.parts, &variablePart{
		typ: varTypeIdent,
		s:   t.Val,
	})

	p.Consume() // we consumed the first identifier of the variable name

variableLoop:
	for p.Remaining() > 0 {
		t = p.Current()

		if p.Match(TokenSymbol, ".") != nil {
			// Next variable part (can be either NUMBER or IDENT)
			t2 := p.Current()
			if t2 != nil {
				switch t2.Typ {
				case TokenIdentifier:
					resolver.parts = append(resolver.parts, &variablePart{
						typ: varTypeIdent,
						s:   t2.Val,
					})
					p.Consume() // consume: IDENT
					continue variableLoop
				case TokenNumber:
					i, err := strconv.Atoi(t2.Val)
					if err != nil {
						return nil, p.Error(err.Error(), t2)
					}
					resolver.parts = append(resolver.parts, &variablePart{
						typ: varTypeInt,
						i:   i,
					})
					p.Consume() // consume: NUMBER
					continue variableLoop
				default:
					return nil, p.Error("This token is not allowed within a variable name.", t2)
				}
			} else {
				// EOF
				return nil, p.Error("Unexpected EOF, expected either IDENTIFIER or NUMBER after DOT.",
					p.last_token)
			}
		} else if p.Match(TokenSymbol, "(") != nil {
			// Function call
			// FunctionName '(' Comma-separated list of expressions ')'
			part := resolver.parts[len(resolver.parts)-1]
			part.is_function_call = true
		argumentLoop:
			for {
				if p.Remaining() == 0 {
					return nil, p.Error("Unexpected EOF, expected function call argument list.", p.last_token)
				}

				if p.Peek(TokenSymbol, ")") == nil {
					// No closing bracket, so we're parsing an expression
					expr_arg, err := p.ParseExpression()
					if err != nil {
						return nil, err
					}
					part.calling_args = append(part.calling_args, expr_arg)

					if p.Match(TokenSymbol, ")") != nil {
						// If there's a closing bracket after an expression, we will stop parsing the arguments
						break argumentLoop
					} else {
						// If there's NO closing bracket, there MUST be an comma
						if p.Match(TokenSymbol, ",") == nil {
							return nil, p.Error("Missing comma or closing bracket after argument.", nil)
						}
					}
				} else {
					// We got a closing bracket, so stop parsing arguments
					p.Consume()
					break argumentLoop
				}

			}
			// We're done parsing the function call, next variable part
			continue variableLoop
		}

		// No dot or function call? Then we're done with the variable parsing
		break
	}

	return resolver, nil
}

func (p *Parser) parseVariableOrLiteralWithFilter() (*nodeFilteredVariable, *Error) {
	v := &nodeFilteredVariable{
		location_token: p.Current(),
	}

	// Parse the variable name
	resolver, err := p.parseVariableOrLiteral()
	if err != nil {
		return nil, err
	}
	v.resolver = resolver

	// Parse all the filters
filterLoop:
	for p.Match(TokenSymbol, "|") != nil {
		// Parse one single filter
		filter, err := p.parseFilter()
		if err != nil {
			return nil, err
		}

		// Check sandbox filter restriction
		if _, is_banned := p.template.set.bannedFilters[filter.name]; is_banned {
			return nil, p.Error(fmt.Sprintf("Usage of filter '%s' is not allowed (sandbox restriction active).", filter.name), nil)
		}

		v.filterChain = append(v.filterChain, filter)

		continue filterLoop

		return nil, p.Error("This token is not allowed within a variable.", nil)
	}

	return v, nil
}

func (p *Parser) parseVariableElement() (INode, *Error) {
	node := &nodeVariable{
		location_token: p.Current(),
	}

	p.Consume() // consume '{{'

	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	node.expr = expr

	if p.Match(TokenSymbol, "}}") == nil {
		return nil, p.Error("'}}' expected", nil)
	}

	return node, nil
}
