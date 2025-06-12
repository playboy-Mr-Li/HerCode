package hercodeinterpreter

import (
	"fmt"
	"strings"
)

// 函数调用语句
type FuncCallStmt struct {
	Name      string
	Arguments []Expression
}

func (s *FuncCallStmt) String() string {
	if len(s.Arguments) == 0 {
		return fmt.Sprintf("%s()", s.Name)
	}
	args := make([]string, len(s.Arguments))
	for i, arg := range s.Arguments {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", s.Name, strings.Join(args, ","))

}

func (s *FuncCallStmt) Execute(ctx *Context) (Value, error) {
	// 创建函数调用表达式
	callExpr := &FuncCallExpr{
		Name:      s.Name,
		Arguments: s.Arguments,
	}
	//fmt.Printf("调用函数：%s, 参数: %s", s.Name, s.Arguments)
	// 执行函数调用
	result, _ := callExpr.Eval(ctx)

	if result.Error != nil {
		return Value{}, result.Error
	}
	if result.Type != VoidType {
		return result, result.Error
	}

	// 忽略返回值
	return Value{Type: VoidType}, result.Error
}
