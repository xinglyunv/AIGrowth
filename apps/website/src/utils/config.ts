export interface SiteConfig {
  site_name: string;
  site_title: string;
  site_description: string;
  logo_url: string;
  footer_text: string;
}

let cachedConfig: SiteConfig | null = null;

export async function getSiteConfig(): Promise<SiteConfig> {
  if (cachedConfig) return cachedConfig;
  try {
    const res = await fetch('/api/v1/settings');
    const data = await res.json();
    if (data.success) {
      cachedConfig = data.data;
      return data.data;
    }
  } catch {}
  return {
    site_name: 'AI Growth Engine',
    site_title: 'AI 品牌可见度分析平台',
    site_description: 'AI 品牌可见度分析与增长优化 SaaS 平台',
    logo_url: '',
    footer_text: 'AI Growth Engine. All rights reserved.',
  };
}
