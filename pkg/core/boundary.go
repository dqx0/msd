package core

import "fmt"

const (
	// 系の長さ
	L = 20.0
)

func (p *ParticlePath) CorrectBoundary() {
	p.X = correct(p.X)
	p.Y = correct(p.Y)
	p.Z = correct(p.Z)
}

func correct(array []float64) []float64 {
	for i := 1; i < len(array); i++ {
		if array[i]-array[i-1] > L-0.5 {
			// マイナスからプラスへのジャンプを補正
			fmt.Println("correct -")
			for j := i; j < len(array); j++ {
				array[j] -= L
			}
		} else if array[i]-array[i-1] < 0.5-L {
			// プラスからマイナスへのジャンプを補正
			fmt.Println("correct +")
			for j := i; j < len(array); j++ {
				array[j] += L
			}
		}
	}
	return array
}
