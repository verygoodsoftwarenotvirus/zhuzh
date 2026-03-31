<script lang="ts">
  interface Props {
    id?: string;
    value?: string;
    placeholder?: string;
    disabled?: boolean;
    required?: boolean;
    class?: string;
    dataTestId?: string;
    suggestions?: { id: string; label: string }[];
    loading?: boolean;
    onInput?: (value: string) => void;
    onSelect?: (item: { id: string; label: string }) => void;
    onBlur?: () => void;
  }

  let {
    id,
    value = $bindable(''),
    placeholder,
    disabled = false,
    required = false,
    class: className = '',
    dataTestId,
    suggestions = [],
    loading = false,
    onInput,
    onSelect,
    onBlur,
  }: Props = $props();

  let open = $state(false);
  let focusedIndex = $state(-1);

  function handleInput(e: Event) {
    const target = e.target as HTMLInputElement;
    value = target.value;
    open = true;
    focusedIndex = -1;
    onInput?.(target.value);
  }

  function handleKeydown(e: KeyboardEvent) {
    if (!open || suggestions.length === 0) {
      if (e.key === 'Escape') open = false;
      return;
    }
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      focusedIndex = Math.min(focusedIndex + 1, suggestions.length - 1);
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      focusedIndex = Math.max(focusedIndex - 1, 0);
    } else if (e.key === 'Enter' && focusedIndex >= 0) {
      e.preventDefault();
      select(suggestions[focusedIndex]);
    } else if (e.key === 'Escape') {
      e.preventDefault();
      open = false;
    }
  }

  function select(item: { id: string; label: string }) {
    value = item.label;
    open = false;
    focusedIndex = -1;
    onSelect?.(item);
  }

  function handleBlur() {
    setTimeout(() => {
      open = false;
      focusedIndex = -1;
      onBlur?.();
    }, 150);
  }
</script>

<div class="autocomplete-wrapper">
  <input
    {id}
    type="text"
    {placeholder}
    {required}
    {disabled}
    class="autocomplete-input {className}"
    bind:value
    oninput={handleInput}
    onkeydown={handleKeydown}
    onfocus={() => (open = true)}
    onblur={handleBlur}
    autocomplete="off"
    role="combobox"
    aria-expanded={open}
    aria-autocomplete="list"
    aria-controls="autocomplete-list-{id ?? 'default'}"
    data-testid={dataTestId}
  />
  {#if loading}
    <span class="autocomplete-loading" aria-hidden="true">...</span>
  {/if}
  {#if open && (suggestions.length > 0 || loading)}
    <ul id="autocomplete-list-{id ?? 'default'}" class="autocomplete-list" role="listbox">
      {#each suggestions as item, i}
        <li
          role="option"
          aria-selected={i === focusedIndex}
          class="autocomplete-option"
          class:autocomplete-option--focused={i === focusedIndex}
          onmousedown={(e) => {
            e.preventDefault();
            select(item);
          }}
        >
          {item.label}
        </li>
      {/each}
    </ul>
  {/if}
</div>

<style>
  .autocomplete-wrapper {
    position: relative;
    width: 100%;
  }

  .autocomplete-input {
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

  .autocomplete-input::placeholder {
    color: var(--color-text-muted);
  }

  .autocomplete-input:focus {
    outline: none;
    border-color: var(--color-primary);
    box-shadow: 0 0 0 2px var(--color-primary-muted);
  }

  .autocomplete-loading {
    position: absolute;
    right: var(--space-sm);
    top: 50%;
    transform: translateY(-50%);
    color: var(--color-text-muted);
    font-size: 0.875rem;
  }

  .autocomplete-list {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    margin: 0;
    padding: 0;
    list-style: none;
    border: 1px solid var(--color-border);
    border-top: none;
    border-radius: 0 0 var(--radius-sm) var(--radius-sm);
    background: var(--color-surface);
    max-height: 12rem;
    overflow-y: auto;
    z-index: 50;
    box-shadow: var(--shadow-md);
  }

  .autocomplete-option {
    padding: var(--space-sm) var(--space-md);
    cursor: pointer;
    font-family: var(--font-sans);
    font-size: 1rem;
  }

  .autocomplete-option:hover,
  .autocomplete-option--focused {
    background: var(--color-surface-alt);
  }
</style>
