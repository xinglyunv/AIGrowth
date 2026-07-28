'use client';

import { useEffect, useState, use } from 'react';
import { useRouter } from 'next/navigation';
import {
  User, Mail, Calendar, Shield, Activity, FolderOpen, Zap, Coins, X,
  ChevronLeft, Loader2, ToggleLeft, ToggleRight, KeyRound, AlertTriangle
  , Save
} from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface Project {
  id: string;
  name: string;
  status: string;
  created_at: string;
}

interface UserDetail {
  id: string;
  username: string;
  email: string;
  role: string;
  status: string;
  created_at: string;
  project_count: number;
  ai_usage: number;
  credits: number;
  projects: Project[];
}

export default function UserDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const [user, setUser] = useState<UserDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [actionLoading, setActionLoading] = useState<string | null>(null);
  const [showCreditsModal, setShowCreditsModal] = useState(false);
  const [creditsValue, setCreditsValue] = useState(0);
  const [showAccountModal, setShowAccountModal] = useState(false);
  const [accountForm, setAccountForm] = useState({ username: '', email: '', password: '', credits: 0 });

  useEffect(() => {
    setLoading(true);
    setError('');
    Promise.all([
      apiRequest<UserDetail>(`/users/${id}`),
      apiRequest<Project[]>(`/users/${id}/projects`),
    ])
      .then(([userRes, projectsRes]) => {
        if (userRes.data) {
          setUser({
            ...userRes.data,
            projects: Array.isArray(projectsRes.data) ? projectsRes.data : [],
            project_count: Array.isArray(projectsRes.data) ? projectsRes.data.length : (userRes.data.project_count || 0),
            ai_usage: userRes.data.ai_usage || 0,
          });
        }
      })
      .catch((e) => setError(e instanceof Error ? e.message : '加载失败'))
      .finally(() => setLoading(false));
  }, [id]);

  const handleToggleStatus = async () => {
    if (!user) return;
    const newStatus = user.status === 'active' ? 'disabled' : 'active';
    if (!confirm(`确定${newStatus === 'disabled' ? '禁用' : '启用'}该用户？`)) return;
    setActionLoading('status');
    try {
      await apiRequest(`/users/${id}`, { method: 'PUT', body: JSON.stringify({ status: newStatus }) });
      setUser((prev) => prev ? { ...prev, status: newStatus } : prev);
    } catch (e) {
      setError(e instanceof Error ? e.message : '操作失败');
    }
    setActionLoading(null);
  };

  const handleResetPassword = () => {
    const newPassword = prompt('请输入新密码：');
    if (!newPassword) return;
    setActionLoading('password');
    apiRequest(`/users/${id}`, { method: 'PUT', body: JSON.stringify({ password: newPassword }) })
      .catch(() => {})
      .finally(() => setActionLoading(null));
  };

  const handleUpdateCredits = async () => {
    if (creditsValue < 0) return;
    setActionLoading('credits');
    try {
      await apiRequest(`/users/${id}`, { method: 'PUT', body: JSON.stringify({ credits: creditsValue }) });
      setUser((prev) => prev ? { ...prev, credits: creditsValue } : prev);
      setShowCreditsModal(false);
    } catch {
      alert('保存失败');
    }
    setActionLoading(null);
  };

  const openAccountModal = () => {
    if (!user) return;
    setAccountForm({ username: user.username, email: user.email, password: '', credits: user.credits || 0 });
    setShowAccountModal(true);
  };

  const handleUpdateAccount = async () => {
    if (!user || !accountForm.username.trim() || !accountForm.email.trim() || accountForm.credits < 0) return;
    setActionLoading('account');
    try {
      const payload: Record<string, string | number> = {
        username: accountForm.username.trim(),
        email: accountForm.email.trim(),
        credits: accountForm.credits,
      };
      if (accountForm.password.trim()) payload.password = accountForm.password.trim();
      await apiRequest(`/users/${id}`, { method: 'PUT', body: JSON.stringify(payload) });
      setUser((prev) => prev ? { ...prev, username: accountForm.username.trim(), email: accountForm.email.trim(), credits: accountForm.credits } : prev);
      setShowAccountModal(false);
    } catch (e) {
      alert(e instanceof Error ? e.message : '保存失败');
    } finally {
      setActionLoading(null);
    }
  };

  if (loading) {
    return <div className="flex items-center justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-gray-400" /></div>;
  }

  if (error) {
    return <div className="p-12 text-center">
      <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-sm text-red-700 inline-block">{error}</div>
    </div>;
  }

  if (!user) {
    return <div className="text-center py-20 text-gray-400">用户不存在</div>;
  }

  const roleLabels: Record<string, string> = { admin: '管理员', superadmin: '超级管理员', user: '用户' };

  const infoCard = 'bg-white rounded-lg border border-gray-200 p-5';
  const labelClass = 'text-xs text-gray-400 mb-0.5';
  const valueClass = 'text-sm font-medium text-gray-900';

  return (
    <div>
      <button onClick={() => router.push('/users')}
        className="flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 mb-4">
        <ChevronLeft className="h-4 w-4" />
        返回用户列表
      </button>

      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <User className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">用户详情</h2>
        </div>
        <div className="flex items-center gap-2">
          <button onClick={handleToggleStatus} disabled={actionLoading !== null}
            className={`flex items-center gap-1.5 px-4 py-2 text-sm font-medium rounded-md transition-colors ${
              user.status === 'active'
                ? 'bg-yellow-50 text-yellow-700 hover:bg-yellow-100 border border-yellow-200'
                : 'bg-green-50 text-green-700 hover:bg-green-100 border border-green-200'
            } disabled:opacity-60`}>
            {actionLoading === 'status' ? <Loader2 className="h-4 w-4 animate-spin" /> :
              user.status === 'active' ? <ToggleLeft className="h-4 w-4" /> : <ToggleRight className="h-4 w-4" />}
            {user.status === 'active' ? '禁用用户' : '启用用户'}
          </button>
          <button onClick={handleResetPassword} disabled={actionLoading !== null}
            className="flex items-center gap-1.5 px-4 py-2 text-sm font-medium rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200 transition-colors disabled:opacity-60">
            {actionLoading === 'password' ? <Loader2 className="h-4 w-4 animate-spin" /> : <KeyRound className="h-4 w-4" />}
            重置密码
          </button>
          <button onClick={() => { setCreditsValue(user.credits || 0); setShowCreditsModal(true); }} disabled={actionLoading !== null}
            className="flex items-center gap-1.5 px-4 py-2 text-sm font-medium rounded-md bg-gray-100 text-gray-700 hover:bg-gray-200 transition-colors disabled:opacity-60">
            <Coins className="h-4 w-4" />
            编辑额度
          </button>
          <button onClick={openAccountModal} disabled={actionLoading !== null}
            className="flex items-center gap-1.5 px-4 py-2 text-sm font-medium rounded-md bg-primary text-white hover:bg-primary/90 transition-colors disabled:opacity-60">
            编辑账户
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <div className={infoCard}>
            <h3 className="text-base font-semibold text-gray-900 mb-4">账户信息</h3>
            <div className="grid grid-cols-2 gap-4">
              <div><p className={labelClass}>用户名</p><p className={valueClass}>{user.username}</p></div>
              <div><p className={labelClass}>邮箱</p><p className="text-sm text-gray-600 flex items-center gap-1"><Mail className="h-3.5 w-3.5 text-gray-400" />{user.email}</p></div>
              <div><p className={labelClass}>角色</p><p className="flex items-center gap-1"><Shield className="h-3.5 w-3.5 text-gray-400" /><span className={`text-sm font-medium px-2 py-0.5 rounded-full ${user.role === 'admin' || user.role === 'superadmin' ? 'bg-purple-50 text-purple-700' : 'bg-gray-100 text-gray-600'}`}>{roleLabels[user.role] || user.role}</span></p></div>
              <div><p className={labelClass}>状态</p><span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium ${user.status === 'active' ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'}`}>{user.status === 'active' ? '正常' : '已禁用'}</span></div>
               <div><p className={labelClass}>注册时间</p><p className={valueClass}><Calendar className="h-3.5 w-3.5 inline text-gray-400 mr-1" />{new Date(user.created_at).toLocaleString('zh-CN')}</p></div>
              <div><p className={labelClass}>用户 ID</p><p className="text-sm font-mono text-gray-500">{user.id}</p></div>
            </div>
          </div>

          <div className={infoCard}>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-base font-semibold text-gray-900">项目列表</h3>
              <span className="text-xs text-gray-400">共 {user.projects.length} 个项目</span>
            </div>
            {user.projects.length === 0 ? (
              <p className="text-sm text-gray-400 text-center py-6">暂无项目</p>
            ) : (
              <div className="divide-y divide-gray-100">
                {user.projects.map((p) => (
                  <div key={p.id} className="flex items-center justify-between py-3">
                    <div className="flex items-center gap-3">
                      <FolderOpen className="h-4 w-4 text-gray-400" />
                      <span className="text-sm font-medium text-gray-900">{p.name}</span>
                    </div>
                    <div className="flex items-center gap-3">
                      <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${p.status === 'active' ? 'bg-green-50 text-green-700' : 'bg-gray-100 text-gray-500'}`}>{p.status === 'active' ? '进行中' : '已停用'}</span>
                      <span className="text-xs text-gray-400">{p.created_at}</span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        <div className="space-y-6">
          <div className={infoCard}>
            <h3 className="text-base font-semibold text-gray-900 mb-4">使用概览</h3>
            <div className="space-y-4">
              <div className="flex items-center justify-between p-3 bg-blue-50 rounded-md">
                <div className="flex items-center gap-2">
                  <FolderOpen className="h-4 w-4 text-blue-600" />
                  <span className="text-sm text-blue-700">项目数</span>
                </div>
                <span className="text-lg font-bold text-blue-700">{user.project_count}</span>
              </div>
              <div className="flex items-center justify-between p-3 bg-purple-50 rounded-md">
                <div className="flex items-center gap-2">
                  <Zap className="h-4 w-4 text-purple-600" />
                  <span className="text-sm text-purple-700">AI 调用次数</span>
                </div>
                <span className="text-lg font-bold text-purple-700">{user.ai_usage}</span>
              </div>
              <div className="flex items-center justify-between p-3 bg-green-50 rounded-md">
                <div className="flex items-center gap-2">
                  <Activity className="h-4 w-4 text-green-600" />
                  <span className="text-sm text-green-700">活跃度</span>
                </div>
                <span className="text-lg font-bold text-green-700">{user.ai_usage > 0 ? '活跃' : '不活跃'}</span>
              </div>
              <div className="flex items-center justify-between p-3 bg-amber-50 rounded-md">
                <div className="flex items-center gap-2">
                  <Coins className="h-4 w-4 text-amber-600" />
                  <span className="text-sm text-amber-700">额度</span>
                </div>
                <span className="text-lg font-bold text-amber-700">{user.credits ?? 0}</span>
              </div>
            </div>
          </div>

          <div className={`${infoCard} border-yellow-200 bg-yellow-50`}>
            <div className="flex items-start gap-2">
              <AlertTriangle className="h-4 w-4 text-yellow-600 mt-0.5 shrink-0" />
              <div>
                <p className="text-sm font-medium text-yellow-800">操作提示</p>
                <p className="text-xs text-yellow-700 mt-1">禁用用户后，该用户将无法登录系统。重置密码会向用户邮箱发送密码重置链接。</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {showCreditsModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40" onClick={() => setShowCreditsModal(false)}>
          <div className="bg-white rounded-xl shadow-xl w-full max-w-sm mx-4" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-center justify-between p-5 border-b border-gray-200">
              <h3 className="text-base font-semibold text-gray-900">编辑用户额度</h3>
              <button onClick={() => setShowCreditsModal(false)} className="p-1 text-gray-400 hover:text-gray-600"><X className="h-4 w-4" /></button>
            </div>
            <div className="p-5 space-y-4">
              <div className="text-sm text-gray-600">
                用户：<span className="font-medium text-gray-900">{user.username}</span>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">额度</label>
                <input type="number" min={0} className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary text-gray-900"
                  value={creditsValue} onChange={(e) => setCreditsValue(parseInt(e.target.value) || 0)} />
              </div>
            </div>
            <div className="flex items-center justify-end gap-3 p-5 border-t border-gray-200">
              <button onClick={() => setShowCreditsModal(false)} className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900">取消</button>
              <button onClick={handleUpdateCredits} disabled={actionLoading === 'credits'}
                className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 disabled:opacity-60">
                {actionLoading === 'credits' ? <Loader2 className="h-4 w-4 animate-spin" /> : <Coins className="h-4 w-4" />}
                保存
              </button>
            </div>
          </div>
        </div>
      )}

      {showAccountModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40" onClick={() => setShowAccountModal(false)}>
          <div className="bg-white rounded-xl shadow-xl w-full max-w-md mx-4" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-center justify-between p-5 border-b border-gray-200">
              <h3 className="text-base font-semibold text-gray-900">编辑账户信息</h3>
              <button onClick={() => setShowAccountModal(false)} className="p-1 text-gray-400 hover:text-gray-600"><X className="h-4 w-4" /></button>
            </div>
            <div className="p-5 space-y-4">
              <label className="block text-sm font-medium text-gray-700">用户名
                <input value={accountForm.username} onChange={(e) => setAccountForm((p) => ({ ...p, username: e.target.value }))}
                  className="mt-1.5 w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary" />
              </label>
              <label className="block text-sm font-medium text-gray-700">邮箱
                <input type="email" value={accountForm.email} onChange={(e) => setAccountForm((p) => ({ ...p, email: e.target.value }))}
                  className="mt-1.5 w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary" />
              </label>
              <label className="block text-sm font-medium text-gray-700">新密码（留空保持不变）
                <input type="password" value={accountForm.password} onChange={(e) => setAccountForm((p) => ({ ...p, password: e.target.value }))}
                  className="mt-1.5 w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary" />
              </label>
              <label className="block text-sm font-medium text-gray-700">额度
                <input type="number" min={0} value={accountForm.credits} onChange={(e) => setAccountForm((p) => ({ ...p, credits: parseInt(e.target.value) || 0 }))}
                  className="mt-1.5 w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary" />
              </label>
            </div>
            <div className="flex items-center justify-end gap-3 p-5 border-t border-gray-200">
              <button onClick={() => setShowAccountModal(false)} className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900">取消</button>
              <button onClick={handleUpdateAccount} disabled={actionLoading === 'account' || !accountForm.username.trim() || !accountForm.email.trim()}
                className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 disabled:opacity-60">
                {actionLoading === 'account' ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
                保存
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
