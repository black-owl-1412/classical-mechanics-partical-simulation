package main

import (
	"image/color"
	"log"
	"physics-sim/particles"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	width  = 800
	height = 600
)

type Game struct {
	dt         float64
	Particales []*particles.Particle
}

func (g *Game) Update() error {
	n := len(g.Particales)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {

			particles.Applyforce(g.Particales[j], g.Particales[i], g.dt)
		}

	}
	for i := 0; i < len(g.Particales); i++ {
		for j := i + 1; j < len(g.Particales); j++ {
			particles.CollidePair(g.Particales[i], g.Particales[j])
		}
	}

	for i := range g.Particales {
		p := g.Particales[i]
		p.Position[0] += p.Momentum[0] / p.Mass * g.dt
		p.Position[1] += p.Momentum[1] / p.Mass * g.dt
	}

	for i := range g.Particales {
		p := g.Particales[i]

		particles.Keepinframe(p, width, height)

	}

	return nil

}
func (g *Game) Draw(screen *ebiten.Image) {
	n := len(g.Particales)
	for i := 0; i < n; i++ {
		x := g.Particales[i].Position[0]
		y := g.Particales[i].Position[1]
		r := g.Particales[i].Radius
		color := g.Particales[i].Color

		vector.DrawFilledCircle(screen, float32(x), float32(y), float32(r), color, true)
	}

}
func (g *Game) Layout(outsidewidth, outsideheight int) (int, int) {
	return width, height
}

func main() {

	g := &Game{}
	g.dt = 1. / 60

	samplesparticles := particles.GenerateParticles(50, 0, width, 0, height)
	orange := color.RGBA{R: 255, G: 127, B: 0, A: 255}
	samplesparticles[0].Color = orange

	g.Particales = samplesparticles
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("simulator")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
