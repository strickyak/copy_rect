// Hint1:  go run copy_rect.go  echo HI image? foo  load ~/Downloads/8321629-439281-rye-bread-face.jpg nando  image? nando model? nando bounds? nando new 641 481 z image? z bounds? z  model? z cp nando 100 200 50 70 z 300 100 jpeg  z /tmp/a.jpg

// Hint2: for x in /tmp/d/z00000001/00000*.jpg ; do ~/gocode/src/github.com/strickyak/copy_rect/copy_rect load $x  a , bounds? a , new 640 480 b , cp a 300 300 640 480 b 0 0 , jpeg b /tmp/a.jpg  ; done

// Hint3: mplayer -ao none -vo jpeg:outdir=/tmp/d:subdirs=z ~/Downloads/video.flv

package main

import (
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
)
import . "fmt"

var Store = make(map[string]image.Image)

func CheckPtr(a interface{}) {
	if a == nil {
		panic("was nil")
	}
}
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
func CheckOk(ok bool) {
	if !ok {
		panic("not ok")
	}
}

func Atoi(s string) int {
	z, err := strconv.Atoi(s)
	Check(err)
	return z
}
func Rect(min, max image.Point) image.Rectangle {
	return image.Rectangle{Min: min, Max: max}
}
func Pt(x, y int) image.Point {
	return image.Point{X: x, Y: y}
}

func Do(a []string) []string {
	switch a[0] {
	case "jpeg":
		im := Store[a[1]]
		w, err := os.Create(a[2])
		Check(err)
		err = jpeg.Encode(w, im, nil)
		Check(err)
		err = w.Close()
		Check(err)
		return a[3:]
	case "png":
		im := Store[a[1]]
		w, err := os.Create(a[2])
		Check(err)
		err = png.Encode(w, im)
		Check(err)
		err = w.Close()
		Check(err)
		return a[3:]
	case "cp":
		im := Store[a[1]]
		CheckPtr(im)
		x, y := Atoi(a[2]), Atoi(a[3])
		w, h := Atoi(a[4]), Atoi(a[5])
		im2, ok := Store[a[6]].(draw.Image)
		CheckOk(ok)
		CheckPtr(im2)
		x2, y2 := Atoi(a[7]), Atoi(a[8])
		draw.FloydSteinberg.Draw(im2, Rect(Pt(x2, y2), Pt(x2+w, y2+h)), im, Pt(x, y))
		return a[9:]
	case "new":
		r := Rect(Pt(0, 0), Pt(Atoi(a[1]), Atoi(a[2])))
		im := image.NewRGBA(r)
		Store[a[3]] = im
		return a[4:]
	case "load":
		f, err := os.Open(a[1])
		Check(err)
		im, s, err := image.Decode(f)
		Check(err)
		Printf("Load %q format %q -> %q\n", a[1], s, a[2])
		err = f.Close()
		Check(err)
		Store[a[2]] = im
		return a[3:]
	case "bounds?":
		im := Store[a[1]]
		b := im.Bounds()
		Printf("Bounds %q = min (%d, %d) max (%d, %d)\n", a[1], b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)
		return a[2:]
	case "model?":
		im := Store[a[1]]
		Printf("ColorModel %q = (%T) %v\n", a[1], im.ColorModel(), im.ColorModel())
		return a[2:]
	case "image?":
		im := Store[a[1]]
		Printf("Image %q = (%T)\n", a[1], im)
		return a[2:]
	case ",":
		// NOP
		return a[1:]
	case "echo":
		Printf("ECHO: %q\n", a[1])
		return a[2:]
	}
	panic("UNKNOWN COMMAND: " + a[0])
}

func main() {
	a := os.Args[1:]
	for len(a) > 0 {
		Printf("DOING %v\n", a)
		a = Do(a)
	}
}
