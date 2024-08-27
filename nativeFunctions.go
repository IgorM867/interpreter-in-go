package main

import "fmt"

func nativePrint(args []RuntimeVal, env *Env) RuntimeVal {

	for _, arg := range args {
		fmt.Printf("%v ", arg)
	}

	return NullVal{}
}
func nativePrintln(args []RuntimeVal, env *Env) RuntimeVal {

	for _, arg := range args {
		fmt.Printf("%v ", arg)
	}

	fmt.Print("\n")

	return NullVal{}
}
