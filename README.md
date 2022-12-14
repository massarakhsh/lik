# Lik

Управление динамическими объектами

- [Lik](#lik)
  - [Examples](#examples)
  - [Interfaces](#interfaces)
  - [Interface LikItem](#interface-likitem)
    - [BuildItem(data interface{}) Itemer](#builditemdata-interface-itemer)
    - [item.Clone() Itemer](#itemclone-itemer)
    - [item.IsBool() bool](#itemisbool-bool)
    - [item.IsInt() bool](#itemisint-bool)
    - [item.IsFloat() bool](#itemisfloat-bool)
    - [item.IsString() bool](#itemisstring-bool)
    - [item.IsList() bool](#itemislist-bool)
    - [item.IsSet() bool](#itemisset-bool)
    - [item.ToBool() bool](#itemtobool-bool)
    - [item.ToInt() int64](#itemtoint-int64)
    - [item.ToFloat() float64](#itemtofloat-float64)
    - [item.ToString() string](#itemtostring-string)
    - [item.ToList() Lister](#itemtolist-lister)
    - [item.ToSet() Seter](#itemtoset-seter)
    - [item.Serialize() string](#itemserialize-string)
    - [item.Format(prefix string) string](#itemformatprefix-string-string)
  - [Interface LikSet](#interface-likset)
    - [BuildSet(data ...interface{}) Set](#buildsetdata-interface-set)
    - [BuildStringSet(data ...string) Set](#buildstringsetdata-string-set)
    - [set.Count() int](#setcount-int)
    - [set.Seek(key string) int](#setseekkey-string-int)
    - [set.IsItem(path string) bool](#setisitempath-string-bool)
    - [set.GetItem(path string) Itemer](#setgetitempath-string-itemer)
    - [set.GetBool(path string) bool](#setgetboolpath-string-bool)
    - [set.GetInt(path string) int64](#setgetintpath-string-int64)
    - [set.GetFloat(path string) float64](#setgetfloatpath-string-float64)
    - [set.GetString(path string) string](#setgetstringpath-string-string)
    - [set.GetList(path string) Lister](#setgetlistpath-string-lister)
    - [set.GetSet(path string) Seter](#setgetsetpath-string-seter)
    - [set.DelItem(path string) bool](#setdelitempath-string-bool)
    - [set.SetValue(path string, val interface{}) bool](#setsetvaluepath-string-val-interface-bool)
    - [set.SetValues(vals ...interface{})](#setsetvaluesvals-interface)
    - [set.AddSet(path string) Seter](#setaddsetpath-string-seter)
    - [set.AddList(path string) Lister](#setaddlistpath-string-lister)
    - [set.DelPos(pos int) bool](#setdelpospos-int-bool)
    - [set.Merge(set Seter)](#setmergeset-seter)
    - [set.ToJson() string](#settojson-string)
    - [set.Values() \[\]SetElm](#setvalues-setelm)
    - [set.Keys() \[\]string](#setkeys-string)
    - [set.SortKeys() \[\]string](#setsortkeys-string)
    - [set.Self() \*DItemSet](#setself-ditemset)
    - [set.SetString(key string, val string)](#setsetstringkey-string-val-string)
  - [Interface LikList](#interface-liklist)
    - [BuildList(data ...interface{}) List](#buildlistdata-interface-list)
    - [list.Count() int](#listcount-int)
    - [list.GetItem(idx int) Itemer](#listgetitemidx-int-itemer)
    - [list.GetBool(idx int) bool](#listgetboolidx-int-bool)
    - [list.GetInt(idx int) int64](#listgetintidx-int-int64)
    - [list.GetFloat(idx int) float64](#listgetfloatidx-int-float64)
    - [list.GetString(idx int) string](#listgetstringidx-int-string)
    - [list.GetList(idx int) Lister](#listgetlistidx-int-lister)
    - [list.GetSet(idx int) Seter](#listgetsetidx-int-seter)
    - [list.AddItems(vals ...interface{})](#listadditemsvals-interface)
    - [list.AddItemers(vals \[\]Itemer)](#listadditemersvals-itemer)
    - [list.InsertItem(val interface{}, idx int)](#listinsertitemval-interface-idx-int)
    - [list.AddItemSet(vals ...interface{}) Seter](#listadditemsetvals-interface-seter)
    - [list.SetValue(idx int, val interface{}) bool](#listsetvalueidx-int-val-interface-bool)
    - [list.DelItem(idx int) bool](#listdelitemidx-int-bool)
    - [list.SwapItem(pos1 int, pos2 int)](#listswapitempos1-int-pos2-int)
    - [list.Values() \[\]Itemer](#listvalues-itemer)
    - [list.Self() \*DItemList](#listself-ditemlist)

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
set.SetValue("power/value", 17.5)
tags := set.AddList("tags")
tags.AddItems("cache", "reset")
```

Объект `set`, результат:

``` json
{
  "id": "stream37",
  "url": "http://12.34.56.78:9090",
  "power": {
    "value": {
      17.5
    }
  },
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

### item.IsBool() bool

Возвращает истину, если интерфейс указывает на объект булевского типа (bool)

### item.IsInt() bool

Возвращает истину, если интерфейс указывает на объект типа целое число (int64)

### item.IsFloat() bool

Возвращает истину, если интерфейс указывает на объект типа плавающее число (float64)

### item.IsString() bool

Возвращает истину, если интерфейс указывает на объект типа строка (string)

### item.IsList() bool

Возвращает истину, если интерфейс указывает на объект типа список (Lister)

### item.IsSet() bool

Возвращает истину, если интерфейс указывает на объект типа структура (Seter)

### item.ToBool() bool

Преобразовает объект к булевскому типу и возвращяет указатель на результат.

- целые и плавающие числа возвращают true если они не равны 0
- строки возвращают `true`, если они равны "true", `false` если "false", иначе `true`, если строка не пустая
- массивы и структуры всегда вовращают `false`

### item.ToInt() int64

Преобразовает объект к целому типу и возвращяет указатель на результат.

- плавающие числа округляются до ближайшего целого
- строки возвращают число, если они корректно представляют число, иначе 0
- массивы и структуры всегда вовращают 0

### item.ToFloat() float64

Преобразовает объект к плавающему типу и возвращяет указатель на результат.

- строки возвращают число, если они корректно представляют число, иначе 0
- массивы и структуры всегда вовращают 0

### item.ToString() string

Преобразовает объект к строковому типу и возвращяет указатель на результат.

- булевские значения преобразовываются в "true" или "false"
- целые и плавающие числа преобразуются в своё строковое представление
- массивы и структуры всегда вовращают ""

### item.ToList() Lister

Преобразовает объект к типу массива (Lister) и возвращяет указатель на результат.

Если объект не относится к типу Lister, возвращается nil

### item.ToSet() Seter

Преобразовает объект к типу структуры (Seter) и возвращяет указатель на результат.

Если объект не относится к типу Seter, возвращается nil

### item.Serialize() string

Сериализация, преобразование динамического объекта в строку, аналогичную сериализованному JSON объекту.

- Объекты с простым скалярным значением будут представлены простым строковым значением, например, `25`, `true`, `"Hello"`
- Массивы будут представлены строкой, обрамлённой квадратными скобками, например, `[12, 25, 0, true]`
- Структуры будут представлены строкой, обрамлённой фигурным скобками, например, `{"id":"bran","volume":25}`

### item.Format(prefix string) string

Форматированная сериализация, аналогично функции `item.Serialize()`, но результат форматирован как
многострочный текст для удобного восприятия человеком.

`prefix` - строка, которая будет добавлена в начало всех строк результата, может быть пустой.

## Interface LikSet

Интерфейс динамического объекта типа "структура"

Структура содержит произвольное количество именованных полей, значение каждого из которых
является динамическим объектом. Таким образом, вложенные структуры образуют иерархию.
К внутренним полям можно обращаться, указывая полное составное имя так же, как образуются
имена вложенных файлов, через символ следжа, например:

``` go
set.SetValue("key", "0123456789")
set.setValue("options/main/autorun", true)
```

### BuildSet(data ...interface{}) Set

Создаётся новый объект-структура и возвращается его интерфейс

В параметрах, если они указаны, можно задать поля и значения, которыми инициализируется структура.
Можно либо последовательно указывать имя и значение поля, либо в одной строке в форме `имя=значение`, например:

``` go
set0 := lik.BuildSet()
set1 := lik.BuildSet("alpha=1", "beta=2", "autostart=true")
set2 := lik.BuildSet("id", myId, "value", set1)
```

Как значения могут быть указаны:

- целое число int, uint, int32, uint32, int64, uint64 или производный от них тип
- плавающее число float32, float64 или пропизводный
- строка string или производный
- булевский bool или производный
- интерфейс Itemer или объект совместимого типа
- интерфейс Lister или объект совместимого типа
- интерфейс Seter или объект совместимого типа

### BuildStringSet(data ...string) Set

Аналогично функции `BuildSet`, но в параметрах указываются только строки

### set.Count() int

Возвращает количество полей в структуре

### set.Seek(key string) int

Находит, какую позицию занимает поле с именем `key` в списке всех полей, если не найдено -1

### set.IsItem(path string) bool

Проверяет, присутствует ли в структуре поле с именем `path`

В качестве имени может быть указан полный путь поля, как и во всех следующих функциях.

### set.GetItem(path string) Itemer

Возвращает интерфейс объекта в поле с именем `path`, если не найдено - nil

### set.GetBool(path string) bool

Возвращает значение поля `path` как булевское значение, преобразование аналогично `ToBool`, если нет - `false`

### set.GetInt(path string) int64

Возвращает значение поля `path` как целое значение, преобразование аналогично `ToInt`, если нет - `0`

### set.GetFloat(path string) float64

Возвращает значение поля `path` как значение с плавающей точкой, преобразование аналогично `ToFloat`, если нет - `0`

### set.GetString(path string) string

Возвращает значение поля `path` как значение строки, преобразование аналогично `ToString`, если нет - пустая строка

### set.GetList(path string) Lister

Возвращает значение поля `path` как массив, если нет или не массив, возвращается nil

### set.GetSet(path string) Seter

Возвращает значение поля `path` как струкура, если нет или не структура, возвращается nil

### set.DelItem(path string) bool

Удаляет поле с именем `path`, возвращается true, если изменения были внесены

### set.SetValue(path string, val interface{}) bool

Устанавливает поле с именем `path` в значение `val`, можно указать любое значение, допустмое в функции BuildItem()

### set.SetValues(vals ...interface{})

Устанавливает несколько значений, синтаксис аналогичен BuildSet()

### set.AddSet(path string) Seter

Добавляет структуру в поле с путём `path`, возвращает интерфейс на новую структуру

### set.AddList(path string) Lister

Добавляет массив в поле с путём `path`, возвращает интерфейс на новый массив

### set.DelPos(pos int) bool

Удаляет поле по относительной позиции `pos`

### set.Merge(set Seter)

Объединяет структуры, добавляя в текущую все поля из структуры `set`, одноимённые поля заменяются.

### set.ToJson() string

Аналогично функции item.Serialize()

### set.Values() []SetElm

Возвращает массив всех полей струкуры в формате []ItSet

### set.Keys() []string

Возвращает массив всех имён полей в формате []string

### set.SortKeys() []string

Возвращает отсортированный массив всех имён полей в формате []string

### set.Self() *DItemSet

Возвращает указатель на объект, которому принадлежит интерфейс

### set.SetString(key string, val string)

Устанавливает поле `key` в значение `val`, при этом автомаически преобразовываются значения `true`, `false`,
а также целые т плавающие числа

## Interface LikList

Интерфейс динамического объекта типа "массив"

Массив содержит упорядоченный набор неименованных полей, значение каждого из которых
является динамическим объектом. Обращение к элементам массива происходит по индексу,
числу, начинающемуся с 0.

При указании путей в динамическом объекте индексы могут присутствовать в позициях отдельных имён.

``` go
list.SetValue(5, "0123456789")
set.setValue("options/pool/5/autorun", true)
```

### BuildList(data ...interface{}) List

Создаётся новый объект-массив и возвращается его интерфейс

В параметрах, если они указаны, можно указать объекты, которые заносятся в массив

``` go
list0 := lik.BuildList()
list1 := lik.BuildList(a, b, c, 16, "hoo")
list2 := lik.BuildList("alpha", "beta", list1)
```

Как значения могут быть указаны:

- целое число int, uint, int32, uint32, int64, uint64 или производный от них тип
- плавающее число float32, float64 или пропизводный
- строка string или производный
- булевский bool или производный
- интерфейс Itemer или объект совместимого типа
- интерфейс Lister или объект совместимого типа
- интерфейс Seter или объект совместимого типа

### list.Count() int

Возвращает текущее количество элементов массива

### list.GetItem(idx int) Itemer

Возвращает объект, элемент номер `idx`

### list.GetBool(idx int) bool

Возвращает элемент массива как булевское

### list.GetInt(idx int) int64

Возвращает элемент массива как целое

### list.GetFloat(idx int) float64

Возвращает элемент массива как плавающее

### list.GetString(idx int) string

Возвращает элемент массива как строку

### list.GetList(idx int) Lister

Возвращает элемент массива как массив, если это не массив, то nil

### list.GetSet(idx int) Seter

Возвращает элемент массива как структуру, если это не структура, то nil

### list.AddItems(vals ...interface{})

Добавляет в массив новые элементы

### list.AddItemers(vals []Itemer)

Добавляет в массив новые элементы

### list.InsertItem(val interface{}, idx int)

Вставляет новый элемент по указанному индексу

### list.AddItemSet(vals ...interface{}) Seter

Создаёт новую струкуру и добавляет в массив, синтаксис аналогичен BuildSet()

### list.SetValue(idx int, val interface{}) bool

Устанавливает значение элемента массива, возвращает true, если были внесены изменения

### list.DelItem(idx int) bool

Удаляет элемент по номеру

### list.SwapItem(pos1 int, pos2 int)

Меняет местами два элемента

### list.Values() []Itemer

Возвращает массив элементов массива

### list.Self() *DItemList

Возвращает указатель на объект, которому принадлежит интерфейс
