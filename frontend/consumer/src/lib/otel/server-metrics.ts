/**
 * Server-side request metrics. Uses the global MeterProvider set by initServerOtel()
 * when OTEL_COLLECTOR_GRPC_URL is configured. No-op otherwise.
 * Kept minimal: count by status_class only, one global duration histogram (no path labels).
 */

import { metrics } from '@opentelemetry/api';
import { shouldRecordPathForMetrics } from './utils';

const METER_NAME = 'consumer-web-server';
const METER_VERSION = '1.0';

let _requestCounter: { add: (value: number, attributes?: Record<string, string>) => void } | null = null;
let _requestDurationHistogram: { record: (value: number) => void } | null = null;

function getInstruments(): {
  requestCounter: NonNullable<typeof _requestCounter>;
  requestDuration: NonNullable<typeof _requestDurationHistogram>;
} | null {
  if (_requestCounter && _requestDurationHistogram) {
    return { requestCounter: _requestCounter, requestDuration: _requestDurationHistogram };
  }
  const meter = metrics.getMeter(METER_NAME, METER_VERSION);
  _requestCounter = meter.createCounter('consumer_web_http_requests_total', {
    description: 'Total HTTP requests',
  });
  _requestDurationHistogram = meter.createHistogram('consumer_web_http_request_duration_ms', {
    description: 'HTTP request duration in milliseconds',
    unit: 'ms',
  });
  return { requestCounter: _requestCounter, requestDuration: _requestDurationHistogram };
}

/**
 * Record a completed HTTP request for metrics.
 */
export function recordRequest(pathname: string, statusCode: number, durationMs: number): void {
  if (!shouldRecordPathForMetrics(pathname)) return;
  const instruments = getInstruments();
  if (!instruments) return;
  const statusClass = statusCode >= 500 ? '5xx' : statusCode >= 400 ? '4xx' : statusCode >= 300 ? '3xx' : '2xx';
  instruments.requestCounter.add(1, { status_class: statusClass });
  instruments.requestDuration.record(durationMs);
}
