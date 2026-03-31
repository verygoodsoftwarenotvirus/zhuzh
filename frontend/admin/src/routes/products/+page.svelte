<script lang="ts">
  import { Heading, Text, Link, Card } from '@zhuzh/ui';

  let { data } = $props();
</script>

<Heading level={1}>Products</Heading>
<p class="subtitle">List and manage products</p>

{#if data?.error}
  <p class="error">{data.error}</p>
{:else if data?.products && data.products.length > 0}
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
        {#each data.products as product (product.id)}
          <tr>
            <td><code>{(product as Record<string, unknown>).id ?? '-'}</code></td>
            <td>{(product as Record<string, unknown>).name ?? '-'}</td>
            <td><Link href="/products/{(product as Record<string, unknown>).id}">View</Link></td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{:else}
  <Card>
    <Text>No products found.</Text>
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
