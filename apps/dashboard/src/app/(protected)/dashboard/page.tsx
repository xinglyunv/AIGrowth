'use client';
import { useEffect, useState } from 'react';
import { useAuth } from '@/lib/auth';
import { apiRequest } from '@/lib/api';
import Link from 'next/link';
import {
  FolderOpen, Search, Zap, TrendingUp, Clock,
  ArrowRight, Plus, Loader2, AlertCircle,
  Activity, ChevronRight,
} from 'lucide-react';

interface DashboardStats {
  projects?: number;
  tasks?: number;
  completed?: number;
  visibilityScore?: number;
  total_projects?: number;
  total_tasks?: number;
  completed_tasks?: number;
  avg_visibility_score?: number;
}

interface TaskItem {
  id: string;
  project_id: string;
  project_name: string;
  model: string;
  status: string;
  created_at: string;
}

const STATUS_LABEL: Record<string, string> = {
  pending: '待处理',
  running: '运行中',
  completed: '已完成',
  failed: '失败',
};

const STATUS_BADGE: Record<string, string> = {
  pending: 'text-gray-600 bg-gray-100',
  running: 'text-blue-700 bg-blue-100',
  completed: 'text-green-700 bg-green-100',
  failed: 'text-red-700 bg-red-100',
};

export default function DashboardPage() {
  const { user } = useAuth();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [credits, setCredits] = useState<number | null>(null);
  const [recentTasks, setRecentTasks] = useState<TaskItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    Promise.all([
      apiRequest<DashboardStats>('/dashboard/stats'),
      apiRequest<{ credits: number }>('/billing/credits').catch(() => null),
      apiRequest<{ data: TaskItem[]; total: number }>('/tasks?limit=5'),
    ])
      .then(([statsRes, creditsRes, tasksRes]) => {
        setStats(statsRes.data);
        if (creditsRes) setCredits(creditsRes.data.credits);
        setRecentTasks(tasksRes.data.data || []);
      })
      .catch((err) => {
        setError(err instanceof Error ? err.message : '加载数据失败');
      })
      .finally(() => setLoading(false));
  }, []);

  const todayTaskCount = recentTasks.filter((t) => {
    const d = new Date(t.created_at);
    const now = new Date();
    return d.toDateString() === now.toDateString();
  }).length;

  const scoreValue = stats?.visibilityScore ?? stats?.avg_visibility_score;
  const scoreColor = scoreValue !== undefined
    ? scoreValue >= 60 ? 'text-green-600'
      : scoreValue >= 30 ? 'text-amber-600' : 'text-red-600'
    : 'text-gray-400';

  const statCards = [
    {
      title: '品牌项目数',       value: stats?.projects ?? stats?.total_projects ?? 0,
      icon: FolderOpen, color: 'text-blue-600', bg: 'bg-blue-50',
      link: '/projects', linkText: '查看项目',
    },
    {
      title: 'AI 检测次数',       value: stats?.tasks ?? stats?.total_tasks ?? 0,
      icon: Search, color: 'text-purple-600', bg: 'bg-purple-50',
      link: '/tasks', linkText: '查看任务',
    },
    {
      title: '剩余额度', value: credits,
      icon: Zap, color: 'text-amber-600', bg: 'bg-amber-50',
      link: '/billing', linkText: '充值',
    },
    {
      title: '可见度分数', value: scoreValue, suffix: '',
      icon: TrendingUp, color: scoreColor, bg: 'bg-green-50',
      link: '/tasks', linkText: '查看报告',
    },
    {
      title: '今日检测', value: todayTaskCount,
      icon: Clock, color: 'text-rose-600', bg: 'bg-rose-50',
      link: '/tasks', linkText: '查看详情',
    },
  ];

  const quickActions = [
    {
      title: '新建项目', desc: '添加新品牌开始 AI 可见度分析',
      icon: Plus, bg: 'bg-blue-50', color: 'text-blue-600', href: '/projects/new',
    },
    {
      title: '新建检测', desc: '对已有项目执行 AI 可见度检测',
      icon: Activity, bg: 'bg-purple-50', color: 'text-purple-600', href: '/projects',
    },
    {
      title: '查看账单', desc: '查看剩余额度和充值记录',
      icon: Zap, bg: 'bg-amber-50', color: 'text-amber-600', href: '/billing',
    },
  ];

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="flex items-center gap-3 text-gray-400">
          <Loader2 className="animate-spin" size={24} />
          <span className="text-sm">加载仪表盘数据...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-4">
        <div className="p-3 rounded-full bg-red-50">
          <AlertCircle size={24} className="text-red-500" />
        </div>
        <p className="text-sm text-gray-600">{error}</p>
        <button
          onClick={() => window.location.reload()}
          className="px-4 py-2 text-sm font-medium text-white bg-primary rounded-lg hover:bg-blue-700 transition-colors"
        >
          重新加载
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">
            欢迎回来，{user?.username}
          </h1>
          <p className="text-sm text-gray-500 mt-1">您的品牌 AI 可见度分析概览</p>
        </div>
        <Link
          href="/projects/new"
          className="inline-flex items-center gap-2 px-4 py-2 bg-primary text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors shadow-sm"
        >
          <Plus size={16} />
          新建项目
        </Link>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-4">
        {statCards.map((card) => {
          const Icon = card.icon;
          const displayVal = card.value !== undefined && card.value !== null ? card.value : '-';
          return (
            <div
              key={card.title}
              className="bg-white rounded-xl border border-gray-100 shadow-sm p-5 hover:shadow-md transition-shadow"
            >
              <div className="flex items-center justify-between mb-3">
                <span className="text-sm font-medium text-gray-500">{card.title}</span>
                <div className={`p-2 rounded-lg ${card.bg}`}>
                  <Icon size={18} className={card.color} />
                </div>
              </div>
              <div className={`text-2xl font-bold ${card.color} mb-3`}>
                {displayVal}
                {card.suffix}
              </div>
              <Link
                href={card.link}
                className="inline-flex items-center gap-1 text-xs font-medium text-gray-400 hover:text-primary transition-colors"
              >
                {card.linkText}
                <ArrowRight size={12} />
              </Link>
            </div>
          );
        })}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 bg-white rounded-xl border border-gray-100 shadow-sm p-5">
          <h2 className="text-base font-semibold text-gray-900 mb-4">最近活动</h2>
          {recentTasks.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-100">
                    <th className="text-left py-3 px-2 text-gray-500 font-medium">项目</th>
                    <th className="text-left py-3 px-2 text-gray-500 font-medium">模型</th>
                    <th className="text-left py-3 px-2 text-gray-500 font-medium">状态</th>
                    <th className="text-right py-3 px-2 text-gray-500 font-medium">日期</th>
                  </tr>
                </thead>
                <tbody>
                  {recentTasks.map((task) => (
                    <tr key={task.id} className="border-b border-gray-50 hover:bg-gray-50 transition-colors">
                      <td className="py-3 px-2">
                        <Link href={`/tasks/${task.id}`} className="font-medium text-gray-900 hover:text-primary">
                          {task.project_name}
                        </Link>
                      </td>
                      <td className="py-3 px-2 text-gray-500">
                        {task.model}
                      </td>
                      <td className="py-3 px-2">
                        <span className={`inline-flex px-2 py-0.5 rounded text-xs font-medium ${STATUS_BADGE[task.status] || 'text-gray-600 bg-gray-100'}`}>
                          {STATUS_LABEL[task.status] || task.status}
                        </span>
                      </td>
                      <td className="py-3 px-2 text-right text-gray-500">
                        {new Date(task.created_at).toLocaleDateString('zh-CN')}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="text-center py-12">
              <div className="p-3 rounded-full bg-gray-50 inline-flex mb-3">
                <Clock size={24} className="text-gray-300" />
              </div>
              <p className="text-sm text-gray-400 mb-3">暂无活动记录</p>
              <Link
                href="/projects"
                className="text-sm font-medium text-primary hover:text-blue-700"
              >
                创建项目开始 AI 检测
              </Link>
            </div>
          )}
        </div>

        <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-5">
          <h2 className="text-base font-semibold text-gray-900 mb-4">快速操作</h2>
          <div className="space-y-3">
            {quickActions.map((action) => {
              const Icon = action.icon;
              return (
                <Link
                  key={action.title}
                  href={action.href}
                  className="flex items-center gap-3 p-3 rounded-lg border border-gray-100 hover:border-primary/30 hover:bg-blue-50/50 transition-all"
                >
                  <div className={`p-2 rounded-lg ${action.bg}`}>
                    <Icon size={16} className={action.color} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900">{action.title}</p>
                    <p className="text-xs text-gray-500 truncate">{action.desc}</p>
                  </div>
                  <ChevronRight size={14} className="text-gray-300 flex-shrink-0" />
                </Link>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
