package core

type IParticle interface {
	AppendPath(particle ParticlePath)
	Append(x, y, z float64)
	GetMsd() <-chan Msd
	CorrectBoundary()
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

func (p *ParticlePath) AppendPath(particle ParticlePath) {
	p.X = append(p.X, particle.X...)
	p.Y = append(p.Y, particle.Y...)
	p.Z = append(p.Z, particle.Z...)
}

func (p *ParticlePath) Append(x, y, z float64) {
	p.X = append(p.X, x)
	p.Y = append(p.Y, y)
	p.Z = append(p.Z, z)
}

func (p *ParticlePath) calculateMSD(t int) float64 {
	n := len(p.X)
	if t >= n {
		return 0.0
	}

	// https://qiita.com/Authns/items/def59166dfd49975e9ba
	var sum1 float64 = 0 // Σx(i)² for i=1 to n-t
	var sum2 float64 = 0 // Σx(i)² for i=t+1 to n
	var crossTerm float64 = 0

	for i := 0; i < n-t; i++ {
		sum1 += p.X[i]*p.X[i] + p.Y[i]*p.Y[i] + p.Z[i]*p.Z[i]
	}

	for i := t; i < n; i++ {
		sum2 += p.X[i]*p.X[i] + p.Y[i]*p.Y[i] + p.Z[i]*p.Z[i]
	}

	for i := 0; i < n-t; i++ {
		crossTerm += p.X[i]*p.X[i+t] + p.Y[i]*p.Y[i+t] + p.Z[i]*p.Z[i+t]
	// 2つの累積和と交差項の計算
	// https://qiita.com/Authns/items/def59166dfd49975e9ba
	for i := 0; i < n; i++ {
		S_n += p.X[i]*p.X[i] + p.Y[i]*p.Y[i] + p.Z[i]*p.Z[i]
		if i < t {
			S_t += p.X[i]*p.X[i] + p.Y[i]*p.Y[i] + p.Z[i]*p.Z[i]
		}
		if i+t < n {
			S_n_t += p.X[i]*p.X[i] + p.Y[i]*p.Y[i] + p.Z[i]*p.Z[i]
			crossTerm += p.X[i]*p.X[i+t] + p.Y[i]*p.Y[i+t] + p.Z[i]*p.Z[i+t]
		}
	}

	MSD := (sum1 + sum2 - 2*crossTerm) / float64(n-t)

	return MSD
}

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
			msdTotal := p.calculateMSD(t - 1)
			// 時間ステップtとMSDをチャネルに送信
			ch <- Msd{Time: t - 1, MSD: msdTotal}
		}
	}()
	return ch
}
