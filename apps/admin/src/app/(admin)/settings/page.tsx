'use client';

import { useEffect, useState } from 'react';
import { Settings as SettingsIcon, Save, Loader2 } from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface SiteConfig {
  site_name: string;
  site_title: string;
  site_description: string;
  site_theme: string;
  logo_url: string;
  footer_text: string;
  contact_email: string;
  contact_address: string;
  contact_phone: string;
  working_hours: string;
  stat_companies: string;
  stat_models: string;
  stat_reports: string;
  stat_accuracy: string;
  hero_tagline: string;
  hero_title: string;
  hero_subtitle: string;
  allow_registration: string;
  smtp_host: string;
  smtp_port: string;
  smtp_user: string;
  smtp_password: string;
  smtp_from: string;
  sms_provider: string;
  sms_access_key: string;
  sms_secret_key: string;
  sms_sign_name: string;
  smsbao_username: string;
  smsbao_password: string;
  navigation: { items: { label: string; href: string }[]; login_label: string; login_href: string; cta_label: string; cta_href: string };
  hero: { tagline: string; title: string; subtitle: string; primary_label: string; primary_href: string; secondary_label: string; secondary_href: string };
  features: { eyebrow: string; title: string; description: string; href: string; highlights: string[]; icon: string }[];
  stats: { items: { value: string; label: string }[] };
  dashboard: { score: string; score_delta: string; mentions: string; recommendation_rate: string; competitor_gap: string; trend_label: string };
  seo: { title: string; description: string; og_image: string };
  theme: { primary: string; accent: string; surface: string; text: string; muted: string };
}

const DEFAULT_CONFIG: SiteConfig = {
  site_name: '', site_title: '', site_description: '', site_theme: 'syntro', logo_url: '', footer_text: '',
  contact_email: '', contact_address: '', contact_phone: '', working_hours: '',
  stat_companies: '', stat_models: '', stat_reports: '', stat_accuracy: '',
  hero_tagline: '', hero_title: '', hero_subtitle: '',
  allow_registration: 'true',
  smtp_host: '', smtp_port: '', smtp_user: '', smtp_password: '', smtp_from: '',
  sms_provider: '', sms_access_key: '', sms_secret_key: '', sms_sign_name: '', smsbao_username: '', smsbao_password: '',
  navigation: { items: [], login_label: '登录', login_href: '/login', cta_label: '开始分析', cta_href: '/register' },
  hero: { tagline: '', title: '', subtitle: '', primary_label: '', primary_href: '/register', secondary_label: '', secondary_href: '#features' },
  features: [],
  stats: { items: [] },
  dashboard: { score: '82.4', score_delta: '+18.6%', mentions: '1,284', recommendation_rate: '64.8%', competitor_gap: '-12.4', trend_label: '过去 30 天' },
  seo: { title: '', description: '', og_image: '' },
  theme: { primary: '#6d5efc', accent: '#2dd4bf', surface: '#f8fafc', text: '#0f172a', muted: '#64748b' },
};

function mergeConfig(data: Partial<SiteConfig>): SiteConfig {
  return {
    ...DEFAULT_CONFIG,
    ...data,
    navigation: { ...DEFAULT_CONFIG.navigation, ...data.navigation },
    hero: { ...DEFAULT_CONFIG.hero, ...data.hero },
    features: data.features?.length ? data.features : DEFAULT_CONFIG.features,
    stats: { ...DEFAULT_CONFIG.stats, ...data.stats },
    dashboard: { ...DEFAULT_CONFIG.dashboard, ...data.dashboard },
    seo: { ...DEFAULT_CONFIG.seo, ...data.seo },
    theme: { ...DEFAULT_CONFIG.theme, ...data.theme },
  };
}

export default function SettingsPage() {
  const [form, setForm] = useState<SiteConfig>(DEFAULT_CONFIG);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState('');

  const loadSettings = () => {
    setLoading(true);
    setError('');
    apiRequest<SiteConfig>('/settings')
      .then((r) => { if (r.data) setForm(mergeConfig(r.data)); })
      .catch((err) => { setError(err?.message || '加载设置失败，请重试'); })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    loadSettings();
  }, []);

  const update = <K extends keyof SiteConfig>(key: K, value: SiteConfig[K]) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  };

  const validateConfig = () => {
    if (!form.site_name.trim()) return '请填写站点名称';
    if (!form.hero.title.trim() || !form.hero.subtitle.trim()) return '请填写完整的 Hero 标题和描述';
    if (!form.hero.primary_href.trim() || !form.hero.secondary_href.trim()) return '请填写 Hero 按钮链接';
    if (form.features.length < 3 || form.features.some((feature) => !feature.title.trim() || !feature.description.trim())) return '请完整填写至少三张功能卡片';
    const colorPattern = /^#[0-9a-f]{6}$/i;
    if (Object.values(form.theme).some((color) => !colorPattern.test(color))) return '主题颜色必须使用六位十六进制格式，例如 #6d5efc';
    return '';
  };

  const handleSave = () => {
    const validationError = validateConfig();
    if (validationError) {
      setError(validationError);
      setSaved(false);
      return;
    }
    setSaving(true);
    setError('');
    setSaved(false);
    apiRequest('/settings', { method: 'PUT', body: JSON.stringify(form) })
      .then(() => { setSaved(true); setTimeout(() => setSaved(false), 3000); })
      .catch((err) => { setError(err?.message || '保存失败，请重试'); })
      .finally(() => setSaving(false));
  };

  const inputClass = 'w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary text-gray-900';
  const labelClass = 'block text-sm font-medium text-gray-700 mb-1.5';
  const sectionCard = 'bg-white rounded-lg border border-gray-200 p-6 mb-6';

  if (loading) {
    return <div className="flex items-center justify-center py-20"><Loader2 className="h-5 w-5 animate-spin text-gray-400" /><p className="ml-2 text-gray-400">加载中...</p></div>;
  }

  return (
    <div>
      <div className="flex items-center gap-3 mb-6">
        <SettingsIcon className="h-6 w-6 text-gray-700" />
        <h2 className="text-xl font-bold text-gray-900">站点设置</h2>
      </div>

      {saved && <div className="mb-4 px-4 py-3 bg-green-50 border border-green-200 rounded-md text-sm text-green-700">设置已保存成功</div>}
      {error && <div className="mb-4 flex items-center justify-between gap-4 px-4 py-3 bg-red-50 border border-red-200 rounded-md text-sm text-red-700"><span>{error}</span><button type="button" onClick={loadSettings} className="shrink-0 rounded-md border border-red-300 px-3 py-1.5 font-medium hover:bg-red-100">重新加载</button></div>}

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">基本设置</h3>
        <div className="space-y-4">
          <div>
            <label className={labelClass}>官网主题</label>
            <select className={inputClass} value={form.site_theme} onChange={(e) => update('site_theme', e.target.value)}>
              <option value="syntro">Syntro Astro（默认）</option>
              <option value="saas-landing">SaaS Landing Template</option>
              <option value="saas-ui">SaaS UI Next.js</option>
              <option value="hikari">Hikari SaaS</option>
              <option value="cruip">Cruip Open React</option>
              <option value="nextly">Nextly Landing Page</option>
            </select>
            <p className="text-xs text-gray-500 mt-1.5">每个选项对应独立的首页布局、组件节奏和视觉语言。</p>
          </div>
          <div>
            <label className={labelClass}>站点名称</label>
            <input type="text" className={inputClass} value={form.site_name} onChange={(e) => update('site_name', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>网站标题</label>
            <input type="text" className={inputClass} value={form.site_title} onChange={(e) => update('site_title', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>网站描述</label>
            <textarea className={`${inputClass} resize-none`} rows={2} value={form.site_description} onChange={(e) => update('site_description', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>页脚文本</label>
            <input type="text" className={inputClass} value={form.footer_text} onChange={(e) => update('footer_text', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>Logo URL</label>
            <input type="text" className={inputClass} placeholder="输入 Logo 图片 URL" value={form.logo_url} onChange={(e) => update('logo_url', e.target.value)} />
            {form.logo_url && (
              <div className="mt-2 inline-block border border-gray-200 rounded-md p-2">
                <img src={form.logo_url} alt="Logo 预览" className="max-h-12 object-contain" onError={(e) => { (e.target as HTMLImageElement).style.display = 'none'; }} />
              </div>
            )}
          </div>
        </div>
      </div>

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">首页结构化内容</h3>
        <p className="text-xs text-gray-500 mb-4">六套官网主题共用这组内容配置，保存后会同步影响导航、Hero、功能卡片和统计模块。</p>
        <div className="space-y-4">
          <div>
            <label className={labelClass}>Hero 标签</label>
            <input className={inputClass} value={form.hero.tagline} onChange={(e) => setForm((prev) => ({ ...prev, hero: { ...prev.hero, tagline: e.target.value } }))} />
          </div>
          <div>
            <label className={labelClass}>Hero 标题</label>
            <input className={inputClass} value={form.hero.title} onChange={(e) => setForm((prev) => ({ ...prev, hero: { ...prev.hero, title: e.target.value } }))} />
          </div>
          <div>
            <label className={labelClass}>Hero 描述</label>
            <textarea className={`${inputClass} resize-none`} rows={3} value={form.hero.subtitle} onChange={(e) => setForm((prev) => ({ ...prev, hero: { ...prev.hero, subtitle: e.target.value } }))} />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div><label className={labelClass}>主按钮文字</label><input className={inputClass} value={form.hero.primary_label} onChange={(e) => setForm((prev) => ({ ...prev, hero: { ...prev.hero, primary_label: e.target.value } }))} /></div>
            <div><label className={labelClass}>主按钮链接</label><input className={inputClass} value={form.hero.primary_href} onChange={(e) => setForm((prev) => ({ ...prev, hero: { ...prev.hero, primary_href: e.target.value } }))} /></div>
            <div><label className={labelClass}>次按钮文字</label><input className={inputClass} value={form.hero.secondary_label} onChange={(e) => setForm((prev) => ({ ...prev, hero: { ...prev.hero, secondary_label: e.target.value } }))} /></div>
            <div><label className={labelClass}>次按钮链接</label><input className={inputClass} value={form.hero.secondary_href} onChange={(e) => setForm((prev) => ({ ...prev, hero: { ...prev.hero, secondary_href: e.target.value } }))} /></div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            {form.features.slice(0, 3).map((feature, index) => <div key={index} className="col-span-2 rounded-md border border-gray-100 p-4">
              <p className="mb-3 text-xs font-semibold text-gray-500">功能卡片 {index + 1}</p>
              <div className="grid grid-cols-2 gap-4">
                <input className={inputClass} placeholder="标题" value={feature.title} onChange={(e) => setForm((prev) => ({ ...prev, features: prev.features.map((item, itemIndex) => itemIndex === index ? { ...item, title: e.target.value } : item) }))} />
                <input className={inputClass} placeholder="跳转链接" value={feature.href} onChange={(e) => setForm((prev) => ({ ...prev, features: prev.features.map((item, itemIndex) => itemIndex === index ? { ...item, href: e.target.value } : item) }))} />
                <textarea className={`${inputClass} col-span-2 resize-none`} rows={2} placeholder="描述" value={feature.description} onChange={(e) => setForm((prev) => ({ ...prev, features: prev.features.map((item, itemIndex) => itemIndex === index ? { ...item, description: e.target.value } : item) }))} />
              </div>
            </div>)}
          </div>
        </div>
      </div>

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">主题视觉变量</h3>
        <div className="grid grid-cols-2 gap-4">
          {(['primary', 'accent', 'surface', 'text', 'muted'] as const).map((key) => <div key={key}>
            <label className={labelClass}>{key}</label>
            <input className={inputClass} value={form.theme[key]} onChange={(e) => setForm((prev) => ({ ...prev, theme: { ...prev.theme, [key]: e.target.value } }))} />
          </div>)}
        </div>
      </div>

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">导航与行动</h3>
        <div className="space-y-4">
          {form.navigation.items.map((item, index) => <div key={index} className="grid grid-cols-2 gap-4">
            <input className={inputClass} placeholder={`导航文字 ${index + 1}`} value={item.label} onChange={(e) => setForm((prev) => ({ ...prev, navigation: { ...prev.navigation, items: prev.navigation.items.map((navItem, navIndex) => navIndex === index ? { ...navItem, label: e.target.value } : navItem) } }))} />
            <input className={inputClass} placeholder="导航链接" value={item.href} onChange={(e) => setForm((prev) => ({ ...prev, navigation: { ...prev.navigation, items: prev.navigation.items.map((navItem, navIndex) => navIndex === index ? { ...navItem, href: e.target.value } : navItem) } }))} />
          </div>)}
          <div className="grid grid-cols-2 gap-4 border-t border-gray-100 pt-4">
            <div><label className={labelClass}>登录按钮文字</label><input className={inputClass} value={form.navigation.login_label} onChange={(e) => setForm((prev) => ({ ...prev, navigation: { ...prev.navigation, login_label: e.target.value } }))} /></div>
            <div><label className={labelClass}>登录链接</label><input className={inputClass} value={form.navigation.login_href} onChange={(e) => setForm((prev) => ({ ...prev, navigation: { ...prev.navigation, login_href: e.target.value } }))} /></div>
            <div><label className={labelClass}>主行动文字</label><input className={inputClass} value={form.navigation.cta_label} onChange={(e) => setForm((prev) => ({ ...prev, navigation: { ...prev.navigation, cta_label: e.target.value } }))} /></div>
            <div><label className={labelClass}>主行动链接</label><input className={inputClass} value={form.navigation.cta_href} onChange={(e) => setForm((prev) => ({ ...prev, navigation: { ...prev.navigation, cta_href: e.target.value } }))} /></div>
          </div>
        </div>
      </div>

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">统计数字</h3>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className={labelClass}>服务企业数</label>
            <input type="text" className={inputClass} value={form.stat_companies} onChange={(e) => update('stat_companies', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>AI 模型数</label>
            <input type="text" className={inputClass} value={form.stat_models} onChange={(e) => update('stat_models', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>分析报告数</label>
            <input type="text" className={inputClass} value={form.stat_reports} onChange={(e) => update('stat_reports', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>数据准确率</label>
            <input type="text" className={inputClass} value={form.stat_accuracy} onChange={(e) => update('stat_accuracy', e.target.value)} />
          </div>
        </div>
      </div>

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">联系方式</h3>
        <div className="space-y-4">
          <div>
            <label className={labelClass}>联系邮箱</label>
            <input type="email" className={inputClass} value={form.contact_email} onChange={(e) => update('contact_email', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>联系电话</label>
            <input type="text" className={inputClass} value={form.contact_phone} onChange={(e) => update('contact_phone', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>公司地址</label>
            <input type="text" className={inputClass} value={form.contact_address} onChange={(e) => update('contact_address', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>工作时间</label>
            <input type="text" className={inputClass} value={form.working_hours} onChange={(e) => update('working_hours', e.target.value)} />
          </div>
        </div>
      </div>

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">注册与安全</h3>
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <label className={labelClass}>允许新用户注册</label>
              <p className="text-xs text-gray-500">关闭后用户将无法自行注册新账户</p>
            </div>
            <label className="relative inline-flex items-center cursor-pointer">
              <input
                type="checkbox"
                className="sr-only peer"
                checked={form.allow_registration === 'true'}
                onChange={(e) => update('allow_registration', e.target.checked ? 'true' : 'false')}
              />
              <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
            </label>
          </div>
        </div>
      </div>

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">SMTP 邮件配置</h3>
        <p className="text-xs text-gray-500 mb-4">配置 SMTP 服务器用于发送验证码和系统通知邮件</p>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className={labelClass}>SMTP 主机</label>
            <input type="text" className={inputClass} placeholder="smtp.example.com" value={form.smtp_host} onChange={(e) => update('smtp_host', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>SMTP 端口</label>
            <input type="text" className={inputClass} placeholder="587" value={form.smtp_port} onChange={(e) => update('smtp_port', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>SMTP 用户名</label>
            <input type="text" className={inputClass} placeholder="noreply@example.com" value={form.smtp_user} onChange={(e) => update('smtp_user', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>SMTP 密码</label>
            <input type="password" className={inputClass} placeholder="********" value={form.smtp_password} onChange={(e) => update('smtp_password', e.target.value)} />
          </div>
          <div className="col-span-2">
            <label className={labelClass}>发件人地址</label>
            <input type="email" className={inputClass} placeholder="noreply@example.com" value={form.smtp_from} onChange={(e) => update('smtp_from', e.target.value)} />
          </div>
        </div>
      </div>

      <div className={sectionCard}>
        <h3 className="text-base font-semibold text-gray-900 mb-4">短信配置</h3>
        <p className="text-xs text-gray-500 mb-4">配置短信服务商用于发送验证码和通知短信</p>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className={labelClass}>短信服务商</label>
            <select className={inputClass} value={form.sms_provider} onChange={(e) => update('sms_provider', e.target.value)}>
              <option value="">请选择</option>
              <option value="aliyun">阿里云短信</option>
              <option value="smsbao">短信宝</option>
              <option value="twilio">Twilio</option>
              <option value="tencent">腾讯云短信</option>
            </select>
          </div>
          <div>
            <label className={labelClass}>短信签名</label>
            <input type="text" className={inputClass} placeholder="AI Growth Engine" value={form.sms_sign_name} onChange={(e) => update('sms_sign_name', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>Access Key</label>
            <input type="text" className={inputClass} placeholder="输入 Access Key" value={form.sms_access_key} onChange={(e) => update('sms_access_key', e.target.value)} />
          </div>
          <div>
            <label className={labelClass}>Secret Key</label>
            <input type="password" className={inputClass} placeholder="********" value={form.sms_secret_key} onChange={(e) => update('sms_secret_key', e.target.value)} />
          </div>
        </div>
        {form.sms_provider === 'smsbao' && <div className="mt-4 border-t border-gray-100 pt-4">
          <p className="text-sm font-medium text-gray-700 mb-3">短信宝配置</p>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className={labelClass}>短信宝用户名</label>
              <input type="text" className={inputClass} placeholder="输入短信宝用户名" value={form.smsbao_username} onChange={(e) => update('smsbao_username', e.target.value)} />
            </div>
            <div>
              <label className={labelClass}>短信宝密码</label>
              <input type="password" className={inputClass} placeholder="********" value={form.smsbao_password} onChange={(e) => update('smsbao_password', e.target.value)} />
            </div>
          </div>
        </div>}
      </div>

      <button onClick={handleSave} disabled={saving}
        className="flex items-center gap-1.5 px-5 py-2.5 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 transition-colors disabled:opacity-60">
        {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
         {saving ? '保存中...' : '保存设置'}
      </button>
    </div>
  );
}
