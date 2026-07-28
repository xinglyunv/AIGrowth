'use client';

import { useEffect, useState, useCallback } from 'react';
import {
  Cpu, Plus, Info, CheckCircle2, XCircle, Loader2,
  Edit3, Trash2, RefreshCw, Search, X, Zap
} from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface AIModel {
  id: string;
  name: string;
  provider: string;
  model: string;
  base_url: string;
  api_key: string;
  enabled: boolean;
  description: string;
  is_system: boolean;
  last_tested_at: string | null;
  last_test_status: string;
}

interface DiscoveredModel {
  id: string;
  object: string;
  owned_by: string;
}

const PROVIDERS = ['OpenAI', 'Anthropic', 'Google', '深度求索', '月之暗面', '阿里云', '百度', '自定义'];

export default function ModelsPage() {
  const [models, setModels] = useState<AIModel[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editing, setEditing] = useState<AIModel | null>(null);
  const [saving, setSaving] = useState(false);
  const [testingId, setTestingId] = useState<string | null>(null);
  const [testMsg, setTestMsg] = useState<{id: string; ok: boolean; msg: string} | null>(null);

  const [form, setForm] = useState({ name: '', provider: '', model: '', base_url: '', api_key: '', enabled: true, description: '' });
  const [discovering, setDiscovering] = useState(false);
  const [discovered, setDiscovered] = useState<DiscoveredModel[]>([]);
  const [discoverError, setDiscoverError] = useState('');

  const load = useCallback(() => {
    setLoading(true);
    apiRequest<AIModel[]>('/models')
      .then((r) => { if (Array.isArray(r.data)) setModels(r.data); })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => { load(); }, [load]);

  const openCreate = () => {
    setEditing(null);
    setForm({ name: '', provider: '', model: '', base_url: '', api_key: '', enabled: true, description: '' });
    setDiscovered([]);
    setShowForm(true);
  };

  const openEdit = (m: AIModel) => {
    setEditing(m);
    setForm({ name: m.name, provider: m.provider, model: m.model, base_url: m.base_url, api_key: m.api_key || '', enabled: m.enabled, description: m.description });
    setDiscovered([]);
    setShowForm(true);
  };

  const handleDiscover = async () => {
    if (!form.base_url) return;
    setDiscovering(true);
    setDiscovered([]);
    setDiscoverError('');
    try {
      const body: Record<string, string> = { base_url: form.base_url, api_key: form.api_key };
      // Pass model_id so backend can look up the real key when key is masked
      if (editing && form.api_key.includes('****')) {
        body.model_id = editing.id;
      }
      const res = await apiRequest<DiscoveredModel[]>('/models/discover', {
        method: 'POST',
        body: JSON.stringify(body),
      });
      if (Array.isArray(res.data)) {
        setDiscovered(res.data);
        if (res.data.length === 0) setDiscoverError('接口返回了空模型列表，请检查 API Base URL 和权限。');
      }
    } catch (error) {
      setDiscovered([]);
      setDiscoverError(error instanceof Error ? error.message : '自动获取模型失败，请检查 API Base URL 和 API Key。');
    }
    finally { setDiscovering(false); }
  };

  const pickDiscovered = (d: DiscoveredModel) => {
    setForm((prev) => ({ ...prev, model: d.id }));
  };

  const handleSave = async () => {
    if (!form.name || !form.provider || !form.model || !form.base_url) return;
    setSaving(true);
    try {
      const body: Record<string, unknown> = { ...form };
      if (editing) {
        // Don't send empty api_key on edit (keeps existing key)
        if (!body.api_key) delete body.api_key;
        await apiRequest(`/models/${editing.id}`, { method: 'PUT', body: JSON.stringify(body) });
      } else {
        await apiRequest('/models', { method: 'POST', body: JSON.stringify(body) });
      }
      setShowForm(false);
      load();
    } catch { alert('保存失败'); }
    finally { setSaving(false); }
  };

  const handleDelete = async (m: AIModel) => {
    if (!confirm(`确定删除模型 "${m.name}"？`)) return;
    try {
      await apiRequest(`/models/${m.id}`, { method: 'DELETE' });
      load();
    } catch { alert('删除失败'); }
  };

  const handleTest = async (m: AIModel) => {
    setTestingId(m.id);
    setTestMsg(null);
    try {
      // Don't send masked key; backend will use the stored key
      const body: Record<string, string> = {};
      if (m.api_key && !m.api_key.includes('****')) {
        body.api_key = m.api_key;
      }
      const res = await apiRequest<{success: boolean; status: string; message: string}>(`/models/${m.id}/test`, {
        method: 'POST',
        body: JSON.stringify(body),
      });
      setTestMsg({ id: m.id, ok: res.data.success, msg: res.data.message });
      load();
    } catch { setTestMsg({ id: m.id, ok: false, msg: '测试请求失败' }); }
    finally { setTestingId(null); }
  };

  const inputClass = 'w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary text-gray-900';
  const labelClass = 'block text-sm font-medium text-gray-700 mb-1.5';

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Cpu className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">AI 模型管理</h2>
        </div>
        <button onClick={openCreate}
          className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 transition-colors">
          <Plus className="h-4 w-4" />
          添加模型
        </button>
      </div>

      <div className="mb-4 px-4 py-3 bg-blue-50 border border-blue-200 rounded-md text-sm text-blue-700 flex items-start gap-2">
        <Info className="h-4 w-4 mt-0.5 shrink-0" />
        <span>在此管理可用的 AI 分析模型。添加模型时输入 API Base URL 和 API Key 后点击"自动获取"，系统会自动拉取可用模型列表。</span>
      </div>

      {/* Form Modal */}
      {showForm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40" onClick={() => setShowForm(false)}>
          <div className="bg-white rounded-xl shadow-xl w-full max-w-lg mx-4 max-h-[90vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-center justify-between p-5 border-b border-gray-200">
              <h3 className="text-base font-semibold text-gray-900">{editing ? '编辑模型' : '添加模型'}</h3>
              <button onClick={() => setShowForm(false)} className="p-1 text-gray-400 hover:text-gray-600"><X className="h-4 w-4" /></button>
            </div>
            <div className="p-5 space-y-4">
              <div>
                <label className={labelClass}>模型名称 *</label>
                <input type="text" className={inputClass} placeholder="如 GPT-4o" value={form.name}
                  onChange={(e) => setForm((p) => ({ ...p, name: e.target.value }))} />
              </div>
              <div>
                <label className={labelClass}>供应商 *</label>
                <select className={inputClass} value={form.provider}
                  onChange={(e) => setForm((p) => ({ ...p, provider: e.target.value }))}>
                  <option value="">选择供应商</option>
                  {PROVIDERS.map((p) => <option key={p} value={p}>{p}</option>)}
                </select>
              </div>
              <div>
                <label className={labelClass}>API Base URL *</label>
                <input type="text" className={inputClass} placeholder="如 https://api.openai.com/v1" value={form.base_url}
                  onChange={(e) => setForm((p) => ({ ...p, base_url: e.target.value, model: '' }))} />
              </div>
              <div>
                <label className={labelClass}>API Key</label>
                <div className="flex gap-2">
                  <input type="password" className={`${inputClass} flex-1`} placeholder={editing ? '留空则保持原有 Key' : '输入 API Key 以自动发现模型'} value={form.api_key}
                    onChange={(e) => setForm((p) => ({ ...p, api_key: e.target.value }))} />
                  <button onClick={handleDiscover} disabled={discovering || !form.base_url || (!form.api_key && !editing)}
                    className="flex items-center gap-1 px-3 py-2 text-sm font-medium rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200 disabled:opacity-50 shrink-0">
                    {discovering ? <Loader2 className="h-4 w-4 animate-spin" /> : <Search className="h-4 w-4" />}
                    自动获取
                  </button>
                </div>
              </div>

              {discovered.length > 0 && (
                <div>
                  <label className={labelClass}>可用模型（点击选择）</label>
                  <div className="max-h-40 overflow-y-auto border border-gray-200 rounded-md divide-y divide-gray-100">
                    {discovered.map((d) => (
                      <button key={d.id} type="button"
                        className={`w-full text-left px-3 py-2 text-sm hover:bg-gray-50 transition flex items-center gap-2 ${form.model === d.id ? 'bg-primary/5 text-primary font-medium' : 'text-gray-700'}`}
                        onClick={() => pickDiscovered(d)}>
                        <Cpu className="h-3.5 w-3.5 shrink-0 text-gray-400" />
                        <span>{d.id}</span>
                        {d.owned_by && <span className="text-xs text-gray-400 ml-auto">{d.owned_by}</span>}
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {discoverError && (
                <div className="px-3 py-2 rounded-md bg-red-50 border border-red-200 text-sm text-red-700">
                  {discoverError}
                </div>
              )}

              <div>
                <label className={labelClass}>模型 ID *</label>
                <input type="text" className={inputClass} placeholder="如 gpt-4o" value={form.model}
                  onChange={(e) => setForm((p) => ({ ...p, model: e.target.value }))} />
              </div>
              <div>
                <label className={labelClass}>描述</label>
                <textarea className={`${inputClass} resize-none`} rows={2} placeholder="描述模型能力特点" value={form.description}
                  onChange={(e) => setForm((p) => ({ ...p, description: e.target.value }))} />
              </div>
              <label className="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" checked={form.enabled}
                  onChange={(e) => setForm((p) => ({ ...p, enabled: e.target.checked }))}
                  className="rounded border-gray-300 text-primary focus:ring-primary" />
                <span className="text-sm text-gray-700">创建后立即启用</span>
              </label>
            </div>
            <div className="flex items-center justify-end gap-3 p-5 border-t border-gray-200">
              <button onClick={() => setShowForm(false)} className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900">
                取消
              </button>
              <button onClick={handleSave} disabled={saving || !form.name || !form.provider || !form.model || !form.base_url}
                className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 disabled:opacity-60">
                {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
                {editing ? '保存修改' : '添加模型'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Model List */}
      {loading ? (
        <div className="flex items-center justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-gray-400" /></div>
      ) : models.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-16 text-gray-400">
          <Cpu className="h-12 w-12 mb-3" />
          <p>暂无模型配置</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {models.map((m) => (
            <div key={m.id}
              className="bg-white rounded-lg border border-gray-200 p-5 hover:shadow-md hover:border-gray-300 transition-all duration-200">
              <div className="flex items-start justify-between mb-3">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <h3 className="text-base font-semibold text-gray-900 truncate">{m.name}</h3>
                    {m.is_system && <span className="px-1.5 py-0.5 bg-gray-100 text-gray-500 text-[10px] font-medium rounded">系统</span>}
                  </div>
                  <p className="text-xs text-gray-400 mt-0.5">{m.provider} · {m.model}</p>
                </div>
                <div className="flex items-center gap-1 shrink-0 ml-2">
                  {m.last_test_status === 'success' && <CheckCircle2 className="h-3.5 w-3.5 text-green-500" />}
                  {m.last_test_status === 'failed' && <XCircle className="h-3.5 w-3.5 text-red-500" />}
                  <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium ${
                    m.enabled ? 'bg-green-50 text-green-700' : 'bg-gray-100 text-gray-500'
                  }`}>
                    {m.enabled ? '已启用' : '已禁用'}
                  </span>
                </div>
              </div>
              {m.description && <p className="text-sm text-gray-600 mb-3 leading-relaxed">{m.description}</p>}
              <div className="bg-gray-50 rounded-md px-3 py-2 text-xs font-mono text-gray-500 break-all mb-3 flex items-center gap-2">
                <span>{m.base_url}</span>
                <span className="text-gray-300">/</span>
                <span className="text-primary font-semibold">{m.model}</span>
              </div>

              {testMsg && testMsg.id === m.id && (
                <div className={`mb-3 px-3 py-2 rounded-md text-xs ${
                  testMsg.ok ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'
                }`}>{testMsg.msg}</div>
              )}

              <div className="flex items-center justify-between text-xs text-gray-400 pt-2 border-t border-gray-100">
                <span>
                  上次检测：<span className="text-gray-500">{m.last_tested_at ? new Date(m.last_tested_at).toLocaleString('zh-CN') : '未检测'}</span>
                  {m.last_tested_at && m.last_test_status === 'success' && <CheckCircle2 className="h-3 w-3 inline ml-1 text-green-500" />}
                  {m.last_tested_at && m.last_test_status === 'failed' && <XCircle className="h-3 w-3 inline ml-1 text-red-500" />}
                </span>
                <div className="flex items-center gap-2">
                  <button onClick={() => handleTest(m)} disabled={testingId === m.id}
                    className="flex items-center gap-1 text-primary hover:text-primary/80 font-medium disabled:opacity-50">
                    {testingId === m.id ? <Loader2 className="h-3 w-3 animate-spin" /> : <Zap className="h-3 w-3" />}
                    测试连接
                  </button>
                  <button onClick={() => openEdit(m)}
                    className="p-1 text-gray-400 hover:text-gray-600"><Edit3 className="h-3.5 w-3.5" /></button>
                  <button onClick={() => handleDelete(m)}
                    className="p-1 text-gray-400 hover:text-red-500"><Trash2 className="h-3.5 w-3.5" /></button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
