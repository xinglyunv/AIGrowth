'use client';

import { useEffect, useState, useCallback } from 'react';
import { Key, Plus, X, Loader2, Save, Eye, EyeOff, ChevronLeft, ChevronRight } from 'lucide-react';
import { apiRequest } from '@/lib/api';
import { toDateTimeLocal, toRFC3339 } from '@/lib/cdk';

interface CDKCode {
  id: string;
  code: string;
  credits: number;
  max_uses: number;
  used_count: number;
  is_active: boolean;
  expires_at: string | null;
  created_at: string;
  updated_at: string;
}

interface CDKUsage {
  id: string;
  user_id: string;
  username: string;
  used_at: string;
}

interface PaginatedData<T> {
  data: T[];
  total: number;
}

const PAGE_SIZE = 20;

export default function CDKPage() {
  const [codes, setCodes] = useState<CDKCode[]>([]);
  const [total, setTotal] = useState(0);
  const [offset, setOffset] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreate, setShowCreate] = useState(false);
  const [showEdit, setShowEdit] = useState<CDKCode | null>(null);
  const [showUsages, setShowUsages] = useState<string | null>(null);
  const [usages, setUsages] = useState<CDKUsage[]>([]);
  const [saving, setSaving] = useState(false);
  const [createForm, setCreateForm] = useState({ code: '', credits: 0, max_uses: 0, expires_at: '' });
  const [editForm, setEditForm] = useState({ max_uses: 0, is_active: true, expires_at: '' });

  const fetchCodes = useCallback(() => {
    setLoading(true);
    setError('');
    apiRequest<PaginatedData<CDKCode>>(`/cdk?offset=${offset}&limit=${PAGE_SIZE}`)
      .then((res) => {
        if (res.data) {
          setCodes(res.data.data ?? []);
          setTotal(res.data.total ?? 0);
        }
      })
      .catch((e) => setError(e instanceof Error ? e.message : '加载失败'))
      .finally(() => setLoading(false));
  }, [offset]);

  useEffect(() => { fetchCodes(); }, [fetchCodes]);

  const totalPages = Math.ceil(total / PAGE_SIZE);
  const currentPage = Math.floor(offset / PAGE_SIZE) + 1;

  const handleCreate = () => {
    setSaving(true);
    apiRequest('/cdk', {
      method: 'POST',
      body: JSON.stringify({
        code: createForm.code.trim(),
        credits: createForm.credits,
        max_uses: createForm.max_uses,
        expires_at: toRFC3339(createForm.expires_at),
      }),
    })
      .then(() => {
        setShowCreate(false);
        setCreateForm({ code: '', credits: 0, max_uses: 0, expires_at: '' });
        fetchCodes();
      })
      .catch((e) => setError(e instanceof Error ? e.message : '创建失败'))
      .finally(() => setSaving(false));
  };

  const handleEdit = () => {
    if (!showEdit) return;
    setSaving(true);
    apiRequest(`/cdk/${showEdit.id}`, {
      method: 'PUT',
      body: JSON.stringify({
        max_uses: editForm.max_uses,
        is_active: editForm.is_active,
        expires_at: toRFC3339(editForm.expires_at),
      }),
    })
      .then(() => {
        setShowEdit(null);
        fetchCodes();
      })
      .catch((e) => setError(e instanceof Error ? e.message : '保存失败'))
      .finally(() => setSaving(false));
  };

  const toggleUsages = (id: string) => {
    if (showUsages === id) {
      setShowUsages(null);
      setUsages([]);
      return;
    }
    setShowUsages(id);
    setUsages([]);
    apiRequest<CDKUsage[]>(`/cdk/${id}/usages`)
      .then((res) => { if (Array.isArray(res.data)) setUsages(res.data); })
      .catch(() => setUsages([]));
  };

  const openEditModal = (code: CDKCode) => {
    setShowEdit(code);
    setEditForm({
      max_uses: code.max_uses,
      is_active: code.is_active,
       expires_at: toDateTimeLocal(code.expires_at),
    });
  };

  const inputClass = 'w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary text-gray-900';
  const labelClass = 'block text-sm font-medium text-gray-700 mb-1.5';

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Key className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">CDK 管理</h2>
        </div>
        <button onClick={() => setShowCreate(true)}
          className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90">
          <Plus className="h-4 w-4" />
          新建 CDK
        </button>
      </div>

      {error && <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-sm text-red-700">{error}</div>}

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 bg-gray-50 text-left">
                <th className="px-5 py-3 font-semibold text-gray-600">CDK 代码</th>
                <th className="px-5 py-3 font-semibold text-gray-600">额度</th>
                <th className="px-5 py-3 font-semibold text-gray-600">使用次数</th>
                <th className="px-5 py-3 font-semibold text-gray-600">状态</th>
                <th className="px-5 py-3 font-semibold text-gray-600">过期时间</th>
                <th className="px-5 py-3 font-semibold text-gray-600">创建时间</th>
                <th className="px-5 py-3 font-semibold text-gray-600">操作</th>
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={7} className="px-5 py-12 text-center text-gray-400">
                    <Loader2 className="animate-spin inline h-5 w-5 mr-2" />加载中...
                  </td>
                </tr>
              ) : codes.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-5 py-12 text-center text-gray-400">暂无 CDK</td>
                </tr>
              ) : (
                codes.map((item) => (
                  <>
                    <tr key={item.id} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                      <td className="px-5 py-3 font-mono text-sm text-gray-900">{item.code}</td>
                      <td className="px-5 py-3 text-gray-700">{item.credits}</td>
                      <td className="px-5 py-3 text-gray-700">{item.used_count} / {item.max_uses === 0 ? '不限' : item.max_uses}</td>
                      <td className="px-5 py-3">
                        <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${item.is_active ? 'bg-green-50 text-green-700' : 'bg-gray-100 text-gray-500'}`}>
                          {item.is_active ? '启用' : '禁用'}
                        </span>
                      </td>
                      <td className="px-5 py-3 text-gray-500 text-xs">
                        {item.expires_at ? new Date(item.expires_at).toLocaleDateString('zh-CN') : '永不过期'}
                      </td>
                      <td className="px-5 py-3 text-gray-500 text-xs">
                        {new Date(item.created_at).toLocaleDateString('zh-CN')}
                      </td>
                      <td className="px-5 py-3">
                        <div className="flex items-center gap-2">
                          <button onClick={() => openEditModal(item)}
                            className="px-2 py-1 text-xs font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded transition-colors">
                            编辑
                          </button>
                          <button onClick={() => toggleUsages(item.id)}
                            className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded transition-colors">
                            {showUsages === item.id ? <EyeOff className="h-3 w-3" /> : <Eye className="h-3 w-3" />}
                            使用记录
                          </button>
                        </div>
                      </td>
                    </tr>
                    {showUsages === item.id && (
                      <tr key={`${item.id}-usages`}>
                        <td colSpan={7} className="px-5 py-3 bg-gray-50">
                          {usages.length === 0 ? (
                            <p className="text-xs text-gray-400 text-center py-2">暂无使用记录</p>
                          ) : (
                            <table className="w-full text-xs">
                              <thead>
                                <tr className="text-left text-gray-500 border-b border-gray-200">
                                  <th className="px-3 py-1.5 font-medium">用户</th>
                                  <th className="px-3 py-1.5 font-medium">使用时间</th>
                                </tr>
                              </thead>
                              <tbody>
                                {usages.map((u) => (
                                  <tr key={u.id} className="border-b border-gray-100">
                                    <td className="px-3 py-1.5 text-gray-700">{u.username}</td>
                                    <td className="px-3 py-1.5 text-gray-500">{new Date(u.used_at).toLocaleString('zh-CN')}</td>
                                  </tr>
                                ))}
                              </tbody>
                            </table>
                          )}
                        </td>
                      </tr>
                    )}
                  </>
                ))
              )}
            </tbody>
          </table>
        </div>
        {totalPages > 1 && (
          <div className="flex items-center justify-between px-5 py-3 border-t border-gray-200 bg-gray-50">
            <span className="text-xs text-gray-500">共 {total} 条，第 {currentPage}/{totalPages} 页</span>
            <div className="flex items-center gap-2">
              <button onClick={() => setOffset((p) => Math.max(0, p - PAGE_SIZE))} disabled={offset === 0}
                className="p-1 rounded text-gray-600 hover:bg-gray-200 disabled:opacity-40 disabled:cursor-not-allowed">
                <ChevronLeft className="h-4 w-4" />
              </button>
              <button onClick={() => setOffset((p) => p + PAGE_SIZE)} disabled={offset + PAGE_SIZE >= total}
                className="p-1 rounded text-gray-600 hover:bg-gray-200 disabled:opacity-40 disabled:cursor-not-allowed">
                <ChevronRight className="h-4 w-4" />
              </button>
            </div>
          </div>
        )}
      </div>

      {showCreate && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40" onClick={() => setShowCreate(false)}>
          <div className="bg-white rounded-xl shadow-xl w-full max-w-md mx-4" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-center justify-between p-5 border-b border-gray-200">
              <h3 className="text-base font-semibold text-gray-900">新建 CDK</h3>
              <button onClick={() => setShowCreate(false)} className="p-1 text-gray-400 hover:text-gray-600"><X className="h-4 w-4" /></button>
            </div>
            <div className="p-5 space-y-4">
              <div>
                 <label className={labelClass}>CDK 代码（可选）</label>
                <input type="text" className={inputClass} value={createForm.code}
                  onChange={(e) => setCreateForm((p) => ({ ...p, code: e.target.value }))}
                   placeholder="留空自动生成" />
              </div>
              <div>
                <label className={labelClass}>额度</label>
                <input type="number" min={0} className={inputClass} value={createForm.credits}
                  onChange={(e) => setCreateForm((p) => ({ ...p, credits: parseInt(e.target.value) || 0 }))} />
              </div>
              <div>
                <label className={labelClass}>最大使用次数</label>
                <input type="number" min={0} className={inputClass} value={createForm.max_uses}
                  onChange={(e) => setCreateForm((p) => ({ ...p, max_uses: parseInt(e.target.value) || 0 }))} />
                <p className="text-xs text-gray-400 mt-1">0 表示不限</p>
              </div>
              <div>
                <label className={labelClass}>过期时间</label>
                <input type="datetime-local" className={inputClass} value={createForm.expires_at}
                  onChange={(e) => setCreateForm((p) => ({ ...p, expires_at: e.target.value }))} />
                <p className="text-xs text-gray-400 mt-1">留空表示永不过期</p>
              </div>
            </div>
            <div className="flex items-center justify-end gap-3 p-5 border-t border-gray-200">
              <button onClick={() => setShowCreate(false)} className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900">取消</button>
               <button onClick={handleCreate} disabled={saving}
                className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 disabled:opacity-60">
                {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
                创建
              </button>
            </div>
          </div>
        </div>
      )}

      {showEdit && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40" onClick={() => setShowEdit(null)}>
          <div className="bg-white rounded-xl shadow-xl w-full max-w-md mx-4" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-center justify-between p-5 border-b border-gray-200">
              <h3 className="text-base font-semibold text-gray-900">编辑 CDK</h3>
              <button onClick={() => setShowEdit(null)} className="p-1 text-gray-400 hover:text-gray-600"><X className="h-4 w-4" /></button>
            </div>
            <div className="p-5 space-y-4">
              <div className="bg-gray-50 rounded-md p-3">
                <p className="text-xs text-gray-500">CDK 代码</p>
                <p className="text-sm font-mono font-medium text-gray-900">{showEdit.code}</p>
              </div>
              <div className="bg-gray-50 rounded-md p-3">
                <p className="text-xs text-gray-500">额度</p>
                <p className="text-sm font-medium text-gray-900">{showEdit.credits}</p>
              </div>
              <div>
                <label className={labelClass}>最大使用次数</label>
                <input type="number" min={0} className={inputClass} value={editForm.max_uses}
                  onChange={(e) => setEditForm((p) => ({ ...p, max_uses: parseInt(e.target.value) || 0 }))} />
                <p className="text-xs text-gray-400 mt-1">0 表示不限</p>
              </div>
              <div>
                <label className={labelClass}>过期时间</label>
                <input type="datetime-local" className={inputClass} value={editForm.expires_at}
                  onChange={(e) => setEditForm((p) => ({ ...p, expires_at: e.target.value }))} />
                <p className="text-xs text-gray-400 mt-1">留空表示永不过期</p>
              </div>
              <label className="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" checked={editForm.is_active}
                  onChange={(e) => setEditForm((p) => ({ ...p, is_active: e.target.checked }))}
                  className="rounded border-gray-300 text-primary focus:ring-primary" />
                <span className="text-sm text-gray-700">启用</span>
              </label>
            </div>
            <div className="flex items-center justify-end gap-3 p-5 border-t border-gray-200">
              <button onClick={() => setShowEdit(null)} className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900">取消</button>
              <button onClick={handleEdit} disabled={saving}
                className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 disabled:opacity-60">
                {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
                保存修改
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
