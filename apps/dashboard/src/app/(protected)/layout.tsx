'use client';
import { useAuth } from '@/lib/auth';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import { ReactNode, useState, useEffect, useCallback } from 'react';
import { apiRequest } from '@/lib/api';
import { LayoutDashboard, FolderOpen, BarChart3, FileText, Settings, LogOut, Search, Menu, X, Bell, CreditCard, UsersRound } from 'lucide-react';

const navItems = [
  { href: '/dashboard', label: '首页', icon: LayoutDashboard, disabled: false },
  { href: '/projects', label: '品牌项目', icon: FolderOpen, disabled: false },
  { href: '/tasks', label: 'AI 检测', icon: Search, disabled: false },
  { href: '/analytics', label: '数据分析', icon: BarChart3, disabled: false },
  { href: '/reports', label: '报告中心', icon: FileText, disabled: false },
  { href: '/billing', label: '充值中心', icon: CreditCard, disabled: false },
  { href: '/spaces', label: '团队空间', icon: UsersRound, disabled: false },
  { href: '/profile', label: '设置', icon: Settings, disabled: false },
];

export default function ProtectedLayout({ children }: { children: ReactNode }) {
  const { user, loading, logout, hasToken } = useAuth();
  const router = useRouter();
  const pathname = usePathname();
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [unreadCount, setUnreadCount] = useState(0);

  const fetchUnreadCount = useCallback(() => {
    apiRequest<{ count: number }>('/notifications/unread-count')
      .then((res) => setUnreadCount(res.data.count ?? 0))
      .catch(() => {});
  }, []);

  useEffect(() => {
    if (!loading && user) {
      fetchUnreadCount();
      const interval = setInterval(fetchUnreadCount, 30000);
      return () => clearInterval(interval);
    }
  }, [loading, user, fetchUnreadCount]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <p className="text-gray-500">加载中...</p>
      </div>
    );
  }

  if (!user && !hasToken) {
    router.replace('/login');
    return null;
  }

  // Still loading user data but has a token (retry in progress)
  if (!user && hasToken) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <p className="text-gray-500">加载中...</p>
      </div>
    );
  }

  const handleLogout = () => {
    logout();
    router.replace('/login');
  };

  return (
    <div className="min-h-screen flex">
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-20 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      <aside
        className={`
          fixed lg:static inset-y-0 left-0 z-30 w-60 bg-white border-r border-gray-200
          transform transition-transform duration-200 ease-in-out
          ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
          lg:translate-x-0 lg:flex flex-col shrink-0
        `}
      >
        <div className="p-5 border-b border-gray-100">
          <span className="text-lg font-bold text-primary">AI Growth Engine</span>
        </div>
        <nav className="flex-1 p-3 space-y-1">
          {navItems.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(item.href + '/');
            if (item.disabled) {
              return (
                <span
                  key={item.href}
                  className="flex items-center gap-3 px-3 py-2 rounded-btn text-sm text-gray-400 cursor-not-allowed"
                >
                  <item.icon size={18} />
                  {item.label}
                </span>
              );
            }
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={() => setSidebarOpen(false)}
                className={`flex items-center gap-3 px-3 py-2 rounded-btn text-sm font-medium transition-colors ${
                  isActive
                    ? 'bg-primary/10 text-primary'
                    : 'text-gray-700 hover:bg-gray-100'
                }`}
              >
                <item.icon size={18} />
                {item.label}
              </Link>
            );
          })}
        </nav>
        <div className="p-3 border-t border-gray-100 space-y-1">
          <Link
            href="/notifications"
            onClick={() => setSidebarOpen(false)}
            className={`flex items-center gap-3 px-3 py-2 rounded-btn text-sm font-medium transition-colors ${
              pathname === '/notifications'
                ? 'bg-primary/10 text-primary'
                : 'text-gray-700 hover:bg-gray-100'
            }`}
          >
            <Bell size={18} />
            <span>通知</span>
            {unreadCount > 0 && (
              <span className="ml-auto inline-flex items-center justify-center min-w-[20px] h-5 px-1.5 rounded-full bg-red-500 text-white text-xs font-medium">
                {unreadCount > 99 ? '99+' : unreadCount}
              </span>
            )}
          </Link>
          <button
            onClick={handleLogout}
            className="flex items-center gap-3 px-3 py-2 rounded-btn text-sm font-medium text-gray-600 hover:bg-gray-100 w-full transition-colors"
          >
            <LogOut size={18} />
            退出登录
          </button>
        </div>
      </aside>

      <div className="flex-1 flex flex-col min-w-0">
        <header className="lg:hidden flex items-center justify-between p-4 border-b border-gray-200 bg-white">
          <button onClick={() => setSidebarOpen(true)} className="p-1 rounded-btn hover:bg-gray-100">
            <Menu size={22} className="text-gray-700" />
          </button>
          <div className="flex items-center gap-3">
            <Link href="/notifications" className="relative p-1 rounded-btn hover:bg-gray-100">
              <Bell size={20} className="text-gray-600" />
              {unreadCount > 0 && (
                <span className="absolute -top-0.5 -right-0.5 inline-flex items-center justify-center min-w-[16px] h-4 px-1 rounded-full bg-red-500 text-white text-[10px] font-medium">
                  {unreadCount > 99 ? '99+' : unreadCount}
                </span>
              )}
            </Link>
            <span className="text-lg font-bold text-primary">AI Growth Engine</span>
          </div>
          <div className="w-8" />
        </header>
        <main className="flex-1 p-4 lg:p-8 overflow-auto">{children}</main>
      </div>
    </div>
  );
}
