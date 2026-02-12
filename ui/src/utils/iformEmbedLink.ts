export interface IFormEmbedLinkParams {
  worldId: string;
  channelId: string;
  formId: string;
  width?: number;
  height?: number;
}

export interface ParsedSingleIFormEmbedLink extends IFormEmbedLinkParams {
  rawLink: string;
}

const IFORM_LINK_EXACT_REGEX = /^https?:\/\/[^\s<>"']*#\/([a-zA-Z0-9_-]+)\/([a-zA-Z0-9_-]+)\?([^\s#]+)$/;

const normalizeInput = (value: string) => value.replace(/&amp;/gi, '&').trim();

const parsePositiveInt = (value?: string | null): number | undefined => {
  if (!value) {
    return undefined;
  }
  const parsed = Number.parseInt(value, 10);
  if (!Number.isFinite(parsed) || parsed <= 0) {
    return undefined;
  }
  return parsed;
};

const resolveLinkBase = (base?: string): string => {
  const trimmed = (base || '').trim();
  if (trimmed) {
    return trimmed.replace(/\/+$/, '');
  }
  if (typeof window === 'undefined') {
    return '';
  }
  return window.location.origin;
};

const extractLinkBase = (rawLink: string): string | undefined => {
  const normalized = normalizeInput(rawLink);
  const matched = normalized.match(/^(https?:\/\/[^\s<>"']*?)#\/[a-zA-Z0-9_-]+\/[a-zA-Z0-9_-]+\?[^\s#]+$/);
  if (!matched?.[1]) {
    return undefined;
  }
  return matched[1].replace(/\/+$/, '');
};

export function generateIFormEmbedLink(
  params: IFormEmbedLinkParams,
  options?: { base?: string },
): string {
  const base = resolveLinkBase(options?.base);
  const search = new URLSearchParams({
    iform: params.formId,
  });
  if (params.width && params.width > 0) {
    search.set('w', String(Math.round(params.width)));
  }
  if (params.height && params.height > 0) {
    search.set('h', String(Math.round(params.height)));
  }
  return `${base}/#/${params.worldId}/${params.channelId}?${search.toString()}`;
}

export function parseIFormEmbedLink(url: string): IFormEmbedLinkParams | null {
  if (!url || typeof url !== 'string') {
    return null;
  }
  const normalized = normalizeInput(url);
  const match = normalized.match(IFORM_LINK_EXACT_REGEX);
  if (!match) {
    return null;
  }
  const [, worldId, channelId, queryString] = match;
  if (!worldId || !channelId || !queryString) {
    return null;
  }
  const search = new URLSearchParams(queryString);
  const formId = (search.get('iform') || '').trim();
  if (!formId) {
    return null;
  }
  return {
    worldId,
    channelId,
    formId,
    width: parsePositiveInt(search.get('w')),
    height: parsePositiveInt(search.get('h')),
  };
}

export function updateIFormEmbedLinkSize(url: string, width: number, height: number): string | null {
  const parsed = parseIFormEmbedLink(url);
  if (!parsed) {
    return null;
  }
  const nextWidth = Math.max(1, Math.round(width));
  const nextHeight = Math.max(1, Math.round(height));
  return generateIFormEmbedLink(
    {
      worldId: parsed.worldId,
      channelId: parsed.channelId,
      formId: parsed.formId,
      width: nextWidth,
      height: nextHeight,
    },
    { base: extractLinkBase(url) },
  );
}

export function isIFormEmbedLink(url: string): boolean {
  return parseIFormEmbedLink(url) !== null;
}

export function parseSingleIFormEmbedLinkText(text: string): ParsedSingleIFormEmbedLink | null {
  if (!text || typeof text !== 'string') {
    return null;
  }
  const normalized = normalizeInput(text).replace(/\u00a0/g, ' ').trim();
  if (!normalized || /\s/.test(normalized)) {
    return null;
  }
  const parsed = parseIFormEmbedLink(normalized);
  if (!parsed) {
    return null;
  }
  return {
    ...parsed,
    rawLink: normalized,
  };
}
