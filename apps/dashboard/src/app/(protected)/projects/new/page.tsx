'use client';
import { useState, FormEvent, KeyboardEvent } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { apiRequest } from '@/lib/api';
import { ArrowLeft, Loader2, X } from 'lucide-react';

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

export default function NewProjectPage() {
  const router = useRouter();

  const [name, setName] = useState('');
  const [industry, setIndustry] = useState('');
  const [description, setDescription] = useState('');
  const [keywords, setKeywords] = useState<string[]>([]);
  const [keywordInput, setKeywordInput] = useState('');
  const [brandIntro, setBrandIntro] = useState('');
  const [productIntro, setProductIntro] = useState('');
  const [serviceIntro, setServiceIntro] = useState('');

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState('');

  const addKeyword = (value: string) => {
    const trimmed = value.trim();
    if (trimmed && !keywords.includes(trimmed)) {
      setKeywords([...keywords, trimmed]);
    }
    setKeywordInput('');
  };

  const removeKeyword = (keyword: string) => {
    setKeywords(keywords.filter((k) => k !== keyword));
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      addKeyword(keywordInput);
    }
  };

  const validate = (): boolean => {
    const errs: Record<string, string> = {};
    if (!name.trim()) errs.name = '项目名称不能为空';
    if (keywords.length === 0) errs.keywords = '至少添加 1 个关键词';
    setErrors(errs);
    return Object.keys(errs).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setSubmitError('');
    if (!validate()) return;

    setSubmitting(true);
    try {
      const res = await apiRequest<{ id: string }>('/projects', {
        method: 'POST',
        body: JSON.stringify({
          name: name.trim(),
          industry: industry || undefined,
          description: description.trim() || undefined,
          keywords,
          brand_intro: brandIntro.trim() || undefined,
          product_intro: productIntro.trim() || undefined,
          service_intro: serviceIntro.trim() || undefined,
        }),
      });
      router.push(`/projects/${res.data.id}`);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : '创建项目失败';
      setSubmitError(message);
    } finally {
      setSubmitting(false);
    }
  };

  const inputClass = (field: string) =>
    `w-full px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors ${
      errors[field] ? 'border-red-300 bg-red-50' : 'border-gray-300'
    }`;

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

      <h1 className="text-2xl font-bold text-gray-900 mb-6">创建项目</h1>

      <div className="max-w-3xl">
        {submitError && (
          <div className="mb-6 p-4 rounded-lg bg-red-50 border border-red-200 text-sm text-red-700">
            {submitError}
          </div>
        )}

        <form onSubmit={handleSubmit} className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
                项目名称 <span className="text-red-500">*</span>
              </label>
              <input
                id="name"
                type="text"
                value={name}
                onChange={(e) => {
                  setName(e.target.value);
                  if (errors.name) setErrors({ ...errors, name: '' });
                }}
                placeholder="例如：XX 品牌 AI 可见度分析"
                className={inputClass('name')}
              />
              {errors.name && <p className="mt-1 text-xs text-red-600">{errors.name}</p>}
            </div>

            <div>
              <label htmlFor="industry" className="block text-sm font-medium text-gray-700 mb-1">
                行业 <span className="text-red-500">*</span>
              </label>
              <select
                id="industry"
                value={industry}
                onChange={(e) => setIndustry(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors bg-white"
              >
                {INDUSTRY_OPTIONS.map((opt) => (
                  <option key={opt.value} value={opt.value} disabled={opt.value === ''}>
                    {opt.label}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
                描述
              </label>
              <textarea
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="简要描述品牌业务范围和市场定位"
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors resize-none"
              />
              <p className="mt-1 text-xs text-gray-400 text-right">{description.length} 字</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                关键词 <span className="text-red-500">*</span>
              </label>
              <div
                className={`flex flex-wrap gap-2 p-2 border rounded-lg min-h-[42px] focus-within:ring-2 focus-within:ring-primary/20 focus-within:border-primary transition-colors ${
                  errors.keywords ? 'border-red-300 bg-red-50' : 'border-gray-300'
                }`}
              >
                {keywords.map((kw) => (
                  <span
                    key={kw}
                    className="inline-flex items-center gap-1 px-2 py-1 bg-blue-50 text-blue-700 text-xs font-medium rounded-full"
                  >
                    {kw}
                    <button type="button" onClick={() => removeKeyword(kw)} className="hover:text-blue-900">
                      <X size={12} />
                    </button>
                  </span>
                ))}
                <input
                  type="text"
                  value={keywordInput}
                  onChange={(e) => setKeywordInput(e.target.value)}
                  onKeyDown={handleKeyDown}
                  onBlur={() => { if (keywordInput.trim()) addKeyword(keywordInput); }}
                  placeholder={keywords.length === 0 ? '输入关键词后按 Enter 添加' : '继续添加...'}
                  className="flex-1 min-w-[120px] border-none outline-none text-sm bg-transparent py-1"
                />
              </div>
              {errors.keywords && <p className="mt-1 text-xs text-red-600">{errors.keywords}</p>}
              {!errors.keywords && keywords.length > 0 && (
                <p className="mt-1 text-xs text-gray-400">{keywords.length} 个关键词</p>
              )}
            </div>
          </div>

          <div className="border-t border-gray-100 pt-6">
            <h3 className="text-sm font-semibold text-gray-900 mb-1">品牌信息</h3>
            <p className="text-xs text-gray-500 mb-4">选填，完善品牌信息有助于提升 AI 分析质量</p>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label htmlFor="brandIntro" className="block text-sm font-medium text-gray-700 mb-1">
                  品牌介绍
                </label>
                <textarea
                  id="brandIntro"
                  value={brandIntro}
                  onChange={(e) => setBrandIntro(e.target.value)}
                  placeholder="品牌定位、核心理念、市场定位等"
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors resize-none"
                />
                <p className="mt-1 text-xs text-gray-400 text-right">{brandIntro.length} 字</p>
              </div>

              <div>
                <label htmlFor="productIntro" className="block text-sm font-medium text-gray-700 mb-1">
                  产品介绍
                </label>
                <textarea
                  id="productIntro"
                  value={productIntro}
                  onChange={(e) => setProductIntro(e.target.value)}
                  placeholder="核心产品功能、优势特点等"
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors resize-none"
                />
                <p className="mt-1 text-xs text-gray-400 text-right">{productIntro.length} 字</p>
              </div>
            </div>

            <div className="mt-4">
              <label htmlFor="serviceIntro" className="block text-sm font-medium text-gray-700 mb-1">
                服务介绍
              </label>
              <textarea
                id="serviceIntro"
                value={serviceIntro}
                onChange={(e) => setServiceIntro(e.target.value)}
                placeholder="服务模式、服务范围、服务体系等"
                rows={2}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors resize-none"
              />
              <p className="mt-1 text-xs text-gray-400 text-right">{serviceIntro.length} 字</p>
            </div>
          </div>

          <div className="flex items-center gap-3 pt-2">
            <button
              type="submit"
              disabled={submitting}
              className="inline-flex items-center gap-2 px-6 py-2 bg-primary text-white text-sm font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {submitting && <Loader2 size={14} className="animate-spin" />}
              {submitting ? '创建中...' : '创建项目'}
            </button>
            <Link
              href="/projects"
              className="px-6 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-50 transition-colors"
            >
              取消
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}
