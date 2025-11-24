'use client';

import React from 'react';
import { useDroppable } from '@dnd-kit/core';

export function DroppableCanvas() {
  const { setNodeRef } = useDroppable({
    id: 'canvas-dropzone',
  });

  return (
    <div
      ref={setNodeRef}
      className="absolute inset-0 bg-slate-900/10 flex items-center justify-center"
    >
      <div className="text-center">
        <h2 className="text-white text-xl mb-2">画布区域</h2>
        <p className="text-slate-400">从左侧拖拽组件到此处</p>
      </div>
    </div>
  );
}