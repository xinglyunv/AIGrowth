'use client';
import { useEffect, useState } from 'react';
import { apiRequest } from '@/lib/api';
import { Card, Grid, BarChart, LineChart, DonutChart } from '@tremor/react';
import { BarChart3, TrendingUp, Target, ThumbsUp, Calendar } from 'lucide-react';

interface Stats {
  total_projects: number;
  total_tasks: number;
  completed_tasks: number;
  avg_visibility_score: number;
}

interface TaskItem {
  id: string;
  project_id: string;
  project_name: string;
  model: string;
  status: string;
  questions_count: number;
  completed_count: number;
  created_at: string;
}

interface ReportData {
  ai_visibility_score: number;
  brand_mentions: number;
  total_questions: number;
  answers: { question: string; sentiment: string; mentions_brand: boolean }[];
}

function formatShortDate(dateStr: string): string {
  const d = new Date(dateStr);
  return `${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
}

function getDateKey(dateStr: string): string {
  const d = new Date(dateStr);
  return `${d.getFullYear()}-${d.getMonth()}-${d.getDate()}`;
}

export default function AnalyticsPage() {
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [timeRange, setTimeRange] = useState('7d');

  const [visibilityData, setVisibilityData] = useState<{ date: string; score: number }[]>([]);
  const [taskTrendData, setTaskTrendData] = useState<{ date: string; completed: number; total: number }[]>([]);
  const [questionsData, setQuestionsData] = useState<{ question: string; mentions: number }[]>([]);
  const [sentimentData, setSentimentData] = useState<{ name: string; value: number }[]>([]);

  useEffect(() => {
    Promise.all([
      apiRequest<Stats>('/dashboard/stats'),
      apiRequest<TaskItem[]>('/tasks?limit=50'),
    ])
      .then(([statsRes, tasksRes]) => {
        setStats(statsRes.data);
         const taskPayload = tasksRes.data as unknown;
         const tasks = Array.isArray(taskPayload)
           ? taskPayload
           : (taskPayload as { data?: TaskItem[] } | null)?.data ?? [];

        const completedTasks = tasks.filter((t) => t.status === 'completed');
        const recentCompleted = completedTasks.slice(0, 10);

        const dateMap: Record<string, { completed: number; total: number; scores: number[] }> = {};
        tasks.forEach((t) => {
          const key = getDateKey(t.created_at);
          if (!dateMap[key]) dateMap[key] = { completed: 0, total: 0, scores: [] };
          dateMap[key].total += 1;
          if (t.status === 'completed') dateMap[key].completed += 1;
        });

        const sortedDates = Object.keys(dateMap).sort().slice(-7);
        const trendData = sortedDates.map((key) => ({
          date: formatShortDate(key),
          completed: dateMap[key].completed,
          total: dateMap[key].total,
        }));
        setTaskTrendData(trendData);

        const scoreDateMap: Record<string, number[]> = {};
        const questionMentionMap: Record<string, number> = {};
        let totalPositive = 0;
        let totalNeutral = 0;
        let totalNegative = 0;

        const reportPromises = recentCompleted.map((t) =>
          apiRequest<ReportData>(`/tasks/${t.id}/report`).catch(() => null)
        );

         Promise.all(reportPromises).then((reports) => {
           reports.forEach((r, index) => {
            if (!r || !r.data) return;
            const report = r.data;
             const scoreKey = getDateKey(recentCompleted[index]?.created_at || '');
            if (!scoreDateMap[scoreKey]) scoreDateMap[scoreKey] = [];
            scoreDateMap[scoreKey].push(report.ai_visibility_score);

            (report.answers || []).forEach((a) => {
              if (a.mentions_brand) {
                questionMentionMap[a.question] = (questionMentionMap[a.question] || 0) + 1;
              }
              if (a.sentiment === 'positive') totalPositive++;
              else if (a.sentiment === 'negative') totalNegative++;
              else totalNeutral++;
            });
          });

          const visibilitySortedKeys = Object.keys(scoreDateMap).sort().slice(-7);
          const visData = visibilitySortedKeys.map((key) => {
            const scores = scoreDateMap[key];
            const avg = scores.length > 0 ? Math.round(scores.reduce((a, b) => a + b, 0) / scores.length) : 0;
            return { date: formatShortDate(key), score: avg };
          });
          if (visData.length === 0) {
            visData.push({ date: formatShortDate(new Date().toISOString()), score: statsRes.data.avg_visibility_score });
          }
          setVisibilityData(visData);

          const sortedQuestions = Object.entries(questionMentionMap)
            .sort((a, b) => b[1] - a[1])
            .slice(0, 6)
            .map(([question, mentions]) => ({ question, mentions }));
          setQuestionsData(sortedQuestions.length > 0 ? sortedQuestions : []);

          const total = totalPositive + totalNeutral + totalNegative;
          if (total > 0) {
            setSentimentData([
              { name: '正面', value: Math.round((totalPositive / total) * 100) },
              { name: '中性', value: Math.round((totalNeutral / total) * 100) },
              { name: '负面', value: Math.round((totalNegative / total) * 100) },
            ]);
          } else {
            setSentimentData([]);
          }
        });
      })
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : '加载数据失败';
        setError(message);
      })
      .finally(() => setLoading(false));
  }, []);

  const scoreColor = stats && stats.avg_visibility_score >= 60
    ? 'text-green-600'
    : stats && stats.avg_visibility_score >= 30
    ? 'text-yellow-600'
    : 'text-red-600';

  const completionRate = stats && stats.total_tasks > 0
    ? Math.round((stats.completed_tasks / stats.total_tasks) * 100)
    : 0;

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-2xl font-bold text-gray-900">数据分析</h1>
        <div className="flex items-center gap-2">
          <Calendar size={16} className="text-gray-400" />
          <select
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value)}
            className="text-sm border border-gray-200 rounded-btn px-3 py-1.5 bg-white text-gray-700 focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
          >
            <option value="7d">近 7 天</option>
            <option value="30d">近 30 天</option>
            <option value="90d">近 90 天</option>
          </select>
        </div>
      </div>

      {error && (
        <div className="mb-6 p-4 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
          {error}
        </div>
      )}

      <Grid numItemsSm={2} numItemsLg={4} className="gap-4 mb-8">
        <Card className="bg-white rounded-card border border-gray-100 p-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-gray-500">品牌项目数</span>
            <div className="p-2.5 rounded-lg bg-blue-50">
              <BarChart3 size={20} className="text-blue-600" />
            </div>
          </div>
          <div className="text-3xl font-bold text-gray-900">
            {loading ? '-' : stats?.total_projects ?? 0}
          </div>
          <p className="text-xs text-gray-400 mt-1">已创建的品牌分析项目</p>
        </Card>
        <Card className="bg-white rounded-card border border-gray-100 p-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-gray-500">检测任务总数</span>
            <div className="p-2.5 rounded-lg bg-purple-50">
              <Target size={20} className="text-purple-600" />
            </div>
          </div>
          <div className="text-3xl font-bold text-gray-900">
            {loading ? '-' : stats?.total_tasks ?? 0}
          </div>
          <p className="text-xs text-gray-400 mt-1">已执行的 AI 可见度检测</p>
        </Card>
        <Card className="bg-white rounded-card border border-gray-100 p-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-gray-500">已完成任务</span>
            <div className="p-2.5 rounded-lg bg-green-50">
              <TrendingUp size={20} className="text-green-600" />
            </div>
          </div>
          <div className="text-3xl font-bold text-gray-900">
            {loading ? '-' : stats?.completed_tasks ?? 0}
          </div>
          <p className="text-xs text-gray-400 mt-1">{completionRate}% 完成率</p>
        </Card>
        <Card className="bg-white rounded-card border border-gray-100 p-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-gray-500">平均可见度</span>
            <div className="p-2.5 rounded-lg bg-orange-50">
              <ThumbsUp size={20} className="text-orange-600" />
            </div>
          </div>
          <div className={`text-3xl font-bold ${scoreColor}`}>
            {loading ? '-' : stats ? `${stats.avg_visibility_score}%` : '0%'}
          </div>
          <p className="text-xs text-gray-400 mt-1">品牌在各 AI 模型中的可见度评分</p>
        </Card>
      </Grid>

      <Grid numItemsSm={1} numItemsLg={2} className="gap-6 mb-8">
        <Card className="bg-white rounded-card border border-gray-100 p-5">
          <h2 className="text-base font-semibold text-gray-900 mb-4">品牌可见度趋势</h2>
          {visibilityData.length > 0 ? (
            <BarChart
              data={visibilityData}
              index="date"
              categories={['score']}
              colors={['blue']}
              valueFormatter={(v) => `${v}%`}
              className="h-72"
              showAnimation
            />
          ) : (
            <div className="h-72 flex items-center justify-center text-sm text-gray-400">暂无可见度数据</div>
          )}
        </Card>
        <Card className="bg-white rounded-card border border-gray-100 p-5">
          <h2 className="text-base font-semibold text-gray-900 mb-4">检测任务完成趋势</h2>
          {taskTrendData.length > 0 ? (
            <LineChart
              data={taskTrendData}
              index="date"
              categories={['completed', 'total']}
              colors={['green', 'gray']}
              valueFormatter={(v) => `${v}`}
              className="h-72"
              showAnimation
            />
          ) : (
            <div className="h-72 flex items-center justify-center text-sm text-gray-400">暂无任务数据</div>
          )}
        </Card>
      </Grid>

      <Grid numItemsSm={1} numItemsLg={2} className="gap-6">
        <Card className="bg-white rounded-card border border-gray-100 p-5">
          <h2 className="text-base font-semibold text-gray-900 mb-4">常见问题分析</h2>
          {questionsData.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-100">
                    <th className="text-left py-2 px-1 font-medium text-gray-500">问题主题</th>
                    <th className="text-right py-2 px-1 font-medium text-gray-500 w-24">品牌提及数</th>
                  </tr>
                </thead>
                <tbody>
                  {questionsData.map((q) => (
                    <tr key={q.question} className="border-b border-gray-50 last:border-0">
                      <td className="py-2.5 px-1 text-gray-900">{q.question}</td>
                      <td className="py-2.5 px-1 text-right">
                        <span className="inline-flex items-center justify-center min-w-[2rem] px-2 py-0.5 rounded-full text-xs font-medium bg-blue-50 text-blue-700">
                          {q.mentions}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="h-60 flex items-center justify-center text-sm text-gray-400">暂无问题数据</div>
          )}
        </Card>
        <Card className="bg-white rounded-card border border-gray-100 p-5">
          <h2 className="text-base font-semibold text-gray-900 mb-4">品牌提及情感分布</h2>
          {sentimentData.length > 0 ? (
            <>
              <div className="flex items-center justify-center h-72">
                <DonutChart
                  data={sentimentData}
                  category="value"
                  index="name"
                  colors={['green', 'gray', 'red']}
                  valueFormatter={(v) => `${v}%`}
                  className="h-60 w-full"
                  showAnimation
                />
              </div>
              <div className="flex items-center justify-center gap-6 mt-2">
                {sentimentData.map((s) => (
                  <div key={s.name} className="flex items-center gap-1.5 text-sm">
                    <span
                      className={`w-2.5 h-2.5 rounded-full ${
                        s.name === '正面' ? 'bg-green-500' : s.name === '中性' ? 'bg-gray-400' : 'bg-red-500'
                      }`}
                    />
                    <span className="text-gray-600">{s.name}</span>
                    <span className="font-medium text-gray-900">{s.value}%</span>
                  </div>
                ))}
              </div>
            </>
          ) : (
            <div className="h-72 flex items-center justify-center text-sm text-gray-400">暂无情感数据</div>
          )}
        </Card>
      </Grid>
    </div>
  );
}
