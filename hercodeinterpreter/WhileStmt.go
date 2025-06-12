package hercodeinterpreter

import "fmt"

// 循环语句
type WhileStmt struct {
	Condition Expression
	Body      []Statement
}

func (s *WhileStmt) String() string {
	var bodyStr string
	for _, stmt := range s.Body {
		bodyStr += "    " + stmt.String() + "\n"
	}
	return fmt.Sprintf("while %s {\n%s}\n", s.Condition.String(), bodyStr)
}

func (s *WhileStmt) Execute(ctx *Context) (Value, error) {
	for {
		condVal, err := s.Condition.Eval(ctx)
		if err != nil {
			return Value{}, err
		}
		if condVal.Type == ErrorType {
			return condVal, nil
		}

		if condVal.Type != BoolType {
			return Value{Type: ErrorType, Error: fmt.Errorf("条件表达式必须为布尔类型")}, nil
		}

		if !condVal.Bool {
			break
		}

		for _, stmt := range s.Body {
			result, err := stmt.Execute(ctx)
			if err != nil {
				return Value{}, err
			}
			// 如果有返回值，则返回
			if result.Type != VoidType {
				return result, nil
			}
		}
	}
	return Value{Type: VoidType}, nil
}
