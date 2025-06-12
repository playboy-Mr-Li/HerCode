package hercodeinterpreter

import (
	"strings"
)

// HerCode函数
type HerCodeFunction struct {
	Name       string
	Parameters []string
	Statements []Statement
	ReturnType ValueType
}

func NewHerCodeFunction(name string) *HerCodeFunction {
	return &HerCodeFunction{
		Name:       name,
		Parameters: []string{},
		Statements: []Statement{},
	}
}

func (h *HerCodeFunction) String() string {
	r := strings.Builder{}
	r.WriteString(h.Name + "(")
	for i, param := range h.Parameters {
		if i < len(h.Parameters)-1 {
			r.WriteString(param + ", ")
		} else {
			r.WriteString(param)
		}
	}
	r.WriteString(") {\n")
	for _, stmt := range h.Statements {
		r.WriteString("  ")
		r.WriteString(stmt.String() + "\n")

	}
	r.WriteString("}\n")
	return r.String()
}

func (h *HerCodeFunction) SetName(name string) {
	h.Name = name
}

func (h *HerCodeFunction) SetStatements(statements []Statement) {
	h.Statements = statements
}

func (h *HerCodeFunction) SetParameters(parameters []string) {
	h.Parameters = parameters
}
