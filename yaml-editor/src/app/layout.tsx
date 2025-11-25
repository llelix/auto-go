import type { Metadata } from "next";
import "./globals.css";

// 使用系统字体，避免网络下载问题
const systemFonts = {
  sans: {
    variable: "--font-sans",
    style: {
      fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif',
    }
  },
  mono: {
    variable: "--font-mono",
    style: {
      fontFamily: '"SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas, "Courier New", monospace',
    }
  }
};

export const metadata: Metadata = {
  title: "YAML 可视化编辑器",
  description: "基于Next.js的YAML可视化编辑器，专为auto-go自动化测试工具设计",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body
        className={`${systemFonts.sans.variable} ${systemFonts.mono.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}