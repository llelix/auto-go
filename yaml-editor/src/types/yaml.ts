// YAML任务配置类型定义
export interface TaskAction {
  type: 'wait_appear' | 'wait_disappear' | 'fill' | 'click' | 'select' | 'hover' | 'drag_drop' | 'get_text' | 'get_attribute';
  selector?: string;
  target?: string; // 用于drag_drop
  attribute?: string; // 用于get_attribute
  value?: string | number;
  timeout?: number;
  output_key?: string;
  error_message: string;
}

export interface Task {
  name: string;
  url: string;
  wait_time?: number;
  screenshot?: boolean;
  actions: TaskAction[];
}

export interface YamlConfig {
  tasks: Task[];
}

// Canvas节点类型
export interface CanvasNode {
  id: string;
  type: 'task' | 'action';
  position: { x: number; y: number };
  data: {
    label: string;
    task?: Task;
    action?: TaskAction;
    icon?: string;
  };
}

// 连接线类型
export interface CanvasEdge {
  id: string;
  source: string;
  target: string;
  type?: 'default';
}

// 组件库项目类型
export interface ComponentLibraryItem {
  id: string;
  name: string;
  description: string;
  icon: string;
  category: 'basic' | 'interaction' | 'verification' | 'extraction';
  template: TaskAction;
}

// 编辑器状态类型
export interface EditorState {
  tasks: Task[];
  nodes: CanvasNode[];
  edges: CanvasEdge[];
  selectedNodeId?: string;
  yamlOutput: string;
  isDirty: boolean;
}