package main

import (
	"fmt"
	"github.com/playboy-Mr-Li/HerCode/hercodeinterpreter"
	"github.com/playboy-Mr-Li/HerCode/itype"
	"github.com/playboy-Mr-Li/HerCode/readfile"
)

var F itype.Flag

func main() {
	itype.PaseFlag(&F)

	b, err := readfile.ReadFile(F.FileName)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	script := string(b)

	//fmt.Println(script)

	interpreter := hercodeinterpreter.NewHerCodeInterpreter()

	// 解析脚本
	fmt.Println("Her Code is Compiling...")
	if err := interpreter.Parse(script); err != nil {
		fmt.Printf("解析错误: %v\n", err)
		return
	}
	if F.Debug {
		fmt.Println("Her Code is Debugging...")
		// 打印解析结果
		interpreter.PrintFunctions()
	}

	// 执行程序
	fmt.Println("Her Code is Running...")
	vals, errs := interpreter.Execute()
	if len(errs) > 0 {
		for _, err := range errs {
			if err != nil {
				fmt.Printf("执行错误: %v\n", err)
			}
		}
	}
	if len(vals) > 0 {
		for _, val := range vals {
			if val.Error != nil {
				fmt.Printf("执行错误: %v\n", val.Error)
			}
		}

	}

}
