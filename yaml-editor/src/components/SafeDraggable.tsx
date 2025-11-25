'use client';

import React from 'react';

interface DraggableComponentProps {
  id: string;
  type: string;
  icon: string;
  title: string;
  description: string;
  isDraggable?: boolean;
}

// 安全的可拖拽组件，避免SSR问题
export function SafeDraggable({ 
  id, 
  type, 
  icon, 
  title, 
  description,
  isDraggable = true 
}: DraggableComponentProps) {
  const [isClient, setIsClient] = React.useState(false);

  React.useEffect(() => {
    setIsClient(true);
  }, []);

  const handleDragStart = (e: React.DragEvent) => {
    if (isDraggable && isClient) {
      e.dataTransfer?.setData('text/plain', JSON.stringify({
        id,
        type,
        title,
        description
      }));
      e.dataTransfer.effectAllowed = 'copy';
    }
  };

  if (!isClient) {
    return (
      <div className="flex items-center gap-3 p-3 bg-slate-600 rounded-lg border border-slate-500">
        <div className="text-2xl">{icon}</div>
        <div className="flex-1 min-w-0">
          <div className="text-white font-medium text-sm">{title}</div>
          <div className="text-slate-400 text-xs mt-1">{description}</div>
        </div>
      </div>
    );
  }

  return (
    <div
      draggable={isDraggable}
      onDragStart={handleDragStart}
      className="flex items-center gap-3 p-3 bg-slate-700 hover:bg-slate-600 rounded-lg cursor-move transition-all border border-slate-600 hover:border-slate-500"
    >
      <div className="text-2xl">{icon}</div>
      <div className="flex-1 min-w-0">
        <div className="text-white font-medium text-sm">{title}</div>
        <div className="text-slate-400 text-xs mt-1">{description}</div>
      </div>
    </div>
  );
}