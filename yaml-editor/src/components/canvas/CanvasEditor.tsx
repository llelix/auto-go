'use client';

import React from 'react';

export function CanvasEditor() {
  return (
    <div className="w-full h-full flex items-center justify-center">
      <div className="text-center">
        <h2 className="text-white text-xl mb-2">画布区域</h2>
        <p className="text-slate-400">从左侧拖拽组件到此处</p>
      </div>
    </div>
  );
}