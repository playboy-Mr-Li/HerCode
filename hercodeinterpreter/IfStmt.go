package hercodeinterpreter

import (
	"fmt"
	"strings"
)

// 条件语句
type IfStmt struct {
	Condition  Expression
	ThenBranch []Statement
	ElseBranch []Statement
}

func (s *IfStmt) String() string {
	r := strings.Builder{}
	r.WriteString("if " + s.Condition.String() + " {\n")
	if s.ThenBranch == nil {
		return ""
	} else {
		for _, stmt := range s.ThenBranch {
			r.WriteString("    ")

			r.WriteString(stmt.String() + "\n")

		}
		r.WriteString("}\n")
	}

	if s.ElseBranch != nil {
		r.WriteString("else {\n")
		for _, stmt := range s.ElseBranch {
			r.WriteString("    ")
			r.WriteString(stmt.String() + "\n")

		}
		r.WriteString("}\n")
	}
	return r.String()
}

func (s *IfStmt) Execute(ctx *Context) (Value, error) {
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

	if condVal.Bool {
		for _, stmt := range s.ThenBranch {
			result, err := stmt.Execute(ctx)
			if err != nil {
				return Value{}, err
			}
			// 如果有返回值，则返回
			if result.Type != VoidType {
				return result, nil
			}
		}
	} else if s.ElseBranch != nil {
		for _, stmt := range s.ElseBranch {
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
