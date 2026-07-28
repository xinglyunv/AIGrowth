'use client';
import { useState, useEffect, useRef } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { apiRequest } from '@/lib/api';
import { ArrowLeft, FlaskConical, Loader2 } from 'lucide-react';

interface Project {
  id: string;
  name: string;
}

interface ModelOption {
  id: string;
  name: string;
  model: string;
  provider: string;
  enabled: boolean;
}

export default function NewTaskPage() {
  const params = useParams();
  const router = useRouter();
  const projectId = params.id as string;

  const [project, setProject] = useState<Project | null>(null);
  const [models, setModels] = useState<ModelOption[]>([]);
  const [selectedModelIds, setSelectedModelIds] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [modelsLoading, setModelsLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [progress, setProgress] = useState(0);
  const [countdown, setCountdown] = useState(0);
  const pollingRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    apiRequest<Project>(`/projects/${projectId}`)
      .then((res) => {
        const { name } = res.data as unknown as Project;
        setProject({ id: projectId, name });
      })
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : '加载项目信息失败';
        setError(message);
      })
      .finally(() => setLoading(false));
    apiRequest<ModelOption[]>('/models')
      .then((res) => {
        if (Array.isArray(res.data)) {
          const enabled = res.data.filter((m) => m.enabled !== false);
          setModels(enabled);
          if (enabled.length > 0) setSelectedModelIds([enabled[0].model]);
        }
      })
      .catch(() => {})
      .finally(() => setModelsLoading(false));
  }, [projectId]);

  const handleSubmit = async () => {
    setSubmitting(true);
    setError('');
    setProgress(0);

    try {
      const modelParam = selectedModelIds.join(',');
      const createRes = await apiRequest<{ id: string }>('/tasks', {
        method: 'POST',
        body: JSON.stringify({ project_id: projectId, model: modelParam }),
      });
      const taskId = createRes.data.id;

      // Start execution (returns immediately as async)
      await apiRequest(`/tasks/${taskId}/execute`, { method: 'POST' });

      // Estimate: ~20s per model × 5 questions
      const estimatedSeconds = selectedModelIds.length * 20;
      setCountdown(estimatedSeconds);

      // Poll for task status every 2 seconds
      pollingRef.current = setInterval(async () => {
        try {
          const res = await apiRequest<{ status: string; completed_count: number; questions_count: number; error_message?: string }>(`/tasks/${taskId}`);
          if (res.data.status === 'completed') {
            if (pollingRef.current) clearInterval(pollingRef.current);
            router.push(`/tasks/${taskId}`);
          } else if (res.data.status === 'failed') {
            if (pollingRef.current) clearInterval(pollingRef.current);
            setError(res.data.error_message || '任务执行失败');
            setSubmitting(false);
          } else {
            const pct = res.data.questions_count > 0
              ? Math.round((res.data.completed_count / res.data.questions_count) * 100)
              : 0;
            setProgress(pct);
          }
        } catch {
          // Ignore polling errors, keep trying
        }
      }, 2000);

      // Countdown timer
      const timer = setInterval(() => {
        setCountdown((prev) => {
          if (prev <= 0) return 0;
          return prev - 1;
        });
      }, 1000);

      return () => {
        clearInterval(timer);
        if (pollingRef.current) clearInterval(pollingRef.current);
      };
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : '创建任务失败';
      setError(message);
      setSubmitting(false);
    }
  };

  // Cleanup polling on unmount
  useEffect(() => {
    return () => {
      if (pollingRef.current) clearInterval(pollingRef.current);
    };
  }, []);

  if (loading) {
    return (
      <div className="animate-pulse space-y-6">
        <div className="h-5 w-24 bg-gray-200 rounded" />
        <div className="h-8 w-48 bg-gray-200 rounded" />
        <div className="h-48 bg-gray-100 rounded-card" />
      </div>
    );
  }

  if (error && !project) {
    return (
      <div>
        <Link href={`/projects/${projectId}`} className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 transition-colors mb-6">
          <ArrowLeft size={16} />
          返回项目详情
        </Link>
        <div className="p-4 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">{error}</div>
      </div>
    );
  }

  return (
    <div>
      <Link href={`/projects/${projectId}`} className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 transition-colors mb-6">
        <ArrowLeft size={16} />
        返回项目详情
      </Link>

      <h1 className="text-2xl font-bold text-gray-900 mb-6">创建 AI 检测任务</h1>

      <div className="bg-white rounded-card shadow-sm border border-gray-100 p-6 max-w-xl">
        {/* 项目名称 */}
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-1.5">项目</label>
          <p className="text-sm text-gray-900 bg-gray-50 rounded-btn px-3 py-2 border border-gray-200">
            {project?.name || '-'}
          </p>
        </div>

        {/* 模型选择 */}
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-1.5">选择模型（支持多选进行对比分析）</label>
          {modelsLoading ? (
            <p className="text-sm text-gray-400 flex items-center gap-2"><Loader2 className="h-4 w-4 animate-spin" />加载模型列表...</p>
          ) : models.length === 0 ? (
            <p className="text-sm text-red-500">暂无可用的 AI 模型，请联系管理员配置</p>
          ) : (
            <div className="space-y-2">
              {models.map((m) => (
                <label
                  key={m.id}
                  className="flex items-center gap-3 px-3 py-2.5 rounded-btn border border-gray-200 hover:border-primary/30 hover:bg-primary/5 transition-colors cursor-pointer"
                >
                  <input
                    type="checkbox"
                    checked={selectedModelIds.includes(m.model)}
                    onChange={(e) => {
                      if (e.target.checked) {
                        setSelectedModelIds([...selectedModelIds, m.model]);
                      } else {
                        setSelectedModelIds(selectedModelIds.filter(id => id !== m.model));
                      }
                    }}
                    className="rounded border-gray-300 text-primary focus:ring-primary/30"
                  />
                  <div className="flex-1 min-w-0">
                    <span className="text-sm font-medium text-gray-900">{m.name}</span>
                    <span className="text-xs text-gray-500 ml-2">({m.provider})</span>
                  </div>
                </label>
              ))}
               {selectedModelIds.length > 0 && (
                 <p className="text-xs text-gray-500 mt-1">
                   将使用 {selectedModelIds.length} 个模型进行对比分析，预计消耗 {selectedModelIds.length} 积分
                 </p>
               )}
            </div>
          )}
        </div>

        {error && !submitting && (
          <div className="mb-6 p-3 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
            {error}
          </div>
        )}

        <button
          onClick={handleSubmit}
          disabled={submitting || modelsLoading || selectedModelIds.length === 0}
          className="inline-flex items-center gap-2 px-6 py-2.5 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <FlaskConical size={16} />
          {submitting ? '正在创建并执行...' : '开始检测'}
        </button>
      </div>

      {/* Full-screen loading overlay */}
      {submitting && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
          <div className="bg-white rounded-xl shadow-2xl p-10 max-w-sm w-full mx-4 text-center">
            <Loader2 className="h-10 w-10 animate-spin text-primary mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-gray-900 mb-2">AI 检测任务执行中</h3>
            <p className="text-sm text-gray-500 mb-4">
              正在使用 {selectedModelIds.length} 个模型进行分析，预计剩余 {countdown} 秒...
            </p>
            <div className="w-full bg-gray-100 rounded-full h-2 overflow-hidden">
              <div
                className="bg-primary h-full rounded-full transition-all duration-500 ease-out"
                style={{ width: `${Math.max(progress, 5)}%` }}
              />
            </div>
            <p className="text-xs text-gray-400 mt-2">{progress}% 完成</p>
          </div>
        </div>
      )}
    </div>
  );
}
