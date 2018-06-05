package main

import (
	"strings"
)

var Names = map[string]bool{}

func RegisterFunc(op *Operator) {
	op.Name = strings.ToLower(op.Name)
	//	log.Print("Register function: " + op.Name)
	OpRegister(op)
	Names[op.Name] = true
}

func IsFunction(str string) bool {
	str = strings.ToLower(str)
	//	log.Print("is Function: " + str)
	//	log.Print(Names)
	return Names[str]
}
