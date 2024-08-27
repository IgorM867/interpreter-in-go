package main

import (
	"fmt"
	"main/colors"
)

type RuntimeVal interface {
	getType() string
	String() string
}
type NullVal struct{}
type NumberVal struct {
	value int64
}
type StringVaL struct {
	value string
}
type BooleanVal struct {
	value bool
}
type Object struct {
	properties map[string]RuntimeVal
}
type FunctionCall func(args []RuntimeVal, env *Env) RuntimeVal
type NativeFn struct {
	call FunctionCall
}
type Function struct {
	name           string
	parameters     []string
	declarationEnv *Env
	body           []Stmt
}

func (lhs NumberVal) binaryOperation(operator string, rhs NumberVal) NumberVal {
	switch operator {
	case "+":
		return NumberVal{value: lhs.value + rhs.value}
	case "-":
		return NumberVal{value: lhs.value - rhs.value}
	case "*":
		return NumberVal{value: lhs.value * rhs.value}
	case "/":
		if rhs.value == 0 {
			panic("Cannod divide by 0!")
		}
		return NumberVal{value: lhs.value / rhs.value}
	case "%":
		return NumberVal{value: lhs.value % rhs.value}
	default:
		panic(fmt.Sprintf("invalid operator: %s\n", operator))
	}
}
func (lhs StringVaL) binaryOperation(operator string, rhs StringVaL) StringVaL {
	switch operator {
	case "+":
		return StringVaL{value: lhs.value + rhs.value}
	default:
		panic(fmt.Sprintf("invalid string operatorion: %s\n", operator))
	}
}
func (NullVal) getType() string {
	return "null"
}
func (NumberVal) getType() string {
	return "number"
}
func (StringVaL) getType() string {
	return "string"
}
func (BooleanVal) getType() string {
	return "boolean"
}
func (Object) getType() string {
	return "Object"
}
func (NativeFn) getType() string {
	return "Function"
}
func (Function) getType() string {
	return "Function"
}
func (num NumberVal) String() string {
	return colors.GreenString(fmt.Sprintf("%d", num.value))
}
func (str StringVaL) String() string {
	return colors.YellowString(fmt.Sprintf(`"%s"`, str.value))
}
func (NullVal) String() string {
	return colors.MagentaString("null")
}
func (boolean BooleanVal) String() string {
	return colors.MagentaString(fmt.Sprintf("%v", boolean.value))
}
func (obj Object) String() string {
	str := "{"
	count := len(obj.properties)

	for key, value := range obj.properties {
		count--
		str += fmt.Sprintf(" %v: %v", key, value)
		if count > 0 {
			str += ","
		}
	}

	str += " }"
	return str
}
func (NativeFn) String() string {
	return "[Function]"
}
func (Function) String() string {
	return "[Function]"
}
func compareTypes(val1, val2 RuntimeVal) bool {
	return fmt.Sprintf("%T", val1) == fmt.Sprintf("%T", val2)
}
