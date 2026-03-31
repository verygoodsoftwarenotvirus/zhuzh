/**
 * Helpers to get the current OpenTelemetry span and add attributes/events.
 * Use in server-side code (load functions, form actions, +page.server.ts, +layout.server.ts)
 * where the HTTP instrumentation has created an active request span.
 *
 * Example:
 *   import { getCurrentSpan, setSpanAttribute } from '$lib/otel/span';
 *
 *   export const load = async ({ route }) => {
 *     setSpanAttribute('app.route.id', route.id);
 *     getCurrentSpan()?.addEvent('load.started', { route: route.id });
 *     // ...
 *   }
 *
 * On the client there is no automatic request span; getCurrentSpan() will be undefined
 * unless you start a span (e.g. in layout) and run code inside its context.
 */

import { trace } from '@opentelemetry/api';

export type SpanAttributeValue = string | number | boolean | string[] | number[] | boolean[];

/**
 * Returns the currently active span, if any. On the server this is the HTTP request span
 * created by HttpInstrumentation. Safe to call from client (returns undefined if no span).
 */
export function getCurrentSpan(): ReturnType<typeof trace.getActiveSpan> {
  return trace.getActiveSpan();
}

/**
 * Set a single attribute on the current span. No-op if there is no active span.
 */
export function setSpanAttribute(key: string, value: SpanAttributeValue): void {
  const span = trace.getActiveSpan();
  span?.setAttribute(key, value);
}
