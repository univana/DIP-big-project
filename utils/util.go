package utils

import "math"

//Min :返回最小值
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

	//计算前将RGB值转化为浮点数
	r := float64(R)
	g := float64(G)
	b := float64(B)

	sum := r + g + b

	I = sum / 3.0

	temp1 := 2.0*r - g - b
	temp2 := 2.0 * math.Sqrt(math.Pow((r-g), 2)+(r-b)*(g-b))
	if temp2 == 0 {
		temp2 = 0.0002
	}
	theta := math.Acos(temp1 / temp2)
	if G >= B {
		H = theta
	} else {
		H = 2.0*math.Pi - theta
	}
	S = 1.0 - 3.0*math.Min(math.Min(r, g), b)/(r+g+b)

	return H, S, I

}

//HSI2RGB :HSI转化为RGB
func HSI2RGB(H float64, S float64, I float64) (R uint8, G uint8, B uint8) {

	if H < (math.Pi * 2.0 / 3.0) {
		/* H<120° */

		B = uint8(I * (1 - S))
		R = uint8(I * (1.0 + (S*math.Cos(H))/math.Cos(math.Pi/3.0-H)))
		G = uint8(I*3) - R - B
	} else if H < (math.Pi * 4.0 / 3.0) {
		/* H<240° */
		H -= math.Pi * 2.0 / 3.0
		R = uint8(I * (1 - S))
		G = uint8(I * (1.0 + (S*math.Cos(H))/math.Cos(math.Pi/3.0-H)))
		B = uint8(I*3) - R - G
	} else {
		H -= math.Pi * 4.0 / 3.0
		G = uint8(I * (1 - S))
		B = uint8(I * (1.0 + (S*math.Cos(H))/math.Cos(math.Pi/3.0-H)))
		R = uint8(I*3) - G - B
	}

	return
}
