'use client';

import React, { useState, useEffect } from 'react';
import { Copy, Check, FileText, AlertCircle } from 'lucide-react';
import { useEditorStore } from '@/lib/store/editor-store';
import { YamlParser } from '@/lib/yaml-parser/parser';

export function YamlPreview() {
  const { yamlOutput, tasks } = useEditorStore();
  const [copied, setCopied] = useState(false);
  const [validation, setValidation] = useState<{ valid: boolean; errors: string[] } | null>(null);

  // 验证配置
  useEffect(() => {
    if (tasks.length > 0) {
      const result = YamlParser.validateConfig(tasks);
      // 使用 setTimeout 避免同步调用 setState
      const timeoutId = setTimeout(() => {
        setValidation(result);
      }, 0);
      
      return () => clearTimeout(timeoutId);
    } else {
      // 使用 setTimeout 避免同步调用 setState
      const timeoutId = setTimeout(() => {
        setValidation(null);
      }, 0);
      
      return () => clearTimeout(timeoutId);
    }
  }, [tasks]);

  // 复制到剪贴板
  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(yamlOutput);
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

      {/* 验证状态 */}
      {validation && (
        <div className={`p-3 border-b border-slate-700 ${validation.valid ? 'bg-green-900/20' : 'bg-red-900/20'}`}>
          <div className="flex items-center gap-2">
            <AlertCircle className={`w-4 h-4 ${validation.valid ? 'text-green-400' : 'text-red-400'}`} />
            <span className={`text-sm ${validation.valid ? 'text-green-400' : 'text-red-400'}`}>
              {validation.valid ? '配置验证通过' : `发现 ${validation.errors.length} 个错误`}
            </span>
          </div>
          {!validation.valid && validation.errors.length > 0 && (
            <div className="mt-2 space-y-1">
              {validation.errors.slice(0, 3).map((error, index) => (
                <div key={index} className="text-xs text-red-300 ml-6">
                  • {error}
                </div>
              ))}
              {validation.errors.length > 3 && (
                <div className="text-xs text-red-300 ml-6">
                  ... 还有 {validation.errors.length - 3} 个错误
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {/* 代码内容 */}
      <div className="flex-1 overflow-hidden">
        {yamlOutput ? (
          <pre className="h-full overflow-auto p-4 text-sm font-mono text-slate-300 bg-slate-950">
            <code>{yamlOutput}</code>
          </pre>
        ) : (
          <div className="h-full flex items-center justify-center">
            <div className="text-center">
              <FileText className="w-12 h-12 text-slate-600 mx-auto mb-3" />
              <div className="text-slate-400">暂无YAML内容</div>
              <div className="text-slate-500 text-sm mt-1">
                拖拽组件到画布开始创建配置
              </div>
            </div>
          </div>
        )}
      </div>

      {/* 统计信息 */}
      {tasks.length > 0 && (
        <div className="p-3 border-t border-slate-700 bg-slate-800/50">
          <div className="text-xs text-slate-400 space-y-1">
            <div>任务数量: {tasks.length}</div>
            <div>操作总数: {tasks.reduce((sum, task) => sum + (task.actions?.length || 0), 0)}</div>
          </div>
        </div>
      )}
    </div>
  );
}