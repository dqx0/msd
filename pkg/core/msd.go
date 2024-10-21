package core

type IParticle interface {
	Append(particle ParticlePath)
	GetMsd() <-chan Msd
}

type ParticlePath struct {
	X, Y, Z []float64
}

type Msd struct {
	Time int
	MSD  float64
}

func NewParticle() IParticle {
	return &ParticlePath{
		X: []float64{},
		Y: []float64{},
		Z: []float64{},
	}
}

func (p *ParticlePath) Append(particle ParticlePath) {
	p.X = append(p.X, particle.X...)
	p.Y = append(p.Y, particle.Y...)
	p.Z = append(p.Z, particle.Z...)
}

func calculateMSD(p *ParticlePath, t int) float64 {
	n := len(p.X)
	if t >= n {
		return 0.0
	}

	var S_n float64 = 0
	var S_t float64 = 0
	var S_n_t float64 = 0
	var crossTerm float64 = 0

	// 2つの累積和と交差項の計算
	// https://qiita.com/Authns/items/def59166dfd49975e9ba
	for i := 0; i < n; i++ {
		S_n += p.X[i]*p.X[i] + p.Y[i]*p.Y[i] + p.Z[i]*p.Z[i]
		if i < t {
			S_t += p.X[i]*p.X[i] + p.Y[i]*p.Y[i] + p.Z[i]*p.Z[i]
		}
		if i+t < n {
			S_n_t += p.X[i+t]*p.X[i+t] + p.Y[i+t]*p.Y[i+t] + p.Z[i+t]*p.Z[i+t]
			crossTerm += p.X[i]*p.X[i+t] + p.Y[i]*p.Y[i+t] + p.Z[i]*p.Z[i+t]
		}
	}

	// MSDの計算
	MSD := (S_n + S_n_t - S_t - 2*crossTerm) / float64(n-t)

	return MSD
}

// GetMsd は平均二乗変位を計算してチャネルに送信します
func (p *ParticlePath) GetMsd() <-chan Msd {
	ch := make(chan Msd)
	go func() {
		defer close(ch)
		n := len(p.X)
		if n == 0 {
			ch <- Msd{Time: 0, MSD: 0.0}
			return
		}
		for t := 1; t <= n; t++ {
			// 3次元空間でのMSDの計算
			msdTotal := calculateMSD(p, t-1)
			// 時間ステップtとMSDをチャネルに送信
			ch <- Msd{Time: t - 1, MSD: msdTotal}
		}
	}()
	return ch
}
