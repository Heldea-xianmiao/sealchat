import type { TutorialModule } from '@/stores/onboarding'

/**
 * 所有教程模块定义
 */
export const TUTORIAL_MODULES: TutorialModule[] = [
    // ========== 基础功能 (Basic) ==========
    {
        id: 'chat-basics',
        title: '发送消息',
        description: '学习如何在聊天频道中发送文字消息',
        category: 'basic',
        estimatedTime: 30,
        steps: [
            {
                id: 'chat-basics-1',
                title: '输入消息',
                content: '在底部输入框中输入你想说的话，按 Enter 或点击发送按钮即可发送。',
                target: '.chat-input-main',
                placement: 'top',
                highlight: true,
            },
            {
                id: 'chat-basics-2',
                title: '消息列表',
                content: '发送的消息会显示在中央聊天区域，最新消息显示在底部。',
                target: '.messages-list',
                placement: 'center',
            },
            {
                id: 'chat-basics-3',
                title: '消息操作',
                content: '右键点击消息可以进行回复、引用、编辑、删除等操作。',
                target: '.message-row',
                placement: 'right',
            },
        ],
    },
    {
        id: 'identity-switcher',
        title: '角色切换',
        description: '创建和切换角色身份，自定义头像和颜色',
        category: 'basic',
        estimatedTime: 45,
        steps: [
            {
                id: 'identity-1',
                title: '角色切换器',
                content: '点击左上角的头像和名称区域，可以切换当前使用的角色身份。',
                target: '.identity-switcher',
                placement: 'bottom',
                highlight: true,
            },
            {
                id: 'identity-2',
                title: '创建新角色',
                content: '选择「创建新角色」可以创建拥有独特名称、头像和颜色的身份。',
                placement: 'center',
            },
            {
                id: 'identity-3',
                title: '管理角色',
                content: '选择「管理角色」可以编辑或删除已创建的角色身份。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'ic-ooc-toggle',
        title: 'IC/OOC 模式',
        description: '角色内与角色外模式的区别与切换',
        category: 'basic',
        estimatedTime: 30,
        steps: [
            {
                id: 'ic-ooc-1',
                title: '什么是 IC/OOC？',
                content: 'IC (In Character) 表示「角色内」发言，以你扮演的角色身份说话。\nOOC (Out of Character) 表示「角色外」发言，以玩家身份交流。',
                placement: 'center',
            },
            {
                id: 'ic-ooc-2',
                title: '切换模式',
                content: '点击此开关可以在 IC 和 OOC 模式之间切换。',
                target: '.ic-ooc-toggle',
                placement: 'top',
                highlight: true,
            },
            {
                id: 'ic-ooc-3',
                title: '视觉区分',
                content: 'IC 消息和 OOC 消息在显示上有不同的背景色区分，方便阅读。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'display-settings',
        title: '显示设置',
        description: '主题切换、布局模式、字体大小等自定义',
        category: 'basic',
        estimatedTime: 40,
        steps: [
            {
                id: 'display-1',
                title: '打开显示设置',
                content: '点击工具栏中的「显示设置」按钮可以打开设置面板。',
                target: '[data-tour="display-settings"]',
                placement: 'bottom',
                highlight: true,
            },
            {
                id: 'display-2',
                title: '日夜主题',
                content: '可以在日间和夜间模式之间切换，也可以创建自定义主题。',
                placement: 'center',
            },
            {
                id: 'display-3',
                title: '布局模式',
                content: '气泡模式类似微信，紧凑模式适合阅读大量文字。',
                placement: 'center',
            },
            {
                id: 'display-4',
                title: '排版设置',
                content: '可以调整字体大小、行高、字间距等，打造舒适的阅读体验。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'dice-tray',
        title: '骰子托盘',
        description: '内置骰点功能，支持常见骰子表达式',
        category: 'basic',
        estimatedTime: 35,
        steps: [
            {
                id: 'dice-1',
                title: '打开骰子托盘',
                content: '点击输入区上方的骰子图标可以打开骰子托盘。',
                target: '[data-tour="dice-tray"]',
                placement: 'top',
                highlight: true,
            },
            {
                id: 'dice-2',
                title: '快捷骰子',
                content: '托盘中预置了常用骰子（d4, d6, d8, d10, d12, d20），点击即可投掷。',
                placement: 'center',
            },
            {
                id: 'dice-3',
                title: '自定义表达式',
                content: '可以输入如 2d6+3 的表达式进行自定义骰点。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'emoji-panel',
        title: '表情与收藏',
        description: '收藏表情包、使用画廊功能',
        category: 'basic',
        estimatedTime: 30,
        steps: [
            {
                id: 'emoji-1',
                title: '打开表情面板',
                content: '点击输入框旁的表情图标可以打开表情面板。',
                target: '[data-tour="emoji-panel"]',
                placement: 'top',
                highlight: true,
            },
            {
                id: 'emoji-2',
                title: '添加收藏',
                content: '右键点击消息中的图片可以将其添加到表情收藏。',
                placement: 'center',
            },
            {
                id: 'emoji-3',
                title: '画廊功能',
                content: '画廊可以管理更大的图片集合，适合存放角色立绘等素材。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'channel-favorites',
        title: '频道收藏',
        description: '快捷收藏常用频道',
        category: 'basic',
        estimatedTime: 20,
        steps: [
            {
                id: 'fav-1',
                title: '频道收藏栏',
                content: '屏幕上方的收藏栏可以快速访问常用频道。',
                target: '.channel-favorite-bar',
                placement: 'bottom',
                highlight: true,
            },
            {
                id: 'fav-2',
                title: '管理收藏',
                content: '点击工具栏中的收藏图标可以管理收藏的频道。',
                target: '[data-tour="favorites"]',
                placement: 'bottom',
            },
        ],
    },

    // ========== 社交功能 (Social) ==========
    {
        id: 'world-lobby',
        title: '世界大厅',
        description: '浏览、搜索、加入公开世界',
        category: 'social',
        estimatedTime: 40,
        steps: [
            {
                id: 'world-1',
                title: '进入世界大厅',
                content: '点击侧边栏顶部的世界选择器，可以看到「世界大厅」选项。',
                target: '.world-selector',
                placement: 'right',
                highlight: true,
            },
            {
                id: 'world-2',
                title: '浏览公开世界',
                content: '切换到「发现」标签可以浏览所有公开的世界。',
                placement: 'center',
            },
            {
                id: 'world-3',
                title: '搜索与加入',
                content: '使用搜索框查找感兴趣的世界，点击「加入」即可成为成员。',
                placement: 'center',
            },
            {
                id: 'world-4',
                title: '邀请码加入',
                content: '如果有邀请码，可以输入邀请码直接加入私有世界。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'channel-tree',
        title: '频道导航',
        description: '世界下的频道树结构',
        category: 'social',
        estimatedTime: 25,
        steps: [
            {
                id: 'channel-1',
                title: '频道列表',
                content: '进入世界后，左侧显示该世界下的所有频道。',
                target: '.channel-tree',
                placement: 'right',
                highlight: true,
            },
            {
                id: 'channel-2',
                title: '频道分类',
                content: '频道可以按分类折叠展开，点击频道名称切换到该频道。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'private-chat',
        title: '私聊功能',
        description: '添加好友、发起私聊',
        category: 'social',
        estimatedTime: 30,
        steps: [
            {
                id: 'private-1',
                title: '切换到私聊',
                content: '点击侧边栏的「私聊」标签切换到私聊列表。',
                target: '.private-chat-tab',
                placement: 'right',
                highlight: true,
            },
            {
                id: 'private-2',
                title: '发起私聊',
                content: '点击用户头像或从成员列表中选择用户，可以发起私聊。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'member-list',
        title: '成员列表',
        description: '查看频道/世界成员',
        category: 'social',
        estimatedTime: 20,
        steps: [
            {
                id: 'member-1',
                title: '查看成员',
                content: '点击工具栏中的成员图标可以查看当前频道的成员列表。',
                target: '[data-tour="members"]',
                placement: 'bottom',
                highlight: true,
            },
            {
                id: 'member-2',
                title: '成员操作',
                content: '点击成员可以查看资料、发起私聊等操作。',
                placement: 'center',
            },
        ],
    },

    // ========== 进阶功能 (Advanced) ==========
    {
        id: 'message-search',
        title: '消息搜索',
        description: '搜索历史消息，支持高级筛选',
        category: 'advanced',
        estimatedTime: 40,
        steps: [
            {
                id: 'search-1',
                title: '打开搜索面板',
                content: '点击工具栏中的搜索图标可以打开搜索面板。',
                target: '[data-tour="search"]',
                placement: 'bottom',
                highlight: true,
            },
            {
                id: 'search-2',
                title: '关键词搜索',
                content: '输入关键词即可搜索当前频道的历史消息。',
                placement: 'center',
            },
            {
                id: 'search-3',
                title: '高级筛选',
                content: '展开高级选项可以按时间、发送者、IC/OOC 等条件筛选。',
                placement: 'center',
            },
            {
                id: 'search-4',
                title: '跳转定位',
                content: '点击搜索结果可以直接跳转到该消息在聊天记录中的位置。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'message-archive',
        title: '消息归档',
        description: '归档和管理重要消息',
        category: 'advanced',
        estimatedTime: 30,
        steps: [
            {
                id: 'archive-1',
                title: '打开归档抽屉',
                content: '点击工具栏中的归档图标可以查看已归档的消息。',
                target: '[data-tour="archive"]',
                placement: 'bottom',
                highlight: true,
            },
            {
                id: 'archive-2',
                title: '归档消息',
                content: '右键点击消息选择「归档」可以将消息保存到归档列表。',
                placement: 'center',
            },
            {
                id: 'archive-3',
                title: '管理归档',
                content: '可以搜索、恢复或删除归档的消息。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'message-export',
        title: '日志导出',
        description: '导出聊天记录为多种格式',
        category: 'advanced',
        estimatedTime: 35,
        steps: [
            {
                id: 'export-1',
                title: '打开导出对话框',
                content: '点击工具栏中的导出图标可以导出聊天记录。',
                target: '[data-tour="export"]',
                placement: 'bottom',
                highlight: true,
            },
            {
                id: 'export-2',
                title: '选择格式',
                content: '支持导出为纯文本、HTML、以及海豹染色器格式。',
                placement: 'center',
            },
            {
                id: 'export-3',
                title: '时间范围',
                content: '可以指定导出的时间范围，或导出全部记录。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'message-import',
        title: '日志导入',
        description: '导入外部聊天记录',
        category: 'advanced',
        estimatedTime: 40,
        steps: [
            {
                id: 'import-1',
                title: '打开导入对话框',
                content: '点击工具栏中的导入图标可以导入外部聊天记录。',
                target: '[data-tour="import"]',
                placement: 'bottom',
                highlight: true,
            },
            {
                id: 'import-2',
                title: '选择格式模板',
                content: '选择与源文件匹配的解析模板。',
                placement: 'center',
            },
            {
                id: 'import-3',
                title: '角色映射',
                content: '将日志中的角色名映射到系统中的用户身份。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'keyword-highlight',
        title: '术语高亮',
        description: '世界术语词条、悬浮提示',
        category: 'advanced',
        estimatedTime: 30,
        steps: [
            {
                id: 'keyword-1',
                title: '术语高亮',
                content: '世界管理员定义的术语词条会在消息中自动高亮显示。',
                placement: 'center',
            },
            {
                id: 'keyword-2',
                title: '悬浮提示',
                content: '鼠标悬停在高亮术语上可以查看解释说明。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'shortcuts',
        title: '快捷键管理',
        description: '自定义键盘快捷操作',
        category: 'advanced',
        estimatedTime: 25,
        steps: [
            {
                id: 'shortcuts-1',
                title: '快捷键设置',
                content: '在显示设置中可以自定义工具栏按钮的快捷键。',
                placement: 'center',
            },
            {
                id: 'shortcuts-2',
                title: '常用快捷键',
                content: '可以为显示设置、表情面板、骰子托盘等功能设置快捷键。',
                placement: 'center',
            },
        ],
    },
    {
        id: 'custom-theme',
        title: '自定义主题',
        description: '创建个性化配色方案',
        category: 'advanced',
        estimatedTime: 45,
        steps: [
            {
                id: 'theme-1',
                title: '启用自定义主题',
                content: '在显示设置中开启「自定义主题」开关。',
                placement: 'center',
            },
            {
                id: 'theme-2',
                title: '编辑主题',
                content: '点击编辑按钮可以打开主题编辑器，调整各项颜色。',
                placement: 'center',
            },
            {
                id: 'theme-3',
                title: '导入导出',
                content: '可以导入预设主题或导出自己的主题分享给他人。',
                placement: 'center',
            },
        ],
    },
]

/**
 * 分类信息
 */
export const TUTORIAL_CATEGORIES = [
    { id: 'basic' as const, label: '🗨️ 基础功能', description: '入门必备的核心功能' },
    { id: 'social' as const, label: '👥 社交功能', description: '世界和频道探索' },
    { id: 'advanced' as const, label: '⚙️ 进阶功能', description: '提升效率的高级功能' },
]

/**
 * 推荐给新用户的模块（快速开始/推荐入门）
 * 包含：全部基础功能、全部社交功能、搜索/归档/导出/术语高亮/快捷键
 */
export const RECOMMENDED_MODULES = [
    // 全部基础功能
    'chat-basics',
    'identity-switcher',
    'ic-ooc-toggle',
    'display-settings',
    'dice-tray',
    'emoji-panel',
    'channel-favorites',
    // 全部社交功能
    'world-lobby',
    'channel-tree',
    'private-chat',
    'member-list',
    // 部分进阶功能
    'message-search',
    'message-archive',
    'message-export',
    'keyword-highlight',
    'shortcuts',
]

/**
 * 根据分类 ID 获取分类信息
 */
export function getCategoryInfo(categoryId: string) {
    return TUTORIAL_CATEGORIES.find((c) => c.id === categoryId)
}

/**
 * 根据模块 ID 获取模块
 */
export function getModuleById(moduleId: string) {
    return TUTORIAL_MODULES.find((m) => m.id === moduleId)
}

/**
 * 格式化时长（秒 -> 可读文本）
 */
export function formatDuration(seconds: number): string {
    if (seconds < 60) return `${seconds}秒`
    return `${Math.round(seconds / 60)}分钟`
}
