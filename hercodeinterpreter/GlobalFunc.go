package hercodeinterpreter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ==================== 解释器实现 ====================

// 解析表达式
func parseExpression(exprStr string, LineNum int) (Expression, error) {
	exprStr = strings.TrimSpace(exprStr)
	// 检查是否为空表达式
	if exprStr == "" {
		return nil, fmt.Errorf("行 %d, 空表达式", LineNum)
	}
	// 调试输出
	//fmt.Printf("解析表达式: %s\n", exprStr)

	// 尝试解析为字面量或变量引用
	if expr, err := parseSimpleExpression(exprStr, LineNum); err == nil {
		return expr, nil
	}

	// 尝试解析为函数调用
	if expr, err := parseFunctionCall(exprStr, LineNum); err == nil {
		return expr, nil
	}

	// 尝试解析为二元运算
	operators := []string{"==", "!=", "<=", ">=", "<", ">", "+", "-", "*", "/", "%"}
	for _, op := range operators {
		parts := splitByOperator(exprStr, op)
		if len(parts) > 1 {
			left, err := parseExpression(strings.TrimSpace(parts[0]), LineNum)
			if err != nil {
				return nil, fmt.Errorf("解析左侧表达式失败: %v", err)
			}

			right, err := parseExpression(strings.TrimSpace(strings.Join(parts[1:], op)), LineNum)
			if err != nil {
				return nil, fmt.Errorf("解析右侧表达式失败: %v", err)
			}

			return &BinOpExpr{Left: left, Operator: op, Right: right}, nil
		}
	}

	return nil, fmt.Errorf("行 %d, 无法解析表达式: %s", LineNum, exprStr)
}

// 解析简单表达式（字面量、变量）
func parseSimpleExpression(exprStr string, lineNum int) (Expression, error) {
	// 数字字面量
	if num, err := strconv.ParseFloat(exprStr, 64); err == nil {
		return &LiteralExpr{Value{Type: NumberType, Num: num}}, nil
	}

	// 字符串字面量
	if len(exprStr) >= 2 && exprStr[0] == '"' && exprStr[len(exprStr)-1] == '"' {
		return &LiteralExpr{Value{Type: StringType, Str: exprStr[1 : len(exprStr)-1]}}, nil
	}

	// 布尔字面量
	if exprStr == "true" {
		return &LiteralExpr{Value{Type: BoolType, Bool: true}}, nil
	}
	if exprStr == "false" {
		return &LiteralExpr{Value{Type: BoolType, Bool: false}}, nil
	}

	// 变量引用
	if regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(exprStr) {
		return &VarRefExpr{Name: exprStr}, nil
	}

	return nil, fmt.Errorf("行 %d, 不是简单表达式", lineNum)
}

// 解析函数调用
func parseFunctionCall(exprStr string, lineNum int) (Expression, error) {
	// 函数调用
	if matches := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\((.*)\)$`).FindStringSubmatch(exprStr); matches != nil {
		funcName := matches[1]
		argsStr := matches[2]

		// 解析参数
		var args []Expression
		if argsStr != "" {
			argStrs := splitArguments(argsStr)
			for _, argStr := range argStrs {
				argExpr, err := parseExpression(strings.TrimSpace(argStr), lineNum)
				if err != nil {
					return nil, err
				}
				args = append(args, argExpr)
			}
		}

		return &FuncCallExpr{Name: funcName, Arguments: args}, nil
	}

	return nil, fmt.Errorf("行 %d, 不是函数调用", lineNum)
}

// 解析语句
// 解析语句 - 增强条件语句处理
func parseStatement(stmtStr string, lineNum int) (Statement, error) {
	stmtStr = strings.TrimSpace(stmtStr)

	// 处理 if 语句
	if strings.HasPrefix(stmtStr, "if ") {
		condStr := strings.TrimSpace(strings.TrimPrefix(stmtStr, "if"))
		// 确保移除末尾冒号
		condStr = strings.TrimSuffix(condStr, ":")
		//fmt.Printf("解析 if 条件: %s\n", condStr)

		condExpr, err := parseExpression(condStr, lineNum)
		if err != nil {
			return nil, fmt.Errorf("行 %d, 解析条件表达式失败: %v", lineNum, err)
		}
		return &IfStmt{
			Condition:  condExpr,
			ThenBranch: []Statement{},
		}, nil
	}

	// 处理 while 语句
	if strings.HasPrefix(stmtStr, "while ") {
		condStr := strings.TrimSpace(strings.TrimPrefix(stmtStr, "while"))
		// 确保移除末尾冒号
		condStr = strings.TrimSuffix(condStr, ":")
		//fmt.Printf("解析 while 条件: %s\n", condStr)

		condExpr, err := parseExpression(condStr, lineNum)
		if err != nil {
			return nil, fmt.Errorf("行 %d, 解析条件表达式失败: %v", lineNum, err)
		}
		return &WhileStmt{
			Condition: condExpr,
			Body:      []Statement{},
		}, nil
	}

	// 赋值语句
	if matches := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*(.*)$`).FindStringSubmatch(stmtStr); matches != nil {
		varName := matches[1]
		expr, err := parseExpression(strings.TrimSpace(matches[2]), lineNum)
		if err != nil {
			return nil, fmt.Errorf("行 %d, 解析赋值表达式错误: %v", lineNum, err)
		}
		return &AssignStmt{VarName: varName, Expr: expr}, nil
	}

	// Say语句
	if strings.HasPrefix(stmtStr, "say ") {
		exprStr := strings.TrimSpace(strings.TrimPrefix(stmtStr, "say"))
		expr, err := parseExpression(exprStr, lineNum)
		if err != nil {
			return nil, err
		}
		return &SayStmt{Expr: expr}, nil
	}

	// 变量声明
	if strings.HasPrefix(stmtStr, "var ") {
		parts := strings.SplitN(strings.TrimPrefix(stmtStr, "var "), "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("行 %d, 无效的变量声明: %s", lineNum, stmtStr)
		}

		varName := strings.TrimSpace(parts[0])
		expr, err := parseExpression(strings.TrimSpace(parts[1]), lineNum)
		if err != nil {
			return nil, err
		}

		return &VarDeclStmt{VarName: varName, Expr: expr}, nil
	}

	// 返回语句
	if strings.HasPrefix(stmtStr, "return ") {
		exprStr := strings.TrimSpace(strings.TrimPrefix(stmtStr, "return"))
		expr, err := parseExpression(exprStr, lineNum)
		if err != nil {
			return nil, err
		}
		return &ReturnStmt{Expr: expr}, nil
	}

	// 函数调用
	if regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\(.*\)$`).MatchString(stmtStr) {
		expr, err := parseExpression(stmtStr, lineNum)
		if err != nil {
			return nil, err
		}
		fnCall, ok := expr.(*FuncCallExpr)
		if !ok {
			return nil, fmt.Errorf("行 %d, 无效的函数调用: %s", lineNum, stmtStr)
		}
		return &FuncCallStmt{Name: fnCall.Name, Arguments: fnCall.Arguments}, nil
	}

	// 变量引用（作为函数调用）
	if regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString(stmtStr) {
		return &FuncCallStmt{Name: stmtStr, Arguments: []Expression{}}, nil
	}

	return nil, fmt.Errorf("行 %d, 无法解析语句: %s", lineNum, stmtStr)
}
