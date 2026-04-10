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

    <div class="grid gap-5 xl:grid-cols-[minmax(0,1.15fr)_340px]">
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

          <p class="mt-3 rounded-[16px] border border-slate-200/80 bg-slate-50/80 px-4 py-3 text-xs leading-6 text-slate-500">
            {{ copy.forms.authorProfile.githubIntegrationHint }}
          </p>

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

        <div v-else-if="activeSection === 'about-page'" :class="surfaceClass">
          <div>
            <p class="text-sm font-semibold text-slate-900">{{ copy.forms.aboutPage.title }}</p>
            <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.aboutPage.subtitle }}</p>
          </div>

          <div class="mt-6 flex flex-wrap items-center justify-between gap-3">
            <div>
              <div class="flex flex-wrap items-center gap-2">
                <p class="text-sm font-semibold text-slate-900">{{ copy.forms.aboutPage.introCardsTitle }}</p>
                <span class="rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] font-medium text-slate-500">
                  {{ aboutIntroCards.length }} / {{ maxAboutIntroCards }}
                </span>
              </div>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.aboutPage.introCardsSubtitle }}</p>
            </div>
            <button
              type="button"
              class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
              :aria-label="copy.forms.aboutPage.addIntroButtonAria"
              :disabled="loading || savingAboutPage || aboutIntroCards.length >= maxAboutIntroCards"
              @click="addAboutIntroCard"
            >
              {{ copy.forms.aboutPage.addIntroButton }}
            </button>
          </div>

          <div class="mt-4 space-y-3">
            <div
              v-for="(card, index) in aboutIntroCards"
              :key="`about-intro-${index}`"
              class="rounded-[20px] border border-slate-200/80 bg-white/75 p-4"
            >
              <div class="grid gap-3 md:grid-cols-[minmax(0,1fr)_auto]">
                <input
                  v-model="aboutIntroCards[index].title"
                  type="text"
                  :placeholder="copy.forms.aboutPage.introTitlePlaceholder"
                  :class="inputClass"
                  :aria-label="`${copy.forms.aboutPage.introCardsTitle}${index + 1}`"
                  :disabled="loading || savingAboutPage"
                />
                <button
                  type="button"
                  class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                  :aria-label="copy.forms.aboutPage.removeButtonAria"
                  :disabled="loading || savingAboutPage"
                  @click="removeAboutIntroCard(index)"
                >
                  {{ copy.forms.aboutPage.removeButton }}
                </button>
              </div>
              <textarea
                v-model="aboutIntroCards[index].description"
                rows="3"
                :placeholder="copy.forms.aboutPage.introDescriptionPlaceholder"
                :class="`${textareaClass} mt-3`"
                :aria-label="`${copy.forms.aboutPage.introCardsTitle}${index + 1} description`"
                :disabled="loading || savingAboutPage"
              />
            </div>

            <p v-if="aboutIntroCards.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
              {{ copy.forms.aboutPage.introEmptyHint }}
            </p>
          </div>

          <div class="mt-8 flex flex-wrap items-center justify-between gap-3">
            <div>
              <div class="flex flex-wrap items-center gap-2">
                <p class="text-sm font-semibold text-slate-900">{{ copy.forms.aboutPage.milestonesTitle }}</p>
                <span class="rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] font-medium text-slate-500">
                  {{ aboutMilestones.length }} / {{ maxAboutMilestones }}
                </span>
              </div>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.aboutPage.milestonesSubtitle }}</p>
            </div>
            <button
              type="button"
              class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
              :aria-label="copy.forms.aboutPage.addMilestoneButtonAria"
              :disabled="loading || savingAboutPage || aboutMilestones.length >= maxAboutMilestones"
              @click="addAboutMilestone"
            >
              {{ copy.forms.aboutPage.addMilestoneButton }}
            </button>
          </div>

          <div class="mt-4 space-y-3">
            <div
              v-for="(item, index) in aboutMilestones"
              :key="`about-milestone-${index}`"
              class="rounded-[20px] border border-slate-200/80 bg-white/75 p-4"
            >
              <div class="grid gap-3 md:grid-cols-[160px_minmax(0,1fr)_auto]">
                <input
                  v-model="aboutMilestones[index].year"
                  type="text"
                  :placeholder="copy.forms.aboutPage.milestoneYearPlaceholder"
                  :class="inputClass"
                  :aria-label="`${copy.forms.aboutPage.milestonesTitle}${index + 1} year`"
                  :disabled="loading || savingAboutPage"
                />
                <input
                  v-model="aboutMilestones[index].title"
                  type="text"
                  :placeholder="copy.forms.aboutPage.milestoneTitlePlaceholder"
                  :class="inputClass"
                  :aria-label="`${copy.forms.aboutPage.milestonesTitle}${index + 1} title`"
                  :disabled="loading || savingAboutPage"
                />
                <button
                  type="button"
                  class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                  :aria-label="copy.forms.aboutPage.removeButtonAria"
                  :disabled="loading || savingAboutPage"
                  @click="removeAboutMilestone(index)"
                >
                  {{ copy.forms.aboutPage.removeButton }}
                </button>
              </div>
              <textarea
                v-model="aboutMilestones[index].summary"
                rows="3"
                :placeholder="copy.forms.aboutPage.milestoneSummaryPlaceholder"
                :class="`${textareaClass} mt-3`"
                :aria-label="`${copy.forms.aboutPage.milestonesTitle}${index + 1} summary`"
                :disabled="loading || savingAboutPage"
              />
              <input
                v-model="aboutMilestones[index].result"
                type="text"
                :placeholder="copy.forms.aboutPage.milestoneResultPlaceholder"
                :class="`${inputClass} mt-3`"
                :aria-label="`${copy.forms.aboutPage.milestonesTitle}${index + 1} result`"
                :disabled="loading || savingAboutPage"
              />
            </div>

            <p v-if="aboutMilestones.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
              {{ copy.forms.aboutPage.milestoneEmptyHint }}
            </p>
          </div>

          <div class="mt-8 flex flex-wrap items-center justify-between gap-3">
            <div>
              <div class="flex flex-wrap items-center gap-2">
                <p class="text-sm font-semibold text-slate-900">{{ copy.forms.aboutPage.capabilityGroupsTitle }}</p>
                <span class="rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] font-medium text-slate-500">
                  {{ aboutCapabilityGroups.length }} / {{ maxAboutCapabilityGroups }}
                </span>
              </div>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.aboutPage.capabilityGroupsSubtitle }}</p>
            </div>
            <button
              type="button"
              class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
              :aria-label="copy.forms.aboutPage.addCapabilityButtonAria"
              :disabled="loading || savingAboutPage || aboutCapabilityGroups.length >= maxAboutCapabilityGroups"
              @click="addAboutCapabilityGroup"
            >
              {{ copy.forms.aboutPage.addCapabilityButton }}
            </button>
          </div>

          <div class="mt-4 space-y-3">
            <div
              v-for="(group, index) in aboutCapabilityGroups"
              :key="`about-capability-${index}`"
              class="rounded-[20px] border border-slate-200/80 bg-white/75 p-4"
            >
              <div class="grid gap-3 md:grid-cols-[minmax(0,1fr)_auto]">
                <input
                  v-model="aboutCapabilityGroups[index].title"
                  type="text"
                  :placeholder="copy.forms.aboutPage.capabilityTitlePlaceholder"
                  :class="inputClass"
                  :aria-label="`${copy.forms.aboutPage.capabilityGroupsTitle}${index + 1} title`"
                  :disabled="loading || savingAboutPage"
                />
                <button
                  type="button"
                  class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                  :aria-label="copy.forms.aboutPage.removeButtonAria"
                  :disabled="loading || savingAboutPage"
                  @click="removeAboutCapabilityGroup(index)"
                >
                  {{ copy.forms.aboutPage.removeButton }}
                </button>
              </div>
              <textarea
                v-model="aboutCapabilityGroups[index].desc"
                rows="3"
                :placeholder="copy.forms.aboutPage.capabilityDescPlaceholder"
                :class="`${textareaClass} mt-3`"
                :aria-label="`${copy.forms.aboutPage.capabilityGroupsTitle}${index + 1} description`"
                :disabled="loading || savingAboutPage"
              />
              <input
                v-model="aboutCapabilityGroups[index].stackText"
                type="text"
                :placeholder="copy.forms.aboutPage.capabilityStackPlaceholder"
                :class="`${inputClass} mt-3`"
                :aria-label="`${copy.forms.aboutPage.capabilityGroupsTitle}${index + 1} stack`"
                :disabled="loading || savingAboutPage"
              />
            </div>

            <p v-if="aboutCapabilityGroups.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
              {{ copy.forms.aboutPage.capabilityEmptyHint }}
            </p>
          </div>

          <div class="mt-8 flex flex-wrap items-center justify-between gap-3">
            <div>
              <div class="flex flex-wrap items-center gap-2">
                <p class="text-sm font-semibold text-slate-900">{{ copy.forms.aboutPage.featuredProjectsTitle }}</p>
                <span class="rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] font-medium text-slate-500">
                  {{ aboutFeaturedProjects.length }} / {{ maxAboutFeaturedProjects }}
                </span>
              </div>
              <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.aboutPage.featuredProjectsSubtitle }}</p>
            </div>
            <button
              type="button"
              class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
              :aria-label="copy.forms.aboutPage.addProjectButtonAria"
              :disabled="loading || savingAboutPage || aboutFeaturedProjects.length >= maxAboutFeaturedProjects"
              @click="addAboutFeaturedProject"
            >
              {{ copy.forms.aboutPage.addProjectButton }}
            </button>
          </div>

          <div class="mt-4 space-y-3">
            <div
              v-for="(project, index) in aboutFeaturedProjects"
              :key="`about-project-${index}`"
              class="rounded-[20px] border border-slate-200/80 bg-white/75 p-4"
            >
              <div class="grid gap-3 md:grid-cols-[minmax(0,1fr)_auto]">
                <input
                  v-model="aboutFeaturedProjects[index].name"
                  type="text"
                  :placeholder="copy.forms.aboutPage.projectNamePlaceholder"
                  :class="inputClass"
                  :aria-label="`${copy.forms.aboutPage.featuredProjectsTitle}${index + 1} name`"
                  :disabled="loading || savingAboutPage"
                />
                <button
                  type="button"
                  class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                  :aria-label="copy.forms.aboutPage.removeButtonAria"
                  :disabled="loading || savingAboutPage"
                  @click="removeAboutFeaturedProject(index)"
                >
                  {{ copy.forms.aboutPage.removeButton }}
                </button>
              </div>
              <input
                v-model="aboutFeaturedProjects[index].focus"
                type="text"
                :placeholder="copy.forms.aboutPage.projectFocusPlaceholder"
                :class="`${inputClass} mt-3`"
                :aria-label="`${copy.forms.aboutPage.featuredProjectsTitle}${index + 1} focus`"
                :disabled="loading || savingAboutPage"
              />
              <textarea
                v-model="aboutFeaturedProjects[index].role"
                rows="3"
                :placeholder="copy.forms.aboutPage.projectRolePlaceholder"
                :class="`${textareaClass} mt-3`"
                :aria-label="`${copy.forms.aboutPage.featuredProjectsTitle}${index + 1} role`"
                :disabled="loading || savingAboutPage"
              />
              <input
                v-model="aboutFeaturedProjects[index].metric"
                type="text"
                :placeholder="copy.forms.aboutPage.projectMetricPlaceholder"
                :class="`${inputClass} mt-3`"
                :aria-label="`${copy.forms.aboutPage.featuredProjectsTitle}${index + 1} metric`"
                :disabled="loading || savingAboutPage"
              />
              <input
                v-model="aboutFeaturedProjects[index].href"
                type="text"
                :placeholder="copy.forms.aboutPage.projectHrefPlaceholder"
                :class="`${inputClass} mt-3`"
                :aria-label="`${copy.forms.aboutPage.featuredProjectsTitle}${index + 1} href`"
                :disabled="loading || savingAboutPage"
              />
            </div>

            <p v-if="aboutFeaturedProjects.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
              {{ copy.forms.aboutPage.projectEmptyHint }}
            </p>
          </div>

          <div class="mt-8 grid gap-6 xl:grid-cols-2">
            <div>
              <div class="flex flex-wrap items-center justify-between gap-3">
                <div>
                  <div class="flex flex-wrap items-center gap-2">
                    <p class="text-sm font-semibold text-slate-900">{{ copy.forms.aboutPage.monthlyGoalsTitle }}</p>
                    <span class="rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] font-medium text-slate-500">
                      {{ aboutMonthlyGoals.length }} / {{ maxAboutMonthlyGoals }}
                    </span>
                  </div>
                  <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.aboutPage.monthlyGoalsSubtitle }}</p>
                </div>
                <button
                  type="button"
                  class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
                  :aria-label="copy.forms.aboutPage.addMonthlyGoalButtonAria"
                  :disabled="loading || savingAboutPage || aboutMonthlyGoals.length >= maxAboutMonthlyGoals"
                  @click="addAboutMonthlyGoal"
                >
                  {{ copy.forms.aboutPage.addMonthlyGoalButton }}
                </button>
              </div>

              <div class="mt-4 space-y-2.5">
                <div
                  v-for="(goal, index) in aboutMonthlyGoals"
                  :key="`about-goal-${index}`"
                  class="grid gap-2 md:grid-cols-[1fr_auto]"
                >
                  <input
                    v-model="aboutMonthlyGoals[index]"
                    type="text"
                    :placeholder="copy.forms.aboutPage.monthlyGoalPlaceholder"
                    :class="inputClass"
                    :aria-label="`${copy.forms.aboutPage.monthlyGoalsTitle}${index + 1}`"
                    :disabled="loading || savingAboutPage"
                  />
                  <button
                    type="button"
                    class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                    :aria-label="copy.forms.aboutPage.removeButtonAria"
                    :disabled="loading || savingAboutPage"
                    @click="removeAboutMonthlyGoal(index)"
                  >
                    {{ copy.forms.aboutPage.removeButton }}
                  </button>
                </div>

                <p v-if="aboutMonthlyGoals.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
                  {{ copy.forms.aboutPage.monthlyGoalsEmptyHint }}
                </p>
              </div>
            </div>

            <div>
              <div class="flex flex-wrap items-center justify-between gap-3">
                <div>
                  <div class="flex flex-wrap items-center gap-2">
                    <p class="text-sm font-semibold text-slate-900">{{ copy.forms.aboutPage.listeningTitle }}</p>
                    <span class="rounded-full border border-slate-200 bg-white px-2.5 py-1 text-[11px] font-medium text-slate-500">
                      {{ aboutListeningNow.length }} / {{ maxAboutListeningNow }}
                    </span>
                  </div>
                  <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.aboutPage.listeningSubtitle }}</p>
                </div>
                <button
                  type="button"
                  class="rounded-full border border-miku/25 bg-miku/10 px-3.5 py-2 text-xs font-medium text-miku transition hover:bg-miku/15 disabled:cursor-not-allowed disabled:opacity-45"
                  :aria-label="copy.forms.aboutPage.addListeningButtonAria"
                  :disabled="loading || savingAboutPage || aboutListeningNow.length >= maxAboutListeningNow"
                  @click="addAboutListeningItem"
                >
                  {{ copy.forms.aboutPage.addListeningButton }}
                </button>
              </div>

              <div class="mt-4 space-y-2.5">
                <div
                  v-for="(item, index) in aboutListeningNow"
                  :key="`about-listening-${index}`"
                  class="grid gap-2 md:grid-cols-[1fr_auto]"
                >
                  <input
                    v-model="aboutListeningNow[index]"
                    type="text"
                    :placeholder="copy.forms.aboutPage.listeningPlaceholder"
                    :class="inputClass"
                    :aria-label="`${copy.forms.aboutPage.listeningTitle}${index + 1}`"
                    :disabled="loading || savingAboutPage"
                  />
                  <button
                    type="button"
                    class="rounded-[16px] border border-slate-200 bg-white/85 px-3 py-2 text-sm text-slate-500 transition hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-45"
                    :aria-label="copy.forms.aboutPage.removeButtonAria"
                    :disabled="loading || savingAboutPage"
                    @click="removeAboutListeningItem(index)"
                  >
                    {{ copy.forms.aboutPage.removeButton }}
                  </button>
                </div>

                <p v-if="aboutListeningNow.length === 0" class="rounded-[18px] border border-dashed border-slate-200 bg-white/60 px-4 py-4 text-xs leading-6 text-slate-500">
                  {{ copy.forms.aboutPage.listeningEmptyHint }}
                </p>
              </div>
            </div>
          </div>

          <div class="mt-8">
            <p class="text-sm font-semibold text-slate-900">{{ copy.forms.aboutPage.signatureTitle }}</p>
            <p class="mt-1 text-xs leading-6 text-slate-500">{{ copy.forms.aboutPage.signatureSubtitle }}</p>
          </div>

          <div class="mt-4 grid gap-3">
            <textarea
              v-model="aboutSignatureDescription"
              rows="4"
              :placeholder="copy.forms.aboutPage.signatureDescriptionPlaceholder"
              :class="textareaClass"
              :aria-label="copy.forms.aboutPage.signatureTitle"
              :disabled="loading || savingAboutPage"
            />
            <input
              v-model="aboutSignatureFooter"
              type="text"
              :placeholder="copy.forms.aboutPage.signatureFooterPlaceholder"
              :class="inputClass"
              :aria-label="`${copy.forms.aboutPage.signatureTitle} footer`"
              :disabled="loading || savingAboutPage"
            />
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

          <div v-else-if="activeSection === 'about-page'" class="mt-5 space-y-4">
            <div class="rounded-[20px] border border-slate-200/80 bg-[linear-gradient(145deg,rgba(255,255,255,0.92),rgba(229,255,252,0.9),rgba(248,244,255,0.86))] p-4">
              <div class="grid gap-3 sm:grid-cols-2">
                <div class="rounded-[16px] border border-white/80 bg-white/78 px-3 py-3">
                  <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutIntroLabel }}</p>
                  <p class="mt-2 text-lg font-semibold text-slate-900">{{ aboutPagePreview.introCards.length }}</p>
                </div>
                <div class="rounded-[16px] border border-white/80 bg-white/78 px-3 py-3">
                  <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutTimelineLabel }}</p>
                  <p class="mt-2 text-lg font-semibold text-slate-900">{{ aboutPagePreview.milestones.length }}</p>
                </div>
                <div class="rounded-[16px] border border-white/80 bg-white/78 px-3 py-3">
                  <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutCapabilitiesLabel }}</p>
                  <p class="mt-2 text-lg font-semibold text-slate-900">{{ aboutPagePreview.capabilityGroups.length }}</p>
                </div>
                <div class="rounded-[16px] border border-white/80 bg-white/78 px-3 py-3">
                  <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutProjectsLabel }}</p>
                  <p class="mt-2 text-lg font-semibold text-slate-900">{{ aboutPagePreview.featuredProjects.length }}</p>
                </div>
                <div class="rounded-[16px] border border-white/80 bg-white/78 px-3 py-3">
                  <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutGoalsLabel }}</p>
                  <p class="mt-2 text-lg font-semibold text-slate-900">{{ aboutPagePreview.monthlyGoals.length }}</p>
                </div>
                <div class="rounded-[16px] border border-white/80 bg-white/78 px-3 py-3">
                  <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutListeningLabel }}</p>
                  <p class="mt-2 text-lg font-semibold text-slate-900">{{ aboutPagePreview.listeningNow.length }}</p>
                </div>
              </div>
            </div>

            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutIntroLabel }}</p>
              <div class="mt-3 space-y-3">
                <div
                  v-for="card in aboutPagePreview.introCards.slice(0, 2)"
                  :key="`${card.title}-${card.description}`"
                  class="rounded-[16px] border border-slate-200 bg-white px-3 py-3"
                >
                  <p class="text-sm font-semibold text-slate-900">{{ card.title }}</p>
                  <p class="mt-1 text-xs leading-6 text-slate-500">{{ card.description }}</p>
                </div>
                <p
                  v-if="aboutPagePreview.introCards.length === 0"
                  class="rounded-[16px] border border-dashed border-slate-200 bg-slate-50/70 px-3 py-4 text-xs leading-6 text-slate-500"
                >
                  {{ copy.actions.emptyValue }}
                </p>
              </div>
            </div>

            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutTimelineLabel }}</p>
              <div class="mt-3 space-y-3">
                <div
                  v-for="item in aboutPagePreview.milestones.slice(0, 2)"
                  :key="`${item.year}-${item.title}`"
                  class="relative overflow-hidden rounded-[16px] border border-slate-200 bg-white px-3 py-3"
                >
                  <span class="absolute left-0 top-0 h-full w-1 bg-[linear-gradient(180deg,rgba(57,197,187,0.85),rgba(192,132,252,0.75))]" />
                  <div class="pl-2">
                    <div class="flex flex-wrap items-center justify-between gap-2">
                      <p class="text-sm font-semibold text-slate-900">{{ item.title }}</p>
                      <span class="rounded-full border border-miku/20 bg-miku/10 px-2 py-0.5 text-[11px] font-medium text-miku">
                        {{ item.year }}
                      </span>
                    </div>
                    <p class="mt-2 text-xs leading-6 text-slate-500">{{ item.summary }}</p>
                    <p class="mt-2 text-[11px] text-slate-400">{{ item.result }}</p>
                  </div>
                </div>
                <p
                  v-if="aboutPagePreview.milestones.length === 0"
                  class="rounded-[16px] border border-dashed border-slate-200 bg-slate-50/70 px-3 py-4 text-xs leading-6 text-slate-500"
                >
                  {{ copy.actions.emptyValue }}
                </p>
              </div>
            </div>

            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutCapabilitiesLabel }}</p>
              <div class="mt-3 space-y-3">
                <div
                  v-for="(group, index) in aboutPagePreview.capabilityGroups.slice(0, 2)"
                  :key="`${group.title}-${group.desc}`"
                  class="rounded-[16px] border border-slate-200 bg-white px-3 py-3"
                >
                  <p class="text-sm font-semibold text-slate-900">{{ group.title }}</p>
                  <p class="mt-1 text-xs leading-6 text-slate-500">{{ group.desc }}</p>
                  <div class="mt-3 flex flex-wrap gap-1.5">
                    <span
                      v-for="item in group.stack"
                      :key="`${group.title}-${item}`"
                      class="rounded-full border px-2.5 py-1 text-[11px] font-medium"
                      :class="aboutPreviewCapabilityChipClass(index)"
                    >
                      {{ item }}
                    </span>
                  </div>
                </div>
                <p
                  v-if="aboutPagePreview.capabilityGroups.length === 0"
                  class="rounded-[16px] border border-dashed border-slate-200 bg-slate-50/70 px-3 py-4 text-xs leading-6 text-slate-500"
                >
                  {{ copy.actions.emptyValue }}
                </p>
              </div>
            </div>

            <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
              <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutProjectsLabel }}</p>
              <div class="mt-3 space-y-3">
                <div
                  v-for="project in aboutPagePreview.featuredProjects.slice(0, 2)"
                  :key="`${project.name}-${project.href}`"
                  class="rounded-[16px] border border-slate-200 bg-white px-3 py-3"
                >
                  <p class="text-sm font-semibold text-slate-900">{{ project.name }}</p>
                  <p class="mt-1 text-xs text-slate-500">{{ project.focus }}</p>
                  <p class="mt-2 text-xs leading-6 text-slate-500">{{ project.role }}</p>
                  <p class="mt-2 rounded-[12px] border border-slate-200 bg-slate-50/80 px-2.5 py-2 text-[11px] leading-5 text-slate-500">
                    {{ project.metric }}
                  </p>
                  <p class="mt-2 truncate text-[11px] text-miku">{{ project.href }}</p>
                </div>
                <p
                  v-if="aboutPagePreview.featuredProjects.length === 0"
                  class="rounded-[16px] border border-dashed border-slate-200 bg-slate-50/70 px-3 py-4 text-xs leading-6 text-slate-500"
                >
                  {{ copy.actions.emptyValue }}
                </p>
              </div>
            </div>

            <div class="grid gap-4">
              <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
                <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutGoalsLabel }}</p>
                <div class="mt-3 flex flex-wrap gap-2">
                  <span
                    v-for="goal in aboutPagePreview.monthlyGoals"
                    :key="goal"
                    class="rounded-full border border-slate-200 bg-white px-3 py-1 text-xs font-medium text-slate-600"
                  >
                    {{ goal }}
                  </span>
                  <span
                    v-if="aboutPagePreview.monthlyGoals.length === 0"
                    class="rounded-full border border-dashed border-slate-200 bg-slate-50/70 px-3 py-1 text-xs font-medium text-slate-500"
                  >
                    {{ copy.actions.emptyValue }}
                  </span>
                </div>

                <p class="mt-4 text-xs font-medium text-slate-500">{{ copy.preview.aboutListeningLabel }}</p>
                <div class="mt-3 flex flex-wrap gap-2">
                  <span
                    v-for="item in aboutPagePreview.listeningNow"
                    :key="item"
                    class="rounded-full border border-miku/20 bg-miku/10 px-3 py-1 text-xs font-medium text-miku"
                  >
                    {{ item }}
                  </span>
                  <span
                    v-if="aboutPagePreview.listeningNow.length === 0"
                    class="rounded-full border border-dashed border-slate-200 bg-slate-50/70 px-3 py-1 text-xs font-medium text-slate-500"
                  >
                    {{ copy.actions.emptyValue }}
                  </span>
                </div>
              </div>

              <div class="rounded-[20px] border border-slate-200/80 bg-white/80 p-4">
                <p class="text-xs font-medium text-slate-500">{{ copy.preview.aboutSignatureLabel }}</p>
                <blockquote class="mt-2 rounded-[16px] border border-slate-200 bg-slate-50/80 px-3 py-3 text-xs leading-6 text-slate-500">
                  {{ aboutPagePreview.signature.description || copy.actions.emptyValue }}
                  <p class="mt-2 text-[11px] text-slate-400">
                    {{ aboutPagePreview.signature.footer || copy.actions.emptyValue }}
                  </p>
                </blockquote>
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
import { computed, onMounted, ref, watch } from 'vue'

import { adminCopy } from '../../content/copy'
import {
  normalizeAboutPageSettings,
  type AboutCapabilityGroup,
  type AboutFeaturedProject,
  type AboutIntroCard,
  type AboutMilestone,
  type AboutPageSettings,
} from '../../lib/about-page'
import { api, ApiError } from '../../lib/api'
import { DEFAULT_PUBLIC_AUTHOR_AVATAR_URL, DEFAULT_PUBLIC_AVATAR_URL } from '../../lib/default-assets'
import {
  resolveAuthorAvatarURL,
} from '../../lib/author-profile'
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
  aboutPageSettings,
  hydrateAboutPageSettings,
  resetAboutPageSettings,
  saveAboutPageSettings,
} from '../../stores/aboutPage'
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
import {
  readAdminSettingsSectionFromLocation,
  setAdminSettingsSection,
  type AdminSettingsSectionKey,
} from '../../stores/adminSettings'
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

type SectionKey = AdminSettingsSectionKey

type ProbeState = 'idle' | 'loading' | 'ok' | 'error'

interface ProbeStatus {
  state: ProbeState
  message: string
  detail: string
}

interface AboutCapabilityGroupDraft {
  title: string
  desc: string
  stackText: string
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
const maxAboutIntroCards = 6
const maxAboutMilestones = 8
const maxAboutCapabilityGroups = 6
const maxAboutFeaturedProjects = 6
const maxAboutMonthlyGoals = 8
const maxAboutListeningNow = 8

const homeHeroStore = useStore(homeHeroSettings)
const homeAssetsStore = useStore(homeAssetsSettings)
const siteProfileStore = useStore(siteProfileSettings)
const siteFooterStore = useStore(siteFooterSettings)
const authorProfileStore = useStore(authorProfileSettings)
const aboutPageStore = useStore(aboutPageSettings)
const siteIntegrationsStore = useStore(siteIntegrationsSettings)
const auth = useStore(authState)

const activeSection = ref<SectionKey>('site-profile')
const loading = ref(false)
const savingHomeHero = ref(false)
const savingHomeAssets = ref(false)
const savingSiteProfile = ref(false)
const savingAuthorProfile = ref(false)
const savingAboutPage = ref(false)
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
const authorAvatarUrl = ref(DEFAULT_PUBLIC_AUTHOR_AVATAR_URL)
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

const aboutIntroCards = ref<AboutIntroCard[]>([])
const aboutMilestones = ref<AboutMilestone[]>([])
const aboutCapabilityGroups = ref<AboutCapabilityGroupDraft[]>([])
const aboutFeaturedProjects = ref<AboutFeaturedProject[]>([])
const aboutMonthlyGoals = ref<string[]>([])
const aboutListeningNow = ref<string[]>([])
const aboutSignatureDescription = ref('')
const aboutSignatureFooter = ref('')

const integrationsGithubUsername = ref('')
const integrationsWeatherLocation = ref('')
const integrationsShowWeather = ref(true)
const integrationsShowMusic = ref(true)
const integrationsShowClock = ref(true)

const icpText = ref('')
const icpLink = ref('')
const customTexts = ref<string[]>([])

const profileDisplayName = ref('')
const profileAvatarURL = ref(DEFAULT_PUBLIC_AVATAR_URL)
const accountUsername = ref('')
const accountEmail = ref('')
const accountPassword = ref('')

const homeHeroSnapshot = ref('')
const homeAssetsSnapshot = ref('')
const siteProfileSnapshot = ref('')
const authorProfileSnapshot = ref('')
const aboutPageSnapshot = ref('')
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

function splitStackText(value: string) {
  return value
    .split(/[\n,，]+/g)
    .map((item) => trimText(item))
    .filter((item, index, list) => item.length > 0 && list.indexOf(item) === index)
}

function cleanAboutIntroDraftCards(items: AboutIntroCard[]) {
  const result: AboutIntroCard[] = []
  for (const item of items) {
    const title = trimText(item.title)
    const description = trimText(item.description)
    if (!title || !description) {
      continue
    }
    if (result.some((current) => current.title === title && current.description === description)) {
      continue
    }
    result.push({ title, description })
    if (result.length >= maxAboutIntroCards) {
      break
    }
  }
  return result
}

function cleanAboutMilestonesDraft(items: AboutMilestone[]) {
  const result: AboutMilestone[] = []
  for (const item of items) {
    const year = trimText(item.year)
    const title = trimText(item.title)
    const summary = trimText(item.summary)
    const resultText = trimText(item.result)
    if (!year || !title || !summary || !resultText) {
      continue
    }
    if (result.some((current) => current.year === year && current.title === title)) {
      continue
    }
    result.push({
      year,
      title,
      summary,
      result: resultText,
    })
    if (result.length >= maxAboutMilestones) {
      break
    }
  }
  return result
}

function cleanAboutCapabilityDraftGroups(items: AboutCapabilityGroupDraft[]) {
  const result: AboutCapabilityGroup[] = []
  for (const item of items) {
    const title = trimText(item.title)
    const desc = trimText(item.desc)
    const stack = splitStackText(item.stackText)
    if (!title || !desc || stack.length === 0) {
      continue
    }
    if (result.some((current) => current.title === title && current.desc === desc)) {
      continue
    }
    result.push({
      title,
      desc,
      stack,
    })
    if (result.length >= maxAboutCapabilityGroups) {
      break
    }
  }
  return result
}

function cleanAboutFeaturedProjectsDraft(items: AboutFeaturedProject[]) {
  const result: AboutFeaturedProject[] = []
  for (const item of items) {
    const name = trimText(item.name)
    const focus = trimText(item.focus)
    const role = trimText(item.role)
    const metric = trimText(item.metric)
    const href = trimText(item.href)
    if (!name || !focus || !role || !metric || !href) {
      continue
    }
    if (result.some((current) => current.name === name && current.href === href)) {
      continue
    }
    result.push({
      name,
      focus,
      role,
      metric,
      href,
    })
    if (result.length >= maxAboutFeaturedProjects) {
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

function getAboutPageDraft(): AboutPageSettings {
  return {
    introCards: cleanAboutIntroDraftCards(aboutIntroCards.value),
    milestones: cleanAboutMilestonesDraft(aboutMilestones.value),
    capabilityGroups: cleanAboutCapabilityDraftGroups(aboutCapabilityGroups.value),
    featuredProjects: cleanAboutFeaturedProjectsDraft(aboutFeaturedProjects.value),
    monthlyGoals: cleanTextItems(aboutMonthlyGoals.value, maxAboutMonthlyGoals),
    listeningNow: cleanTextItems(aboutListeningNow.value, maxAboutListeningNow),
    signature: {
      description: trimText(aboutSignatureDescription.value),
      footer: trimText(aboutSignatureFooter.value),
    },
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

function serializeAboutPage(settings: AboutPageSettings) {
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

function syncAboutPageFromStore() {
  aboutIntroCards.value = aboutPageStore.value.introCards.map((item) => ({ ...item }))
  aboutMilestones.value = aboutPageStore.value.milestones.map((item) => ({ ...item }))
  aboutCapabilityGroups.value = aboutPageStore.value.capabilityGroups.map((item) => ({
    title: item.title,
    desc: item.desc,
    stackText: item.stack.join(', '),
  }))
  aboutFeaturedProjects.value = aboutPageStore.value.featuredProjects.map((item) => ({ ...item }))
  aboutMonthlyGoals.value = [...aboutPageStore.value.monthlyGoals]
  aboutListeningNow.value = [...aboutPageStore.value.listeningNow]
  aboutSignatureDescription.value = aboutPageStore.value.signature.description
  aboutSignatureFooter.value = aboutPageStore.value.signature.footer
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
    profileAvatarURL.value = current.avatar || DEFAULT_PUBLIC_AVATAR_URL
    accountUsername.value = current.username || ''
    accountEmail.value = current.email || ''
  }

  try {
    const me = await api.get<AdminProfilePayload>('/auth/me')
    profileDisplayName.value = trimText(me.display_name || '') || me.username
    profileAvatarURL.value = trimText(me.avatar_url || '') || DEFAULT_PUBLIC_AVATAR_URL
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
  aboutPageSnapshot.value = serializeAboutPage(getAboutPageDraft())
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

function addAboutIntroCard() {
  if (aboutIntroCards.value.length >= maxAboutIntroCards) {
    return
  }
  aboutIntroCards.value.push({ title: '', description: '' })
}

function addAboutMilestone() {
  if (aboutMilestones.value.length >= maxAboutMilestones) {
    return
  }
  aboutMilestones.value.push({ year: '', title: '', summary: '', result: '' })
}

function addAboutCapabilityGroup() {
  if (aboutCapabilityGroups.value.length >= maxAboutCapabilityGroups) {
    return
  }
  aboutCapabilityGroups.value.push({ title: '', desc: '', stackText: '' })
}

function addAboutFeaturedProject() {
  if (aboutFeaturedProjects.value.length >= maxAboutFeaturedProjects) {
    return
  }
  aboutFeaturedProjects.value.push({ name: '', focus: '', role: '', metric: '', href: '' })
}

function addAboutMonthlyGoal() {
  if (aboutMonthlyGoals.value.length >= maxAboutMonthlyGoals) {
    return
  }
  aboutMonthlyGoals.value.push('')
}

function addAboutListeningItem() {
  if (aboutListeningNow.value.length >= maxAboutListeningNow) {
    return
  }
  aboutListeningNow.value.push('')
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

function removeAboutIntroCard(index: number) {
  aboutIntroCards.value = aboutIntroCards.value.filter((_, currentIndex) => currentIndex !== index)
}

function removeAboutMilestone(index: number) {
  aboutMilestones.value = aboutMilestones.value.filter((_, currentIndex) => currentIndex !== index)
}

function removeAboutCapabilityGroup(index: number) {
  aboutCapabilityGroups.value = aboutCapabilityGroups.value.filter((_, currentIndex) => currentIndex !== index)
}

function removeAboutFeaturedProject(index: number) {
  aboutFeaturedProjects.value = aboutFeaturedProjects.value.filter((_, currentIndex) => currentIndex !== index)
}

function removeAboutMonthlyGoal(index: number) {
  aboutMonthlyGoals.value = aboutMonthlyGoals.value.filter((_, currentIndex) => currentIndex !== index)
}

function removeAboutListeningItem(index: number) {
  aboutListeningNow.value = aboutListeningNow.value.filter((_, currentIndex) => currentIndex !== index)
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

function aboutPreviewCapabilityChipClass(index: number) {
  if (index === 0) {
    return 'border-miku/35 bg-miku-soft text-miku'
  }
  if (index === 1) {
    return 'border-slate-200 bg-white/70 text-slate-600'
  }
  return 'border-[#c084fc]/35 bg-[#f3e8ff] text-[#8b5cf6]'
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
      hydrateAboutPageSettings(),
      hydrateSiteIntegrationsSettings(),
      hydrateSiteFooterSettings(),
      loadAdminProfile(),
    ])
    syncHomeHeroFromStore()
    syncHomeAssetsFromStore()
    syncSiteProfileFromStore()
    syncAuthorProfileFromStore()
    syncAboutPageFromStore()
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
const aboutPageDirty = computed(() => serializeAboutPage(getAboutPageDraft()) !== aboutPageSnapshot.value)
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
    key: 'about-page' as SectionKey,
    ...copy.sections.aboutPage,
    dirty: aboutPageDirty.value,
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
  if (activeSection.value === 'about-page') {
    return aboutPageDirty.value
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
  if (activeSection.value === 'about-page') {
    return savingAboutPage.value
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

const authorAvatarPreview = computed(() => resolveAuthorAvatarURL(authorAvatarUrl.value))
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

const aboutPagePreview = computed(() => normalizeAboutPageSettings(getAboutPageDraft()))

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

const profileAvatarPreview = computed(() => trimText(profileAvatarURL.value) || DEFAULT_PUBLIC_AVATAR_URL)
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

function readSectionFromLocation(): SectionKey | null {
  if (typeof window === 'undefined') {
    return null
  }

  return readAdminSettingsSectionFromLocation(window.location.search)
}

function syncSectionToLocation(section: SectionKey) {
  if (typeof window === 'undefined') {
    return
  }

  const url = new URL(window.location.href)
  const current = url.searchParams.get('section')

  if (current === section) {
    return
  }

  if (!current && section === 'site-profile') {
    return
  }

  url.searchParams.set('section', section)
  window.history.replaceState(window.history.state, '', `${url.pathname}${url.search}${url.hash}`)
}

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

async function saveAboutPageSection() {
  savingAboutPage.value = true
  try {
    await saveAboutPageSettings(getAboutPageDraft())
    syncAboutPageFromStore()
    aboutPageSnapshot.value = serializeAboutPage(getAboutPageDraft())
    markSaved(copy.sections.aboutPage.title)
    showToast(copy.toasts.aboutPageSaved, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.toasts.aboutPageFailed
    showToast(message, 'error')
  } finally {
    savingAboutPage.value = false
  }
}

async function resetAboutPageSection() {
  savingAboutPage.value = true
  try {
    await resetAboutPageSettings()
    syncAboutPageFromStore()
    aboutPageSnapshot.value = serializeAboutPage(getAboutPageDraft())
    markSaved(copy.sections.aboutPage.title)
    showToast(copy.toasts.aboutPageReset, 'success')
  } catch (err) {
    const message = err instanceof ApiError ? err.message : copy.toasts.aboutPageFailed
    showToast(message, 'error')
  } finally {
    savingAboutPage.value = false
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

  if (activeSection.value === 'about-page') {
    await saveAboutPageSection()
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

  if (activeSection.value === 'about-page') {
    await resetAboutPageSection()
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

watch(activeSection, (next, prev) => {
  if (next === prev) {
    return
  }
  setAdminSettingsSection(next)
  syncSectionToLocation(next)
})

onMounted(() => {
  const sectionFromLocation = readSectionFromLocation()
  if (sectionFromLocation) {
    activeSection.value = sectionFromLocation
  }
  setAdminSettingsSection(activeSection.value)
  void load()
})
</script>
