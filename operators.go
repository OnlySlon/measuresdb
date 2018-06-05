package main

import (
	"log"
	"strings"
)

var Ops = map[string]*Operator{}

const (
	L = 0
	R = 1
	A = 0
	B = 1
	C = 2
	D = 3
	E = 4
	F = 5
)

type Operator struct {
	Name          string
	Precedence    int
	Associativity int
	Args          int
	Operation     func(args []float64) float64
}

func OpRegister(op *Operator) {
	op.Name = strings.ToLower(op.Name)
	Ops[op.Name] = op
}

func IsOperator(str string) bool {
	str = strings.ToLower(str)
	_, exist := Ops[str]
	return exist
}

func FindOperatorFromString(str string) *Operator {
	str = strings.ToLower(str)
	op, exist := Ops[str]
	if exist {
		return op
	}
	log.Print(Ops)
	return nil
}
