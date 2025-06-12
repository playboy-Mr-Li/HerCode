package hercodeinterpreter

import (
	"fmt"
	"math"
	"strings"
)

// 函数调用表达式
type FuncCallExpr struct {
	Name      string
	Arguments []Expression
}

func (e *FuncCallExpr) String() string {
	var m []string
	for _, arg := range e.Arguments {
		m = append(m, arg.String())
	}
	return fmt.Sprintf("%s(%v)}", e.Name, strings.Join(m, ","))

}

func (e *FuncCallExpr) Eval(ctx *Context) (Value, error) {
	// 查找函数
	//fmt.Printf("正在执行函数：%s 参数：%v\n", e.Name, e.Arguments)

	fn, ok := ctx.GetFunc(e.Name)
	if !ok {
		return Value{Type: ErrorType, Error: fmt.Errorf("函数未定义: %s", e.Name)}, nil
	}

	// 计算参数值
	args := make([]Value, len(e.Arguments))
	for i, argExpr := range e.Arguments {
		argVal, err := argExpr.Eval(ctx)
		if err != nil {
			return Value{}, err
		}
		if argVal.Type == ErrorType {
			return argVal, nil
		}
		args[i] = argVal
	}

	// 处理内置函数
	if e.Name == "len" {
		if len(args) != 1 {
			return Value{Type: ErrorType, Error: fmt.Errorf("len() 需要1个参数")}, nil
		}
		if args[0].Type != StringType {
			return Value{Type: ErrorType, Error: fmt.Errorf("len() 需要字符串参数")}, nil
		}
		return Value{Type: NumberType, Num: float64(len(args[0].Str))}, nil
	}

	if e.Name == "substr" {
		if len(args) < 2 || len(args) > 3 {
			return Value{Type: ErrorType, Error: fmt.Errorf("substr() 需要2-3个参数")}, nil
		}

		if args[0].Type != StringType {
			return Value{Type: ErrorType, Error: fmt.Errorf("substr() 第一个参数必须是字符串")}, nil
		}

		if args[1].Type != NumberType {
			return Value{Type: ErrorType, Error: fmt.Errorf("substr() 第二个参数必须是数字")}, nil
		}

		start := int(args[1].Num)
		if start < 0 || start >= len(args[0].Str) {
			return Value{Type: ErrorType, Error: fmt.Errorf("substr() 起始位置超出范围")}, nil
		}

		end := len(args[0].Str)
		if len(args) == 3 {
			if args[2].Type != NumberType {
				return Value{Type: ErrorType, Error: fmt.Errorf("substr() 第三个参数必须是数字")}, nil
			}
			end = int(args[2].Num)
			if end < start || end > len(args[0].Str) {
				return Value{Type: ErrorType, Error: fmt.Errorf("substr() 结束位置超出范围")}, nil
			}
		}

		return Value{Type: StringType, Str: args[0].Str[start:end]}, nil
	}

	if e.Name == "sqrt" {
		if len(args) != 1 {
			return Value{Type: ErrorType, Error: fmt.Errorf("sqrt() 需要1个参数")}, nil
		}
		if args[0].Type != NumberType {
			return Value{Type: ErrorType, Error: fmt.Errorf("sqrt() 需要数字参数")}, nil
		}
		if args[0].Num < 0 {
			return Value{Type: ErrorType, Error: fmt.Errorf("sqrt() 参数不能为负数")}, nil
		}
		return Value{Type: NumberType, Num: math.Sqrt(args[0].Num)}, nil
	}

	// 创建新的执行上下文
	localCtx := NewContext(ctx)

	// 设置参数
	for i, param := range fn.Parameters {
		if i < len(args) {
			localCtx.SetVar(param, args[i])
		} else {
			localCtx.SetVar(param, Value{Type: VoidType})
		}
	}

	// 执行函数体
	for _, stmt := range fn.Statements {

		result, err := stmt.Execute(localCtx)
		//defer fmt.Printf("执行函数：%s 参数：[%v] 结果：[%v], 错误：[%v]\n", fn.Name, args, result, err)
		if err != nil {
			return Value{}, err
		}

		// 如果有返回值，则返回
		if result.Type != VoidType {
			return result, nil
		}
	}

	return Value{Type: VoidType}, nil
}
