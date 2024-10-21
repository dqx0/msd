package core_test

import (
	"testing"

	core "github.com/dqx0/msd/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestNewParticle(t *testing.T) {
	particle := core.NewParticle()
	assert.NotNil(t, particle)
}

func TestGetMsd(t *testing.T) {
	particle := core.NewParticle()

	particle.Append(core.ParticlePath{
		X: []float64{1.0, 2.1, 3.0, 4.0, 5.0},
		Y: []float64{2.0, 3.0, 4.0, 5.0, 6.0},
		Z: []float64{1.5, 2.5, 3.5, 4.5, 5.5},
	})

	_ = particle.GetMsd()

	assert.NoError(t, nil)
}
