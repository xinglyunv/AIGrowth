'use client';
import { useEffect, useState } from 'react';
import { apiRequest } from '@/lib/api';
import Link from 'next/link';
import { Card } from '@tremor/react';
import { FileText, ArrowRight, Calendar, BarChart3, MessageSquare, CheckCircle2, ChevronDown, ChevronUp, ExternalLink } from 'lucide-react';

interface Task {
  id: string;
  project_id: string;
  project_name: string;
  model: string;
  status: string;
  questions_count: number;
  completed_count: number;
  created_at: string;
}

interface Report {
  task: { model: string; status: string };
  project: { name: string; industry: string };
  visibility_score: number;
  brand_mentions: number;
  total_questions: number;
  answers: { question: string; sentiment: string; mentions_brand: boolean; rank_position: number }[];
  recommendations: string[];
}

interface TaskWithReport extends Task {
  report?: Report | null;
}

interface GroupedTasks {
  [projectName: string]: TaskWithReport[];
}

export default function ReportsPage() {
  const [tasks, setTasks] = useState<TaskWithReport[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [loadingReports, setLoadingReports] = useState<Set<string>>(new Set());

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await apiRequest<Task[] | { data: Task[] }>('/tasks?limit=100');
        const payload = res.data;
        const tasks = Array.isArray(payload) ? payload : payload?.data ?? [];
        const taskData: TaskWithReport[] = tasks
          .filter((t) => t.status === 'completed')
          .map((t) => ({ ...t }));
        setTasks(taskData);
      } catch (err: unknown) {
        const message = err instanceof Error ? err.message : '加载报告列表失败';
        setError(message);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const loadReport = async (taskId: string) => {
    if (loadingReports.has(taskId)) return;
    setLoadingReports((prev) => new Set(prev).add(taskId));
    try {
      const res = await apiRequest<Report>(`/tasks/${taskId}/report`);
      setTasks((prev) => prev.map((t) => (t.id === taskId ? { ...t, report: res.data } : t)));
    } catch {
      setTasks((prev) => prev.map((t) => (t.id === taskId ? { ...t, report: null } : t)));
    } finally {
      setLoadingReports((prev) => {
        const next = new Set(prev);
        next.delete(taskId);
        return next;
      });
    }
  };

  const toggleExpand = (taskId: string) => {
    if (expandedId === taskId) {
      setExpandedId(null);
    } else {
      setExpandedId(taskId);
      const task = tasks.find((t) => t.id === taskId);
      if (task && task.report === undefined) {
        loadReport(taskId);
      }
    }
  };

  const grouped = tasks.reduce<GroupedTasks>((acc, task) => {
    const name = task.project_name || '未分组';
    if (!acc[name]) acc[name] = [];
    acc[name].push(task);
    return acc;
  }, {});

  const scoreColor = (score: number) => {
    if (score >= 60) return 'bg-green-100 text-green-700';
    if (score >= 30) return 'bg-yellow-100 text-yellow-700';
    return 'bg-red-100 text-red-700';
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">报告中心</h1>
          <p className="text-sm text-gray-500 mt-1">按品牌项目分组的 AI 可见度检测报告</p>
        </div>
        <Link
          href="/tasks"
          className="flex items-center gap-2 px-4 py-2 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 transition-colors"
        >
          <BarChart3 size={16} />
          创建检测任务
        </Link>
      </div>

      {error && (
        <div className="mb-6 p-4 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
          {error}
        </div>
      )}

      {loading ? (
        <div className="space-y-4">
          {[1, 2].map((i) => (
            <div key={i} className="bg-white rounded-card border border-gray-100 p-5">
              <div className="h-6 w-40 bg-gray-100 rounded animate-pulse mb-4" />
              <div className="space-y-3">
                {[1, 2, 3].map((j) => (
                  <div key={j} className="h-24 bg-gray-50 rounded-lg animate-pulse" />
                ))}
              </div>
            </div>
          ))}
        </div>
      ) : tasks.length === 0 ? (
        <Card className="bg-white rounded-card border border-gray-100 p-5">
          <div className="text-center py-16">
            <FileText className="mx-auto h-12 w-12 text-gray-300 mb-4" />
            <p className="text-gray-500 mb-4">暂无已完成报告</p>
            <Link href="/tasks" className="text-sm font-medium text-primary hover:text-blue-700">
              前往创建检测任务
            </Link>
          </div>
        </Card>
      ) : (
        <div className="space-y-6">
          {Object.entries(grouped).map(([projectName, projectTasks]) => (
            <div key={projectName}>
              <div className="flex items-center gap-3 mb-3">
                <div className="w-1 h-6 bg-primary rounded-full" />
                <h2 className="text-base font-semibold text-gray-900">{projectName}</h2>
                <span className="text-xs text-gray-400 bg-gray-100 px-2 py-0.5 rounded-full">
                  {projectTasks.length} 份报告
                </span>
              </div>
              <div className="space-y-3">
                {projectTasks.map((task) => {
                  const report = task.report;
                   const score = report ? report.visibility_score : null;
                  const mentions = report ? report.brand_mentions : null;
                  const answered = task.completed_count ?? task.questions_count ?? 0;
                  const total = task.questions_count ?? 0;
                  const isExpanded = expandedId === task.id;
                  const isLoadingReport = loadingReports.has(task.id);

                  return (
                    <div key={task.id}>
                      <div
                        className="block bg-white rounded-card border border-gray-100 p-5 hover:border-primary/30 hover:shadow-sm transition-all cursor-pointer"
                        onClick={() => toggleExpand(task.id)}
                      >
                        <div className="flex items-start justify-between mb-3">
                          <div className="flex items-center gap-2">
                            <span className="text-sm font-medium text-gray-900">{task.model}</span>
                            {score !== null && (
                              <span className={`px-2 py-0.5 rounded text-xs font-medium ${scoreColor(score)}`}>
                                {score}分
                              </span>
                            )}
                          </div>
                          <div className="flex items-center gap-1 text-xs font-medium text-primary">
                            {isExpanded ? '收起详情' : '查看报告'}
                            {isExpanded ? <ChevronUp size={12} /> : <ChevronDown size={12} />}
                          </div>
                        </div>

                        <div className="flex flex-wrap items-center gap-x-5 gap-y-2 text-xs text-gray-500">
                          <span className="flex items-center gap-1">
                            <Calendar size={14} className="text-gray-400" />
                            {new Date(task.created_at).toLocaleDateString('zh-CN')}
                          </span>
                          <span className="flex items-center gap-1">
                            <CheckCircle2 size={14} className="text-gray-400" />
                            已回答 {answered}/{total}
                          </span>
                          {mentions !== null && (
                            <span className="flex items-center gap-1">
                              <MessageSquare size={14} className="text-gray-400" />
                              品牌提及 {mentions} 次
                            </span>
                          )}
                          {score !== null && (
                            <span className="flex items-center gap-1">
                              <BarChart3 size={14} className="text-gray-400" />
                              可见度评分 {score}
                            </span>
                          )}
                        </div>
                      </div>

                      {isExpanded && (
                        <div className="bg-gray-50 rounded-card border border-gray-100 border-t-0 p-5">
                          {isLoadingReport ? (
                            <div className="flex items-center justify-center py-8">
                              <div className="animate-spin w-5 h-5 border-2 border-primary border-t-transparent rounded-full" />
                              <span className="ml-2 text-sm text-gray-500">加载报告中...</span>
                            </div>
                          ) : report ? (
                            <div className="space-y-4">
                              <div className="flex gap-4">
                                <div className="flex-1 bg-white rounded-lg p-4 border border-gray-100">
                                  <p className="text-xs text-gray-500 mb-1">可见度分数</p>
                                  <p className={`text-xl font-bold ${score && score >= 60 ? 'text-green-600' : score && score >= 30 ? 'text-yellow-600' : 'text-red-600'}`}>
                                     {report.visibility_score}%
                                  </p>
                                </div>
                                <div className="flex-1 bg-white rounded-lg p-4 border border-gray-100">
                                  <p className="text-xs text-gray-500 mb-1">品牌提及</p>
                                  <p className="text-xl font-bold text-gray-900">{report.brand_mentions}</p>
                                </div>
                                <div className="flex-1 bg-white rounded-lg p-4 border border-gray-100">
                                  <p className="text-xs text-gray-500 mb-1">问题总数</p>
                                  <p className="text-xl font-bold text-gray-900">{report.total_questions}</p>
                                </div>
                              </div>

                              {report.answers && report.answers.length > 0 && (
                                <div>
                                  <h3 className="text-sm font-medium text-gray-900 mb-2">回答摘要</h3>
                                  <div className="space-y-2 max-h-48 overflow-y-auto">
                                    {report.answers.slice(0, 5).map((a, i) => (
                                      <div key={i} className="bg-white rounded-lg p-3 border border-gray-100 text-sm">
                                        <p className="text-gray-900 font-medium mb-1">{a.question}</p>
                                        <div className="flex items-center gap-3 text-xs text-gray-500">
                                          <span className={a.mentions_brand ? 'text-green-600' : 'text-gray-400'}>
                                            {a.mentions_brand ? '已提及' : '未提及'}
                                          </span>
                                          <span>{a.sentiment === 'positive' ? '正面' : a.sentiment === 'negative' ? '负面' : '中性'}</span>
                                          <span>排名 #{a.rank_position}</span>
                                        </div>
                                      </div>
                                    ))}
                                  </div>
                                </div>
                              )}

                              <Link
                                href={`/tasks/${task.id}`}
                                className="inline-flex items-center gap-1 text-sm font-medium text-primary hover:text-blue-700"
                              >
                                查看完整报告
                                <ExternalLink size={14} />
                              </Link>
                            </div>
                          ) : (
                            <div className="text-center py-8 text-sm text-gray-400">
                              报告数据加载失败
                            </div>
                          )}
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
