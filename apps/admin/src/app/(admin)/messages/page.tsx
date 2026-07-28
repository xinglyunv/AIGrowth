'use client';

import { useEffect, useState, useCallback } from 'react';
import { MessageSquare, Mail, Trash2, RefreshCw, ChevronDown, ChevronUp } from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface Message {
  id: number;
  name: string;
  email: string;
  subject: string;
  message: string;
  is_read: boolean;
  created_at: string;
}

interface MessagesData {
  messages: Message[];
  unread: number;
}

export default function MessagesPage() {
  const [data, setData] = useState<MessagesData>({ messages: [], unread: 0 });
  const [loading, setLoading] = useState(true);
  const [expandedId, setExpandedId] = useState<number | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    apiRequest<MessagesData>('/contact')
      .then((r) => { if (r.data) setData(r.data); })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => { load(); }, [load]);

  const handleMarkRead = async (id: number) => {
    try {
      await apiRequest(`/contact/${id}/read`, { method: 'PUT' });
      setData((prev) => ({
        ...prev,
        messages: prev.messages.map((m) => m.id === id ? { ...m, is_read: true } : m),
        unread: Math.max(0, prev.unread - 1),
      }));
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这条留言吗？')) return;
    try {
      await apiRequest(`/contact/${id}`, { method: 'DELETE' });
      const msg = data.messages.find((m) => m.id === id);
      setData((prev) => ({
        ...prev,
        messages: prev.messages.filter((m) => m.id !== id),
        unread: msg && !msg.is_read ? prev.unread - 1 : prev.unread,
      }));
      if (expandedId === id) setExpandedId(null);
    } catch {}
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <MessageSquare className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">留言管理</h2>
          {data.unread > 0 && (
            <span className="px-2 py-0.5 bg-red-50 text-red-600 text-xs font-medium rounded-full">
              {data.unread} 条未读
            </span>
          )}
        </div>
        <button onClick={load} disabled={loading}
          className="flex items-center gap-1.5 px-3 py-2 text-sm text-gray-600 hover:text-gray-900 transition-colors">
          <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          刷新
        </button>
      </div>

      {loading && data.messages.length === 0 ? (
        <div className="flex items-center justify-center py-20"><p className="text-gray-400">加载中...</p></div>
      ) : data.messages.length === 0 ? (
        <div className="bg-white rounded-lg border border-gray-200 p-12 text-center">
          <Mail className="h-10 w-10 text-gray-300 mx-auto mb-3" />
          <p className="text-gray-400 text-sm">暂无留言</p>
          <p className="text-gray-300 text-xs mt-1">来自官网联系表单的留言将显示在这里</p>
        </div>
      ) : (
        <div className="space-y-3">
          {data.messages.map((msg) => (
            <div key={msg.id} className={`bg-white rounded-lg border ${msg.is_read ? 'border-gray-200' : 'border-primary/30 bg-primary/5'} p-5`}>
              <div className="flex items-start justify-between gap-4">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-3 mb-2">
                    <span className="text-sm font-semibold text-gray-900">{msg.name}</span>
                    <span className="text-xs text-gray-400">{msg.email}</span>
                    {!msg.is_read && <span className="px-1.5 py-0.5 bg-primary/10 text-primary text-[10px] font-medium rounded">未读</span>}
                  </div>
                  <p className="text-sm text-gray-700 font-medium truncate">{msg.subject || '(无主题)'}</p>
                  <p className="text-xs text-gray-400 mt-1">{new Date(msg.created_at).toLocaleString('zh-CN')}</p>
                  {expandedId === msg.id && (
                    <div className="mt-3 p-3 bg-gray-50 rounded-md">
                      <p className="text-sm text-gray-700 whitespace-pre-wrap">{msg.message}</p>
                    </div>
                  )}
                </div>
                <div className="flex items-center gap-2 shrink-0">
                  <button
                    onClick={() => { setExpandedId(expandedId === msg.id ? null : msg.id); if (!msg.is_read) handleMarkRead(msg.id); }}
                    className="p-1.5 text-gray-400 hover:text-gray-600 transition-colors"
                    title={expandedId === msg.id ? '收起' : '展开'}>
                    {expandedId === msg.id ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                  </button>
                  <button onClick={() => handleDelete(msg.id)}
                    className="p-1.5 text-gray-400 hover:text-red-500 transition-colors" title="删除">
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
