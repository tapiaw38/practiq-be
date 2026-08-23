package imagecompose

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"
)

func pngOf(w, h int, c color.Color) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func jpegOf(w, h int, c color.Color) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func decode(t *testing.T, raw []byte) image.Image {
	t.Helper()
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	return img
}

func isWhite(img image.Image, x, y int) bool {
	r, g, b, a := img.At(x, y).RGBA()
	return r == 0xffff && g == 0xffff && b == 0xffff && a == 0xffff
}

func TestStackStaysWithinThePixelBudget(t *testing.T) {
	teacher := pngOf(2480, 3508, color.RGBA{R: 0xff, A: 0xff})
	student := pngOf(1000, 500, color.RGBA{B: 0xff, A: 0xff})

	b := decode(t, mustStack(t, teacher, student)).Bounds()

	if got := b.Dx() * b.Dy(); got > maxSheetPixels {
		t.Fatalf("sheet exceeds the pixel budget: %d", got)
	}
	if b.Dx() > maxSheetWidth || b.Dx() < 1 {
		t.Fatalf("unexpected sheet width %d", b.Dx())
	}
}

func TestStackKeepsTheStudentHalfLegible(t *testing.T) {
	teacher := pngOf(2480, 3508, color.White)
	student := pngOf(1000, 500, color.White)

	b := decode(t, mustStack(t, teacher, student)).Bounds()
	studentH := scaledHeight(decode(t, student), b.Dx())

	if share := float64(studentH) / float64(b.Dy()); share < 0.10 {
		t.Fatalf("student half is %.1f%% of the sheet, too small to survive downscaling", share*100)
	}
}

func TestStackClampsTinyAndHugeStudentPages(t *testing.T) {
	if w := sheetWidth(200); w != minSheetWidth {
		t.Fatalf("a tiny canvas should be lifted to %d, got %d", minSheetWidth, w)
	}
	if w := sheetWidth(4000); w != maxSheetWidth {
		t.Fatalf("a huge canvas should be capped at %d, got %d", maxSheetWidth, w)
	}
	if w := sheetWidth(1100); w != 1100 {
		t.Fatalf("a canvas already in range should be kept, got %d", w)
	}
}

func TestStackPaintsTransparencyWhite(t *testing.T) {
	teacher := pngOf(1000, 400, color.RGBA{R: 0xff, A: 0xff})
	student := pngOf(1000, 400, color.RGBA{})

	img := decode(t, mustStack(t, teacher, student))
	b := img.Bounds()

	if !isWhite(img, b.Dx()/2, b.Dy()-3) {
		r, g, bl, a := img.At(b.Dx()/2, b.Dy()-3).RGBA()
		t.Fatalf("student area should be white, got rgba(%d,%d,%d,%d)", r, g, bl, a)
	}
}

func TestStackSeparatesThePagesWithARule(t *testing.T) {
	white := pngOf(1000, 400, color.White)

	img := decode(t, mustStack(t, white, white))
	b := img.Bounds()

	found := false
	for y := b.Min.Y; y < b.Max.Y && !found; y++ {
		if !isWhite(img, b.Dx()/2, y) {
			found = true
		}
	}
	if !found {
		t.Fatal("two white pages must not merge into one sheet: rule missing")
	}
}

func TestStackAcceptsMixedFormats(t *testing.T) {
	if _, err := StackVertically(jpegOf(900, 400, color.White), pngOf(900, 400, color.White)); err != nil {
		t.Fatalf("jpeg over png: %v", err)
	}
}

func TestStackDegradesInsteadOfFailing(t *testing.T) {
	page := pngOf(100, 100, color.White)

	if out, err := StackVertically(nil, page); err != nil || !bytes.Equal(out, page) {
		t.Fatal("a missing teacher page should yield the student's page unchanged")
	}
	if out, err := StackVertically(page, nil); err != nil || !bytes.Equal(out, page) {
		t.Fatal("a missing student page should yield the teacher's page unchanged")
	}
	if _, err := StackVertically(nil, nil); err == nil {
		t.Fatal("two missing pages should error rather than return an empty image")
	}
	if _, err := StackVertically([]byte("not an image"), page); err == nil {
		t.Fatal("undecodable input should error so the caller can fall back to text")
	}
}

func mustStack(t *testing.T, top, bottom []byte) []byte {
	t.Helper()
	out, err := StackVertically(top, bottom)
	if err != nil {
		t.Fatalf("stack: %v", err)
	}
	return out
}

func TestStudentStrokeSurvivesModelDownscale(t *testing.T) {
	teacher := pngOf(2480, 3508, color.White)
	student := pngOf(1000, 500, color.White)

	out, err := StackVertically(teacher, student)
	if err != nil {
		t.Fatal(err)
	}
	img, _, _ := image.Decode(bytes.NewReader(out))
	b := img.Bounds()

	longest := b.Dy()
	if b.Dx() > longest {
		longest = b.Dx()
	}
	factor := 1024.0 / float64(longest)
	if factor > 1 {
		factor = 1
	}
	strokePx := 3.0 * factor

	t.Logf("sheet=%dx%d  factor=%.3f  trazo 3px -> %.2fpx", b.Dx(), b.Dy(), factor, strokePx)
	if strokePx < 1.5 {
		t.Fatalf("stroke collapses to %.2fpx after downscale", strokePx)
	}
}
