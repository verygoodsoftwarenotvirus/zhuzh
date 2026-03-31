<script lang="ts">
  import { Heading, Text, Link, Card } from '@zhuzh/ui';

  let { data } = $props();
  const client = data?.client as Record<string, unknown> | null;
</script>

{#if data?.error}
  <p class="error">{data.error}</p>
{:else if client}
  <Heading level={1}>OAuth2 Client</Heading>
  <Card>
    <dl class="detail-list">
      <dt>ID</dt>
      <dd><code>{client.id ?? '-'}</code></dd>
      <dt>Client ID</dt>
      <dd>{client.clientId ?? '-'}</dd>
      <dt>Name</dt>
      <dd>{client.name ?? '-'}</dd>
    </dl>
  </Card>
  <p><Link href="/oauth2_clients">Back to OAuth2 clients</Link></p>
{:else}
  <Text>Client not found.</Text>
{/if}

<style>
  .error {
    color: var(--color-error, #b91c1c);
  }
  .detail-list {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: var(--space-xs) var(--space-lg);
  }
  .detail-list dt {
    color: var(--color-text-muted);
  }
</style>
