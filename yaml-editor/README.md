# YAML 可视化编辑器

一个基于 Next.js 的 YAML 可视化编辑器，专门为 auto-go 自动化测试工具设计，支持通过拖拽方式创建和编辑任务配置文件。

## 功能特性

- 🎨 **可视化编辑器**: 通过拖拽组件创建测试流程
- 📝 **实时预览**: 实时显示生成的 YAML 代码
- 🔄 **文件操作**: 支持导入/导出 YAML 文件
- ✅ **配置验证**: 自动验证配置格式正确性
- 🎯 **组件库**: 丰富的自动化操作组件
- 🖱️ **交互友好**: 直观的拖拽界面和连接线

## 技术栈

- **Next.js 14** - React 框架
- **TypeScript** - 类型安全
- **Tailwind CSS** - 样式框架
- **@dnd-kit** - 拖拽功能
- **js-yaml** - YAML 解析
- **Zustand** - 状态管理
- **shadcn/ui** - UI 组件库

## 项目结构

```
src/
├── app/                    # Next.js App Router
├── components/
│   ├── ui/                # UI 组件
│   ├── canvas/            # 画布相关组件
│   ├── sidebar/           # 侧边栏组件
│   └── preview/           # 预览组件
├── lib/
│   ├── yaml-parser/       # YAML 解析器
│   ├── drag-drop/         # 拖拽逻辑
│   └── store/             # 状态管理
└── types/                 # TypeScript 类型定义
```

## 快速开始

1. 安装依赖：
```bash
npm install
```

2. 启动开发服务器：
```bash
npm run dev
```

3. 打开浏览器访问 `http://localhost:3000`

## 使用说明

### 基本操作

1. **添加组件**: 从左侧组件库拖拽组件到画布
2. **编辑属性**: 双击节点编辑组件属性
3. **连接流程**: 拖拽节点边缘连接操作流程
4. **预览代码**: 右侧实时显示 YAML 代码
5. **导出文件**: 点击工具栏导出按钮保存配置

### 支持的操作类型

- **基础操作**: 等待出现、等待消失
- **交互操作**: 填写、点击、选择、悬停、拖拽
- **验证操作**: 获取文本、获取属性

## 项目截图

![编辑器界面](screenshots/editor.png)

## 开发
task_sample.yaml是样例

### 可用命令

```bash
npm run dev          # 启动开发服务器
npm run build        # 构建生产版本
npm run start        # 启动生产服务器
npm run lint         # 代码检查
```

### 自定义组件

要添加新的操作组件，请在 `lib/drag-drop/component-library.ts` 中配置：

```typescript
{
  id: 'custom_action',
  name: '自定义操作',
  description: '操作描述',
  icon: '🎯',
  category: 'interaction',
  template: {
    type: 'custom_action',
    selector: '',
    error_message: '操作失败'
  }
}
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License