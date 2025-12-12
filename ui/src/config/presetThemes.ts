/**
 * 预设主题配置 (Preset Theme Configurations)
 * 
 * 这个文件包含了Sealchat的预设自定义主题配置。
 * 用户可以通过自定义主题面板快速导入这些预设。
 */

import type { CustomThemeColors } from '@/stores/display'

export interface PresetTheme {
    id: string
    name: string
    description: string
    colors: CustomThemeColors
}

/**
 * 豆沙绿护眼主题 (Dusha Green Eye-Care Theme)
 * 
 * 参考菠萝平台的护眼豆沙绿颜色设计
 * 使用低饱和度暖绿色系，减少视觉疲劳
 */
export const dushaGreenTheme: PresetTheme = {
    id: 'preset_dusha_green',
    name: '豆沙绿护眼',
    description: '低饱和度暖绿色系，减少视觉疲劳，适合长时间阅读',
    colors: {
        // 背景色
        bgSurface: '#f6f9f4',      // 主背景 - 最浅豆沙绿
        bgElevated: '#fcfdf9',     // 卡片/弹窗背景
        bgInput: '#edf3e8',        // 输入框背景
        bgHeader: '#f8faf6',       // 顶栏背景
        // 文字色
        textPrimary: '#2a3b21',    // 主文字 - 深绿
        textSecondary: '#557542',  // 次要文字
        // 聊天区域
        chatIcBg: '#f4f8f0',       // 场内消息背景
        chatOocBg: '#fcfdf9',      // 场外消息背景
        chatStageBg: '#edf3e8',    // 聊天舞台背景
        chatPreviewBg: '#dce8d4',  // 预览背景
        chatPreviewDot: '#c5d9b8', // 预览圆点
        // 边框色
        borderMute: '#dce8d4',     // 淡边框
        borderStrong: '#a8c494',   // 强边框
        // 强调色
        primaryColor: '#6d9255',   // 主题色 - 深绿
        primaryColorHover: '#8bae72', // 悬停色
        // 术语高亮
        keywordBg: 'rgba(139, 174, 114, 0.2)',  // 高亮背景
        keywordBorder: '#8bae72',              // 下划线色
    },
}

/**
 * 所有预设主题列表
 */
export const presetThemes: PresetTheme[] = [
    dushaGreenTheme,
]

/**
 * 根据ID获取预设主题
 */
export const getPresetThemeById = (id: string): PresetTheme | undefined => {
    return presetThemes.find(theme => theme.id === id)
}
