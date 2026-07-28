'use client';

import { useEffect, useState } from 'react';
import { Settings, Save } from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface SystemSettings {
  siteName: string;
  siteDescription: string;
  adminEmail: string;
  defaultModel: string;
  maxQuestions: number;
  taskTimeout: string;
  emailNotification: boolean;
  taskNotification: boolean;
}

interface AIModel {
  id: string;
  name: string;
  model: string;
  enabled: boolean;
}

const DEFAULT_SETTINGS: SystemSettings = {
  siteName: 'AI Growth Engine',
  siteDescription: 'AI 品牌可见度分析与增长优化平台',
  adminEmail: 'admin@aigrowthengine.com',
  defaultModel: '',
  maxQuestions: 5,
  taskTimeout: '30m',
  emailNotification: true,
  taskNotification: true,
};

const TIMEOUT_OPTIONS = [
  { value: '15m', label: '15 分钟' },
  { value: '30m', label: '30 分钟' },
  { value: '1h', label: '1 小时' },
  { value: '2h', label: '2 小时' },
  { value: '6h', label: '6 小时' },
];

export default function SystemPage() {
  const [settings, setSettings] = useState<SystemSettings>(DEFAULT_SETTINGS);
  const [models, setModels] = useState<AIModel[]>([]);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    Promise.all([apiRequest<SystemSettings>('/settings'), apiRequest<AIModel[]>('/models')])
      .then(([settingsResponse, modelsResponse]) => {
        if (settingsResponse.data) setSettings(settingsResponse.data);
        if (Array.isArray(modelsResponse.data)) setModels(modelsResponse.data.filter((model) => model.enabled));
      })
      .catch(() => {});
  }, []);

  const modelOptions = models.map((model) => ({ value: model.model, label: model.name }));
  if (settings.defaultModel && !modelOptions.some((option) => option.value === settings.defaultModel)) {
    modelOptions.unshift({ value: settings.defaultModel, label: `${settings.defaultModel}（当前配置）` });
  }

  const update = <K extends keyof SystemSettings>(key: K, value: SystemSettings[K]) => {
    setSettings((prev) => ({ ...prev, [key]: value }));
  };

  const handleSave = () => {
    apiRequest('/settings', {
      method: 'PUT',
      body: JSON.stringify(settings),
    })
      .then(() => {
        setSaved(true);
        setTimeout(() => setSaved(false), 3000);
      })
      .catch(() => {
        setSaved(true);
        setTimeout(() => setSaved(false), 3000);
      });
  };

  const inputClass =
    'w-full border border-gray-300 rounded-md px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary';

  const labelClass = 'block text-sm font-medium text-gray-700 mb-1.5';
  const sectionCard = 'bg-white rounded-lg border border-gray-200 p-5 mb-4';

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Settings className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">系统设置</h2>
        </div>
      </div>

      {saved && (
        <div className="mb-4 px-4 py-3 bg-green-50 border border-green-200 rounded-md text-sm text-green-700">
          设置已保存
        </div>
      )}

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">基本设置</h3>
        <div className="space-y-4">
          <div>
            <label className={labelClass}>站点名称</label>
            <input
              type="text"
              className={inputClass}
              value={settings.siteName}
              onChange={(e) => update('siteName', e.target.value)}
            />
          </div>
          <div>
            <label className={labelClass}>站点描述</label>
            <textarea
              className={`${inputClass} resize-none`}
              rows={2}
              value={settings.siteDescription}
              onChange={(e) => update('siteDescription', e.target.value)}
            />
          </div>
          <div>
            <label className={labelClass}>管理员邮箱</label>
            <input
              type="email"
              className={inputClass}
              value={settings.adminEmail}
              onChange={(e) => update('adminEmail', e.target.value)}
            />
          </div>
        </div>
      </div>

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">分析设置</h3>
        <div className="space-y-4">
          <div>
            <label className={labelClass}>默认 AI 模型</label>
            <select
              className={inputClass}
              value={settings.defaultModel}
              onChange={(e) => update('defaultModel', e.target.value)}
              disabled={modelOptions.length === 0}
            >
              {modelOptions.length === 0 && <option value="">暂无已启用模型</option>}
              {modelOptions.map((o) => (
                <option key={o.value} value={o.value}>
                  {o.label}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className={labelClass}>每个任务最大问题数</label>
            <input
              type="number"
              min={1}
              max={50}
              className={inputClass}
              value={settings.maxQuestions}
              onChange={(e) => update('maxQuestions', parseInt(e.target.value) || 1)}
            />
          </div>
          <div>
            <label className={labelClass}>任务超时时间</label>
            <select
              className={inputClass}
              value={settings.taskTimeout}
              onChange={(e) => update('taskTimeout', e.target.value)}
            >
              {TIMEOUT_OPTIONS.map((o) => (
                <option key={o.value} value={o.value}>
                  {o.label}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">通知设置</h3>
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-700">邮件通知</p>
              <p className="text-xs text-gray-400 mt-0.5">系统通知将发送到管理员邮箱</p>
            </div>
            <button
              onClick={() => update('emailNotification', !settings.emailNotification)}
              className={`relative w-10 h-5 rounded-full transition-colors ${
                settings.emailNotification ? 'bg-primary' : 'bg-gray-300'
              }`}
            >
              <span
                className={`absolute top-0.5 left-0.5 w-4 h-4 bg-white rounded-full shadow transition-transform ${
                  settings.emailNotification ? 'translate-x-5' : 'translate-x-0'
                }`}
              />
            </button>
          </div>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-700">任务完成通知</p>
              <p className="text-xs text-gray-400 mt-0.5">分析任务完成后发送通知</p>
            </div>
            <button
              onClick={() => update('taskNotification', !settings.taskNotification)}
              className={`relative w-10 h-5 rounded-full transition-colors ${
                settings.taskNotification ? 'bg-primary' : 'bg-gray-300'
              }`}
            >
              <span
                className={`absolute top-0.5 left-0.5 w-4 h-4 bg-white rounded-full shadow transition-transform ${
                  settings.taskNotification ? 'translate-x-5' : 'translate-x-0'
                }`}
              />
            </button>
          </div>
        </div>
      </div>

      <button
        onClick={handleSave}
        className="flex items-center gap-1.5 px-5 py-2.5 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 transition-colors"
      >
        <Save className="h-4 w-4" />
        保存设置
      </button>
    </div>
  );
}
