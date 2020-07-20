package main

import (
	"fmt"
	"os"
)

type State struct {
	x, y, vx, vy float64
}

func (state *State) next(G, dt float64) {
	f1 := state.x + state.y
	f2 := state.y - state.x

	f1sq := f1 * f1 * f1 * f1
	f2sq := f2 * f2 * f1 * f1

	if f1 < 0 {
		f1sq *= -1
	}
	if f2 < 0 {
		f2sq *= -1
	}

	ad1 := G * f1sq * dt
	ad2 := G * f2sq * dt

	a1x, a1y := -ad1, -ad1
	a2x, a2y := ad2, -ad2

	state.x += state.vx * dt
	state.y += state.vy * dt

	state.vx += a1x + a2x
	state.vy += a1y + a2y
}

func main() {
	s := State{30, 0, 0, 1}

	dst, _ := os.Create("traj.txt")
	defer dst.Close()

	for meter := 0; meter < 10000; meter += 1 {
		fmt.Fprintf(dst, "%g %g\n", s.x, s.y)
		s.next(0.3e-5, 0.01)
	}
}
