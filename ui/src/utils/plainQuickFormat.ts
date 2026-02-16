const inlineCodePattern = /`([^`\n]+)`/g;
const codeFenceLiteralPattern = /```([\s\S]*?)```/g;
const linkPattern = /\[([^\]\n]+)\]\((https?:\/\/[^\s)]+)\)/gi;
const boldPattern = /\*\*([^\n*][^*\n]*?)\*\*/g;
const italicPattern = /(^|[^*])\*([^*\n]+)\*/g;

const normalizeNewlines = (value: string) => value.replace(/\r\n/g, '\n').replace(/\r/g, '\n');

const isSafeHttpUrl = (value: string) => {
  const normalized = value.replace(/&amp;/g, '&').trim();
  if (!/^https?:\/\//i.test(normalized)) {
    return false;
  }
  try {
    const parsed = new URL(normalized);
    return parsed.protocol === 'http:' || parsed.protocol === 'https:';
  } catch {
    return false;
  }
};

const processInlineFromEscaped = (escapedInput: string) => {
  let text = escapedInput;
  const codeTokens: Array<{ token: string; html: string }> = [];
  const linkTokens: Array<{ token: string; html: string }> = [];

  text = text.replace(inlineCodePattern, (_, code: string) => {
    const token = `__QF_CODE_${codeTokens.length}__`;
    codeTokens.push({ token, html: `<code>${code}</code>` });
    return token;
  });

  text = text.replace(linkPattern, (full: string, label: string, url: string) => {
    if (!isSafeHttpUrl(url)) {
      return full;
    }
    const token = `__QF_LINK_${linkTokens.length}__`;
    linkTokens.push({
      token,
      html: `<a href="${url}" class="text-blue-500" target="_blank" rel="noopener noreferrer">${label}</a>`,
    });
    return token;
  });

  text = text.replace(boldPattern, '<strong>$1</strong>');
  text = text.replace(italicPattern, (_match, prefix: string, body: string) => `${prefix}<em>${body}</em>`);

  linkTokens.forEach((entry) => {
    text = text.split(entry.token).join(entry.html);
  });

  codeTokens.forEach((entry) => {
    text = text.split(entry.token).join(entry.html);
  });

  return text;
};

export const renderQuickFormatHtmlFromEscaped = (escapedInput: string) => {
  if (!escapedInput) {
    return '';
  }

  let text = normalizeNewlines(escapedInput);
  const fenceTokens: Array<{ token: string; html: string }> = [];

  text = text.replace(codeFenceLiteralPattern, (segment: string) => {
    const token = `__QF_FENCE_${fenceTokens.length}__`;
    fenceTokens.push({ token, html: segment });
    return token;
  });

  text = processInlineFromEscaped(text);

  fenceTokens.forEach((entry) => {
    text = text.split(entry.token).join(entry.html);
  });

  text = text.replace(/\n/g, '<br />');

  return text;
};
