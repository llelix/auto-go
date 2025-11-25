'use client';

import React from 'react';
import { SafeDraggable } from '../SafeDraggable';

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
          
          <SafeDraggable
            id="wait-appear"
            type="wait_appear"
            icon="⏳"
            title="等待出现"
            description="等待元素出现在页面上"
          />

          <SafeDraggable
            id="click"
            type="click"
            icon="👆"
            title="点击"
            description="点击页面元素"
          />

          <SafeDraggable
            id="fill"
            type="fill"
            icon="✏️"
            title="填写"
            description="在输入框中填写内容"
          />
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