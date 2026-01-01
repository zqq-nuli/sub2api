/**
 * Orders API endpoints
 * Handles order creation and listing for users
 */

import { apiClient } from './client'

export interface Order {
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

export interface CreateOrderRequest {
  product_id: number
  payment_method: string
  return_url?: string
}

export interface CreateOrderResponse {
  order: Order
  payment_url: string
  payment_params?: Record<string, string>  // 表单参数（用于POST提交）
}

export interface GetOrdersResponse {
  orders: Order[]
  total: number
  page: number
  limit: number
}

/**
 * Create a recharge order
 * @param request - Order creation request
 * @returns Created order and payment URL
 */
export async function createOrder(request: CreateOrderRequest): Promise<CreateOrderResponse> {
  const { data } = await apiClient.post<CreateOrderResponse>('/orders', request)
  return data
}

/**
 * Get order by order number
 * @param orderNo - Order number
 * @returns Order details
 */
export async function getOrderByNo(orderNo: string): Promise<Order> {
  const { data } = await apiClient.get<Order>(`/orders/${orderNo}`)
  return data
}

/**
 * Get user's orders
 * @param page - Page number
 * @param limit - Items per page
 * @returns Paginated list of orders
 */
export async function getOrders(page: number = 1, limit: number = 20): Promise<GetOrdersResponse> {
  const { data } = await apiClient.get<GetOrdersResponse>('/orders', {
    params: { page, limit }
  })
  return data
}

export const ordersAPI = {
  createOrder,
  getOrderByNo,
  getOrders
}

export default ordersAPI
