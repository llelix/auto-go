'use client';

import { useState } from 'react';
import { ComponentSidebar } from './sidebar/ComponentSidebar';
import { DroppableCanvas } from './canvas/DroppableCanvas';
import { YamlPreview } from './preview/YamlPreview';
import { Toolbar } from './ui/Toolbar';

interface DroppedItem {
  id: string;
  type: string;
  title: string;
  description: string;
  x: number;
  y: number;
}

export function YamlEditor() {
  const [droppedItems, setDroppedItems] = useState<DroppedItem[]>([]);

  // 解析YAML文件
  const parseYamlFile = (content: string): DroppedItem[] => {
    try {
      // 简单的YAML解析器，提取action项
      const lines = content.split('\n');
      const items: DroppedItem[] = [];
      let index = 0;

      lines.forEach((line) => {
        if (line.trim().startsWith('- type:')) {
          const type = line.trim().replace('- type:', '').trim().replace(/"/g, '');
          
          // 根据type创建对应的组件
          const componentMap = {
            'wait_appear': { title: '等待出现', description: '等待元素出现在页面上', icon: '⏳' },
            'click': { title: '点击', description: '点击页面元素', icon: '👆' },
            'fill': { title: '填写', description: '在输入框中填写内容', icon: '✏️' }
          };

          const config = componentMap[type as keyof typeof componentMap];
          if (config) {
            items.push({
              id: `${type}-imported-${index++}`,
              type,
              title: config.title,
              description: config.description,
              x: 100 + (items.length % 3) * 220,
              y: 50 + Math.floor(items.length / 3) * 120
            });
          }
        }
      });

      return items;
    } catch (error) {
      console.error('YAML解析失败:', error);
      return [];
    }
  };

  // 处理文件导入
  const handleImport = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      if (content) {
        const parsedItems = parseYamlFile(content);
        if (parsedItems.length > 0) {
          setDroppedItems(parsedItems);
          alert(`成功导入 ${parsedItems.length} 个组件`);
        } else {
          alert('未能从文件中识别到有效的组件配置');
        }
      }
    };
    reader.onerror = () => {
      alert('文件读取失败');
    };
    reader.readAsText(file);
  };

  // 处理文件导出
  const handleExport = () => {
    const yamlContent = generateYamlFromItems(droppedItems);
    const blob = new Blob([yamlContent], { type: 'text/yaml' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'autogo-config.yaml';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  // 根据组件生成YAML
  const generateYamlFromItems = (items: DroppedItem[]): string => {
    if (items.length === 0) {
      return `# AutoGo 任务配置文件
# 空配置

- name: "示例任务"
  url: "https://example.com"
  wait_time: 3
  screenshot: true
  actions: []`;
    }

    const yamlActions = items.map(item => {
      let actionConfig = `    - type: "${item.type}"`;
      
      switch (item.type) {
        case 'wait_appear':
          actionConfig += `
      selector: "#element"
      timeout: 5
      error_message: "等待元素出现失败"`;
          break;
        case 'click':
          actionConfig += `
      selector: "#button"
      wait_before: 1`;
          break;
        case 'fill':
          actionConfig += `
      selector: "#input"
      value: "示例文本"`;
          break;
        default:
          actionConfig += `
      selector: "#element"`;
      }
      
      return actionConfig;
    }).join('\n');

    return `# AutoGo 任务配置文件
# 通过拖拽组件创建的任务配置

- name: "自动化任务"
  url: "https://example.com"
  wait_time: 3
  screenshot: true
  actions:
${yamlActions}`;
  };

  return (
    <div className="h-screen w-full bg-slate-900 flex flex-col overflow-hidden">
      {/* 顶部工具栏 */}
      <Toolbar onImport={handleImport} onExport={handleExport} />

      {/* 主编辑区域 */}
      <div className="flex-1 flex overflow-hidden">
        {/* 左侧组件库 */}
        <div className="w-80 bg-slate-800 border-r border-slate-700 flex-shrink-0">
          <ComponentSidebar />
        </div>

        {/* 中间画布区域 */}
        <div className="flex-1 bg-slate-950 relative overflow-hidden">
          <DroppableCanvas 
            items={droppedItems} 
            onItemsChange={setDroppedItems}
          />
        </div>

        {/* 右侧预览面板 */}
        <div className="w-96 bg-slate-800 border-l border-slate-700 flex-shrink-0">
          <YamlPreview items={droppedItems} />
        </div>
      </div>
    </div>
  );
}