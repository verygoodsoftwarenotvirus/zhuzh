<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    title?: string;
    collapsible?: boolean;
    expanded?: boolean;
    class?: string;
    children?: Snippet;
  }

  let { title, collapsible = false, expanded = $bindable(true), class: className = '', children }: Props = $props();
</script>

<div class="card {className}">
  {#if title}
    {#if collapsible}
      <button
        type="button"
        class="card-header card-header--clickable"
        aria-expanded={expanded}
        onclick={() => (expanded = !expanded)}
      >
        <span class="card-title">{title}</span>
        <span class="card-chevron" aria-hidden="true">{expanded ? '▼' : '▶'}</span>
      </button>
    {:else}
      <div class="card-header">
        <span class="card-title">{title}</span>
      </div>
    {/if}
  {/if}
  {#if !collapsible || expanded}
    <div class="card-body">
      {#if children}
        {@render children()}
      {/if}
    </div>
  {/if}
</div>

<style>
  .card {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-surface);
    overflow: hidden;
    box-shadow: var(--shadow-sm);
  }

  .card-header {
    width: 100%;
    padding: var(--space-md);
    border: none;
    border-bottom: 1px solid var(--color-border);
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--color-surface-alt);
    transition: background var(--transition-fast);
    font: inherit;
    text-align: left;
    cursor: default;
  }

  .card-header--clickable {
    cursor: pointer;
  }

  .card-header--clickable:hover {
    background: var(--color-border);
  }

  .card-title {
    font-weight: var(--font-weight-medium);
    font-size: 1rem;
  }

  .card-chevron {
    font-size: 0.75rem;
    color: var(--color-text-muted);
  }

  .card-body {
    padding: var(--space-md);
  }
</style>
