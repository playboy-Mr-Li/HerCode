package hercodeinterpreter

import (
	"fmt"
	"strings"
)

// 值类型
type ValueType int

const (
	NumberType ValueType = iota
	StringType
	BoolType
	VoidType
	SliceType
	MapType
	FunctionType
	ErrorType
)

// 值结构
type Value struct {
	Type  ValueType
	Num   float64
	Str   string
	Bool  bool
	Slice []Value
	Map   map[string]Value
	Func  *HerCodeFunction
	Error error
}

func (e Value) String() string {
	switch e.Type {
	case NumberType:
		return fmt.Sprintf("%f", e.Num)
	case StringType:
		return fmt.Sprintf("%s", e.Str)
	case BoolType:
		return fmt.Sprintf("%t", e.Bool)
	case MapType:
		var m []string
		for k, v := range e.Map {
			m = append(m, k+": "+v.String())

		}
		return strings.Join(m, ", ")

	case SliceType:
		var s []string
		for _, v := range e.Slice {
			s = append(s, v.String())
		}
		return "[" + strings.Join(s, ", ") + "]"

	case FunctionType:
		return fmt.Sprintf("%s", e.String())
	default:
		return ""

	}
}

// 表达式接口
type Expression interface {
	String() string
	Eval(ctx *Context) (Value, error)
}

// 语句接口
type Statement interface {
	String() string
	Execute(ctx *Context) (Value, error)
}
