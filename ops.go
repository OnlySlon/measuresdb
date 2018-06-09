package main

import (
	"math"
)

func deg2rad(a float64) float64 {
	return (math.Pi / 180) * a
}

func rad2deg(a float64) float64 {
	return a * math.Pi / 180
}

func linear2db(a float64) float64 {
	return 10 * math.Log10(a)
}

func db2linear(a float64) float64 {
	return math.Pow(10, (a / 10))
}

var MagValues [8]float64
var PhaseValues [8]float64

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

func SetMagPhase(idx int, mag float64, phase float64) {
	MagValues[idx] = mag
	PhaseValues[idx] = phase
}

func OpsRegister() {

	var (
		phdelta = &Operator{
			Name:          "phdelta",
			Precedence:    0,
			Associativity: L,
			Args:          2,
			Operation: func(args []float64) float64 {
				if (math.Abs(args[0] - args[1])) < 180 {
					return math.Abs(args[0] - args[1])
				} else {
					return 360 - math.Abs(args[0]-args[1])
				}
			},
		}
		torad = &Operator{
			Name: "	",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return (math.Pi / 180) * args[0]
			},
		}
		sin = &Operator{
			Name:          "Sin",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Sin(args[0])
			},
		}

		sind = &Operator{
			Name:          "SinD",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Sin(deg2rad(args[0]))
			},
		}

		asin = &Operator{
			Name:          "ASin",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Asin(args[0])
			},
		}

		asind = &Operator{
			Name:          "ASinD",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Asin(deg2rad(args[0]))
			},
		}

		cos = &Operator{
			Name:          "Cos",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Cos(args[0])
			},
		}

		cosd = &Operator{
			Name:          "CosD",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Cos(deg2rad(args[0]))
			},
		}

		acosd = &Operator{
			Name:          "ACosD",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Acos(deg2rad(args[0]))
			},
		}
		tan = &Operator{
			Name:          "Tan",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Tan(args[0])
			},
		}
		tand = &Operator{
			Name: "TanD	",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Tan(deg2rad(args[0]))
			},
		}
		atn = &Operator{
			Name:          "Atn",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Atan(args[0])
			},
		}
		atn2 = &Operator{
			Name:          "Atn2",
			Precedence:    0,
			Associativity: L,
			Args:          2,
			Operation: func(args []float64) float64 {
				return math.Atan2(args[0], args[1])
			},
		}
		atnd = &Operator{
			Name:          "AtnD",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Atan(deg2rad(args[0]))
			},
		}
		atn2d = &Operator{
			Name:          "Atn2D",
			Precedence:    0,
			Associativity: L,
			Args:          2,
			Operation: func(args []float64) float64 {
				return math.Atan2(deg2rad(args[0]), deg2rad(args[1]))
			},
		}
		log = &Operator{
			Name:          "Log",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Log2(args[0])
			},
		}
		ln = &Operator{
			Name:          "Ln",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Log(args[0])
			},
		}
		log10 = &Operator{
			Name:          "Log10",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Log10(args[0])
			},
		}
		sqr = &Operator{
			Name:          "Sqr",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Sqrt(args[0])
			},
		}
		mag = &Operator{
			Name:          "Mag",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				//				fmt.Print("get MAG idx=", args[0])
				return MagValues[int(args[0])]
			},
		}
		phd = &Operator{
			Name:          "Phd",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				//				fmt.Print("get PHD idx=", args[0])
				return PhaseValues[int(args[0])]
			},
		}
		abs = &Operator{
			Name:          "abs",
			Precedence:    0,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return math.Abs(args[0])
			},
		}
		sum = &Operator{
			Name:          "sum",
			Precedence:    0,
			Associativity: L,
			Args:          2,
			Operation: func(args []float64) float64 {
				return args[0] + args[1]
			},
		}
	)
	RegisterFunc(phdelta)
	RegisterFunc(torad)
	RegisterFunc(sin)
	RegisterFunc(sind)
	RegisterFunc(asin)
	RegisterFunc(asind)
	RegisterFunc(cos)
	RegisterFunc(cosd)
	RegisterFunc(acosd)
	RegisterFunc(tan)
	RegisterFunc(tand)
	RegisterFunc(atn)
	RegisterFunc(atn2)
	RegisterFunc(atnd)
	RegisterFunc(atn2d)
	RegisterFunc(log)
	RegisterFunc(ln)
	RegisterFunc(log10)
	RegisterFunc(sqr)
	RegisterFunc(mag)
	RegisterFunc(phd)
	RegisterFunc(abs)
	RegisterFunc(sum)

	var (
		add = &Operator{
			Name:          "+",
			Precedence:    1,
			Associativity: L,
			Args:          2,
			Operation: func(args []float64) float64 {
				return args[0] + args[1]
			},
		}
		sub = &Operator{
			Name:          "-",
			Precedence:    1,
			Associativity: L,
			Args:          2,
			Operation: func(args []float64) float64 {
				return args[0] - args[1]
			},
		}
		neg = &Operator{
			Name:          "neg",
			Precedence:    2,
			Associativity: L,
			Args:          1,
			Operation: func(args []float64) float64 {
				return 0 - args[0]
			},
		}
		mul = &Operator{
			Name:          "*",
			Precedence:    2,
			Associativity: L,
			Args:          2,
			Operation: func(args []float64) float64 {
				return args[0] * args[1]
			},
		}
		div = &Operator{
			Name:          "/",
			Precedence:    2,
			Associativity: L,
			Args:          2,
			Operation: func(args []float64) float64 {
				return args[0] / args[1]
			},
		}
		mod = &Operator{
			Name:          "%",
			Precedence:    2,
			Associativity: L,
			Args:          2,
			Operation: func(args []float64) float64 {
				return math.Mod(args[0], args[1])
			},
		}
		pow = &Operator{
			Name:          "^",
			Precedence:    3,
			Associativity: R,
			Args:          2,
			Operation: func(args []float64) float64 {
				return math.Pow(args[0], args[1])
			},
		}
	)
	OpRegister(add)
	OpRegister(sub)
	OpRegister(neg)
	OpRegister(pow)
	OpRegister(mul)
	OpRegister(mod)
	OpRegister(div)

	var (
		cA = &Constant{
			Name:  "A",
			Value: A,
		}
		cB = &Constant{
			Name:  "B",
			Value: B,
		}
		cC = &Constant{
			Name:  "C",
			Value: C,
		}
		cD = &Constant{
			Name:  "D",
			Value: D,
		}
	)

	ConstRegister(cA)
	ConstRegister(cB)
	ConstRegister(cC)
	ConstRegister(cD)
}
