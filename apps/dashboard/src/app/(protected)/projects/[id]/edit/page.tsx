'use client';
import { useState, FormEvent, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { apiRequest } from '@/lib/api';
import { ArrowLeft } from 'lucide-react';

interface Project {
  id: string;
  name: string;
  website: string;
  industry: string;
  status: string;
  description: string;
  keywords: string | string[];
  service_area: string;
  target_users: string;
  created_at: string;
  updated_at: string;
}

const INDUSTRY_OPTIONS = [
  { value: '', label: '请选择行业' },
  { value: 'SaaS / 软件服务', label: 'SaaS / 软件服务' },
  { value: '人工智能 / AI', label: '人工智能 / AI' },
  { value: '电商 / 零售', label: '电商 / 零售' },
  { value: '在线教育 / 知识付费', label: '在线教育 / 知识付费' },
  { value: '金融科技 / 支付', label: '金融科技 / 支付' },
  { value: '医疗健康 / 生物科技', label: '医疗健康 / 生物科技' },
  { value: '制造 / 工业', label: '制造 / 工业' },
  { value: '游戏 / 娱乐', label: '游戏 / 娱乐' },
  { value: '媒体 / 内容平台', label: '媒体 / 内容平台' },
  { value: '出行 / 物流', label: '出行 / 物流' },
  { value: '房产 / 家居', label: '房产 / 家居' },
  { value: '餐饮 / 本地生活', label: '餐饮 / 本地生活' },
  { value: '新能源汽车', label: '新能源汽车' },
  { value: '企业服务 / B2B', label: '企业服务 / B2B' },
  { value: '法律 / 咨询', label: '法律 / 咨询' },
  { value: '人力资源 / 招聘', label: '人力资源 / 招聘' },
  { value: '农业 / 食品', label: '农业 / 食品' },
  { value: '其他', label: '其他' },
];

function SkeletonForm() {
  return (
    <div className="bg-white rounded-card shadow-sm border border-gray-100 p-6 animate-pulse space-y-5">
      <div className="space-y-1">
        <div className="h-4 w-20 bg-gray-200 rounded" />
        <div className="h-10 w-full bg-gray-100 rounded-btn" />
      </div>
      <div className="space-y-1">
        <div className="h-4 w-12 bg-gray-200 rounded" />
        <div className="h-10 w-full bg-gray-100 rounded-btn" />
      </div>
      <div className="space-y-1">
        <div className="h-4 w-12 bg-gray-200 rounded" />
        <div className="h-10 w-full bg-gray-100 rounded-btn" />
      </div>
      <div className="space-y-1">
        <div className="h-4 w-16 bg-gray-200 rounded" />
        <div className="h-20 w-full bg-gray-100 rounded-btn" />
      </div>
      <div className="space-y-1">
        <div className="h-4 w-20 bg-gray-200 rounded" />
        <div className="h-10 w-full bg-gray-100 rounded-btn" />
      </div>
      <div className="flex gap-3 pt-2">
        <div className="h-10 w-24 bg-gray-200 rounded-btn" />
        <div className="h-10 w-20 bg-gray-100 rounded-btn" />
      </div>
    </div>
  );
}

export default function EditProjectPage() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;

  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState('');
  const [name, setName] = useState('');
  const [website, setWebsite] = useState('');
  const [industry, setIndustry] = useState('');
  const [description, setDescription] = useState('');
  const [keywords, setKeywords] = useState('');
  const [serviceArea, setServiceArea] = useState('');
  const [targetUsers, setTargetUsers] = useState('');
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    apiRequest<Project>(`/projects/${id}`)
      .then((res) => {
        const p = res.data;
        setName(p.name || '');
        setWebsite(p.website || '');
        setIndustry(p.industry || '');
        setDescription(p.description || '');
        setKeywords(Array.isArray(p.keywords) ? (p.keywords as string[]).join(', ') : (p.keywords || ''));
        setServiceArea(p.service_area || '');
        setTargetUsers(p.target_users || '');
      })
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : '加载项目信息失败';
        setFetchError(message);
      })
      .finally(() => setLoading(false));
  }, [id]);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');

    if (!name.trim()) {
      setError('项目名称不能为空');
      return;
    }
    if (!industry) {
      setError('请选择行业');
      return;
    }

    setSubmitting(true);
    try {
      await apiRequest(`/projects/${id}`, {
        method: 'PUT',
        body: JSON.stringify({
          name: name.trim(),
          website: website.trim() || undefined,
          industry,
          description: description.trim() || undefined,
          keywords: keywords ? keywords.split(',').map((k: string) => k.trim()).filter(Boolean) : undefined,
          service_area: serviceArea.trim() || undefined,
          target_users: targetUsers.trim() || undefined,
        }),
      });
      router.push(`/projects/${id}`);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : '更新项目失败';
      setError(message);
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div>
        <div className="h-5 w-24 bg-gray-200 rounded animate-pulse mb-6" />
        <div className="max-w-2xl">
          <SkeletonForm />
        </div>
      </div>
    );
  }

  if (fetchError) {
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
          {fetchError}
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <Link
          href={`/projects/${id}`}
          className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 transition-colors"
        >
          <ArrowLeft size={16} />
          返回项目详情
        </Link>
      </div>

      <h1 className="text-2xl font-bold text-gray-900 mb-6">编辑项目</h1>

      <div className="max-w-2xl">
        {error && (
          <div className="mb-6 p-4 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="bg-white rounded-card shadow-sm border border-gray-100 p-6 space-y-5">
          <div>
            <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
              项目名称 <span className="text-red-500">*</span>
            </label>
            <input
              id="name"
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors"
            />
          </div>

          <div>
            <label htmlFor="website" className="block text-sm font-medium text-gray-700 mb-1">
              网站
            </label>
            <input
              id="website"
              type="text"
              value={website}
              onChange={(e) => setWebsite(e.target.value)}
              placeholder="例如: example.com"
              className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors"
            />
          </div>

          <div>
            <label htmlFor="industry" className="block text-sm font-medium text-gray-700 mb-1">
              行业 <span className="text-red-500">*</span>
            </label>
            <select
              id="industry"
              required
              value={industry}
              onChange={(e) => setIndustry(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors bg-white"
            >
              {INDUSTRY_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value} disabled={opt.value === ''}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
              描述
            </label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors resize-none"
            />
          </div>

          <div>
            <label htmlFor="keywords" className="block text-sm font-medium text-gray-700 mb-1">
              关键词
            </label>
            <input
              id="keywords"
              type="text"
              value={keywords}
              onChange={(e) => setKeywords(e.target.value)}
              placeholder="多个关键词用逗号分隔，例如: AI搜索, 品牌可见度"
              className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors"
            />
          </div>

          <div>
            <label htmlFor="serviceArea" className="block text-sm font-medium text-gray-700 mb-1">
              服务区域
            </label>
            <input
              id="serviceArea"
              type="text"
              value={serviceArea}
              onChange={(e) => setServiceArea(e.target.value)}
              placeholder="例如: 全国, 华东"
              className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors"
            />
          </div>

          <div>
            <label htmlFor="targetUsers" className="block text-sm font-medium text-gray-700 mb-1">
              目标用户
            </label>
            <input
              id="targetUsers"
              type="text"
              value={targetUsers}
              onChange={(e) => setTargetUsers(e.target.value)}
              placeholder="例如: 企业HR, 中小型企业"
              className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors"
            />
          </div>

          <div className="flex items-center gap-3 pt-2">
            <button
              type="submit"
              disabled={submitting}
              className="px-6 py-2 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {submitting ? '保存中...' : '保存修改'}
            </button>
            <Link
              href={`/projects/${id}`}
              className="px-6 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-btn hover:bg-gray-50 transition-colors"
            >
              取消
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
