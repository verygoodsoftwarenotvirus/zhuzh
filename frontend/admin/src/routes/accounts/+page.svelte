<script lang="ts">
  import { Heading, Text, Link, Card } from '@zhuzh/ui';

  let { data } = $props();
</script>

<Heading level={1}>Accounts</Heading>
<p class="subtitle">Search and manage accounts</p>

{#if data?.error}
  <p class="error">{data.error}</p>
{:else if data?.accounts && data.accounts.length > 0}
  <div class="table-wrap">
    <table class="data-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Name</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each data.accounts as account (account.id)}
          <tr>
            <td><code>{account.id ?? '-'}</code></td>
            <td>{account.name ?? '-'}</td>
            <td><Link href="/accounts/{account.id}">View</Link></td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{:else}
  <Card>
    <Text>No accounts found.</Text>
  </Card>
{/if}

<style>
  .subtitle {
    color: var(--color-text-muted);
    margin-bottom: var(--space-lg);
  }
  .error {
    color: var(--color-error, #b91c1c);
  }
  .table-wrap {
    overflow-x: auto;
  }
  .data-table {
    width: 100%;
    border-collapse: collapse;
  }
  .data-table th,
  .data-table td {
    padding: var(--space-sm) var(--space-md);
    text-align: left;
    border-bottom: 1px solid var(--color-border);
  }
</style>
