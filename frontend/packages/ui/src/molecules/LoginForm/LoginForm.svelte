<script lang="ts">
  import type { Snippet } from 'svelte';
  import { enhance, applyAction } from '$app/forms';
  import FormField from '../FormField/FormField.svelte';
  import Input from '../../atoms/Input/Input.svelte';
  import Button from '../../atoms/Button/Button.svelte';
  import Alert from '../../layout/Alert/Alert.svelte';

  interface Props {
    action?: string;
    method?: 'POST' | 'GET';
    username?: string;
    /** When true, TOTP field is always visible (e.g. admin login where 2FA is required). */
    showTotp?: boolean;
    /** When true with showTotp, the TOTP field is required and labeled as such. */
    totpRequired?: boolean;
    error?: string;
    submitLabel?: string;
    passkeySlot?: Snippet;
  }

  let {
    action = '?/login',
    method = 'POST',
    username = '',
    showTotp = false,
    totpRequired = false,
    error,
    submitLabel = 'Sign In',
    passkeySlot,
  }: Props = $props();
</script>

<form
  {method}
  {action}
  use:enhance={() => {
    return async ({ result }) => {
      if (result.type === 'redirect') {
        window.location.href = result.location;
      } else {
        await applyAction(result);
      }
    };
  }}
>
  {#if error}
    <Alert variant="error">{error}</Alert>
  {/if}
  <FormField id="username" label="Username" required>
    <Input
      id="username"
      name="username"
      type="text"
      autocomplete="username"
      value={username}
      required
      dataTestId="login-username"
    />
  </FormField>
  <FormField id="password" label="Password" required>
    <Input
      id="password"
      name="password"
      type="password"
      autocomplete="current-password"
      required
      dataTestId="login-password"
    />
  </FormField>
  {#if showTotp}
    <FormField
      id="totpToken"
      label={totpRequired ? 'Authentication code (2FA)' : 'TOTP (if enabled)'}
      required={totpRequired}
    >
      <Input
        id="totpToken"
        name="totpToken"
        type="text"
        autocomplete="one-time-code"
        placeholder={totpRequired ? 'Enter 6-digit code' : 'Optional'}
        required={totpRequired}
        dataTestId="login-totp-token"
      />
    </FormField>
  {/if}
  <Button type="submit">{submitLabel}</Button>
  {#if passkeySlot}
    {@render passkeySlot()}
  {/if}
</form>
