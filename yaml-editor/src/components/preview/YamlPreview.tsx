'use client';

import { useState } from 'react';
import { Copy, Check, FileText } from 'lucide-react';
import { DroppedItem } from '../canvas/DroppableCanvas';

interface YamlPreviewProps {
  items: DroppedItem[];
}

export function YamlPreview({ items }: YamlPreviewProps) {
  const [copied, setCopied] = useState(false);

  // 生成YAML内容
  const generateYaml = () => {
    if (items.length === 0) {
      return `# AutoGo 任务配置文件
# 从左侧拖拽组件开始创建配置

- name: "示例任务"
  url: "https://example.com"
  wait_time: 3
  screenshot: true
  actions:
    - type: "wait_appear"
      selector: "#element"
      timeout: 5
      error_message: "等待元素出现失败"`;
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

  const yamlContent = generateYaml();

  // 复制到剪贴板
  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(yamlContent);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error('复制失败:', error);
    }
  };

  return (
    <div className="h-full flex flex-col">
      {/* 标题栏 */}
      <div className="p-4 border-b border-slate-700 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <FileText className="w-5 h-5 text-green-400" />
          <h2 className="text-white font-semibold">YAML 预览</h2>
        </div>
        <button
          onClick={handleCopy}
          className="flex items-center gap-2 px-3 py-1.5 bg-slate-700 hover:bg-slate-600 text-white rounded-md transition-colors text-sm"
          title="复制到剪贴板"
        >
          {copied ? (
            <>
              <Check className="w-4 h-4" />
              已复制
            </>
          ) : (
            <>
              <Copy className="w-4 h-4" />
              复制
            </>
          )}
        </button>
      </div>

      {/* 代码内容 */}
      <div className="flex-1 overflow-hidden">
        <pre className="h-full overflow-auto p-4 text-sm font-mono text-slate-300 bg-slate-950">
          <code>{yamlContent}</code>
        </pre>
      </div>

      {/* 统计信息 */}
      {items.length > 0 && (
        <div className="p-3 border-t border-slate-700 bg-slate-800/50">
          <div className="text-xs text-slate-400">
            <div>组件数量: {items.length}</div>
          </div>
        </div>
      )}
    </div>
  );
}