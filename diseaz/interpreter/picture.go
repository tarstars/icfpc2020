package interpreter

import (
	"image"
	"image/color"
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

type Picture map[image.Point]bool

func (pic Picture) Eval(c Context) (Token, bool) {
	return pic, false
}

func (pic Picture) String() string {
	var pp []string
	for p := range pic {
		pp = append(pp, p.String())
	}
	return "[" + strings.Join(pp, " ") + "]"
}

func (pic Picture) Draw(x, y int) image.Point {
	r := image.Pt(x, y)
	pic[r] = true
	return r
}

func (pic Picture) DrawPicture(other Picture) {
	for p, v := range other {
		if !v {
			continue
		}
		pic[p] = true
	}
}

func (pic Picture) ColorModel() color.Model {
	return color.RGBAModel
}

func (pic Picture) Bounds() image.Rectangle {
	var b *image.Rectangle
	for p := range pic {
		if b == nil {
			b = &image.Rectangle{
				Min: image.Point{
					X: p.X,
					Y: p.Y,
				},
				Max: image.Point{
					X: p.X + 1,
					Y: p.Y + 1,
				},
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
	if pic[image.Pt(x, y)] {
		return White
	}
	return Black
}
