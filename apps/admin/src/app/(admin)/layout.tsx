'use client';

import { useEffect, useState } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import Link from 'next/link';
import { Users, ShoppingCart, Cpu, Settings, FileText, LogOut, Shield, Menu, X, LayoutDashboard, MessageSquare, ClipboardList, FileText as FileTextIcon, CreditCard, BarChart3, Key } from 'lucide-react';
import { useAuth } from '@/lib/auth';

const sidebarMenu = [
  { href: '/', label: '仪表盘', icon: LayoutDashboard },
  { href: '/users', label: '用户管理', icon: Users },
  { href: '/orders', label: '订单管理', icon: ShoppingCart },
  { href: '/cdk', label: 'CDK 管理', icon: Key },
  { href: '/plans', label: '套餐管理', icon: CreditCard },
  { href: '/models', label: 'AI 模型管理', icon: Cpu },
  { href: '/models/stats', label: '调用统计', icon: BarChart3 },
  { href: '/tasks', label: '任务管理', icon: ClipboardList },
  { href: '/prompts', label: '提示词管理', icon: FileTextIcon },
  { href: '/payment', label: '支付配置', icon: CreditCard },
  { href: '/settings', label: '站点设置', icon: Settings },
  { href: '/messages', label: '留言管理', icon: MessageSquare },
  { href: '/logs', label: '日志管理', icon: FileText },
];

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const { admin, loading, logout, hasToken } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  useEffect(() => {
    if (!loading && !admin && !hasToken) router.push('/login');
  }, [loading, admin, hasToken, router]);

  if (loading) return <div className="min-h-screen flex items-center justify-center bg-gray-50"><p className="text-gray-500">加载中...</p></div>;
  if (!admin && !hasToken) return null;
  if (!admin && hasToken) return <div className="min-h-screen flex items-center justify-center bg-gray-50"><p className="text-gray-500">加载中...</p></div>;

  const handleLogout = () => { logout(); router.push('/login'); };

  const isActive = (href: string) => {
    const p = pathname ?? '';
    if (href === '/') return p === '/' || p === '';
    return p.startsWith(href);
  };

  return (
    <div className="min-h-screen flex bg-gray-50">
      {sidebarOpen && (
        <div className="fixed inset-0 bg-black/50 z-20 lg:hidden" onClick={() => setSidebarOpen(false)} />
      )}

      <aside className={`
        fixed lg:static inset-y-0 left-0 z-30 w-60 bg-gray-900 text-white flex flex-col shrink-0
        transform transition-transform duration-200 ease-in-out
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        lg:translate-x-0 lg:flex
      `}>
        <div className="px-5 py-5 border-b border-gray-700 flex items-center justify-between">
          <div>
            <h1 className="text-lg font-bold tracking-tight">AI Growth Engine</h1>
            <p className="text-xs text-gray-400 mt-1">管理后台</p>
          </div>
          <button onClick={() => setSidebarOpen(false)} className="lg:hidden p-1 text-gray-400 hover:text-white">
            <X size={18} />
          </button>
        </div>
        <nav className="flex-1 px-3 py-4 space-y-1">
          {sidebarMenu.map((item) => {
            const active = isActive(item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={() => setSidebarOpen(false)}
                className={`flex items-center gap-3 px-3 py-2 text-sm rounded-md transition-colors ${
                  active ? 'bg-gray-700 text-white' : 'text-gray-300 hover:text-white hover:bg-gray-700'
                }`}
              >
                <item.icon className="h-4 w-4" />
                {item.label}
              </Link>
            );
          })}
        </nav>
        <div className="px-3 py-4 border-t border-gray-700">
          <button onClick={handleLogout}
            className="w-full flex items-center gap-3 px-3 py-2 text-sm text-gray-300 hover:text-white hover:bg-gray-700 rounded-md transition-colors">
            <LogOut className="h-4 w-4" />
            退出登录
          </button>
        </div>
      </aside>

      <div className="flex-1 flex flex-col min-w-0">
        <header className="lg:hidden bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between">
          <span className="text-sm font-bold text-gray-900">AI Growth Engine</span>
          <button onClick={() => setSidebarOpen(true)} className="p-1">
            <Menu size={20} />
          </button>
        </header>

        <header className="hidden lg:flex bg-white border-b border-gray-200 px-6 py-3 items-center justify-end shrink-0">
          <div className="flex items-center gap-3 text-sm text-gray-700">
            <div className="flex items-center gap-1.5">
              <Shield className="h-4 w-4 text-primary" />
              <span className="font-medium">{admin?.username}</span>
            </div>
            <span className="text-gray-300">|</span>
            <span className="px-2 py-0.5 bg-blue-50 text-primary text-xs rounded-full font-medium">
              {admin?.role === 'superadmin' ? '超级管理员' : admin?.role === 'admin' ? '管理员' : admin?.role}
            </span>
          </div>
        </header>

        <main className="flex-1 p-4 lg:p-6 overflow-auto">{children}</main>
      </div>
    </div>
  );
}
