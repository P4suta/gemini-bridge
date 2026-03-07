// src/config/bindings.ts

export interface Env {
  /** R2バケット: ビルド済み静的アセット */
  readonly CONTENT: R2Bucket;

  /** KV名前空間: メタデータキャッシュ */
  readonly METADATA: KVNamespace;

  /** D1データベース: 記事データ・解析データ */
  readonly DB: D1Database;

  /** Workers AI: 推論エンジン */
  readonly AI: Ai;

  /** Analytics Engine: カスタム解析 */
  readonly ANALYTICS: AnalyticsEngineDataset;

  /** 環境変数 */
  readonly SITE_URL: string;
  readonly SITE_TITLE: string;
}
