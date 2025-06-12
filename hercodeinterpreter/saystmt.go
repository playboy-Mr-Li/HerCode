package hercodeinterpreter

import "fmt"

// Say语句
type SayStmt struct {
	Expr Expression
}

func (s *SayStmt) String() string {
	return fmt.Sprintf("say %v", s.Expr.String())
}
func (s *SayStmt) Execute(ctx *Context) (Value, error) {
	val, err := s.Expr.Eval(ctx)
	if err != nil {
		return Value{}, err
	}

	if val.Type == ErrorType {
		fmt.Printf("错误: %v\n", val.Error)
	} else {
		switch val.Type {
		case NumberType:
			fmt.Println(val.Num)
		case StringType:
			fmt.Println(val.Str)
		case BoolType:
			fmt.Println(val.Bool)
		default:
			fmt.Println()
		}
	}
	return Value{Type: VoidType}, nil
}
