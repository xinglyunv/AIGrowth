'use client';
import { useState, FormEvent, useEffect } from 'react';
import { useAuth } from '@/lib/auth';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { apiRequest } from '@/lib/api';

export default function RegisterPage() {
  const { register } = useAuth();
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [companyName, setCompanyName] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [regAllowed, setRegAllowed] = useState<boolean | null>(null);

  useEffect(() => {
    apiRequest<Record<string, string>>('/settings')
      .then((res) => {
        if (res.data && res.data.allow_registration === 'false') {
          setRegAllowed(false);
        } else {
          setRegAllowed(true);
        }
      })
      .catch(() => setRegAllowed(true));
  }, []);

  const validate = (): boolean => {
    const errors: Record<string, string> = {};
    if (!email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      errors.email = '请输入有效的邮箱地址';
    }
    if (!password || password.length < 8) {
      errors.password = '密码至少需要8个字符';
    }
    if (password !== confirmPassword) {
      errors.confirmPassword = '两次输入的密码不一致';
    }
    setFieldErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    if (!validate()) return;
    setSubmitting(true);
    try {
      await register({ email, password, username, company_name: companyName || undefined });
      router.push('/dashboard');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : '注册失败';
      setError(message);
    } finally {
      setSubmitting(false);
    }
  };

  const inputClass = (hasError: boolean) =>
    `w-full px-3 py-2 border rounded-btn text-sm focus:outline-none focus:ring-2 transition-colors ${
      hasError
        ? 'border-red-400 focus:ring-red/20 focus:border-red-400'
        : 'border-gray-300 focus:ring-primary/20 focus:border-primary'
    }`;

  if (regAllowed === null) {
    return (
      <div className="bg-white rounded-card shadow-sm border border-gray-100 p-8">
        <p className="text-center text-gray-500">加载中...</p>
      </div>
    );
  }

  if (!regAllowed) {
    return (
      <div className="bg-white rounded-card shadow-sm border border-gray-100 p-8 text-center">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">注册已关闭</h1>
        <p className="text-gray-500 mb-6">当前系统暂未开放注册，请联系管理员。</p>
        <Link href="/login" className="text-primary hover:underline font-medium">
          返回登录
        </Link>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-card shadow-sm border border-gray-100 p-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-2">创建账户</h1>
      <p className="text-sm text-gray-500 mb-6">开始使用 AI Growth Engine</p>

      {error && (
        <div className="mb-4 p-3 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} noValidate className="space-y-4">
        <div>
          <label htmlFor="username" className="block text-sm font-medium text-gray-700 mb-1">
            用户名
          </label>
          <input
            id="username"
            type="text"
            required
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="请输入用户名"
            className={inputClass(false)}
          />
        </div>

        <div>
          <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
            邮箱
          </label>
          <input
            id="email"
            type="text"
            inputMode="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="you@example.com"
            className={inputClass(!!fieldErrors.email)}
          />
          {fieldErrors.email && <p className="mt-1 text-xs text-red-600">{fieldErrors.email}</p>}
        </div>

        <div>
          <label htmlFor="companyName" className="block text-sm font-medium text-gray-700 mb-1">
            公司名称 <span className="text-gray-400 font-normal">(选填)</span>
          </label>
          <input
            id="companyName"
            type="text"
            value={companyName}
            onChange={(e) => setCompanyName(e.target.value)}
            placeholder="请输入公司名称"
            className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors"
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
            placeholder="至少8个字符"
            className={inputClass(!!fieldErrors.password)}
          />
          {fieldErrors.password && <p className="mt-1 text-xs text-red-600">{fieldErrors.password}</p>}
        </div>

        <div>
          <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 mb-1">
            确认密码
          </label>
          <input
            id="confirmPassword"
            type="password"
            required
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            placeholder="请再次输入密码"
            className={inputClass(!!fieldErrors.confirmPassword)}
          />
          {fieldErrors.confirmPassword && <p className="mt-1 text-xs text-red-600">{fieldErrors.confirmPassword}</p>}
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full py-2.5 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {submitting ? '创建中...' : '创建账户'}
        </button>
      </form>

      <p className="mt-6 text-center text-sm text-gray-500">
        已有账户？
        <Link href="/login" className="text-primary hover:underline font-medium">
          立即登录
        </Link>
      </p>
    </div>
  );
}
