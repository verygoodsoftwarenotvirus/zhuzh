/**
 * Server-side OpenTelemetry tracing and metrics for SvelteKit (Node).
 * Import this module from hooks.server.ts so it runs at startup.
 * When OTEL_COLLECTOR_GRPC_URL is unset, tracing and metrics are no-op.
 * HTTP server instrumentation creates a span per incoming request so traces appear in Tempo.
 */

import { NodeSDK } from '@opentelemetry/sdk-node';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { OTLPMetricExporter } from '@opentelemetry/exporter-metrics-otlp-grpc';
import { HttpInstrumentation } from '@opentelemetry/instrumentation-http';
import { PeriodicExportingMetricReader } from '@opentelemetry/sdk-metrics';
import { Resource } from '@opentelemetry/resources';
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';
import { env } from '$env/dynamic/private';
import { OPS_PATH_PREFIX } from './utils';

const SERVICE_NAME = 'consumer-web';

const METRIC_EXPORT_INTERVAL_MS = 60_000;

let initialized = false;

export function initServerOtel(): void {
  if (initialized) return;
  initialized = true;

  const collectorGrpcUrl = env.OTEL_COLLECTOR_GRPC_URL?.trim();
  if (!collectorGrpcUrl) {
    return;
  }

  const endpoint =
    collectorGrpcUrl.startsWith('http://') || collectorGrpcUrl.startsWith('https://')
      ? collectorGrpcUrl
      : `http://${collectorGrpcUrl}`;

  const traceExporter = new OTLPTraceExporter({ url: endpoint });
  const metricExporter = new OTLPMetricExporter({ url: endpoint });
  const metricReader = new PeriodicExportingMetricReader({
    exporter: metricExporter,
    exportIntervalMillis: METRIC_EXPORT_INTERVAL_MS,
  });
  const resource = new Resource({ [ATTR_SERVICE_NAME]: SERVICE_NAME });

  function shouldIgnoreRequestForTracing(req: { url?: string }): boolean {
    const pathname = req.url?.split('?')[0] ?? '';
    return pathname === OPS_PATH_PREFIX || pathname.startsWith(`${OPS_PATH_PREFIX}/`);
  }

  const sdk = new NodeSDK({
    resource,
    traceExporter,
    metricReader,
    instrumentations: [
      new HttpInstrumentation({
        ignoreIncomingRequestHook: shouldIgnoreRequestForTracing,
      }),
    ],
    // Rely on collector for sampling in prod; trace everything when collector is configured for now.
    sampler: undefined,
  });

  sdk.start();
}
