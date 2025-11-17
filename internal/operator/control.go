package operator

import (
	"encoding/json"
	"fmt"
)

// ControlType 定义流程控制类型
const (
	ControlTypeForLoop     = "for"
	ControlTypeIfCondition = "if"
	ControlTypeElseCondition = "else"
	ControlTypeWhileLoop   = "while"
)

// ControlFlowType 定义控制流类型
const (
	ControlFlowBreak    = "break"
	ControlFlowContinue = "continue"
)

// ControlNode 流程控制节点基类
type ControlNode struct {
	Type     string      `json:"type"`     // 控制类型："for", "if", "else", "while"
	Children []NodeItem  `json:"children"` // 子节点，可以是Action或ControlNode
	
	// 循环参数
	Variable  string `json:"variable,omitempty"`  // 循环变量名
	From      int    `json:"from,omitempty"`       // 起始值
	To        int    `json:"to,omitempty"`         // 结束值
	Step      int    `json:"step,omitempty"`       // 步长，默认1
	Condition string `json:"condition,omitempty"`  // 条件表达式
}

// ForLoop 定义for循环结构（已弃用，参数直接集成到ControlNode中）
type ForLoop struct {
	ControlNode
}

// WhileLoop 定义while循环结构（已弃用，只保留for循环）
type WhileLoop struct {
	ControlNode
	Condition string `json:"condition"` // 循环条件表达式
	MaxIter   int    `json:"max_iter"`  // 最大迭代次数，防止无限循环
}

// IfCondition 定义条件分支结构
type IfCondition struct {
	ControlNode
	Condition string `json:"condition"` // 条件表达式
}

// NodeItem 定义节点项，可以是Action或ControlNode
type NodeItem struct {
	Action      *Action      `json:"action,omitempty" yaml:"action,omitempty"`       // 基本操作
	ControlNode *ControlNode `json:"control_node,omitempty" yaml:"control_node,omitempty"` // 控制节点
}

// ExecutionContext 执行上下文，用于存储变量和控制流状态
type ExecutionContext struct {
	Variables    map[string]interface{} `json:"variables"`     // 变量表
	ControlFlow  *ControlFlow           `json:"control_flow"`  // 控制流状态
	OutputValues map[string]string      `json:"output_values"` // 输出值存储
}

// ControlFlow 控制流状态
type ControlFlow struct {
	BreakSignal    bool   `json:"break_signal"`    // break信号
	ContinueSignal bool   `json:"continue_signal"` // continue信号
	CurrentLoop    string `json:"current_loop"`    // 当前循环标识
}

// IsValid 验证控制节点是否有效
func (cn *ControlNode) IsValid() error {
	switch cn.Type {
	case ControlTypeForLoop:
		return cn.validateForLoop()
	case ControlTypeIfCondition:
		return cn.validateIfCondition()
	case ControlTypeElseCondition:
		return cn.validateElseCondition()
	case ControlTypeWhileLoop:
		return cn.validateWhileLoop()
	default:
		return fmt.Errorf("不支持的流程控制类型: %s", cn.Type)
	}
}

// validateForLoop 验证for循环配置
func (cn *ControlNode) validateForLoop() error {
	// 这里可以添加具体的验证逻辑
	return nil
}

// validateIfCondition 验证条件分支配置
func (cn *ControlNode) validateIfCondition() error {
	// 这里可以添加具体的验证逻辑
	return nil
}

// validateWhileLoop 验证while循环配置
func (cn *ControlNode) validateWhileLoop() error {
	// 这里可以添加具体的验证逻辑
	return nil
}

// validateElseCondition 验证else分支配置
func (cn *ControlNode) validateElseCondition() error {
	// 验证else节点必须有子节点
	if len(cn.Children) == 0 {
		return fmt.Errorf("else分支必须包含至少一个子节点")
	}
	// 在实际执行时，else节点需要检查前一个节点是否为if节点
	// 这个验证在解码时无法完成，需要在执行时进行
	return nil
}

// NewExecutionContext 创建新的执行上下文
func NewExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		Variables:    make(map[string]interface{}),
		OutputValues: make(map[string]string),
		ControlFlow: &ControlFlow{
			BreakSignal:    false,
			ContinueSignal: false,
			CurrentLoop:    "",
		},
	}
}

// SetVariable 设置变量值
func (ec *ExecutionContext) SetVariable(name string, value interface{}) {
	ec.Variables[name] = value
}

// GetVariable 获取变量值
func (ec *ExecutionContext) GetVariable(name string) interface{} {
	if val, exists := ec.Variables[name]; exists {
		return val
	}
	return nil
}

// SignalBreak 发送break信号
func (ec *ExecutionContext) SignalBreak() {
	ec.ControlFlow.BreakSignal = true
}

// SignalContinue 发送continue信号
func (ec *ExecutionContext) SignalContinue() {
	ec.ControlFlow.ContinueSignal = true
}

// ResetControlFlow 重置控制流状态
func (ec *ExecutionContext) ResetControlFlow() {
	ec.ControlFlow.BreakSignal = false
	ec.ControlFlow.ContinueSignal = false
}

// MarshalJSON NodeItem的自定义JSON序列化
func (ni *NodeItem) MarshalJSON() ([]byte, error) {
	if ni.Action != nil {
		return json.Marshal(ni.Action)
	}
	if ni.ControlNode != nil {
		return json.Marshal(ni.ControlNode)
	}
	return []byte("null"), nil
}

// UnmarshalJSON NodeItem的自定义JSON反序列化
func (ni *NodeItem) UnmarshalJSON(data []byte) error {
	// 首先尝试解析为ControlNode（因为某些类型如"if", "else"需要优先作为控制节点处理）
	var controlNode ControlNode
	if err := json.Unmarshal(data, &controlNode); err == nil && controlNode.Type != "" {
		// 检查是否为支持的控制节点类型
		if controlNode.Type == ControlTypeForLoop || controlNode.Type == ControlTypeIfCondition || controlNode.Type == ControlTypeElseCondition || controlNode.Type == ControlTypeWhileLoop {
			ni.ControlNode = &controlNode
			return nil
		}
	}

	// 尝试解析为Action
	var action Action
	if err := json.Unmarshal(data, &action); err == nil && action.Type != "" {
		ni.Action = &action
		return nil
	}

	return fmt.Errorf("无法解析节点项，既不是有效的Action也不是ControlNode")
}

// UnmarshalYAML NodeItem的自定义YAML反序列化
func (ni *NodeItem) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// 首先尝试解析为ControlNode
	var controlNode ControlNode
	if err := unmarshal(&controlNode); err == nil && controlNode.Type != "" {
		// 检查是否为支持的控制节点类型
		if controlNode.Type == ControlTypeForLoop || controlNode.Type == ControlTypeIfCondition || controlNode.Type == ControlTypeElseCondition || controlNode.Type == ControlTypeWhileLoop {
			ni.ControlNode = &controlNode
			return nil
		}
	}

	// 尝试解析为Action
	var action Action
	if err := unmarshal(&action); err == nil && action.Type != "" {
		ni.Action = &action
		return nil
	}

	return fmt.Errorf("无法解析节点项，既不是有效的Action也不是ControlNode")
}

// IsAction 检查是否为Action节点
func (ni *NodeItem) IsAction() bool {
	return ni.Action != nil
}

// IsControlNode 检查是否为ControlNode节点
func (ni *NodeItem) IsControlNode() bool {
	return ni.ControlNode != nil
}