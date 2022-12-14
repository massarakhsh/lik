# Библиотека lik

Управление динамическими объектами

- [Библиотека lik](#библиотека-lik)
  - [Examples](#examples)
  - [Interfaces](#interfaces)
  - [Interface LikItem](#interface-likitem)
    - [BuildItem(data interface{}) Itemer](#builditemdata-interface-itemer)
    - [item.Clone() Itemer](#itemclone-itemer)
    - [item.Format(prefix string) string](#itemformatprefix-string-string)
    - [IsBool() bool](#isbool-bool)
    - [IsInt() bool](#isint-bool)
    - [IsFloat() bool](#isfloat-bool)
    - [IsString() bool](#isstring-bool)
    - [IsList() bool](#islist-bool)
    - [IsSet() bool](#isset-bool)
    - [Serialize() string](#serialize-string)
    - [ToBool() bool](#tobool-bool)
    - [ToInt() int64](#toint-int64)
    - [ToFloat() float64](#tofloat-float64)
    - [ToString() string](#tostring-string)
    - [ToList() Lister](#tolist-lister)
    - [ToSet() Seter](#toset-seter)
  - [Interface LikSet](#interface-likset)
  - [Interface LikList](#interface-liklist)

Библиотека предоставляет несколько интерфейсов, позволяющих организовать сложные многоуровневые динамические объекты,
объединяющие простые скалярные значения, структуры из именованных полей - объектов и массивы объектов.

Отличительной особенностью библиотеки является одновременная поддержка следующих возможностей:

- Работа с объектом как с единым целым, включая клонирование, сериализацию в строку json и восстановление из строки
- Интерпретацию объекта как файловой системы в памяти с доступом по пути, например, `/data/info/list/8/id`
- Получение интерфейсов подобъектов и работа непосредственно с ними
- Обильный синтаксический сахар для работы со структурами



## Examples

``` go
set := lik.BuildSet("id", id, "url", url, "active=true")
set.SetValue("power", 17.5)
tags := set.AddList("tags")
tags.AddItems("cache", "reset")
```

Объект `set`, результат:

``` json
{
  "id": "stream37",
  "url": "http://12.34.56.78:9090",
  "power": 17.5,
  "active": true,
  "tags": [
    "cache",
    "reset"
  ]
}
```

## Interfaces

- LikItem - Общий базовый интерфейс динамических объектов
- LikSet - Интерфейс структур
- LikList - Интерфейс массивов

## Interface LikItem

Базовый интерфейс для всех динамических объектов

### BuildItem(data interface{}) Itemer

Создаётся новый динамический объект и возвращается его интерфейс

В параметре `data` может быть указано:

- целое число int, uint, int32, uint32, int64, uint64 или производный от них тип
- плавающее число float32, float64 или пропизводный
- строка string или производный
- булевский bool или производный
- интерфейс Itemer или объект совместимого типа
- интерфейс Lister или объект совместимого типа
- интерфейс Seter или объект совместимого типа

### item.Clone() Itemer

Клонировать экземпляр объекта, создаётся точная копия, тип результата - интерфейс Itemer,
можно привести к необходимому типу.

``` go
iamitem := data.Clone()
iamset := iamitem.ToSet()
```

### item.Format(prefix string) string

Форматированная сериализация, аналогично функции `item.Serialize()`, но результат форматирован как
многострочный текст для удобного восприятия человеком.

`prefix` - строка, которая будет добавлена в начало всех строк результата, может быть пустой.

### IsBool() bool

Возвращает истину, если интерфейс указывает на объект булевского типа (bool)

### IsInt() bool

Возвращает истину, если интерфейс указывает на объект типа целое число (int64)

### IsFloat() bool

Возвращает истину, если интерфейс указывает на объект типа плавающее число (float64)

### IsString() bool

Возвращает истину, если интерфейс указывает на объект типа строка (string)

### IsList() bool

Возвращает истину, если интерфейс указывает на объект типа список (Lister)

### IsSet() bool

Возвращает истину, если интерфейс указывает на объект типа структура (Seter)

### Serialize() string

Сериализация, преобразование динамического объекта в строку, аналогичную сериализованному JSON объекту.

- Объекты с простым скалярным значением будут представлены простым строковым значением, например, `25`, `true`, `"Hello"`
- Массивы будут представлены строкой, обрамлённой квадратными скобками, например, `[12, 25, 0, true]`
- Структуры будут представлены строкой, обрамлённой фигурным скобками, например, `{"id":"bran","volume":25}`

### ToBool() bool

Преобразовает объект к булевскому типу и возвращяет указатель на результат.

- целые и плавающие числа возвращают true если они не равны 0
- строки возвращают `true`, если они равны "true", `false` если "false", иначе `true`, если строка не пустая
- массивы и структуры всегда вовращают `false`

### ToInt() int64

Преобразовает объект к целому типу и возвращяет указатель на результат.

- плавающие числа округляются до ближайшего целого
- строки возвращают число, если они корректно представляют число, иначе 0
- массивы и структуры всегда вовращают 0

### ToFloat() float64

Преобразовает объект к плавающему типу и возвращяет указатель на результат.

- строки возвращают число, если они корректно представляют число, иначе 0
- массивы и структуры всегда вовращают 0

### ToString() string

Преобразовает объект к строковому типу и возвращяет указатель на результат.

- булевские значения преобразовываются в "true" или "false"
- целые и плавающие числа преобразуются в своё строковое представление
- массивы и структуры всегда вовращают ""

### ToList() Lister

Преобразовает объект к типу массива (Lister) и возвращяет указатель на результат.

Если объект не относится к типу Lister, возвращается nil

### ToSet() Seter

Преобразовает объект к типу структуры (Seter) и возвращяет указатель на результат.

Если объект не относится к типу Seter, возвращается nil

## Interface LikSet

	Count() int
	Seek(key string) int
	IsItem(path string) bool
	GetItem(path string) Itemer
	GetBool(path string) bool
	GetInt(path string) int64
	GetFloat(path string) float64
	GetString(path string) string
	GetList(path string) Lister
	GetSet(path string) Seter
	GetIDB(path string) IDB
	DelItem(path string) bool
	SetValue(path string, val interface{}) bool
	SetValues(vals ...interface{})
	AddSet(path string) Seter
	AddList(path string) Lister
	DelPos(pos int) bool
	Merge(set Seter)
	ToJson() string
	Values() []SetElm
	Keys() []string
	SortKeys() []string
	Self() *DItemSet
	SetString(key string, val string)

## Interface LikList

	Count() int
	GetItem(idx int) Itemer
	GetBool(idx int) bool
	GetInt(idx int) int64
	GetFloat(idx int) float64
	GetString(idx int) string
	GetList(idx int) Lister
	GetSet(idx int) Seter
	GetIDB(idx int) IDB
	AddItems(vals ...interface{})
	AddItemers(vals []Itemer)
	InsertItem(val interface{}, idx int)
	AddItemSet(vals ...interface{}) Seter
	SetValue(idx int, val interface{}) bool
	DelItem(idx int) bool
	SwapItem(pos1 int, pos2 int)
	ToCsv(dlm string) string
	Values() []Itemer
	Self() *DItemList
