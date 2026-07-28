'use client';
import { useState, FormEvent } from 'react';
import { useAuth } from '@/lib/auth';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

type LoginMode = 'email' | 'phone';

export default function LoginPage() {
  const { login } = useAuth();
  const router = useRouter();
  const [mode, setMode] = useState<LoginMode>('email');
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    const value = identifier.trim();
    if (!value) {
      setError(mode === 'email' ? '请输入邮箱地址' : '请输入手机号');
      return;
    }
    if (mode === 'email' && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
      setError('请输入有效的邮箱地址');
      return;
    }
    if (!password) {
      setError('请输入密码');
      return;
    }
    setSubmitting(true);
    try {
      await login(value, password, mode);
      router.push('/dashboard');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : '登录失败';
      setError(message);
    } finally {
      setSubmitting(false);
    }
  };

  const inputClass = 'w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors';

  return (
    <div className="bg-white rounded-card shadow-sm border border-gray-100 p-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-2">登录 AI Growth Engine</h1>
      <p className="text-sm text-gray-500 mb-6">输入您的凭据以访问账户</p>

      {error && (
        <div className="mb-4 p-3 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
          {error}
        </div>
      )}

      {/* Tab switcher */}
      <div className="flex mb-6 border-b border-gray-200">
        <button
          type="button"
          onClick={() => { setMode('email'); setError(''); }}
          className={`pb-2 px-4 text-sm font-medium border-b-2 transition-colors ${
            mode === 'email'
              ? 'border-primary text-primary'
              : 'border-transparent text-gray-500 hover:text-gray-700'
          }`}
        >
          邮箱登录
        </button>
        <button
          type="button"
          onClick={() => { setMode('phone'); setError(''); }}
          className={`pb-2 px-4 text-sm font-medium border-b-2 transition-colors ${
            mode === 'phone'
              ? 'border-primary text-primary'
              : 'border-transparent text-gray-500 hover:text-gray-700'
          }`}
        >
          手机号登录
        </button>
      </div>

      <form onSubmit={handleSubmit} noValidate className="space-y-4">
        <div>
          <label htmlFor="identifier" className="block text-sm font-medium text-gray-700 mb-1">
            {mode === 'email' ? '邮箱' : '手机号'}
          </label>
          <input
            id="identifier"
            type={mode === 'email' ? 'text' : 'text'}
            inputMode={mode === 'email' ? 'email' : 'tel'}
            required
            value={identifier}
            onChange={(e) => setIdentifier(e.target.value)}
            placeholder={mode === 'email' ? 'you@example.com' : '13800138000'}
            className={inputClass}
          />
        </div>

        <div>
          <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
            密码
          </label>
          <input
            id="password"
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="请输入密码"
            className={inputClass}
          />
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full py-2.5 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {submitting ? '登录中...' : '登录'}
        </button>
      </form>

      <div className="mt-6 flex items-center justify-between text-sm">
        <Link href="/register" className="text-primary hover:underline">
          创建账户
        </Link>
        <Link href="/forgot-password" className="text-gray-500 hover:text-primary transition-colors">
          忘记密码？
        </Link>
      </div>
    </div>
  );
}
