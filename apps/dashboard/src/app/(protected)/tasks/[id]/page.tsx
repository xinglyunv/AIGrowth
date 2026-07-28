'use client';
import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { apiRequest } from '@/lib/api';
import { ArrowLeft, FlaskConical, Brain, Building2, ThumbsUp, ThumbsDown, Minus, Trophy, TrendingUp, ChevronDown, ChevronUp, BarChart3 } from 'lucide-react';

interface Answer {
  id: string;
  question: string;
  answer: string;
  answer_summary: string;
  model: string;
  mentions_brand: boolean;
  sentiment: string;
  rank_position: number;
  analysis?: {
    entities?: Array<{
      name: string;
      role: string;
      mentioned: boolean;
      advantages?: string[];
      disadvantages?: string[];
    }>;
    comparison?: Record<string, unknown>;
  } | null;
}

interface Report {
  task: {
    project_id: string;
    model: string;
    status: string;
  };
  project: {
    id: string;
    name: string;
    industry: string;
  };
  visibility_score: number;
  brand_mentions: number;
  total_questions: number;
  answers: Answer[];
  recommendations: string[];
}

interface ComparisonAnswer {
  id: string;
  question: string;
  answer: string;
  model: string;
  brand_mentioned: boolean;
  sentiment: string;
  rank_position: number | null;
  analysis: Record<string, unknown> | null;
}

interface ModelComparisonResult {
  model: string;
  status: string;
  answers: ComparisonAnswer[];
  score: number;
}

interface ComparisonReport {
  task: Report;
  results: ModelComparisonResult[];
}

interface Competitor {
  id: string;
  name: string;
  mention_count: number;
  rank_position: number;
  advantages: string;
  analysis: Record<string, unknown> | null;
}

const STATUS_CONFIG: Record<string, { label: string; color: string }> = {
  pending: { label: '待处理', color: 'bg-gray-100 text-gray-600 border-gray-200' },
  running: { label: '运行中', color: 'bg-blue-100 text-blue-700 border-blue-200' },
  completed: { label: '已完成', color: 'bg-green-100 text-green-700 border-green-200' },
  failed: { label: '失败', color: 'bg-red-100 text-red-700 border-red-200' },
};

function getScoreColor(score: number): string {
  if (score >= 60) return 'text-green-600';
  if (score >= 30) return 'text-yellow-600';
  return 'text-red-600';
}

function getScoreBg(score: number): string {
  if (score >= 60) return 'bg-green-50 border-green-200';
  if (score >= 30) return 'bg-yellow-50 border-yellow-200';
  return 'bg-red-50 border-red-200';
}

function getSentimentIcon(sentiment: string) {
  switch (sentiment) {
    case 'positive': return <ThumbsUp size={14} className="text-green-600" />;
    case 'negative': return <ThumbsDown size={14} className="text-red-600" />;
    default: return <Minus size={14} className="text-gray-400" />;
  }
}

function getSentimentLabel(sentiment: string): string {
  switch (sentiment) {
    case 'positive': return '正面';
    case 'negative': return '负面';
    default: return '中性';
  }
}

export default function TaskDetailPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  const [report, setReport] = useState<Report | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [competitors, setCompetitors] = useState<Competitor[]>([]);
  const [competitorsLoading, setCompetitorsLoading] = useState(false);
  const [showCompetitors, setShowCompetitors] = useState(false);
  const [comparison, setComparison] = useState<ComparisonReport | null>(null);
  const [comparisonLoading, setComparisonLoading] = useState(false);
  const [showComparison, setShowComparison] = useState(false);

  useEffect(() => {
      apiRequest<Report>(`/tasks/${id}/report`)
      .then((res) => setReport(res.data))
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : '加载报告失败';
        setError(message);
      })
      .finally(() => setLoading(false));
  }, [id]);

  const isMultiModel = report?.task.model?.includes(',');

  const loadComparison = async () => {
    if (comparisonLoading) return;
    if (comparison) {
      setShowComparison(!showComparison);
      return;
    }
    setComparisonLoading(true);
    setShowComparison(true);
    try {
      const res = await apiRequest<ComparisonReport>(`/tasks/${id}/comparison`);
      setComparison(res.data);
    } catch {
      setComparison(null);
    } finally {
      setComparisonLoading(false);
    }
  };

  const loadCompetitors = async () => {
    if (competitorsLoading || competitors.length > 0) {
      setShowCompetitors(!showCompetitors);
      return;
    }
    setCompetitorsLoading(true);
    setShowCompetitors(true);
    try {
      const projectId = report?.project.id;
      if (!projectId) return;
      const res = await apiRequest<Competitor[] | { competitors: Competitor[] }>(`/projects/${projectId}/competitors`);
      const data = Array.isArray(res.data) ? res.data : (res.data as { competitors: Competitor[] }).competitors || [];
      setCompetitors(data);
    } catch {
      setCompetitors([]);
    } finally {
      setCompetitorsLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="animate-pulse space-y-6">
        <div className="h-5 w-24 bg-gray-200 rounded" />
        <div className="h-8 w-48 bg-gray-200 rounded" />
        <div className="h-40 bg-gray-100 rounded-card" />
        <div className="h-60 bg-gray-100 rounded-card" />
      </div>
    );
  }

  if (error) {
    return (
      <div>
        <Link href="/tasks" className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 transition-colors mb-6">
          <ArrowLeft size={16} />
          返回任务列表
        </Link>
        <div className="p-4 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">{error}</div>
      </div>
    );
  }

  if (!report) {
    return (
      <div>
        <Link href="/tasks" className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 transition-colors mb-6">
          <ArrowLeft size={16} />
          返回任务列表
        </Link>
        <div className="p-4 rounded-btn bg-yellow-50 border border-yellow-200 text-sm text-yellow-700">报告不存在</div>
      </div>
    );
  }

  const statusCfg = STATUS_CONFIG[report.task.status] || { label: report.task.status, color: 'bg-gray-100 text-gray-600 border-gray-200' };

  return (
    <div>
      <Link href="/tasks" className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 transition-colors mb-6">
        <ArrowLeft size={16} />
        返回任务列表
      </Link>

      <div className="bg-white rounded-card shadow-sm border border-gray-100 p-6 mb-6">
        <div className="flex items-start justify-between">
          <div>
            <div className="flex items-center gap-3 mb-2">
             <h1 className="text-2xl font-bold text-gray-900">{report.project.name}</h1>
              <span className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-medium border ${statusCfg.color}`}>
                {statusCfg.label}
              </span>
            </div>
            <div className="flex items-center gap-4 text-sm text-gray-500">
              <span className="inline-flex items-center gap-1.5">
                <Building2 size={14} />
                 {report.project.industry || '-'}
              </span>
              <span className="inline-flex items-center gap-1.5">
                <Brain size={14} />
                 {report.task.model}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div className={`rounded-card border p-6 mb-6 ${getScoreBg(report.visibility_score)}`}>
        <div className="flex items-center gap-8">
          <div className="text-center">
            <p className="text-sm text-gray-500 mb-1">AI 可见度分数</p>
             <p className={`text-5xl font-bold ${getScoreColor(report.visibility_score)}`}>
               {report.visibility_score}%
            </p>
            <p className="text-xs text-gray-400 mt-1">
              {report.visibility_score >= 60 ? '品牌在 AI 中有良好的曝光度' : report.visibility_score >= 30 ? '品牌可见度有待提升' : '品牌需要加大 AI 可见度优化'}
            </p>
          </div>
          <div className="flex gap-8">
            <div className="text-center">
              <p className="text-sm text-gray-500 mb-1">品牌提及</p>
              <p className="text-2xl font-semibold text-gray-900">{report.brand_mentions}</p>
              <p className="text-xs text-gray-400 mt-1">
                {report.brand_mentions > 0 ? `共被 AI 提及 ${report.brand_mentions} 次` : '尚未获得 AI 提及'}
              </p>
            </div>
            <div className="text-center">
              <p className="text-sm text-gray-500 mb-1">总问题数</p>
              <p className="text-2xl font-semibold text-gray-900">{report.total_questions}</p>
              <p className="text-xs text-gray-400 mt-1">覆盖 {report.total_questions} 个检测维度</p>
            </div>
          </div>
        </div>
      </div>

      {isMultiModel && (
        <div className="bg-white rounded-card shadow-sm border border-gray-100 mb-6 overflow-hidden">
          <button
            onClick={loadComparison}
            className="w-full px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors"
          >
            <div className="flex items-center gap-2">
              <BarChart3 size={18} className="text-purple-500" />
              <h2 className="text-base font-semibold text-gray-900">模型对比报告</h2>
              <span className="text-xs text-purple-500 bg-purple-50 px-2 py-0.5 rounded-full">多模型综合分析</span>
            </div>
            <div className="flex items-center gap-2 text-sm text-gray-400">
              {comparison && <span className="text-xs bg-purple-50 text-purple-700 px-2 py-0.5 rounded-full">{comparison.results.length} 个模型</span>}
              {showComparison ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
            </div>
          </button>
          {showComparison && (
            <div className="px-6 pb-5 border-t border-gray-100 pt-4">
              {comparisonLoading ? (
                <div className="flex items-center justify-center py-8">
                  <div className="animate-spin w-5 h-5 border-2 border-purple-500 border-t-transparent rounded-full" />
                  <span className="ml-2 text-sm text-gray-500">加载对比数据...</span>
                </div>
              ) : comparison && comparison.results.length > 0 ? (
                <div>
                  <p className="text-sm text-gray-500 mb-4">
                    以下对比显示了不同 AI 模型对 {report.project.name} 品牌的评估差异，
                    帮助您了解各模型的关注点和判断标准。
                  </p>
                  <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-gray-200">
                        <th className="text-left py-2 pr-4 font-semibold text-gray-700">模型</th>
                        <th className="text-center py-2 px-4 font-semibold text-gray-700">可见度分数</th>
                        <th className="text-center py-2 px-4 font-semibold text-gray-700">品牌提及</th>
                        <th className="text-center py-2 px-4 font-semibold text-gray-700">回答数</th>
                        <th className="text-left py-2 pl-4 font-semibold text-gray-700">关键发现</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                      {comparison.results.map((r) => {
                        const positiveCount = r.answers.filter(a => a.sentiment === 'positive').length;
                        const negativeCount = r.answers.filter(a => a.sentiment === 'negative').length;
                        const brandMentionedCount = r.answers.filter(a => a.brand_mentioned).length;
                        const scoreColor = r.score >= 60 ? 'text-green-600' : r.score >= 30 ? 'text-yellow-600' : 'text-red-600';
                        const topRankCount = r.answers.filter(a => a.rank_position != null && a.rank_position <= 3).length;
                        return (
                          <tr key={r.model}>
                            <td className="py-3 pr-4">
                              <span className="font-medium text-gray-900">{r.model}</span>
                            </td>
                            <td className={`py-3 px-4 text-center font-bold text-lg ${scoreColor}`}>
                              <div>{Math.round(r.score)}%</div>
                              <div className="w-full bg-gray-100 rounded-full h-1.5 mt-1">
                                <div className={`h-full rounded-full ${r.score >= 60 ? 'bg-green-500' : r.score >= 30 ? 'bg-yellow-500' : 'bg-red-500'}`} style={{ width: `${Math.round(r.score)}%` }} />
                              </div>
                            </td>
                            <td className="py-3 px-4 text-center text-gray-700">
                              {brandMentionedCount}/{r.answers.length}
                            </td>
                            <td className="py-3 px-4 text-center text-gray-700">
                              {r.answers.length}
                            </td>
                            <td className="py-3 pl-4">
                              <div className="flex flex-wrap gap-1.5">
                                {positiveCount > 0 && (
                                  <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-green-50 text-green-700 border border-green-200">
                                    <ThumbsUp size={12} />{positiveCount} 正面
                                  </span>
                                )}
                                {negativeCount > 0 && (
                                  <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-red-50 text-red-700 border border-red-200">
                                    <ThumbsDown size={12} />{negativeCount} 负面
                                  </span>
                                )}
                                {topRankCount > 0 && (
                                  <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-amber-50 text-amber-700 border border-amber-200">
                                    <Trophy size={12} /> 前3排名: {topRankCount}
                                  </span>
                                )}
                              </div>
                            </td>
                          </tr>
                        );
                      })}
                    </tbody>
                  </table>
                  </div>
                  {(() => {
                    const best = comparison.results.reduce((a, b) => a.score > b.score ? a : b);
                    return (
                      <div className="mt-4 p-4 bg-purple-50 rounded-lg border border-purple-100">
                        <p className="text-sm font-medium text-purple-800">
                          推荐模型: {best.model}（评分 {Math.round(best.score)}%）
                        </p>
                        <p className="text-xs text-purple-600 mt-1">
                          该模型对 {report.project.name} 的品牌识别度最高，建议优先参考其分析结果。
                        </p>
                      </div>
                    );
                  })()}
                </div>
              ) : (
                <div className="text-center py-8 text-sm text-gray-400">暂无对比数据</div>
              )}
            </div>
          )}
        </div>
      )}

      <div className="bg-white rounded-card shadow-sm border border-gray-100 mb-6 overflow-hidden">
        <button
          onClick={loadCompetitors}
          className="w-full px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors"
        >
          <div className="flex items-center gap-2">
            <Trophy size={18} className="text-amber-500" />
            <h2 className="text-base font-semibold text-gray-900">竞争分析</h2>
            <span className="text-xs text-amber-600 bg-amber-50 px-2 py-0.5 rounded-full">竞品监测</span>
          </div>
          <div className="flex items-center gap-2 text-sm text-gray-400">
            {competitors.length > 0 && <span className="text-xs bg-amber-50 text-amber-700 px-2 py-0.5 rounded-full">{competitors.length} 个竞品</span>}
            {showCompetitors ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
          </div>
        </button>
        {showCompetitors && (
          <div className="px-6 pb-5 border-t border-gray-100 pt-4">
            {competitorsLoading ? (
              <div className="flex items-center justify-center py-8">
                <div className="animate-spin w-5 h-5 border-2 border-primary border-t-transparent rounded-full" />
                <span className="ml-2 text-sm text-gray-500">加载竞品数据...</span>
              </div>
            ) : competitors.length === 0 ? (
              <div className="text-center py-8 text-sm text-gray-400">
                <TrendingUp size={24} className="mx-auto mb-2 text-gray-300" />
                <p className="font-medium text-gray-500 mb-1">暂未发现竞品</p>
                <p>AI 模型在回答中尚未提及 {report.project.name} 品牌相关的竞争品牌，这是一个积极信号。</p>
              </div>
            ) : (
              <div>
                <p className="text-sm text-gray-500 mb-4">
                  AI 模型在评估过程中识别到以下竞品品牌，这反映了市场认知中的主要竞争对手。
                </p>
                <div className="space-y-3">
                {competitors.map((comp, i) => (
                  <div key={comp.name} className="flex items-center gap-4 p-4 rounded-lg bg-gray-50 border border-gray-100 hover:border-amber-200 transition-colors">
                    <div className={`w-8 h-8 rounded-full flex items-center justify-center shrink-0 text-sm font-bold ${
                      i === 0 ? 'bg-amber-500 text-white' :
                      i === 1 ? 'bg-gray-400 text-white' :
                      i === 2 ? 'bg-amber-200 text-amber-800' :
                      'bg-gray-200 text-gray-600'
                    }`}>
                      {i + 1}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-semibold text-gray-900">{comp.name}</p>
                      <div className="flex items-center gap-3 mt-1 text-xs text-gray-500">
                        <span className="inline-flex items-center gap-1">
                          <span className="w-2 h-2 rounded-full bg-blue-400" />
                          提及 {comp.mention_count} 次
                        </span>
                        {comp.rank_position > 0 && <span>最佳排名 #{comp.rank_position}</span>}
                        {comp.advantages && <span className="text-green-600 truncate max-w-[200px]">优势: {comp.advantages}</span>}
                      </div>
                    </div>
                    <div className="shrink-0 text-right">
                      <span className="text-xs text-gray-400">关注度</span>
                      <div className="w-16 bg-gray-200 rounded-full h-1.5 mt-1">
                        <div className="bg-amber-500 h-full rounded-full" style={{ width: `${Math.min(100, comp.mention_count * 15)}%` }} />
                      </div>
                    </div>
                  </div>
                ))}
                </div>
                <div className="mt-4 p-4 bg-amber-50 rounded-lg border border-amber-100">
                  <p className="text-sm font-medium text-amber-800">竞争洞察</p>
                  <p className="text-xs text-amber-600 mt-1">
                    {competitors.length > 0 
                      ? `检测到 ${competitors.length} 个竞品品牌正在与您抢占 AI 可见度，重点关注排名第一的 ${competitors[0]?.name}。`
                      : '当前 AI 生态中您的品牌具有较好的独立认知度。'}
                  </p>
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      <div className="bg-white rounded-card shadow-sm border border-gray-100 mb-6 overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <div>
            <h2 className="text-base font-semibold text-gray-900">AI 回答分析</h2>
            <p className="text-xs text-gray-500 mt-0.5">每个问题对应的 AI 回答及品牌表现评估</p>
          </div>
          <span className="text-xs text-gray-400 bg-gray-50 px-2 py-1 rounded-full">{report.answers.length} 条回答</span>
        </div>
        <div className="divide-y divide-gray-50">
          {report.answers.map((answer, i) => {
            const entities = answer.analysis?.entities || [];
            const mentionedCompetitors = entities.filter(e => e.role === 'competitor' && e.mentioned);
            const brandMentioned = answer.mentions_brand;
            const brandName = entities.find(e => e.role === 'target')?.name || report.project.name;
            return (
            <div key={i} className="px-6 py-5 hover:bg-gray-50/50 transition-colors">
              <div className="flex items-start gap-3">
                <div className={`shrink-0 w-2 h-2 rounded-full mt-2 ${
                  brandMentioned ? 'bg-green-500' : 'bg-gray-300'
                }`} title={brandMentioned ? '品牌被提及' : '品牌未被提及'} />
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1.5">
                    <span className="text-xs font-medium text-gray-400 bg-gray-100 px-2 py-0.5 rounded">Q{i + 1}</span>
                    <p className="text-sm font-medium text-gray-900">{answer.question}</p>
                  </div>
                  <p className="text-sm text-gray-600 mb-3 line-clamp-2">{answer.answer_summary}</p>
                  <div className="flex items-center gap-4 text-xs flex-wrap">
                    <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium ${
                      brandMentioned
                        ? 'bg-green-50 text-green-700 border border-green-200'
                        : 'bg-gray-50 text-gray-500 border border-gray-200'
                    }`}>
                      {getSentimentIcon(answer.sentiment)}
                      {brandMentioned ? `${brandName} 被提及` : '未提及品牌'}
                    </span>
                    <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium ${
                      answer.sentiment === 'positive' ? 'bg-green-50 text-green-700 border border-green-200' :
                      answer.sentiment === 'negative' ? 'bg-red-50 text-red-700 border border-red-200' :
                      'bg-gray-50 text-gray-500 border border-gray-200'
                    }`}>
                      情感: {getSentimentLabel(answer.sentiment)}
                    </span>
                    <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-blue-50 text-blue-700 border border-blue-200">
                      排名: #{answer.rank_position || '-'}
                    </span>
                    {mentionedCompetitors.length > 0 && (
                      <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-amber-50 text-amber-700 border border-amber-200">
                        竞品提及: {mentionedCompetitors.map(c => c.name).join('、')}
                      </span>
                    )}
                  </div>
                  <p className="text-xs text-gray-400 mt-2.5 leading-relaxed">
                    AI 评估: {answer.answer_summary?.slice(0, 150)}{(answer.answer_summary?.length ?? 0) > 150 ? '...' : ''}
                  </p>
                </div>
              </div>
            </div>
            );
          })}
          {report.answers.length === 0 && (
            <div className="px-6 py-12 text-center text-sm text-gray-500">
              <FlaskConical size={24} className="mx-auto mb-2 text-gray-300" />
              <p className="font-medium mb-1">暂无回答数据</p>
              <p>完成 AI 检测任务后，详细的回答分析将在此展示。</p>
            </div>
          )}
        </div>
      </div>

      {report.recommendations && report.recommendations.length > 0 && (
        <div className="bg-white rounded-card shadow-sm border border-gray-100 overflow-hidden">
          <div className="px-6 py-4 border-b border-gray-100">
            <div className="flex items-center gap-2">
              <TrendingUp size={18} className="text-green-500" />
              <h2 className="text-base font-semibold text-gray-900">优化建议</h2>
              <span className="text-xs text-green-600 bg-green-50 px-2 py-0.5 rounded-full">提升指南</span>
            </div>
            <p className="text-xs text-gray-500 mt-1">
              基于 AI 检测结果，为您量身定制的品牌可见度优化策略
            </p>
          </div>
          <div className="px-6 py-5 space-y-4">
            {report.recommendations.map((rec, i) => {
              const colors = [
                { bg: 'bg-green-50', border: 'border-green-200', text: 'text-green-700', icon: 'bg-green-500' },
                { bg: 'bg-blue-50', border: 'border-blue-200', text: 'text-blue-700', icon: 'bg-blue-500' },
                { bg: 'bg-purple-50', border: 'border-purple-200', text: 'text-purple-700', icon: 'bg-purple-500' },
                { bg: 'bg-amber-50', border: 'border-amber-200', text: 'text-amber-700', icon: 'bg-amber-500' },
              ];
              const c = colors[i % colors.length];
              return (
                <div key={i} className={`flex items-start gap-4 p-4 rounded-lg ${c.bg} border ${c.border}`}>
                  <div className={`shrink-0 w-7 h-7 rounded-full ${c.icon} text-white text-xs font-bold flex items-center justify-center mt-0.5`}>
                    {i + 1}
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className={`text-sm font-medium ${c.text}`}>建议 {i + 1}</p>
                    <p className="text-sm text-gray-700 mt-1 leading-relaxed">{rec}</p>
                  </div>
                  <div className="shrink-0">
                    <ThumbsUp size={16} className={c.text} />
                  </div>
                </div>
              );
            })}
          </div>
          <div className="px-6 py-4 bg-gray-50 border-t border-gray-100">
            <p className="text-xs text-gray-500 text-center">
              这些建议基于 {report.total_questions} 个维度的 AI 检测结果生成，建议每季度复查并调整策略。
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
