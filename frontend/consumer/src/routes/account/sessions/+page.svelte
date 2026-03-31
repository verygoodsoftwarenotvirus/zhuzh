<script lang="ts">
  import { enhance } from '$app/forms';
  import { PageContainer, Button, Alert, Link } from '@zhuzh/ui';
  import type { UserSession } from '@zhuzh/api-client/auth/auth_messages';

  let { data } = $props();
  const sessions = $derived((data?.sessions ?? []) as UserSession[]);
  const error = $derived(data?.error as string | null | undefined);
  const revoked = $derived(data?.revoked ?? false);
  const revokedAll = $derived(data?.revokedAll ?? false);
  const hasOtherSessions = $derived(sessions.some((s) => !s.isCurrent));

  const errorMessages: Record<string, string> = {
    invalid: 'Invalid request.',
    revoke_failed: 'Failed to revoke session. Please try again.',
    revoke_all_failed: 'Failed to revoke sessions. Please try again.',
    server: 'Something went wrong. Please try again.',
  };
  const displayError = $derived(error ? (errorMessages[error] ?? 'Something went wrong.') : null);

  function formatDateTime(d: Date | undefined): string {
    if (!d) return '';
    return new Date(d).toLocaleDateString(undefined, {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
    });
  }
</script>

<PageContainer>
  <h1>Sessions</h1>
  <p><Link href="/account/settings" class="back-link">Back to Account Settings</Link></p>

  {#if revoked}
    <Alert variant="info">Session revoked successfully.</Alert>
  {/if}
  {#if revokedAll}
    <Alert variant="info">All other sessions have been revoked.</Alert>
  {/if}
  {#if displayError}
    <Alert variant="error">{displayError}</Alert>
  {/if}

  <div class="sessions-content">
    {#if hasOtherSessions}
      <div class="revoke-all-section">
        <form method="POST" action="?/revoke-all" use:enhance>
          <Button type="submit" variant="default" class="revoke-all-btn">Revoke all other sessions</Button>
        </form>
      </div>
    {/if}

    {#if sessions.length > 0}
      <div class="session-list">
        {#each sessions as session (session.id)}
          <div class="session-card" class:session-current={session.isCurrent}>
            <div class="session-info">
              <div class="session-header">
                <span class="session-device">{session.deviceName || 'Unknown device'}</span>
                {#if session.isCurrent}
                  <span class="current-badge">Current session</span>
                {/if}
              </div>
              <div class="session-details">
                {#if session.loginMethod}
                  <span>Login: {session.loginMethod}</span>
                {/if}
                {#if session.clientIp}
                  <span>IP: {session.clientIp}</span>
                {/if}
              </div>
              {#if session.userAgent}
                <div class="session-ua">{session.userAgent}</div>
              {/if}
              <div class="session-timestamps">
                {#if session.createdAt}
                  <span>Created {formatDateTime(session.createdAt)}</span>
                {/if}
                {#if session.lastActiveAt}
                  <span>Last active {formatDateTime(session.lastActiveAt)}</span>
                {/if}
                {#if session.expiresAt}
                  <span>Expires {formatDateTime(session.expiresAt)}</span>
                {/if}
              </div>
            </div>
            {#if !session.isCurrent}
              <form method="POST" action="?/revoke" use:enhance class="revoke-form">
                <input type="hidden" name="session_id" value={session.id} />
                <Button type="submit" variant="default" class="revoke-btn">Revoke</Button>
              </form>
            {/if}
          </div>
        {/each}
      </div>
    {:else}
      <p class="muted">No active sessions found.</p>
    {/if}
  </div>
</PageContainer>

<style>
  .back-link {
    font-size: 0.875rem;
  }
  .sessions-content {
    display: flex;
    flex-direction: column;
    gap: var(--space-lg);
    margin-top: var(--space-lg);
  }
  .revoke-all-section {
    display: flex;
  }
  .revoke-all-btn {
    font-size: 0.875rem;
    color: var(--color-error, #c00);
    border-color: var(--color-error, #c00);
  }
  .session-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-sm);
  }
  .session-card {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: var(--space-md);
    padding: var(--space-md);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-surface);
  }
  .session-current {
    border-left: 3px solid var(--color-primary, #c45c3e);
  }
  .session-info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  .session-header {
    display: flex;
    align-items: center;
    gap: var(--space-sm);
  }
  .session-device {
    font-weight: var(--font-weight-medium);
  }
  .current-badge {
    font-size: 0.75rem;
    padding: 0.125rem 0.5rem;
    border-radius: var(--radius-sm);
    background: var(--color-primary, #c45c3e);
    color: #fff;
  }
  .session-details {
    display: flex;
    gap: var(--space-md);
    font-size: 0.875rem;
    color: var(--color-muted, #666);
  }
  .session-ua {
    font-size: 0.8125rem;
    color: var(--color-muted, #666);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 500px;
  }
  .session-timestamps {
    display: flex;
    gap: var(--space-md);
    font-size: 0.8125rem;
    color: var(--color-muted, #666);
  }
  .revoke-form {
    flex-shrink: 0;
  }
  .revoke-btn {
    font-size: 0.875rem;
    padding: 0.25rem 0.5rem;
    color: var(--color-error, #c00);
    border-color: var(--color-error, #c00);
  }
  .muted {
    font-size: 0.875rem;
    color: var(--color-muted, #666);
  }
</style>
