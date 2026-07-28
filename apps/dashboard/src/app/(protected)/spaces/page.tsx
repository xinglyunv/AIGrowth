'use client';

import { useEffect, useState } from 'react';
import { UserPlus, UsersRound, Loader2 } from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface Space { id: string; name: string; slug: string; status: string; }
interface Member { user_id: string; email: string; username: string; role: string; status: string; }

export default function SpacesPage() {
  const [spaces, setSpaces] = useState<Space[]>([]);
  const [current, setCurrent] = useState<Space | null>(null);
  const [members, setMembers] = useState<Member[]>([]);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const load = async () => {
    setLoading(true);
    try {
      const [spacesRes, currentRes] = await Promise.all([
        apiRequest<Space[]>('/spaces'),
        apiRequest<Space>('/spaces/current'),
      ]);
      const list = Array.isArray(spacesRes.data) ? spacesRes.data : [];
      setSpaces(list);
      const active = currentRes.data || list[0] || null;
      setCurrent(active);
      if (active) {
        const memberRes = await apiRequest<Member[]>(`/spaces/${active.id}/members`);
        setMembers(Array.isArray(memberRes.data) ? memberRes.data : []);
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : '加载团队空间失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const createSpace = async () => {
    if (!name.trim()) return;
    try {
      await apiRequest('/spaces', { method: 'POST', body: JSON.stringify({ name: name.trim() }) });
      setName('');
      await load();
    } catch (e) { setError(e instanceof Error ? e.message : '创建失败'); }
  };

  const inviteMember = async () => {
    if (!current || !email.trim()) return;
    try {
      await apiRequest(`/spaces/${current.id}/members/invite`, { method: 'POST', body: JSON.stringify({ email: email.trim(), role: 'member' }) });
      setEmail('');
      setError('邀请已创建');
    } catch (e) { setError(e instanceof Error ? e.message : '邀请失败'); }
  };

  if (loading) return <div className="flex items-center justify-center py-20"><Loader2 className="h-6 w-6 animate-spin text-gray-400" /></div>;

  return (
    <div className="max-w-5xl mx-auto">
      <div className="flex items-center gap-3 mb-6"><UsersRound className="h-6 w-6 text-primary" /><h1 className="text-2xl font-bold text-gray-900">团队空间</h1></div>
      {error && <div className="mb-4 rounded-md border border-blue-200 bg-blue-50 px-4 py-3 text-sm text-blue-700">{error}</div>}
      <div className="grid gap-6 lg:grid-cols-3">
        <section className="rounded-xl border border-gray-200 bg-white p-5">
          <h2 className="mb-4 font-semibold text-gray-900">我的空间</h2>
          <div className="space-y-2">{spaces.map((space) => <button key={space.id} onClick={async () => { await apiRequest(`/spaces/${space.id}/current`, { method: 'PUT' }); await load(); }} className={`w-full rounded-lg border px-3 py-2 text-left text-sm ${current?.id === space.id ? 'border-primary bg-primary/5 text-primary' : 'border-gray-200 text-gray-700'}`}>{space.name}</button>)}</div>
          <div className="mt-5 flex gap-2"><input value={name} onChange={(e) => setName(e.target.value)} placeholder="新空间名称" className="min-w-0 flex-1 rounded-md border border-gray-300 px-3 py-2 text-sm" /><button onClick={createSpace} className="rounded-md bg-primary px-3 py-2 text-sm font-medium text-white">创建</button></div>
        </section>
        <section className="rounded-xl border border-gray-200 bg-white p-5 lg:col-span-2">
          <div className="mb-4 flex items-center justify-between"><h2 className="font-semibold text-gray-900">{current?.name || '成员'} 成员</h2><span className="text-sm text-gray-400">{members.length} 人</span></div>
          <div className="divide-y divide-gray-100">{members.map((member) => <div key={member.user_id} className="flex items-center justify-between py-3"><div><p className="text-sm font-medium text-gray-900">{member.username || member.email}</p><p className="text-xs text-gray-400">{member.email}</p></div><span className="rounded-full bg-gray-100 px-2 py-1 text-xs text-gray-600">{member.role}</span></div>)}</div>
          <div className="mt-5 flex gap-2 border-t border-gray-100 pt-5"><input type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="成员邮箱" className="min-w-0 flex-1 rounded-md border border-gray-300 px-3 py-2 text-sm" /><button onClick={inviteMember} disabled={!current} className="flex items-center gap-1 rounded-md bg-primary px-3 py-2 text-sm font-medium text-white disabled:opacity-50"><UserPlus className="h-4 w-4" />邀请成员</button></div>
        </section>
      </div>
    </div>
  );
}
