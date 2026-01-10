/**
 * 消息链接工具函数
 * 用于生成和解析消息跳转链接
 */

export interface MessageLinkParams {
  worldId: string
  channelId: string
  messageId: string
}

/**
 * 生成消息的完整链接
 */
export function generateMessageLink(
  params: MessageLinkParams,
  options?: { base?: string }
): string {
  const { worldId, channelId, messageId } = params
  const base = resolveMessageLinkBase(options?.base)
  return `${base}/#/${worldId}/${channelId}?msg=${messageId}`
}

function resolveMessageLinkBase(base?: string): string {
  const trimmed = (base || '').trim()
  if (trimmed) {
    return trimmed.replace(/\/+$/, '')
  }
  if (typeof window === 'undefined') {
    return ''
  }
  return window.location.origin
}

/**
 * 匹配消息链接路径的正则表达式
 * 格式: /#/{worldId}/{channelId}?msg={messageId}
 * 忽略域名，方便服务器迁移
 * 注意: \? 转义问号
 */
const MESSAGE_LINK_PATH_REGEX = /#\/([a-zA-Z0-9_-]+)\/([a-zA-Z0-9_-]+)\?msg=([a-zA-Z0-9_-]+)/

/**
 * 解析消息链接，返回 worldId, channelId, messageId
 * 仅匹配路径格式，忽略域名
 */
export function parseMessageLink(url: string): MessageLinkParams | null {
  if (!url) return null

  // 直接匹配路径部分 #/{worldId}/{channelId}?msg={messageId}
  const match = url.match(MESSAGE_LINK_PATH_REGEX)
  if (!match) return null

  const [, worldId, channelId, messageId] = match
  if (!worldId || !channelId || !messageId) return null

  return { worldId, channelId, messageId }
}

/**
 * 检查 URL 是否为消息链接格式
 * 仅检查路径格式，不检查域名
 */
export function isLocalMessageLink(url: string): boolean {
  return parseMessageLink(url) !== null
}

/**
 * 消息链接的正则表达式（用于在纯文本中匹配链接）
 * 匹配格式: http(s)://domain/#/{worldId}/{channelId}?msg={messageId}
 */
export const MESSAGE_LINK_REGEX =
  /https?:\/\/[^\s<>"]*#\/[a-zA-Z0-9_-]+\/[a-zA-Z0-9_-]+\?msg=[a-zA-Z0-9_-]+/g

/**
 * 带自定义标题的消息链接正则表达式
 * 匹配格式: [自定义标题](http(s)://domain/#/{worldId}/{channelId}?msg={messageId})
 */
export const TITLED_MESSAGE_LINK_REGEX =
  /\[([^\]]+)\]\((https?:\/\/[^\s<>"()]*#\/[a-zA-Z0-9_-]+\/[a-zA-Z0-9_-]+\?msg=[a-zA-Z0-9_-]+)\)/g

export interface TitledMessageLink {
  title: string
  url: string
  params: MessageLinkParams
}

/**
 * 解析带标题的消息链接
 */
export function parseTitledMessageLink(text: string): TitledMessageLink | null {
  TITLED_MESSAGE_LINK_REGEX.lastIndex = 0
  const match = TITLED_MESSAGE_LINK_REGEX.exec(text)
  if (!match) return null

  const [, title, url] = match
  const params = parseMessageLink(url)
  if (!params) return null

  return { title, url, params }
}
