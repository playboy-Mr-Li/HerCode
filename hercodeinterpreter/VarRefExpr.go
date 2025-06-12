package hercodeinterpreter

import "fmt"

// 变量引用表达式
type VarRefExpr struct {
	Name string
}

func (e *VarRefExpr) Eval(ctx *Context) (Value, error) {
	val, ok := ctx.GetVar(e.Name)
	if !ok {
		return Value{Type: ErrorType, Error: fmt.Errorf("变量未定义: %s", e.Name)}, nil
	}
	return val, nil
}

func (e *VarRefExpr) String() string {
	return fmt.Sprintf("%s", e.Name)
}
