package operator

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// ControlExecutor 流程控制执行器
type ControlExecutor struct {
	TaskManager    *TaskManager
	Context        *ExecutionContext
	LoopStack      []string                  // 循环栈，用于嵌套循环管理
	ScopeVariables map[string]map[string]any // 嵌套作用域变量存储
}

// NewControlExecutor 创建新的控制执行器
func NewControlExecutor(tm *TaskManager) *ControlExecutor {
	return &ControlExecutor{
		TaskManager:    tm,
		Context:        NewExecutionContext(),
		LoopStack:      make([]string, 0),
		ScopeVariables: make(map[string]map[string]any),
	}
}

// ExecuteNodeItems 执行节点项序列
func (ce *ControlExecutor) ExecuteNodeItems(items []NodeItem) error {
	for _, item := range items {
		// 检查控制流信号
		if ce.Context.ControlFlow.BreakSignal {
			ce.Context.ResetControlFlow()
			break
		}

		if ce.Context.ControlFlow.ContinueSignal {
			ce.Context.ResetControlFlow()
			continue
		}

		if err := ce.ExecuteNodeItem(item); err != nil {
			return err
		}

		// 操作间添加短暂延迟，提高执行稳定性
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

// ExecuteNodeItem 执行单个节点项
func (ce *ControlExecutor) ExecuteNodeItem(item NodeItem) error {
	if item.IsAction() {
		return ce.executeAction(item.Action)
	}

	if item.IsControlNode() {
		return ce.executeControlNode(item.ControlNode)
	}

	return fmt.Errorf("无效的节点项，既不是Action也不是ControlNode")
}

// replaceVariables 替换字符串中的模板变量
func (ce *ControlExecutor) replaceVariables(input string) string {
	result := input
	
	// 替换 {{variable}} 格式的变量
	for key, value := range ce.Context.Variables {
		placeholder := "{{" + key + "}}"
		var stringValue string
		
		switch v := value.(type) {
		case string:
			stringValue = v
		case int:
			stringValue = strconv.Itoa(v)
		case float64:
			stringValue = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			stringValue = strconv.FormatBool(v)
		default:
			stringValue = fmt.Sprintf("%v", v)
		}
		
		result = strings.ReplaceAll(result, placeholder, stringValue)
	}
	
	return result
}

// executeAction 执行单个动作
func (ce *ControlExecutor) executeAction(action *Action) error {
	// 检查是否为控制流操作
	if action.Type == "break" {
		ce.HandleBreak()
		return nil
	}

	if action.Type == "continue" {
		ce.HandleContinue()
		return nil
	}

	// 替换模板变量
	selector := ce.replaceVariables(action.Selector)
	target := ce.replaceVariables(action.Target)
	value := ce.replaceVariables(action.Value)
	errorMessage := ce.replaceVariables(action.ErrorMessage)

	// 执行常规动作
	var err error
	switch action.Type {
	case ActionClick:
		err = ce.TaskManager.BrowserManager.Click(selector)

	case ActionFill:
		if value == "" {
			err = fmt.Errorf("fill操作需要提供value参数")
		} else {
			err = ce.TaskManager.BrowserManager.FillForm(map[string]string{selector: value})
		}

	case ActionHover:
		err = ce.TaskManager.BrowserManager.Hover(selector)

	case ActionSelect:
		if value == "" {
			err = fmt.Errorf("select操作需要提供value参数")
		} else {
			err = ce.TaskManager.BrowserManager.SelectOption(selector, value)
		}

	case ActionScroll:
		err = ce.TaskManager.BrowserManager.ScrollToElement(selector)

	case ActionRightClick:
		err = ce.TaskManager.BrowserManager.RightClick(selector)

	case ActionDragDrop:
		if target == "" {
			err = fmt.Errorf("drag_drop操作需要提供target参数")
		} else {
			err = ce.TaskManager.BrowserManager.DragAndDrop(selector, target)
		}

	case ActionWaitAppear:
		timeout := time.Duration(10) * time.Second
		if action.Timeout > 0 {
			timeout = time.Duration(action.Timeout) * time.Second
		}
		err = ce.TaskManager.BrowserManager.WaitForSelector(selector, timeout)

	case ActionWaitDisappear:
		timeout := time.Duration(10) * time.Second
		if action.Timeout > 0 {
			timeout = time.Duration(action.Timeout) * time.Second
		}
		err = ce.TaskManager.BrowserManager.WaitForElementDisappear(selector, timeout)

	case ActionGetText:
		text, getTextErr := ce.TaskManager.BrowserManager.GetText(selector)
		if getTextErr != nil {
			err = getTextErr
		} else {
			log.Printf("📝 获取元素文本: %s = '%s'", selector, text)
			if action.OutputKey != "" {
				ce.Context.SetVariable(action.OutputKey, text)
				log.Printf("📋 文本已存储到变量: %s", action.OutputKey)
			}
		}

	case ActionGetAttribute:
		if action.Attribute == "" {
			err = fmt.Errorf("get_attribute操作需要提供attribute参数")
		} else {
			attr, getAttrErr := ce.TaskManager.BrowserManager.GetAttribute(selector, action.Attribute)
			if getAttrErr != nil {
				err = getAttrErr
			} else {
				log.Printf("🏷️ 获取元素属性: %s.%s = '%s'", selector, action.Attribute, attr)
				if action.OutputKey != "" {
					ce.Context.SetVariable(action.OutputKey, attr)
					log.Printf("📋 属性值已存储到变量: %s", action.OutputKey)
				}
			}
		}

	default:
		err = fmt.Errorf("不支持的操作类型: %s", action.Type)
	}

	if err != nil {
		if errorMessage != "" {
			return fmt.Errorf("操作失败: %s", errorMessage)
		}
		return fmt.Errorf("操作失败: %s - %v", action.Type, err)
	}

	return nil
}

// executeControlNode 执行控制节点
func (ce *ControlExecutor) executeControlNode(node *ControlNode) error {
	// 验证控制节点
	if err := node.IsValid(); err != nil {
		return err
	}

	switch node.Type {
	case ControlTypeForLoop:
		return ce.executeForLoop(node)
	case ControlTypeIfCondition:
		return ce.executeIfCondition(node)
	case ControlTypeElseCondition:
		return ce.executeElseCondition(node)
	default:
		return fmt.Errorf("不支持的控制节点类型: %s", node.Type)
	}
}

// executeForLoop 执行for循环（使用单层循环结构）
func (ce *ControlExecutor) executeForLoop(node *ControlNode) error {
	log.Printf("🔄 开始执行for循环")

	// 解析循环参数
	var loopVar string
	var start, end int

	// 从node中解析参数
	loopVar = node.Variable
	if loopVar == "" {
		loopVar = "i"
	}
	
	start = node.From
	end = node.To

	log.Printf("🔄 循环参数: 变量=%s, 起始=%d, 结束=%d", loopVar, start, end)

	// 使用真正的单层循环结构，避免嵌套
	currentIndex := start
	
	for currentIndex <= end {
		// 检查控制流信号
		if ce.Context.ControlFlow.BreakSignal {
			ce.Context.ResetControlFlow()
			break
		}

		// 设置循环变量
		ce.Context.SetVariable(loopVar, currentIndex)
		log.Printf("🔄 循环迭代: %s = %d", loopVar, currentIndex)

		// 执行子节点序列
		for childIndex := 0; childIndex < len(node.Children); {
			// 检查控制流信号
			if ce.Context.ControlFlow.BreakSignal {
				ce.Context.ResetControlFlow()
				break
			}
			
			if ce.Context.ControlFlow.ContinueSignal {
				ce.Context.ResetControlFlow()
				continue
			}
			
			if err := ce.ExecuteNodeItem(node.Children[childIndex]); err != nil {
				return err
			}
			
			childIndex++
			
			// 操作间添加短暂延迟
			time.Sleep(500 * time.Millisecond)
		}

		// 检查continue信号
		if ce.Context.ControlFlow.ContinueSignal {
			ce.Context.ResetControlFlow()
			currentIndex++
			continue
		}

		currentIndex++
	}

	log.Printf("🔄 for循环执行完成")
	return nil
}

// executeIfCondition 执行条件分支（单层循环结构）
func (ce *ControlExecutor) executeIfCondition(node *ControlNode) error {
	log.Printf("❓ 开始执行条件判断")

	// 如果有子节点，直接执行（实际应用中应该评估条件表达式）
	// 这里简化处理，执行if分支
	if len(node.Children) > 0 {
		log.Printf("✅ 条件为真，执行if分支")
		// 使用单层循环结构执行子节点
		for idx := 0; idx < len(node.Children); idx++ {
			// 检查控制流信号
			if ce.Context.ControlFlow.BreakSignal {
				ce.Context.ResetControlFlow()
				break
			}
			
			if ce.Context.ControlFlow.ContinueSignal {
				ce.Context.ResetControlFlow()
				continue
			}
			
			if err := ce.ExecuteNodeItem(node.Children[idx]); err != nil {
				return err
			}
			
			// 操作间添加短暂延迟
			time.Sleep(500 * time.Millisecond)
		}
	} else {
		log.Printf("❌ 条件为假或无子节点，跳过if分支")
	}

	log.Printf("❓ 条件判断执行完成")
	return nil
}

// executeElseCondition 执行else分支（单层循环结构）
func (ce *ControlExecutor) executeElseCondition(node *ControlNode) error {
	log.Printf("❓ 开始执行else分支")

	// 检查else分支是否有对应的if前置节点
	// 这个验证应该在执行时进行，因为解码时无法确定执行顺序
	
	// 执行else分支的子节点
	if len(node.Children) > 0 {
		log.Printf("✅ 执行else分支")
		// 使用单层循环结构执行子节点
		for idx := 0; idx < len(node.Children); idx++ {
			// 检查控制流信号
			if ce.Context.ControlFlow.BreakSignal {
				ce.Context.ResetControlFlow()
				break
			}
			
			if ce.Context.ControlFlow.ContinueSignal {
				ce.Context.ResetControlFlow()
				continue
			}
			
			if err := ce.ExecuteNodeItem(node.Children[idx]); err != nil {
				return err
			}
			
			// 操作间添加短暂延迟
			time.Sleep(500 * time.Millisecond)
		}
	} else {
		log.Printf("❌ else分支无子节点，跳过")
	}

	log.Printf("❓ else分支执行完成")
	return nil
}

// EvaluateCondition 评估条件表达式
func (ce *ControlExecutor) EvaluateCondition(conditionExpr string) (bool, error) {
	if conditionExpr == "" {
		return true, nil
	}

	result, err := EvaluateBoolean(conditionExpr, ce.Context)
	if err != nil {
		return false, fmt.Errorf("条件表达式评估失败: %w", err)
	}

	log.Printf("❓ 条件表达式 '%s' 评估结果: %v", conditionExpr, result)
	return result, nil
}

// SetVariable 设置变量
func (ce *ControlExecutor) SetVariable(name string, value interface{}) {
	ce.Context.SetVariable(name, value)
}

// GetVariable 获取变量
func (ce *ControlExecutor) GetVariable(name string) interface{} {
	return ce.Context.GetVariable(name)
}

// ResetContext 重置执行上下文
func (ce *ControlExecutor) ResetContext() {
	ce.Context = NewExecutionContext()
}

// PrintVariables 打印所有变量
func (ce *ControlExecutor) PrintVariables() {
	log.Printf("📊 当前变量状态:")
	for name, value := range ce.Context.Variables {
		log.Printf("  %s = %v", name, value)
	}
}

// 嵌套循环管理方法

// PushLoop 推入循环到栈
func (ce *ControlExecutor) PushLoop(loopID string) {
	ce.LoopStack = append(ce.LoopStack, loopID)
	ce.Context.ControlFlow.CurrentLoop = loopID
	// 创建新的作用域
	ce.ScopeVariables[loopID] = make(map[string]any)
}

// PopLoop 从栈中弹出循环
func (ce *ControlExecutor) PopLoop() {
	if len(ce.LoopStack) > 0 {
		// 移除当前作用域
		lastLoop := ce.LoopStack[len(ce.LoopStack)-1]
		delete(ce.ScopeVariables, lastLoop)

		ce.LoopStack = ce.LoopStack[:len(ce.LoopStack)-1]

		// 更新当前循环
		if len(ce.LoopStack) > 0 {
			ce.Context.ControlFlow.CurrentLoop = ce.LoopStack[len(ce.LoopStack)-1]
		} else {
			ce.Context.ControlFlow.CurrentLoop = ""
		}
	}
}

// GetCurrentLoop 获取当前循环ID
func (ce *ControlExecutor) GetCurrentLoop() string {
	if len(ce.LoopStack) > 0 {
		return ce.LoopStack[len(ce.LoopStack)-1]
	}
	return ""
}

// IsInLoop 检查是否在循环中
func (ce *ControlExecutor) IsInLoop() bool {
	return len(ce.LoopStack) > 0
}

// SetVariableInScope 在作用域中设置变量
func (ce *ControlExecutor) SetVariableInScope(variableName string, value any) {
	if ce.IsInLoop() {
		currentLoop := ce.GetCurrentLoop()
		if scope, exists := ce.ScopeVariables[currentLoop]; exists {
			scope[variableName] = value
		}
	}
	// 同时设置到全局上下文
	ce.Context.SetVariable(variableName, value)
}

// GetVariableFromScope 从作用域获取变量
func (ce *ControlExecutor) GetVariableFromScope(variableName string) any {
	// 首先检查当前作用域
	if ce.IsInLoop() {
		currentLoop := ce.GetCurrentLoop()
		if scope, exists := ce.ScopeVariables[currentLoop]; exists {
			if val, ok := scope[variableName]; ok {
				return val
			}
		}
	}
	// 回退到全局上下文
	return ce.Context.GetVariable(variableName)
}

// HandleBreak 处理break信号
func (ce *ControlExecutor) HandleBreak() {
	if ce.IsInLoop() {
		ce.Context.SignalBreak()
		log.Printf("🛑 发送break信号，退出当前循环: %s", ce.GetCurrentLoop())
	} else {
		log.Printf("⚠️  break信号在循环外无效")
	}
}

// HandleContinue 处理continue信号
func (ce *ControlExecutor) HandleContinue() {
	if ce.IsInLoop() {
		ce.Context.SignalContinue()
		log.Printf("⏭️  发送continue信号，跳过当前循环剩余部分: %s", ce.GetCurrentLoop())
	} else {
		log.Printf("⚠️  continue信号在循环外无效")
	}
}
