import { create } from 'zustand';
import { EditorState, Task, CanvasNode, CanvasEdge } from '@/types/yaml';
import { YamlParser } from '@/lib/yaml-parser/parser';

interface EditorStore extends EditorState {
  // Actions
  setTasks: (tasks: Task[]) => void;
  addTask: (task: Task) => void;
  updateTask: (taskId: string, task: Task) => void;
  deleteTask: (taskId: string) => void;
  
  // Canvas actions
  setNodes: (nodes: CanvasNode[]) => void;
  addNode: (node: CanvasNode) => void;
  updateNode: (nodeId: string, node: Partial<CanvasNode>) => void;
  deleteNode: (nodeId: string) => void;
  setEdges: (edges: CanvasEdge[]) => void;
  addEdge: (edge: CanvasEdge) => void;
  deleteEdge: (edgeId: string) => void;
  
  // Selection
  setSelectedNode: (nodeId?: string) => void;
  
  // File operations
  loadYamlFile: (yamlString: string) => void;
  exportYaml: () => void;
  resetEditor: () => void;
  
  // Utility
  updateYamlOutput: () => void;
  markDirty: () => void;
  markClean: () => void;
}

export const useEditorStore = create<EditorStore>((set, get) => ({
  // Initial state
  tasks: [],
  nodes: [],
  edges: [],
  selectedNodeId: undefined,
  yamlOutput: '',
  isDirty: false,

  // Task actions
  setTasks: (tasks) => {
    set({ tasks });
    get().updateYamlOutput();
    get().markDirty();
  },

  addTask: (task) => {
    const tasks = [...get().tasks, task];
    set({ tasks });
    get().updateYamlOutput();
    get().markDirty();
  },

  updateTask: (taskId, updatedTask) => {
    const tasks = get().tasks.map(task => 
      task.name === taskId ? updatedTask : task
    );
    set({ tasks });
    get().updateYamlOutput();
    get().markDirty();
  },

  deleteTask: (taskId) => {
    const tasks = get().tasks.filter(task => task.name !== taskId);
    set({ tasks });
    get().updateYamlOutput();
    get().markDirty();
  },

  // Canvas actions
  setNodes: (nodes) => {
    set({ nodes });
  },

  addNode: (node) => {
    const nodes = [...get().nodes, node];
    set({ nodes });
    get().markDirty();
  },

  updateNode: (nodeId, updatedNode) => {
    const nodes = get().nodes.map(node => 
      node.id === nodeId ? { ...node, ...updatedNode } : node
    );
    set({ nodes });
    get().markDirty();
  },

  deleteNode: (nodeId) => {
    const nodes = get().nodes.filter(node => node.id !== nodeId);
    const edges = get().edges.filter(edge => 
      edge.source !== nodeId && edge.target !== nodeId
    );
    set({ nodes, edges });
    get().markDirty();
  },

  setEdges: (edges) => {
    set({ edges });
  },

  addEdge: (edge) => {
    const edges = [...get().edges, edge];
    set({ edges });
    get().markDirty();
  },

  deleteEdge: (edgeId) => {
    const edges = get().edges.filter(edge => edge.id !== edgeId);
    set({ edges });
    get().markDirty();
  },

  // Selection
  setSelectedNode: (nodeId) => {
    set({ selectedNodeId: nodeId });
  },

  // File operations
  loadYamlFile: (yamlString) => {
    try {
      const tasks = YamlParser.parseYaml(yamlString);
      set({ 
        tasks, 
        isDirty: false,
        selectedNodeId: undefined
      });
      get().updateYamlOutput();
    } catch (error) {
      console.error('加载YAML失败:', error);
      throw error;
    }
  },

  exportYaml: () => {
    try {
      YamlParser.downloadFile(get().tasks);
      get().markClean();
    } catch (error) {
      console.error('导出YAML失败:', error);
      throw error;
    }
  },

  resetEditor: () => {
    set({
      tasks: [],
      nodes: [],
      edges: [],
      selectedNodeId: undefined,
      yamlOutput: '',
      isDirty: false
    });
  },

  // Utility
  updateYamlOutput: () => {
    try {
      const yamlOutput = YamlParser.stringifyYaml(get().tasks);
      set({ yamlOutput });
    } catch (error) {
      console.error('YAML生成失败:', error);
      set({ yamlOutput: '// YAML生成错误' });
    }
  },

  markDirty: () => {
    set({ isDirty: true });
  },

  markClean: () => {
    set({ isDirty: false });
  }
}));