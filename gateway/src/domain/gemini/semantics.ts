// domain/gemini/semantics.ts
import type { GeminiResponse, GeminiStatusCode } from './types';

export interface HttpMapping {
  readonly httpStatus: number;
  readonly headers: Record<string, string>;
  readonly transformBody: boolean;
}

const STATUS_MAP: Record<GeminiStatusCode, number> = {
  10: 200, 11: 200,
  20: 200,
  30: 302, 31: 301,
  40: 503, 41: 503, 42: 502, 43: 502, 44: 429,
  50: 410, 51: 404, 52: 410, 53: 502, 59: 400,
  60: 401, 61: 403, 62: 403,
};

export function mapGeminiToHttp(gemini: GeminiResponse): HttpMapping {
  const httpStatus = STATUS_MAP[gemini.status];
  const headers: Record<string, string> = {
    'X-Gemini-Status': String(gemini.status),
    'X-Gemini-Meta': gemini.meta,
  };

  // 3x REDIRECT: Locationヘッダーにリダイレクト先を設定
  if (gemini.status === 30 || gemini.status === 31) {
    headers['Location'] = translateGeminiUrl(gemini.meta);
  }

  // 44 SLOW DOWN: Retry-Afterヘッダーに秒数を設定
  if (gemini.status === 44) {
    headers['Retry-After'] = gemini.meta;
  }

  // 4x TEMPORARY FAILURE: デフォルトRetry-After
  if (gemini.status === 40 || gemini.status === 41) {
    headers['Retry-After'] = '300';
  }

  // 60 CLIENT CERTIFICATE: WWW-Authenticateヘッダー
  if (gemini.status === 60) {
    headers['WWW-Authenticate'] =
      'GeminiCert realm="Gemini client certificate required"';
  }

  return {
    httpStatus,
    headers,
    transformBody: gemini.status === 20,
  };
}

/**
 * gemini:// URLをhttps:// URLに変換する。
 * 自サイト内リンクは直接変換、外部Geminiリンクはプロキシ経由に変換。
 */
function translateGeminiUrl(url: string): string {
  if (!url.startsWith('gemini://')) return url;

  const parsed = new URL(url);
  // 自サイトのGeminiアドレスの場合は直接HTTPS変換
  // 外部の場合はそのまま返す（将来的にプロキシ経由に拡張可能）
  return 'https://' + parsed.host + parsed.pathname + parsed.search;
}
