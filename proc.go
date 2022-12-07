package lik

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	DayWeek = [8]string{"Воскресенье", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}
)

func GetArgs(args []string) (Seter, bool) {
	maper := BuildSet()
	isok := true
	for na := 0; na < len(args); na++ {
		key := args[na]
		val := ""
		if strings.HasPrefix(key, "--") {
			if pos := strings.Index(key, "="); pos >= 0 {
				val = key[pos+1:]
				key = key[:pos]
			}
		} else if strings.HasPrefix(key, "-") && na+1 < len(args) {
			na++
			val = args[na]
		}
		for strings.HasPrefix(key, "-") {
			key = key[1:]
		}
		if len(key) > 0 {
			key = strings.ToLower(key)
		} else {
			key = strings.ToLower(val)
			val = ""
		}
		if len(key) > 0 {
			if vali, ok := StrToIntIf(val); ok {
				maper.SetValue(key, vali)
			} else {
				maper.SetValue(key, val)
			}
		}
	}
	return maper, isok
}

func StrToIntIf(str string) (int, bool) {
	val, ok := 0, false
	if vl, err := strconv.ParseInt(str, 0, 32); err == nil {
		val, ok = int(vl), true
	}
	return val, ok
}
func StrToInt(str string) int {
	if val, ok := StrToIntIf(str); ok {
		return val
	}
	return 0
}
func StrToInt64If(str string) (int64, bool) {
	val, ok := int64(0), false
	if vl, err := strconv.ParseInt(str, 0, 64); err == nil {
		val, ok = vl, true
	}
	return val, ok
}
func StrToInt64(str string) int64 {
	if val, ok := StrToInt64If(str); ok {
		return val
	}
	return 0
}

func StrToIDB(str string) IDB {
	return IDB(StrToInt(str))
}

func IntToStr(val int) string {
	return fmt.Sprint(val)
}
func IDBToStr(val IDB) string {
	return fmt.Sprint(int(val))
}

func StrToFloatIf(str string) (float64, bool) {
	str = strings.Replace(str, ",", ".", -1)
	if val, err := strconv.ParseFloat(str, 64); err == nil {
		return val, true
	}
	return 0, false
}
func StrToFloat(str string) float64 {
	if val, ok := StrToFloatIf(str); ok {
		return val
	}
	return 0
}

func FloatToStr(val float64) string {
	return fmt.Sprint(val)
}

func FloatFromBytes(val []byte) float64 {
	bits := binary.LittleEndian.Uint64(val)
	res := math.Float64frombits(bits)
	return res
}

func FloatToBytes(val float64) []byte {
	bits := math.Float64bits(val)
	res := make([]byte, 8)
	binary.LittleEndian.PutUint64(res, bits)
	return res
}

func StringFromXS(data string) string {
	str := ""
	bts := []byte(data)
	for b := 0; b+3 < len(bts); b += 4 {
		bth := ByteFromX2(bts[b], bts[b+1])
		btl := ByteFromX2(bts[b+2], bts[b+3])
		str += fmt.Sprintf("%c", int(btl)+int(bth)*256)
	}
	return str
}
func IntFromXS(data string) int {
	return StrToInt(StringFromXS(data))
}
func IDBFromXS(data string) IDB {
	return IDB(IntFromXS(data))
}

func BytesFromXS(data string) []byte {
	raw := []byte{}
	bts := []byte(data)
	for b := 0; b+1 < len(bts); b += 2 {
		bt := ByteFromX2(bts[b], bts[b+1])
		raw = append(raw, bt)
	}
	return raw
}

func ByteFromX2(ch byte, cl byte) byte {
	bt := byte(0)
	if cl >= '0' && cl <= '9' {
		bt += (cl - '0')
	} else if cl >= 'a' && cl <= 'f' {
		bt += (cl - 'a' + 10)
	} else if cl >= 'A' && cl <= 'F' {
		bt += (cl - 'A' + 10)
	}
	if ch >= '0' && ch <= '9' {
		bt += 16 * (ch - '0')
	} else if ch >= 'a' && ch <= 'f' {
		bt += 16 * (ch - 'a' + 10)
	} else if ch >= 'A' && ch <= 'F' {
		bt += 16 * (ch - 'A' + 10)
	}
	return bt
}

func StringToXS(data string) string {
	text := ""
	runes := []rune(data)
	for _, rn := range runes {
		text += fmt.Sprintf("%04x", rn)
	}
	return text
}

func BytesToXS(data []byte) string {
	text := ""
	for _, bt := range data {
		text += fmt.Sprintf("%02x", bt)
	}
	return text
}

func IntToMoney(val int) string {
	isneg := (val < 0)
	if isneg {
		val = -val
	}
	data := ""
	for ndg := 0; val > 0; ndg++ {
		if ndg > 0 && ndg%3 == 0 {
			data = " " + data
		}
		data = fmt.Sprintf("%d%s", val%10, data)
		val /= 10
	}
	if isneg {
		data = "-" + data
	}
	return data
}

func IsIt(a string, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}

func LimitString(text string, size int) string {
	if size >= 3 {
		runes := []rune(text)
		if real := len(runes); real > size {
			end := (size - 3) / 2
			begin := size - 3 - end
			text = string(runes[:begin]) + "..." + string(runes[real-end:])
		}
	}
	return text
}

func LengthString(text string) int {
	runes := []rune(text)
	return len(runes)
}

func SubString(text string, pos int, size int) string {
	result := ""
	runes := []rune(text)
	total := len(runes)
	if pos < 0 {
		size += pos
		pos = 0
	}
	if pos > total {
		pos = total
	}
	if pos+size > total {
		size = total - pos
	}
	if size > 0 {
		result = string(runes[pos : pos+size])
	}
	return result
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func RegExCompare(text string, reg string) bool {
	return regexp.MustCompile(reg).MatchString(text)
}

func RegExParse(text string, reg string) []string {
	return regexp.MustCompile(reg).FindStringSubmatch(text)
}

func RegExReplace(text string, reg string, repl string) string {
	return regexp.MustCompile(reg).ReplaceAllString(text, repl)
}

func GetFirstExt(path string) (string, string) {
	if match := RegExParse(path, "^/(.*)"); match != nil {
		path = match[1]
	}
	name, ext := path, ""
	if match := RegExParse(path, "^([^/]*)/(.*)"); match != nil {
		name = match[1]
		ext = match[2]
	}
	return name, ext
}

func PathToNames(path string) []string {
	names := strings.Split(path, "/")
	for len(names) > 0 && names[0] == "" {
		names = names[1:]
	}
	return names
}

func CompareVersion(ver1 string, ver2 string) int {
	cp := 0
	for cp == 0 && (ver1 != "" || ver2 != "") {
		v1 := ""
		i1 := 0
		ok1 := false
		if match := RegExParse(ver1, "^([^\\.]*)\\.(.*)$"); match != nil {
			v1 = match[1]
			ver1 = match[2]
		} else {
			v1 = ver1
			ver1 = ""
		}
		if match := RegExParse(v1, "^(\\d+)"); match != nil {
			i1 = StrToInt(match[1])
			ok1 = true
		}
		v2 := ""
		i2 := 0
		ok2 := false
		if match := RegExParse(ver2, "^([^\\.]*)\\.(.*)$"); match != nil {
			v2 = match[1]
			ver2 = match[2]
		} else {
			v2 = ver2
			ver2 = ""
		}
		if match := RegExParse(v2, "^(\\d+)"); match != nil {
			i2 = StrToInt(match[1])
			ok2 = true
		}
		if ok1 && ok2 {
			if i1 > i2 {
				cp = 1
			} else if i1 < i2 {
				cp = -1
			}
		} else if ok1 {
			cp = 1
		} else if ok2 {
			cp = -1
		}
		if cp == 0 {
			if v1 > v2 {
				cp = 1
			} else if v1 < v2 {
				cp = -1
			}
		}
	}
	return cp
}

// RuTransiltMap описывает замены русских букв на английские при транслитерации. Некоторые буквы
// заменяются ни на одну, а на две или три буквы латинского алфавита. А мягкий знак вообще исчезает.
// Но такова обычная распространенная схема транслитерации.
var RuTransiltMap = map[rune]string{
	'а': "a",
	'б': "b",
	'в': "v",
	'г': "g",
	'д': "d",
	'е': "e",
	'ё': "yo",
	'ж': "zh",
	'з': "z",
	'и': "i",
	'й': "j",
	'к': "k",
	'л': "l",
	'м': "m",
	'н': "n",
	'о': "o",
	'п': "p",
	'р': "r",
	'с': "s",
	'т': "t",
	'у': "u",
	'ф': "f",
	'х': "h",
	'ц': "c",
	'ч': "ch",
	'ш': "sh",
	'щ': "sch",
	'ъ': "'",
	'ы': "y",
	'ь': "",
	'э': "e",
	'ю': "ju",
	'я': "ja",
}

func Transliterate(text string) string {
	return TransliterateMap(text, RuTransiltMap)
}

func TransliterateMap(text string, translitMap map[rune]string) string {
	result := ""
	runes := []rune(text)
	for index := 0; index < len(runes); index++ {
		run := runes[index]
		if run >= '0' && run <= '9' ||
			run >= 'a' && run <= 'z' ||
			run >= 'A' && run <= 'Z' ||
			run == '-' || run == '_' || run == '.' {
			result += string(run)
		} else if str, ok := translitMap[unicode.ToLower(run)]; ok {
			if !unicode.IsUpper(run) {
				result += str
			} else if (len(runes) > index+1 && unicode.IsUpper(runes[index+1])) ||
				(index > 0 && unicode.IsUpper(runes[index-1])) {
				result += strings.ToUpper(str)
			} else {
				result += strings.Title(str)
			}
		} else {
			result += "_"
		}
	}
	return result
}

func IfString(cond bool, yes string, no string) string {
	if cond {
		return yes
	} else {
		return no
	}
}
