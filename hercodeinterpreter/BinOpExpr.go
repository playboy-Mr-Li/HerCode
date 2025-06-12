package hercodeinterpreter

import "fmt"

// 二元运算表达式
type BinOpExpr struct {
	Left     Expression
	Operator string
	Right    Expression
	lineNum  int
}

func (e *BinOpExpr) String() string {
	return fmt.Sprintf("%s %s %s", e.Left, e.Operator, e.Right)
}
func (e *BinOpExpr) Eval(ctx *Context) (Value, error) {

	// 处理赋值操作
	if e.Operator == "=" {
		// 左侧必须是变量引用
		leftVar, ok := e.Left.(*VarRefExpr)
		if !ok {
			return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 赋值操作左侧必须是变量", e.lineNum)}, nil
		}

		// 计算右侧值
		rightVal, err := e.Right.Eval(ctx)
		if err != nil {
			return Value{}, err
		}
		if rightVal.Type == ErrorType {
			return rightVal, nil
		}

		// 设置变量值
		ctx.SetVar(leftVar.Name, rightVal)
		return rightVal, nil
	}

	leftVal, err := e.Left.Eval(ctx)
	if err != nil {
		return Value{}, err
	}
	if leftVal.Type == ErrorType {
		return leftVal, nil
	}

	rightVal, err := e.Right.Eval(ctx)
	if err != nil {
		return Value{}, err
	}
	if rightVal.Type == ErrorType {
		return rightVal, nil
	}

	switch e.Operator {
	case "+":
		if leftVal.Type == NumberType && rightVal.Type == NumberType {
			return Value{Type: NumberType, Num: leftVal.Num + rightVal.Num}, nil
		}
		if leftVal.Type == StringType || rightVal.Type == StringType {
			return Value{Type: StringType, Str: fmt.Sprintf("%v%v", leftVal, rightVal)}, nil
		}
		return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 类型不匹配: %s + %s", e.lineNum, leftVal.Type, rightVal.Type)}, nil

	case "-":
		if leftVal.Type == NumberType && rightVal.Type == NumberType {
			return Value{Type: NumberType, Num: leftVal.Num - rightVal.Num}, nil
		}
		return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 类型不匹配: %s - %s", e.lineNum, leftVal.Type, rightVal.Type)}, nil

	case "*":
		if leftVal.Type == NumberType && rightVal.Type == NumberType {
			return Value{Type: NumberType, Num: leftVal.Num * rightVal.Num}, nil
		}
		return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 类型不匹配: %s * %s", e.lineNum, leftVal.Type, rightVal.Type)}, nil

	case "/":
		if leftVal.Type == NumberType && rightVal.Type == NumberType {
			if rightVal.Num == 0 {
				return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 除以零错误", e.lineNum)}, nil
			}
			return Value{Type: NumberType, Num: leftVal.Num / rightVal.Num}, nil
		}
		return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 类型不匹配: %s / %s", e.lineNum, leftVal.Type, rightVal.Type)}, nil

	case "%":
		if leftVal.Type == NumberType && rightVal.Type == NumberType {
			if rightVal.Num == 0 {
				return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 取模运算除以零错误", e.lineNum)}, nil
			}
			return Value{Type: NumberType, Num: float64(int(leftVal.Num) % int(rightVal.Num))}, nil
		}
		return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 类型不匹配: %s %% %s", e.lineNum, leftVal.Type, rightVal.Type)}, nil

	case "==":
		if leftVal.Type == rightVal.Type {
			switch leftVal.Type {
			case NumberType:
				return Value{Type: BoolType, Bool: leftVal.Num == rightVal.Num}, nil
			case StringType:
				return Value{Type: BoolType, Bool: leftVal.Str == rightVal.Str}, nil
			case BoolType:
				return Value{Type: BoolType, Bool: leftVal.Bool == rightVal.Bool}, nil
			}
		}
		return Value{Type: BoolType, Bool: false}, nil

	case "!=":
		if leftVal.Type == rightVal.Type {
			switch leftVal.Type {
			case NumberType:
				return Value{Type: BoolType, Bool: leftVal.Num != rightVal.Num}, nil
			case StringType:
				return Value{Type: BoolType, Bool: leftVal.Str != rightVal.Str}, nil
			case BoolType:
				return Value{Type: BoolType, Bool: leftVal.Bool != rightVal.Bool}, nil
			}
		}
		return Value{Type: BoolType, Bool: true}, nil

	case "<":
		if leftVal.Type == NumberType && rightVal.Type == NumberType {
			return Value{Type: BoolType, Bool: leftVal.Num < rightVal.Num}, nil
		}
		if leftVal.Type == StringType && rightVal.Type == StringType {
			return Value{Type: BoolType, Bool: leftVal.Str < rightVal.Str}, nil
		}
		return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 类型不匹配: %s < %s", e.lineNum, leftVal.Type, rightVal.Type)}, nil

	case ">":
		if leftVal.Type == NumberType && rightVal.Type == NumberType {
			return Value{Type: BoolType, Bool: leftVal.Num > rightVal.Num}, nil
		}
		if leftVal.Type == StringType && rightVal.Type == StringType {
			return Value{Type: BoolType, Bool: leftVal.Str > rightVal.Str}, nil
		}
		return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 类型不匹配: %s > %s", e.lineNum, leftVal.Type, rightVal.Type)}, nil

	case "<=":
		if leftVal.Type == NumberType && rightVal.Type == NumberType {
			return Value{Type: BoolType, Bool: leftVal.Num <= rightVal.Num}, nil
		}
		if leftVal.Type == StringType && rightVal.Type == StringType {
			return Value{Type: BoolType, Bool: leftVal.Str <= rightVal.Str}, nil
		}
		return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 类型不匹配: %s <= %s", e.lineNum, leftVal.Type, rightVal.Type)}, nil

	case ">=":
		if leftVal.Type == NumberType && rightVal.Type == NumberType {
			return Value{Type: BoolType, Bool: leftVal.Num >= rightVal.Num}, nil
		}
		if leftVal.Type == StringType && rightVal.Type == StringType {
			return Value{Type: BoolType, Bool: leftVal.Str >= rightVal.Str}, nil
		}
		return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 类型不匹配: %s >= %s", e.lineNum, leftVal.Type, rightVal.Type)}, nil

	default:
		return Value{Type: ErrorType, Error: fmt.Errorf("行 %d, 未知运算符: %s", e.lineNum, e.Operator)}, nil
	}
}
