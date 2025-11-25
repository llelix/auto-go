'use client';

import React, { useState } from 'react';
import { SafeDropZone } from '../SafeDropZone';

export interface DroppedItem {
  id: string;
  type: string;
  title: string;
  description: string;
  x: number;
  y: number;
}

interface DroppableCanvasProps {
  items: DroppedItem[];
  onItemsChange: (items: DroppedItem[]) => void;
}

function SimpleActionNode({ item }: { item: DroppedItem }) {
  const getTypeIcon = (type: string) => {
    const icons = {
      wait_appear: '⏳',
      click: '👆',
      fill: '✏️',
    };
    return icons[type as keyof typeof icons] || '⚡';
  };

  const getTypeColor = (type: string) => {
    const colors = {
      wait_appear: 'bg-yellow-600',
      click: 'bg-green-600',
      fill: 'bg-blue-600',
    };
    return colors[type as keyof typeof colors] || 'bg-slate-600';
  };

  return (
    <div className="bg-slate-800 border-2 border-slate-600 rounded-lg shadow-lg hover:border-slate-500 transition-all cursor-pointer" style={{ width: 200, minHeight: 60 }}>
      <div className="flex items-center gap-2 p-3 border-b border-slate-700">
        <div className={`w-8 h-8 ${getTypeColor(item.type)} rounded flex items-center justify-center text-sm text-white`}>
          {getTypeIcon(item.type)}
        </div>
        <div className="flex-1 min-w-0">
          <div className="text-white font-medium text-sm">{item.title}</div>
          <div className="text-slate-400 text-xs">{item.type}</div>
        </div>
      </div>
      <div className="p-2">
        <div className="text-xs text-slate-300">{item.description}</div>
      </div>
    </div>
  );
}

export function DroppableCanvas({ items, onItemsChange }: DroppableCanvasProps) {
  const [isClient, setIsClient] = useState(false);

  React.useEffect(() => {
    setIsClient(true);
  }, []);

  const handleDrop = (data: any) => {
    const newItem: DroppedItem = {
      id: `${data.id}-${Date.now()}`,
      type: data.type || 'unknown',
      title: data.title || '未知操作',
      description: data.description || '',
      x: 100 + (items.length % 3) * 220,
      y: 50 + Math.floor(items.length / 3) * 120,
    };

    onItemsChange([...items, newItem]);
  };

  if (!isClient) {
    return (
      <div className="absolute inset-0 bg-slate-900/10">
        <div className="h-full flex items-center justify-center">
          <div className="text-center">
            <h2 className="text-white text-xl mb-2">画布区域</h2>
            <p className="text-slate-400">从左侧拖拽组件到此处</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <SafeDropZone onDrop={handleDrop}>
      <div className="absolute inset-0 bg-slate-900/10">
        {items.length === 0 ? (
          <div className="h-full flex items-center justify-center">
            <div className="text-center">
              <h2 className="text-white text-xl mb-2">画布区域</h2>
              <p className="text-slate-400">从左侧拖拽组件到此处</p>
            </div>
          </div>
        ) : (
          <div className="relative w-full h-full">
            {items.map(item => (
              <div
                key={item.id}
                style={{
                  position: 'absolute',
                  left: `${item.x}px`,
                  top: `${item.y}px`,
                }}
              >
                <SimpleActionNode item={item} />
              </div>
            ))}
          </div>
        )}
      </div>
    </SafeDropZone>
  );
}