package material

import (
	"math"
	"math/rand"

	"github.com/rrothenb/pbr/pkg/geom"
	"github.com/rrothenb/pbr/pkg/render"
	"github.com/rrothenb/pbr/pkg/rgb"
	"github.com/rrothenb/pbr/pkg/surface"
)

type Grid struct {
	base    surface.Material
	line    surface.Material
	spacing float64
	radius  float64
}

func NewGrid(base, line surface.Material, tiles int, thickness float64) *Grid {
	return &Grid{
		base:    base,
		line:    line,
		spacing: 1.0 / float64(tiles),
		radius:  1.0 / float64(tiles) * thickness,
	}
}

func (g *Grid) At(u, v float64, in, norm geom.Dir, rnd *rand.Rand) (normal geom.Dir, bsdf render.BSDF) {
	du := math.Mod(u, g.spacing)
	dv := math.Mod(v, g.spacing)
	if du < g.radius || dv < g.radius {
		return g.line.At(u, v, in, norm, rnd)
	}
	return g.base.At(u, v, in, norm, rnd)
}

func (g *Grid) Light() rgb.Energy {
	return rgb.Black
}

func (g *Grid) Transmit() rgb.Energy {
	return rgb.Black
}
