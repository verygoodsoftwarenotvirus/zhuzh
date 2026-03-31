<script lang="ts">
  import { Heading, Text, Link, Card } from '@zhuzh/ui';

  let { data } = $props();
</script>

<Heading level={1}>Users</Heading>
<p class="subtitle">Search and manage users</p>

{#if data?.error}
  <p class="error">{data.error}</p>
{:else if data?.users && data.users.length > 0}
  <div class="table-wrap">
    <table class="data-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Username</th>
          <th>First name</th>
          <th>Last name</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each data.users as user (user.id)}
          <tr>
            <td><code>{user.id ?? '-'}</code></td>
            <td>{user.username ?? '-'}</td>
            <td>{user.firstName ?? '-'}</td>
            <td>{user.lastName ?? '-'}</td>
            <td><Link href="/users/{user.id}">View</Link></td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{:else}
  <Card>
    <Text>No users found.</Text>
  </Card>
{/if}

<style>
  .subtitle {
    color: var(--color-text-muted);
    margin-bottom: var(--space-lg);
  }
  .error {
    color: var(--color-error, #b91c1c);
    margin-bottom: var(--space-md);
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
  .data-table th {
    font-weight: 600;
    color: var(--color-text-muted);
  }
  .data-table code {
    font-size: 0.875em;
  }
</style>
