package hercodeinterpreter

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

// HerCode解释器
type HerCodeInterpreter struct {
	Functions    map[string]*HerCodeFunction
	StartFunc    string
	GlobalCtx    *Context
	funcStack    []*HerCodeFunction
	blockStack   []Statement
	currentBlock *Statement // 当前处理的块（if 或 while）
	inElseBranch bool       // 标记当前是否在 else 分支中

}

// 创建新解释器
func NewHerCodeInterpreter() *HerCodeInterpreter {
	return &HerCodeInterpreter{
		Functions: make(map[string]*HerCodeFunction),
		GlobalCtx: NewContext(nil),
	}
}

// 注册内置函数
func (h *HerCodeInterpreter) registerBuiltinFunctions() {
	// len 函数
	h.GlobalCtx.SetFunc("len", &HerCodeFunction{
		Name:       "len",
		Parameters: []string{"str"},
		ReturnType: NumberType,
	})

	// substr 函数
	h.GlobalCtx.SetFunc("substr", &HerCodeFunction{
		Name:       "substr",
		Parameters: []string{"str", "start", "end"},
		ReturnType: StringType,
	})

	// sqrt 函数
	h.GlobalCtx.SetFunc("sqrt", &HerCodeFunction{
		Name:       "sqrt",
		Parameters: []string{"num"},
		ReturnType: NumberType,
	})
}

// 解析HerCode脚本
func (h *HerCodeInterpreter) Parse(script string) error {
	scanner := bufio.NewScanner(strings.NewReader(script))
	currentFunc := ""
	currentFuncStatements := []Statement{}
	h.blockStack = []Statement{}
	h.funcStack = []*HerCodeFunction{}
	h.currentBlock = nil
	h.inElseBranch = false
	// 正则表达式
	//funcRegex := regexp.MustCompile(`^function\s+([a-zA-Z_][a-zA-Z0-9_]*)\(([^)]*)\):`)
	//funcRegex := regexp.MustCompile(`^function\s+([a-zA-Z_][a-zA-Z0-9_]*)\(([^)]*)\)\s*:\s*(.*)`)
	//funcRegex := regexp.MustCompile(`^function\s+([a-zA-Z_][a-zA-Z0-9_]*)(?:\s+([a-zA-Z_][a-zA-Z0-9_]*))*\s*:`)
	//ifRegex := regexp.MustCompile(`^if\s+(.*):`)
	//elseRegex := regexp.MustCompile(`^else?`)
	//endifRegex := regexp.MustCompile(`^endif`)
	//whileRegex := regexp.MustCompile(`^当\s+(.*)`)
	//endwhileRegex := regexp.MustCompile(`^行了细狗`)
	startRegex := regexp.MustCompile(`^start:`)
	//varRegex := regexp.MustCompile(`^var\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*(.*)$`)
	//sayRegex := regexp.MustCompile(`^say\s+(.*)`)
	returnRegex := regexp.MustCompile(`^return\s+(.*)`)

	//endRegex := regexp.MustCompile(`^end`)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		//debug := func() {
		//	fmt.Printf("行 %d: %s\n", lineNum, line)
		//}
		//debug()
		line = cleanComment(line)
		line = strings.TrimSpace(line)

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		//line = cleanQuotes(line)

		// 处理函数定义
		if strings.HasPrefix(line, "function ") {
			funcName, params, err := parseFunctionDefinition(line, lineNum)
			if err != nil {
				return fmt.Errorf("行 %d: %v", lineNum, err)
			}

			// 结束之前的函数
			if currentFunc != "" {
				// 确保我们保存的是正确的函数信息
				h.GlobalCtx.SetFunc(currentFunc, &HerCodeFunction{
					Name: currentFunc,
					Parameters: func() []string {
						if len(h.funcStack) == 0 {
							return []string{}
						}
						return h.funcStack[len(h.funcStack)-1].Parameters
					}(), // 获取当前函数的实际参数
					Statements: currentFuncStatements,
				})
			}

			currentFunc = funcName
			currentFuncStatements = []Statement{}
			// 创建新函数
			newFunc := &HerCodeFunction{
				Name:       funcName,
				Parameters: params,
			}
			h.GlobalCtx.SetFunc(funcName, newFunc)
			h.funcStack = append(h.funcStack, newFunc)
			h.currentBlock = nil
			continue
		}

		// 处理start入口
		if startRegex.MatchString(line) {
			currentFunc = "start"
			currentFuncStatements = []Statement{}
			h.GlobalCtx.SetFunc("start", &HerCodeFunction{
				Name: "start",
			})
			h.funcStack = append(h.funcStack, h.GlobalCtx.Functions["start"])
			h.currentBlock = nil
			continue
		}

		// 处理else语句 - 关键修复
		if strings.HasPrefix(line, "else") {
			if len(h.blockStack) == 0 {
				return fmt.Errorf("行 %d: else 没有匹配的 if", lineNum)
			}

			// 获取栈顶元素
			topStmt := h.blockStack[len(h.blockStack)-1]
			ifStmt, ok := topStmt.(*IfStmt)
			if !ok {
				return fmt.Errorf("行 %d: else 必须紧跟在 if 之后", lineNum)
			}

			// 创建新的Else分支
			ifStmt.ElseBranch = []Statement{}

			// 设置当前块为 else 分支
			h.currentBlock = &topStmt
			h.inElseBranch = true // 标记当前在 else 分支中
			continue
		}

		// 处理endif语句
		if strings.HasPrefix(line, "endif") {
			if len(h.blockStack) == 0 {
				return fmt.Errorf("行 %极: endif 没有匹配的 if", lineNum)
			}

			// 弹出栈顶元素
			ifStmt := h.blockStack[len(h.blockStack)-1]
			h.blockStack = h.blockStack[:len(h.blockStack)-1]

			// 重置当前块和分支状态
			h.currentBlock = nil
			h.inElseBranch = false

			// 添加到当前函数或块
			if len(h.blockStack) > 0 {
				currentBlock := h.blockStack[len(h.blockStack)-1]
				switch b := currentBlock.(type) {
				case *IfStmt:
					// 根据是否在 else 分支中添加
					if h.inElseBranch {
						b.ElseBranch = append(b.ElseBranch, ifStmt)
					} else {
						b.ThenBranch = append(b.ThenBranch, ifStmt)
					}
				case *WhileStmt:
					b.Body = append(b.Body, ifStmt)
				}
			} else {
				currentFuncStatements = append(currentFuncStatements, ifStmt)
			}
			continue
		}

		if strings.HasPrefix(line, "endwhile") {
			if len(h.blockStack) == 0 {
				return fmt.Errorf("行 %d: endwhile 没有匹配的 while", lineNum)
			}

			// 弹出栈顶元素
			whileStmt := h.blockStack[len(h.blockStack)-1]
			h.blockStack = h.blockStack[:len(h.blockStack)-1]

			// 重置当前块
			h.currentBlock = nil
			h.inElseBranch = false

			// 添加到当前函数或块 - 使用 inElseBranch 判断位置
			if len(h.blockStack) > 0 {
				currentBlock := h.blockStack[len(h.blockStack)-1]
				switch b := currentBlock.(type) {
				case *IfStmt:
					// 根据是否在 else 分支中添加
					if h.inElseBranch {
						b.ElseBranch = append(b.ElseBranch, whileStmt)
					} else {
						b.ThenBranch = append(b.ThenBranch, whileStmt)
					}
				case *WhileStmt:
					b.Body = append(b.Body, whileStmt)
				}
			} else {
				currentFuncStatements = append(currentFuncStatements, whileStmt)
			}
			continue
		}

		// 处理返回语句
		if matches := returnRegex.FindStringSubmatch(line); matches != nil {
			exprStr := strings.TrimSpace(matches[1])
			expr, err := parseExpression(exprStr, lineNum)
			if err != nil {
				return fmt.Errorf("行 %d: %v", lineNum, err)
			}

			stmt := &ReturnStmt{Expr: expr}

			// 添加到当前块或函数体
			if h.currentBlock != nil {
				switch b := (*h.currentBlock).(type) {
				case *IfStmt:
					if len(b.ElseBranch) > 0 {
						b.ElseBranch = append(b.ElseBranch, stmt)
					} else {
						b.ThenBranch = append(b.ThenBranch, stmt)
					}
				case *WhileStmt:
					b.Body = append(b.Body, stmt)
				}
			} else {
				currentFuncStatements = append(currentFuncStatements, stmt)
			}
			continue
		}

		// 处理end语句
		if line == "end" {
			if len(h.funcStack) > 0 {
				fn := h.funcStack[len(h.funcStack)-1]
				fn.Statements = currentFuncStatements

				h.funcStack = h.funcStack[:len(h.funcStack)-1]

				if len(h.funcStack) > 0 {
					currentFunc = h.funcStack[len(h.funcStack)-1].Name
				} else {
					currentFunc = ""
				}
				currentFuncStatements = []Statement{}
				h.currentBlock = nil
			}
			continue
		}

		// 处理函数体内的语句
		if currentFunc != "" {
			stmt, err := parseStatement(line, lineNum)
			if err != nil {
				return fmt.Errorf("行 %d: %v", lineNum, err)
			}

			// 如果语句是控制结构，添加到块堆栈
			switch s := stmt.(type) {
			case *IfStmt:
				h.blockStack = append(h.blockStack, s)
				h.currentBlock = &stmt
				h.inElseBranch = false // 重置分支状态
			case *WhileStmt:
				h.blockStack = append(h.blockStack, s)
				h.currentBlock = &stmt
				h.inElseBranch = false // 重置分支状态
			default:
				// 添加到当前块或函数体
				if h.currentBlock != nil {
					switch b := (*h.currentBlock).(type) {
					case *IfStmt:
						// 根据是否在 else 分支中添加
						if h.inElseBranch {
							b.ElseBranch = append(b.ElseBranch, stmt)
						} else {
							b.ThenBranch = append(b.ThenBranch, stmt)
						}
					case *WhileStmt:
						b.Body = append(b.Body, stmt)
					}
				} else {
					currentFuncStatements = append(currentFuncStatements, stmt)
				}
			}
		}
	}

	if len(h.funcStack) > 0 {
		return fmt.Errorf("函数 %s 未结束", h.funcStack[0].Name)
	}

	return nil
}

// 执行HerCode程序
func (h *HerCodeInterpreter) Execute() ([]Value, []error) {
	// 检查入口函数是否存在
	startFunc, exists := h.GlobalCtx.GetFunc("start")
	if !exists {
		return nil, []error{fmt.Errorf("入口函数 start 未定义")}
	}

	//for _, fn := range h.GlobalCtx.Functions {
	//	fmt.Printf("  函数名: %s\n", fn.Name)
	//	fmt.Printf("  参数: %v\n", fn.Parameters)
	//	fmt.Printf("  返回值: %v\n", fn.ReturnType)
	//	fmt.Println("  函数体:")
	//	for i, stmt := range fn.Statements {
	//		fmt.Printf("    %d: %s\n", i+1, stmt)
	//	}
	//	fmt.Println()
	//}
	var errs []error
	var vals []Value
	for _, stmt := range startFunc.Statements {
		val, err := stmt.Execute(h.GlobalCtx)
		vals = append(vals, val)
		errs = append(errs, err)

	}
	// 执行入口函数
	//_, err := startFunc.Statements[0].Execute(h.GlobalCtx)
	//return err
	return vals, errs
}

// 打印解析的函数信息
func (h *HerCodeInterpreter) PrintFunctions() {
	fmt.Println("解析到的函数:")
	for name, fn := range h.GlobalCtx.Functions {
		fmt.Printf("  函数名: %s\n", name)
		fmt.Printf("  参数: %v\n", fn.Parameters)
		fmt.Printf("  返回值: %v\n", fn.ReturnType)
		fmt.Println("  函数体:")
		for i, stmt := range fn.Statements {
			fmt.Printf("    %d: %s\n", i+1, stmt)
		}
		fmt.Println()
	}
}
