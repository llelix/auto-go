'use client';

import React, { useState } from 'react';

interface SafeDropZoneProps {
  children: React.ReactNode;
  onDrop: (data: any) => void;
}

export function SafeDropZone({ children, onDrop }: SafeDropZoneProps) {
  const [isDragOver, setIsDragOver] = useState(false);
  const [isClient, setIsClient] = useState(false);

  React.useEffect(() => {
    setIsClient(true);
  }, []);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'copy';
    setIsDragOver(true);
  };

  const handleDragLeave = () => {
    setIsDragOver(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
    
    if (!isClient) return;

    try {
      const data = e.dataTransfer?.getData('text/plain');
      if (data) {
        const parsedData = JSON.parse(data);
        onDrop(parsedData);
      }
    } catch (error) {
      console.error('Drop error:', error);
    }
  };

  return (
    <div
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      className={`relative ${isDragOver ? 'bg-slate-800/20' : ''} transition-colors duration-200`}
    >
      {children}
      {isDragOver && (
        <div className="absolute inset-0 border-4 border-dashed border-blue-400 rounded-lg pointer-events-none">
          <div className="flex items-center justify-center h-full">
            <div className="text-blue-400 text-lg font-medium">释放以添加组件</div>
          </div>
        </div>
      )}
    </div>
  );
}