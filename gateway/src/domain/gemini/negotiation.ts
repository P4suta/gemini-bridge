// domain/gemini/negotiation.ts
export type ContentFormat = 'html' | 'gemtext' | 'json';

export function negotiateFormat(acceptHeader: string | null): ContentFormat {
  if (!acceptHeader || acceptHeader === '*/*') return 'html';

  const mediaTypes = parseAcceptHeader(acceptHeader);

  for (const { type } of mediaTypes) {
    if (type === 'text/gemini') return 'gemtext';
    if (type === 'text/html') return 'html';
    if (type === 'application/json') return 'json';
  }

  return 'html';
}

interface MediaType {
  type: string;
  quality: number;
}

function parseAcceptHeader(header: string): MediaType[] {
  return header
    .split(',')
    .map((part) => {
      const [type, ...params] = part.trim().split(';');
      const qParam = params.find((p) => p.trim().startsWith('q='));
      const quality = qParam ? parseFloat(qParam.split('=')[1]) : 1.0;
      return { type: type.trim(), quality };
    })
    .sort((a, b) => b.quality - a.quality);
}
