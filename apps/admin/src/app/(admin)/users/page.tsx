'use client';

import { useEffect, useState, useCallback } from 'react';
import {
  Users as UsersIcon, Search, Loader2, Shield, ShieldOff,
  ToggleLeft, ToggleRight, X, Check
} from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface User {
  id: string;
  username: string;
  email: string;
  phone: string;
  company_name: string;
  role: string;
  status: string;
  email_verified: boolean;
  last_login_at: string | null;
  created_at: string;
  updated_at: string;
  credits: number;
}

const roleLabels: Record<string, string> = {
  superadmin: '超级管理员',
  admin: '管理员',
  user: '普通用户',
};
const statusLabels: Record<string, string> = { active: '正常', disabled: '已禁用' };
const ROLES = ['user', 'admin', 'superadmin'];

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');

  const [editModal, setEditModal] = useState<User | null>(null);
  const [editRole, setEditRole] = useState('');
  const [editUsername, setEditUsername] = useState('');
  const [editEmail, setEditEmail] = useState('');
  const [editCredits, setEditCredits] = useState(0);
  const [saving, setSaving] = useState(false);

  const load = useCallback(() => {
    setLoading(true);
    setError('');
    apiRequest<User[]>('/users')
      .then((r) => { if (Array.isArray(r.data)) setUsers(r.data); })
      .catch((e) => setError(e instanceof Error ? e.message : '加载失败'))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => { load(); }, [load]);

  const toggleStatus = async (u: User) => {
    const newStatus = u.status === 'active' ? 'disabled' : 'active';
    try {
      await apiRequest(`/users/${u.id}`, {
        method: 'PUT',
        body: JSON.stringify({ status: newStatus }),
      });
      load();
    } catch { alert('操作失败'); }
  };

  const openEdit = (u: User) => {
    setEditModal(u);
    setEditRole(u.role);
    setEditUsername(u.username);
    setEditEmail(u.email);
    setEditCredits(u.credits || 0);
  };

  const saveUser = async () => {
    if (!editModal || !editUsername.trim() || !editEmail.trim() || editCredits < 0) return;
    const payload = {
      username: editUsername.trim(),
      email: editEmail.trim(),
      role: editRole,
      credits: editCredits,
    };
    const unchanged = payload.username === editModal.username
      && payload.email === editModal.email
      && payload.role === editModal.role
      && payload.credits === (editModal.credits || 0);
    if (unchanged) { setEditModal(null); return; }
    setSaving(true);
    try {
      await apiRequest(`/users/${editModal.id}`, {
        method: 'PUT',
        body: JSON.stringify(payload),
      });
      setEditModal(null);
      load();
    } catch (e) { setError(e instanceof Error ? e.message : '保存失败'); }
    finally { setSaving(false); }
  };

  const displayUsers = search
    ? users.filter(
        (u) =>
          u.username.includes(search) ||
          u.email.toLowerCase().includes(search.toLowerCase())
      )
    : users;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <UsersIcon className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">用户管理</h2>
        </div>
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
          <input
            type="text"
            placeholder="搜索用户名或邮箱..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9 pr-4 py-2 border border-gray-300 rounded-md text-sm
              focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary
              text-gray-900 placeholder-gray-400 w-64"
          />
        </div>
      </div>

      {loading && <div className="p-12 text-center text-gray-400"><Loader2 className="animate-spin inline h-5 w-5" /> 加载中...</div>}

      {error && <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-sm text-red-700">{error}</div>}

      {!loading && !error && (
        <>
          <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-200 bg-gray-50 text-left">
                    <th className="px-4 py-3 font-semibold text-gray-600">用户名</th>
                    <th className="px-4 py-3 font-semibold text-gray-600">邮箱</th>
                    <th className="px-4 py-3 font-semibold text-gray-600">所属公司</th>
                    <th className="px-4 py-3 font-semibold text-gray-600">状态</th>
                    <th className="px-4 py-3 font-semibold text-gray-600">注册时间</th>
                    <th className="px-4 py-3 font-semibold text-gray-600 w-20">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {displayUsers.length === 0 ? (
                    <tr>
                      <td colSpan={6} className="px-4 py-12 text-center text-gray-400">暂无数据</td>
                    </tr>
                  ) : (
                    displayUsers.map((u) => (
                      <tr key={u.id} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                        <td className="px-4 py-3">
                          <div className="font-medium text-gray-900">{u.username}</div>
                          {u.company_name && <div className="text-xs text-gray-400">{u.company_name}</div>}
                        </td>
                        <td className="px-4 py-3 text-gray-600">{u.email}</td>
                        <td className="px-4 py-3 text-gray-600">{u.company_name || '-'}</td>
                        <td className="px-4 py-3">
                          <span
                            className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium
                              ${u.status === 'active'
                                ? 'bg-green-50 text-green-700'
                                : 'bg-red-50 text-red-700'
                              }`}>
                            {statusLabels[u.status] || u.status}
                          </span>
                        </td>
                        <td className="px-4 py-3 text-gray-500 text-xs">
                          {new Date(u.created_at).toLocaleDateString('zh-CN')}
                        </td>
                        <td className="px-4 py-3">
                          <div className="flex items-center gap-2">
                            <button onClick={() => toggleStatus(u)}
                              className="inline-flex items-center gap-1 text-xs font-medium text-gray-500 hover:text-primary transition-colors"
                              title={u.status === 'active' ? '禁用用户' : '启用用户'}>
                              {u.status === 'active'
                                ? <ToggleRight className="h-4 w-4 text-green-500" />
                                : <ToggleLeft className="h-4 w-4 text-red-400" />}
                              {u.status === 'active' ? '已启用' : '已禁用'}
                            </button>
                            <button onClick={() => openEdit(u)}
                              className="inline-flex items-center gap-1 text-xs font-medium text-gray-500 hover:text-primary transition-colors">
                              <Shield className="h-3.5 w-3.5" />
                              编辑
                            </button>
                          </div>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </div>

          <div className="mt-3 text-xs text-gray-400">
            共 {displayUsers.length} 条记录
            {search && ` — 搜索 "${search}"`}
          </div>
        </>
      )}

      {/* Role Edit Modal */}
      {editModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40" onClick={() => setEditModal(null)}>
          <div className="bg-white rounded-xl shadow-xl w-full max-w-md mx-4" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-center justify-between p-5 border-b border-gray-200">
              <h3 className="text-base font-semibold text-gray-900">编辑用户信息</h3>
              <button onClick={() => setEditModal(null)} className="p-1 text-gray-400 hover:text-gray-600">
                <X className="h-4 w-4" />
              </button>
            </div>
            <div className="p-5 space-y-4">
              <div className="grid grid-cols-1 gap-4">
                <label className="block text-sm font-medium text-gray-700">
                  用户名
                  <input value={editUsername} onChange={(e) => setEditUsername(e.target.value)}
                    className="mt-1.5 w-full rounded-md border border-gray-300 px-3 py-2 text-sm font-normal text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary" />
                </label>
                <label className="block text-sm font-medium text-gray-700">
                  邮箱
                  <input type="email" value={editEmail} onChange={(e) => setEditEmail(e.target.value)}
                    className="mt-1.5 w-full rounded-md border border-gray-300 px-3 py-2 text-sm font-normal text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary" />
                </label>
                <label className="block text-sm font-medium text-gray-700">
                  额度
                  <input type="number" min={0} value={editCredits} onChange={(e) => setEditCredits(parseInt(e.target.value, 10) || 0)}
                    className="mt-1.5 w-full rounded-md border border-gray-300 px-3 py-2 text-sm font-normal text-gray-900 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary" />
                </label>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">角色</label>
                <div className="space-y-2">
                  {ROLES.map((r) => (
                    <label key={r}
                      className={`flex items-center gap-3 p-3 rounded-md border cursor-pointer transition-colors
                        ${editRole === r
                          ? 'border-primary bg-primary/5'
                          : 'border-gray-200 hover:border-gray-300'
                        }`}>
                      <input type="radio" name="role" value={r} checked={editRole === r}
                        onChange={() => setEditRole(r)} className="text-primary focus:ring-primary" />
                      <div>
                        <div className="text-sm font-medium text-gray-900">{roleLabels[r] || r}</div>
                        <div className="text-xs text-gray-400">
                          {r === 'superadmin' ? '全部权限' : r === 'admin' ? '管理后台操作权限' : '普通用户，无管理权限'}
                        </div>
                      </div>
                    </label>
                  ))}
                </div>
              </div>
            </div>
            <div className="flex items-center justify-end gap-3 p-5 border-t border-gray-200">
              <button onClick={() => setEditModal(null)}
                className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900">取消</button>
              <button onClick={saveUser} disabled={saving || !editUsername.trim() || !editEmail.trim() || editCredits < 0}
                className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 disabled:opacity-60">
                {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Check className="h-4 w-4" />}
                保存
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
