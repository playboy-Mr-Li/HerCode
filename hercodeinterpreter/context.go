package hercodeinterpreter

// 上下文环境
type Context struct {
	Variables map[string]Value
	Functions map[string]*HerCodeFunction
	Parent    *Context
}

func NewContext(parent *Context) *Context {
	return &Context{
		Variables: make(map[string]Value),
		Functions: make(map[string]*HerCodeFunction), // 确保这里初始化了 Functions 字段
		Parent:    parent,
	}
}

func (c *Context) GetVar(name string) (Value, bool) {
	val, ok := c.Variables[name]
	if !ok && c.Parent != nil {
		return c.Parent.GetVar(name)
	}
	return val, ok
}

func (c *Context) SetVar(name string, value Value) {
	c.Variables[name] = value
}
func (c *Context) GetFunc(name string) (*HerCodeFunction, bool) {
	fn, ok := c.Functions[name]
	if !ok && c.Parent != nil {
		return c.Parent.GetFunc(name)
	}
	return fn, ok
}

func (c *Context) SetFunc(name string, fn *HerCodeFunction) {
	// 确保 Functions 映射存在
	if c.Functions == nil {
		c.Functions = make(map[string]*HerCodeFunction)
	}
	c.Functions[name] = fn
}

// ==================== 解释器实现 ====================

func (c *Context) GlobalFunc(name string) (*HerCodeFunction, bool) {
	// 在实际实现中，这里应该访问全局函数表
	// 简化处理：假设函数在全局上下文中
	if c.Parent != nil {
		return c.Parent.GlobalFunc(name)
	}
	return nil, false
}
