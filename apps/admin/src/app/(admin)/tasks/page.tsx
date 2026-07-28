'use client';

import { useEffect, useState, useCallback } from 'react';
import {
  ClipboardList, ChevronDown, ChevronUp, Loader2, RefreshCw,
  RotateCw, Trash2, Search
} from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface Task {
  id: string;
  project_name: string;
  username: string;
  status: string;
  model: string;
  created_at: string;
  progress: number;
  ai_answers: { question: string; answer: string }[];
}



const STATUS_OPTIONS = ['all', 'running', 'completed', 'failed', 'pending'] as const;
const statusLabels: Record<string, string> = { all: '全部', running: '进行中', completed: '已完成', failed: '失败', pending: '等待中' };
const statusColors: Record<string, string> = {
  running: 'bg-blue-50 text-blue-700', completed: 'bg-green-50 text-green-700',
  failed: 'bg-red-50 text-red-700', pending: 'bg-gray-100 text-gray-600',
};

export default function TasksPage() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [search, setSearch] = useState('');
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    setError('');
    apiRequest<Task[]>('/tasks')
      .then((r) => { if (Array.isArray(r.data)) setTasks(r.data); })
      .catch((e) => setError(e instanceof Error ? e.message : '加载失败'))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => { load(); }, [load]);

  const handleRetry = async (id: string) => {
    try { await apiRequest(`/tasks/${id}/retry`, { method: 'POST' }); load(); }
    catch { load(); }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('确定删除此任务？')) return;
    try { await apiRequest(`/tasks/${id}`, { method: 'DELETE' }); load(); }
    catch { load(); }
  };

  const displayTasks = tasks
    .filter((t) => statusFilter === 'all' || t.status === statusFilter)
    .filter((t) => !search || t.project_name.includes(search) || (t.username ?? '').includes(search) || t.id.includes(search));

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <ClipboardList className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">任务管理</h2>
          <span className="text-sm text-gray-400 font-normal">共 {tasks.length} 个任务</span>
        </div>
        <button onClick={load} disabled={loading}
          className="flex items-center gap-1.5 px-3 py-2 text-sm text-gray-600 hover:text-gray-900">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          刷新
        </button>
      </div>

      <div className="flex items-center gap-3 mb-4 flex-wrap">
        {STATUS_OPTIONS.map((s) => (
          <button key={s} onClick={() => setStatusFilter(s)}
            className={`px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
              statusFilter === s ? 'bg-primary text-white' : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50'
            }`}>
            {statusLabels[s]}
          </button>
        ))}
        <div className="relative ml-auto">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
          <input type="text" placeholder="搜索项目/用户/ID..." value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9 pr-4 py-1.5 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary text-gray-900 placeholder-gray-400 w-56" />
        </div>
      </div>

      {error && <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-sm text-red-700">{error}</div>}

      {loading && tasks.length === 0 ? (
        <div className="flex items-center justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-gray-400" /></div>
      ) : !loading && !error && displayTasks.length === 0 ? (
        <div className="bg-white rounded-lg border border-gray-200 p-12 text-center">
          <ClipboardList className="h-10 w-10 text-gray-300 mx-auto mb-3" />
          <p className="text-gray-400 text-sm">暂无任务</p>
        </div>
      ) : !error && (
        <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200 bg-gray-50 text-left">
                  <th className="px-4 py-3 font-semibold text-gray-600">任务 ID</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">项目</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">用户</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">状态</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">模型</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">创建时间</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">进度</th>
                  <th className="px-4 py-3 font-semibold text-gray-600">操作</th>
                </tr>
              </thead>
              <tbody>
                {displayTasks.map((t) => (
                  <>
                    <tr key={t.id} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                      <td className="px-4 py-3">
                        <button onClick={() => setExpandedId(expandedId === t.id ? null : t.id)}
                          className="flex items-center gap-1 font-mono text-xs text-gray-500 hover:text-gray-700">
                          {expandedId === t.id ? <ChevronUp className="h-3 w-3" /> : <ChevronDown className="h-3 w-3" />}
                          {t.id.slice(0, 8)}...
                        </button>
                      </td>
                      <td className="px-4 py-3 font-medium text-gray-900">{t.project_name}</td>
                      <td className="px-4 py-3 text-gray-600">{t.username ?? '-'}</td>
                      <td className="px-4 py-3">
                        <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium ${statusColors[t.status] || 'bg-gray-100 text-gray-600'}`}>
                          {statusLabels[t.status] || t.status}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-gray-600">{t.model}</td>
                      <td className="px-4 py-3 text-gray-500 text-xs">{t.created_at}</td>
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <div className="w-20 bg-gray-200 rounded-full h-1.5">
                            <div className={`h-1.5 rounded-full ${t.status === 'failed' ? 'bg-red-400' : t.status === 'completed' ? 'bg-green-400' : 'bg-blue-400'}`}
                              style={{ width: `${t.progress ?? 0}%` }} />
                          </div>
                          <span className="text-xs text-gray-500">{t.progress ?? 0}%</span>
                        </div>
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-1">
                          {t.status === 'failed' && (
                            <button onClick={() => handleRetry(t.id)}
                              className="p-1.5 text-gray-400 hover:text-blue-600 transition-colors" title="重试">
                              <RotateCw className="h-3.5 w-3.5" />
                            </button>
                          )}
                          <button onClick={() => handleDelete(t.id)}
                            className="p-1.5 text-gray-400 hover:text-red-500 transition-colors" title="删除">
                            <Trash2 className="h-3.5 w-3.5" />
                          </button>
                        </div>
                      </td>
                    </tr>
                    {expandedId === t.id && (
                      <tr key={`${t.id}-expanded`}>
                        <td colSpan={8} className="px-4 py-3 bg-gray-50 border-b border-gray-100">
                          <div className="text-sm">
                            <p className="text-xs text-gray-400 mb-2 font-medium">AI 分析结果：</p>
                            {(t.ai_answers?.length ?? 0) === 0 ? (
                              <p className="text-xs text-gray-400">暂无分析结果</p>
                            ) : (
                              <div className="space-y-2">
                                {t.ai_answers.map((a, idx) => (
                                  <div key={idx} className="bg-white rounded-md p-3 border border-gray-200">
                                    <p className="text-xs font-medium text-gray-700 mb-1">Q: {a.question}</p>
                                    <p className="text-xs text-gray-600">A: {a.answer}</p>
                                  </div>
                                ))}
                              </div>
                            )}
                          </div>
                        </td>
                      </tr>
                    )}
                  </>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      <div className="mt-3 text-xs text-gray-400">
        共 {displayTasks.length} 条记录{search ? ` — 搜索"${search}"` : ''}
      </div>
    </div>
  );
}
