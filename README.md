# HerCode 解释器

HerCode 是一个专为女性编程学习者设计的轻量级解释型编程语言，具有简洁的语法和友好的错误提示。本项目是一个用 Go 语言实现的 HerCode 解释器，支持运行 HerCode 脚本文件。

## 功能特性

- **简洁的语法**：专为编程初学者设计，语法直观易学
- **多种数据类型**：支持数字、字符串、布尔值等基本数据类型
- **控制结构**：支持 if/else 条件判断和 while 循环
- **函数支持**：支持函数定义和调用，包括递归调用
- **内置函数**：提供 len、substr、sqrt 等实用内置函数
- **友好的错误提示**：详细的错误信息和行号定位
- **跨平台**：基于 Go 语言实现，可在 Windows、macOS 和 Linux 上运行

## 语法示例

### 基本结构
```hercode
函数定义
function 函数名 参数1 参数2:
# 函数体
end

程序入口
start:
# 主程序代码
end
```

### 变量声明与赋值
```hercode
var 变量名 = 值
变量名 = 新值
```

### 控制结构
```hercode
if 语句
if 条件:
# 条件成立时执行的代码
else:
# 条件不成立时执行的代码
endif

while 循环
while 条件:
# 循环体
endwhile
```

### 函数调用
```hercode
函数名(参数1, 参数2)
```

### 输入输出
```hercode
say "Hello, World!" # 输出内容
```

## 安装与运行

### 前提条件
安装 Go 1.16 或更高版本

### 安装步骤
1. 克隆仓库：
   git clone https://github.com/playboy-Mr-Li/HerCode.git
2. cd hercode-interpreter
3. 构建项目：
   go build -o hercode
4. 运行 HerCode 脚本
   - 运行示例脚本
     ./hercode -f examples/hello.hc
   - 运行自定义脚本
     ./hercode -f path/to/your/script.hc

## 示例脚本

### hello.her
```hercode
简单的问候程序
function greet name:
say "Hello, " + name + "!"
end

start:
greet("Her World")
end
```


## 内置函数

| 函数名            | 描述         | 示例                             |
|-------------------|--------------|----------------------------------|
| len(str)          | 返回字符串长度   | len("hello") → 5                 |
| substr(str, start, end) | 返回子字符串   | substr("hello", 1, 3) → "el"    |
| sqrt(num)         | 计算平方根     | sqrt(25) → 5                     |
| print(value)      | 打印值（不换行） | print("Hello")                   |

## 错误处理

解释器提供详细的错误信息，包括错误类型和发生位置：

- 解析错误: 行 10: 变量未定义: unknown_var
- 执行错误: 行 15: 除以零错误


## 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本项目
2. 创建新分支：git checkout -b feature/your-feature
3. 提交更改：git commit -am 'Add some feature'
4. 推送到分支：git push origin feature/your-feature
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请见 LICENSE 文件。

## 灵感来源

HerCode 语言设计灵感来自对女性编程学习者的研究，旨在创造一个更友好、更包容的编程学习环境。语法设计上特别考虑了以下几点：

- 减少符号使用：使用自然语言关键词代替传统编程符号
- 明确的块结构：使用 end 明确标记代码块结束
- 友好的错误提示：提供详细的错误解释和修复建议
- 鼓励性命名：如 you_can_do_this 等鼓励性函数名

## TODO

[ ] 函数返回值实现还未成功
[ ] if while 不支持嵌套

我们希望通过 HerCode 语言，让更多女性学习者能够轻松入门编程，享受编程的乐趣！
