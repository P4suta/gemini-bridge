// domain/gemini/types.ts

/** Geminiステータスコードのカテゴリ */
export type GeminiStatusCategory =
  | 'INPUT'
  | 'SUCCESS'
  | 'REDIRECT'
  | 'TEMPORARY_FAILURE'
  | 'PERMANENT_FAILURE'
  | 'CLIENT_CERTIFICATE';

/** Geminiステータスコード（2桁整数） */
export type GeminiStatusCode =
  | 10 | 11
  | 20
  | 30 | 31
  | 40 | 41 | 42 | 43 | 44
  | 50 | 51 | 52 | 53 | 59
  | 60 | 61 | 62;

/** Geminiレスポンスの抽象表現 */
export interface GeminiResponse {
  readonly status: GeminiStatusCode;
  readonly meta: string;
  readonly body?: ReadableStream | null;
}

/** ステータスコードからカテゴリを導出 */
export function statusCategory(code: GeminiStatusCode): GeminiStatusCategory {
  const tens = Math.floor(code / 10);
  switch (tens) {
    case 1: return 'INPUT';
    case 2: return 'SUCCESS';
    case 3: return 'REDIRECT';
    case 4: return 'TEMPORARY_FAILURE';
    case 5: return 'PERMANENT_FAILURE';
    case 6: return 'CLIENT_CERTIFICATE';
    default: throw new Error(`Unknown Gemini status category: ${code}`);
  }
}
