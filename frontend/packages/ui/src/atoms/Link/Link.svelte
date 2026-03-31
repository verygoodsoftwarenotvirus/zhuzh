<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    href: string;
    class?: string;
    children?: Snippet;
  }

  let { href, class: className = '', children }: Props = $props();
</script>

<a {href} class="link {className}">
  {#if children}
    {@render children()}
  {/if}
</a>

<style>
  .link {
    color: var(--color-primary);
    text-decoration: none;
    font-family: var(--font-sans);
    transition:
      color var(--transition-fast),
      text-decoration-color var(--transition-fast);
  }

  .link:hover {
    text-decoration: underline;
  }
</style>
