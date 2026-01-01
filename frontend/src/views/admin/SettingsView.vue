<template>
  <AppLayout>
    <div class="mx-auto max-w-4xl space-y-6">
      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
      </div>

      <!-- Settings Form -->
      <form v-else @submit.prevent="saveSettings" class="space-y-6">
        <!-- Admin API Key Settings -->
        <div class="card">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('admin.settings.adminApiKey.title') }}
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {{ t('admin.settings.adminApiKey.description') }}
            </p>
          </div>
          <div class="space-y-4 p-6">
            <!-- Security Warning -->
            <div
              class="rounded-lg border border-amber-200 bg-amber-50 p-4 dark:border-amber-800 dark:bg-amber-900/20"
            >
              <div class="flex items-start">
                <svg
                  class="mt-0.5 h-5 w-5 flex-shrink-0 text-amber-500"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    fill-rule="evenodd"
                    d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                    clip-rule="evenodd"
                  />
                </svg>
                <p class="ml-3 text-sm text-amber-700 dark:text-amber-300">
                  {{ t('admin.settings.adminApiKey.securityWarning') }}
                </p>
              </div>
            </div>

            <!-- Loading State -->
            <div v-if="adminApiKeyLoading" class="flex items-center gap-2 text-gray-500">
              <div class="h-4 w-4 animate-spin rounded-full border-b-2 border-primary-600"></div>
              {{ t('common.loading') }}
            </div>

            <!-- No Key Configured -->
            <div v-else-if="!adminApiKeyExists" class="flex items-center justify-between">
              <span class="text-gray-500 dark:text-gray-400">
                {{ t('admin.settings.adminApiKey.notConfigured') }}
              </span>
              <button
                type="button"
                @click="createAdminApiKey"
                :disabled="adminApiKeyOperating"
                class="btn btn-primary btn-sm"
              >
                <svg
                  v-if="adminApiKeyOperating"
                  class="mr-1 h-4 w-4 animate-spin"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                  ></circle>
                  <path
                    class="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                {{
                  adminApiKeyOperating
                    ? t('admin.settings.adminApiKey.creating')
                    : t('admin.settings.adminApiKey.create')
                }}
              </button>
            </div>

            <!-- Key Exists -->
            <div v-else class="space-y-4">
              <div class="flex items-center justify-between">
                <div>
                  <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {{ t('admin.settings.adminApiKey.currentKey') }}
                  </label>
                  <code
                    class="rounded bg-gray-100 px-2 py-1 font-mono text-sm text-gray-900 dark:bg-dark-700 dark:text-gray-100"
                  >
                    {{ adminApiKeyMasked }}
                  </code>
                </div>
                <div class="flex gap-2">
                  <button
                    type="button"
                    @click="regenerateAdminApiKey"
                    :disabled="adminApiKeyOperating"
                    class="btn btn-secondary btn-sm"
                  >
                    {{
                      adminApiKeyOperating
                        ? t('admin.settings.adminApiKey.regenerating')
                        : t('admin.settings.adminApiKey.regenerate')
                    }}
                  </button>
                  <button
                    type="button"
                    @click="deleteAdminApiKey"
                    :disabled="adminApiKeyOperating"
                    class="btn btn-secondary btn-sm text-red-600 hover:text-red-700 dark:text-red-400"
                  >
                    {{ t('admin.settings.adminApiKey.delete') }}
                  </button>
                </div>
              </div>

              <!-- Newly Generated Key Display -->
              <div
                v-if="newAdminApiKey"
                class="space-y-3 rounded-lg border border-green-200 bg-green-50 p-4 dark:border-green-800 dark:bg-green-900/20"
              >
                <p class="text-sm font-medium text-green-700 dark:text-green-300">
                  {{ t('admin.settings.adminApiKey.keyWarning') }}
                </p>
                <div class="flex items-center gap-2">
                  <code
                    class="flex-1 select-all break-all rounded border border-green-300 bg-white px-3 py-2 font-mono text-sm dark:border-green-700 dark:bg-dark-800"
                  >
                    {{ newAdminApiKey }}
                  </code>
                  <button
                    type="button"
                    @click="copyNewKey"
                    class="btn btn-primary btn-sm flex-shrink-0"
                  >
                    {{ t('admin.settings.adminApiKey.copyKey') }}
                  </button>
                </div>
                <p class="text-xs text-green-600 dark:text-green-400">
                  {{ t('admin.settings.adminApiKey.usage') }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Registration Settings -->
        <div class="card">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('admin.settings.registration.title') }}
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {{ t('admin.settings.registration.description') }}
            </p>
          </div>
          <div class="space-y-5 p-6">
            <!-- Enable Registration -->
            <div class="flex items-center justify-between">
              <div>
                <label class="font-medium text-gray-900 dark:text-white">{{
                  t('admin.settings.registration.enableRegistration')
                }}</label>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  {{ t('admin.settings.registration.enableRegistrationHint') }}
                </p>
              </div>
              <Toggle v-model="form.registration_enabled" />
            </div>

            <!-- Email Verification -->
            <div
              class="flex items-center justify-between border-t border-gray-100 pt-4 dark:border-dark-700"
            >
              <div>
                <label class="font-medium text-gray-900 dark:text-white">{{
                  t('admin.settings.registration.emailVerification')
                }}</label>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  {{ t('admin.settings.registration.emailVerificationHint') }}
                </p>
              </div>
              <Toggle v-model="form.email_verify_enabled" />
            </div>
          </div>
        </div>

        <!-- Cloudflare Turnstile Settings -->
        <div class="card">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('admin.settings.turnstile.title') }}
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {{ t('admin.settings.turnstile.description') }}
            </p>
          </div>
          <div class="space-y-5 p-6">
            <!-- Enable Turnstile -->
            <div class="flex items-center justify-between">
              <div>
                <label class="font-medium text-gray-900 dark:text-white">{{
                  t('admin.settings.turnstile.enableTurnstile')
                }}</label>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  {{ t('admin.settings.turnstile.enableTurnstileHint') }}
                </p>
              </div>
              <Toggle v-model="form.turnstile_enabled" />
            </div>

            <!-- Turnstile Keys - Only show when enabled -->
            <div
              v-if="form.turnstile_enabled"
              class="border-t border-gray-100 pt-4 dark:border-dark-700"
            >
              <div class="grid grid-cols-1 gap-6">
                <div>
                  <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {{ t('admin.settings.turnstile.siteKey') }}
                  </label>
                  <input
                    v-model="form.turnstile_site_key"
                    type="text"
                    class="input font-mono text-sm"
                    placeholder="0x4AAAAAAA..."
                  />
                  <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                    {{ t('admin.settings.turnstile.siteKeyHint') }}
                    <a
                      href="https://dash.cloudflare.com/"
                      target="_blank"
                      class="text-primary-600 hover:text-primary-500"
                      >Cloudflare Dashboard</a
                    >
                  </p>
                </div>
                <div>
                  <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {{ t('admin.settings.turnstile.secretKey') }}
                  </label>
                  <input
                    v-model="form.turnstile_secret_key"
                    type="password"
                    class="input font-mono text-sm"
                    placeholder="0x4AAAAAAA..."
                  />
                  <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                    {{ t('admin.settings.turnstile.secretKeyHint') }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Default Settings -->
        <div class="card">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('admin.settings.defaults.title') }}
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {{ t('admin.settings.defaults.description') }}
            </p>
          </div>
          <div class="p-6">
            <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
              <div>
                <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('admin.settings.defaults.defaultBalance') }}
                </label>
                <input
                  v-model.number="form.default_balance"
                  type="number"
                  step="0.01"
                  min="0"
                  class="input"
                  placeholder="0.00"
                />
                <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.settings.defaults.defaultBalanceHint') }}
                </p>
              </div>
              <div>
                <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('admin.settings.defaults.defaultConcurrency') }}
                </label>
                <input
                  v-model.number="form.default_concurrency"
                  type="number"
                  min="1"
                  class="input"
                  placeholder="1"
                />
                <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.settings.defaults.defaultConcurrencyHint') }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Site Settings -->
        <div class="card">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('admin.settings.site.title') }}
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {{ t('admin.settings.site.description') }}
            </p>
          </div>
          <div class="space-y-6 p-6">
            <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
              <div>
                <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('admin.settings.site.siteName') }}
                </label>
                <input
                  v-model="form.site_name"
                  type="text"
                  class="input"
                  :placeholder="t('admin.settings.site.siteNamePlaceholder')"
                />
                <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.settings.site.siteNameHint') }}
                </p>
              </div>
              <div>
                <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('admin.settings.site.siteSubtitle') }}
                </label>
                <input
                  v-model="form.site_subtitle"
                  type="text"
                  class="input"
                  :placeholder="t('admin.settings.site.siteSubtitlePlaceholder')"
                />
                <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.settings.site.siteSubtitleHint') }}
                </p>
              </div>
            </div>

            <!-- API Base URL -->
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('admin.settings.site.apiBaseUrl') }}
              </label>
              <input
                v-model="form.api_base_url"
                type="text"
                class="input font-mono text-sm"
                :placeholder="t('admin.settings.site.apiBaseUrlPlaceholder')"
              />
              <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                {{ t('admin.settings.site.apiBaseUrlHint') }}
              </p>
            </div>

            <!-- Contact Info -->
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('admin.settings.site.contactInfo') }}
              </label>
              <input
                v-model="form.contact_info"
                type="text"
                class="input"
                :placeholder="t('admin.settings.site.contactInfoPlaceholder')"
              />
              <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                {{ t('admin.settings.site.contactInfoHint') }}
              </p>
            </div>

            <!-- Doc URL -->
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('admin.settings.site.docUrl') }}
              </label>
              <input
                v-model="form.doc_url"
                type="url"
                class="input font-mono text-sm"
                :placeholder="t('admin.settings.site.docUrlPlaceholder')"
              />
              <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                {{ t('admin.settings.site.docUrlHint') }}
              </p>
            </div>

            <!-- Site Logo Upload -->
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('admin.settings.site.siteLogo') }}
              </label>
              <div class="flex items-start gap-6">
                <!-- Logo Preview -->
                <div class="flex-shrink-0">
                  <div
                    class="flex h-20 w-20 items-center justify-center overflow-hidden rounded-xl border-2 border-dashed border-gray-300 bg-gray-50 dark:border-dark-600 dark:bg-dark-800"
                    :class="{ 'border-solid': form.site_logo }"
                  >
                    <img
                      v-if="form.site_logo"
                      :src="form.site_logo"
                      alt="Site Logo"
                      class="h-full w-full object-contain"
                    />
                    <svg
                      v-else
                      class="h-8 w-8 text-gray-400 dark:text-dark-500"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="1.5"
                        d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                      />
                    </svg>
                  </div>
                </div>
                <!-- Upload Controls -->
                <div class="flex-1 space-y-3">
                  <div class="flex items-center gap-3">
                    <label class="btn btn-secondary btn-sm cursor-pointer">
                      <input
                        type="file"
                        accept="image/*"
                        class="hidden"
                        @change="handleLogoUpload"
                      />
                      <svg
                        class="mr-1.5 h-4 w-4"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                        />
                      </svg>
                      {{ t('admin.settings.site.uploadImage') }}
                    </label>
                    <button
                      v-if="form.site_logo"
                      type="button"
                      @click="form.site_logo = ''"
                      class="btn btn-secondary btn-sm text-red-600 hover:text-red-700 dark:text-red-400"
                    >
                      <svg
                        class="mr-1.5 h-4 w-4"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        />
                      </svg>
                      {{ t('admin.settings.site.remove') }}
                    </button>
                  </div>
                  <p class="text-xs text-gray-500 dark:text-gray-400">
                    {{ t('admin.settings.site.logoHint') }}
                  </p>
                  <p v-if="logoError" class="text-xs text-red-500">{{ logoError }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- SMTP Settings - Only show when email verification is enabled -->
        <div v-if="form.email_verify_enabled" class="card">
          <div
            class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700"
          >
            <div>
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('admin.settings.smtp.title') }}
              </h2>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {{ t('admin.settings.smtp.description') }}
              </p>
            </div>
            <button
              type="button"
              @click="testSmtpConnection"
              :disabled="testingSmtp"
              class="btn btn-secondary btn-sm"
            >
              <svg v-if="testingSmtp" class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              {{
                testingSmtp
                  ? t('admin.settings.smtp.testing')
                  : t('admin.settings.smtp.testConnection')
              }}
            </button>
          </div>
          <div class="space-y-6 p-6">
            <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
              <div>
                <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('admin.settings.smtp.host') }}
                </label>
                <input
                  v-model="form.smtp_host"
                  type="text"
                  class="input"
                  :placeholder="t('admin.settings.smtp.hostPlaceholder')"
                />
              </div>
              <div>
                <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('admin.settings.smtp.port') }}
                </label>
                <input
                  v-model.number="form.smtp_port"
                  type="number"
                  min="1"
                  max="65535"
                  class="input"
                  :placeholder="t('admin.settings.smtp.portPlaceholder')"
                />
              </div>
              <div>
                <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('admin.settings.smtp.username') }}
                </label>
                <input
                  v-model="form.smtp_username"
                  type="text"
                  class="input"
                  :placeholder="t('admin.settings.smtp.usernamePlaceholder')"
                />
              </div>
              <div>
                <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('admin.settings.smtp.password') }}
                </label>
                <input
                  v-model="form.smtp_password"
                  type="password"
                  class="input"
                  :placeholder="t('admin.settings.smtp.passwordPlaceholder')"
                />
                <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.settings.smtp.passwordHint') }}
                </p>
              </div>
              <div>
                <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('admin.settings.smtp.fromEmail') }}
                </label>
                <input
                  v-model="form.smtp_from_email"
                  type="email"
                  class="input"
                  :placeholder="t('admin.settings.smtp.fromEmailPlaceholder')"
                />
              </div>
              <div>
                <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('admin.settings.smtp.fromName') }}
                </label>
                <input
                  v-model="form.smtp_from_name"
                  type="text"
                  class="input"
                  :placeholder="t('admin.settings.smtp.fromNamePlaceholder')"
                />
              </div>
            </div>

            <!-- Use TLS Toggle -->
            <div
              class="flex items-center justify-between border-t border-gray-100 pt-4 dark:border-dark-700"
            >
              <div>
                <label class="font-medium text-gray-900 dark:text-white">{{
                  t('admin.settings.smtp.useTls')
                }}</label>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  {{ t('admin.settings.smtp.useTlsHint') }}
                </p>
              </div>
              <Toggle v-model="form.smtp_use_tls" />
            </div>
          </div>
        </div>

        <!-- Send Test Email - Only show when email verification is enabled -->
        <div v-if="form.email_verify_enabled" class="card">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('admin.settings.testEmail.title') }}
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {{ t('admin.settings.testEmail.description') }}
            </p>
          </div>
          <div class="p-6">
            <div class="flex items-end gap-4">
              <div class="flex-1">
                <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('admin.settings.testEmail.recipientEmail') }}
                </label>
                <input
                  v-model="testEmailAddress"
                  type="email"
                  class="input"
                  :placeholder="t('admin.settings.testEmail.recipientEmailPlaceholder')"
                />
              </div>
              <button
                type="button"
                @click="sendTestEmail"
                :disabled="sendingTestEmail || !testEmailAddress"
                class="btn btn-secondary"
              >
                <svg
                  v-if="sendingTestEmail"
                  class="h-4 w-4 animate-spin"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                  ></circle>
                  <path
                    class="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                {{
                  sendingTestEmail
                    ? t('admin.settings.testEmail.sending')
                    : t('admin.settings.testEmail.sendTestEmail')
                }}
              </button>
            </div>
          </div>
        </div>

        <!-- SSO Settings Card -->
        <div class="card">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              SSO登录设置
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              配置OpenID Connect (OIDC)单点登录，支持自动创建用户和邮箱域名限制
            </p>
          </div>
          <div class="p-6">
            <!-- SSO Settings List -->
            <div class="space-y-4">
              <!-- 启用SSO登录 -->
              <div class="flex items-center gap-4 border-b border-gray-100 pb-4 dark:border-dark-700">
                <div class="w-48 flex-shrink-0">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">启用SSO登录</div>
                </div>
                <div class="flex-shrink-0">
                  <Toggle
                    :model-value="form.sso_enabled"
                    @update:model-value="handleSSOSettingChange('sso_enabled', $event)"
                  />
                  <span class="ml-2 text-sm text-primary-600 dark:text-primary-400">
                    {{ form.sso_enabled ? '开启' : '关闭' }}
                  </span>
                </div>
                <div class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                  允许用户使用SSO (OpenID Connect)登录
                </div>
              </div>

              <!-- 启用密码登录 -->
              <div class="flex items-center gap-4 border-b border-gray-100 pb-4 dark:border-dark-700">
                <div class="w-48 flex-shrink-0">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">
                    启用密码登录
                    <span class="ml-1 text-xs text-red-500">* 必填</span>
                  </div>
                </div>
                <div class="flex-shrink-0">
                  <Toggle
                    :model-value="form.password_login_enabled"
                    @update:model-value="handleSSOSettingChange('password_login_enabled', $event)"
                  />
                  <span class="ml-2 text-sm text-primary-600 dark:text-primary-400">
                    {{ form.password_login_enabled ? '开启' : '关闭' }}
                  </span>
                </div>
                <div class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                  允许用户使用密码登录（至少保留一种登录方式）
                </div>
              </div>

              <!-- Issuer URL -->
              <div class="flex items-center gap-4 border-b border-gray-100 pb-4 dark:border-dark-700">
                <div class="w-48 flex-shrink-0">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">
                    Issuer URL
                    <span class="ml-1 text-xs text-red-500">* 必填</span>
                  </div>
                </div>
                <div class="flex-1">
                  <a
                    v-if="form.sso_issuer_url"
                    :href="form.sso_issuer_url"
                    target="_blank"
                    class="text-sm text-primary-600 hover:underline dark:text-primary-400"
                  >
                    {{ form.sso_issuer_url }}
                    <svg class="ml-1 inline h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                    </svg>
                  </a>
                  <span v-else class="text-sm text-gray-400 dark:text-gray-500">未配置</span>
                </div>
                <div class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                  OIDC Provider的Issuer URL（例如：https://accounts.google.com）
                </div>
                <button
                  @click="handleEditSSOSetting('sso_issuer_url')"
                  class="flex-shrink-0 text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
                >
                  编辑
                </button>
              </div>

              <!-- Client ID -->
              <div class="flex items-center gap-4 border-b border-gray-100 pb-4 dark:border-dark-700">
                <div class="w-48 flex-shrink-0">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">
                    Client ID
                    <span class="ml-1 text-xs text-red-500">* 必填</span>
                  </div>
                </div>
                <div class="flex-1">
                  <span v-if="form.sso_client_id" class="text-sm text-gray-700 dark:text-gray-300">
                    {{ form.sso_client_id }}
                  </span>
                  <span v-else class="text-sm text-gray-400 dark:text-gray-500">未配置</span>
                </div>
                <div class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                  OIDC应用的Client ID
                </div>
                <button
                  @click="handleEditSSOSetting('sso_client_id')"
                  class="flex-shrink-0 text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
                >
                  编辑
                </button>
              </div>

              <!-- Client Secret -->
              <div class="flex items-center gap-4 border-b border-gray-100 pb-4 dark:border-dark-700">
                <div class="w-48 flex-shrink-0">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">
                    Client Secret
                    <span class="ml-1 text-xs text-red-500">* 必填</span>
                  </div>
                </div>
                <div class="flex-1">
                  <span v-if="form.sso_client_secret" class="text-sm text-gray-500 dark:text-gray-400">
                    ••••••••（已加密）
                  </span>
                  <span v-else class="text-sm text-gray-400 dark:text-gray-500">未配置</span>
                </div>
                <div class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                  OIDC应用的Client Secret（已加密存储）
                </div>
                <button
                  @click="handleEditSSOSetting('sso_client_secret')"
                  class="flex-shrink-0 text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
                >
                  编辑
                </button>
              </div>

              <!-- Redirect URI -->
              <div class="flex items-center gap-4 border-b border-gray-100 pb-4 dark:border-dark-700">
                <div class="w-48 flex-shrink-0">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">
                    Redirect URI
                    <span class="ml-1 text-xs text-red-500">* 必填</span>
                  </div>
                </div>
                <div class="flex-1">
                  <a
                    v-if="form.sso_redirect_uri"
                    :href="form.sso_redirect_uri"
                    target="_blank"
                    class="text-sm text-primary-600 hover:underline dark:text-primary-400"
                  >
                    {{ form.sso_redirect_uri }}
                    <svg class="ml-1 inline h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                    </svg>
                  </a>
                  <span v-else class="text-sm text-gray-400 dark:text-gray-500">未配置</span>
                </div>
                <div class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                  SSO回调地址（例如：http://localhost:8080/auth/sso/callback）
                </div>
                <button
                  @click="handleEditSSOSetting('sso_redirect_uri')"
                  class="flex-shrink-0 text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
                >
                  编辑
                </button>
              </div>

              <!-- 允许的邮箱域名 -->
              <div class="flex items-center gap-4 border-b border-gray-100 pb-4 dark:border-dark-700">
                <div class="w-48 flex-shrink-0">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">允许的邮箱域名</div>
                </div>
                <div class="flex-1">
                  <span
                    v-if="form.sso_allowed_domains && form.sso_allowed_domains.length > 0"
                    class="text-sm text-gray-700 dark:text-gray-300"
                  >
                    {{ form.sso_allowed_domains.join(', ') }}
                  </span>
                  <span v-else class="text-sm text-gray-500 dark:text-gray-400">无限制</span>
                </div>
                <div class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                  限制允许SSO登录的邮箱域名（留空则不限制）
                </div>
                <button
                  @click="handleEditSSOSetting('sso_allowed_domains')"
                  class="flex-shrink-0 text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
                >
                  编辑
                </button>
              </div>

              <!-- 自动创建用户 -->
              <div class="flex items-center gap-4 border-b border-gray-100 pb-4 dark:border-dark-700">
                <div class="w-48 flex-shrink-0">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">自动创建用户</div>
                </div>
                <div class="flex-shrink-0">
                  <Toggle
                    :model-value="form.sso_auto_create_user"
                    @update:model-value="handleSSOSettingChange('sso_auto_create_user', $event)"
                  />
                  <span class="ml-2 text-sm text-primary-600 dark:text-primary-400">
                    {{ form.sso_auto_create_user ? '开启' : '关闭' }}
                  </span>
                </div>
                <div class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                  首次SSO登录时自动创建用户账号
                </div>
              </div>

              <!-- 最小信任等级 -->
              <div class="flex items-center gap-4 pb-4">
                <div class="w-48 flex-shrink-0">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">最小信任等级</div>
                </div>
                <div class="flex-1">
                  <span class="text-sm text-gray-700 dark:text-gray-300">
                    {{ form.sso_min_trust_level }}
                  </span>
                </div>
                <div class="flex-1 text-sm text-gray-500 dark:text-gray-400">
                  Discourse论坛最小信任等级要求（0-4，0表示无限制）
                </div>
                <button
                  @click="handleEditSSOSetting('sso_min_trust_level')"
                  class="flex-shrink-0 text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
                >
                  编辑
                </button>
              </div>
            </div>

            <!-- Test SSO Connection Button -->
            <div class="mt-6 flex justify-end">
              <button
                type="button"
                @click="testSSOConnection"
                :disabled="testingSso || !form.sso_issuer_url"
                class="btn btn-secondary"
              >
                <svg
                  v-if="testingSso"
                  class="mr-2 h-4 w-4 animate-spin"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                  ></circle>
                  <path
                    class="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                {{ testingSso ? '测试中...' : '测试SSO连接' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Epay Settings Card -->
        <div class="card">
          <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('admin.settings.epay.title') }}
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {{ t('admin.settings.epay.description') }}
            </p>
          </div>
          <div class="space-y-5 p-6">
            <!-- Enable Epay -->
            <div class="flex items-center justify-between">
              <div>
                <label class="font-medium text-gray-900 dark:text-white">{{
                  t('admin.settings.epay.enableEpay')
                }}</label>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  {{ t('admin.settings.epay.enableEpayHint') }}
                </p>
              </div>
              <Toggle v-model="form.epay_enabled" />
            </div>

            <!-- Epay Configuration - Only show when enabled -->
            <div
              v-if="form.epay_enabled"
              class="border-t border-gray-100 pt-4 dark:border-dark-700"
            >
              <div class="grid grid-cols-1 gap-6">
                <div>
                  <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {{ t('admin.settings.epay.apiUrl') }}
                  </label>
                  <input
                    v-model="form.epay_api_url"
                    type="text"
                    class="input font-mono text-sm"
                    placeholder="https://pay.example.com"
                  />
                  <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                    {{ t('admin.settings.epay.apiUrlHint') }}
                  </p>
                </div>
                <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
                  <div>
                    <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                      {{ t('admin.settings.epay.merchantId') }}
                    </label>
                    <input
                      v-model="form.epay_merchant_id"
                      type="text"
                      class="input font-mono text-sm"
                      placeholder="1000"
                    />
                    <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                      {{ t('admin.settings.epay.merchantIdHint') }}
                    </p>
                  </div>
                  <div>
                    <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                      {{ t('admin.settings.epay.merchantKey') }}
                    </label>
                    <input
                      v-model="form.epay_merchant_key"
                      type="password"
                      class="input font-mono text-sm"
                      placeholder="********"
                    />
                    <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                      {{ t('admin.settings.epay.merchantKeyHint') }}
                    </p>
                  </div>
                </div>
                <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
                  <div>
                    <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                      {{ t('admin.settings.epay.notifyUrl') }}
                    </label>
                    <input
                      v-model="form.epay_notify_url"
                      type="text"
                      class="input font-mono text-sm"
                      readonly
                    />
                    <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                      {{ t('admin.settings.epay.notifyUrlHint') }}
                    </p>
                  </div>
                  <div>
                    <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                      {{ t('admin.settings.epay.returnUrl') }}
                    </label>
                    <input
                      v-model="form.epay_return_url"
                      type="text"
                      class="input font-mono text-sm"
                      readonly
                    />
                    <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">
                      {{ t('admin.settings.epay.returnUrlHint') }}
                    </p>
                  </div>
                </div>

                <!-- Payment Channels Configuration -->
                <div class="border-t border-gray-100 pt-4 dark:border-dark-700">
                  <div class="mb-3 flex items-center justify-between">
                    <div>
                      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                        {{ t('admin.settings.epay.paymentChannels') }}
                      </label>
                      <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                        {{ t('admin.settings.epay.paymentChannelsHint') }}
                      </p>
                    </div>
                    <button
                      type="button"
                      @click="showAddChannelDialog = true"
                      class="btn btn-secondary btn-sm"
                    >
                      <svg class="mr-1.5 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                      </svg>
                      {{ t('admin.settings.epay.addChannel') }}
                    </button>
                  </div>
                  <div class="overflow-hidden rounded-lg border border-gray-200 dark:border-dark-600">
                    <table class="min-w-full divide-y divide-gray-200 dark:divide-dark-600">
                      <thead class="bg-gray-50 dark:bg-dark-700">
                        <tr>
                          <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">
                            {{ t('admin.settings.epay.channelEnabled') }}
                          </th>
                          <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">
                            {{ t('admin.settings.epay.channelName') }}
                          </th>
                          <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">
                            {{ t('admin.settings.epay.channelDisplayName') }}
                          </th>
                          <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">
                            {{ t('admin.settings.epay.channelEpayType') }}
                          </th>
                          <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">
                            {{ t('common.actions') }}
                          </th>
                        </tr>
                      </thead>
                      <tbody class="divide-y divide-gray-200 bg-white dark:divide-dark-600 dark:bg-dark-800">
                        <tr v-for="(channel, index) in form.payment_channels" :key="channel.key">
                          <td class="px-4 py-3">
                            <Toggle v-model="form.payment_channels[index].enabled" />
                          </td>
                          <td class="px-4 py-3">
                            <span class="font-mono text-sm text-gray-900 dark:text-gray-100">{{ channel.key }}</span>
                          </td>
                          <td class="px-4 py-3">
                            <input
                              v-model="form.payment_channels[index].display_name"
                              type="text"
                              class="input w-32 text-sm"
                              :placeholder="channel.key"
                            />
                          </td>
                          <td class="px-4 py-3">
                            <input
                              v-model="form.payment_channels[index].epay_type"
                              type="text"
                              class="input w-24 font-mono text-sm"
                              placeholder="epay"
                            />
                          </td>
                          <td class="px-4 py-3">
                            <button
                              type="button"
                              @click="deletePaymentChannel(index)"
                              class="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                              :title="t('common.delete')"
                            >
                              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                              </svg>
                            </button>
                          </td>
                        </tr>
                        <tr v-if="form.payment_channels.length === 0">
                          <td colspan="5" class="px-4 py-8 text-center text-gray-500 dark:text-gray-400">
                            {{ t('admin.settings.epay.noChannels') }}
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                  <p class="mt-2 text-xs text-amber-600 dark:text-amber-400">
                    {{ t('admin.settings.epay.epayTypeNote') }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- SSO Edit Dialog -->
        <SettingEditDialog
          :show="showSsoEditDialog"
          :title="editingSettingLabel"
          :label="editingSettingLabel"
          :description="editingSettingDescription"
          :input-type="editingSettingType"
          :value="editingSettingValue"
          :placeholder="editingSettingPlaceholder"
          @close="showSsoEditDialog = false"
          @save="saveSSOSetting"
        />

        <!-- Add Payment Channel Dialog -->
        <Teleport to="body">
          <div
            v-if="showAddChannelDialog"
            class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
            @click.self="showAddChannelDialog = false"
          >
            <div class="w-full max-w-md rounded-lg bg-white p-6 shadow-xl dark:bg-dark-800">
              <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('admin.settings.epay.addChannel') }}
              </h3>
              <div class="space-y-4">
                <div>
                  <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {{ t('admin.settings.epay.channelKey') }}
                    <span class="ml-1 text-xs text-red-500">*</span>
                  </label>
                  <input
                    v-model="newChannel.key"
                    type="text"
                    class="input font-mono text-sm"
                    :placeholder="t('admin.settings.epay.channelKeyPlaceholder')"
                  />
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    {{ t('admin.settings.epay.channelKeyHint') }}
                  </p>
                </div>
                <div>
                  <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {{ t('admin.settings.epay.channelDisplayName') }}
                    <span class="ml-1 text-xs text-red-500">*</span>
                  </label>
                  <input
                    v-model="newChannel.display_name"
                    type="text"
                    class="input text-sm"
                    :placeholder="t('admin.settings.epay.channelDisplayNamePlaceholder')"
                  />
                </div>
                <div>
                  <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {{ t('admin.settings.epay.channelEpayType') }}
                  </label>
                  <input
                    v-model="newChannel.epay_type"
                    type="text"
                    class="input font-mono text-sm"
                    placeholder="epay"
                  />
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    {{ t('admin.settings.epay.epayTypeHint') }}
                  </p>
                </div>
                <div>
                  <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
                    {{ t('admin.settings.epay.channelIcon') }}
                  </label>
                  <input
                    v-model="newChannel.icon"
                    type="text"
                    class="input font-mono text-sm"
                    placeholder="credit-card"
                  />
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    {{ t('admin.settings.epay.channelIconHint') }}
                  </p>
                </div>
              </div>
              <div class="mt-6 flex justify-end gap-3">
                <button
                  type="button"
                  @click="showAddChannelDialog = false"
                  class="btn btn-secondary"
                >
                  {{ t('common.cancel') }}
                </button>
                <button
                  type="button"
                  @click="addPaymentChannel"
                  :disabled="!newChannel.key || !newChannel.display_name"
                  class="btn btn-primary"
                >
                  {{ t('common.add') }}
                </button>
              </div>
            </div>
          </div>
        </Teleport>

        <!-- Save Button -->
        <div class="flex justify-end">
          <button type="submit" :disabled="saving" class="btn btn-primary">
            <svg v-if="saving" class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            {{ saving ? t('admin.settings.saving') : t('admin.settings.saveSettings') }}
          </button>
        </div>
      </form>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '@/api'
import type { SystemSettings, PaymentChannel } from '@/api/admin/settings'
import AppLayout from '@/components/layout/AppLayout.vue'
import Toggle from '@/components/common/Toggle.vue'
import SettingEditDialog from '@/components/admin/SettingEditDialog.vue'
import { useAppStore } from '@/stores'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const saving = ref(false)
const testingSmtp = ref(false)
const sendingTestEmail = ref(false)
const testEmailAddress = ref('')
const logoError = ref('')

// Admin API Key 状态
const adminApiKeyLoading = ref(true)
const adminApiKeyExists = ref(false)
const adminApiKeyMasked = ref('')
const adminApiKeyOperating = ref(false)
const newAdminApiKey = ref('')

// SSO Settings 状态
const testingSso = ref(false)
const showSsoEditDialog = ref(false)
const editingSettingKey = ref('')
const editingSettingLabel = ref('')
const editingSettingDescription = ref('')
const editingSettingType = ref<'text' | 'password' | 'number' | 'url' | 'array' | 'textarea'>('text')
const editingSettingValue = ref<any>('')
const editingSettingPlaceholder = ref('')

// Payment Channel 状态
const showAddChannelDialog = ref(false)
const newChannel = reactive({
  key: '',
  display_name: '',
  epay_type: 'epay',
  icon: 'credit-card'
})

const form = reactive<SystemSettings>({
  registration_enabled: true,
  email_verify_enabled: false,
  default_balance: 0,
  default_concurrency: 1,
  site_name: 'Sub2API',
  site_logo: '',
  site_subtitle: 'Subscription to API Conversion Platform',
  api_base_url: '',
  contact_info: '',
  doc_url: '',
  smtp_host: '',
  smtp_port: 587,
  smtp_username: '',
  smtp_password: '',
  smtp_from_email: '',
  smtp_from_name: '',
  smtp_use_tls: true,
  // Cloudflare Turnstile
  turnstile_enabled: false,
  turnstile_site_key: '',
  turnstile_secret_key: '',
  // SSO Settings
  sso_enabled: false,
  password_login_enabled: true,
  sso_issuer_url: '',
  sso_client_id: '',
  sso_client_secret: '',
  sso_redirect_uri: '',
  sso_allowed_domains: [],
  sso_auto_create_user: true,
  sso_min_trust_level: 0,
  // Epay Settings
  epay_enabled: false,
  epay_api_url: '',
  epay_merchant_id: '',
  epay_merchant_key: '',
  epay_notify_url: '',
  epay_return_url: '',
  // Payment Channels
  payment_channels: [
    { key: 'alipay', display_name: '支付宝', epay_type: 'epay', icon: 'alipay', enabled: true, sort_order: 1 },
    { key: 'wxpay', display_name: '微信支付', epay_type: 'epay', icon: 'wechat', enabled: true, sort_order: 2 },
    { key: 'epusdt', display_name: 'USDT', epay_type: 'epay', icon: 'usdt', enabled: true, sort_order: 3 }
  ] as PaymentChannel[]
})

function handleLogoUpload(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  logoError.value = ''

  if (!file) return

  // Check file size (300KB = 307200 bytes)
  const maxSize = 300 * 1024
  if (file.size > maxSize) {
    logoError.value = t('admin.settings.site.logoSizeError', {
      size: (file.size / 1024).toFixed(1)
    })
    input.value = ''
    return
  }

  // Check file type
  if (!file.type.startsWith('image/')) {
    logoError.value = t('admin.settings.site.logoTypeError')
    input.value = ''
    return
  }

  // Convert to base64
  const reader = new FileReader()
  reader.onload = (e) => {
    form.site_logo = e.target?.result as string
  }
  reader.onerror = () => {
    logoError.value = t('admin.settings.site.logoReadError')
  }
  reader.readAsDataURL(file)

  // Reset input
  input.value = ''
}

async function loadSettings() {
  loading.value = true
  try {
    const settings = await adminAPI.settings.getSettings()
    Object.assign(form, settings)
  } catch (error: any) {
    appStore.showError(
      t('admin.settings.failedToLoad') + ': ' + (error.message || t('common.unknownError'))
    )
  } finally {
    loading.value = false
  }
}

async function saveSettings() {
  saving.value = true
  try {
    await adminAPI.settings.updateSettings(form)
    // Refresh cached public settings so sidebar/header update immediately
    await appStore.fetchPublicSettings(true)
    appStore.showSuccess(t('admin.settings.settingsSaved'))
  } catch (error: any) {
    appStore.showError(
      t('admin.settings.failedToSave') + ': ' + (error.message || t('common.unknownError'))
    )
  } finally {
    saving.value = false
  }
}

async function testSmtpConnection() {
  testingSmtp.value = true
  try {
    const result = await adminAPI.settings.testSmtpConnection({
      smtp_host: form.smtp_host,
      smtp_port: form.smtp_port,
      smtp_username: form.smtp_username,
      smtp_password: form.smtp_password,
      smtp_use_tls: form.smtp_use_tls
    })
    // API returns { message: "..." } on success, errors are thrown as exceptions
    appStore.showSuccess(result.message || t('admin.settings.smtpConnectionSuccess'))
  } catch (error: any) {
    appStore.showError(
      t('admin.settings.failedToTestSmtp') + ': ' + (error.message || t('common.unknownError'))
    )
  } finally {
    testingSmtp.value = false
  }
}

async function sendTestEmail() {
  if (!testEmailAddress.value) {
    appStore.showError(t('admin.settings.testEmail.enterRecipientHint'))
    return
  }

  sendingTestEmail.value = true
  try {
    const result = await adminAPI.settings.sendTestEmail({
      email: testEmailAddress.value,
      smtp_host: form.smtp_host,
      smtp_port: form.smtp_port,
      smtp_username: form.smtp_username,
      smtp_password: form.smtp_password,
      smtp_from_email: form.smtp_from_email,
      smtp_from_name: form.smtp_from_name,
      smtp_use_tls: form.smtp_use_tls
    })
    // API returns { message: "..." } on success, errors are thrown as exceptions
    appStore.showSuccess(result.message || t('admin.settings.testEmailSent'))
  } catch (error: any) {
    appStore.showError(
      t('admin.settings.failedToSendTestEmail') + ': ' + (error.message || t('common.unknownError'))
    )
  } finally {
    sendingTestEmail.value = false
  }
}

// Admin API Key 方法
async function loadAdminApiKey() {
  adminApiKeyLoading.value = true
  try {
    const status = await adminAPI.settings.getAdminApiKey()
    adminApiKeyExists.value = status.exists
    adminApiKeyMasked.value = status.masked_key
  } catch (error: any) {
    console.error('Failed to load admin API key status:', error)
  } finally {
    adminApiKeyLoading.value = false
  }
}

async function createAdminApiKey() {
  adminApiKeyOperating.value = true
  try {
    const result = await adminAPI.settings.regenerateAdminApiKey()
    newAdminApiKey.value = result.key
    adminApiKeyExists.value = true
    adminApiKeyMasked.value = result.key.substring(0, 10) + '...' + result.key.slice(-4)
    appStore.showSuccess(t('admin.settings.adminApiKey.keyGenerated'))
  } catch (error: any) {
    appStore.showError(error.message || t('common.error'))
  } finally {
    adminApiKeyOperating.value = false
  }
}

async function regenerateAdminApiKey() {
  if (!confirm(t('admin.settings.adminApiKey.regenerateConfirm'))) return
  await createAdminApiKey()
}

async function deleteAdminApiKey() {
  if (!confirm(t('admin.settings.adminApiKey.deleteConfirm'))) return
  adminApiKeyOperating.value = true
  try {
    await adminAPI.settings.deleteAdminApiKey()
    adminApiKeyExists.value = false
    adminApiKeyMasked.value = ''
    newAdminApiKey.value = ''
    appStore.showSuccess(t('admin.settings.adminApiKey.keyDeleted'))
  } catch (error: any) {
    appStore.showError(error.message || t('common.error'))
  } finally {
    adminApiKeyOperating.value = false
  }
}

function copyNewKey() {
  navigator.clipboard
    .writeText(newAdminApiKey.value)
    .then(() => {
      appStore.showSuccess(t('admin.settings.adminApiKey.keyCopied'))
    })
    .catch(() => {
      appStore.showError(t('common.copyFailed'))
    })
}

// SSO Settings Methods
async function handleSSOSettingChange(key: string, value: any) {
  try {
    // 验证逻辑：至少保留一种登录方式
    if (key === 'sso_enabled' || key === 'password_login_enabled') {
      const newSSO = key === 'sso_enabled' ? value : form.sso_enabled
      const newPwd = key === 'password_login_enabled' ? value : form.password_login_enabled

      if (!newSSO && !newPwd) {
        appStore.showError('至少需要启用一种登录方式（SSO或密码登录）')
        return
      }
    }

    // 调用单项更新API
    await adminAPI.settings.updateSingleSetting(key, value)

    // 更新本地表单值
    ;(form as any)[key] = value

    // 如果更新的是公开配置，刷新缓存
    await appStore.fetchPublicSettings(true)

    appStore.showSuccess('设置已保存')
  } catch (error: any) {
    appStore.showError('保存失败: ' + (error.message || t('common.unknownError')))
    // 恢复旧值
    await loadSettings()
  }
}

function handleEditSSOSetting(key: string) {
  // 配置项元数据映射
  const settingsMeta: Record<string, { label: string; description: string; type: 'text' | 'password' | 'number' | 'url' | 'array' | 'textarea'; placeholder: string }> = {
    sso_issuer_url: {
      label: 'Issuer URL',
      description: 'OIDC Provider的Issuer URL',
      type: 'url',
      placeholder: 'https://accounts.google.com'
    },
    sso_client_id: {
      label: 'Client ID',
      description: 'OIDC应用的Client ID',
      type: 'text',
      placeholder: 'your-client-id.apps.googleusercontent.com'
    },
    sso_client_secret: {
      label: 'Client Secret',
      description: 'OIDC应用的Client Secret（已加密存储）',
      type: 'password',
      placeholder: '留空则不修改'
    },
    sso_redirect_uri: {
      label: 'Redirect URI',
      description: 'SSO回调地址',
      type: 'url',
      placeholder: 'http://localhost:8080/auth/sso/callback'
    },
    sso_allowed_domains: {
      label: '允许的邮箱域名',
      description: '限制允许SSO登录的邮箱域名（留空则不限制）',
      type: 'array',
      placeholder: 'example.com, company.com'
    },
    sso_min_trust_level: {
      label: '最小信任等级',
      description: 'Discourse论坛最小信任等级要求（0-4）',
      type: 'number',
      placeholder: '0'
    }
  }

  const meta = settingsMeta[key]
  if (!meta) return

  editingSettingKey.value = key
  editingSettingLabel.value = meta.label
  editingSettingDescription.value = meta.description
  editingSettingType.value = meta.type
  editingSettingValue.value = (form as any)[key]
  editingSettingPlaceholder.value = meta.placeholder

  showSsoEditDialog.value = true
}

async function saveSSOSetting(value: any) {
  try {
    const key = editingSettingKey.value

    // 对于密码字段，如果为空则不更新
    if (editingSettingType.value === 'password' && !value) {
      showSsoEditDialog.value = false
      return
    }

    await adminAPI.settings.updateSingleSetting(key, value)

    // 更新本地表单值
    ;(form as any)[key] = value

    // 如果更新的是公开配置，刷新缓存
    await appStore.fetchPublicSettings(true)

    showSsoEditDialog.value = false
    appStore.showSuccess('设置已保存')
  } catch (error: any) {
    appStore.showError('保存失败: ' + (error.message || t('common.unknownError')))
  }
}

async function testSSOConnection() {
  if (!form.sso_issuer_url) {
    appStore.showError('请先配置Issuer URL')
    return
  }

  testingSso.value = true
  try {
    const result = await adminAPI.settings.testSSOConnection({
      issuer_url: form.sso_issuer_url
    })
    appStore.showSuccess(`SSO配置有效！Issuer: ${result.issuer}`)
  } catch (error: any) {
    appStore.showError(
      '测试失败: ' + (error.message || t('common.unknownError'))
    )
  } finally {
    testingSso.value = false
  }
}

// Payment Channel Methods
function addPaymentChannel() {
  // 验证 key 是否唯一
  const existingKeys = form.payment_channels.map(c => c.key.toLowerCase())
  if (existingKeys.includes(newChannel.key.toLowerCase())) {
    appStore.showError(t('admin.settings.epay.channelKeyExists'))
    return
  }

  // 添加新渠道
  const maxSortOrder = form.payment_channels.reduce((max, c) => Math.max(max, c.sort_order || 0), 0)
  form.payment_channels.push({
    key: newChannel.key.toLowerCase(),
    display_name: newChannel.display_name,
    epay_type: newChannel.epay_type || 'epay',
    icon: newChannel.icon || 'credit-card',
    enabled: true,
    sort_order: maxSortOrder + 1
  })

  // 重置表单并关闭对话框
  newChannel.key = ''
  newChannel.display_name = ''
  newChannel.epay_type = 'epay'
  newChannel.icon = 'credit-card'
  showAddChannelDialog.value = false

  appStore.showSuccess(t('admin.settings.epay.channelAdded'))
}

function deletePaymentChannel(index: number) {
  const channel = form.payment_channels[index]
  if (!confirm(t('admin.settings.epay.confirmDeleteChannel', { name: channel.display_name }))) {
    return
  }

  form.payment_channels.splice(index, 1)
  appStore.showSuccess(t('admin.settings.epay.channelDeleted'))
}

onMounted(() => {
  loadSettings()
  loadAdminApiKey()

  // Initialize Epay callback URLs based on current location
  const baseUrl = window.location.origin
  form.epay_notify_url = `${baseUrl}/api/v1/payment/notify/epay`
  form.epay_return_url = `${baseUrl}/recharge/result`
})
</script>
