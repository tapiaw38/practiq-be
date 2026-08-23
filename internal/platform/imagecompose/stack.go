// Package imagecompose joins two pages into the single image an assistant can
// grade: the teacher's page above, the student's work below.
//
// Gillie's message carries one image_content field, so a page and the answer to
// it cannot travel as two attachments. The frontend already solved the same
// problem for the chat assistant by stacking them client-side; this is the
// server-side counterpart for the notebook, which is graded without a browser
// in the loop.
package imagecompose

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
)

const (
	// Space between the two pages, in the surface colour, so the boundary is
	// unmistakable without drawing a border over either one.
	gapHeight = 24
	// A thin rule inside the gap: the two pages are often both mostly white,
	// and white-on-white space alone can read as one continuous sheet.
	ruleHeight = 2
)

var ruleColor = color.RGBA{R: 0x94, G: 0xa3, B: 0xb8, A: 0xff}

// StackVertically returns a PNG with top drawn above bottom on a white sheet.
//
// Neither image is rescaled. A resampler would soften handwriting, which is the
// one thing the grader has to read, so the canvas takes the wider of the two and
// each page is centred at its native size.
func StackVertically(top, bottom []byte) ([]byte, error) {
	if len(top) == 0 && len(bottom) == 0 {
		return nil, errors.New("imagecompose: nothing to stack")
	}
	if len(top) == 0 {
		return bottom, nil
	}
	if len(bottom) == 0 {
		return top, nil
	}

	topImg, _, err := image.Decode(bytes.NewReader(top))
	if err != nil {
		return nil, err
	}
	bottomImg, _, err := image.Decode(bytes.NewReader(bottom))
	if err != nil {
		return nil, err
	}

	topB := topImg.Bounds()
	bottomB := bottomImg.Bounds()

	width := max(topB.Dx(), bottomB.Dx())
	height := topB.Dy() + gapHeight + bottomB.Dy()
	if width <= 0 || height <= 0 {
		return nil, errors.New("imagecompose: empty source image")
	}

	out := image.NewRGBA(image.Rect(0, 0, width, height))
	// Both sources may carry transparency, and the page the student saw was
	// white; without this they would composite onto black.
	draw.Draw(out, out.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	drawAt(out, topImg, (width-topB.Dx())/2, 0)

	ruleY := topB.Dy() + (gapHeight-ruleHeight)/2
	draw.Draw(
		out,
		image.Rect(0, ruleY, width, ruleY+ruleHeight),
		image.NewUniform(ruleColor),
		image.Point{},
		draw.Src,
	)

	drawAt(out, bottomImg, (width-bottomB.Dx())/2, topB.Dy()+gapHeight)

	var buf bytes.Buffer
	if err := png.Encode(&buf, out); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func drawAt(dst *image.RGBA, src image.Image, x, y int) {
	b := src.Bounds()
	draw.Draw(
		dst,
		image.Rect(x, y, x+b.Dx(), y+b.Dy()),
		src,
		b.Min,
		draw.Over,
	)
}
