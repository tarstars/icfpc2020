package interpreter

import (
	"fmt"
	"image"
	"image/color"
	"sort"
	"strings"
)

var (
	White = color.RGBA{
		R: 0xff,
		G: 0xff,
		B: 0xff,
		A: 0xff,
	}
	Black = color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0xff,
	}
)

type Picture struct {
	Pts  map[Point]bool `json:"-"`
	Draw [][]Point      `json:""`
}

func NewPicture(pts ...Point) *Picture {
	pic := &Picture{
		Pts: make(map[Point]bool),
	}
	if len(pts) > 0 {
		pic.DrawPts(pts...)
	}
	return pic
}

func (pic Picture) Eval(c Context) (Token, bool) {
	return pic, false
}

func (pic Picture) String() string {
	var pp []string
	for _, p := range pic.Serial() {
		pp = append(pp, p.String())
	}
	return "{" + strings.Join(pp, " ") + "}"
}

func (pic Picture) Galaxy() string {
	var pp []string
	for _, p := range pic.Serial() {
		pp = append(pp, p.String())
	}
	return "{" + strings.Join(pp, " ") + "}"
}

func (pic Picture) Serial() []Point {
	r := []Point{}
	for p, v := range pic.Pts {
		if !v {
			continue
		}
		r = append(r, p)
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i].Lt(r[j])
	})
	return r
}

func (pic *Picture) DrawPts(pts ...Point) *Picture {
	if len(pts) == 0 {
		pic.Draw = append(pic.Draw, []Point{})
		return pic
	}
	// log.Printf("DrawPts(%p) %#v", pic, pts)
	// log.Printf("  Picture[before] %s", pic)
	// log.Printf("  Draw[before] %#v", pic.Draw)
	pic.Draw = append(pic.Draw, pts)
	for _, p := range pts {
		pic.Pts[p] = true
	}
	// log.Printf("  Picture[after] %s", pic)
	// log.Printf("  Draw[after] %#v", pic.Draw)

	return pic
}

func (pic *Picture) DrawPicture(other *Picture) {
	// log.Printf("DrawPicture(%p) %s", pic, other)
	// log.Printf("  Picture[before] %s", pic)
	// log.Printf("  Draw[before] %#v", pic.Draw)
	pic.Draw = append(pic.Draw, other.Draw...)
	for p, v := range other.Pts {
		if !v {
			continue
		}
		pic.Pts[p] = true
	}
	// log.Printf("  Picture[after] %s", pic)
	// log.Printf("  Draw[after] %#v", pic.Draw)
}

func (pic Picture) ColorModel() color.Model {
	return color.RGBAModel
}

func (pic Picture) Bounds() image.Rectangle {
	var b *image.Rectangle
	for p := range pic.Pts {
		if b == nil {
			b = &image.Rectangle{
				Min: p.ToImg(),
				Max: p.ToImg().Add(image.Pt(1, 1)),
			}
			continue
		}
		if p.X < b.Min.X {
			b.Min.X = p.X
		}
		if p.X+1 > b.Max.X {
			b.Max.X = p.X + 1
		}
		if p.Y < b.Min.Y {
			b.Min.Y = p.Y
		}
		if p.Y+1 > b.Max.Y {
			b.Max.Y = p.Y + 1
		}
	}
	if b == nil {
		b = &image.Rectangle{}
	}
	return *b
}

func (pic Picture) At(x, y int) color.Color {
	if pic.Pts[Pt(x, y)] {
		return White
	}
	return Black
}

type Point struct {
	X int `json:""`
	Y int `json:""`
}

func Pt(x, y int) Point {
	return Point{X: x, Y: y}
}

func (p Point) String() string {
	return fmt.Sprintf("[%d, %d]", p.X, p.Y)
}

func (p Point) ToImg() image.Point {
	return image.Pt(p.X, p.Y)
}

func (p Point) Lt(q Point) bool {
	if p.X != q.X {
		return p.X < q.X
	}
	return p.Y < q.Y
}
