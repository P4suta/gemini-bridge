// src/application/serve-content.ts
import type { ContentFormat } from '../domain/gemini/negotiation';
import { negotiateFormat } from '../domain/gemini/negotiation';
import { mapGeminiToHttp } from '../domain/gemini/semantics';
import type { GeminiResponse } from '../domain/gemini/types';

export interface ContentRepository {
  getHtml(slug: string): Promise<string | null>;
  getGemtext(slug: string): Promise<string | null>;
  getMeta(slug: string): Promise<PostMeta | null>;
}

export interface ContentCache {
  get(key: string): Promise<string | null>;
  set(key: string, value: string, ttlSeconds: number): Promise<void>;
}

export interface PostMeta {
  slug: string;
  title: string;
  date: string;
  tags: string[];
  language: string;
  description: string;
}

export interface ServeContentResult {
  body: string;
  contentType: string;
  status: number;
  headers: Record<string, string>;
}

/**
 * コンテンツ配信ユースケース。
 * Content Negotiationとgemini-to-HTTPセマンティクスマッピングを統合する。
 */
export async function serveContent(
  slug: string,
  acceptHeader: string | null,
  repository: ContentRepository,
  cache: ContentCache,
): Promise<ServeContentResult> {
  const format = negotiateFormat(acceptHeader);

  // キャッシュチェック
  const cacheKey = `content:${format}:${slug}`;
  const cached = await cache.get(cacheKey);
  if (cached) {
    const parsed = JSON.parse(cached) as ServeContentResult;
    return parsed;
  }

  // コンテンツ取得
  const content = await fetchContent(slug, format, repository);

  if (!content) {
    // Gemini 51 NOT FOUND -> HTTP 404
    const geminiResponse: GeminiResponse = {
      status: 51,
      meta: 'Content not found',
    };
    const mapping = mapGeminiToHttp(geminiResponse);
    return {
      body: '<h1>404 Not Found</h1><p>The requested Gemini content was not found.</p>',
      contentType: 'text/html; charset=utf-8',
      status: mapping.httpStatus,
      headers: mapping.headers,
    };
  }

  // Gemini 20 SUCCESS -> HTTP 200
  const geminiResponse: GeminiResponse = {
    status: 20,
    meta: format === 'gemtext' ? 'text/gemini; charset=utf-8' : 'text/html; charset=utf-8',
  };
  const mapping = mapGeminiToHttp(geminiResponse);

  const result: ServeContentResult = {
    body: content.body,
    contentType: content.contentType,
    status: mapping.httpStatus,
    headers: {
      ...mapping.headers,
      'Vary': 'Accept',
      'Link': `</posts/${slug}/index.gmi>; rel="alternate"; type="text/gemini"`,
    },
  };

  // キャッシュ保存（1時間）
  await cache.set(cacheKey, JSON.stringify(result), 3600);

  return result;
}

async function fetchContent(
  slug: string,
  format: ContentFormat,
  repository: ContentRepository,
): Promise<{ body: string; contentType: string } | null> {
  switch (format) {
    case 'html': {
      const html = await repository.getHtml(slug);
      return html ? { body: html, contentType: 'text/html; charset=utf-8' } : null;
    }
    case 'gemtext': {
      const gmi = await repository.getGemtext(slug);
      return gmi ? { body: gmi, contentType: 'text/gemini; charset=utf-8' } : null;
    }
    case 'json': {
      const meta = await repository.getMeta(slug);
      return meta ? { body: JSON.stringify(meta), contentType: 'application/json; charset=utf-8' } : null;
    }
  }
}
