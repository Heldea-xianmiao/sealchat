/**
 * 消息链接渲染工具
 * 将消息链接转换为 Discord 风格的跳转标记
 */

import { parseMessageLink } from './messageLink'

export interface MessageLinkRenderInfo {
  url: string
  worldId: string
  channelId: string
  messageId: string
  worldName: string
  channelName: string
  isCurrentWorld: boolean
  customTitle?: string
}

/**
 * 解析并获取链接的渲染信息
 */
export function resolveMessageLinkInfo(
  url: string,
  context: {
    currentWorldId: string
    worldMap: Record<string, { name?: string }>
    findChannelById: (id: string) => { name?: string } | null
  },
  customTitle?: string
): MessageLinkRenderInfo | null {
  const params = parseMessageLink(url)
  if (!params) return null

  const { worldId, channelId, messageId } = params
  const isCurrentWorld = worldId === context.currentWorldId

  // 获取世界名称
  let worldName = '未知世界'
  const worldInfo = context.worldMap[worldId]
  if (worldInfo?.name) {
    worldName = worldInfo.name
  }

  // 获取频道名称
  let channelName = '消息'
  if (isCurrentWorld) {
    const channelInfo = context.findChannelById(channelId)
    if (channelInfo?.name) {
      channelName = channelInfo.name
    }
  }

  return {
    url,
    worldId,
    channelId,
    messageId,
    worldName,
    channelName,
    isCurrentWorld,
    customTitle,
  }
}

/**
 * 生成链接的 HTML 显示内容
 * 自定义标题: #自定义标题 › 📝
 * 本世界: #频道名 › 📝
 * 跨世界: #世界名 › 📝
 */
export function renderMessageLinkHtml(info: MessageLinkRenderInfo): string {
  const displayName = info.customTitle || (info.isCurrentWorld ? info.channelName : info.worldName)
  // 使用简单的消息图标 SVG
  const icon = `<svg class="msg-link-icon" viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><path d="M20 2H4c-1.1 0-2 .9-2 2v18l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H6l-2 2V4h16v12z"/></svg>`

  return `<a href="${escapeHtml(info.url)}" class="message-jump-link" data-world-id="${escapeHtml(info.worldId)}" data-channel-id="${escapeHtml(info.channelId)}" data-message-id="${escapeHtml(info.messageId)}" data-is-current-world="${info.isCurrentWorld}"><span class="message-jump-link__hash">#</span><span class="message-jump-link__name">${escapeHtml(displayName)}</span><span class="message-jump-link__separator">›</span>${icon}</a>`
}

/**
 * HTML 转义
 */
function escapeHtml(text: string): string {
  const map: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;',
  }
  return text.replace(/[&<>"']/g, (char) => map[char] || char)
}
