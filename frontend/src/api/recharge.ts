/**
 * Recharge Products API endpoints
 * Handles recharge product listing for users
 */

import { apiClient } from './client'

export interface RechargeProduct {
  id: number
  name: string
  amount: number
  balance: number
  bonus_balance: number
  total_balance: number
  description: string
  sort_order: number
  is_hot: boolean
  discount_label: string
  status: string
  created_at: string
  updated_at: string
}

export interface PaymentChannel {
  key: string
  display_name: string
  icon: string
}

/**
 * Get active recharge products
 * @returns List of active recharge products
 */
export async function getProducts(): Promise<RechargeProduct[]> {
  const { data } = await apiClient.get<RechargeProduct[]>('/recharge/products')
  return data
}

/**
 * Get available payment channels
 * @returns List of enabled payment channels
 */
export async function getPaymentChannels(): Promise<PaymentChannel[]> {
  const { data } = await apiClient.get<PaymentChannel[]>('/payment/channels')
  return data
}

export const rechargeAPI = {
  getProducts,
  getPaymentChannels
}

export default rechargeAPI
