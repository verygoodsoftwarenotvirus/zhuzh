<script lang="ts">
  import { Heading, Button } from '@zhuzh/ui';

  let { form } = $props();
</script>

<Heading level={1}>Analytics Test</Heading>
<p class="subtitle">Send a test analytics event</p>

<form method="POST" action="?/default">
  <label for="event_name">Event name</label>
  <input
    id="event_name"
    name="event_name"
    type="text"
    value="admin_analytics_test"
    placeholder="admin_analytics_test"
  />
  <Button type="submit">Send test event</Button>
</form>

{#if form?.success}
  <p class="result success">{form.message ?? 'Event sent.'}</p>
{:else if form?.error}
  <p class="result error">{form.error}</p>
{/if}

<style>
  .subtitle {
    color: var(--color-text-muted);
    margin-bottom: var(--space-lg);
  }
  form {
    display: flex;
    flex-direction: column;
    gap: var(--space-md);
    max-width: 20rem;
  }
  label {
    font-weight: 500;
  }
  input {
    padding: var(--space-sm) var(--space-md);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
  }
  .result {
    margin-top: var(--space-lg);
  }
  .result.success {
    color: var(--color-success, green);
  }
  .result.error {
    color: var(--color-error, #b91c1c);
  }
</style>
