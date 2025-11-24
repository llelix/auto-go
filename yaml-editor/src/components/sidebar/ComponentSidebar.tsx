'use client';

import React from 'react';

export function ComponentSidebar() {
  return (
    <div className="h-full flex flex-col">
      {/* 标题 */}
      <div className="p-4 border-b border-slate-700">
        <h2 className="text-white font-semibold text-lg">组件库</h2>
        <p className="text-slate-400 text-sm mt-1">拖拽组件到画布中创建操作</p>
      </div>

      {/* 组件列表 */}
      <div className="flex-1 overflow-y-auto p-4">
        <div className="space-y-4">
          <div className="text-slate-300 font-medium text-sm">基础操作</div>
          
          <div className="flex items-center gap-3 p-3 bg-slate-700 hover:bg-slate-600 rounded-lg cursor-move transition-all border border-slate-600 hover:border-slate-500">
            <div className="text-2xl">⏳</div>
            <div className="flex-1 min-w-0">
              <div className="text-white font-medium text-sm">等待出现</div>
              <div className="text-slate-400 text-xs mt-1">等待元素出现在页面上</div>
            </div>
          </div>

          <div className="flex items-center gap-3 p-3 bg-slate-700 hover:bg-slate-600 rounded-lg cursor-move transition-all border border-slate-600 hover:border-slate-500">
            <div className="text-2xl">👆</div>
            <div className="flex-1 min-w-0">
              <div className="text-white font-medium text-sm">点击</div>
              <div className="text-slate-400 text-xs mt-1">点击页面元素</div>
            </div>
          </div>

          <div className="flex items-center gap-3 p-3 bg-slate-700 hover:bg-slate-600 rounded-lg cursor-move transition-all border border-slate-600 hover:border-slate-500">
            <div className="text-2xl">✏️</div>
            <div className="flex-1 min-w-0">
              <div className="text-white font-medium text-sm">填写</div>
              <div className="text-slate-400 text-xs mt-1">在输入框中填写内容</div>
            </div>
          </div>
        </div>
      </div>

      {/* 使用提示 */}
      <div className="p-4 border-t border-slate-700 bg-slate-750">
        <div className="text-slate-400 text-xs space-y-1">
          <div>💡 提示：</div>
          <div>• 拖拽组件到画布添加操作</div>
          <div>• 双击节点编辑属性</div>
          <div>• 拖拽节点边缘连接流程</div>
        </div>
      </div>
    </div>
  );
}