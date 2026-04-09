<template>
  <section class="space-y-5">
    <div :class="heroClass">
      <div class="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
        <div class="max-w-3xl">
          <p class="text-[11px] uppercase tracking-[0.26em] text-miku/80">{{ copy.overview.badge }}</p>
          <h1 class="mt-3 text-3xl font-semibold tracking-[0.02em] text-slate-900">{{ copy.page.title }}</h1>
          <p class="mt-2 max-w-2xl text-sm leading-7 text-slate-600">{{ copy.page.subtitle }}</p>
          <p class="mt-3 max-w-2xl text-sm leading-7 text-slate-500">{{ copy.overview.description }}</p>
        </div>

        <div class="grid gap-3 sm:grid-cols-2">
          <div class="rounded-[22px] border border-white/70 bg-white/78 px-4 py-4 shadow-[0_12px_28px_rgba(15,23,42,0.08)] backdrop-blur-xl">
            <p class="text-xs uppercase tracking-[0.18em] text-slate-400">{{ copy.overview.mergedLabel }}</p>
            <div class="mt-3 flex flex-wrap gap-2">
              <span
                v-for="item in copy.overview.mergedItems"
                :key="item"
                class="rounded-full border border-miku/20 bg-miku/10 px-3 py-1 text-xs font-medium text-miku"
              >
                {{ item }}
              </span>
            </div>
          </div>

          <div class="rounded-[22px] border border-white/70 bg-white/78 px-4 py-4 shadow-[0_12px_28px_rgba(15,23,42,0.08)] backdrop-blur-xl">
            <p class="text-xs uppercase tracking-[0.18em] text-slate-400">{{ copy.overview.fallbackLabel }}</p>
            <p class="mt-3 text-sm font-semibold text-slate-900">{{ copy.overview.fallbackValue }}</p>
          </div>
        </div>
      </div>
    </div>

    <div class="grid gap-5 xl:grid-cols-[240px_minmax(0,1.15fr)_340px]">
      <div class="space-y-5 xl:sticky xl:top-6 xl:self-start">
        <div :class="surfaceClass">
          <p class="text-[11px] uppercase tracking-[0.24em] text-slate-400">{{ copy.sections.navTitle }}</p>
          <div class="mt-4 space-y-2">
            <button
              v-for="section in sections"
              :key="section.key"
              type="button"
              class="w-full rounded-[20px] border px-4 py-3 text-left transition duration-200"
              :class="section.key === activeSection
                ? 'border-miku/35 bg-[linear-gradient(135deg,rgba(57,197,187,0.12),rgba(255,255,255,0.95))] shadow-[0_14px_24px_rgba(57,197,187,0.12)]'
                : 'border-white/70 bg-white/68 hover:border-slate-200 hover:bg-white/85'"
              @click="activeSection = section.key"
            >
              <div class="flex items-start justify-between gap-3">
                <div>
                  <p class="text-[11px] uppercase tracking-[0.24em] text-slate-400">{{ section.eyebrow }}</p>
                  <p class="mt-2 text-sm font-semibold text-slate-900">{{ section.title }}</p>
                </div>
                <span
                  v-if="section.dirty"
                  class="rounded-full border border-amber-200 bg-amber-50 px-2 py-1 text-[11px] font-medium text-amber-600"
                >
                  {{ copy.sections.pendingBadge }}
                </span>
              </div>
              <p class="mt-2 text-xs leading-6 text-slate-500">{{ section.description }}</p>
            </button>
          </div>
        </div>

        <div :class="surfaceClass">
          <p class="text-sm font-semibold text-slate-900">{{ copy.overview.oldEntryTitle }}</p>
          <p class="mt-2 text-xs leading-6 text-slate-500">{{ copy.overview.oldEntrySubtitle }}</p>
          <div class="mt-4 space-y-2">
            <a
              v-for="link in copy.overview.oldEntryLinks"
              :key="link.href"
              :href="link.href"
              class="flex items-center justify-between rounded-[18px] border border-slate-200/80 bg-white/72 px-3.5 py-3 text-sm text-slate-600 transition hover:border-miku/30 hover:text-miku"
            >
              <span>{{ link.label }}</span>
              <span aria-hidden="true">›</span>
            </a>
          </div>
        </div>
      </div>

      <div class="space-y-5">
        <div :class="surfaceClass">
          <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div class="max-w-2xl">
              <p class="text-[11px] uppercase tracking-[0.24em] text-slate-400">{{ currentSection.eyebrow }}</p>
              <h2 class="mt-2 text-2xl font-semibold text-slate-900">{{ currentSection.title }}</h2>
              <p class="mt-2 text-sm leading-7 text-slate-600">{{ currentSection.description }}</p>
            </div>

            <div class="max-w-sm rounded-[20px] border border-slate-200/80 bg-white/72 px-4 py-3">
              <p class="text-[11px] uppercase tracking-[0.22em] text-slate-400">{{ copy.sections.affectedLabel }}</p>
              <div class="mt-3 flex flex-wrap gap-2">
                <span
                  v-for="item in currentSection.impact"
                  :key="item"
                  class="rounded-full border border-slate-200 bg-white px-3 py-1 text-xs font-medium text-slate-600"
                >
                  {{ item }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeSection === 'site-profile'" :class="surfaceClass">
          <div>
            <p class="text-sm font-semibold text-slate-900">{{ footerCopy.siteProfile.brandTitle }}</p>
            <p class="mt-1 text-xs leading-6 text-slate-500">{{ footerCopy.siteProfile.brandHint }}</p>
          </div>

          <div class="mt-5 grid gap-3 md:grid-cols-2">
            <input
              v-model="brandText"
              type="text"
              :placeholder="footerCopy.siteProfile.brandTextPlaceholder"
              :class="inputClass"
              :aria-label="footerCopy.siteProfile.brandTextLabel"
              :disabled="loading || savingSiteProfile"
            />
            <input
              v-model="siteTitle"
              type="text"
              :placeholder="footerCopy.siteProfile.siteTitlePlaceholder"
              :class="inputClass"
              :aria-label="footerCopy.siteProfile.siteTitleLabel"
              :disabled="loading || savingSiteProfile"
            />
            <input
              v-model="logoAlt"
              type="text"
              :placeholder="footerCopy.siteProfile.logoAltPlaceholder"
              :class="`${inputClass} md:col-span-2`"
              :aria-label="footerCopy.siteProfile.logoAltLabel"
              :disabled="loading || savingSiteProfile"
            />
          </div>

          <div class="mt-6">
            <p class="text-sm font-semibold text-slate-900">{{ footerCopy.siteProfile.seoTitle }}</p>
            <p class="mt-1 text-xs leading-6 text-slate-500">{{ footerCopy.siteProfile.seoHint }}</p>
          </div>

          <div class="mt-4 grid gap-3">
            <input
              v-model="siteUrl"
              type="text"
              :placeholder="footerCopy.siteProfile.siteUrlPlaceholder"
              :class="inputClass"
              :aria-label="footerCopy.siteProfile.siteUrlLabel"
              :disabled="loading || savingSiteProfile"
            />
            <textarea
              v-model="defaultDescription"
              rows="4"
              :placeholder="footerCopy.siteProfile.defaultDescriptionPlaceholder"
              :class="textareaClass"
              :aria-label="footerCopy.siteProfile.defaultDescriptionLabel"
              :disabled="loading || savingSiteProfile"
            />
            <input
              v-model="defaultSocialImage"
              type="text"
              :placeholder="footerCopy.siteProfile.defaultSocialImagePlaceholder"
              :class="inputClass"
              :aria-label="footerCopy.siteProfile.defaultSocialImageLabel"
              :disabled="loading || savingSiteProfile"
            />
          </div>
        </div>

        <div v-else-if="activeSection === 'home-hero'" :class="surfaceClass">
          <div>
            <p class="text-sm font-semibold text-slate-900">{{ copy.forms.homeHero.title }}</p>
            <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.homeHero.subtitle }}</p>
          </div>

          <div class="mt-5 grid gap-3">
            <input
              v-model="homeHeroTitle"
              type="text"
              :placeholder="copy.forms.homeHero.heroTitlePlaceholder"
              :class="inputClass"
              :aria-label="copy.forms.homeHero.heroTitleLabel"
              :disabled="loading || savingHomeHero"
            />
            <textarea
              v-model="homeHeroSubtitle"
              rows="4"
              :placeholder="copy.forms.homeHero.heroSubtitlePlaceholder"
              :class="textareaClass"
              :aria-label="copy.forms.homeHero.heroSubtitleLabel"
              :disabled="loading || savingHomeHero"
            />
          </div>
        </div>

        <div v-else-if="activeSection === 'home-assets'" :class="surfaceClass">
          <div class="flex flex-wrap items-center justify-between gap-3">
            <div>
              <p class="text-sm font-semibold text-slate-900">{{ copy.forms.homeAssets.title }}</p>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.homeAssets.subtitle }}</p>
            </div>

            <button
              type="button"
              class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
              :aria-label="copy.forms.homeAssets.addButtonAria"
              :disabled="loading || savingHomeAssets || homeHeroImages.length >= maxHomeAssetImages"
              @click="addHeroImage"
            >
              {{ copy.forms.homeAssets.addButton }}
            </button>
          </div>

          <div class="mt-4 space-y-2.5">
            <div
              v-for="(image, index) in homeHeroImages"
              :key="`settings-hero-image-${index}`"
              class="grid gap-2 md:grid-cols-[1fr_auto]"
            >
              <input
                v-model="homeHeroImages[index]"
                type="text"
                :placeholder="copy.forms.homeAssets.heroImagePlaceholder"
                :class="inputClass"
                :aria-label="`${copy.forms.homeAssets.heroImageLabelPrefix}${index + 1}`"
                :disabled="loading || savingHomeAssets"
              />
              <button
                type="button"
                class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                :aria-label="copy.forms.homeAssets.removeButtonAria"
                :disabled="loading || savingHomeAssets"
                @click="removeHeroImage(index)"
              >
                {{ copy.forms.homeAssets.removeButton }}
              </button>
            </div>

            <p v-if="homeHeroImages.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
              {{ copy.forms.homeAssets.emptyHint }}
            </p>
          </div>
        </div>

        <div v-else-if="activeSection === 'author-profile'" :class="surfaceClass">
          <div class="grid gap-4 lg:grid-cols-[auto_1fr] lg:items-start">
            <img
              :src="authorAvatarPreview"
              :alt="copy.preview.avatarAlt"
              class="h-20 w-20 rounded-[24px] border border-white/80 object-cover shadow-[0_14px_28px_rgba(15,23,42,0.08)]"
            />
            <div>
              <p class="text-sm font-semibold text-slate-900">{{ copy.forms.authorProfile.basicsTitle }}</p>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.authorProfile.basicsSubtitle }}</p>
            </div>
          </div>

          <div class="mt-5 grid gap-3 md:grid-cols-2">
            <input
              v-model="authorDisplayName"
              type="text"
              :placeholder="copy.forms.authorProfile.displayNamePlaceholder"
              :class="inputClass"
              :aria-label="copy.forms.authorProfile.displayNameLabel"
              :disabled="loading || savingAuthorProfile"
            />
            <input
              v-model="authorAvatarUrl"
              type="text"
              :placeholder="copy.forms.authorProfile.avatarPlaceholder"
              :class="inputClass"
              :aria-label="copy.forms.authorProfile.avatarLabel"
              :disabled="loading || savingAuthorProfile"
            />
            <input
              v-model="authorRole"
              type="text"
              :placeholder="copy.forms.authorProfile.rolePlaceholder"
              :class="inputClass"
              :aria-label="copy.forms.authorProfile.roleLabel"
              :disabled="loading || savingAuthorProfile"
            />
            <input
              v-model="authorLocation"
              type="text"
              :placeholder="copy.forms.authorProfile.locationPlaceholder"
              :class="inputClass"
              :aria-label="copy.forms.authorProfile.locationLabel"
              :disabled="loading || savingAuthorProfile"
            />
            <input
              v-model="authorSince"
              type="text"
              :placeholder="copy.forms.authorProfile.sincePlaceholder"
              :class="inputClass"
              :aria-label="copy.forms.authorProfile.sinceLabel"
              :disabled="loading || savingAuthorProfile"
            />
            <input
              v-model="authorContactEmail"
              type="email"
              :placeholder="copy.forms.authorProfile.contactEmailPlaceholder"
              :class="inputClass"
              :aria-label="copy.forms.authorProfile.contactEmailLabel"
              :disabled="loading || savingAuthorProfile"
            />
          </div>

          <div class="mt-4 grid gap-3">
            <textarea
              v-model="authorBio"
              rows="3"
              :placeholder="copy.forms.authorProfile.bioPlaceholder"
              :class="textareaClass"
              :aria-label="copy.forms.authorProfile.bioLabel"
              :disabled="loading || savingAuthorProfile"
            />
            <textarea
              v-model="authorAboutDescription"
              rows="4"
              :placeholder="copy.forms.authorProfile.aboutDescriptionPlaceholder"
              :class="textareaClass"
              :aria-label="copy.forms.authorProfile.aboutDescriptionLabel"
              :disabled="loading || savingAuthorProfile"
            />
            <textarea
              v-model="authorQuote"
              rows="3"
              :placeholder="copy.forms.authorProfile.quotePlaceholder"
              :class="textareaClass"
              :aria-label="copy.forms.authorProfile.quoteLabel"
              :disabled="loading || savingAuthorProfile"
            />
          </div>

          <div class="mt-6 flex flex-wrap items-center justify-between gap-3">
            <div>
              <p class="text-sm font-semibold text-slate-900">{{ copy.forms.authorProfile.skillsTitle }}</p>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.authorProfile.skillsSubtitle }}</p>
            </div>

            <button
              type="button"
              class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
              :aria-label="copy.forms.authorProfile.addSkillButtonAria"
              :disabled="loading || savingAuthorProfile || authorSkills.length >= maxAuthorSkills"
              @click="addAuthorSkill"
            >
              {{ copy.forms.authorProfile.addSkillButton }}
            </button>
          </div>

          <div class="mt-4 space-y-2.5">
            <div
              v-for="(skill, index) in authorSkills"
              :key="`author-skill-${index}`"
              class="grid gap-2 md:grid-cols-[1fr_auto]"
            >
              <input
                v-model="authorSkills[index]"
                type="text"
                :placeholder="copy.forms.authorProfile.skillPlaceholder"
                :class="inputClass"
                :aria-label="`${copy.forms.authorProfile.skillsTitle}${index + 1}`"
                :disabled="loading || savingAuthorProfile"
              />
              <button
                type="button"
                class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                :aria-label="copy.forms.authorProfile.removeButtonAria"
                :disabled="loading || savingAuthorProfile"
                @click="removeAuthorSkill(index)"
              >
                {{ copy.forms.authorProfile.removeButton }}
              </button>
            </div>

            <p v-if="authorSkills.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
              {{ copy.forms.authorProfile.emptySkillsHint }}
            </p>
          </div>

          <div class="mt-6 flex flex-wrap items-center justify-between gap-3">
            <div>
              <p class="text-sm font-semibold text-slate-900">{{ copy.forms.authorProfile.nowTitle }}</p>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.authorProfile.nowSubtitle }}</p>
            </div>

            <button
              type="button"
              class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
              :aria-label="copy.forms.authorProfile.addNowButtonAria"
              :disabled="loading || savingAuthorProfile || authorNowItems.length >= maxAuthorNowItems"
              @click="addAuthorNowItem"
            >
              {{ copy.forms.authorProfile.addNowButton }}
            </button>
          </div>

          <div class="mt-4 space-y-2.5">
            <div
              v-for="(item, index) in authorNowItems"
              :key="`author-now-${index}`"
              class="grid gap-2 md:grid-cols-[1fr_auto]"
            >
              <input
                v-model="authorNowItems[index]"
                type="text"
                :placeholder="copy.forms.authorProfile.nowPlaceholder"
                :class="inputClass"
                :aria-label="`${copy.forms.authorProfile.nowTitle}${index + 1}`"
                :disabled="loading || savingAuthorProfile"
              />
              <button
                type="button"
                class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                :aria-label="copy.forms.authorProfile.removeButtonAria"
                :disabled="loading || savingAuthorProfile"
                @click="removeAuthorNowItem(index)"
              >
                {{ copy.forms.authorProfile.removeButton }}
              </button>
            </div>

            <p v-if="authorNowItems.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
              {{ copy.forms.authorProfile.emptyNowHint }}
            </p>
          </div>

          <div class="mt-6 flex flex-wrap items-center justify-between gap-3">
            <div>
              <p class="text-sm font-semibold text-slate-900">{{ copy.forms.authorProfile.socialTitle }}</p>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.authorProfile.socialSubtitle }}</p>
            </div>

            <button
              type="button"
              class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
              :aria-label="copy.forms.authorProfile.addSocialButtonAria"
              :disabled="loading || savingAuthorProfile || authorSocialLinks.length >= maxAuthorSocialLinks"
              @click="addAuthorSocialLink"
            >
              {{ copy.forms.authorProfile.addSocialButton }}
            </button>
          </div>

          <div class="mt-4 space-y-2.5">
            <div
              v-for="(link, index) in authorSocialLinks"
              :key="`author-social-${index}`"
              class="grid gap-2 md:grid-cols-[minmax(0,0.9fr)_minmax(0,1.3fr)_minmax(0,0.8fr)_auto]"
            >
              <input
                v-model="authorSocialLinks[index].label"
                type="text"
                :placeholder="copy.forms.authorProfile.socialLabelPlaceholder"
                :class="inputClass"
                :aria-label="`${copy.forms.authorProfile.socialTitle}${index + 1}`"
                :disabled="loading || savingAuthorProfile"
              />
              <input
                v-model="authorSocialLinks[index].href"
                type="text"
                :placeholder="copy.forms.authorProfile.socialHrefPlaceholder"
                :class="inputClass"
                :aria-label="`${copy.forms.authorProfile.socialTitle}${index + 1} href`"
                :disabled="loading || savingAuthorProfile"
              />
              <input
                v-model="authorSocialLinks[index].iconKey"
                type="text"
                :placeholder="copy.forms.authorProfile.socialIconPlaceholder"
                :class="inputClass"
                :aria-label="`${copy.forms.authorProfile.socialTitle}${index + 1} icon`"
                :disabled="loading || savingAuthorProfile"
              />
              <button
                type="button"
                class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                :aria-label="copy.forms.authorProfile.removeButtonAria"
                :disabled="loading || savingAuthorProfile"
                @click="removeAuthorSocialLink(index)"
              >
                {{ copy.forms.authorProfile.removeButton }}
              </button>
            </div>

            <p v-if="authorSocialLinks.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
              {{ copy.forms.authorProfile.emptySocialHint }}
            </p>
          </div>

          <div class="mt-6 flex flex-wrap items-center justify-between gap-3">
            <div>
              <p class="text-sm font-semibold text-slate-900">{{ copy.forms.authorProfile.contactLinksTitle }}</p>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.authorProfile.contactLinksSubtitle }}</p>
            </div>

            <button
              type="button"
              class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
              :aria-label="copy.forms.authorProfile.addContactButtonAria"
              :disabled="loading || savingAuthorProfile || authorContactLinks.length >= maxAuthorContactLinks"
              @click="addAuthorContactLink"
            >
              {{ copy.forms.authorProfile.addContactButton }}
            </button>
          </div>

          <div class="mt-4 space-y-2.5">
            <div
              v-for="(link, index) in authorContactLinks"
              :key="`author-contact-${index}`"
              class="grid gap-2 md:grid-cols-[minmax(0,0.95fr)_minmax(0,1.45fr)_auto]"
            >
              <input
                v-model="authorContactLinks[index].label"
                type="text"
                :placeholder="copy.forms.authorProfile.contactLabelPlaceholder"
                :class="inputClass"
                :aria-label="`${copy.forms.authorProfile.contactLinksTitle}${index + 1}`"
                :disabled="loading || savingAuthorProfile"
              />
              <input
                v-model="authorContactLinks[index].href"
                type="text"
                :placeholder="copy.forms.authorProfile.contactHrefPlaceholder"
                :class="inputClass"
                :aria-label="`${copy.forms.authorProfile.contactLinksTitle}${index + 1} href`"
                :disabled="loading || savingAuthorProfile"
              />
              <button
                type="button"
                class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                :aria-label="copy.forms.authorProfile.removeButtonAria"
                :disabled="loading || savingAuthorProfile"
                @click="removeAuthorContactLink(index)"
              >
                {{ copy.forms.authorProfile.removeButton }}
              </button>
            </div>

            <p v-if="authorContactLinks.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
              {{ copy.forms.authorProfile.emptyContactHint }}
            </p>
          </div>
        </div>

        <div v-else-if="activeSection === 'site-integrations'" :class="surfaceClass">
          <div>
            <p class="text-sm font-semibold text-slate-900">{{ copy.forms.siteIntegrations.title }}</p>
            <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.siteIntegrations.subtitle }}</p>
          </div>

          <div class="mt-5 grid gap-3 md:grid-cols-2">
            <input
              v-model="integrationsGithubUsername"
              type="text"
              :placeholder="copy.forms.siteIntegrations.githubUsernamePlaceholder"
              :class="inputClass"
              :aria-label="copy.forms.siteIntegrations.githubUsernameLabel"
              :disabled="loading || savingSiteIntegrations || refreshingIntegrationStatus"
            />
            <input
              v-model="integrationsWeatherLocation"
              type="text"
              :placeholder="copy.forms.siteIntegrations.weatherLocationPlaceholder"
              :class="inputClass"
              :aria-label="copy.forms.siteIntegrations.weatherLocationLabel"
              :disabled="loading || savingSiteIntegrations || refreshingIntegrationStatus"
            />
          </div>

          <div class="mt-6">
            <p class="text-sm font-semibold text-slate-900">{{ copy.forms.siteIntegrations.widgetsTitle }}</p>
            <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.siteIntegrations.widgetsSubtitle }}</p>
          </div>

          <div class="mt-4 grid gap-3 sm:grid-cols-3">
            <label
              v-for="widget in widgetCards"
              :key="widget.key"
              class="flex cursor-pointer flex-col gap-3 rounded-[20px] border p-4 transition"
              :class="widgetCardClass(widget.enabled)"
            >
              <div class="flex items-start justify-between gap-3">
                <div>
                  <p class="text-sm font-semibold">{{ widget.label }}</p>
                  <p class="mt-1 text-xs leading-6">{{ widget.hint }}</p>
                </div>
                <input
                  v-if="widget.key === 'weather'"
                  v-model="integrationsShowWeather"
                  type="checkbox"
                  class="mt-1 h-4 w-4 rounded border-slate-300 text-miku focus:ring-miku/25"
                  :disabled="loading || savingSiteIntegrations"
                />
                <input
                  v-else-if="widget.key === 'music'"
                  v-model="integrationsShowMusic"
                  type="checkbox"
                  class="mt-1 h-4 w-4 rounded border-slate-300 text-miku focus:ring-miku/25"
                  :disabled="loading || savingSiteIntegrations"
                />
                <input
                  v-else
                  v-model="integrationsShowClock"
                  type="checkbox"
                  class="mt-1 h-4 w-4 rounded border-slate-300 text-miku focus:ring-miku/25"
                  :disabled="loading || savingSiteIntegrations"
                />
              </div>

              <span class="inline-flex items-center rounded-full px-2.5 py-1 text-[11px] font-medium" :class="widget.enabled ? 'bg-miku/12 text-miku' : 'bg-slate-100 text-slate-500'">
                {{ widget.enabled ? copy.preview.widgetOnLabel : copy.preview.widgetOffLabel }}
              </span>
            </label>
          </div>

          <div class="mt-6 flex flex-wrap items-center justify-between gap-3">
            <div>
              <p class="text-sm font-semibold text-slate-900">{{ copy.forms.siteIntegrations.statusTitle }}</p>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.siteIntegrations.statusSubtitle }}</p>
            </div>

            <button
              type="button"
              class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
              :disabled="loading || savingSiteIntegrations || refreshingIntegrationStatus"
              @click="refreshIntegrationStatus"
            >
              {{ refreshingIntegrationStatus ? copy.forms.siteIntegrations.refreshingButton : copy.forms.siteIntegrations.refreshButton }}
            </button>
          </div>

          <div class="mt-4 grid gap-3 md:grid-cols-2">
            <div
              v-for="item in integrationStatusCards"
              :key="item.title"
              class="rounded-[20px] border p-4 transition"
              :class="probeToneClass(item.state)"
            >
              <div class="flex items-center justify-between gap-3">
                <p class="text-sm font-semibold">{{ item.title }}</p>
                <span class="rounded-full px-2.5 py-1 text-[11px] font-medium" :class="probeToneClass(item.state)">
                  {{
                    item.state === 'ok'
                      ? copy.preview.statusSuccess
                      : item.state === 'error'
                        ? copy.preview.statusError
                        : item.state === 'loading'
                          ? copy.preview.statusLoading
                          : copy.preview.statusIdle
                  }}
                </span>
              </div>
              <p class="mt-3 text-sm font-medium">{{ item.message }}</p>
              <p class="mt-1 text-xs leading-6 opacity-85">{{ item.detail }}</p>
            </div>
          </div>
        </div>

        <div v-else-if="activeSection === 'site-footer'" :class="surfaceClass">
          <div>
            <p class="text-sm font-semibold text-slate-900">{{ footerCopy.form.icpTitle }}</p>
            <p class="mt-1 text-xs leading-6 text-slate-500">{{ footerCopy.form.icpHint }}</p>
          </div>

          <div class="mt-5 grid gap-3 md:grid-cols-2">
            <input
              v-model="icpText"
              type="text"
              :placeholder="footerCopy.form.icpPlaceholder"
              :class="inputClass"
              :aria-label="footerCopy.form.icpLabel"
              :disabled="loading || savingFooter"
            />
            <input
              v-model="icpLink"
              type="text"
              :placeholder="footerCopy.form.icpLinkPlaceholder"
              :class="inputClass"
              :aria-label="footerCopy.form.icpLinkLabel"
              :disabled="loading || savingFooter"
            />
          </div>

          <div class="mt-6 flex flex-wrap items-center justify-between gap-3">
            <div>
              <p class="text-sm font-semibold text-slate-900">{{ footerCopy.form.customTitle }}</p>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ footerCopy.form.customHint }}</p>
            </div>

            <button
              type="button"
              class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
              :aria-label="footerCopy.form.addButtonAria"
              :disabled="loading || savingFooter || customTexts.length >= maxCustomTexts"
              @click="addLine"
            >
              {{ footerCopy.form.addButton }}
            </button>
          </div>

          <div class="mt-4 space-y-2.5">
            <div
              v-for="(line, index) in customTexts"
              :key="`settings-footer-line-${index}`"
              class="grid gap-2 md:grid-cols-[1fr_auto]"
            >
              <input
                v-model="customTexts[index]"
                type="text"
                :placeholder="footerCopy.form.customLinePlaceholder"
                :class="inputClass"
                :aria-label="`${footerCopy.form.customLineLabelPrefix}${index + 1}`"
                :disabled="loading || savingFooter"
              />
              <button
                type="button"
                class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                :aria-label="footerCopy.form.removeButtonAria"
                :disabled="loading || savingFooter"
                @click="removeLine(index)"
              >
                {{ footerCopy.form.removeButton }}
              </button>
            </div>

            <p v-if="customTexts.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
              {{ footerCopy.form.emptyCustomHint }}
            </p>
          </div>
        </div>

        <div v-else :class="surfaceClass">
          <div class="grid gap-4 lg:grid-cols-[auto_1fr] lg:items-start">
            <img
              :src="profileAvatarPreview"
              :alt="copy.preview.avatarAlt"
              class="h-20 w-20 rounded-[24px] border border-white/80 object-cover shadow-[0_14px_28px_rgba(15,23,42,0.08)]"
            />
            <div>
              <p class="text-sm font-semibold text-slate-900">{{ profileCopy.profile.title }}</p>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ profileCopy.profile.subtitle }}</p>
            </div>
          </div>

          <div class="mt-5 grid gap-3 md:grid-cols-2">
            <input
              v-model="profileDisplayName"
              type="text"
              :placeholder="profileCopy.profile.displayNamePlaceholder"
              :class="inputClass"
              :aria-label="profileCopy.profile.displayNameLabel"
              :disabled="loading || savingAdminProfile"
            />
            <input
              v-model="profileAvatarURL"
              type="text"
              :placeholder="profileCopy.profile.avatarPlaceholder"
              :class="inputClass"
              :aria-label="profileCopy.profile.avatarLabel"
              :disabled="loading || savingAdminProfile"
            />
          </div>

          <div class="mt-6">
            <p class="text-sm font-semibold text-slate-900">{{ profileCopy.account.title }}</p>
            <p class="mt-1 text-xs leading-6 text-slate-500">{{ profileCopy.account.subtitle }}</p>
          </div>

          <div class="mt-4 grid gap-3 md:grid-cols-2">
            <input
              v-model="accountUsername"
              type="text"
              :placeholder="profileCopy.account.usernamePlaceholder"
              :class="inputClass"
              :aria-label="profileCopy.account.usernameLabel"
              :disabled="loading || savingAdminProfile"
            />
            <input
              v-model="accountEmail"
              type="email"
              :placeholder="profileCopy.account.emailPlaceholder"
              :class="inputClass"
              :aria-label="profileCopy.account.emailLabel"
              :disabled="loading || savingAdminProfile"
            />
            <input
              v-model="accountPassword"
              type="password"
              :placeholder="profileCopy.account.passwordPlaceholder"
              :class="`${inputClass} md:col-span-2`"
              :aria-label="profileCopy.account.passwordLabel"
              autocomplete="new-password"
              :disabled="loading || savingAdminProfile"
            />
          </div>
        </div>
      </div>

      <div class="space-y-5 xl:sticky xl:top-6 xl:self-start">
        <div :class="surfaceClass">
          <p class="text-[11px] uppercase tracking-[0.24em] text-slate-400">{{ copy.preview.panelTitle }}</p>
          <h3 class="mt-2 text-lg font-semibold text-slate-900">{{ currentSection.previewTitle }}</h3>
          <p class="mt-1 text-xs leading-6 text-slate-500">{{ currentSection.previewSubtitle }}</p>

          <div v-if="activeSection === 'site-profile'" class="mt-5 space-y-4">
            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs uppercase tracking-[0.2em] text-slate-400">SEO</p>
              <p class="mt-3 text-sm font-semibold text-slate-900">{{ copy.preview.seoLabel }}</p>
              <p class="mt-1 text-sm text-slate-600">{{ sitePreviewTitle }}</p>
              <p class="mt-3 text-xs font-medium text-slate-500">{{ copy.preview.urlLabel }}</p>
              <p class="mt-1 break-all text-xs leading-6 text-slate-500">{{ sitePreviewURL }}</p>
              <p class="mt-3 text-xs leading-6 text-slate-500">{{ sitePreviewDescription }}</p>
            </div>

            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs font-medium text-slate-500">{{ copy.preview.socialImageLabel }}</p>
              <img
                v-if="socialImagePreview"
                :src="socialImagePreview"
                alt=""
                class="mt-3 h-36 w-full rounded-[18px] border border-slate-200 object-cover"
              />
              <p v-else class="mt-3 rounded-[18px] border border-dashed border-slate-200 px-4 py-6 text-center text-xs text-slate-400">
                {{ copy.actions.emptyValue }}
              </p>
            </div>
          </div>

          <div v-else-if="activeSection === 'home-hero'" class="mt-5 space-y-4">
            <div class="rounded-[20px] border border-slate-200/80 bg-[linear-gradient(145deg,rgba(255,255,255,0.92),rgba(229,255,252,0.92),rgba(248,244,255,0.88))] p-5">
              <p class="text-[11px] uppercase tracking-[0.24em] text-miku/80">{{ copy.preview.heroBadge }}</p>
              <p class="mt-4 text-xs font-medium text-slate-500">{{ copy.preview.heroTitleLabel }}</p>
              <p class="mt-2 text-2xl font-semibold tracking-[0.03em] text-slate-900">{{ homeHeroTitlePreview }}</p>
              <p class="mt-4 text-xs font-medium text-slate-500">{{ copy.preview.heroSubtitleLabel }}</p>
              <p class="mt-2 rounded-[18px] border border-white/80 bg-white/72 px-4 py-3 text-sm leading-7 text-slate-600">
                {{ homeHeroSubtitlePreview }}
              </p>
            </div>
          </div>

          <div v-else-if="activeSection === 'home-assets'" class="mt-5 space-y-4">
            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs font-medium text-slate-500">{{ copy.preview.heroImagesLabel }}</p>
              <p class="mt-2 text-sm font-semibold text-slate-900">{{ homeAssetsCountPreview }}</p>
            </div>

            <div class="grid gap-3 sm:grid-cols-2">
              <div
                v-for="(image, index) in homeHeroImagesPreview"
                :key="`hero-assets-preview-${index}`"
                class="overflow-hidden rounded-[20px] border border-slate-200/80 bg-white/80"
              >
                <img
                  :src="image"
                  :alt="`${copy.forms.homeAssets.heroImageLabelPrefix}${index + 1}`"
                  class="h-32 w-full object-cover"
                />
                <div class="px-3 py-2">
                  <p class="truncate text-xs text-slate-500">{{ image }}</p>
                </div>
              </div>
            </div>
          </div>

          <div v-else-if="activeSection === 'author-profile'" class="mt-5 space-y-4">
            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <div class="flex items-start gap-3">
                <img
                  :src="authorAvatarPreview"
                  :alt="copy.preview.avatarAlt"
                  class="h-16 w-16 rounded-[18px] border border-slate-200 object-cover"
                />
                <div class="min-w-0">
                  <p class="text-sm font-semibold text-slate-900">{{ authorDisplayNamePreview }}</p>
                  <p class="mt-1 text-xs text-slate-500">{{ authorRolePreview }}</p>
                  <p class="mt-2 text-xs leading-6 text-slate-500">{{ authorMetaPreview }}</p>
                </div>
              </div>
            </div>

            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs font-medium text-slate-500">{{ copy.preview.authorBioLabel }}</p>
              <p class="mt-2 text-sm leading-7 text-slate-600">{{ authorBioPreview }}</p>
              <p class="mt-4 text-xs font-medium text-slate-500">{{ copy.preview.authorAboutLabel }}</p>
              <p class="mt-2 text-sm leading-7 text-slate-600">{{ authorAboutPreview }}</p>
              <blockquote class="mt-4 rounded-[16px] border border-slate-200 bg-slate-50/80 px-3 py-2.5 text-xs leading-6 text-slate-500">
                {{ authorQuotePreview }}
              </blockquote>
            </div>

            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs font-medium text-slate-500">{{ copy.preview.authorTagsLabel }}</p>
              <div class="mt-3 flex flex-wrap gap-2">
                <span
                  v-for="item in authorSkillsPreview"
                  :key="`author-preview-skill-${item}`"
                  class="rounded-full border border-slate-200 bg-white px-3 py-1 text-xs font-medium text-slate-600"
                >
                  {{ item }}
                </span>
              </div>

              <p class="mt-4 text-xs font-medium text-slate-500">{{ copy.preview.authorSocialLabel }}</p>
              <div class="mt-3 flex flex-wrap gap-2">
                <span
                  v-for="item in authorSocialLinksPreview"
                  :key="`${item.label}-${item.href}`"
                  class="rounded-full border border-miku/20 bg-miku/10 px-3 py-1 text-xs font-medium text-miku"
                >
                  {{ item.label }}
                </span>
                <span
                  v-if="authorSocialLinksPreview.length === 0"
                  class="rounded-full border border-slate-200 bg-white px-3 py-1 text-xs font-medium text-slate-500"
                >
                  {{ copy.actions.emptyValue }}
                </span>
              </div>
            </div>
          </div>

          <div v-else-if="activeSection === 'site-integrations'" class="mt-5 space-y-4">
            <div
              v-for="item in integrationStatusCards"
              :key="`preview-${item.title}`"
              class="rounded-[20px] border p-4 transition"
              :class="probeToneClass(item.state)"
            >
              <div class="flex items-center justify-between gap-3">
                <p class="text-sm font-semibold">{{ item.title }}</p>
                <span class="rounded-full px-2.5 py-1 text-[11px] font-medium" :class="probeToneClass(item.state)">
                  {{
                    item.state === 'ok'
                      ? copy.preview.statusSuccess
                      : item.state === 'error'
                        ? copy.preview.statusError
                        : item.state === 'loading'
                          ? copy.preview.statusLoading
                          : copy.preview.statusIdle
                  }}
                </span>
              </div>
              <p class="mt-3 text-sm font-medium">{{ item.message }}</p>
              <p class="mt-1 text-xs leading-6 opacity-85">{{ item.detail }}</p>
            </div>

            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs font-medium text-slate-500">{{ copy.preview.widgetStatusTitle }}</p>
              <div class="mt-3 space-y-2">
                <div
                  v-for="widget in widgetCards"
                  :key="`preview-widget-${widget.key}`"
                  class="flex items-center justify-between rounded-[16px] border px-3 py-2"
                  :class="widgetCardClass(widget.enabled)"
                >
                  <span class="text-sm font-medium">{{ widget.label }}</span>
                  <span class="text-xs">{{ widget.enabled ? copy.preview.widgetOnLabel : copy.preview.widgetOffLabel }}</span>
                </div>
              </div>
            </div>
          </div>

          <div v-else-if="activeSection === 'site-footer'" class="mt-5 space-y-3">
            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-sm font-semibold uppercase tracking-[0.16em] text-slate-700">{{ footerBrandPreview }}</p>
              <p class="mt-3 text-xs font-medium text-slate-500">{{ footerCopy.form.icpLabel }}</p>
              <p class="mt-1 text-sm text-slate-600">{{ footerICPPreview }}</p>
            </div>

            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs font-medium text-slate-500">{{ copy.preview.footerLabel }}</p>
              <div class="mt-3 space-y-2">
                <p
                  v-for="(line, index) in footerLinesPreview"
                  :key="`preview-footer-${index}`"
                  class="rounded-[16px] border border-slate-200 bg-white px-3 py-2 text-xs leading-6 text-slate-600"
                >
                  {{ line }}
                </p>
              </div>
            </div>
          </div>

          <div v-else class="mt-5 space-y-4">
            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <div class="flex items-center gap-3">
                <img
                  :src="profileAvatarPreview"
                  :alt="copy.preview.avatarAlt"
                  class="h-14 w-14 rounded-[18px] border border-slate-200 object-cover"
                />
                <div>
                  <p class="text-sm font-semibold text-slate-900">{{ profileDisplayNamePreview }}</p>
                  <p class="mt-1 text-xs text-slate-500">{{ copy.preview.accountLabel }}</p>
                  <p class="mt-1 text-xs text-slate-500">{{ accountIdentityPreview }}</p>
                </div>
              </div>
            </div>

            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs font-medium text-slate-500">{{ copy.sections.affectedLabel }}</p>
              <div class="mt-3 flex flex-wrap gap-2">
                <span
                  v-for="item in currentSection.impact"
                  :key="`admin-impact-${item}`"
                  class="rounded-full border border-slate-200 bg-white px-3 py-1 text-xs font-medium text-slate-600"
                >
                  {{ item }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div :class="surfaceClass">
          <p class="text-sm font-semibold text-slate-900">{{ copy.preview.boundaryTitle }}</p>
          <p class="mt-2 text-xs leading-6 text-slate-500">{{ copy.preview.boundaryDescription }}</p>

          <div class="mt-4 rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-slate-400">{{ copy.preview.impactTitle }}</p>
            <p class="mt-2 text-xs leading-6 text-slate-500">{{ copy.preview.impactDescription }}</p>
          </div>
        </div>
      </div>
    </div>

    <div class="sticky bottom-4 z-20">
      <div class="flex flex-col gap-3 rounded-[24px] border border-white/80 bg-white/86 px-5 py-4 shadow-[0_22px_44px_rgba(15,23,42,0.12)] backdrop-blur-xl lg:flex-row lg:items-center lg:justify-between">
        <div class="min-w-0">
          <p class="text-[11px] uppercase tracking-[0.24em] text-slate-400">{{ copy.actions.statusScopePrefix }}</p>
          <p class="mt-1 text-sm font-semibold text-slate-900">{{ currentSection.title }}</p>
          <p class="mt-1 text-xs leading-6" :class="statusClass">{{ statusText }}</p>
        </div>

        <div class="flex flex-wrap items-center gap-3">
          <button
            v-if="canResetCurrentSection"
            type="button"
            class="rounded-full border border-slate-200 bg-white px-4 py-2 text-sm text-slate-600 transition hover:border-miku/30 hover:text-miku disabled:cursor-not-allowed disabled:opacity-45"
            :disabled="loading || isCurrentSaving"
            @click="resetCurrentSection"
          >
            {{ copy.actions.resetButton }}
          </button>

          <MikuButton
            type="button"
            variant="solid"
            :disabled="loading || isCurrentSaving || !currentSectionDirty"
            @click="saveCurrentSection"
          >
            {{ isCurrentSaving ? copy.actions.savingButton : copy.actions.saveButton }}
          </MikuButton>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useStore } from '@nanostores/vue'
import { computed, onMounted, ref } from 'vue'

import { adminCopy } from '../../content/copy'
import { api, ApiError } from '../../lib/api'
import type {
  AuthorContactLink,
  AuthorProfileSettings,
  AuthorSocialLink,
} from '../../lib/author-profile'
import type { SiteIntegrationsSettings } from '../../lib/site-integrations'
import {
  hydrateAuth,
  updateMyAccount,
  updateMyProfile,
  authState,
} from '../../stores/auth'
import {
  authorProfileSettings,
  hydrateAuthorProfileSettings,
  resetAuthorProfileSettings,
  saveAuthorProfileSettings,
} from '../../stores/authorProfile'
import {
  hydrateSiteFooterSettings,
  resetSiteFooterSettings,
  saveSiteFooterSettings,
  siteFooterSettings,
  type SiteFooterSettings,
} from '../../stores/siteFooter'
import {
  hydrateHomeAssetsSettings,
  resetHomeAssetsSettings,
  saveHomeAssetsSettings,
  homeAssetsSettings,
  type HomeAssetsSettings,
} from '../../stores/homeAssets'
import {
  hydrateHomeHeroSettings,
  resetHomeHeroSettings,
  saveHomeHeroSettings,
  homeHeroSettings,
  type HomeHeroSettings,
} from '../../stores/homeHero'
import {
  hydrateSiteIntegrationsSettings,
  resetSiteIntegrationsSettings,
  saveSiteIntegrationsSettings,
  siteIntegrationsSettings,
} from '../../stores/siteIntegrations'
import {
  hydrateSiteProfileSettings,
  resetSiteProfileSettings,
  saveSiteProfileSettings,
  siteProfileSettings,
  type SiteProfileSettings,
} from '../../stores/siteProfile'
import { showToast } from '../../stores/ui'
import MikuButton from '../ui/MikuButton.vue'

interface AdminProfilePayload {
  username: string
  email: string
  display_name?: string
  avatar_url?: string
}

interface GitHubProbeResponse {
  profile?: {
    name?: string
    total_repos?: number | string
  }
}

interface WeatherProbeResponse {
  temp?: string
  desc?: string
  location?: string
}

type SectionKey =
  | 'site-profile'
  | 'home-hero'
  | 'home-assets'
  | 'author-profile'
  | 'site-integrations'
  | 'site-footer'
  | 'admin-profile'

type ProbeState = 'idle' | 'loading' | 'ok' | 'error'

interface ProbeStatus {
  state: ProbeState
  message: string
  detail: string
}

const copy = adminCopy.settingsCenter
const footerCopy = adminCopy.footerManager
const profileCopy = adminCopy.profileManager

const heroClass = 'rounded-[30px] border border-white/75 bg-[linear-gradient(135deg,rgba(255,255,255,0.94),rgba(229,255,252,0.92),rgba(248,244,255,0.9))] p-6 shadow-[0_24px_50px_rgba(15,23,42,0.08)] backdrop-blur-xl'
const surfaceClass = 'rounded-[24px] border border-white/75 bg-white/80 p-5 shadow-[0_18px_40px_rgba(15,23,42,0.08)] backdrop-blur-xl'
const inputClass = 'w-full rounded-[18px] border border-slate-200/90 bg-white/90 px-3.5 py-3 text-sm text-slate-700 outline-none transition focus:border-miku/45 focus:ring-2 focus:ring-miku/10 disabled:cursor-not-allowed disabled:opacity-60'
const textareaClass = `${inputClass} min-h-[110px] resize-y`
const maxCustomTexts = 8
const maxHomeAssetImages = 8
const maxAuthorSkills = 8
const maxAuthorNowItems = 8
const maxAuthorSocialLinks = 6
const maxAuthorContactLinks = 6

const homeHeroStore = useStore(homeHeroSettings)
const homeAssetsStore = useStore(homeAssetsSettings)
const siteProfileStore = useStore(siteProfileSettings)
const siteFooterStore = useStore(siteFooterSettings)
const authorProfileStore = useStore(authorProfileSettings)
const siteIntegrationsStore = useStore(siteIntegrationsSettings)
const auth = useStore(authState)

const activeSection = ref<SectionKey>('site-profile')
const loading = ref(false)
const savingHomeHero = ref(false)
const savingHomeAssets = ref(false)
const savingSiteProfile = ref(false)
const savingAuthorProfile = ref(false)
const savingSiteIntegrations = ref(false)
const savingFooter = ref(false)
const savingAdminProfile = ref(false)
const refreshingIntegrationStatus = ref(false)
const lastSavedAt = ref('')
const lastSavedSection = ref('')

const brandText = ref('')
const siteTitle = ref('')
const logoAlt = ref('')
const siteUrl = ref('')
const defaultDescription = ref('')
const defaultSocialImage = ref('')

const homeHeroTitle = ref('')
const homeHeroSubtitle = ref('')
const homeHeroImages = ref<string[]>([])

const authorDisplayName = ref('')
const authorAvatarUrl = ref('/picture/author.jpg')
const authorRole = ref('')
const authorBio = ref('')
const authorAboutDescription = ref('')
const authorLocation = ref('')
const authorSince = ref('')
const authorSkills = ref<string[]>([])
const authorNowItems = ref<string[]>([])
const authorQuote = ref('')
const authorContactEmail = ref('')
const authorSocialLinks = ref<AuthorSocialLink[]>([])
const authorContactLinks = ref<AuthorContactLink[]>([])

const integrationsGithubUsername = ref('')
const integrationsWeatherLocation = ref('')
const integrationsShowWeather = ref(true)
const integrationsShowMusic = ref(true)
const integrationsShowClock = ref(true)

const icpText = ref('')
const icpLink = ref('')
const customTexts = ref<string[]>([])

const profileDisplayName = ref('')
const profileAvatarURL = ref('/picture/author.jpg')
const accountUsername = ref('')
const accountEmail = ref('')
const accountPassword = ref('')

const homeHeroSnapshot = ref('')
const homeAssetsSnapshot = ref('')
const siteProfileSnapshot = ref('')
const authorProfileSnapshot = ref('')
const siteIntegrationsSnapshot = ref('')
const siteFooterSnapshot = ref('')
const adminProfileSnapshot = ref('')

const githubProbe = ref<ProbeStatus>({
  state: 'idle',
  message: copy.preview.statusIdle,
  detail: copy.forms.siteIntegrations.githubMissingHint,
})

const weatherProbe = ref<ProbeStatus>({
  state: 'idle',
  message: copy.preview.statusIdle,
  detail: copy.forms.siteIntegrations.weatherDefaultHint,
})

function trimText(input: string) {
  return input.trim()
}

function cleanFooterLines(lines: string[]) {
  return lines
    .map((line) => trimText(line))
    .filter((line) => line.length > 0)
    .slice(0, maxCustomTexts)
}

function cleanHeroImages(images: string[]) {
  return images
    .map((image) => trimText(image))
    .filter((image, index, list) => image.length > 0 && list.indexOf(image) === index)
    .slice(0, maxHomeAssetImages)
}

function cleanTextItems(items: string[], limit: number) {
  return items
    .map((item) => trimText(item))
    .filter((item, index, list) => item.length > 0 && list.indexOf(item) === index)
    .slice(0, limit)
}

function cleanAuthorSocialDraftLinks(items: AuthorSocialLink[]) {
  const result: AuthorSocialLink[] = []
  for (const item of items) {
    const label = trimText(item.label)
    const href = trimText(item.href)
    const iconKey = trimText(item.iconKey)
    if (!label || !href) {
      continue
    }
    if (result.some((current) => current.label === label && current.href === href)) {
      continue
    }
    result.push({ label, href, iconKey })
    if (result.length >= maxAuthorSocialLinks) {
      break
    }
  }
  return result
}

function cleanAuthorContactDraftLinks(items: AuthorContactLink[]) {
  const result: AuthorContactLink[] = []
  for (const item of items) {
    const label = trimText(item.label)
    const href = trimText(item.href)
    if (!label || !href) {
      continue
    }
    if (result.some((current) => current.label === label && current.href === href)) {
      continue
    }
    result.push({ label, href })
    if (result.length >= maxAuthorContactLinks) {
      break
    }
  }
  return result
}

function formatTime(date: Date) {
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
}

function getSiteProfileDraft(): SiteProfileSettings {
  return {
    brandText: trimText(brandText.value),
    siteTitle: trimText(siteTitle.value),
    logoAlt: trimText(logoAlt.value),
    siteUrl: trimText(siteUrl.value),
    defaultDescription: trimText(defaultDescription.value),
    defaultSocialImage: trimText(defaultSocialImage.value),
  }
}

function getHomeHeroDraft(): HomeHeroSettings {
  return {
    heroTitle: trimText(homeHeroTitle.value),
    heroSubtitle: trimText(homeHeroSubtitle.value),
  }
}

function getHomeAssetsDraft(): HomeAssetsSettings {
  return {
    heroImages: cleanHeroImages(homeHeroImages.value),
  }
}

function getAuthorProfileDraft(): AuthorProfileSettings {
  return {
    displayName: trimText(authorDisplayName.value),
    avatarUrl: trimText(authorAvatarUrl.value),
    role: trimText(authorRole.value),
    bio: trimText(authorBio.value),
    aboutDescription: trimText(authorAboutDescription.value),
    location: trimText(authorLocation.value),
    since: trimText(authorSince.value),
    skills: cleanTextItems(authorSkills.value, maxAuthorSkills),
    nowItems: cleanTextItems(authorNowItems.value, maxAuthorNowItems),
    quote: trimText(authorQuote.value),
    contactEmail: trimText(authorContactEmail.value),
    socialLinks: cleanAuthorSocialDraftLinks(authorSocialLinks.value),
    contactLinks: cleanAuthorContactDraftLinks(authorContactLinks.value),
  }
}

function getSiteIntegrationsDraft(): SiteIntegrationsSettings {
  return {
    githubUsername: trimText(integrationsGithubUsername.value),
    weatherLocation: trimText(integrationsWeatherLocation.value),
    showWeather: integrationsShowWeather.value,
    showMusic: integrationsShowMusic.value,
    showClock: integrationsShowClock.value,
  }
}

function getSiteFooterDraft(): SiteFooterSettings {
  return {
    icpText: trimText(icpText.value),
    icpLink: trimText(icpLink.value),
    customTexts: cleanFooterLines(customTexts.value),
  }
}

function getAdminProfileDraft() {
  return {
    displayName: trimText(profileDisplayName.value),
    avatarURL: trimText(profileAvatarURL.value),
    username: trimText(accountUsername.value),
    email: trimText(accountEmail.value),
    password: trimText(accountPassword.value),
  }
}

function serializeSiteProfile(settings: SiteProfileSettings) {
  return JSON.stringify(settings)
}

function serializeHomeHero(settings: HomeHeroSettings) {
  return JSON.stringify(settings)
}

function serializeHomeAssets(settings: HomeAssetsSettings) {
  return JSON.stringify(settings)
}

function serializeAuthorProfile(settings: AuthorProfileSettings) {
  return JSON.stringify(settings)
}

function serializeSiteIntegrations(settings: SiteIntegrationsSettings) {
  return JSON.stringify(settings)
}

function serializeSiteFooter(settings: SiteFooterSettings) {
  return JSON.stringify(settings)
}

function serializeAdminProfile() {
  return JSON.stringify(getAdminProfileDraft())
}

function syncSiteProfileFromStore() {
  brandText.value = siteProfileStore.value.brandText
  siteTitle.value = siteProfileStore.value.siteTitle
  logoAlt.value = siteProfileStore.value.logoAlt
  siteUrl.value = siteProfileStore.value.siteUrl
  defaultDescription.value = siteProfileStore.value.defaultDescription
  defaultSocialImage.value = siteProfileStore.value.defaultSocialImage
}

function syncHomeHeroFromStore() {
  homeHeroTitle.value = homeHeroStore.value.heroTitle
  homeHeroSubtitle.value = homeHeroStore.value.heroSubtitle
}

function syncHomeAssetsFromStore() {
  homeHeroImages.value = [...homeAssetsStore.value.heroImages]
}

function syncAuthorProfileFromStore() {
  authorDisplayName.value = authorProfileStore.value.displayName
  authorAvatarUrl.value = authorProfileStore.value.avatarUrl
  authorRole.value = authorProfileStore.value.role
  authorBio.value = authorProfileStore.value.bio
  authorAboutDescription.value = authorProfileStore.value.aboutDescription
  authorLocation.value = authorProfileStore.value.location
  authorSince.value = authorProfileStore.value.since
  authorSkills.value = [...authorProfileStore.value.skills]
  authorNowItems.value = [...authorProfileStore.value.nowItems]
  authorQuote.value = authorProfileStore.value.quote
  authorContactEmail.value = authorProfileStore.value.contactEmail
  authorSocialLinks.value = authorProfileStore.value.socialLinks.map((item) => ({ ...item }))
  authorContactLinks.value = authorProfileStore.value.contactLinks.map((item) => ({ ...item }))
}

function syncSiteIntegrationsFromStore() {
  integrationsGithubUsername.value = siteIntegrationsStore.value.githubUsername
  integrationsWeatherLocation.value = siteIntegrationsStore.value.weatherLocation
  integrationsShowWeather.value = siteIntegrationsStore.value.showWeather
  integrationsShowMusic.value = siteIntegrationsStore.value.showMusic
  integrationsShowClock.value = siteIntegrationsStore.value.showClock
}

function syncFooterFromStore() {
  icpText.value = siteFooterStore.value.icpText
  icpLink.value = siteFooterStore.value.icpLink
  customTexts.value = [...siteFooterStore.value.customTexts]
}

async function loadAdminProfile() {
  const hydratedUser = await hydrateAuth()
  const current = hydratedUser ?? auth.value.user

  if (current) {
    profileDisplayName.value = current.name || current.username || ''
    profileAvatarURL.value = current.avatar || '/picture/author.jpg'
    accountUsername.value = current.username || ''
    accountEmail.value = current.email || ''
  }

  try {
    const me = await api.get<AdminProfilePayload>('/auth/me')
    profileDisplayName.value = trimText(me.display_name || '') || me.username
    profileAvatarURL.value = trimText(me.avatar_url || '') || '/picture/author.jpg'
    accountUsername.value = trimText(me.username || '')
    accountEmail.value = trimText(me.email || '')
  } catch {
    // keep store fallback
  }

  accountPassword.value = ''
}

function commitSnapshots() {
  homeHeroSnapshot.value = serializeHomeHero(getHomeHeroDraft())
  homeAssetsSnapshot.value = serializeHomeAssets(getHomeAssetsDraft())
  siteProfileSnapshot.value = serializeSiteProfile(getSiteProfileDraft())
  authorProfileSnapshot.value = serializeAuthorProfile(getAuthorProfileDraft())
  siteIntegrationsSnapshot.value = serializeSiteIntegrations(getSiteIntegrationsDraft())
  siteFooterSnapshot.value = serializeSiteFooter(getSiteFooterDraft())
  adminProfileSnapshot.value = serializeAdminProfile()
}

function markSaved(sectionTitle: string) {
  lastSavedSection.value = sectionTitle
  lastSavedAt.value = formatTime(new Date())
}

function addLine() {
  if (customTexts.value.length >= maxCustomTexts) {
    return
  }
  customTexts.value.push('')
}

function addHeroImage() {
  if (homeHeroImages.value.length >= maxHomeAssetImages) {
    return
  }
  homeHeroImages.value.push('')
}

function addAuthorSkill() {
  if (authorSkills.value.length >= maxAuthorSkills) {
    return
  }
  authorSkills.value.push('')
}

function addAuthorNowItem() {
  if (authorNowItems.value.length >= maxAuthorNowItems) {
    return
  }
  authorNowItems.value.push('')
}

function addAuthorSocialLink() {
  if (authorSocialLinks.value.length >= maxAuthorSocialLinks) {
    return
  }
  authorSocialLinks.value.push({ label: '', href: '', iconKey: 'link' })
}

function addAuthorContactLink() {
  if (authorContactLinks.value.length >= maxAuthorContactLinks) {
    return
  }
  authorContactLinks.value.push({ label: '', href: '' })
}

function removeHeroImage(index: number) {
  homeHeroImages.value = homeHeroImages.value.filter((_, currentIndex) => currentIndex !== index)
}

function removeLine(index: number) {
  customTexts.value = customTexts.value.filter((_, currentIndex) => currentIndex !== index)
}

function removeAuthorSkill(index: number) {
  authorSkills.value = authorSkills.value.filter((_, currentIndex) => currentIndex !== index)
}

function removeAuthorNowItem(index: number) {
  authorNowItems.value = authorNowItems.value.filter((_, currentIndex) => currentIndex !== index)
}

function removeAuthorSocialLink(index: number) {
  authorSocialLinks.value = authorSocialLinks.value.filter((_, currentIndex) => currentIndex !== index)
}

function removeAuthorContactLink(index: number) {
  authorContactLinks.value = authorContactLinks.value.filter((_, currentIndex) => currentIndex !== index)
}

function probeToneClass(state: ProbeState) {
  if (state === 'ok') {
    return 'border-miku/20 bg-miku/10 text-miku'
  }
  if (state === 'error') {
    return 'border-red-200 bg-red-50 text-red-500'
  }
  if (state === 'loading') {
    return 'border-amber-200 bg-amber-50 text-amber-600'
  }
  return 'border-slate-200 bg-slate-50 text-slate-500'
}

function widgetCardClass(enabled: boolean) {
  return enabled
    ? 'border-miku/25 bg-[linear-gradient(135deg,rgba(57,197,187,0.12),rgba(255,255,255,0.95))] text-slate-700'
    : 'border-slate-200 bg-white/75 text-slate-500'
}

async function probeGitHubStatus() {
  const username = trimText(integrationsGithubUsername.value)
  if (!username) {
    githubProbe.value = {
      state: 'idle',
      message: copy.preview.statusIdle,
      detail: copy.forms.siteIntegrations.githubMissingHint,
    }
    return
  }

  githubProbe.value = {
    state: 'loading',
    message: copy.preview.statusLoading,
    detail: username,
  }

  try {
    const data = await api.get<GitHubProbeResponse>(`/github/profile?username=${encodeURIComponent(username)}`)
    const name = trimText(data.profile?.name || '') || username
    const totalRepos = Number(data.profile?.total_repos ?? 0)
    githubProbe.value = {
      state: 'ok',
      message: name,
      detail: totalRepos > 0 ? `${totalRepos} 个公开仓库已可读取` : 'GitHub 资料读取正常',
    }
  } catch (err) {
    const detail = err instanceof ApiError ? err.message : copy.preview.statusError
    githubProbe.value = {
      state: 'error',
      message: copy.preview.statusError,
      detail,
    }
  }
}

async function probeWeatherStatus() {
  const location = trimText(integrationsWeatherLocation.value)
  weatherProbe.value = {
    state: 'loading',
    message: copy.preview.statusLoading,
    detail: location || copy.forms.siteIntegrations.weatherDefaultHint,
  }

  try {
    const query = location ? `?location=${encodeURIComponent(location)}` : ''
    const data = await api.get<WeatherProbeResponse>(`/weather${query}`)
    const displayLocation = trimText(data.location || '') || location || 'Default'
    const temp = trimText(data.temp || '')
    weatherProbe.value = {
      state: 'ok',
      message: temp ? `${displayLocation} · ${temp}°C` : displayLocation,
      detail: trimText(data.desc || '') || (location ? '天气服务读取正常' : '已使用后端默认地点'),
    }
  } catch (err) {
    const detail = err instanceof ApiError ? err.message : copy.preview.statusError
    weatherProbe.value = {
      state: 'error',
      message: copy.preview.statusError,
      detail,
    }
  }
}

async function refreshIntegrationStatus() {
  if (refreshingIntegrationStatus.value) {
    return
  }

  refreshingIntegrationStatus.value = true
  try {
    await Promise.all([
      probeGitHubStatus(),
      probeWeatherStatus(),
    ])
  } finally {
    refreshingIntegrationStatus.value = false
  }
}

async function load() {
  loading.value = true
  try {
    await Promise.all([
      hydrateHomeHeroSettings(),
      hydrateHomeAssetsSettings(),
      hydrateSiteProfileSettings(),
      hydrateAuthorProfileSettings(),
      hydrateSiteIntegrationsSettings(),
      hydrateSiteFooterSettings(),
      loadAdminProfile(),
    ])
    syncHomeHeroFromStore()
    syncHomeAssetsFromStore()
    syncSiteProfileFromStore()
    syncAuthorProfileFromStore()
    syncSiteIntegrationsFromStore()
    syncFooterFromStore()
    commitSnapshots()
  } finally {
    loading.value = false
  }

  await refreshIntegrationStatus()
}

const homeHeroDirty = computed(() => serializeHomeHero(getHomeHeroDraft()) !== homeHeroSnapshot.value)
const homeAssetsDirty = computed(() => serializeHomeAssets(getHomeAssetsDraft()) !== homeAssetsSnapshot.value)
const siteProfileDirty = computed(() => serializeSiteProfile(getSiteProfileDraft()) !== siteProfileSnapshot.value)
const authorProfileDirty = computed(() => serializeAuthorProfile(getAuthorProfileDraft()) !== authorProfileSnapshot.value)
const siteIntegrationsDirty = computed(() => serializeSiteIntegrations(getSiteIntegrationsDraft()) !== siteIntegrationsSnapshot.value)
const siteFooterDirty = computed(() => serializeSiteFooter(getSiteFooterDraft()) !== siteFooterSnapshot.value)
const adminProfileDirty = computed(() => serializeAdminProfile() !== adminProfileSnapshot.value)

const sections = computed(() => ([
  {
    key: 'site-profile' as SectionKey,
    ...copy.sections.siteProfile,
    dirty: siteProfileDirty.value,
  },
  {
    key: 'home-hero' as SectionKey,
    ...copy.sections.homeHero,
    dirty: homeHeroDirty.value,
  },
  {
    key: 'home-assets' as SectionKey,
    ...copy.sections.homeAssets,
    dirty: homeAssetsDirty.value,
  },
  {
    key: 'author-profile' as SectionKey,
    ...copy.sections.authorProfile,
    dirty: authorProfileDirty.value,
  },
  {
    key: 'site-integrations' as SectionKey,
    ...copy.sections.siteIntegrations,
    dirty: siteIntegrationsDirty.value,
  },
  {
    key: 'site-footer' as SectionKey,
    ...copy.sections.siteFooter,
    dirty: siteFooterDirty.value,
  },
  {
    key: 'admin-profile' as SectionKey,
    ...copy.sections.adminProfile,
    dirty: adminProfileDirty.value,
  },
]))

const currentSection = computed(() => {
  return sections.value.find((section) => section.key === activeSection.value) ?? sections.value[0]
})

const currentSectionDirty = computed(() => {
  if (activeSection.value === 'site-profile') {
    return siteProfileDirty.value
  }
  if (activeSection.value === 'home-hero') {
    return homeHeroDirty.value
  }
  if (activeSection.value === 'home-assets') {
    return homeAssetsDirty.value
  }
  if (activeSection.value === 'author-profile') {
    return authorProfileDirty.value
  }
  if (activeSection.value === 'site-integrations') {
    return siteIntegrationsDirty.value
  }
  if (activeSection.value === 'site-footer') {
    return siteFooterDirty.value
  }
  return adminProfileDirty.value
})

const canResetCurrentSection = computed(() => activeSection.value !== 'admin-profile')

const isCurrentSaving = computed(() => {
  if (activeSection.value === 'site-profile') {
    return savingSiteProfile.value
  }
  if (activeSection.value === 'home-hero') {
    return savingHomeHero.value
  }
  if (activeSection.value === 'home-assets') {
    return savingHomeAssets.value
  }
  if (activeSection.value === 'author-profile') {
    return savingAuthorProfile.value
  }
  if (activeSection.value === 'site-integrations') {
    return savingSiteIntegrations.value
  }
  if (activeSection.value === 'site-footer') {
    return savingFooter.value
  }
  return savingAdminProfile.value
})

const statusText = computed(() => {
  if (loading.value) {
    return copy.actions.statusLoading
  }
  if (currentSectionDirty.value) {
    return copy.actions.statusPending
  }
  if (lastSavedAt.value && lastSavedSection.value) {
    return `${copy.actions.statusSavedPrefix} · ${lastSavedSection.value} · ${lastSavedAt.value}`
  }
  return copy.actions.statusIdle
})

const statusClass = computed(() => {
  if (loading.value) {
    return 'text-slate-500'
  }
  if (currentSectionDirty.value) {
    return 'text-amber-600'
  }
  if (lastSavedAt.value) {
    return 'text-miku'
  }
  return 'text-slate-500'
})

const sitePreviewTitle = computed(() => {
  const titles = Array.from(new Set([
    trimText(siteTitle.value),
    trimText(brandText.value),
  ].filter((item) => item.length > 0)))

  return titles.length > 0 ? titles.join(' · ') : copy.actions.emptyValue
})

const sitePreviewURL = computed(() => trimText(siteUrl.value) || copy.actions.emptyValue)
const sitePreviewDescription = computed(() => trimText(defaultDescription.value) || copy.actions.emptyValue)
const socialImagePreview = computed(() => trimText(defaultSocialImage.value) || '')

const homeHeroTitlePreview = computed(() => trimText(homeHeroTitle.value) || copy.actions.emptyValue)
const homeHeroSubtitlePreview = computed(() => trimText(homeHeroSubtitle.value) || copy.actions.emptyValue)
const homeHeroImagesPreview = computed(() => {
  const images = cleanHeroImages(homeHeroImages.value)
  return images.length > 0 ? images : homeAssetsStore.value.heroImages
})
const homeAssetsCountPreview = computed(() => {
  return `${copy.preview.heroImageCountPrefix}${homeHeroImagesPreview.value.length}${copy.preview.heroImageCountSuffix}`
})

const authorAvatarPreview = computed(() => trimText(authorAvatarUrl.value) || '/picture/author.jpg')
const authorDisplayNamePreview = computed(() => trimText(authorDisplayName.value) || copy.actions.emptyValue)
const authorRolePreview = computed(() => trimText(authorRole.value) || copy.actions.emptyValue)
const authorMetaPreview = computed(() => {
  const parts = [
    trimText(authorLocation.value),
    trimText(authorSince.value),
    trimText(authorContactEmail.value),
  ].filter((item) => item.length > 0)

  return parts.length > 0 ? parts.join(' · ') : copy.actions.emptyValue
})
const authorBioPreview = computed(() => trimText(authorBio.value) || copy.actions.emptyValue)
const authorAboutPreview = computed(() => trimText(authorAboutDescription.value) || copy.actions.emptyValue)
const authorQuotePreview = computed(() => trimText(authorQuote.value) || copy.actions.emptyValue)
const authorSkillsPreview = computed(() => {
  const skills = cleanTextItems(authorSkills.value, maxAuthorSkills)
  return skills.length > 0 ? skills : [copy.actions.emptyValue]
})
const authorNowItemsPreview = computed(() => {
  const items = cleanTextItems(authorNowItems.value, maxAuthorNowItems)
  return items.length > 0 ? items : [copy.actions.emptyValue]
})
const authorSocialLinksPreview = computed(() => {
  return cleanAuthorSocialDraftLinks(authorSocialLinks.value)
})

const integrationStatusCards = computed(() => ([
  {
    title: copy.preview.githubStatusTitle,
    ...githubProbe.value,
  },
  {
    title: copy.preview.weatherStatusTitle,
    ...weatherProbe.value,
  },
]))

const widgetCards = computed(() => ([
  {
    key: 'weather',
    label: copy.forms.siteIntegrations.showWeatherLabel,
    enabled: integrationsShowWeather.value,
    hint: copy.forms.siteIntegrations.showWeatherHint,
  },
  {
    key: 'music',
    label: copy.forms.siteIntegrations.showMusicLabel,
    enabled: integrationsShowMusic.value,
    hint: copy.forms.siteIntegrations.showMusicHint,
  },
  {
    key: 'clock',
    label: copy.forms.siteIntegrations.showClockLabel,
    enabled: integrationsShowClock.value,
    hint: copy.forms.siteIntegrations.showClockHint,
  },
]))

const footerBrandPreview = computed(() => {
  return trimText(siteTitle.value)
    || trimText(brandText.value)
    || trimText(siteProfileStore.value.siteTitle)
    || trimText(siteProfileStore.value.brandText)
    || copy.actions.emptyValue
})

const footerICPPreview = computed(() => trimText(icpText.value) || copy.actions.emptyValue)
const footerLinesPreview = computed(() => {
  const lines = cleanFooterLines(customTexts.value)
  return lines.length > 0 ? lines : [copy.actions.emptyValue]
})

const profileAvatarPreview = computed(() => trimText(profileAvatarURL.value) || '/picture/author.jpg')
const profileDisplayNamePreview = computed(() => {
  return trimText(profileDisplayName.value)
    || auth.value.user?.name
    || auth.value.user?.username
    || copy.actions.emptyValue
})

const accountIdentityPreview = computed(() => {
  const email = trimText(accountEmail.value)
  const username = trimText(accountUsername.value)
  if (email && username) {
    return `${username} · ${email}`
  }
  return email || username || copy.actions.emptyValue
})

function validateSiteProfileSection() {
  return [
    brandText.value,
    siteTitle.value,
    siteUrl.value,
    defaultDescription.value,
    defaultSocialImage.value,
  ].every((item) => trimText(item).length > 0)
}

function validateHomeHeroSection() {
  return [
    homeHeroTitle.value,
    homeHeroSubtitle.value,
  ].every((item) => trimText(item).length > 0)
}

function validateHomeAssetsSection() {
  return cleanHeroImages(homeHeroImages.value).length > 0
}

function validateAuthorProfileSection() {
  return [
    authorDisplayName.value,
    authorRole.value,
    authorBio.value,
    authorAboutDescription.value,
  ].every((item) => trimText(item).length > 0)
}

async function saveSiteProfileSection() {
  if (!validateSiteProfileSection()) {
    showToast(footerCopy.siteProfile.validationRequired, 'error')
    return
  }

  savingSiteProfile.value = true
  try {
    await saveSiteProfileSettings(getSiteProfileDraft())
    syncSiteProfileFromStore()
    siteProfileSnapshot.value = serializeSiteProfile(getSiteProfileDraft())
    markSaved(copy.sections.siteProfile.title)
    showToast(footerCopy.siteProfileToast.saveSuccess, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : footerCopy.siteProfileToast.saveFailed
    showToast(message, 'error')
  } finally {
    savingSiteProfile.value = false
  }
}

async function resetSiteProfileSection() {
  savingSiteProfile.value = true
  try {
    await resetSiteProfileSettings()
    syncSiteProfileFromStore()
    siteProfileSnapshot.value = serializeSiteProfile(getSiteProfileDraft())
    markSaved(copy.sections.siteProfile.title)
    showToast(footerCopy.siteProfileToast.resetSuccess, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : footerCopy.siteProfileToast.saveFailed
    showToast(message, 'error')
  } finally {
    savingSiteProfile.value = false
  }
}

async function saveHomeHeroSection() {
  if (!validateHomeHeroSection()) {
    showToast(copy.forms.homeHero.validationRequired, 'error')
    return
  }

  savingHomeHero.value = true
  try {
    await saveHomeHeroSettings(getHomeHeroDraft())
    syncHomeHeroFromStore()
    homeHeroSnapshot.value = serializeHomeHero(getHomeHeroDraft())
    markSaved(copy.sections.homeHero.title)
    showToast(copy.toasts.homeHeroSaved, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.toasts.homeHeroFailed
    showToast(message, 'error')
  } finally {
    savingHomeHero.value = false
  }
}

async function resetHomeHeroSection() {
  savingHomeHero.value = true
  try {
    await resetHomeHeroSettings()
    syncHomeHeroFromStore()
    homeHeroSnapshot.value = serializeHomeHero(getHomeHeroDraft())
    markSaved(copy.sections.homeHero.title)
    showToast(copy.toasts.homeHeroReset, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.toasts.homeHeroFailed
    showToast(message, 'error')
  } finally {
    savingHomeHero.value = false
  }
}

async function saveHomeAssetsSection() {
  if (!validateHomeAssetsSection()) {
    showToast(copy.forms.homeAssets.validationRequired, 'error')
    return
  }

  savingHomeAssets.value = true
  try {
    await saveHomeAssetsSettings(getHomeAssetsDraft())
    syncHomeAssetsFromStore()
    homeAssetsSnapshot.value = serializeHomeAssets(getHomeAssetsDraft())
    markSaved(copy.sections.homeAssets.title)
    showToast(copy.toasts.homeAssetsSaved, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.toasts.homeAssetsFailed
    showToast(message, 'error')
  } finally {
    savingHomeAssets.value = false
  }
}

async function resetHomeAssetsSection() {
  savingHomeAssets.value = true
  try {
    await resetHomeAssetsSettings()
    syncHomeAssetsFromStore()
    homeAssetsSnapshot.value = serializeHomeAssets(getHomeAssetsDraft())
    markSaved(copy.sections.homeAssets.title)
    showToast(copy.toasts.homeAssetsReset, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.toasts.homeAssetsFailed
    showToast(message, 'error')
  } finally {
    savingHomeAssets.value = false
  }
}

async function saveAuthorProfileSection() {
  if (!validateAuthorProfileSection()) {
    showToast(copy.forms.authorProfile.validationRequired, 'error')
    return
  }

  savingAuthorProfile.value = true
  try {
    await saveAuthorProfileSettings(getAuthorProfileDraft())
    syncAuthorProfileFromStore()
    authorProfileSnapshot.value = serializeAuthorProfile(getAuthorProfileDraft())
    markSaved(copy.sections.authorProfile.title)
    showToast(copy.toasts.authorProfileSaved, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.toasts.authorProfileFailed
    showToast(message, 'error')
  } finally {
    savingAuthorProfile.value = false
  }
}

async function resetAuthorProfileSection() {
  savingAuthorProfile.value = true
  try {
    await resetAuthorProfileSettings()
    syncAuthorProfileFromStore()
    authorProfileSnapshot.value = serializeAuthorProfile(getAuthorProfileDraft())
    markSaved(copy.sections.authorProfile.title)
    showToast(copy.toasts.authorProfileReset, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.toasts.authorProfileFailed
    showToast(message, 'error')
  } finally {
    savingAuthorProfile.value = false
  }
}

async function saveSiteIntegrationsSection() {
  savingSiteIntegrations.value = true
  try {
    await saveSiteIntegrationsSettings(getSiteIntegrationsDraft())
    syncSiteIntegrationsFromStore()
    siteIntegrationsSnapshot.value = serializeSiteIntegrations(getSiteIntegrationsDraft())
    markSaved(copy.sections.siteIntegrations.title)
    showToast(copy.toasts.siteIntegrationsSaved, 'success')
    await refreshIntegrationStatus()
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.toasts.siteIntegrationsFailed
    showToast(message, 'error')
  } finally {
    savingSiteIntegrations.value = false
  }
}

async function resetSiteIntegrationsSection() {
  savingSiteIntegrations.value = true
  try {
    await resetSiteIntegrationsSettings()
    syncSiteIntegrationsFromStore()
    siteIntegrationsSnapshot.value = serializeSiteIntegrations(getSiteIntegrationsDraft())
    markSaved(copy.sections.siteIntegrations.title)
    showToast(copy.toasts.siteIntegrationsReset, 'success')
    await refreshIntegrationStatus()
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.toasts.siteIntegrationsFailed
    showToast(message, 'error')
  } finally {
    savingSiteIntegrations.value = false
  }
}

async function saveSiteFooterSection() {
  savingFooter.value = true
  try {
    await saveSiteFooterSettings(getSiteFooterDraft())
    syncFooterFromStore()
    siteFooterSnapshot.value = serializeSiteFooter(getSiteFooterDraft())
    markSaved(copy.sections.siteFooter.title)
    showToast(footerCopy.toast.saveSuccess, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : footerCopy.toast.saveFailed
    showToast(message, 'error')
  } finally {
    savingFooter.value = false
  }
}

async function resetSiteFooterSection() {
  savingFooter.value = true
  try {
    await resetSiteFooterSettings()
    syncFooterFromStore()
    siteFooterSnapshot.value = serializeSiteFooter(getSiteFooterDraft())
    markSaved(copy.sections.siteFooter.title)
    showToast(footerCopy.toast.resetSuccess, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : footerCopy.toast.saveFailed
    showToast(message, 'error')
  } finally {
    savingFooter.value = false
  }
}

async function saveAdminProfileSection() {
  const draft = getAdminProfileDraft()

  if (!draft.displayName) {
    showToast(profileCopy.profile.emptyNameError, 'error')
    return
  }
  if (!draft.username) {
    showToast(profileCopy.account.emptyUsernameError, 'error')
    return
  }
  if (!draft.email) {
    showToast(profileCopy.account.emptyEmailError, 'error')
    return
  }

  savingAdminProfile.value = true
  try {
    await updateMyProfile(draft.displayName, draft.avatarURL)
    const result = await updateMyAccount(draft.username, draft.email, draft.password)
    accountPassword.value = ''
    adminProfileSnapshot.value = serializeAdminProfile()
    markSaved(copy.sections.adminProfile.title)

    if (result.sessionRevoked) {
      showToast(profileCopy.account.sessionRevokedSuccess, 'success')
      window.setTimeout(() => {
        window.location.replace('/login')
      }, 360)
      return
    }

    await loadAdminProfile()
    adminProfileSnapshot.value = serializeAdminProfile()
    showToast(copy.toasts.adminSectionSaved, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.toasts.adminSectionFailed
    showToast(message, 'error')
  } finally {
    savingAdminProfile.value = false
  }
}

async function saveCurrentSection() {
  if (!currentSectionDirty.value) {
    return
  }

  if (activeSection.value === 'site-profile') {
    await saveSiteProfileSection()
    return
  }

  if (activeSection.value === 'home-hero') {
    await saveHomeHeroSection()
    return
  }

  if (activeSection.value === 'home-assets') {
    await saveHomeAssetsSection()
    return
  }

  if (activeSection.value === 'author-profile') {
    await saveAuthorProfileSection()
    return
  }

  if (activeSection.value === 'site-integrations') {
    await saveSiteIntegrationsSection()
    return
  }

  if (activeSection.value === 'site-footer') {
    await saveSiteFooterSection()
    return
  }

  await saveAdminProfileSection()
}

async function resetCurrentSection() {
  if (activeSection.value === 'site-profile') {
    await resetSiteProfileSection()
    return
  }

  if (activeSection.value === 'home-hero') {
    await resetHomeHeroSection()
    return
  }

  if (activeSection.value === 'home-assets') {
    await resetHomeAssetsSection()
    return
  }

  if (activeSection.value === 'author-profile') {
    await resetAuthorProfileSection()
    return
  }

  if (activeSection.value === 'site-integrations') {
    await resetSiteIntegrationsSection()
    return
  }

  if (activeSection.value === 'site-footer') {
    await resetSiteFooterSection()
  }
}

onMounted(() => {
  void load()
})
</script>
