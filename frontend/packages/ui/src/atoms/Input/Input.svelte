<script lang="ts">
  interface Props {
    id?: string;
    name?: string;
    type?: string;
    value?: string;
    placeholder?: string;
    required?: boolean;
    minlength?: number;
    autocomplete?: string;
    disabled?: boolean;
    class?: string;
    dataTestId?: string;
    oninput?: (e: Event & { currentTarget: HTMLInputElement }) => void;
  }

  let {
    id,
    name,
    type = 'text',
    value = $bindable(''),
    placeholder,
    required = false,
    minlength,
    autocomplete,
    disabled = false,
    class: className = '',
    dataTestId,
    oninput,
  }: Props = $props();
</script>

<input
  {id}
  {name}
  {type}
  {placeholder}
  {required}
  {minlength}
  autocomplete={autocomplete as HTMLInputElement['autocomplete']}
  {disabled}
  class="input {className}"
  {value}
  data-testid={dataTestId}
  {oninput}
/>

<style>
  .input {
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

  .input::placeholder {
    color: var(--color-text-muted);
  }

  .input:focus {
    outline: none;
    border-color: var(--color-primary);
    box-shadow: 0 0 0 2px var(--color-primary-muted);
  }
</style>
