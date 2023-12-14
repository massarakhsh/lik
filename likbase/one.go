package likbase

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/massarakhsh/lik"
	"reflect"
	"time"
)

//	Дескриптор базы данных
var ODB *gorm.DB

//	Ядро объекта базы данных
type One struct {
	Id        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

//	Интерфейс объекта базы данных
type Oner interface {
	Table() string
	GetId() lik.IDB
	Create(datas... interface{}) bool
	Update(datas... interface{}) bool
	Delete()
}

//	Инициализация базы данных
func OpenBase(serv string, base string, user string, pass string) bool {
	args := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, serv, base)
	if ODB,_ = gorm.Open("mysql", args); ODB == nil {
		return false
	}
	return true
}

func dbtable(one Oner) *gorm.DB {
	return ODB.Table(one.Table())
}

func (it *One) GetId() lik.IDB {
	return lik.IDB(it.Id)
}

//	Прочитать объект
func Read(id lik.IDB, one Oner) bool {
	return dbtable(one).First(one, int(id)) != nil
}

//	Обновить объект
func Update(one Oner, datas... interface{}) bool {
	valit := reflect.ValueOf(one).Elem()
	ndm := len(datas)
	for nd := 0; nd < ndm; nd++ {
		data := datas[nd]
		switch key := data.(type) {
		case string:
			var val interface{}
			if match := lik.RegExParse(key, "(.+?)=(.*)"); match != nil {
				key = match[1]
				val = match[2]
			} else if nd + 1 < ndm {
				nd++
				val = datas[nd]
			} else {
				break
			}
			if field := valit.FieldByName(key); field.IsValid() {
				if typ := field.Type().Name(); typ == "string" {
					field.SetString(toString(val))
				} else if typ == "int" {
					field.SetInt(int64(toInt(val)))
				} else if typ == "float" {
					field.SetFloat(toFloat(val))
				}
			} else {
				fmt.Println("Update bad field: ", one.Table(), ": ", key)
				return false
			}
		default:
			fmt.Println("Update ERROR: ", one.Table(), ": ", data)
			return false
		}
	}
	if one.GetId() > 0 {
		return dbtable(one).Save(one) != nil
	} else {
		return dbtable(one).Create(one) != nil
	}
}

//	Удалить объект
func Delete(one Oner) {
	dbtable(one).Delete(one)
}

//	Интефейс в целое
func toInt(data interface{}) int {
	val := 0
	if data != nil {
		switch da := data.(type) {
		case int:
			val = da
		case uint:
			val = int(da)
		case byte:
			val = int(da)
		case string:
			val = lik.StrToInt(da)
		}
	}
	return val
}

//	Интерфейс в плавающее
func toFloat(data interface{}) float64 {
	val := 0.0
	if data != nil {
		switch da := data.(type) {
		case float64:
			val = da
		case int:
			val = float64(da)
		case uint:
			val = float64(da)
		case byte:
			val = float64(da)
		case string:
			val = lik.StrToFloat(da)
		}
	}
	return val
}

//	Интерфейс в строку
func toString(data interface{}) string {
	val := ""
	if data != nil {
		switch da := data.(type) {
		case string:
			val = da
		case float64:
			val = lik.FloatToStr(da)
		case int:
			val = lik.IntToStr(da)
		case uint:
			val = lik.IntToStr(int(da))
		case byte:
			val = lik.IntToStr(int(da))
		}
	}
	return val
}

