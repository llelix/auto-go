import { ComponentLibraryItem } from '@/types/yaml';

export const componentLibrary: ComponentLibraryItem[] = [
  // 基础操作
  {
    id: 'wait_appear',
    name: '等待出现',
    description: '等待元素出现在页面上',
    icon: '⏳',
    category: 'basic',
    template: {
      type: 'wait_appear',
      selector: '',
      timeout: 5,
      error_message: '等待元素出现失败'
    }
  },
  {
    id: 'wait_disappear',
    name: '等待消失',
    description: '等待元素从页面消失',
    icon: '⏸️',
    category: 'basic',
    template: {
      type: 'wait_disappear',
      selector: '',
      timeout: 5,
      error_message: '等待元素消失失败'
    }
  },

  // 交互操作
  {
    id: 'fill',
    name: '填写',
    description: '在输入框中填写内容',
    icon: '✏️',
    category: 'interaction',
    template: {
      type: 'fill',
      selector: '',
      value: '',
      error_message: '填写内容失败'
    }
  },
  {
    id: 'click',
    name: '点击',
    description: '点击页面元素',
    icon: '👆',
    category: 'interaction',
    template: {
      type: 'click',
      selector: '',
      error_message: '点击元素失败'
    }
  },
  {
    id: 'select',
    name: '选择',
    description: '从下拉菜单中选择选项',
    icon: '📋',
    category: 'interaction',
    template: {
      type: 'select',
      selector: '',
      value: '',
      error_message: '选择选项失败'
    }
  },
  {
    id: 'hover',
    name: '悬停',
    description: '鼠标悬停在元素上',
    icon: '🖱️',
    category: 'interaction',
    template: {
      type: 'hover',
      selector: '',
      error_message: '悬停操作失败'
    }
  },
  {
    id: 'drag_drop',
    name: '拖拽',
    description: '拖拽元素到目标位置',
    icon: '🔄',
    category: 'interaction',
    template: {
      type: 'drag_drop',
      selector: '',
      target: '',
      error_message: '拖拽操作失败'
    }
  },

  // 验证操作
  {
    id: 'get_text',
    name: '获取文本',
    description: '获取元素的文本内容',
    icon: '📝',
    category: 'verification',
    template: {
      type: 'get_text',
      selector: '',
      output_key: '',
      error_message: '获取文本失败'
    }
  },
  {
    id: 'get_attribute',
    name: '获取属性',
    description: '获取元素的属性值',
    icon: '🏷️',
    category: 'verification',
    template: {
      type: 'get_attribute',
      selector: '',
      attribute: '',
      output_key: '',
      error_message: '获取属性失败'
    }
  }
];

// 按分类组织组件
export const categorizedComponents = {
  basic: componentLibrary.filter(item => item.category === 'basic'),
  interaction: componentLibrary.filter(item => item.category === 'interaction'),
  verification: componentLibrary.filter(item => item.category === 'verification'),
  extraction: componentLibrary.filter(item => item.category === 'extraction')
};

// 分类名称映射
export const categoryNames = {
  basic: '基础操作',
  interaction: '交互操作',
  verification: '验证操作',
  extraction: '数据提取'
};