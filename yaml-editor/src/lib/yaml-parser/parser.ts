import * as yaml from 'js-yaml';
import { Task } from '@/types/yaml';

export class YamlParser {
  /**
   * 解析YAML字符串为配置对象
   */
  static parseYaml(yamlString: string): Task[] {
    try {
      const parsed = yaml.load(yamlString) as Task[];
      return parsed || [];
    } catch (error) {
      console.error('YAML解析错误:', error);
      throw new Error('YAML格式错误，请检查语法');
    }
  }

  /**
   * 将配置对象转换为YAML字符串
   */
  static stringifyYaml(tasks: Task[]): string {
    try {
      return yaml.dump(tasks, {
        indent: 2,
        lineWidth: 120,
        noRefs: true,
        sortKeys: false,
        styles: {
          '!!null': 'canonical'
        }
      });
    } catch (error) {
      console.error('YAML生成错误:', error);
      throw new Error('YAML生成失败');
    }
  }

  /**
   * 验证任务配置的完整性
   */
  static validateTask(task: Task): { valid: boolean; errors: string[] } {
    const errors: string[] = [];

    if (!task.name || task.name.trim() === '') {
      errors.push('任务名称不能为空');
    }

    if (!task.url || task.url.trim() === '') {
      errors.push('URL不能为空');
    }

    if (!task.actions || task.actions.length === 0) {
      errors.push('至少需要一个操作');
    }

    return {
      valid: errors.length === 0,
      errors
    };
  }

  /**
   * 验证整个配置文件
   */
  static validateConfig(tasks: Task[]): { valid: boolean; errors: string[] } {
    const allErrors: string[] = [];

    if (!tasks || tasks.length === 0) {
      allErrors.push('至少需要一个任务');
    }

    tasks.forEach((task, index) => {
      const validation = this.validateTask(task);
      if (!validation.valid) {
        validation.errors.forEach(error => {
          allErrors.push(`任务${index + 1}: ${error}`);
        });
      }
    });

    return {
      valid: allErrors.length === 0,
      errors: allErrors
    };
  }

  /**
   * 从文件加载YAML
   */
  static async loadFromFile(file: File): Promise<Task[]> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        try {
          const yamlString = e.target?.result as string;
          const tasks = this.parseYaml(yamlString);
          resolve(tasks);
        } catch (error) {
          reject(error);
        }
      };
      reader.onerror = () => reject(new Error('文件读取失败'));
      reader.readAsText(file);
    });
  }

  /**
   * 下载YAML文件
   */
  static downloadFile(tasks: Task[], filename: string = 'tasks.yaml'): void {
    try {
      const yamlString = this.stringifyYaml(tasks);
      const blob = new Blob([yamlString], { type: 'text/yaml' });
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    } catch (error) {
      console.error('下载失败:', error);
      throw new Error('文件下载失败');
    }
  }
}