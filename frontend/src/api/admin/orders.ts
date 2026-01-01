/**
 * Admin Orders API endpoints
 * Handles order management for administrators
 */

import { apiClient } from '../client'

export interface AdminOrder {
  id: number
  order_no: string
  user_id: number
  user_email?: string
  product_id: number | null
  product_name: string
  amount: number
  bonus_amount: number
  actual_amount: number
  payment_method: string
  payment_gateway: string
  trade_no: string
  status: string
  created_at: string
  paid_at: string | null
  expired_at: string
  notes: string
}

export interface AdminListOrdersResponse {
  orders: AdminOrder[]
  total: number
  page: number
  limit: number
}

export interface OrderStatistics {
  total_orders: number
  pending_orders: number
  paid_orders: number
  failed_orders: number
  expired_orders: number
  total_amount: number
  total_balance: number
}

export interface ListOrdersParams {
  page?: number
  limit?: number
  user_id?: number
  status?: string
  payment_method?: string
  order_no?: string
  start_date?: string
  end_date?: string
}

/**
 * List all orders with filters
 * @param params - Query parameters
 * @returns Paginated list of orders
 */
export async function listOrders(params: ListOrdersParams = {}): Promise<AdminListOrdersResponse> {
  const { data } = await apiClient.get<AdminListOrdersResponse>('/admin/orders', {
    params
  })
  return data
}

/**
 * Get order by ID
 * @param id - Order ID
 * @returns Order details
 */
export async function getOrderById(id: number): Promise<AdminOrder> {
  const { data } = await apiClient.get<AdminOrder>(`/admin/orders/${id}`)
  return data
}

/**
 * Update order status
 * @param id - Order ID
 * @param status - New status
 * @param notes - Optional notes
 * @returns Success message
 */
export async function updateOrderStatus(
  id: number,
  status: string,
  notes?: string
): Promise<{ message: string }> {
  const { data } = await apiClient.put<{ message: string }>(`/admin/orders/${id}`, {
    status,
    notes
  })
  return data
}

/**
 * Get order statistics
 * @returns Order statistics
 */
export async function getStatistics(): Promise<OrderStatistics> {
  const { data } = await apiClient.get<OrderStatistics>('/admin/orders/statistics')
  return data
}

export const adminOrdersAPI = {
  listOrders,
  getOrderById,
  updateOrderStatus,
  getStatistics
}

export default adminOrdersAPI
