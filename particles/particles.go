package particles

import (
	"image/color"
	"math"
	"math/rand/v2"
)

var G float64 = 120

type Particle struct {
	Mass     float64    // the mass of the particale
	Position [2]float64 // the position vector
	Momentum [2]float64 // the momentum vector
	Radius   float64    // the radius of the particale
	Color    color.Color
}

func GenerateParticles(n int, minX, maxX, minY, maxY float64) []*Particle {
	particles := make([]*Particle, n)
	for i := 0; i < n; i++ {
		x := minX + rand.Float64()*(maxX-minX)
		y := minY + rand.Float64()*(maxY-minY)
		mass := 1.0 + rand.Float64()*4.0 // mass between 1-5

		particles[i] = &Particle{
			Mass:     mass,
			Position: [2]float64{x, y},
			Momentum: [2]float64{0, 0},
			Radius:   2.5,
			Color:    color.White,
		}
	}
	return particles
}

// Applyforce computes the mutual attraction between p1 and p2
// and *adds* the impulse to each particle's momentum.
// NO position integration happens here – do that once per
// particle after all forces are accumulated.
func Applyforce(p1, p2 *Particle, dt float64) {
	dx := p1.Position[0] - p2.Position[0]
	dy := p1.Position[1] - p2.Position[1]
	r2 := dx*dx + dy*dy

	// softening (avoid singularity)
	const ε2 = 4.0
	r3 := math.Pow(r2+ε2, 1.5)

	// magnitude of force (scalar)
	f := G * p1.Mass * p2.Mass / r3

	// vector components
	fx := f * dx
	fy := f * dy

	// impulse = F*dt
	dpx := fx * dt
	dpy := fy * dt

	// Newton-3: equal and opposite
	p1.Momentum[0] -= dpx
	p1.Momentum[1] -= dpy
	p2.Momentum[0] += dpx
	p2.Momentum[1] += dpy
}

func Keepinframe(p *Particle, w, h float64) {
	wallDamp := 0.90
	if p.Position[0] <= 0 {
		p.Position[0] = 0
		p.Momentum[0] = -p.Momentum[0] * wallDamp
	} else if p.Position[0] >= w {
		p.Position[0] = w
		p.Momentum[0] = -p.Momentum[0] * wallDamp
	}
	// top / bottom  (remember Y grows downward)
	if p.Position[1] <= 0 {
		p.Position[1] = 0
		p.Momentum[1] = -p.Momentum[1] * wallDamp
	} else if p.Position[1] >= h {
		p.Position[1] = h
		p.Momentum[1] = -p.Momentum[1] * wallDamp
	}
}

// collidePair performs an elastic collision between two disks
// (no friction, no spin, equal mass assumed for simplicity).
func CollidePair(p1, p2 *Particle) {
	dx := p2.Position[0] - p1.Position[0]
	dy := p2.Position[1] - p1.Position[1]
	r := p1.Radius + p2.Radius
	d2 := dx*dx + dy*dy

	if d2 >= r*r { // not touching
		return
	}

	// unit normal
	d := math.Sqrt(d2)
	nx := dx / d
	ny := dy / d

	// relative velocity
	dvx := p2.Momentum[0]/p2.Mass - p1.Momentum[0]/p1.Mass
	dvy := p2.Momentum[1]/p2.Mass - p1.Momentum[1]/p1.Mass
	vn := dvx*nx + dvy*ny // velocity along normal

	if vn > 0 { // separating – ignore
		return
	}

	// impulse magnitude (elastic, 1D formula)
	j := 2 * vn / (1/p1.Mass + 1/p2.Mass)

	// apply impulse
	jx := j * nx
	jy := j * ny

	p1.Momentum[0] += jx * p1.Mass
	p1.Momentum[1] += jy * p1.Mass
	p2.Momentum[0] -= jx * p2.Mass
	p2.Momentum[1] -= jy * p2.Mass

	// separate them to avoid overlap (push along normal)
	overlap := r - d
	split := overlap * 0.5 // half each way
	p1.Position[0] -= nx * split
	p1.Position[1] -= ny * split
	p2.Position[0] += nx * split
	p2.Position[1] += ny * split
}
