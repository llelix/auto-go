'use client';

import React, { useRef } from 'react';
import { Upload, Download, Eye, EyeOff, FileText, Trash2 } from 'lucide-react';
import { useEditorStore } from '@/lib/store/editor-store';

interface ToolbarProps {
  onImport: (event: React.ChangeEvent<HTMLInputElement>) => void;
  onExport: () => void;
  onTogglePreview: () => void;
  showPreview: boolean;
}

export function Toolbar({ onImport, onExport, onTogglePreview, showPreview }: ToolbarProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { isDirty, resetEditor } = useEditorStore();

  const handleReset = () => {
    if (isDirty) {
      if (confirm('确定要重置编辑器吗？所有未保存的更改将会丢失。')) {
        resetEditor();
      }
    } else {
      resetEditor();
    }
  };

  return (
    <div className="h-14 bg-slate-800 border-b border-slate-700 flex items-center justify-between px-4">
      {/* 左侧标题和导入导出 */}
      <div className="flex items-center gap-4">
        <div className="flex items-center gap-2">
          <FileText className="w-5 h-5 text-blue-400" />
          <h1 className="text-lg font-semibold text-white">YAML 可视化编辑器</h1>
          {isDirty && (
            <div className="w-2 h-2 bg-orange-500 rounded-full animate-pulse" title="有未保存的更改" />
          )}
        </div>

        <div className="flex items-center gap-2">
          {/* 隐藏的文件输入 */}
          <input
            ref={fileInputRef}
            type="file"
            accept=".yaml,.yml"
            onChange={onImport}
            className="hidden"
          />

          <button
            onClick={() => fileInputRef.current?.click()}
            className="flex items-center gap-2 px-3 py-1.5 bg-slate-700 hover:bg-slate-600 text-white rounded-md transition-colors text-sm"
            title="导入YAML文件"
          >
            <Upload className="w-4 h-4" />
            导入
          </button>

          <button
            onClick={onExport}
            className="flex items-center gap-2 px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded-md transition-colors text-sm"
            title="导出YAML文件"
          >
            <Download className="w-4 h-4" />
            导出
          </button>

          <button
            onClick={handleReset}
            className="flex items-center gap-2 px-3 py-1.5 bg-red-600 hover:bg-red-500 text-white rounded-md transition-colors text-sm"
            title="重置编辑器"
          >
            <Trash2 className="w-4 h-4" />
            重置
          </button>
        </div>
      </div>

      {/* 右侧预览切换 */}
      <div className="flex items-center gap-4">
        <button
          onClick={onTogglePreview}
          className={`flex items-center gap-2 px-3 py-1.5 rounded-md transition-colors text-sm ${
            showPreview 
              ? 'bg-green-600 hover:bg-green-500 text-white' 
              : 'bg-slate-700 hover:bg-slate-600 text-white'
          }`}
          title={showPreview ? '隐藏预览面板' : '显示预览面板'}
        >
          {showPreview ? <Eye className="w-4 h-4" /> : <EyeOff className="w-4 h-4" />}
          {showPreview ? '隐藏预览' : '显示预览'}
        </button>
      </div>
    </div>
  );
}