<script lang="ts">
  interface Option {
    value: string;
    label: string;
  }

  interface Props {
    id?: string;
    name?: string;
    value?: string;
    options: Option[];
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    class?: string;
  }

  let {
    id,
    name,
    value = $bindable(''),
    options,
    placeholder,
    disabled = false,
    required = false,
    class: className = '',
  }: Props = $props();
</script>

<select {id} {name} {required} {disabled} class="select-input {className}" bind:value>
  {#if placeholder}
    <option value="" disabled>{placeholder}</option>
  {/if}
  {#each options as opt}
    <option value={opt.value}>{opt.label}</option>
  {/each}
</select>

<style>
  .select-input {
    width: 100%;
    padding: var(--space-sm);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    font-family: var(--font-sans);
    font-size: var(--font-size-base);
    color: var(--color-text);
    background: var(--color-surface);
    cursor: pointer;
    box-shadow: var(--shadow-sm);
    transition:
      border-color var(--transition-fast),
      box-shadow var(--transition-fast);
  }

  .select-input:focus {
    outline: none;
    border-color: var(--color-primary);
    box-shadow: 0 0 0 2px var(--color-primary-muted);
  }

  .select-input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
