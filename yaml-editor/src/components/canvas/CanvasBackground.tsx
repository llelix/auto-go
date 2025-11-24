'use client';

import React from 'react';

export function CanvasBackground() {
  return (
    <div className="absolute inset-0 bg-slate-950">
      {/* 网格背景 */}
      <svg className="absolute inset-0 w-full h-full">
        <defs>
          <pattern
            id="grid"
            width="20"
            height="20"
            patternUnits="userSpaceOnUse"
          >
            <path
              d="M 20 0 L 0 0 0 20"
              fill="none"
              stroke="#1e293b"
              strokeWidth="1"
            />
          </pattern>
          <pattern
            id="grid-large"
            width="100"
            height="100"
            patternUnits="userSpaceOnUse"
          >
            <path
              d="M 100 0 L 0 0 0 100"
              fill="none"
              stroke="#334155"
              strokeWidth="1"
            />
          </pattern>
        </defs>
        
        {/* 小网格 */}
        <rect width="100%" height="100%" fill="url(#grid)" />
        {/* 大网格 */}
        <rect width="100%" height="100%" fill="url(#grid-large)" />
      </svg>
    </div>
  );
}