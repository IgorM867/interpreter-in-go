package main

import (
	"fmt"
	"os"
)

type Stmt interface {
	evaluate(env *Env) RuntimeVal
}
type Expr interface {
	evaluate(env *Env) RuntimeVal
}

type Program struct {
	body []Stmt
}
type VarDeclaration struct {
	constant   bool
	identifier string
	value      Expr
}
type FunctionDeclaration struct {
	parameters []string
	name       string
	body       []Stmt
}
type IfStmt struct {
	condition   Expr
	body        []Stmt
	alternative []Stmt
}
type WhileStmt struct {
	condition Expr
	body      []Stmt
}
type AssigmentExpr struct {
	assigne Expr
	value   Expr
}
type Property struct {
	key   string
	value Expr
}
type ObjectLiteral struct {
	properties []Property
}
type BooleanExpr struct {
	left     Expr
	right    Expr
	operator string
}
type CallExpr struct {
	args   []Expr
	caller Expr
}
type MemberExpr struct {
	object   Expr
	property Expr
	computed bool
}
type BinaryExpr struct {
	left     Expr
	right    Expr
	operator string
}
type UnaryExpression struct {
	operator string
	operand  Expr
}
type Identifier struct {
	symbol string
}
type NumericLiteral struct {
	value int64
}
type StringLiteral struct {
	value string
}
type ArrayLiteral struct {
	elements []Expr
}

func (p Program) evaluate(env *Env) RuntimeVal {
	var lastEvaluated RuntimeVal = NullVal{}
	for _, stmt := range p.body {
		lastEvaluated = stmt.evaluate(env)
	}
	return lastEvaluated
}
func (v VarDeclaration) evaluate(env *Env) RuntimeVal {
	var value RuntimeVal = NullVal{}
	if v.value != nil {
		value = v.value.evaluate(env)
	}
	val, err := env.declareVar(v.identifier, value, v.constant)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return val
}
func (f FunctionDeclaration) evaluate(env *Env) RuntimeVal {
	fn := Function{
		name:           f.name,
		parameters:     f.parameters,
		declarationEnv: env,
		body:           f.body,
	}
	val, err := env.declareVar(f.name, fn, true)
	if err != nil {
		panic(err)
	}
	return val
}
func (i IfStmt) evaluate(env *Env) RuntimeVal {
	condition, ok := i.condition.evaluate(env).(BooleanVal)
	if !ok {
		fmt.Println("If condition must be a boolean.")
	}

	if condition.value {
		scope := newScope(env)
		for _, stmt := range i.body {
			stmt.evaluate(&scope)
		}
	} else {
		scope := newScope(env)
		for _, stmt := range i.alternative {
			stmt.evaluate(&scope)
		}
	}

	return NullVal{}
}
func (w WhileStmt) evaluate(env *Env) RuntimeVal {
	condition, ok := w.condition.evaluate(env).(BooleanVal)
	if !ok {
		fmt.Println("While condition must be a boolean.")
	}

	for condition.value {
		scope := newScope(env)
		for _, stmt := range w.body {
			stmt.evaluate(&scope)
		}

		condition = w.condition.evaluate(env).(BooleanVal)
	}
	return NullVal{}
}
func (a AssigmentExpr) evaluate(env *Env) RuntimeVal {
	ident, ok := a.assigne.(Identifier)
	if !ok {
		println("Invalid LHS inside assigment expression", a.assigne)
		os.Exit(1)
	}

	return env.assignVar(ident.symbol, a.value.evaluate(env))
}
func (o ObjectLiteral) evaluate(env *Env) RuntimeVal {
	properties := make(map[string]RuntimeVal)

	for _, p := range o.properties {
		var value RuntimeVal

		if p.value == nil {
			value = env.lookupVar(p.key)
		} else {
			value = p.value.evaluate(env)
		}

		properties[p.key] = value
	}

	return Object{properties: properties}
}

func (c CallExpr) evaluate(env *Env) RuntimeVal {

	args := make([]RuntimeVal, len(c.args))

	for i, a := range c.args {
		args[i] = a.evaluate(env)
	}
	function := c.caller.evaluate(env)
	nativeFn, ok := function.(NativeFn)
	if ok {
		return nativeFn.call(args, env)
	}
	fn, ok := function.(Function)
	if ok {
		scope := newScope(fn.declarationEnv)

		if len(args) < len(fn.parameters) {
			fmt.Printf("Function %v expects %v arguments and got only %v\n", fn.name, len(args), len(fn.parameters))
		}
		for i, param := range fn.parameters {
			scope.declareVar(param, args[i], false)
		}

		var result RuntimeVal = NullVal{}

		for _, stmt := range fn.body {
			result = stmt.evaluate(&scope)
		}
		return result
	}

	println("Cannot call value that is not a function.")
	os.Exit(1)
	panic("Unreachable code")
}
func (m MemberExpr) evaluate(env *Env) RuntimeVal {
	obj := m.object.evaluate(env)

	switch obj := obj.(type) {
	case Object:
		if !m.computed {
			propName, ok := m.property.(Identifier)
			if !ok {
				fmt.Println("Invalid property")
				os.Exit(1)
			}
			prop, ok := obj.properties[propName.symbol]
			if !ok {
				fmt.Printf("Propety %v does not exist\n", propName.symbol)
				os.Exit(1)
			}

			return prop
		}

		propName, ok := m.property.evaluate(env).(StringVaL)
		if !ok {
			fmt.Printf("Object property must me of type string. %v\n", propName)
			os.Exit(1)
		}
		prop, ok := obj.properties[propName.value]
		if !ok {
			fmt.Printf("Propety %v does not exist\n", propName.value)
			os.Exit(1)
		}
		return prop
	case Array:
		if !m.computed {
			panic("To get array element you need to use []")
		}
		prop := m.property.evaluate(env)
		index, ok := prop.(NumberVal)
		if !ok {
			panic(fmt.Sprintf("Expected number as an array index and get: %v", prop.getType()))
		}
		if index.value >= int64(len(obj.elements)) {
			panic(fmt.Sprintf("Array index out of bounds. Attempted to access index %v in an array of size %v.", index.value, len(obj.elements)))
		}
		return obj.elements[index.value]

	default:
		panic(fmt.Sprintf("Unsuported member expression: %v is not an object or array\n", obj.getType()))
	}
}
func (u UnaryExpression) evaluate(env *Env) RuntimeVal {

	switch u.operator {
	case "!":
		operand := u.operand.evaluate(env)
		boolean, ok := operand.(BooleanVal)
		if !ok {
			fmt.Printf("invalid operation: operator ! not defined on type %s\n", operand.getType())
			os.Exit(1)
		}
		return BooleanVal{value: !boolean.value}
	default:
		panic(fmt.Sprintf("Not implementet evaluation for this operator: %v\n", u.operator))

	}
}
func (b BooleanExpr) evaluate(env *Env) RuntimeVal {
	lhs := b.left.evaluate(env)
	rhs := b.right.evaluate(env)

	if !compareTypes(lhs, rhs) {
		fmt.Printf("invalid operation: %v %v %v (mismatched types %v and %v)\n", lhs, b.operator, rhs, lhs.getType(), rhs.getType())
		os.Exit(1)
	}

	switch b.operator {
	case "==":
		switch lhs := lhs.(type) {
		case NullVal:
			return BooleanVal{value: true}
		case NumberVal:
			return BooleanVal{value: lhs.value == rhs.(NumberVal).value}
		case StringVaL:
			return BooleanVal{value: lhs.value == rhs.(StringVaL).value}
		case BooleanVal:
			return BooleanVal{value: lhs.value == rhs.(BooleanVal).value}
		default:
			fmt.Printf("This operations is not supported on this type (%s %s %s)", lhs.getType(), b.operator, rhs.getType())
			os.Exit(1)
		}
	case "!=":
		switch lhs := lhs.(type) {
		case NullVal:
			return BooleanVal{value: true}
		case NumberVal:
			return BooleanVal{value: lhs.value != rhs.(NumberVal).value}
		case StringVaL:
			return BooleanVal{value: lhs.value != rhs.(StringVaL).value}
		case BooleanVal:
			return BooleanVal{value: lhs.value != rhs.(BooleanVal).value}
		default:
			fmt.Printf("This operations is not supported on this type (%s %s %s)", lhs.getType(), b.operator, rhs.getType())
			os.Exit(1)
		}
	case ">":
		switch lhs := lhs.(type) {
		case NumberVal:
			return BooleanVal{value: lhs.value > rhs.(NumberVal).value}
		default:
			fmt.Printf("This operations is not supported on this type (%s %s %s)", lhs.getType(), b.operator, rhs.getType())
			os.Exit(1)
		}
	case "<":
		switch lhs := lhs.(type) {
		case NumberVal:
			return BooleanVal{value: lhs.value < rhs.(NumberVal).value}
		default:
			fmt.Printf("This operations is not supported on this type (%s %s %s)", lhs.getType(), b.operator, rhs.getType())
			os.Exit(1)
		}
	case "<=":
		switch lhs := lhs.(type) {
		case NumberVal:
			return BooleanVal{value: lhs.value <= rhs.(NumberVal).value}
		default:
			fmt.Printf("This operations is not supported on this type (%s %s %s)", lhs.getType(), b.operator, rhs.getType())
			os.Exit(1)
		}
	case ">=":
		switch lhs := lhs.(type) {
		case NumberVal:
			return BooleanVal{value: lhs.value >= rhs.(NumberVal).value}
		default:
			fmt.Printf("This operations is not supported on this type (%s %s %s)", lhs.getType(), b.operator, rhs.getType())
			os.Exit(1)
		}
	default:
		fmt.Printf("Invalid operator: %s", b.operator)
		os.Exit(1)
	}

	panic("Unreachable code")
}
func (b BinaryExpr) evaluate(env *Env) RuntimeVal {
	lhs := b.left.evaluate(env)
	rhs := b.right.evaluate(env)

	if !compareTypes(lhs, rhs) {
		fmt.Printf("invalid operation: %v %v %v (mismatched types %v and %v)\n", lhs, b.operator, rhs, lhs.getType(), rhs.getType())
		os.Exit(1)
	}

	switch lhs := lhs.(type) {
	case NumberVal:
		return lhs.binaryOperation(b.operator, rhs.(NumberVal))
	case StringVaL:
		return lhs.binaryOperation(b.operator, rhs.(StringVaL))
	}

	panic(fmt.Sprintf("unsuportet operation: %v %v %v\n", lhs, b.operator, rhs))
}
func (i Identifier) evaluate(env *Env) RuntimeVal {
	val := env.lookupVar(i.symbol)

	return val
}
func (n NumericLiteral) evaluate(_ *Env) RuntimeVal {
	return NumberVal(n)
}
func (s StringLiteral) evaluate(env *Env) RuntimeVal {
	return StringVaL(s)
}
func (a ArrayLiteral) evaluate(env *Env) RuntimeVal {
	elemements := make([]RuntimeVal, len(a.elements))

	for i, elem := range a.elements {
		elemements[i] = elem.evaluate(env)
	}

	return Array{elemements}
}

func (p Program) String() string {
	str := ""
	for _, v := range p.body {
		str += fmt.Sprintf("  %+v \n", v)
	}

	return fmt.Sprintf("Program \n body:[\n%s]", str)
}
func (v VarDeclaration) String() string {
	return fmt.Sprintf("VarDeclaration{identifier: %v, value: %v, constant:%v}", v.identifier, v.value, v.constant)

}
func (a AssigmentExpr) String() string {
	return fmt.Sprintf("AssigmentExpr{assigne: %+v,value: %+v}", a.assigne, a.value)
}
func (m MemberExpr) String() string {
	return fmt.Sprintf("MemberExpr{object:%v, property:%v, computed:%v}", m.object, m.property, m.computed)
}
func (o ObjectLiteral) String() string {
	if len(o.properties) == 0 {
		return "ObjectLiteral{}"
	}
	str := "ObjectLiteral{\n"
	for _, p := range o.properties {
		str += fmt.Sprintf("   %+v\n", p)
	}
	str += "  }"
	return str
}
func (b BinaryExpr) String() string {
	return fmt.Sprintf("BinaryExpr{left:%v right:%v operator:'%v'}", b.left, b.right, b.operator)
}

func (i Identifier) String() string {
	return fmt.Sprintf("Identifier{symbol:'%v'}", i.symbol)
}
func (n NumericLiteral) String() string {
	return fmt.Sprintf("%v", n.value)
}
