'use client';

import React, { useRef } from 'react';
import { Upload, Download } from 'lucide-react';

interface ToolbarProps {
  onImport: (file: File) => void;
  onExport: () => void;
}

export function Toolbar({ onImport, onExport }: ToolbarProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      onImport(file);
    }
    // 重置input值，允许重复选择同一文件
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleImportClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <div className="h-14 bg-slate-800 border-b border-slate-700 flex items-center justify-between px-4">
      <div className="flex items-center gap-2">
        <h1 className="text-lg font-semibold text-white">YAML 可视化编辑器</h1>
      </div>
      <div className="flex items-center gap-4">
        {/* 隐藏的文件输入 */}
        <input
          ref={fileInputRef}
          type="file"
          accept=".yaml,.yml,.json"
          onChange={handleFileChange}
          className="hidden"
        />
        
        <button 
          onClick={handleImportClick}
          className="flex items-center gap-2 px-3 py-1.5 bg-slate-700 hover:bg-slate-600 text-white rounded-md transition-colors text-sm"
        >
          <Upload className="w-4 h-4" />
          导入
        </button>
        
        <button 
          onClick={onExport}
          className="flex items-center gap-2 px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded-md transition-colors text-sm"
        >
          <Download className="w-4 h-4" />
          导出
        </button>
      </div>
    </div>
  );
}