'use client';

import React, { useEffect, useState, useCallback } from 'react';
import { ShoppingCart, Search, ChevronDown, ChevronUp, RefreshCw, Filter, Loader2 } from 'lucide-react';
import { apiRequest } from '@/lib/api';

interface OrderItem {
  name: string;
  price: string;
  quantity: number;
}

interface Order {
  id: string;
  order_no: string;
  user_id: string;
  user_name?: string;
  description?: string;
  plan_id?: string;
  amount: number;
  status: string;
  payment_method: string;
  created_at: string;
  items?: OrderItem[];
}



const STATUS_OPTIONS = ['all', 'paid', 'pending', 'cancelled', 'refunded'];

const statusConfig: Record<string, { label: string; bg: string; text: string }> = {
  paid: { label: '已完成', bg: 'bg-green-50', text: 'text-green-700' },
  pending: { label: '待支付', bg: 'bg-yellow-50', text: 'text-yellow-700' },
  cancelled: { label: '已取消', bg: 'bg-gray-100', text: 'text-gray-600' },
  refunded: { label: '已退款', bg: 'bg-red-50', text: 'text-red-700' },
};

const statusLabels: Record<string, string> = { all: '全部', paid: '已完成', pending: '待支付', cancelled: '已取消', refunded: '已退款' };

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');
  const [sortAsc, setSortAsc] = useState(false);
  const [statusFilter, setStatusFilter] = useState('all');
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    setError('');
    apiRequest<Order[]>('/orders')
      .then((r) => { if (Array.isArray(r.data)) setOrders(r.data); })
      .catch((e) => setError(e instanceof Error ? e.message : '加载失败'))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => { load(); }, [load]);

  const displayOrders = orders
    .filter((o) => statusFilter === 'all' || o.status === statusFilter)
    .filter((o) => !search || o.order_no.includes(search) || (o.user_name ?? '').includes(search) || (o.description ?? '').includes(search) || o.user_id.includes(search))
    .sort((a, b) => sortAsc ? a.created_at.localeCompare(b.created_at) : b.created_at.localeCompare(a.created_at));

  const totalAmount = displayOrders
    .filter((o) => o.status === 'paid')
    .reduce((sum, o) => sum + Number(o.amount || 0), 0);

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <ShoppingCart className="h-6 w-6 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">订单管理</h2>
          <span className="text-sm text-gray-400 font-normal">共 {orders.length} 笔订单</span>
        </div>
        <div className="flex items-center gap-2">
          <button onClick={load}
            className="flex items-center gap-1.5 px-3 py-2 text-sm text-gray-600 hover:text-gray-900">
            <RefreshCw className="h-4 w-4" />
            刷新
          </button>
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
            <input type="text" placeholder="搜索订单号、用户或套餐..." value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-9 pr-4 py-2 border border-gray-300 rounded-md text-sm
                focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary
                text-gray-900 placeholder-gray-400 w-64" />
          </div>
        </div>
      </div>

      {loading && <div className="p-12 text-center text-gray-400"><Loader2 className="animate-spin inline h-5 w-5" /> 加载中...</div>}

      {error && <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md text-sm text-red-700">{error}</div>}

      {!loading && !error && (<>
      <div className="flex items-center gap-2 mb-4 flex-wrap">
        <Filter className="h-4 w-4 text-gray-400" />
        {STATUS_OPTIONS.map((s) => (
          <button key={s} onClick={() => setStatusFilter(s)}
            className={`px-3 py-1.5 text-sm font-medium rounded-md transition-colors ${
              statusFilter === s ? 'bg-primary text-white' : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50'
            }`}>
            {statusLabels[s]}
          </button>
        ))}
      </div>

      <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 bg-gray-50 text-left">
                <th className="px-4 py-3 font-semibold text-gray-600 w-8"></th>
                <th className="px-4 py-3 font-semibold text-gray-600">订单号</th>
                <th className="px-4 py-3 font-semibold text-gray-600">用户</th>
                <th className="px-4 py-3 font-semibold text-gray-600">套餐</th>
                <th className="px-4 py-3 font-semibold text-gray-600">金额</th>
                <th className="px-4 py-3 font-semibold text-gray-600">状态</th>
                <th className="px-4 py-3 font-semibold text-gray-600">
                  <button onClick={() => setSortAsc(!sortAsc)}
                    className="flex items-center gap-1 hover:text-gray-900">
                    创建时间
                    {sortAsc ? <ChevronUp className="h-3 w-3" /> : <ChevronDown className="h-3 w-3" />}
                  </button>
                </th>
              </tr>
            </thead>
            <tbody>
              {displayOrders.length === 0 ? (
                <tr>
                  <td colSpan={7} className="px-4 py-12 text-center text-gray-400">暂无数据</td>
                </tr>
              ) : (
                displayOrders.map((o) => {
                const cfg = statusConfig[o.status];
                const isExpanded = expandedId === o.id;
                return (
                  <React.Fragment key={o.id}>
                    <tr className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                      <td className="px-4 py-3">
                        <button onClick={() => setExpandedId(isExpanded ? null : o.id)}
                          className="p-0.5 text-gray-400 hover:text-gray-600">
                          {isExpanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                        </button>
                      </td>
                       <td className="px-4 py-3 font-mono text-xs text-gray-700">{o.order_no}</td>
                       <td className="px-4 py-3">
                         <div>
                           <p className="text-sm font-medium text-gray-900">{o.user_name || o.user_id}</p>
                           <p className="text-xs text-gray-400">{o.user_id}</p>
                         </div>
                       </td>
                       <td className="px-4 py-3 text-gray-600">{o.description || o.plan_id || '-'}</td>
                       <td className="px-4 py-3 font-semibold text-gray-900">¥{Number(o.amount || 0).toFixed(2)}</td>
                      <td className="px-4 py-3">
                        <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium ${cfg.bg} ${cfg.text}`}>
                          {cfg.label}
                        </span>
                      </td>
                       <td className="px-4 py-3 text-gray-500 text-xs">{o.created_at}</td>
                    </tr>
                    {isExpanded && (
                      <tr key={`${o.id}-expanded`}>
                        <td colSpan={7} className="px-4 py-4 bg-gray-50 border-b border-gray-100">
                          <div className="grid grid-cols-2 gap-4 text-sm">
                            <div>
                              <p className="text-xs text-gray-400 mb-1">支付方式</p>
                              <p className="text-sm text-gray-700">{o.payment_method}</p>
                            </div>
                            <div>
                              <p className="text-xs text-gray-400 mb-1">订单状态</p>
                              <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium ${cfg.bg} ${cfg.text}`}>{cfg.label}</span>
                            </div>
                          </div>
                          {o.items && o.items.length > 0 && (
                            <div className="mt-3">
                              <p className="text-xs text-gray-400 mb-1">订单明细</p>
                              <div className="bg-white rounded-md border border-gray-200 divide-y divide-gray-100">
                                {o.items.map((item, idx) => (
                                  <div key={idx} className="flex items-center justify-between px-3 py-2 text-sm">
                                    <span className="text-gray-700">{item.name}</span>
                                    <div className="flex items-center gap-4">
                                      <span className="text-gray-400">×{item.quantity}</span>
                                      <span className="font-medium text-gray-900">{item.price}</span>
                                    </div>
                                  </div>
                                ))}
                              </div>
                            </div>
                          )}
                        </td>
                      </tr>
                    )}
                  </React.Fragment>
                );
              }))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="mt-3 flex items-center justify-between text-xs text-gray-400">
        <span>共 {displayOrders.length} 条记录{search ? ` — 搜索"${search}"` : ''}</span>
        <span>已完成订单总额：¥{totalAmount.toFixed(2)}</span>
      </div>
      </>)}
    </div>
  );
}
