/**
 * Admin Recharge Products API endpoints
 * Handles recharge product management for administrators
 */

import { apiClient } from '../client'

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

export interface CreateRechargeProductRequest {
  name: string
  amount: number
  balance: number
  bonus_balance?: number
  description?: string
  sort_order?: number
  is_hot?: boolean
  discount_label?: string
}

export interface UpdateRechargeProductRequest {
  name?: string
  amount?: number
  balance?: number
  bonus_balance?: number
  description?: string
  sort_order?: number
  is_hot?: boolean
  discount_label?: string
  status?: string
}

/**
 * List all recharge products
 * @returns List of all recharge products
 */
export async function listProducts(): Promise<RechargeProduct[]> {
  const { data } = await apiClient.get<RechargeProduct[]>('/admin/recharge-products')
  return data
}

/**
 * Get product by ID
 * @param id - Product ID
 * @returns Product details
 */
export async function getProductById(id: number): Promise<RechargeProduct> {
  const { data } = await apiClient.get<RechargeProduct>(`/admin/recharge-products/${id}`)
  return data
}

/**
 * Create a new recharge product
 * @param request - Product creation request
 * @returns Created product
 */
export async function createProduct(request: CreateRechargeProductRequest): Promise<RechargeProduct> {
  const { data } = await apiClient.post<RechargeProduct>('/admin/recharge-products', request)
  return data
}

/**
 * Update a recharge product
 * @param id - Product ID
 * @param request - Product update request
 * @returns Updated product
 */
export async function updateProduct(
  id: number,
  request: UpdateRechargeProductRequest
): Promise<RechargeProduct> {
  const { data } = await apiClient.put<RechargeProduct>(`/admin/recharge-products/${id}`, request)
  return data
}

/**
 * Delete a recharge product
 * @param id - Product ID
 * @returns Success message
 */
export async function deleteProduct(id: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(`/admin/recharge-products/${id}`)
  return data
}

export const adminRechargeProductsAPI = {
  listProducts,
  getProductById,
  createProduct,
  updateProduct,
  deleteProduct
}

export default adminRechargeProductsAPI
