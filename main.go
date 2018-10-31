package main

import (
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/png"
	"os"
)

const MINI_WIDTH = 400
const MINI_HEIGHT = 400
const MINI_ROWS = 5
const MINI_COLS = 5
const GUTTER = 20

type Environment struct {
	gravity int
}

type Particle struct {
	point image.Point
	vel_x int
	vel_y int
	size  int
}

// TODO: DO BOUNDS CHECKING OF PARTICLE WITHIN A MINI_BOX
// IF OUTSIDE OF MINI BOUND, delete? the particle
func updateParticles(env Environment, particles []Particle, timestep int) {
	for i := range particles {
		particles[i].vel_y += env.gravity
		particles[i].point.X += particles[i].vel_x * timestep
		particles[i].point.Y += particles[i].vel_y * timestep
	}
}

type Box struct {
	min_x int
	min_y int
	max_x int
	max_y int
}

func getBoundingBox(position int) Box {
	box := Box{}
	row_idx := position / MINI_COLS
	col_idx := position % MINI_COLS

	box.min_x = (MINI_WIDTH+GUTTER)*col_idx + GUTTER
	box.max_x = box.min_x + MINI_WIDTH
	box.min_y = (MINI_HEIGHT+GUTTER)*row_idx + GUTTER
	box.max_y = box.min_y + MINI_HEIGHT

	//fmt.Printf("%s %s\n", position, box)
	return box
}

func renderMiniImage(img *image.RGBA, position int, particles []Particle) {
	bb := getBoundingBox(position)
	for _, p := range particles {
		for x := 0; x < p.size; x++ {
			for y := 0; y < p.size; y++ {
				img.Set(bb.min_x+p.point.X+x,
					bb.min_y+p.point.Y+y,
					color.RGBA{255, 0, 0, 255})
			}
		}
	}
}

func drawBorder(img *image.RGBA, position int, color color.RGBA) {
	bb := getBoundingBox(position)
	for x := bb.min_x; x < bb.max_x; x++ {
		img.Set(x, bb.min_y, color)
		img.Set(x, bb.max_y, color)
	}

	for y := bb.min_y; y < bb.max_y; y++ {
		img.Set(bb.min_x, y, color)
		img.Set(bb.max_x, y, color)
	}
}

func writeImage(filename string, img image.Image) {
	out, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	png.Encode(out, img)

}

func main() {
	fmt.Println("Hello World!")
	env := Environment{gravity: 1}
	particles := make([]Particle, 1)
	particles[0] = Particle{image.Point{5, 5}, 50, 0, 5}

	total_width := (MINI_WIDTH+GUTTER)*MINI_COLS + GUTTER
	total_height := (MINI_HEIGHT+GUTTER)*MINI_ROWS + GUTTER
	fmt.Printf("%u %u\n", total_width, total_height)
	img := image.NewRGBA(image.Rect(0, 0, total_width, total_height))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	step := 1
	for t := 0; t < MINI_COLS*MINI_ROWS; t += step {
		drawBorder(img, t, color.RGBA{0, 0, 100, 255})
		updateParticles(env, particles, step)
		renderMiniImage(img, t, particles)
	}

	filename := "out.png"

	writeImage(filename, img)
}
