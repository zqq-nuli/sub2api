/**
 * Admin Settings API endpoints
 * Handles system settings management for administrators
 */

import { apiClient } from '../client'

/**
 * Payment channel configuration
 */
export interface PaymentChannel {
  key: string           // Channel identifier: alipay, wxpay, epusdt
  display_name: string  // Display name
  epay_type: string     // Gateway parameter: epay, alipay, wxpay, etc.
  icon: string          // Icon identifier
  enabled: boolean      // Whether enabled
  sort_order: number    // Sort order
}

/**
 * System settings interface
 */
export interface SystemSettings {
  // Registration settings
  registration_enabled: boolean
  email_verify_enabled: boolean
  // Default settings
  default_balance: number
  default_concurrency: number
  // OEM settings
  site_name: string
  site_logo: string
  site_subtitle: string
  api_base_url: string
  contact_info: string
  doc_url: string
  // SMTP settings
  smtp_host: string
  smtp_port: number
  smtp_username: string
  smtp_password: string
  smtp_from_email: string
  smtp_from_name: string
  smtp_use_tls: boolean
  // Cloudflare Turnstile settings
  turnstile_enabled: boolean
  turnstile_site_key: string
  turnstile_secret_key: string
  // SSO settings
  sso_enabled: boolean
  password_login_enabled: boolean
  sso_issuer_url: string
  sso_client_id: string
  sso_client_secret: string
  sso_redirect_uri: string
  sso_allowed_domains: string[]
  sso_auto_create_user: boolean
  sso_min_trust_level: number
  // Epay settings
  epay_enabled: boolean
  epay_api_url: string
  epay_merchant_id: string
  epay_merchant_key: string
  epay_notify_url: string
  epay_return_url: string
  // Payment channels
  payment_channels: PaymentChannel[]
}

/**
 * Get all system settings
 * @returns System settings
 */
export async function getSettings(): Promise<SystemSettings> {
  const { data } = await apiClient.get<SystemSettings>('/admin/settings')
  return data
}

/**
 * Update system settings
 * @param settings - Partial settings to update
 * @returns Updated settings
 */
export async function updateSettings(settings: Partial<SystemSettings>): Promise<SystemSettings> {
  const { data } = await apiClient.put<SystemSettings>('/admin/settings', settings)
  return data
}

/**
 * Test SMTP connection request
 */
export interface TestSmtpRequest {
  smtp_host: string
  smtp_port: number
  smtp_username: string
  smtp_password: string
  smtp_use_tls: boolean
}

/**
 * Test SMTP connection with provided config
 * @param config - SMTP configuration to test
 * @returns Test result message
 */
export async function testSmtpConnection(config: TestSmtpRequest): Promise<{ message: string }> {
  const { data } = await apiClient.post<{ message: string }>('/admin/settings/test-smtp', config)
  return data
}

/**
 * Send test email request
 */
export interface SendTestEmailRequest {
  email: string
  smtp_host: string
  smtp_port: number
  smtp_username: string
  smtp_password: string
  smtp_from_email: string
  smtp_from_name: string
  smtp_use_tls: boolean
}

/**
 * Send test email with provided SMTP config
 * @param request - Email address and SMTP config
 * @returns Test result message
 */
export async function sendTestEmail(request: SendTestEmailRequest): Promise<{ message: string }> {
  const { data } = await apiClient.post<{ message: string }>(
    '/admin/settings/send-test-email',
    request
  )
  return data
}

/**
 * Admin API Key status response
 */
export interface AdminApiKeyStatus {
  exists: boolean
  masked_key: string
}

/**
 * Get admin API key status
 * @returns Status indicating if key exists and masked version
 */
export async function getAdminApiKey(): Promise<AdminApiKeyStatus> {
  const { data } = await apiClient.get<AdminApiKeyStatus>('/admin/settings/admin-api-key')
  return data
}

/**
 * Regenerate admin API key
 * @returns The new full API key (only shown once)
 */
export async function regenerateAdminApiKey(): Promise<{ key: string }> {
  const { data } = await apiClient.post<{ key: string }>('/admin/settings/admin-api-key/regenerate')
  return data
}

/**
 * Delete admin API key
 * @returns Success message
 */
export async function deleteAdminApiKey(): Promise<{ message: string }> {
  const { data} = await apiClient.delete<{ message: string }>('/admin/settings/admin-api-key')
  return data
}

/**
 * Test SSO connection request
 */
export interface TestSSORequest {
  issuer_url: string
}

/**
 * Test SSO connection with provided issuer URL
 * @param config - SSO configuration to test
 * @returns Test result message
 */
export async function testSSOConnection(config: TestSSORequest): Promise<{ message: string; issuer: string }> {
  const { data } = await apiClient.post<{ message: string; issuer: string }>('/admin/settings/test-sso', config)
  return data
}

/**
 * Update single setting (for real-time save)
 * @param key - Setting key
 * @param value - Setting value
 * @returns Success message
 */
export async function updateSingleSetting(key: string, value: any): Promise<{ message: string }> {
  const { data } = await apiClient.patch<{ message: string }>(`/admin/settings/${key}`, { value })
  return data
}

export const settingsAPI = {
  getSettings,
  updateSettings,
  testSmtpConnection,
  sendTestEmail,
  getAdminApiKey,
  regenerateAdminApiKey,
  deleteAdminApiKey,
  testSSOConnection,
  updateSingleSetting
}

export default settingsAPI
