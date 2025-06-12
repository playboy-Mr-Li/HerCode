package hercodeinterpreter

// 字面量表达式
type LiteralExpr struct {
	Value Value
}

func (e *LiteralExpr) String() string {

	return e.Value.String()

}

func (e *LiteralExpr) Eval(ctx *Context) (Value, error) {
	return e.Value, nil
}
