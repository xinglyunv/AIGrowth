'use client';

import { useEffect, useState, useCallback } from 'react';
import {
  FileText, Plus, Loader2, Edit3, CheckCircle2, XCircle,
  RefreshCw, X, Send, Save
} from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface Prompt {
  id: string;
  name: string;
  purpose: string;
  version: string;
  status: 'draft' | 'published';
  content: string;
  updated_at: string;
}



const STATUS_OPTIONS = ['all', 'published', 'draft'];
const statusLabels: Record<string, string> = { all: '全部', published: '已发布', draft: '草稿' };

export default function PromptsPage() {
  const [prompts, setPrompts] = useState<Prompt[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [showEditor, setShowEditor] = useState(false);
  const [editing, setEditing] = useState<Prompt | null>(null);
  const [saving, setSaving] = useState(false);
  const [form, setForm] = useState({ name: '', purpose: '', content: '' });

  const load = useCallback(() => {
    setLoading(true);
    setError('');
    apiRequest<Prompt[]>('/prompts')
      .then((r) => { if (Array.isArray(r.data)) setPrompts(r.data); })
      .catch((e) => setError(e instanceof Error ? e.message : '加载失败'))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => { load(); }, [load]);

  const openCreate = () => {
    setEditing(null);
    setForm({ name: '', purpose: '', content: '' });
    setShowEditor(true);
  };

  const openEdit = (p: Prompt) => {
    setEditing(p);
    setForm({ name: p.name, purpose: p.purpose, content: p.content });
    setShowEditor(true);
  };

  const handleSave = async () => {
    if (!form.name || !form.content) return;
    setSaving(true);
    try {
      if (editing) {
        await apiRequest(`/prompts/${editing.id}`, { method: 'PUT', body: JSON.stringify(form) });
      } else {
        await apiRequest('/prompts', { method: 'POST', body: JSON.stringify(form) });
      }
      setShowEditor(false);
      load();
    } catch { load(); }
    setSaving(false);
  };

  const handlePublish = async (p: Prompt) => {
    try {
      await apiRequest(`/prompts/${p.id}/publish`, { method: 'POST' });
      load();
    } catch { load(); }
  };

  const inputClass = 'w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary text-gray-900';
  const labelClass = 'block text-sm font-medium text-gray-700 mb-1.5';

  const displayPrompts = prompts
    .filter((p) => statusFilter === 'all' || p.status === statusFilter);

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <FileText className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">提示词管理</h2>
          <span className="text-sm text-gray-400 font-normal">共 {prompts.length} 个模板</span>
        </div>
        <div className="flex items-center gap-2">
          <button onClick={load} disabled={loading}
            className="flex items-center gap-1.5 px-3 py-2 text-sm text-gray-600 hover:text-gray-900">
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
            刷新
          </button>
          <button onClick={openCreate}
            className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90">
            <Plus className="h-4 w-4" />
            新建模板
          </button>
        </div>
      </div>

      <div className="flex items-center gap-3 mb-4">
        {STATUS_OPTIONS.map((s) => (
          <button key={s} onClick={() => setStatusFilter(s)}
            className={`px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
              statusFilter === s ? 'bg-primary text-white' : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50'
            }`}>
            {statusLabels[s]}
          </button>
        ))}
      </div>

      {error && <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-sm text-red-700">{error}</div>}

      {loading && prompts.length === 0 ? (
        <div className="flex items-center justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-gray-400" /></div>
      ) : !loading && !error && displayPrompts.length === 0 ? (
        <div className="bg-white rounded-lg border border-gray-200 p-12 text-center">
          <FileText className="h-10 w-10 text-gray-300 mx-auto mb-3" />
          <p className="text-gray-400 text-sm">暂无提示词模板</p>
        </div>
      ) : !error && (
        <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200 bg-gray-50 text-left">
                  <th className="px-4 py-3 font-semibold text-gray-600">名称</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">用途</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">版本</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">状态</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">最后更新</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">操作</th>
                </tr>
              </thead>
              <tbody>
                {displayPrompts.map((p) => (
                  <tr key={p.id} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                    <td className="px-4 py-3 font-medium text-gray-900">{p.name}</td>
                    <td className="px-4 py-3 text-gray-600 max-w-xs truncate">{p.purpose}</td>
                    <td className="px-4 py-3">
                      <span className="font-mono text-xs text-gray-500 bg-gray-100 px-1.5 py-0.5 rounded">{p.version}</span>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium ${
                        p.status === 'published' ? 'bg-green-50 text-green-700' : 'bg-yellow-50 text-yellow-700'
                      }`}>
                        {p.status === 'published' ? <CheckCircle2 className="h-3 w-3" /> : <XCircle className="h-3 w-3" />}
                        {p.status === 'published' ? '已发布' : '草稿'}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-gray-500 text-xs">{p.updated_at}</td>
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-1">
                        <button onClick={() => openEdit(p)}
                          className="p-1.5 text-gray-400 hover:text-gray-600 transition-colors" title="编辑">
                          <Edit3 className="h-3.5 w-3.5" />
                        </button>
                        {p.status === 'draft' && (
                          <button onClick={() => handlePublish(p)}
                            className="p-1.5 text-gray-400 hover:text-green-600 transition-colors" title="发布">
                            <Send className="h-3.5 w-3.5" />
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      <div className="mt-3 text-xs text-gray-400">
        共 {displayPrompts.length} 条记录
      </div>

      {showEditor && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40" onClick={() => setShowEditor(false)}>
          <div className="bg-white rounded-xl shadow-xl w-full max-w-2xl mx-4 max-h-[90vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-center justify-between p-5 border-b border-gray-200">
              <h3 className="text-base font-semibold text-gray-900">{editing ? '编辑提示词' : '新建提示词模板'}</h3>
              <button onClick={() => setShowEditor(false)} className="p-1 text-gray-400 hover:text-gray-600"><X className="h-4 w-4" /></button>
            </div>
            <div className="p-5 space-y-4">
              <div>
                <label className={labelClass}>模板名称 *</label>
                <input type="text" className={inputClass} placeholder="如：品牌分析主提示词" value={form.name}
                  onChange={(e) => setForm((p) => ({ ...p, name: e.target.value }))} />
              </div>
              <div>
                <label className={labelClass}>用途说明</label>
                <input type="text" className={inputClass} placeholder="简述该模板的使用场景" value={form.purpose}
                  onChange={(e) => setForm((p) => ({ ...p, purpose: e.target.value }))} />
              </div>
              <div>
                <label className={labelClass}>提示词内容 *</label>
                <textarea className={`${inputClass} font-mono text-xs resize-none`} rows={12}
                  placeholder="输入提示词模板内容，使用 {variable} 作为变量占位符"
                  value={form.content} onChange={(e) => setForm((p) => ({ ...p, content: e.target.value }))} />
              </div>
              <div className="p-3 bg-blue-50 rounded-md text-xs text-blue-700">
                使用 {'{变量名}'} 作为动态变量，系统会在运行时替换为实际值。例如：{'{brand_name}'}、{'{website_url}'}
              </div>
            </div>
            <div className="flex items-center justify-end gap-3 p-5 border-t border-gray-200">
              <button onClick={() => setShowEditor(false)} className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900">取消</button>
              <button onClick={handleSave} disabled={saving || !form.name || !form.content}
                className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 disabled:opacity-60">
                {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
                {editing ? '保存修改' : '创建模板'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
