package likpdf

import (
	"github.com/massarakhsh/lik"
)

func (it *PDFFile) textDrawBox(px1, py1, px2, py2 float64, text string, opt string) {
	options := it.setOptions(opt)
	align := options.alignX + options.alignY
	lines := it.textPrepare(px2-px1, text)
	ml := len(lines)
	hline := options.fontSize / 2
	htext := hline * float64(ml)
	var by float64
	if options.alignY == "T" {
		by = py1
	} else if options.alignY == "B" {
		by = py2 - htext
	} else {
		by = (py1 + py2 - htext) / 2
	}
	if by < py1 {
		by = py1
	}
	for nl := 0; nl < ml; nl++ {
		py := by + hline*float64(nl)
		if nl > 0 && py+hline > py2 {
			break
		}
		it.pdf.SetXY(px1, py)
		it.pdf.CellFormat(px2-px1, hline, lines[nl], "", 0, align, false, 0, "")
	}
}

func (it *PDFFile) textPrepare(width float64, text string) []string {
	runes := []rune(text)
	mpos := len(runes)
	lines := []string{}
	line := ""
	word := ""
	for pos := 0; pos < mpos || word != ""; pos++ {
		brk := pos >= mpos
		sym := ""
		if pos < mpos {
			sym = string(runes[pos])
		}
		canbrk := false
		if sym == "\n" {
			sym = ""
			brk = true
			canbrk = true
		} else if sym == "\r" {
			sym = ""
			canbrk = true
		} else if !lik.RegExCompare(sym, "[0-9a-zA-Zа-яА-Я_]") {
			canbrk = true
		}
		if !canbrk {
			word += sym
			sym = ""
		} else {
			line += word
			word = ""
		}
		if !brk && line != "" {
			sall := line + word
			if it.pdf.GetStringWidth(sall) >= width {
				brk = true
			}
		}
		if brk {
			lines = append(lines, line)
			line = ""
		}
		word += sym
	}
	return lines
}
