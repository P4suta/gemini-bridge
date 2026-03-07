// src/domain/gemini/semantics.test.ts
import { describe, it, expect } from 'vitest';
import { mapGeminiToHttp } from './semantics';
import type { GeminiResponse } from './types';

describe('mapGeminiToHttp', () => {
  it('maps 20 SUCCESS to HTTP 200', () => {
    const gemini: GeminiResponse = { status: 20, meta: 'text/gemini' };
    const result = mapGeminiToHttp(gemini);

    expect(result.httpStatus).toBe(200);
    expect(result.headers['X-Gemini-Status']).toBe('20');
    expect(result.headers['X-Gemini-Meta']).toBe('text/gemini');
    expect(result.transformBody).toBe(true);
  });

  it('maps 30 REDIRECT to HTTP 302 with Location header', () => {
    const gemini: GeminiResponse = { status: 30, meta: 'gemini://example.com/new' };
    const result = mapGeminiToHttp(gemini);

    expect(result.httpStatus).toBe(302);
    expect(result.headers['Location']).toBe('https://example.com/new');
  });

  it('maps 31 PERMANENT REDIRECT to HTTP 301', () => {
    const gemini: GeminiResponse = { status: 31, meta: '/new-path' };
    const result = mapGeminiToHttp(gemini);

    expect(result.httpStatus).toBe(301);
    expect(result.headers['Location']).toBe('/new-path');
  });

  it('maps 44 SLOW DOWN to HTTP 429 with Retry-After', () => {
    const gemini: GeminiResponse = { status: 44, meta: '30' };
    const result = mapGeminiToHttp(gemini);

    expect(result.httpStatus).toBe(429);
    expect(result.headers['Retry-After']).toBe('30');
  });

  it('maps 51 NOT FOUND to HTTP 404', () => {
    const gemini: GeminiResponse = { status: 51, meta: 'Page not found' };
    const result = mapGeminiToHttp(gemini);

    expect(result.httpStatus).toBe(404);
  });

  it('maps 60 CLIENT CERT to HTTP 401 with WWW-Authenticate', () => {
    const gemini: GeminiResponse = { status: 60, meta: 'Certificate required' };
    const result = mapGeminiToHttp(gemini);

    expect(result.httpStatus).toBe(401);
    expect(result.headers['WWW-Authenticate']).toContain('GeminiCert');
  });

  it('maps 10 INPUT to HTTP 200', () => {
    const gemini: GeminiResponse = { status: 10, meta: 'Enter your name' };
    const result = mapGeminiToHttp(gemini);

    expect(result.httpStatus).toBe(200);
    expect(result.transformBody).toBe(false);
  });

  it('maps 40 TEMPORARY FAILURE to HTTP 503 with Retry-After', () => {
    const gemini: GeminiResponse = { status: 40, meta: 'Server error' };
    const result = mapGeminiToHttp(gemini);

    expect(result.httpStatus).toBe(503);
    expect(result.headers['Retry-After']).toBe('300');
  });

  it('maps 52 GONE to HTTP 410', () => {
    const gemini: GeminiResponse = { status: 52, meta: 'Content removed' };
    const result = mapGeminiToHttp(gemini);

    expect(result.httpStatus).toBe(410);
  });

  it('maps 59 BAD REQUEST to HTTP 400', () => {
    const gemini: GeminiResponse = { status: 59, meta: 'Bad request' };
    const result = mapGeminiToHttp(gemini);

    expect(result.httpStatus).toBe(400);
  });

  it('preserves X-Gemini-Status and X-Gemini-Meta for all codes', () => {
    const gemini: GeminiResponse = { status: 42, meta: 'CGI error details' };
    const result = mapGeminiToHttp(gemini);

    expect(result.headers['X-Gemini-Status']).toBe('42');
    expect(result.headers['X-Gemini-Meta']).toBe('CGI error details');
  });
});
