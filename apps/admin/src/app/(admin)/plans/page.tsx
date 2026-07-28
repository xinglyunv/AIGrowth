'use client';

import { useEffect, useState } from 'react';
import {
  CreditCard, Plus, X, Loader2, Save, Edit3, Check, Zap, Layers
} from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface Plan {
  id: string;
  name: string;
  code: string;
  description: string;
  monthly_price: number;
  yearly_price: number;
  max_projects: number;
  max_ai_queries: number;
  max_reports: number;
  credits: number;
  features: string;
  popular: boolean;
  is_active: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export default function PlansPage() {
  const [plans, setPlans] = useState<Plan[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showEditor, setShowEditor] = useState(false);
  const [editing, setEditing] = useState<Plan | null>(null);
  const [saving, setSaving] = useState(false);
  const [form, setForm] = useState({
    name: '', monthly_price: 0, yearly_price: 0, max_projects: 0, max_ai_queries: 0,
    credits: 0, features: '', popular: false
  });

  const fetchPlans = () => {
    apiRequest<Plan[]>('/plans')
      .then((res) => { if (Array.isArray(res.data)) setPlans(res.data); })
      .catch((e) => setError(e instanceof Error ? e.message : '加载失败'))
      .finally(() => setLoading(false));
  };

  useEffect(() => { fetchPlans(); }, []);

  const openCreate = () => {
    setEditing(null);
    setForm({ name: '', monthly_price: 0, yearly_price: 0, max_projects: 0, max_ai_queries: 0, credits: 0, features: '', popular: false });
    setShowEditor(true);
  };

  const openEdit = (p: Plan) => {
    setEditing(p);
    setForm({
      name: p.name, monthly_price: p.monthly_price, yearly_price: p.yearly_price,
      max_projects: p.max_projects, max_ai_queries: p.max_ai_queries,
      credits: p.credits, features: p.features, popular: p.popular
    });
    setShowEditor(true);
  };

  const handleSave = () => {
    if (!form.name) return;
    setSaving(true);
    const payload = { ...form, features: form.features };
    const request = editing
      ? apiRequest(`/plans/${editing.id}`, { method: 'PUT', body: JSON.stringify(payload) })
      : apiRequest('/plans', { method: 'POST', body: JSON.stringify(payload) });
    request
      .then(() => {
        setShowEditor(false);
        fetchPlans();
      })
      .catch((e) => setError(e instanceof Error ? e.message : '保存失败'))
      .finally(() => setSaving(false));
  };

  const inputClass = 'w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary text-gray-900';
  const labelClass = 'block text-sm font-medium text-gray-700 mb-1.5';

  const featureList = (features: string) =>
    features.split('\n').filter(Boolean).map((f, i) => (
      <li key={i} className="flex items-start gap-1.5 text-xs text-gray-600">
        <Check className="h-3.5 w-3.5 text-green-500 mt-0.5 shrink-0" />
        {f}
      </li>
    ));

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <CreditCard className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">套餐管理</h2>
        </div>
        <button onClick={openCreate}
          className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90">
          <Plus className="h-4 w-4" />
          新建套餐
        </button>
      </div>

      {loading && <div className="p-12 text-center text-gray-400"><Loader2 className="animate-spin inline h-5 w-5" /> 加载中...</div>}

      {error && <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-sm text-red-700">{error}</div>}

      {!loading && !error && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {plans.map((plan) => (
            <div key={plan.id} className={`bg-white rounded-lg border ${plan.popular ? 'border-primary/40 ring-1 ring-primary/20' : 'border-gray-200'} p-5 relative`}>
              {plan.popular && <span className="absolute -top-2.5 right-4 px-2 py-0.5 bg-primary text-white text-[10px] font-medium rounded-full">推荐</span>}
              <h3 className="text-base font-bold text-gray-900 mb-1">{plan.name}</h3>
              {plan.description && <p className="text-xs text-gray-400 mb-3">{plan.description}</p>}
              <div className="mb-4">
                {plan.monthly_price === 0 && plan.yearly_price === 0 ? (
                  <span className="text-2xl font-bold text-gray-900">免费</span>
                ) : (
                  <div>
                    <div className="flex items-baseline gap-1">
                      <span className="text-2xl font-bold text-gray-900">¥{plan.monthly_price}</span>
                      <span className="text-xs text-gray-400">/月</span>
                    </div>
                    <p className="text-xs text-gray-400 mt-0.5">年付 ¥{plan.yearly_price}（省 ¥{plan.monthly_price * 12 - plan.yearly_price}）</p>
                  </div>
                )}
              </div>
              <div className="space-y-2 mb-4 text-sm">
                <div className="flex items-center gap-2 text-gray-600">
                  <Layers className="h-3.5 w-3.5 text-gray-400" />
                  <span>{plan.max_projects === -1 ? '不限' : plan.max_projects} 个项目</span>
                </div>
                <div className="flex items-center gap-2 text-gray-600">
                  <Zap className="h-3.5 w-3.5 text-gray-400" />
                  <span>每月 {plan.max_ai_queries === -1 ? '不限' : plan.max_ai_queries} 次 AI 分析</span>
                </div>
                <div className="flex items-center gap-2 text-gray-600">
                  <CreditCard className="h-3.5 w-3.5 text-gray-400" />
                  <span>赠送 {plan.credits} 次额度</span>
                </div>
              </div>
              <ul className="space-y-1.5 mb-4">{featureList(plan.features)}</ul>
              <button onClick={() => openEdit(plan)}
                className="w-full flex items-center justify-center gap-1.5 px-3 py-2 text-sm font-medium rounded-md border border-gray-200 text-gray-700 hover:bg-gray-50 transition-colors">
                <Edit3 className="h-3.5 w-3.5" />
                编辑
              </button>
            </div>
          ))}
        </div>
      )}

      {showEditor && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40" onClick={() => setShowEditor(false)}>
          <div className="bg-white rounded-xl shadow-xl w-full max-w-lg mx-4 max-h-[90vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-center justify-between p-5 border-b border-gray-200">
              <h3 className="text-base font-semibold text-gray-900">{editing ? '编辑套餐' : '新建套餐'}</h3>
              <button onClick={() => setShowEditor(false)} className="p-1 text-gray-400 hover:text-gray-600"><X className="h-4 w-4" /></button>
            </div>
            <div className="p-5 space-y-4">
              <div>
                <label className={labelClass}>套餐名称 *</label>
                <input type="text" className={inputClass} value={form.name} onChange={(e) => setForm((p) => ({ ...p, name: e.target.value }))} />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className={labelClass}>月付价格 (¥)</label>
                  <input type="number" min={0} className={inputClass} value={form.monthly_price} onChange={(e) => setForm((p) => ({ ...p, monthly_price: parseInt(e.target.value) || 0 }))} />
                </div>
                <div>
                  <label className={labelClass}>年付价格 (¥)</label>
                  <input type="number" min={0} className={inputClass} value={form.yearly_price} onChange={(e) => setForm((p) => ({ ...p, yearly_price: parseInt(e.target.value) || 0 }))} />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className={labelClass}>最大项目数</label>
                  <input type="number" min={-1} className={inputClass} value={form.max_projects} onChange={(e) => setForm((p) => ({ ...p, max_projects: parseInt(e.target.value) || 0 }))} />
                  <p className="text-xs text-gray-400 mt-1">-1 表示不限</p>
                </div>
                <div>
                  <label className={labelClass}>每月 AI 查询数</label>
                  <input type="number" min={-1} className={inputClass} value={form.max_ai_queries} onChange={(e) => setForm((p) => ({ ...p, max_ai_queries: parseInt(e.target.value) || 0 }))} />
                  <p className="text-xs text-gray-400 mt-1">-1 表示不限</p>
                </div>
              </div>
              <div>
                <label className={labelClass}>赠送额度</label>
                <input type="number" min={0} className={inputClass} value={form.credits} onChange={(e) => setForm((p) => ({ ...p, credits: parseInt(e.target.value) || 0 }))} />
              </div>
              <div>
                <label className={labelClass}>功能列表（每行一个）</label>
                <textarea className={`${inputClass} resize-none`} rows={4} value={form.features}
                  onChange={(e) => setForm((p) => ({ ...p, features: e.target.value }))}
                  placeholder="10 个品牌项目&#10;每月 100 次 AI 分析&#10;深度报告" />
              </div>
              <label className="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" checked={form.popular}
                  onChange={(e) => setForm((p) => ({ ...p, popular: e.target.checked }))}
                  className="rounded border-gray-300 text-primary focus:ring-primary" />
                <span className="text-sm text-gray-700">标记为推荐套餐</span>
              </label>
            </div>
            <div className="flex items-center justify-end gap-3 p-5 border-t border-gray-200">
              <button onClick={() => setShowEditor(false)} className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900">取消</button>
              <button onClick={handleSave} disabled={saving || !form.name}
                className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 disabled:opacity-60">
                {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
                {editing ? '保存修改' : '创建套餐'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
