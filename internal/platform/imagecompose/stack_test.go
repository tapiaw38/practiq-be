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

func TestStackKeepsBothPagesAtNativeSize(t *testing.T) {
	top := pngOf(100, 40, color.RGBA{R: 0xff, A: 0xff})
	bottom := pngOf(60, 30, color.RGBA{B: 0xff, A: 0xff})

	out, err := StackVertically(top, bottom)
	if err != nil {
		t.Fatalf("stack: %v", err)
	}
	got := decode(t, out).Bounds()

	if got.Dx() != 100 {
		t.Fatalf("width should be the wider page, got %d", got.Dx())
	}
	if want := 40 + gapHeight + 30; got.Dy() != want {
		t.Fatalf("height = %d, want %d", got.Dy(), want)
	}
}

// The student's canvas arrives transparent where nothing was drawn; without a
// white sheet underneath it would composite onto black and the near-black ink
// would vanish.
func TestStackPaintsTransparencyWhite(t *testing.T) {
	transparent := pngOf(20, 20, color.RGBA{})
	opaque := pngOf(20, 20, color.RGBA{R: 0xff, A: 0xff})

	out, err := StackVertically(opaque, transparent)
	if err != nil {
		t.Fatalf("stack: %v", err)
	}
	img := decode(t, out)

	// Sample inside the bottom page, which was fully transparent.
	r, g, b, a := img.At(10, 20+gapHeight+10).RGBA()
	if r != 0xffff || g != 0xffff || b != 0xffff || a != 0xffff {
		t.Fatalf("transparent area should be white, got rgba(%d,%d,%d,%d)", r, g, b, a)
	}
}

func TestStackSeparatesThePagesWithARule(t *testing.T) {
	white := pngOf(40, 20, color.White)
	out, err := StackVertically(white, white)
	if err != nil {
		t.Fatalf("stack: %v", err)
	}
	img := decode(t, out)

	ruleY := 20 + (gapHeight-ruleHeight)/2
	r, g, b, _ := img.At(20, ruleY).RGBA()
	if r == 0xffff && g == 0xffff && b == 0xffff {
		t.Fatal("two white pages must not merge into one sheet: rule missing")
	}
}

// The teacher uploads whatever their phone produced, so the two sides are not
// always the same format.
func TestStackAcceptsMixedFormats(t *testing.T) {
	if _, err := StackVertically(jpegOf(50, 20, color.White), pngOf(50, 20, color.White)); err != nil {
		t.Fatalf("jpeg over png: %v", err)
	}
}

func TestStackDegradesInsteadOfFailing(t *testing.T) {
	page := pngOf(10, 10, color.White)

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
