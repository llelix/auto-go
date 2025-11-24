'use client';

import React from 'react';

export function YamlEditor() {
  return (
    <div className="h-screen w-full bg-slate-900 flex flex-col overflow-hidden">
      {/* 顶部工具栏 */}
      <div className="h-14 bg-slate-800 border-b border-slate-700 flex items-center justify-between px-4">
        <div className="flex items-center gap-2">
          <h1 className="text-lg font-semibold text-white">YAML 可视化编辑器</h1>
        </div>
        <div className="flex items-center gap-4">
          <button className="flex items-center gap-2 px-3 py-1.5 bg-slate-700 hover:bg-slate-600 text-white rounded-md transition-colors text-sm">
            导入
          </button>
          <button className="flex items-center gap-2 px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white rounded-md transition-colors text-sm">
            导出
          </button>
        </div>
      </div>

      {/* 主编辑区域 */}
      <div className="flex-1 flex overflow-hidden">
        {/* 左侧组件库 */}
        <div className="w-80 bg-slate-800 border-r border-slate-700 flex-shrink-0">
          <div className="h-full flex flex-col">
            <div className="p-4 border-b border-slate-700">
              <h2 className="text-white font-semibold text-lg">组件库</h2>
              <p className="text-slate-400 text-sm mt-1">拖拽组件到画布中创建操作</p>
            </div>
            
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

            <div className="p-4 border-t border-slate-700 bg-slate-750">
              <div className="text-slate-400 text-xs space-y-1">
                <div>💡 提示：</div>
                <div>• 拖拽组件到画布添加操作</div>
                <div>• 双击节点编辑属性</div>
                <div>• 拖拽节点边缘连接流程</div>
              </div>
            </div>
          </div>
        </div>

        {/* 中间画布区域 */}
        <div className="flex-1 bg-slate-950 relative overflow-hidden flex items-center justify-center">
          <div className="text-center">
            <h2 className="text-white text-xl mb-2">画布区域</h2>
            <p className="text-slate-400">从左侧拖拽组件到此处</p>
          </div>
        </div>

        {/* 右侧预览面板 */}
        <div className="w-96 bg-slate-800 border-l border-slate-700 flex-shrink-0">
          <div className="h-full flex flex-col">
            <div className="p-4 border-b border-slate-700 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <h2 className="text-white font-semibold">YAML 预览</h2>
              </div>
              <button className="flex items-center gap-2 px-3 py-1.5 bg-slate-700 hover:bg-slate-600 text-white rounded-md transition-colors text-sm">
                复制
              </button>
            </div>

            <div className="flex-1 overflow-hidden">
              <div className="h-full overflow-auto p-4 text-sm font-mono text-slate-300 bg-slate-950">
                <pre>
                  <code>{`# AutoGo 任务配置文件
# 从左侧拖拽组件开始创建配置

- name: "示例任务"
  url: "https://example.com"
  wait_time: 3
  screenshot: true
  actions:
    - type: "wait_appear"
      selector: "#element"
      timeout: 5
      error_message: "等待元素出现失败"`}</code>
                </pre>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}