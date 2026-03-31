<script lang="ts">
  interface Props {
    id?: string;
    name?: string;
    value?: number;
    min?: number;
    max?: number;
    step?: number;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    class?: string;
    dataTestId?: string;
  }

  let {
    id,
    name,
    value = $bindable(0),
    min,
    max,
    step = 1,
    placeholder,
    disabled = false,
    required = false,
    class: className = '',
    dataTestId,
  }: Props = $props();

  function handleInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const parsed = parseFloat(target.value);
    if (!isNaN(parsed)) {
      value = parsed;
    }
  }
</script>

<input
  {id}
  {name}
  type="number"
  {min}
  {max}
  {step}
  {placeholder}
  {required}
  {disabled}
  class="number-input {className}"
  bind:value
  oninput={handleInput}
  data-testid={dataTestId}
/>

<style>
  .number-input {
    width: 100%;
    padding: var(--space-sm);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    font-family: var(--font-sans);
    font-size: var(--font-size-base);
    color: var(--color-text);
    box-shadow: var(--shadow-sm);
    transition:
      border-color var(--transition-fast),
      box-shadow var(--transition-fast);
  }

  .number-input::placeholder {
    color: var(--color-text-muted);
  }

  .number-input:focus {
    outline: none;
    border-color: var(--color-primary);
    box-shadow: 0 0 0 2px var(--color-primary-muted);
  }
</style>
