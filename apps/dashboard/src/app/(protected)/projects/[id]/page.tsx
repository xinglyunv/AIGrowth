'use client';
import { useState, useEffect, FormEvent } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { apiRequest } from '@/lib/api';
import { ArrowLeft, Edit, FlaskConical, ExternalLink, Globe, Building2, Tag, MapPin, Users, FileText, Calendar, ChevronRight, Search, LucideIcon, Save, X, BookOpen, Package, HelpCircle, Star, Briefcase } from 'lucide-react';

interface Project {
  id: string;
  name: string;
  website: string;
  industry: string;
  status: string;
  description: string;
  keywords: string;
  service_area: string;
  target_users: string;
  created_at: string;
  updated_at: string;
  brand_intro?: string;
  product_intro?: string;
  service_intro?: string;
  faq?: string;
  advantages?: string;
  cases?: string;
}

interface Task {
  id: string;
  project_name: string;
  model: string;
  status: string;
  questions_count: number;
  completed_count: number;
  created_at: string;
}

interface BrandInfo {
  brand_intro: string;
  product_intro: string;
  service_intro: string;
  faq: string;
  advantages: string;
  cases: string;
}

const STATUS_LABELS: Record<string, string> = {
  active: '运行中',
  inactive: '已暂停',
  draft: '草稿',
};
const STATUS_COLORS: Record<string, string> = {
  active: 'bg-green-100 text-green-700 border-green-200',
  inactive: 'bg-gray-100 text-gray-600 border-gray-200',
  draft: 'bg-yellow-100 text-yellow-700 border-yellow-200',
};

const TASK_STATUS_CONFIG: Record<string, { label: string; color: string }> = {
  pending: { label: '待处理', color: 'bg-gray-100 text-gray-600 border-gray-200' },
  running: { label: '运行中', color: 'bg-blue-100 text-blue-700 border-blue-200' },
  completed: { label: '已完成', color: 'bg-green-100 text-green-700 border-green-200' },
  failed: { label: '失败', color: 'bg-red-100 text-red-700 border-red-200' },
};

const BRAND_FIELDS: { key: keyof BrandInfo; label: string; icon: LucideIcon; placeholder: string }[] = [
  { key: 'brand_intro', label: '品牌介绍', icon: BookOpen, placeholder: '品牌历史、定位、愿景等' },
  { key: 'product_intro', label: '产品介绍', icon: Package, placeholder: '核心产品功能、特点等' },
  { key: 'service_intro', label: '服务介绍', icon: Briefcase, placeholder: '提供的服务内容、服务流程等' },
  { key: 'faq', label: '常见问题', icon: HelpCircle, placeholder: '客户常见问题与解答，用换行分隔' },
  { key: 'advantages', label: '核心优势', icon: Star, placeholder: '品牌核心竞争优势，用换行分隔' },
  { key: 'cases', label: '客户案例', icon: FileText, placeholder: '典型客户案例，用换行分隔' },
];

function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const d = new Date(dateStr);
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' });
}

function InfoRow({ label, value, icon: Icon }: { label: string; value: string | undefined; icon: LucideIcon }) {
  if (!value) return null;
  return (
    <div className="flex items-start gap-3">
      <div className="mt-0.5 shrink-0">
        <Icon size={16} className="text-gray-400" />
      </div>
      <div>
        <p className="text-xs text-gray-500 mb-0.5">{label}</p>
        <p className="text-sm text-gray-900">{value}</p>
      </div>
    </div>
  );
}

function SkeletonCard() {
  return (
    <div className="bg-white rounded-card shadow-sm border border-gray-100 p-6 animate-pulse">
      <div className="h-7 w-48 bg-gray-200 rounded mb-6" />
      <div className="space-y-4">
        <div className="h-5 w-full bg-gray-100 rounded" />
        <div className="h-5 w-3/4 bg-gray-100 rounded" />
        <div className="h-5 w-1/2 bg-gray-100 rounded" />
        <div className="h-5 w-2/3 bg-gray-100 rounded" />
      </div>
    </div>
  );
}

export default function ProjectDetailPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [tasks, setTasks] = useState<Task[]>([]);
  const [tasksLoading, setTasksLoading] = useState(false);

  const [editingBrand, setEditingBrand] = useState(false);
  const [brandForm, setBrandForm] = useState<BrandInfo>({
    brand_intro: '', product_intro: '', service_intro: '', faq: '', advantages: '', cases: '',
  });
  const [brandSaving, setBrandSaving] = useState(false);
  const [brandError, setBrandError] = useState('');
  const [brandSuccess, setBrandSuccess] = useState('');

  useEffect(() => {
    apiRequest<Project>(`/projects/${id}`)
      .then((res) => {
        setProject(res.data);
        setBrandForm({
          brand_intro: res.data.brand_intro || '',
          product_intro: res.data.product_intro || '',
          service_intro: res.data.service_intro || '',
          faq: res.data.faq || '',
          advantages: res.data.advantages || '',
          cases: res.data.cases || '',
        });
      })
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : '加载项目详情失败';
        setError(message);
      })
      .finally(() => setLoading(false));
  }, [id]);

  useEffect(() => {
    apiRequest<Task[]>(`/tasks?project_id=${id}`)
      .then((res) => setTasks(res.data))
      .catch(() => {})
      .finally(() => setTasksLoading(false));
  }, [id]);

  const handleBrandSave = async (e: FormEvent) => {
    e.preventDefault();
    setBrandError('');
    setBrandSuccess('');
    setBrandSaving(true);
    try {
      await apiRequest(`/projects/${id}`, {
        method: 'PUT',
        body: JSON.stringify(brandForm),
      });
      setProject((prev) => prev ? { ...prev, ...brandForm } : prev);
      setBrandSuccess('品牌资料已保存');
      setEditingBrand(false);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : '保存失败';
      setBrandError(message);
    } finally {
      setBrandSaving(false);
    }
  };

  if (loading) {
    return (
      <div>
        <div className="h-5 w-24 bg-gray-200 rounded animate-pulse mb-6" />
        <SkeletonCard />
      </div>
    );
  }

  if (error) {
    return (
      <div>
        <Link
          href="/projects"
          className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 transition-colors mb-6"
        >
          <ArrowLeft size={16} />
          返回项目列表
        </Link>
        <div className="p-4 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
          {error}
        </div>
      </div>
    );
  }

  if (!project) {
    return (
      <div>
        <Link
          href="/projects"
          className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 transition-colors mb-6"
        >
          <ArrowLeft size={16} />
          返回项目列表
        </Link>
        <div className="p-4 rounded-btn bg-yellow-50 border border-yellow-200 text-sm text-yellow-700">
          项目不存在
        </div>
      </div>
    );
  }

  const keywordsList: string[] = Array.isArray(project.keywords)
    ? project.keywords
    : project.keywords
      ? String(project.keywords).split(',').map((k) => k.trim()).filter(Boolean)
      : [];

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <Link
          href="/projects"
          className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 transition-colors"
        >
          <ArrowLeft size={16} />
          返回项目列表
        </Link>
      </div>

      <div className="flex items-start justify-between mb-6">
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-bold text-gray-900">{project.name}</h1>
          <span
            className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-medium border ${
              STATUS_COLORS[project.status] || STATUS_COLORS.draft
            }`}
          >
            {STATUS_LABELS[project.status] || project.status}
          </span>
        </div>
        <div className="flex items-center gap-3">
          <Link
            href={`/projects/${id}/edit`}
            className="inline-flex items-center gap-2 px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-btn hover:bg-gray-50 transition-colors"
          >
            <Edit size={16} />
            编辑
          </Link>
          <Link
            href={`/projects/${id}/tasks/new`}
            className="inline-flex items-center gap-2 px-4 py-2 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 transition-colors"
          >
            <FlaskConical size={16} />
            创建 AI 检测
          </Link>
        </div>
      </div>

      <div className="bg-white rounded-card shadow-sm border border-gray-100 p-6 space-y-5">
        <InfoRow label="网站" value={project.website} icon={Globe} />
        <InfoRow label="行业" value={project.industry} icon={Building2} />
        <InfoRow label="描述" value={project.description} icon={FileText} />
        <InfoRow label="服务区域" value={project.service_area} icon={MapPin} />
        <InfoRow label="目标用户" value={project.target_users} icon={Users} />

        {keywordsList.length > 0 && (
          <div className="flex items-start gap-3">
            <div className="mt-0.5 shrink-0">
              <Tag size={16} className="text-gray-400" />
            </div>
            <div>
              <p className="text-xs text-gray-500 mb-1.5">关键词</p>
              <div className="flex flex-wrap gap-2">
                {keywordsList.map((kw, i) => (
                  <span
                    key={i}
                    className="inline-block px-2.5 py-1 bg-blue-50 text-blue-700 text-xs font-medium rounded-full border border-blue-100"
                  >
                    {kw}
                  </span>
                ))}
              </div>
            </div>
          </div>
        )}

        <div className="flex items-start gap-3">
          <div className="mt-0.5 shrink-0">
            <Calendar size={16} className="text-gray-400" />
          </div>
          <div>
            <p className="text-xs text-gray-500 mb-0.5">创建时间</p>
            <p className="text-sm text-gray-900">{formatDate(project.created_at)}</p>
          </div>
        </div>

        {project.website && (
          <div className="pt-3 border-t border-gray-100">
            <a
              href={project.website.startsWith('http') ? project.website : `https://${project.website}`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1.5 text-sm text-primary hover:underline"
            >
              <ExternalLink size={14} />
              访问网站
            </a>
          </div>
        )}
      </div>

      <div className="mt-8 bg-white rounded-card shadow-sm border border-gray-100 p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-base font-semibold text-gray-900">品牌资料</h2>
          <button
            onClick={() => {
              if (editingBrand) {
                setBrandForm({
                  brand_intro: project.brand_intro || '',
                  product_intro: project.product_intro || '',
                  service_intro: project.service_intro || '',
                  faq: project.faq || '',
                  advantages: project.advantages || '',
                  cases: project.cases || '',
                });
                setBrandError('');
                setBrandSuccess('');
              }
              setEditingBrand(!editingBrand);
            }}
            className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:text-blue-700"
          >
            {editingBrand ? <X size={14} /> : <Edit size={14} />}
            {editingBrand ? '取消' : '编辑'}
          </button>
        </div>

        {brandError && (
          <div className="mb-4 p-3 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
            {brandError}
          </div>
        )}
        {brandSuccess && (
          <div className="mb-4 p-3 rounded-btn bg-green-50 border border-green-200 text-sm text-green-700">
            {brandSuccess}
          </div>
        )}

        {editingBrand ? (
          <form onSubmit={handleBrandSave} className="space-y-4">
            {BRAND_FIELDS.map((field) => (
              <div key={field.key}>
                <label className="block text-sm font-medium text-gray-700 mb-1 flex items-center gap-1.5">
                  <field.icon size={14} className="text-gray-400" />
                  {field.label}
                </label>
                <textarea
                  value={brandForm[field.key]}
                  onChange={(e) => setBrandForm((prev) => ({ ...prev, [field.key]: e.target.value }))}
                  placeholder={field.placeholder}
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors resize-none"
                />
              </div>
            ))}
            <button
              type="submit"
              disabled={brandSaving}
              className="inline-flex items-center gap-2 px-4 py-2 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <Save size={14} />
              {brandSaving ? '保存中...' : '保存品牌资料'}
            </button>
          </form>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {BRAND_FIELDS.map((field) => {
              const value = project[field.key as keyof Project] as string | undefined;
              return (
                <div key={field.key} className={`p-4 rounded-lg border border-gray-100 ${!value ? 'bg-gray-50' : ''}`}>
                  <div className="flex items-center gap-1.5 mb-2">
                    <field.icon size={14} className="text-gray-400" />
                    <span className="text-xs font-medium text-gray-500">{field.label}</span>
                  </div>
                  {value ? (
                    <p className="text-sm text-gray-900 whitespace-pre-wrap">{value}</p>
                  ) : (
                    <p className="text-sm text-gray-400 italic">暂未填写</p>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>

      <div className="mt-8">
        <h2 className="text-base font-semibold text-gray-900 mb-4">AI 检测任务</h2>
        <div className="bg-white rounded-card shadow-sm border border-gray-100 overflow-hidden">
          {tasksLoading ? (
            <div className="p-6 text-center text-sm text-gray-500">加载中...</div>
          ) : tasks.length === 0 ? (
            <div className="p-6 text-center">
              <div className="flex flex-col items-center gap-2">
                <Search size={20} className="text-gray-300" />
                <p className="text-sm text-gray-500">暂无检测任务</p>
                <Link
                  href={`/projects/${id}/tasks/new`}
                  className="mt-1 inline-flex items-center gap-1.5 text-sm text-primary hover:underline"
                >
                  <FlaskConical size={14} />
                  创建首个检测任务
                </Link>
              </div>
            </div>
          ) : (
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50/50">
                  <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">模型</th>
                  <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">状态</th>
                  <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">问题数</th>
                  <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">创建时间</th>
                  <th className="w-10 px-6 py-3" />
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {tasks.map((task) => {
                  const tsc = TASK_STATUS_CONFIG[task.status] || { label: task.status, color: 'bg-gray-100 text-gray-600 border-gray-200' };
                  return (
                    <tr
                      key={task.id}
                      className="hover:bg-gray-50 transition-colors cursor-pointer"
                      onClick={() => router.push(`/tasks/${task.id}`)}
                    >
                      <td className="px-6 py-4 text-sm text-gray-600">{task.model}</td>
                      <td className="px-6 py-4">
                        <span className={`inline-block px-2.5 py-0.5 rounded-full text-xs font-medium border ${tsc.color}`}>
                          {tsc.label}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-600">{task.questions_count ?? '-'}</td>
                      <td className="px-6 py-4 text-sm text-gray-500">{formatDate(task.created_at)}</td>
                      <td className="px-6 py-4"><ChevronRight size={16} className="text-gray-300" /></td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}
