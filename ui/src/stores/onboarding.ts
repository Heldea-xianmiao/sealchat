import { defineStore } from 'pinia'

// ========== 类型定义 ==========

/** 教程模块分类 */
export type TutorialCategory = 'basic' | 'social' | 'advanced' | 'management'

/** 单个引导步骤 */
export interface TutorialStep {
    id: string
    title: string
    content: string
    /** CSS 选择器，用于定位目标元素 */
    target?: string
    /** 提示卡片相对目标的位置 */
    placement?: 'top' | 'bottom' | 'left' | 'right' | 'center'
    /** 是否高亮目标元素 */
    highlight?: boolean
    /** 可选的示意图路径 */
    image?: string
}

/** 单个教程模块 */
export interface TutorialModule {
    id: string
    title: string
    description: string
    category: TutorialCategory
    steps: TutorialStep[]
    /** 可选权限要求（仅对有权限用户显示） */
    requiresPermission?: string
    /** 预计时长（秒） */
    estimatedTime: number
}

/** 用户引导进度 */
export interface OnboardingProgress {
    /** 已完成的模块 ID */
    completedModules: string[]
    /** 当前正在进行的模块 ID */
    currentModule: string | null
    /** 当前步骤索引 */
    currentStep: number
    /** 开始时间戳 */
    startedAt: number | null
    /** 最后活跃时间戳 */
    lastActiveAt: number | null
}

interface OnboardingState {
    /** 是否已完成首次引导（旧用户跳过） */
    hasCompletedInitialOnboarding: boolean
    /** 用户引导进度 */
    progress: OnboardingProgress
    /** 当前是否正在显示引导 */
    isActive: boolean
    /** 是否显示教程中心 */
    showTutorialHub: boolean
    /** 选中的模块（批量学习） */
    selectedModules: string[]
    /** 模块列表（运行时注入） */
    _modules: TutorialModule[]
}

const STORAGE_KEY = 'sealchat_onboarding'

/** 创建默认进度对象 */
function createDefaultProgress(): OnboardingProgress {
    return {
        completedModules: [],
        currentModule: null,
        currentStep: 0,
        startedAt: null,
        lastActiveAt: null,
    }
}

export const useOnboardingStore = defineStore('onboarding', {
    state: (): OnboardingState => ({
        hasCompletedInitialOnboarding: false,
        progress: createDefaultProgress(),
        isActive: false,
        showTutorialHub: false,
        selectedModules: [],
        _modules: [],
    }),

    getters: {
        /** 所有教程模块 */
        modules(): TutorialModule[] {
            return this._modules
        },

        /** 是否为新用户（需要显示引导） */
        isNewUser(): boolean {
            return !this.hasCompletedInitialOnboarding && this.progress.completedModules.length === 0
        },

        /** 是否有未完成的引导 */
        hasIncompleteOnboarding(): boolean {
            return this.progress.currentModule !== null && this.progress.startedAt !== null
        },

        /** 当前模块对象 */
        currentModuleConfig(): TutorialModule | null {
            if (!this.progress.currentModule) return null
            return this._modules.find((m) => m.id === this.progress.currentModule) || null
        },

        /** 当前步骤对象 */
        currentStepConfig(): TutorialStep | null {
            const module = this.currentModuleConfig
            if (!module) return null
            return module.steps[this.progress.currentStep] || null
        },

        /** 当前步骤是否为模块的最后一步 */
        isLastStep(): boolean {
            const module = this.currentModuleConfig
            if (!module) return false
            return this.progress.currentStep >= module.steps.length - 1
        },

        /** 当前步骤是否为模块的第一步 */
        isFirstStep(): boolean {
            return this.progress.currentStep === 0
        },

        /** 模块总数（排除需要权限的模块） */
        totalModulesCount(): number {
            return this._modules.filter((m) => !m.requiresPermission).length
        },

        /** 模块完成百分比 */
        completionPercentage(): number {
            const total = this.totalModulesCount
            if (total === 0) return 0
            return Math.round((this.progress.completedModules.length / total) * 100)
        },

        /** 按分类分组的模块 */
        modulesByCategory(): Record<TutorialCategory, TutorialModule[]> {
            const groups: Record<TutorialCategory, TutorialModule[]> = {
                basic: [],
                social: [],
                advanced: [],
                management: [],
            }
            for (const mod of this._modules) {
                groups[mod.category].push(mod)
            }
            return groups
        },

        /** 是否显示欢迎页 */
        showWelcome(): boolean {
            return this.isActive && this.progress.currentModule === '__welcome'
        },

        /** 是否显示聚光灯引导 */
        showSpotlight(): boolean {
            return (
                this.isActive &&
                this.progress.currentModule !== null &&
                this.progress.currentModule !== '__welcome' &&
                !this.showTutorialHub
            )
        },

        /** 是否显示继续提示 */
        showResumePrompt(): boolean {
            return this.isActive && this.hasIncompleteOnboarding && !this.showSpotlight && !this.showWelcome
        },
    },

    actions: {
        /** 注册模块列表 */
        registerModules(modules: TutorialModule[]) {
            this._modules = modules
        },

        /** 从 localStorage 加载状态 */
        loadFromStorage() {
            try {
                const raw = localStorage.getItem(STORAGE_KEY)
                if (raw) {
                    const data = JSON.parse(raw)
                    this.hasCompletedInitialOnboarding = data.hasCompletedInitialOnboarding ?? false
                    this.progress = {
                        ...createDefaultProgress(),
                        ...data.progress,
                    }
                }
            } catch (e) {
                console.warn('[Onboarding] Failed to load state from storage', e)
            }
        },

        /** 保存状态到 localStorage */
        saveToStorage() {
            try {
                localStorage.setItem(
                    STORAGE_KEY,
                    JSON.stringify({
                        hasCompletedInitialOnboarding: this.hasCompletedInitialOnboarding,
                        progress: this.progress,
                    })
                )
            } catch (e) {
                console.warn('[Onboarding] Failed to save state to storage', e)
            }
        },

        /** 检查并启动引导（登录后调用） */
        checkAndStartOnboarding() {
            this.loadFromStorage()

            // 新用户：显示欢迎页
            if (this.isNewUser) {
                this.showWelcomePage()
                return
            }

            // 有未完成的引导：显示继续提示
            if (this.hasIncompleteOnboarding) {
                this.isActive = true
            }
        },

        /** 显示欢迎页 */
        showWelcomePage() {
            this.isActive = true
            this.showTutorialHub = false
            this.progress.currentModule = '__welcome'
            this.progress.currentStep = 0
        },

        /** 显示教程中心 */
        openTutorialHub() {
            this.isActive = true
            this.showTutorialHub = true
            this.progress.currentModule = null
        },

        /** 关闭教程中心 */
        closeTutorialHub() {
            this.showTutorialHub = false
            if (!this.progress.currentModule || this.progress.currentModule === '__welcome') {
                this.isActive = false
            }
        },

        /** 开始指定模块 */
        startModule(moduleId: string) {
            this.progress.currentModule = moduleId
            this.progress.currentStep = 0
            this.progress.startedAt = Date.now()
            this.progress.lastActiveAt = Date.now()
            this.showTutorialHub = false
            this.isActive = true
            this.saveToStorage()
        },

        /** 开始选中的模块（批量） */
        startSelectedModules() {
            if (this.selectedModules.length === 0) return
            // 按模块在列表中的原始顺序排序
            const orderedIds = this._modules.map((m) => m.id).filter((id) => this.selectedModules.includes(id))
            this.selectedModules = orderedIds
            this.startModule(orderedIds[0])
        },

        /** 下一步 */
        nextStep() {
            const module = this.currentModuleConfig
            if (!module) return

            if (this.progress.currentStep < module.steps.length - 1) {
                this.progress.currentStep++
                this.progress.lastActiveAt = Date.now()
                this.saveToStorage()
            } else {
                this.completeCurrentModule()
            }
        },

        /** 上一步 */
        prevStep() {
            if (this.progress.currentStep > 0) {
                this.progress.currentStep--
                this.progress.lastActiveAt = Date.now()
                this.saveToStorage()
            }
        },

        /** 完成当前模块 */
        completeCurrentModule() {
            const moduleId = this.progress.currentModule
            if (moduleId && moduleId !== '__welcome') {
                if (!this.progress.completedModules.includes(moduleId)) {
                    this.progress.completedModules.push(moduleId)
                }
            }

            // 检查是否有下一个选中的模块
            const currentIdx = this.selectedModules.indexOf(moduleId || '')
            if (currentIdx >= 0 && currentIdx < this.selectedModules.length - 1) {
                this.startModule(this.selectedModules[currentIdx + 1])
            } else {
                this.finish()
            }
        },

        /** 完成引导 */
        finish() {
            this.hasCompletedInitialOnboarding = true
            this.progress.currentModule = null
            this.progress.currentStep = 0
            this.progress.startedAt = null
            this.isActive = false
            this.selectedModules = []
            this.saveToStorage()
        },

        /** 跳过引导 */
        skip() {
            this.hasCompletedInitialOnboarding = true
            this.progress.currentModule = null
            this.progress.startedAt = null
            this.isActive = false
            this.showTutorialHub = false
            this.saveToStorage()
        },

        /** 手动重新开始（打开教程中心） */
        restart() {
            this.openTutorialHub()
        },

        /** 继续未完成的引导 */
        resumeOnboarding() {
            if (this.progress.currentModule && this.progress.currentModule !== '__welcome') {
                this.showTutorialHub = false
                this.isActive = true
            }
        },

        /** 重新开始当前模块 */
        restartCurrentModule() {
            if (this.progress.currentModule) {
                this.progress.currentStep = 0
                this.progress.lastActiveAt = Date.now()
                this.saveToStorage()
            }
        },

        /** 检查模块是否已完成 */
        isModuleCompleted(moduleId: string): boolean {
            return this.progress.completedModules.includes(moduleId)
        },

        /** 切换模块选中状态 */
        toggleModuleSelection(moduleId: string) {
            const idx = this.selectedModules.indexOf(moduleId)
            if (idx >= 0) {
                this.selectedModules.splice(idx, 1)
            } else {
                this.selectedModules.push(moduleId)
            }
        },

        /** 选中推荐模块 */
        selectRecommendedModules(recommendedIds: string[]) {
            this.selectedModules = [...recommendedIds]
        },

        /** 选中所有模块 */
        selectAllModules() {
            this.selectedModules = this._modules.filter((m) => !m.requiresPermission).map((m) => m.id)
        },

        /** 清空选中 */
        clearSelection() {
            this.selectedModules = []
        },
    },
})
