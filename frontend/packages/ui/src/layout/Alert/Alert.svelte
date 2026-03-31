<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'error' | 'info';
    class?: string;
    children?: Snippet;
  }

  let { variant = 'info', class: className = '', children }: Props = $props();
</script>

<div class="alert alert--{variant} {className}" role="alert">
  {#if children}
    {@render children()}
  {/if}
</div>

<style>
  .alert {
    padding: var(--space-md);
    border-radius: var(--radius-sm);
    margin-bottom: var(--space-md);
    font-family: var(--font-sans);
    box-shadow: var(--shadow-sm);
  }

  .alert--error {
    color: var(--color-error);
    background: var(--color-error-muted);
    border: 1px solid var(--color-error);
  }

  .alert--info {
    color: var(--color-text);
    background: var(--color-surface-alt);
    border: 1px solid var(--color-border);
  }
</style>
