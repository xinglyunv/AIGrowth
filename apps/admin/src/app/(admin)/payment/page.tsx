'use client';

import { useEffect, useState } from 'react';
import { CreditCard, Save, Loader2 } from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface PaymentConfig {
  id: string;
  name: string;
  channel: string;
  merchant_id: string;
  merchant_key: string;
  api_url: string;
  notify_url: string;
  return_url: string;
  is_active: boolean;
}

export default function PaymentPage() {
  const [config, setConfig] = useState<PaymentConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    apiRequest<PaymentConfig>('/payment/config')
      .then((res) => { if (res.data) setConfig(res.data); })
      .catch((e) => setError(e instanceof Error ? e.message : '加载失败'))
      .finally(() => setLoading(false));
  }, []);

  const update = <K extends keyof PaymentConfig>(key: K, value: PaymentConfig[K]) => {
    if (!config) return;
    setConfig((prev) => prev ? { ...prev, [key]: value } : prev);
  };

  const handleSave = () => {
    if (!config) return;
    setSaving(true);
    setError('');
    setSaved(false);
    apiRequest('/payment/config', { method: 'PUT', body: JSON.stringify(config) })
      .then(() => { setSaved(true); setTimeout(() => setSaved(false), 3000); })
      .catch((err) => setError(err?.message || '保存失败'))
      .finally(() => setSaving(false));
  };

  const inputClass = 'w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary text-gray-900';
  const labelClass = 'block text-sm font-medium text-gray-700 mb-1.5';

  if (loading) {
    return <div className="flex items-center justify-center py-20"><p className="text-gray-400"><Loader2 className="animate-spin inline h-5 w-5 mr-2" />加载中...</p></div>;
  }

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <CreditCard className="h-6 w-6 text-gray-700" />
        <h2 className="text-xl font-bold text-gray-900">支付配置</h2>
      </div>

      {saved && <div className="mb-4 px-4 py-3 bg-green-50 border border-green-200 rounded-md text-sm text-green-700">配置已保存</div>}
      {error && <div className="mb-4 px-4 py-3 bg-red-50 border border-red-200 rounded-md text-sm text-red-700">{error}</div>}

      <div className="bg-white rounded-lg border border-gray-200 p-6 mb-6">
        <h3 className="text-base font-semibold text-gray-900 mb-4">易支付 (Epay) 配置</h3>
        <div className="space-y-4">
          <div>
            <label className={labelClass}>支付渠道</label>
            <select className={inputClass} value={config?.channel ?? 'alipay'} onChange={(e) => update('channel', e.target.value)}>
              <option value="alipay">支付宝</option>
              <option value="wxpay">微信支付</option>
              <option value="qqpay">QQ 钱包</option>
              <option value="bank">网银支付</option>
            </select>
          </div>
          <div>
            <label className={labelClass}>商户 ID</label>
            <input type="text" className={inputClass} value={config?.merchant_id ?? ''} onChange={(e) => update('merchant_id', e.target.value)} placeholder="输入商户 ID" />
          </div>
          <div>
            <label className={labelClass}>商户密钥</label>
            <input type="password" className={inputClass} value={config?.merchant_key ?? ''} onChange={(e) => update('merchant_key', e.target.value)} placeholder="输入商户密钥" />
          </div>
          <div>
            <label className={labelClass}>API 地址</label>
            <input type="text" className={inputClass} value={config?.api_url ?? ''} onChange={(e) => update('api_url', e.target.value)} placeholder="https://epay.example.com" />
          </div>
          <div>
            <label className={labelClass}>异步通知 URL</label>
            <input type="text" className={inputClass} value={config?.notify_url ?? ''} onChange={(e) => update('notify_url', e.target.value)} placeholder="https://your-site.com/api/v1/billing/notify" />
          </div>
          <div>
            <label className={labelClass}>同步跳转 URL</label>
            <input type="text" className={inputClass} value={config?.return_url ?? ''} onChange={(e) => update('return_url', e.target.value)} placeholder="https://your-site.com/dashboard/billing" />
          </div>
          <label className="flex items-center gap-2 cursor-pointer">
            <input type="checkbox" checked={config?.is_active ?? false}
              onChange={(e) => update('is_active', e.target.checked)}
              className="rounded border-gray-300 text-primary focus:ring-primary" />
            <span className="text-sm text-gray-700">启用支付</span>
          </label>
        </div>
      </div>

      <button onClick={handleSave} disabled={saving || !config}
        className="flex items-center gap-1.5 px-5 py-2.5 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 transition-colors disabled:opacity-60">
        {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
        保存配置
      </button>
    </div>
  );
}
