package hercodeinterpreter

import "fmt"

// 返回语句
type ReturnStmt struct {
	Expr Expression
}

func (s *ReturnStmt) Execute(ctx *Context) (Value, error) {
	val, err := s.Expr.Eval(ctx)
	if err != nil {
		return Value{}, err
	}
	if val.Type == ErrorType {
		return val, nil
	}
	return val, nil
}

func (s *ReturnStmt) String() string {
	return fmt.Sprintf("return %s", s.Expr)
}
