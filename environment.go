package main

import (
	"fmt"
	"os"
)

type Variable struct {
	runtimeVal RuntimeVal
	constant   bool
}

type Env struct {
	parent    *Env
	variables map[string]Variable
}

func createGlobalEnv() Env {
	newEnv := Env{
		parent:    nil,
		variables: make(map[string]Variable),
	}
	newEnv.declareVar("false", BooleanVal{value: false}, true)
	newEnv.declareVar("true", BooleanVal{value: true}, true)
	newEnv.declareVar("null", NullVal{}, true)
	newEnv.declareVar("print", NativeFn{call: nativePrint}, true)
	newEnv.declareVar("println", NativeFn{call: nativePrintln}, true)

	return newEnv
}
func newScope(parent *Env) Env {
	newEnv := Env{
		parent:    parent,
		variables: make(map[string]Variable),
	}
	return newEnv
}
func (env *Env) declareVar(varname string, value RuntimeVal, isConst bool) (RuntimeVal, error) {
	if _, ok := env.variables[varname]; ok {
		return nil, fmt.Errorf("cannot declare variable %s. As it already is defined", varname)
	}

	env.variables[varname] = Variable{runtimeVal: value, constant: isConst}
	return value, nil
}
func (env *Env) assignVar(varname string, value RuntimeVal) RuntimeVal {
	varEnv := env.resolve(varname)

	v := varEnv.variables[varname]
	if v.constant {
		fmt.Println("Cannot reasign constant variable:", varname)
		os.Exit(1)
	}
	varEnv.variables[varname] = Variable{runtimeVal: value, constant: false}

	return value
}
func (env *Env) lookupVar(varname string) RuntimeVal {
	e := env.resolve(varname)

	return e.variables[varname].runtimeVal
}

func (env *Env) resolve(varname string) Env {
	if _, ok := env.variables[varname]; ok {
		return *env
	}
	if env.parent == nil {
		fmt.Printf("Cannot resolve '%s' as it does not exist\n", varname)
		os.Exit(1)
	}

	return env.parent.resolve(varname)
}
