package utils

import "math"

func Min(x uint8, y uint8, z uint8) uint8 {
	min := x
	if y < min {
		min = y
	}
	if z < min {
		min = z
	}
	return min
}

//RGB2HSI :RGB模型转化为HSI模型
func RGB2HSI(R uint8, G uint8, B uint8) (H float64, S float64, I float64) {
	temp1 := float64(2*R - G - B)
	temp2 := 2.0 * math.Sqrt(math.Pow(float64(R-G), 2)+float64((R-B)*(G-B)))
	if temp2 == 0 {
		temp2 = 0.02
	}
	theta := math.Acos(temp1 / temp2)
	if G >= B {
		H = theta
	} else {
		H = 2*math.Pi - theta
	}
	S = 1.0 - 3.0*float64(Min(R, G, B))/float64(R+G+B)

	I = float64(R+G+B) / 3.0
	return H, S, I

}
