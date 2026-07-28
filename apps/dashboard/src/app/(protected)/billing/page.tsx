'use client';

import { useEffect, useState } from 'react';
import { CreditCard, Key, Check, Loader2, Gift, Clock, History } from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface Plan {
  id: string;
  name: string;
  description: string;
  monthly_price: number;
  yearly_price: number;
  max_projects: number;
  max_ai_queries: number;
  credits: number;
  features: string;
  popular: boolean;
}

interface Order {
  id: string;
  plan_name: string;
  amount: number;
  status: string;
  created_at: string;
}

type Tab = 'plans' | 'cdk' | 'history';

export default function BillingPage() {
  const [tab, setTab] = useState<Tab>('plans');
  const [plans, setPlans] = useState<Plan[]>([]);
  const [orders, setOrders] = useState<Order[]>([]);
  const [credits, setCredits] = useState(0);
  const [loading, setLoading] = useState(true);
  const [purchasing, setPurchasing] = useState(false);
  const [selectedPlan, setSelectedPlan] = useState<string | null>(null);
  const [cdkCode, setCdkCode] = useState('');
  const [redeeming, setRedeeming] = useState(false);
  const [cdkResult, setCdkResult] = useState<{ credits: number; message: string } | null>(null);
  const [cdkError, setCdkError] = useState('');

  useEffect(() => {
    Promise.all([
      apiRequest<Plan[]>('/plans'),
      apiRequest<{ credits: number }>('/billing/credits'),
      apiRequest<Order[]>('/billing/orders'),
    ])
      .then(([plansRes, creditsRes, ordersRes]) => {
        if (Array.isArray(plansRes.data)) setPlans(plansRes.data);
        if (creditsRes.data) setCredits(creditsRes.data.credits ?? 0);
        if (Array.isArray(ordersRes.data)) setOrders(ordersRes.data);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const handlePurchase = (planId: string) => {
    setPurchasing(true);
    setSelectedPlan(planId);
    apiRequest<{ pay_url: string }>('/billing/orders', {
      method: 'POST',
      body: JSON.stringify({ plan_id: planId }),
    })
      .then((res) => {
        if (res.data?.pay_url) {
          window.location.href = res.data.pay_url;
        }
      })
      .catch(() => {})
      .finally(() => { setPurchasing(false); setSelectedPlan(null); });
  };

  const handleRedeem = () => {
    if (!cdkCode.trim()) return;
    setRedeeming(true);
    setCdkResult(null);
    setCdkError('');
    apiRequest<{ credits: number; message: string }>('/billing/cdk/redeem', {
      method: 'POST',
      body: JSON.stringify({ code: cdkCode.trim() }),
    })
      .then((res) => {
        if (res.data) {
          setCdkResult(res.data);
          setCredits((p) => p + res.data.credits);
        }
      })
      .catch((err) => setCdkError(err?.message || '兑换失败'))
      .finally(() => setRedeeming(false));
  };

  const tabs: { key: Tab; label: string; icon: typeof CreditCard }[] = [
    { key: 'plans', label: '套餐购买', icon: CreditCard },
    { key: 'cdk', label: 'CDK 兑换', icon: Key },
    { key: 'history', label: '消费记录', icon: History },
  ];

  const featureList = (features: string) =>
    features.split('\n').filter(Boolean).map((f, i) => (
      <li key={i} className="flex items-start gap-1.5 text-xs text-gray-600">
        <Check className="h-3.5 w-3.5 text-green-500 mt-0.5 shrink-0" />
        {f}
      </li>
    ));

  const statusLabel: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    completed: '已完成',
    expired: '已过期',
    failed: '失败',
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <CreditCard className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">账单中心</h2>
        </div>
        <div className="flex items-center gap-2 px-4 py-2 bg-green-50 border border-green-200 rounded-lg">
          <Gift className="h-4 w-4 text-green-600" />
          <span className="text-sm font-medium text-green-700">剩余额度：{loading ? '-' : credits}</span>
        </div>
      </div>

      <div className="flex items-center gap-3 mb-6">
        {tabs.map((t) => (
          <button key={t.key} onClick={() => setTab(t.key)}
            className={`flex items-center gap-1.5 px-4 py-2 text-sm font-medium rounded-md transition-colors ${tab === t.key ? 'bg-primary text-white' : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50'}`}>
            <t.icon className="h-4 w-4" />
            {t.label}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="p-12 text-center text-gray-400"><Loader2 className="animate-spin inline h-5 w-5 mr-2" />加载中...</div>
      ) : tab === 'plans' ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {plans.length === 0 ? (
            <div className="col-span-full p-12 text-center text-gray-400">暂无可购买套餐</div>
          ) : (
            plans.map((plan) => (
              <div key={plan.id} className={`bg-white rounded-lg border ${plan.popular ? 'border-primary/40 ring-1 ring-primary/20' : 'border-gray-200'} p-5 relative flex flex-col`}>
                {plan.popular && <span className="absolute -top-2.5 right-4 px-2 py-0.5 bg-primary text-white text-[10px] font-medium rounded-full">推荐</span>}
                <h3 className="text-base font-bold text-gray-900 mb-1">{plan.name}</h3>
                {plan.description && <p className="text-xs text-gray-400 mb-3">{plan.description}</p>}
                <div className="mb-4">
                  {plan.monthly_price === 0 && plan.yearly_price === 0 ? (
                    <span className="text-2xl font-bold text-gray-900">免费</span>
                  ) : (
                    <div>
                      <div className="flex items-baseline gap-1">
                        <span className="text-2xl font-bold text-gray-900">¥{plan.monthly_price}</span>
                        <span className="text-xs text-gray-400">/月</span>
                      </div>
                      <p className="text-xs text-gray-400 mt-0.5">年付 ¥{plan.yearly_price}</p>
                    </div>
                  )}
                </div>
                <div className="flex-1">
                  <ul className="space-y-1.5 mb-4">{featureList(plan.features)}</ul>
                </div>
                <button onClick={() => handlePurchase(plan.id)} disabled={purchasing}
                  className="w-full flex items-center justify-center gap-1.5 px-3 py-2 text-sm font-medium rounded-md bg-primary text-white hover:bg-primary/90 disabled:opacity-60 transition-colors">
                  {purchasing && selectedPlan === plan.id ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
                  立即购买
                </button>
              </div>
            ))
          )}
        </div>
      ) : tab === 'cdk' ? (
        <div className="max-w-lg">
          <div className="bg-white rounded-lg border border-gray-200 p-6">
            <h3 className="text-base font-semibold text-gray-900 mb-2">CDK 兑换</h3>
            <p className="text-sm text-gray-500 mb-4">输入 CDK 兑换码获取额外额度</p>
            <div className="flex items-center gap-3">
              <input type="text" value={cdkCode}
                onChange={(e) => { setCdkCode(e.target.value); setCdkResult(null); setCdkError(''); }}
                placeholder="输入 CDK 兑换码"
                className="flex-1 rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary text-gray-900" />
              <button onClick={handleRedeem} disabled={redeeming || !cdkCode.trim()}
                className="flex items-center gap-1.5 px-4 py-2 bg-primary text-white text-sm font-medium rounded-md hover:bg-primary/90 disabled:opacity-60 transition-colors shrink-0">
                {redeeming ? <Loader2 className="h-4 w-4 animate-spin" /> : <Gift className="h-4 w-4" />}
                兑换
              </button>
            </div>
            {cdkResult && (
              <div className="mt-4 p-3 bg-green-50 border border-green-200 rounded-md text-sm text-green-700 flex items-center gap-2">
                <Check className="h-4 w-4 shrink-0" />
                {cdkResult.message || `兑换成功，获得 ${cdkResult.credits} 额度`}
              </div>
            )}
            {cdkError && (
              <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-md text-sm text-red-700">{cdkError}</div>
            )}
          </div>
        </div>
      ) : (
        <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
          <div className="px-5 py-4 border-b border-gray-200">
            <h3 className="text-base font-semibold text-gray-900">消费记录</h3>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200 bg-gray-50 text-left">
                  <th className="px-5 py-3 font-semibold text-gray-600">套餐</th>
                  <th className="px-5 py-3 font-semibold text-gray-600">金额</th>
                  <th className="px-5 py-3 font-semibold text-gray-600">状态</th>
                  <th className="px-5 py-3 font-semibold text-gray-600">时间</th>
                </tr>
              </thead>
              <tbody>
                {orders.length === 0 ? (
                  <tr><td colSpan={4} className="px-5 py-12 text-center text-gray-400">暂无订单</td></tr>
                ) : (
                  orders.map((o) => (
                    <tr key={o.id} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                      <td className="px-5 py-3 font-medium text-gray-900">{o.plan_name}</td>
                      <td className="px-5 py-3 text-gray-700">¥{o.amount}</td>
                      <td className="px-5 py-3">
                        <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                          o.status === 'paid' || o.status === 'completed' ? 'bg-green-50 text-green-700' :
                          o.status === 'pending' ? 'bg-yellow-50 text-yellow-700' :
                          'bg-gray-100 text-gray-500'
                        }`}>
                          {statusLabel[o.status] || o.status}
                        </span>
                      </td>
                      <td className="px-5 py-3 text-gray-500 text-xs">{new Date(o.created_at).toLocaleString('zh-CN')}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}
