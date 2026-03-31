<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    id?: string;
    type?: 'button' | 'submit' | 'reset';
    variant?: 'primary' | 'default';
    disabled?: boolean;
    class?: string;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
  }

  let {
    id,
    type = 'button',
    variant = 'primary',
    disabled = false,
    class: className = '',
    onclick,
    children,
  }: Props = $props();
</script>

<button {id} {type} {disabled} {onclick} class="btn btn--{variant} {className}">
  {#if children}
    {@render children()}
  {/if}
</button>

<style>
  .btn {
    padding: var(--space-sm) var(--space-md);
    border: none;
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-family: var(--font-sans);
    font-weight: var(--font-weight-medium);
    box-shadow: var(--shadow-sm);
    transition:
      background var(--transition-fast),
      box-shadow var(--transition-fast);
  }

  .btn:focus-visible {
    outline: none;
    box-shadow: var(--shadow-focus);
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn--primary {
    background: var(--color-primary);
    color: var(--color-surface);
  }

  .btn--primary:hover:not(:disabled) {
    background: var(--color-primary-hover);
  }

  .btn--default {
    background: var(--color-surface-alt);
    color: var(--color-text);
    border: 1px solid var(--color-border);
  }

  .btn--default:hover:not(:disabled) {
    background: var(--color-border);
  }
</style>
