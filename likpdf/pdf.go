package likpdf

import (
	"fmt"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/massarakhsh/lik"
)

type PDFFile struct {
	pdf                          *gofpdf.Fpdf
	width, height                float64
	mleft, mright, mtop, mbottom float64
	aspectRatio                  float64
	options                      PDFOptions
}

type PDFFiler interface {
	GetX() float64
	GetY() float64
	GetXY() (float64, float64)
	SetX(x float64)
	SetY(y float64)
	SetXY(x, y float64)
	SetOptions(opt string)
	AddPage()
	ToPad(x1, y1, x2, y2 float64, pad float64) (float64, float64, float64, float64)
	ToRatio(x1, y1, x2, y2 float64, ratio float64) (float64, float64, float64, float64)
	DrawRect(x1, y1, x2, y2 float64, opt string)
	DrawImage(x1, y1, x2, y2 float64, namefile string, opt string)
	DrawText(x1, y1, x2, y2 float64, text string, opt string)
	SaveToFile(name string) bool
}

type PDFOptions struct {
	alignX     string
	alignY     string
	fR, fG, fB int
	bR, bG, bB int
	fontSize   float64
}

func Create(opt string) PDFFiler {
	it := &PDFFile{}
	orient := "P"
	size := "A4"
	for _, o := range strings.Split(opt, ",") {
		if lik.RegExCompare(o, "^(P|L)$") {
			orient = o
		} else if lik.RegExCompare(o, "^(A|B)\\d+") {
			size = o
		}
	}
	it.pdf = gofpdf.New(orient, "mm", size, "")
	it.initialize()
	return it
}

func (it *PDFFile) initialize() {
	it.width, it.height = it.pdf.GetPageSize()
	it.mleft, it.mright, it.mtop, it.mbottom = it.pdf.GetMargins()
	it.aspectRatio = it.height / it.width
	it.options.bR = 255
	it.options.bG = 255
	it.options.bB = 255
	it.options.fontSize = 10
	dir, _ := os.Getwd()
	fontpath := dir + "/lib/font"
	it.pdf.SetFontLocation(fontpath)
	it.pdf.AddUTF8Font("dejavu", "", "DejaVuSansCondensed.ttf")
	it.pdf.AddUTF8Font("dejavu", "B", "DejaVuSansCondensed-Bold.ttf")
	it.pdf.AddUTF8Font("dejavu", "I", "DejaVuSansCondensed-Oblique.ttf")
	it.pdf.AddUTF8Font("dejavu", "BI", "DejaVuSansCondensed-BoldOblique.ttf")
}

func (it *PDFFile) xFrom(x float64) float64 {
	return (x - it.mleft) / (it.width - it.mright - it.mleft)
}
func (it *PDFFile) yFrom(y float64) float64 {
	return (y - it.mtop) / (it.height - it.mbottom - it.mtop)
}

func (it *PDFFile) xTo(x float64) float64 {
	return it.mleft + x*(it.width-it.mright-it.mleft)
}
func (it *PDFFile) yTo(y float64) float64 {
	return it.mtop + y*(it.height-it.mbottom-it.mtop)
}

func (it *PDFFile) from4(x1, y1, x2, y2 float64) (float64, float64, float64, float64) {
	return it.xFrom(x1), it.yFrom(y1), it.xFrom(x2), it.yFrom(y2)
}

func (it *PDFFile) to4(x1, y1, x2, y2 float64) (float64, float64, float64, float64) {
	return it.xTo(x1), it.yTo(y1), it.xTo(x2), it.yTo(y2)
}

func (it *PDFFile) ToPad(x1, y1, x2, y2 float64, pad float64) (float64, float64, float64, float64) {
	px1 := x1 + (x2-x1)*pad/2
	py1 := y1 + (y2-y1)*pad/2
	px2 := x2 - (x2-x1)*pad/2
	py2 := y2 - (y2-y1)*pad/2
	return px1, py1, px2, py2
}

func (it *PDFFile) ToRatio(x1, y1, x2, y2 float64, ratio float64) (float64, float64, float64, float64) {
	px1, py1, px2, py2 := it.to4(x1, y1, x2, y2)
	pw := px2 - px1
	ph := py2 - py1
	area := ph / pw
	if ratio <= area {
		rh := pw * ratio
		py1 += ph/2 - rh/2
		py2 -= ph/2 - rh/2
	} else {
		rw := ph / ratio
		px1 += pw/2 - rw/2
		px2 -= pw/2 - rw/2
	}
	return it.from4(px1, py1, px2, py2)
}

func (it *PDFFile) buildOptions(opt string) PDFOptions {
	options := it.options
	opts := strings.Split(strings.ToUpper(opt), ",")
	for _, o := range opts {
		if o == "L" || o == "C" || o == "R" {
			options.alignX = o
		} else if o == "T" || o == "M" || o == "B" {
			options.alignY = o
		} else if match := lik.RegExParse(o, "^F([0-9\\.]*)"); match != nil {
			if sz := lik.StrToFloat(match[1]); sz > 0 {
				options.fontSize = sz
			}
		} else if match := lik.RegExParse(o, "^#([0-9a-fA-F]{3,6})"); match != nil {
			code := []byte(match[1])
			if len(code) == 3 {
				rgb := [3]int{0, 0, 0}
				for c := 0; c < 3; c++ {
					cd := int(code[c])
					if cd >= 0x30 && cd <= 0x39 {
						cd = cd - 0x30
					} else if cd >= 0x41 && cd <= 0x46 {
						cd = cd - 0x41 + 10
					} else if cd >= 0x61 && cd <= 0x66 {
						cd = cd - 0x61 + 10
					}
					rgb[c] = cd * 0x11
				}
				options.fR = rgb[0]
				options.fG = rgb[1]
				options.fB = rgb[2]
			}
		}
	}
	return options
}

func (it *PDFFile) setOptions(opt string) PDFOptions {
	options := it.buildOptions(opt)
	it.pdf.SetDrawColor(options.fR, options.fG, options.fB)
	it.pdf.SetFont("dejavu", "", options.fontSize)
	return options
}

func (it *PDFFile) SetOptions(opt string) {
	it.options = it.buildOptions(opt)
}

func (it *PDFFile) SaveToFile(name string) bool {
	if match := lik.RegExParse(name, "(.+)/[^/]*$"); match != nil {
		os.MkdirAll(match[1], os.ModePerm)
	}
	if err := it.pdf.OutputFileAndClose(name); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (it *PDFFile) AddPage() {
	it.pdf.AddPage()
}

func (it *PDFFile) SetX(x float64) {
	it.pdf.SetX(it.xTo(x))
}
func (it *PDFFile) SetY(y float64) {
	it.pdf.SetY(it.yTo(y))
}
func (it *PDFFile) SetXY(x, y float64) {
	it.pdf.SetXY(it.xTo(x), it.yTo(y))
}

func (it *PDFFile) GetX() float64 {
	return it.xFrom(it.pdf.GetX())
}
func (it *PDFFile) GetY() float64 {
	return it.yFrom(it.pdf.GetY())
}
func (it *PDFFile) GetXY() (float64, float64) {
	return it.GetX(), it.GetY()
}

func (it *PDFFile) DrawRect(x1, y1, x2, y2 float64, opt string) {
	px1, py1, px2, py2 := it.to4(x1, y1, x2, y2)
	it.setOptions(opt)
	it.pdf.Rect(px1, py1, px2-px1, py2-py1, "")
}

func (it *PDFFile) DrawImage(x1, y1, x2, y2 float64, namefile string, opt string) {
	iw, ih := 0.0, 0.0
	pdf := gofpdf.New("P", "mm", "A4", "")
	if img := pdf.RegisterImage(namefile, ""); img != nil {
		iw, ih = img.Extent()
	}
	if iw > 0 && ih >= 0 {
		it.setOptions(opt)
		px1, py1, px2, py2 := it.to4(x1, y1, x2, y2)
		pw := px2 - px1
		ph := py2 - py1
		if iw/pw >= ih/ph {
			rh := pw * ih / iw
			py1 += ph/2 - rh/2
			py2 -= ph/2 - rh/2
		} else {
			rw := ph * iw / ih
			px1 += pw/2 - rw/2
			px2 -= pw/2 - rw/2
		}
		it.pdf.Image(namefile, px1, py1, px2-px1, py2-py1, false, "", 0, "")
	}
}

func (it *PDFFile) DrawText(x1, y1, x2, y2 float64, text string, opt string) {
	px1, py1, px2, py2 := it.to4(x1, y1, x2, y2)
	it.textDrawBox(px1, py1, px2, py2, text, opt)
}
