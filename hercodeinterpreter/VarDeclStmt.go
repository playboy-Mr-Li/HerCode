package hercodeinterpreter

import "fmt"

// 变量声明语句
type VarDeclStmt struct {
	VarName string
	Expr    Expression
}

func (s *VarDeclStmt) String() string {
	return fmt.Sprintf("var %s %s", s.VarName, s.Expr.String())
}

func (s *VarDeclStmt) Execute(ctx *Context) (Value, error) {
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
