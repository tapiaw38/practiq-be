package imagecompose

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	stddraw "image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"

	xdraw "golang.org/x/image/draw"
)

const (
	gapHeight  = 24
	ruleHeight = 2

	minSheetWidth  = 900
	maxSheetWidth  = 1400
	maxSheetPixels = 1_800_000

	maxTeacherHeight = 900
)

var ruleColor = color.RGBA{R: 0x94, G: 0xa3, B: 0xb8, A: 0xff}

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

	if topImg.Bounds().Dx() <= 0 || bottomImg.Bounds().Dx() <= 0 ||
		topImg.Bounds().Dy() <= 0 || bottomImg.Bounds().Dy() <= 0 {
		return nil, errors.New("imagecompose: empty source image")
	}

	width := sheetWidth(bottomImg.Bounds().Dx())
	topW, topH := teacherSize(topImg, width)
	bottomH := scaledHeight(bottomImg, width)

	if scale := pixelBudgetScale(width, topH+gapHeight+bottomH); scale < 1 {
		width = maxInt(1, int(float64(width)*scale))
		topW, topH = teacherSize(topImg, width)
		bottomH = scaledHeight(bottomImg, width)
	}

	height := topH + gapHeight + bottomH
	out := image.NewRGBA(image.Rect(0, 0, width, height))
	stddraw.Draw(out, out.Bounds(), image.NewUniform(color.White), image.Point{}, stddraw.Src)

	topX := (width - topW) / 2
	xdraw.CatmullRom.Scale(out, image.Rect(topX, 0, topX+topW, topH), topImg, topImg.Bounds(), xdraw.Over, nil)

	ruleY := topH + (gapHeight-ruleHeight)/2
	stddraw.Draw(
		out,
		image.Rect(0, ruleY, width, ruleY+ruleHeight),
		image.NewUniform(ruleColor),
		image.Point{},
		stddraw.Src,
	)

	bottomTop := topH + gapHeight
	xdraw.CatmullRom.Scale(out, image.Rect(0, bottomTop, width, bottomTop+bottomH), bottomImg, bottomImg.Bounds(), xdraw.Over, nil)

	var buf bytes.Buffer
	if err := png.Encode(&buf, out); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func sheetWidth(studentWidth int) int {
	switch {
	case studentWidth < minSheetWidth:
		return minSheetWidth
	case studentWidth > maxSheetWidth:
		return maxSheetWidth
	default:
		return studentWidth
	}
}

func teacherSize(img image.Image, width int) (int, int) {
	b := img.Bounds()
	w := width
	h := scaledHeight(img, w)
	if h > maxTeacherHeight {
		h = maxTeacherHeight
		w = maxInt(1, int(float64(b.Dx())*float64(h)/float64(b.Dy())))
	}
	return w, h
}

func scaledHeight(img image.Image, width int) int {
	b := img.Bounds()
	return maxInt(1, int(float64(b.Dy())*float64(width)/float64(b.Dx())))
}

func pixelBudgetScale(width, height int) float64 {
	total := width * height
	if total <= maxSheetPixels {
		return 1
	}
	return math.Sqrt(float64(maxSheetPixels) / float64(total))
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
