<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    /** Visual style: body (default), muted, small */
    variant?: 'body' | 'muted' | 'small';
    /** Optional class for layout */
    class?: string;
    children?: Snippet;
  }

  let { variant = 'body', class: className = '', children }: Props = $props();
</script>

<p class="text text--{variant} {className}">
  {#if children}
    {@render children()}
  {/if}
</p>

<style>
  .text {
    font-family: var(--font-sans);
    line-height: var(--line-height-normal);
    margin: 0 0 1em;
  }
  .text:last-child {
    margin-bottom: 0;
  }
  .text--body {
    font-size: var(--font-size-base);
    color: var(--color-text);
  }
  .text--muted {
    font-size: var(--font-size-sm);
    color: var(--color-text-muted);
  }
  .text--small {
    font-size: var(--font-size-sm);
    color: var(--color-text);
  }
</style>
