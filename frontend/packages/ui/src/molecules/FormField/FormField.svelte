<script lang="ts">
  import type { Snippet } from 'svelte';
  import Label from '../../atoms/Label/Label.svelte';

  interface Props {
    id?: string;
    label: string;
    error?: string;
    required?: boolean;
    children?: Snippet;
  }

  let { id, label, error, required = false, children }: Props = $props();
</script>

<div class="form-field">
  <Label for={id}
    >{label}{#if required}<span class="required"> *</span>{/if}</Label
  >
  {#if children}
    {@render children()}
  {/if}
  {#if error}
    <p class="error" role="alert">{error}</p>
  {/if}
</div>

<style>
  .form-field {
    margin-bottom: var(--space-md);
    padding-right: var(--space-md);
  }

  .required {
    color: var(--color-error);
  }

  .error {
    color: var(--color-error);
    margin-top: var(--space-xs);
    margin-bottom: 0;
    font-size: var(--font-size-sm);
  }
</style>
