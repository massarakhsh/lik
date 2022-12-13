# Библиотека lik

Управление динамическими объектами

- [Библиотека lik](#библиотека-lik)
  - [Examples](#examples)
  - [Interfaces](#interfaces)
  - [Interface LikItem](#interface-likitem)
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

Access to a common dynamic object is provided by the LikItem interface.

## Interface LikSet

## Interface LikList
