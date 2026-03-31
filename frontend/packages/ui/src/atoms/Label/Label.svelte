<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    for?: string;
    class?: string;
    children?: Snippet;
  }

  let { for: forId, class: className = '', children }: Props = $props();
</script>

<label for={forId} class="label {className}">
  {#if children}
    {@render children()}
  {/if}
</label>

<style>
  .label {
    display: block;
    margin-bottom: var(--space-xs);
    font-weight: var(--font-weight-medium);
    font-family: var(--font-sans);
    color: var(--color-text);
  }
</style>
