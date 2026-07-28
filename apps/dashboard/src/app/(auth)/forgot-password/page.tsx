'use client';
import { useState, FormEvent } from 'react';
import { apiRequest } from '@/lib/api';
import Link from 'next/link';

type Step = 'email' | 'reset' | 'success';

export default function ForgotPasswordPage() {
  const [step, setStep] = useState<Step>('email');
  const [email, setEmail] = useState('');
  const [code, setCode] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSendCode = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    setSubmitting(true);
    try {
      await apiRequest('/auth/send-code', {
        method: 'POST',
        body: JSON.stringify({ email }),
      });
      setStep('reset');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : '发送验证码失败';
      setError(message);
    } finally {
      setSubmitting(false);
    }
  };

  const handleReset = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    setSubmitting(true);
    try {
      await apiRequest('/auth/reset-password', {
        method: 'POST',
        body: JSON.stringify({ email, code, new_password: newPassword }),
      });
      setStep('success');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : '重置密码失败';
      setError(message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-white rounded-card shadow-sm border border-gray-100 p-8">
      {step === 'email' && (
        <>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">重置密码</h1>
          <p className="text-sm text-gray-500 mb-6">输入邮箱地址以接收验证码</p>

          {error && (
            <div className="mb-4 p-3 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
              {error}
            </div>
          )}

          <form onSubmit={handleSendCode} noValidate className="space-y-4">
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
                className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors"
              />
            </div>
            <button
              type="submit"
              disabled={submitting}
              className="w-full py-2.5 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {submitting ? '发送中...' : '发送验证码'}
            </button>
          </form>

          <p className="mt-6 text-center text-sm text-gray-500">
            <Link href="/login" className="text-primary hover:underline font-medium">
              返回登录
            </Link>
          </p>
        </>
      )}

      {step === 'reset' && (
        <>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">输入验证码</h1>
          <p className="text-sm text-gray-500 mb-6">
            验证码已发送至 <span className="font-medium text-gray-700">{email}</span>
          </p>

          {error && (
            <div className="mb-4 p-3 rounded-btn bg-red-50 border border-red-200 text-sm text-red-700">
              {error}
            </div>
          )}

          <form onSubmit={handleReset} noValidate className="space-y-4">
            <div>
              <label htmlFor="code" className="block text-sm font-medium text-gray-700 mb-1">
                验证码
              </label>
              <input
                id="code"
                type="text"
                required
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="请输入6位验证码"
                className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors"
              />
            </div>
            <div>
              <label htmlFor="newPassword" className="block text-sm font-medium text-gray-700 mb-1">
                新密码
              </label>
              <input
                id="newPassword"
                type="password"
                required
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                placeholder="至少8个字符"
                className="w-full px-3 py-2 border border-gray-300 rounded-btn text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-colors"
              />
            </div>
            <button
              type="submit"
              disabled={submitting}
              className="w-full py-2.5 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {submitting ? '重置中...' : '重置密码'}
            </button>
          </form>

          <button
            onClick={() => setStep('email')}
            className="mt-4 w-full text-sm text-gray-500 hover:text-primary transition-colors"
          >
            更换邮箱
          </button>
        </>
      )}

      {step === 'success' && (
        <div className="text-center">
          <div className="w-12 h-12 bg-success/10 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-6 h-6 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">密码重置成功</h1>
          <p className="text-sm text-gray-500 mb-6">密码已更新，请使用新密码登录。</p>
          <Link
            href="/login"
            className="inline-block w-full py-2.5 bg-primary text-white text-sm font-medium rounded-btn hover:bg-blue-700 transition-colors text-center"
          >
            去登录
          </Link>
        </div>
      )}
    </div>
  );
}
