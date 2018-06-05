package main

var Constants = map[string]float64{}

type Constant struct {
	Name  string
	Value float64
}

func ConstRegister(c *Constant) {
	Constants[c.Name] = c.Value
}

func IsConstant(str string) bool {
	_, exist := Constants[str]
	return exist
}

func GetValue(str string) float64 {
	val, _ := Constants[str]
	return val
}
