'use client';
import { useEffect, useState } from 'react';
import { apiRequest } from '@/lib/api';
import { Bell, CheckCheck, ChevronRight, RefreshCw, Mail, MailOpen, Info, AlertTriangle, CheckCircle2, ExternalLink } from 'lucide-react';

interface Notification {
  id: string;
  title: string;
  content: string;
  type: string;
  read: boolean;
  link?: string;
  created_at: string;
}

const TYPE_CONFIG: Record<string, { label: string; icon: typeof Bell; color: string; bg: string }> = {
  info: { label: '信息', icon: Info, color: 'text-blue-600', bg: 'bg-blue-50' },
  warning: { label: '警告', icon: AlertTriangle, color: 'text-amber-600', bg: 'bg-amber-50' },
  success: { label: '成功', icon: CheckCircle2, color: 'text-green-600', bg: 'bg-green-50' },
  task_completed: { label: '任务完成', icon: CheckCircle2, color: 'text-green-600', bg: 'bg-green-50' },
};

function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const d = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - d.getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 60) return `${mins} 分钟前`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours} 小时前`;
  const days = Math.floor(hours / 24);
  if (days < 7) return `${days} 天前`;
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' });
}

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [markingAll, setMarkingAll] = useState(false);

  const fetchNotifications = () => {
    setLoading(true);
    apiRequest<Notification[]>('/notifications')
      .then((res) => setNotifications(Array.isArray(res.data) ? res.data : []))
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : '加载通知失败';
        setError(message);
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchNotifications();
  }, []);

  const handleMarkRead = async (id: string) => {
    setNotifications((prev) => prev.map((n) => (n.id === id ? { ...n, read: true } : n)));
    try {
      await apiRequest(`/notifications/${id}/read`, { method: 'PUT' });
    } catch {
      setNotifications((prev) => prev.map((n) => (n.id === id ? { ...n, read: false } : n)));
    }
  };

  const handleMarkAllRead = async () => {
    setMarkingAll(true);
    try {
      await apiRequest('/notifications/read-all', { method: 'PUT' });
      setNotifications((prev) => prev.map((n) => ({ ...n, read: true })));
    } catch {
      // silent
    } finally {
      setMarkingAll(false);
    }
  };

  const unreadCount = notifications.filter((n) => !n.read).length;

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">通知</h1>
          <p className="text-sm text-gray-500 mt-1">
            {unreadCount > 0 ? `${unreadCount} 条未读通知` : '所有通知已读'}
          </p>
        </div>
        <div className="flex items-center gap-3">
          {unreadCount > 0 && (
            <button
              onClick={handleMarkAllRead}
              disabled={markingAll}
              className="inline-flex items-center gap-2 px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-btn hover:bg-gray-50 disabled:opacity-50 transition-colors"
            >
              <CheckCheck size={16} />
              {markingAll ? '处理中...' : '全部标为已读'}
            </button>
          )}
          <button
            onClick={fetchNotifications}
            disabled={loading}
            className="inline-flex items-center gap-2 px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-btn hover:bg-gray-50 disabled:opacity-50 transition-colors"
          >
            <RefreshCw size={16} className={loading ? 'animate-spin' : ''} />
            刷新
          </button>
        </div>
      </div>

      {error && (
        <div className="mb-6 p-4 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
          {error}
        </div>
      )}

      {loading && notifications.length === 0 ? (
        <div className="space-y-3">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="bg-white rounded-card border border-gray-100 p-5 animate-pulse">
              <div className="flex items-start gap-4">
                <div className="w-10 h-10 rounded-full bg-gray-100 shrink-0" />
                <div className="flex-1 space-y-2">
                  <div className="h-4 w-40 bg-gray-100 rounded" />
                  <div className="h-3 w-full bg-gray-50 rounded" />
                  <div className="h-3 w-24 bg-gray-50 rounded" />
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : notifications.length === 0 ? (
        <div className="bg-white rounded-card border border-gray-100 p-5">
          <div className="text-center py-16">
            <Bell className="mx-auto h-12 w-12 text-gray-300 mb-4" />
            <p className="text-gray-500">暂无通知</p>
          </div>
        </div>
      ) : (
        <div className="space-y-2">
          {notifications.map((notif) => {
            const typeCfg = TYPE_CONFIG[notif.type] || TYPE_CONFIG.info;
            const TypeIcon = typeCfg.icon;
            return (
              <div
                key={notif.id}
                className={`bg-white rounded-card border transition-all ${
                  notif.read ? 'border-gray-100' : 'border-blue-100 shadow-sm'
                }`}
              >
                <div className="p-5">
                  <div className="flex items-start gap-4">
                    <div className={`p-2 rounded-lg shrink-0 ${typeCfg.bg}`}>
                      <TypeIcon size={18} className={typeCfg.color} />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <h3 className={`text-sm font-medium ${notif.read ? 'text-gray-700' : 'text-gray-900'}`}>
                          {notif.title}
                        </h3>
                        {!notif.read && (
                          <span className="w-2 h-2 rounded-full bg-blue-500 shrink-0" />
                        )}
                        <span className="text-xs text-gray-400 ml-auto shrink-0">
                          {formatDate(notif.created_at)}
                        </span>
                      </div>
                      <p className={`text-sm ${notif.read ? 'text-gray-500' : 'text-gray-600'}`}>
                        {notif.content}
                      </p>
                      <div className="flex items-center gap-3 mt-2">
                        {!notif.read && (
                          <button
                            onClick={() => handleMarkRead(notif.id)}
                            className="text-xs font-medium text-primary hover:text-blue-700 transition-colors"
                          >
                            标为已读
                          </button>
                        )}
                        {notif.link && (
                          <a
                            href={notif.link}
                            className="inline-flex items-center gap-1 text-xs font-medium text-gray-500 hover:text-gray-700 transition-colors"
                          >
                            查看详情
                            <ExternalLink size={12} />
                          </a>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
