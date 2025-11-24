'use client';

import React, { useCallback } from 'react';
import { Settings, Trash2 } from 'lucide-react';
import { CanvasNode } from '@/types/yaml';

interface ActionNodeProps {
  node: CanvasNode;
  isSelected: boolean;
  onSelect: () => void;
}

export function ActionNode({ node, isSelected, onSelect }: ActionNodeProps) {
  const handleDelete = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    // TODO: 实现删除功能
  }, []);

  const handleEdit = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    // TODO: 实现编辑功能
  }, []);

  const getTypeColor = (type: string) => {
    const colors = {
      wait_appear: 'bg-yellow-600',
      wait_disappear: 'bg-orange-600',
      fill: 'bg-blue-600',
      click: 'bg-green-600',
      select: 'bg-purple-600',
      hover: 'bg-pink-600',
      drag_drop: 'bg-indigo-600',
      get_text: 'bg-teal-600',
      get_attribute: 'bg-cyan-600',
    };
    return colors[type as keyof typeof colors] || 'bg-slate-600';
  };

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
        width: 200,
        minHeight: 60,
        zIndex: isSelected ? 10 : 1,
      }}
      onClick={onSelect}
    >
      {/* 节点头部 */}
      <div className="flex items-center justify-between p-2 border-b border-slate-700">
        <div className="flex items-center gap-2">
          <div className={`w-6 h-6 ${getTypeColor(node.data.action?.type || '')} rounded flex items-center justify-center text-xs text-white`}>
            {node.data.icon || '⚡'}
          </div>
          <div>
            <div className="text-white font-medium text-xs">{node.data.label}</div>
            <div className="text-slate-400 text-xs opacity-75">{node.data.action?.type}</div>
          </div>
        </div>
        
        <div className="flex items-center gap-1">
          <button
            onClick={handleEdit}
            className="p-1 hover:bg-slate-700 rounded transition-colors"
            title="编辑操作"
          >
            <Settings className="w-3 h-3 text-slate-400" />
          </button>
          <button
            onClick={handleDelete}
            className="p-1 hover:bg-slate-700 rounded transition-colors"
            title="删除操作"
          >
            <Trash2 className="w-3 h-3 text-red-400" />
          </button>
        </div>
      </div>

      {/* 节点内容 */}
      <div className="p-2">
        {node.data.action?.selector && (
          <div className="text-xs text-slate-300 truncate">
            选择器: {node.data.action.selector}
          </div>
        )}
        {node.data.action?.value && (
          <div className="text-xs text-slate-300 truncate">
            值: {String(node.data.action.value)}
          </div>
        )}
      </div>

      {/* 连接点 */}
      <div className="absolute -left-2 top-1/2 -translate-y-1/2 w-3 h-3 bg-slate-500 rounded-full border-2 border-slate-800 cursor-crosshair" />
      <div className="absolute -right-2 top-1/2 -translate-y-1/2 w-3 h-3 bg-slate-500 rounded-full border-2 border-slate-800 cursor-crosshair" />
    </div>
  );
}