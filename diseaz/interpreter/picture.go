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
	Palette = color.Palette{
		rgba(0xff, 0xff, 0xff, 0x80),
		rgba(0, 0, 0xff, 0x80),
		rgba(0xff, 0, 0xff, 0x80),
		rgba(0, 0xff, 0xff, 0x80),
		rgba(0, 0xff, 0, 0x80),
		rgba(0xff, 0xff, 0, 0x80),
		rgba(0xff, 0, 0, 0x80),
		// rgba(0, 0, 0, 0x80),
	}
)

func rgba(r, g, b, a uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: a}
}

type Picture struct {
	Pts  map[Point]int `json:"-"`
	Draw [][]Point     `json:""`
}

func NewPicture(pts ...Point) *Picture {
	pic := &Picture{
		Pts: make(map[Point]int),
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
	for p := range pic.Pts {
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
	color := len(pic.Draw)
	// log.Printf("DrawPts(%p) %#v", pic, pts)
	// log.Printf("  Picture[before] %s", pic)
	// log.Printf("  Draw[before] %#v", pic.Draw)
	pic.Draw = append(pic.Draw, pts)
	for _, p := range pts {
		if _, found := pic.Pts[p]; found {
			continue
		}
		pic.Pts[p] = color
	}
	// log.Printf("  Picture[after] %s", pic)
	// log.Printf("  Draw[after] %#v", pic.Draw)

	return pic
}

func (pic *Picture) DrawPicture(other *Picture) {
	// log.Printf("DrawPicture(%p) %s", pic, other)
	// log.Printf("  Picture[before] %s", pic)
	// log.Printf("  Draw[before] %#v", pic.Draw)
	color := len(pic.Draw)
	pic.Draw = append(pic.Draw, other.Draw...)
	for p, v := range other.Pts {
		if _, found := pic.Pts[p]; found {
			continue
		}
		pic.Pts[p] = v + color
	}
	// log.Printf("  Picture[after] %s", pic)
	// log.Printf("  Draw[after] %#v", pic.Draw)
}

func (pic Picture) ColorModel() color.Model {
	return color.RGBAModel
}

const box = 5

func (pic Picture) Bounds() image.Rectangle {
	var b *image.Rectangle
	bx := image.Pt(box, box)
	for p := range pic.Pts {
		q := p.ToImg().Mul(box)
		u := q.Add(bx)
		if b == nil {
			b = &image.Rectangle{
				Min: q,
				Max: u,
			}
			continue
		}
		if q.X < b.Min.X {
			b.Min.X = q.X
		}
		if u.X > b.Max.X {
			b.Max.X = u.X
		}
		if q.Y < b.Min.Y {
			b.Min.Y = q.Y
		}
		if u.Y > b.Max.Y {
			b.Max.Y = u.Y
		}
	}
	if b == nil {
		b = &image.Rectangle{}
	}
	return *b
}

func fromBox(x int) int {
	if x < 0 {
		x = x - box + 1
	}
	return x / box
}

func (pic Picture) At(x, y int) color.Color {
	colorIdx, found := pic.Pts[Pt(fromBox(x), fromBox(y))]
	if !found {
		return Black
	}
	return Palette[colorIdx]
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
