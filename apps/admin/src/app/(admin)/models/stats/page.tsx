'use client';

import { useEffect, useState } from 'react';
import {
  BarChart3, PhoneCall, Database, DollarSign, TrendingDown,
  RefreshCw, ChevronLeft, Info, Loader2
} from 'lucide-react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '@/lib/api';

interface ModelStat {
  model: string;
  calls: number;
  tokens: number;
  cost: number;
  success_rate: number;
  avg_latency: string;
}

interface ModelStatsData {
  total_calls: number;
  total_tokens: number;
  total_cost: number;
  success_rate: number;
  breakdown: ModelStat[];
}

const TIME_RANGES = ['today', '7d', '30d'] as const;
const timeLabels: Record<string, string> = { today: '今天', '7d': '近 7 天', '30d': '近 30 天' };

export default function ModelStatsPage() {
  const router = useRouter();
  const [timeRange, setTimeRange] = useState<string>('7d');
  const [summary, setSummary] = useState<{ total_calls: number; total_tokens: number; total_cost: number; failure_rate: number } | null>(null);
  const [breakdown, setBreakdown] = useState<ModelStat[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    setLoading(true);
    setError('');
    apiRequest<ModelStatsData>(`/models/stats?range=${timeRange}`)
      .then((r) => {
        if (r.data) {
          setSummary({ total_calls: r.data.total_calls, total_tokens: r.data.total_tokens, total_cost: r.data.total_cost, failure_rate: r.data.success_rate !== undefined ? 100 - r.data.success_rate : 0 });
          if (Array.isArray(r.data.breakdown)) setBreakdown(r.data.breakdown);
        }
      })
      .catch((e) => setError(e instanceof Error ? e.message : '加载失败'))
      .finally(() => setLoading(false));
  }, [timeRange]);

  const hasData = summary !== null;

  return (
    <div>
      <button onClick={() => router.push('/models')}
        className="flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 mb-4">
        <ChevronLeft className="h-4 w-4" />
        返回模型列表
      </button>

      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <BarChart3 className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">AI 调用统计</h2>
        </div>
        <div className="flex items-center gap-2">
          {TIME_RANGES.map((t) => (
            <button key={t} onClick={() => setTimeRange(t)}
              className={`px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
                timeRange === t ? 'bg-primary text-white' : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50'
              }`}>
              {timeLabels[t]}
            </button>
          ))}
        </div>
      </div>

      {loading && <div className="p-12 text-center text-gray-400"><Loader2 className="animate-spin inline h-5 w-5" /> 加载中...</div>}

      {error && <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-sm text-red-700">{error}</div>}

      {!loading && !error && !hasData && (
        <div className="p-12 text-center text-gray-400">
          <BarChart3 className="h-10 w-10 text-gray-300 mx-auto mb-3" />
          <p className="text-sm">暂无数据</p>
        </div>
      )}

      {!loading && !error && hasData && (
        <>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        {[
          { title: '总调用次数', value: summary!.total_calls.toLocaleString(), icon: PhoneCall, color: 'text-blue-600', bg: 'bg-blue-50' },
          { title: 'Token 消耗', value: (summary!.total_tokens / 10000).toFixed(0) + '万', icon: Database, color: 'text-purple-600', bg: 'bg-purple-50' },
          { title: '总费用', value: `¥${(summary!.total_cost ?? 0).toFixed(2)}`, icon: DollarSign, color: 'text-green-600', bg: 'bg-green-50' },
          { title: '失败率', value: `${(summary!.failure_rate ?? 0)}%`, icon: TrendingDown, color: (summary!.failure_rate ?? 0) > 3 ? 'text-red-600' : 'text-yellow-600', bg: (summary!.failure_rate ?? 0) > 3 ? 'bg-red-50' : 'bg-yellow-50' },
        ].map((c) => (
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

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <div className="px-5 py-4 border-b border-gray-200">
          <h3 className="text-base font-semibold text-gray-900">模型明细</h3>
        </div>
        <div className="overflow-x-auto">
          {breakdown.length === 0 ? (
            <div className="p-12 text-center text-gray-400 text-sm">暂无模型数据</div>
          ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 bg-gray-50 text-left">
                <th className="px-5 py-3 font-semibold text-gray-600">模型</th>
                <th className="px-5 py-3 font-semibold text-gray-600">调用次数</th>
                <th className="px-5 py-3 font-semibold text-gray-600">Token 消耗</th>
                <th className="px-5 py-3 font-semibold text-gray-600">费用</th>
                <th className="px-5 py-3 font-semibold text-gray-600">成功率</th>
                <th className="px-5 py-3 font-semibold text-gray-600">平均延迟</th>
              </tr>
            </thead>
            <tbody>
              {breakdown.map((s) => (
                <tr key={s.model} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                  <td className="px-5 py-3 font-medium text-gray-900">{s.model}</td>
                  <td className="px-5 py-3 text-gray-700">{s.calls.toLocaleString()}</td>
                  <td className="px-5 py-3 text-gray-700">{(s.tokens / 10000).toFixed(0)}万</td>
                  <td className="px-5 py-3 text-gray-700">{`¥${s.cost.toFixed(2)}`}</td>
                  <td className="px-5 py-3">
                    <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${s.success_rate >= 99 ? 'bg-green-50 text-green-700' : s.success_rate >= 97 ? 'bg-yellow-50 text-yellow-700' : 'bg-red-50 text-red-700'}`}>
                      {s.success_rate}%
                    </span>
                  </td>
                  <td className="px-5 py-3 text-gray-500">{s.avg_latency}</td>
                </tr>
              ))}
            </tbody>
          </table>
          )}
        </div>
      </div>
        </>
      )}
    </div>
  );
}
