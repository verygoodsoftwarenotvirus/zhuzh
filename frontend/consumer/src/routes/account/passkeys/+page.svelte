<script lang="ts">
  import { enhance } from '$app/forms';
  import { PageContainer, Button, Alert, Link } from '@zhuzh/ui';
  import type { PasskeyCredential } from '@zhuzh/api-client/auth/auth_service_types';

  let { data } = $props();
  const passkeys = $derived((data?.passkeys ?? []) as PasskeyCredential[]);
  const error = $derived(data?.error as string | null | undefined);
  const deleted = $derived(data?.deleted ?? false);

  const errorMessages: Record<string, string> = {
    invalid: 'Invalid request.',
    delete_failed: 'Failed to remove passkey. Please try again.',
    server: 'Something went wrong. Please try again.',
  };
  const displayError = $derived(error ? (errorMessages[error] ?? 'Something went wrong.') : null);

  function formatDate(d: Date | undefined): string {
    if (!d) return '';
    return new Date(d).toLocaleDateString(undefined, {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  }

  async function addPasskey() {
    const btn = document.getElementById('add-passkey-btn');
    if (!btn || !window.PublicKeyCredential) return;
    btn.setAttribute('disabled', 'true');

    function b64enc(buf: ArrayBuffer): string {
      const b = new Uint8Array(buf);
      let s = '';
      for (let i = 0; i < b.length; i++) s += String.fromCharCode(b[i]);
      return btoa(s).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
    }
    function b64dec(s: string): ArrayBuffer {
      const padded = s.replace(/-/g, '+').replace(/_/g, '/');
      const padded2 = padded + '==='.slice((padded.length + 3) % 4);
      const binary = atob(padded2);
      return Uint8Array.from(binary, (c) => c.charCodeAt(0)).buffer;
    }

    try {
      const optsRes = await fetch('/auth/passkey/registration/options', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({}),
        credentials: 'include',
      });
      if (!optsRes.ok) throw new Error('Failed to get options');
      const opts = await optsRes.json();

      const raw = atob(opts.publicKeyCredentialCreationOptions);
      const obj = JSON.parse(raw);
      const pk = obj.publicKey || obj;
      if (typeof pk.challenge === 'string') pk.challenge = b64dec(pk.challenge);
      if (pk.user && typeof pk.user.id === 'string') pk.user.id = b64dec(pk.user.id);

      const cred = await navigator.credentials.create({ publicKey: pk });
      if (!cred) throw new Error('No credential');

      const pkCred = cred as PublicKeyCredential;
      const r = pkCred.response as AuthenticatorAttestationResponse;
      const attestation = {
        id: pkCred.id,
        rawId: b64enc(pkCred.rawId),
        type: pkCred.type,
        response: {
          clientDataJSON: b64enc(r.clientDataJSON),
          attestationObject: b64enc(r.attestationObject),
        },
      };

      const verifyRes = await fetch('/auth/passkey/registration/verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          attestationResponse: attestation,
          challenge: opts.challenge,
        }),
        credentials: 'include',
      });
      if (!verifyRes.ok) throw new Error('Registration failed');
      window.location.reload();
    } catch (err) {
      btn.removeAttribute('disabled');
      alert(err instanceof Error ? err.message : 'Passkey registration failed');
    }
  }
</script>

<PageContainer>
  <h1>Passkeys</h1>
  <p><Link href="/account/settings" class="back-link">Back to Account Settings</Link></p>

  {#if deleted}
    <Alert variant="info">Passkey removed successfully.</Alert>
  {/if}
  {#if displayError}
    <Alert variant="error">{displayError}</Alert>
  {/if}

  <div class="passkeys-content">
    <div class="add-section">
      <h2>Add passkey</h2>
      <p class="muted">Add a passkey to sign in quickly without a password.</p>
      <Button id="add-passkey-btn" type="button" variant="default" onclick={addPasskey}>Add passkey</Button>
    </div>

    {#if passkeys.length > 0}
      <div class="list-section">
        <h2>Your passkeys</h2>
        <div class="passkey-list">
          {#each passkeys as pk (pk.id)}
            <div class="passkey-card">
              <div class="passkey-info">
                <span class="passkey-name">{pk.friendlyName || 'Passkey'}</span>
                <span class="passkey-details">
                  {formatDate(pk.createdAt)}
                  {#if pk.lastUsedAt}
                    · Last used {formatDate(pk.lastUsedAt)}
                  {/if}
                </span>
              </div>
              <form method="POST" action="?/delete" use:enhance class="delete-form">
                <input type="hidden" name="credential_id" value={pk.id} data-testid="passkey-credential-id" />
                <Button type="submit" variant="default" class="remove-btn">Remove</Button>
              </form>
            </div>
          {/each}
        </div>
      </div>
    {:else}
      <div class="empty-section">
        <p class="muted">No passkeys yet. Add one to sign in quickly without a password.</p>
      </div>
    {/if}
  </div>
</PageContainer>

<style>
  .back-link {
    font-size: 0.875rem;
  }
  .passkeys-content {
    display: flex;
    flex-direction: column;
    gap: var(--space-xl);
    margin-top: var(--space-lg);
  }
  .add-section .muted,
  .empty-section .muted {
    font-size: 0.875rem;
    color: var(--color-muted, #666);
    margin: 0 0 var(--space-md);
  }
  .list-section h2,
  .add-section h2 {
    font-size: 1.125rem;
    font-weight: var(--font-weight-medium);
    margin: 0 0 var(--space-sm);
  }
  .passkey-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
  }
  .passkey-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-md);
    padding: var(--space-md);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-surface);
  }
  .passkey-info {
    flex: 1;
    min-width: 0;
  }
  .passkey-name {
    display: block;
    font-weight: var(--font-weight-medium);
  }
  .passkey-details {
    font-size: 0.875rem;
    color: var(--color-muted, #666);
    margin-top: 0.25rem;
  }
  .delete-form {
    flex-shrink: 0;
  }
  .remove-btn {
    font-size: 0.875rem;
    padding: 0.25rem 0.5rem;
    color: var(--color-error, #c00);
    border-color: var(--color-error, #c00);
  }
</style>
