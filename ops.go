package main

import (
	"math"

	"github.com/alfredxing/calc/operators"
	"github.com/alfredxing/calc/operators/functions"
)

func deg2rad(a float64) float64 {
	return (math.Pi / 180) * a
}

func rad2deg(a float64) float64 {
	return a * math.Pi / 180
}

/*
func ConstRegister() {

	var (
		cPh1 = &constants.Constant{
			Name:  "PhA",
			Value: 0,
		}
	)

}
*/

func OpsRegister() {

	var (
		phdelta = &operators.Operator{
			Name:          "phdelta",
			Precedence:    0,
			Associativity: operators.L,
			Args:          2,
			Operation: func(args []float64) float64 {
				if (math.Abs(args[0] - args[1])) < 180 {
					return math.Abs(args[0] - args[1])
				} else {
					return 360 - math.Abs(args[0]-args[1])
				}
			},
		}
		torad = &operators.Operator{
			Name:          "torad",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return (math.Pi / 180) * args[0]
			},
		}
		sin = &operators.Operator{
			Name:          "Sin",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Sin(args[0])
			},
		}

		sind = &operators.Operator{
			Name:          "SinD",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Sin(deg2rad(args[0]))
			},
		}

		asin = &operators.Operator{
			Name:          "ASin",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Asin(args[0])
			},
		}

		asind = &operators.Operator{
			Name:          "ASinD",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Asin(deg2rad(args[0]))
			},
		}

		cos = &operators.Operator{
			Name:          "Cos",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Cos(args[0])
			},
		}

		cosd = &operators.Operator{
			Name:          "CosD",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Cos(deg2rad(args[0]))
			},
		}

		acosd = &operators.Operator{
			Name:          "ACosD",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Acos(deg2rad(args[0]))
			},
		}
		tan = &operators.Operator{
			Name:          "Tan",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Tan(args[0])
			},
		}
		tand = &operators.Operator{
			Name:          "Tan",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Tan(deg2rad(args[0]))
			},
		}
		atn = &operators.Operator{
			Name:          "Atn",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Atan(args[0])
			},
		}
		atn2 = &operators.Operator{
			Name:          "Atn2",
			Precedence:    0,
			Associativity: operators.L,
			Args:          2,
			Operation: func(args []float64) float64 {
				return math.Atan2(args[0], args[1])
			},
		}
		atnd = &operators.Operator{
			Name:          "AtnD",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Atan(deg2rad(args[0]))
			},
		}
		atn2d = &operators.Operator{
			Name:          "Atn2D",
			Precedence:    0,
			Associativity: operators.L,
			Args:          2,
			Operation: func(args []float64) float64 {
				return math.Atan2(deg2rad(args[0]), deg2rad(args[1]))
			},
		}
		log = &operators.Operator{
			Name:          "Log",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Log2(args[0])
			},
		}
		ln = &operators.Operator{
			Name:          "Ln",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Log(args[0])
			},
		}
		log10 = &operators.Operator{
			Name:          "Log10",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Log10(args[0])
			},
		}
		sqr = &operators.Operator{
			Name:          "Sqr",
			Precedence:    0,
			Associativity: operators.L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Sqrt(args[0])
			},
		}
	)
	functions.Register(phdelta)
	functions.Register(torad)
	functions.Register(sin)
	functions.Register(sind)
	functions.Register(asin)
	functions.Register(asind)
	functions.Register(cos)
	functions.Register(cosd)
	functions.Register(acosd)
	functions.Register(tan)
	functions.Register(tand)
	functions.Register(atn)
	functions.Register(atn2)
	functions.Register(atnd)
	functions.Register(atn2d)
	functions.Register(log)
	functions.Register(ln)
	functions.Register(log10)
	functions.Register(sqr)

}
