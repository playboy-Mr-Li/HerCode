package hercodeinterpreter

import "fmt"

// 赋值语句
type AssignStmt struct {
	VarName string
	Expr    Expression
}

func (s *AssignStmt) Execute(ctx *Context) (Value, error) {
	val, err := s.Expr.Eval(ctx)
	if err != nil {
		return Value{}, err
	}
	if val.Type == ErrorType {
		return val, nil
	}
	ctx.SetVar(s.VarName, val)
	return Value{Type: VoidType}, nil
}
func (s *AssignStmt) String() string {
	return fmt.Sprintf("%s = %s", s.VarName, s.Expr)
}
