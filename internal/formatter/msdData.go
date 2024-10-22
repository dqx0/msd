package formatter

import (
	"strconv"

	"github.com/dqx0/msd/pkg/core"
)

// ある時刻における全粒子のMSD
type MsdData struct {
	ParticleIndex []int
	Msd           []float64
}

// 時刻ごとのMSD
type Data struct {
	Time    []int
	MsdData []MsdData
	Average []float64
}

func Create(particles []core.IParticle) [][]interface{} {
	result := Data{
		Time: make([]int, 0),
		MsdData: []MsdData{
			{
				ParticleIndex: make([]int, 0),
				Msd:           make([]float64, 0),
			},
		},
	}

	// MSDを取得
	particleIndex := 0
	timeIndex := 0
	for _, particle := range particles {
		msdOverall := particle.GetMsd()
		timeIndex = 0
		for timeStepMsd := range msdOverall {
			if particleIndex == 0 {
				result.Time = append(result.Time, timeStepMsd.Time)
				result.MsdData = append(result.MsdData, MsdData{
					ParticleIndex: make([]int, 0),
					Msd:           make([]float64, 0),
				})
			}
			result.MsdData[timeIndex].ParticleIndex = append(result.MsdData[timeIndex].ParticleIndex, particleIndex)
			result.MsdData[timeIndex].Msd = append(result.MsdData[timeIndex].Msd, timeStepMsd.MSD)
			timeIndex++
		}
		particleIndex++
	}

	// MSDの平均値を計算
	result.Average = make([]float64, len(result.Time))
	for timeIndex := 0; timeIndex < len(result.Time); timeIndex++ {
		sum := 0.0
		for particleIndex := 0; particleIndex < len(result.MsdData[timeIndex].ParticleIndex); particleIndex++ {
			sum += result.MsdData[timeIndex].Msd[particleIndex]
		}
		result.Average[timeIndex] = sum / float64(len(result.MsdData[timeIndex].ParticleIndex))
	}

	// スプレッドシート用のデータを作成
	values := make([][]interface{}, particleIndex+2)
	values[0] = []interface{}{"Time"}
	for i := 0; i < particleIndex; i++ {
		values[0] = append(values[0], "particle No."+strconv.Itoa(i+1))
	}
	values[0] = append(values[0], "Average")

	for timeIndex := 0; timeIndex < len(result.Time); timeIndex++ {
		row := make([]interface{}, 0, particleIndex+2)
		row = append(row, result.Time[timeIndex]*100)
		for i := 0; i < particleIndex; i++ {
			row = append(row, result.MsdData[timeIndex].Msd[i])
		}
		row = append(row, result.Average[timeIndex])
		values = append(values, row)
	}

	return values
}
