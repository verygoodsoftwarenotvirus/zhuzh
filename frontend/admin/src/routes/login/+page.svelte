<script lang="ts">
  import { browser } from '$app/environment';
  import { PageContainer, LoginForm, Button } from '@zhuzh/ui';

  let { form } = $props();
  const supportsPasskey = browser && typeof PublicKeyCredential !== 'undefined';

  async function signInWithPasskey() {
    if (!window.PublicKeyCredential) {
      alert('Passkeys are not supported in this browser.');
      return;
    }
    const usernameInput = document.getElementById('username') as HTMLInputElement | null;
    const username = usernameInput?.value?.trim() ?? '';

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
      const optsRes = await fetch('/auth/passkey/authentication/options', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username }),
        credentials: 'include',
      });
      if (!optsRes.ok) throw new Error('Failed to get options');
      const opts = await optsRes.json();

      const raw = atob(opts.publicKeyCredentialRequestOptions);
      const obj = JSON.parse(raw);
      const pk = obj.publicKey || obj;
      if (typeof pk.challenge === 'string') pk.challenge = b64dec(pk.challenge);
      if (pk.allowCredentials) {
        for (let i = 0; i < pk.allowCredentials.length; i++) {
          const c = pk.allowCredentials[i];
          if (typeof c.id === 'string') c.id = b64dec(c.id);
        }
      }

      const cred = await navigator.credentials.get({ publicKey: pk });
      if (!cred) throw new Error('No credential');

      const pkCred = cred as PublicKeyCredential;
      const r = pkCred.response as AuthenticatorAssertionResponse;
      const assertion = {
        id: pkCred.id,
        rawId: b64enc(pkCred.rawId),
        type: pkCred.type,
        response: {
          clientDataJSON: b64enc(r.clientDataJSON),
          authenticatorData: b64enc(r.authenticatorData),
          signature: b64enc(r.signature),
          userHandle: r.userHandle ? b64enc(r.userHandle) : null,
        },
      };

      const verifyRes = await fetch('/auth/passkey/authentication/verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          challenge: opts.challenge,
          username,
          assertionResponse: assertion,
        }),
        credentials: 'include',
      });
      if (!verifyRes.ok) throw new Error('Authentication failed');
      const result = await verifyRes.json();
      if (result.redirect) {
        window.location.href = result.redirect;
      } else {
        window.location.href = '/';
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Passkey sign-in failed');
    }
  }
</script>

<PageContainer narrow>
  <h1>Admin Login</h1>
  <LoginForm action="?/login" username={form?.username ?? ''} error={form?.error} showTotp={true} totpRequired={true}>
    {#snippet passkeySlot()}
      {#if supportsPasskey}
        <div class="passkey-section">
          <p class="divider">or</p>
          <Button type="button" variant="default" onclick={signInWithPasskey}>Sign in with passkey</Button>
        </div>
      {/if}
    {/snippet}
  </LoginForm>
</PageContainer>

<style>
  .passkey-section {
    margin-top: var(--space-lg);
  }
  .divider {
    margin: var(--space-md) 0;
    color: var(--color-text-muted);
    font-size: var(--font-size-sm);
  }
</style>
