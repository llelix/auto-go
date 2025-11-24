'use client';

import React, { useCallback } from 'react';
import { Play, Settings, Trash2 } from 'lucide-react';
import { CanvasNode } from '@/types/yaml';

interface TaskNodeProps {
  node: CanvasNode;
  isSelected: boolean;
  onSelect: () => void;
}

export function TaskNode({ node, isSelected, onSelect }: TaskNodeProps) {
  const handleDelete = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    // TODO: 实现删除功能
  }, []);

  const handleEdit = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    // TODO: 实现编辑功能
  }, []);

  return (
    <div
      className={`
        absolute bg-slate-800 border-2 rounded-lg shadow-lg cursor-pointer transition-all
        ${isSelected 
          ? 'border-blue-400 shadow-blue-400/20' 
          : 'border-slate-600 hover:border-slate-500'
        }
      `}
      style={{
        left: node.position.x,
        top: node.position.y,
        width: 240,
        minHeight: 80,
        zIndex: isSelected ? 10 : 1,
      }}
      onClick={onSelect}
    >
      {/* 节点头部 */}
      <div className="flex items-center justify-between p-3 border-b border-slate-700">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center">
            <Play className="w-4 h-4 text-white" />
          </div>
          <div>
            <div className="text-white font-medium text-sm">任务</div>
            <div className="text-slate-400 text-xs">{node.data.label}</div>
          </div>
        </div>
        
        <div className="flex items-center gap-1">
          <button
            onClick={handleEdit}
            className="p-1.5 hover:bg-slate-700 rounded transition-colors"
            title="编辑任务"
          >
            <Settings className="w-3 h-3 text-slate-400" />
          </button>
          <button
            onClick={handleDelete}
            className="p-1.5 hover:bg-slate-700 rounded transition-colors"
            title="删除任务"
          >
            <Trash2 className="w-3 h-3 text-red-400" />
          </button>
        </div>
      </div>

      {/* 节点内容 */}
      <div className="p-3">
        {node.data.task && (
          <div className="space-y-2">
            <div className="flex justify-between text-xs">
              <span className="text-slate-400">URL:</span>
              <span className="text-slate-300 truncate ml-2">{node.data.task.url}</span>
            </div>
            <div className="flex justify-between text-xs">
              <span className="text-slate-400">操作数:</span>
              <span className="text-slate-300">{node.data.task.actions?.length || 0}</span>
            </div>
          </div>
        )}
      </div>

      {/* 连接点 */}
      <div className="absolute -right-2 top-1/2 -translate-y-1/2 w-4 h-4 bg-blue-500 rounded-full border-2 border-slate-800 cursor-crosshair" />
    </div>
  );
}