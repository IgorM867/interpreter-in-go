package main

import (
	"os"
)

func main() {
	env := createGlobalEnv()

	dat, err := os.ReadFile("test.txt")
	if err != nil {
		panic(err)
	}
	program := produceAst(string(dat))
	program.evaluate(&env)
}
