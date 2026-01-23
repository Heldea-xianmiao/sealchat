import { createImageTokenRegex, isValidAttachmentToken } from './attachmentMarkdown';

interface TipTapNode {
  type: string;
  attrs?: Record<string, any>;
  content?: TipTapNode[];
  text?: string;
  marks?: Array<{ type: string; attrs?: Record<string, any> }>;
}

const ATTACHMENT_URL_PATTERN = /\/api\/v1\/attachment\/([^/?#]+)/i;

const extractAttachmentId = (src: string): string | null => {
  if (!src) return null;
  if (src.startsWith('id:')) return src.slice(3);
  const match = src.match(ATTACHMENT_URL_PATTERN);
  if (match) return match[1];
  return null;
};

export function convertPlainWithImagesToTiptap(text: string): TipTapNode {
  if (!text) {
    return { type: 'doc', content: [{ type: 'paragraph' }] };
  }

  const lines = text.split(/\r?\n/);
  const content: TipTapNode[] = [];
  const imageRegex = createImageTokenRegex();

  for (const line of lines) {
    const paragraph: TipTapNode = { type: 'paragraph' };
    const nodes: TipTapNode[] = [];
    let lastIndex = 0;
    let match: RegExpExecArray | null;

    imageRegex.lastIndex = 0;
    while ((match = imageRegex.exec(line)) !== null) {
      if (match.index > lastIndex) {
        const textPart = line.slice(lastIndex, match.index);
        if (textPart) {
          nodes.push({ type: 'text', text: textPart });
        }
      }

      const [full, alt, token] = match;
      if (!isValidAttachmentToken(token)) {
        if (full) {
          nodes.push({ type: 'text', text: full });
        }
      } else {
        nodes.push({
          type: 'image',
          attrs: {
            src: `/api/v1/attachment/${token}`,
            alt: alt || '',
          },
        });
      }

      lastIndex = match.index + full.length;
    }

    if (lastIndex < line.length) {
      const textPart = line.slice(lastIndex);
      if (textPart) {
        nodes.push({ type: 'text', text: textPart });
      }
    }

    if (nodes.length > 0) {
      paragraph.content = nodes;
    }
    content.push(paragraph);
  }

  if (content.length === 0) {
    return { type: 'doc', content: [{ type: 'paragraph' }] };
  }

  return { type: 'doc', content };
}

export function convertTiptapToPlainWithImages(json: TipTapNode | string): string {
  let parsed: TipTapNode | null = null;
  if (typeof json === 'string') {
    try {
      parsed = JSON.parse(json);
    } catch {
      return json;
    }
  } else {
    parsed = json;
  }

  const extractText = (node: TipTapNode): string => {
    if (!node) return '';

    if (node.text !== undefined) {
      return node.text;
    }

    if (node.type === 'hardBreak') {
      return '\n';
    }

    if (node.type === 'image') {
      const src = String(node.attrs?.src || '');
      const alt = String(node.attrs?.alt || '');
      const token = extractAttachmentId(src);
      if (token && isValidAttachmentToken(token)) {
        return `![${alt}](id:${token})`;
      }
      return alt || '';
    }

    if (node.content && node.content.length > 0) {
      const childTexts = node.content.map((child) => extractText(child));
      const joined = childTexts.join('');
      if (node.type === 'paragraph' || node.type === 'heading' || node.type === 'listItem') {
        return joined + '\n';
      }
      return joined;
    }

    return '';
  };

  const result = parsed ? extractText(parsed) : '';
  return result.replace(/\n+$/, '');
}
