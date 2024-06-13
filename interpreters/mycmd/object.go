package mycmd

import "strconv"

// ObjectType определяет тип объекта в интерпретаторе
type ObjectType string

// Object представляет объект в интерпретаторе
type Object interface {
	Type() ObjectType // Метод возвращает тип объекта
	Inspect() string  // Метод возвращает строковое представление объекта
}

// Пример реализации объекта Integer
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return "INTEGER"
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}

// Пример реализации объекта Boolean
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return "BOOLEAN"
}

func (b *Boolean) Inspect() string {
	return strconv.FormatBool(b.Value)
}

// Пример реализации объекта Null
type Null struct{}

func (n *Null) Type() ObjectType {
	return "NULL"
}

func (n *Null) Inspect() string {
	return "null"
}
