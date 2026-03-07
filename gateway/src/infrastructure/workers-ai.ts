// src/infrastructure/workers-ai.ts

export interface AISummaryResult {
  summary: string;
  language: string;
  model: string;
}

const AI_MODEL = '@cf/meta/llama-3.1-8b-instruct' as const;

export class WorkersAIClient {
  constructor(private readonly ai: Ai) {}

  /**
   * gemtextコンテンツから指定言語の要約を生成する。
   */
  async generateSummary(
    gemtext: string,
    targetLanguage: 'ja' | 'en' | 'zh',
  ): Promise<AISummaryResult> {
    const languageNames: Record<string, string> = {
      ja: 'Japanese',
      en: 'English',
      zh: 'Chinese',
    };

    const summary = await this.runChat(
      `You are a technical blog summarizer. Generate a concise summary (2-3 sentences) in ${languageNames[targetLanguage]}. The input is in Gemini protocol's gemtext format. Focus on the key technical concepts and takeaways.`,
      gemtext,
      256,
    );

    return {
      summary,
      language: targetLanguage,
      model: AI_MODEL,
    };
  }

  /**
   * OGP description用の短い要約を生成する。
   */
  async generateOGPDescription(gemtext: string): Promise<string> {
    return this.runChat(
      'Generate a concise meta description (max 160 characters) in the same language as the input. Focus on what the article is about. Do not include quotes or formatting.',
      gemtext,
      64,
    );
  }

  private async runChat(systemPrompt: string, userContent: string, maxTokens: number): Promise<string> {
    const response = await this.ai.run(AI_MODEL, {
      messages: [
        { role: 'system', content: systemPrompt },
        { role: 'user', content: userContent },
      ],
      max_tokens: maxTokens,
    });

    return (response as { response: string }).response;
  }
}
