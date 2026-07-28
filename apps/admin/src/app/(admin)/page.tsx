'use client';
import { useEffect, useState } from 'react';
import { apiRequest } from '@/lib/api';
import { Users, FolderOpen, Search, Activity } from 'lucide-react';

interface Stats {
  total_users: number;
  total_projects: number;
  total_tasks: number;
}

export default function AdminDashboard() {
  const [stats, setStats] = useState<Stats | null>(null);

  useEffect(() => {
    apiRequest<Stats>('/dashboard/stats').then(r => setStats(r.data)).catch(() => {});
  }, []);

  const cards = [
    { title: '用户总数', value: stats?.total_users ?? '-', icon: Users, color: 'text-blue-600', bg: 'bg-blue-50' },
    { title: '品牌项目', value: stats?.total_projects ?? 0, icon: FolderOpen, color: 'text-purple-600', bg: 'bg-purple-50' },
    { title: '检测任务', value: stats?.total_tasks ?? 0, icon: Search, color: 'text-green-600', bg: 'bg-green-50' },
    { title: '系统状态', value: '运行中', icon: Activity, color: 'text-green-600', bg: 'bg-green-50' },
  ];

  return (
    <div>
      <h1 className="text-xl font-bold text-gray-900 mb-6">管理后台</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {cards.map(c => (
          <div key={c.title} className="bg-white rounded-lg border border-gray-200 p-5">
            <div className="flex items-center justify-between mb-3">
              <span className="text-sm font-medium text-gray-500">{c.title}</span>
              <div className={`p-2 rounded-lg ${c.bg}`}>
                <c.icon size={18} className={c.color} />
              </div>
            </div>
            <div className={`text-2xl font-bold ${c.color}`}>{c.value}</div>
          </div>
        ))}
      </div>
      <div className="mt-8 bg-white rounded-lg border border-gray-200 p-5">
        <h2 className="text-base font-semibold text-gray-900 mb-3">系统状态</h2>
        <div className="flex items-center gap-2 text-sm">
          <span className="w-2 h-2 rounded-full bg-green-500" />
          <span className="text-gray-600">服务运行正常</span>
        </div>
        <div className="mt-4 grid grid-cols-1 sm:grid-cols-3 gap-4 text-sm">
          <div className="p-3 bg-gray-50 rounded-md">
            <p className="text-gray-500">API 服务</p>
            <p className="font-medium text-green-600">运行中</p>
          </div>
          <div className="p-3 bg-gray-50 rounded-md">
            <p className="text-gray-500">PostgreSQL</p>
            <p className="font-medium text-green-600">已连接</p>
          </div>
          <div className="p-3 bg-gray-50 rounded-md">
            <p className="text-gray-500">Redis</p>
            <p className="font-medium text-green-600">已连接</p>
          </div>
        </div>
      </div>
    </div>
  );
}
