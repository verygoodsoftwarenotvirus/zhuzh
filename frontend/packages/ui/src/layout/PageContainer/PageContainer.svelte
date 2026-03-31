<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    narrow?: boolean;
    /** Use for wide content (e.g. recipe editor): starts ~15% from left, max 70vw width */
    wide?: boolean;
    class?: string;
    children?: Snippet;
  }

  let { narrow = false, wide = false, class: className = '', children }: Props = $props();
</script>

<main class="page-container {narrow ? 'page-container--narrow' : ''} {wide ? 'page-container--wide' : ''} {className}">
  {#if children}
    {@render children()}
  {/if}
</main>

<style>
  .page-container {
    max-width: var(--content-max-width);
    margin: var(--space-lg) auto;
    padding: var(--space-md);
  }

  .page-container--narrow {
    max-width: var(--content-max-width-narrow);
  }

  .page-container--wide {
    max-width: 70vw;
    margin-left: 15vw;
    margin-right: auto;
  }

  @media (max-width: 768px) {
    .page-container--wide {
      max-width: 100%;
      margin-left: 0;
      margin-right: 0;
    }
  }
</style>
