'use client';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { apiRequest } from '@/lib/api';
import { ChevronRight, Plus, Search } from 'lucide-react';

interface Task {
  id: string;
  project_id: string;
  project_name: string;
  model: string;
  status: string;
  total_questions: number;
  questions_count: number;
  created_at: string;
}

const STATUS_CONFIG: Record<string, { label: string; color: string }> = {
  pending: { label: '待处理', color: 'bg-gray-100 text-gray-600 border-gray-200' },
  running: { label: '运行中', color: 'bg-blue-100 text-blue-700 border-blue-200' },
  completed: { label: '已完成', color: 'bg-green-100 text-green-700 border-green-200' },
  failed: { label: '失败', color: 'bg-red-100 text-red-700 border-red-200' },
};

function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const d = new Date(dateStr);
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' });
}

function SkeletonRow() {
  return (
    <tr className="animate-pulse">
      <td className="px-6 py-4"><div className="h-4 w-32 bg-gray-200 rounded" /></td>
      <td className="px-6 py-4"><div className="h-4 w-16 bg-gray-200 rounded" /></td>
      <td className="px-6 py-4"><div className="h-5 w-12 bg-gray-200 rounded-full" /></td>
      <td className="px-6 py-4"><div className="h-4 w-8 bg-gray-200 rounded" /></td>
      <td className="px-6 py-4"><div className="h-4 w-24 bg-gray-200 rounded" /></td>
    </tr>
  );
}

export default function TasksPage() {
  const router = useRouter();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    apiRequest<Task[]>('/tasks')
      .then((res) => setTasks(res.data))
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : '加载任务列表失败';
        setError(message);
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">AI 检测任务</h1>
        <Link
          href="/tasks/new"
          className="inline-flex items-center gap-2 px-4 py-2 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 transition-colors"
        >
          <Plus size={16} />
          创建任务
        </Link>
      </div>

      {error && (
        <div className="mb-6 p-4 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
          {error}
        </div>
      )}

      <div className="bg-white rounded-card shadow-sm border border-gray-100 overflow-x-auto">
        <table className="w-full whitespace-nowrap">
          <thead>
            <tr className="border-b border-gray-100 bg-gray-50/50">
              <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">项目名称</th>
              <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">模型</th>
              <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">状态</th>
              <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">问题数</th>
              <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">创建时间</th>
              <th className="w-10 px-6 py-3" />
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-50">
            {loading ? (
              <>
                <SkeletonRow />
                <SkeletonRow />
                <SkeletonRow />
              </>
            ) : tasks.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-6 py-16 text-center">
                  <div className="flex flex-col items-center gap-3">
                    <div className="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center">
                      <Plus size={24} className="text-gray-400" />
                    </div>
                    <p className="text-sm text-gray-500">暂无检测任务，为品牌项目创建首个检测</p>
                    <Link
                      href="/tasks/new"
                      className="mt-1 inline-flex items-center gap-2 px-4 py-2 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 transition-colors"
                    >
                      <Plus size={16} />
                      创建任务
                    </Link>
                  </div>
                </td>
              </tr>
            ) : (
              tasks.map((task) => {
                const statusCfg = STATUS_CONFIG[task.status] || { label: task.status, color: 'bg-gray-100 text-gray-600 border-gray-200' };
                return (
                  <tr
                    key={task.id}
                    className="hover:bg-gray-50 transition-colors cursor-pointer"
                    onClick={() => router.push(`/tasks/${task.id}`)}
                  >
                    <td className="px-6 py-4">
                      <span className="text-sm font-medium text-gray-900">{task.project_name}</span>
                    </td>
                    <td className="px-6 py-4">
                      <span className="text-sm text-gray-600">{task.model}</span>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-medium border ${statusCfg.color}`}>
                        {statusCfg.label}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <span className="text-sm text-gray-600">{task.questions_count ?? '-'}</span>
                    </td>
                    <td className="px-6 py-4">
                      <span className="text-sm text-gray-500">{formatDate(task.created_at)}</span>
                    </td>
                    <td className="px-6 py-4">
                      <ChevronRight size={16} className="text-gray-300" />
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}