package hercodeinterpreter

import (
	"fmt"
	"strings"
)

// cleanComment 函数用于清理字符串中的注释。
// 该函数接受一个字符串作为输入，返回清理掉注释后的字符串。
// 注释被定义为字符串中第一个不在引号内的 '#' 字符及其后面的内容。
// 引号内的 '#' 字符和引号内的内容不会被清理。
func cleanComment(s string) string {
	if len(s) == 0 {
		return ""
	}

	insideQuotes := false
	end := len(s)

	for i := 0; i < len(s); i++ {
		if s[i] == '"' {
			// 计算前面连续的反斜杠数量，用于判断是否为转义引号
			slashCount := 0
			j := i - 1
			for j >= 0 && s[j] == '\\' {
				slashCount++
				j--
			}

			// 如果是偶数个反斜杠，则认为是有效引号切换
			if slashCount%2 == 0 {
				insideQuotes = !insideQuotes
			}
		}

		if s[i] == '#' && !insideQuotes {
			end = i
			break
		}
	}

	return s[:end]
}

// 清理字符串中的引号
func cleanQuotes(s string) string {
	n := len(s)
	if n == 0 {
		return ""
	}

	start := 0
	end := n

	for i := 0; i < n; i++ {
		if s[i] == '"' {
			// Check if this quote is escaped
			if i > 0 && s[i-1] == '\\' {
				// This is an escaped quote, skip it
				continue
			}

			// Found the first unescaped quote
			start = i + 1
			i++

			// Search for the closing quote
			for i < n {
				if s[i] == '"' {
					// Check if this closing quote is escaped
					if i > 0 && s[i-1] == '\\' {
						// Escaped quote, continue
						i++
						continue
					}
					// Found a valid closing quote
					end = i
					break
				}
				i++
			}
			break
		}
	}

	return s[start:end]
}

// 辅助函数：分割参数
func splitArguments(argsStr string) []string {
	var args []string
	var current strings.Builder
	parenDepth := 0

	for _, r := range argsStr {
		switch r {
		case '(':
			parenDepth++
			current.WriteRune(r)
		case ')':
			parenDepth--
			current.WriteRune(r)
		case ',':
			if parenDepth == 0 {
				args = append(args, current.String())
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

// 辅助函数：按运算符分割表达式
func splitByOperator(exprStr, op string) []string {
	var parts []string
	var current strings.Builder
	parenDepth := 0

	for i := 0; i < len(exprStr); i++ {
		c := exprStr[i]
		switch c {
		case '(':
			parenDepth++
			current.WriteByte(c)
		case ')':
			parenDepth--
			current.WriteByte(c)
		default:
			if parenDepth == 0 && strings.HasPrefix(exprStr[i:], op) {
				parts = append(parts, current.String())
				current.Reset()
				i += len(op) - 1
			} else {
				current.WriteByte(c)
			}
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// 解析函数定义
func parseFunctionDefinition(line string, linenumber int) (string, []string, error) {
	// 移除注释
	line = cleanComment(line)
	line = strings.TrimSpace(line)

	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return "", nil, fmt.Errorf("%d 行，函数定义缺少冒号", linenumber)
	}
	line = strings.TrimSuffix(line, ":")

	// 检查是否是函数定义
	if !strings.HasPrefix(line, "function ") {
		return "", nil, fmt.Errorf("%d 行，不是函数定义", linenumber)
	}

	// 移除 function 关键字
	line = strings.TrimSpace(strings.TrimPrefix(line, "function"))

	// 分割函数名和参数
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", nil, fmt.Errorf("%d 行，函数定义格式错误", linenumber)
	}

	// 函数名是第一个部分
	funcName := parts[0]

	// 参数是剩余部分（直到冒号）
	var params []string
	for i := 1; i < len(parts); i++ {
		if strings.HasSuffix(parts[i], ":") {
			// 处理结尾的冒号
			param := strings.TrimSuffix(parts[i], ":")
			if param != "" {
				params = append(params, param)
			}
			break
		}
		params = append(params, parts[i])
	}

	return funcName, params, nil
}
