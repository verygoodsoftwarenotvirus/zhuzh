<script lang="ts">
  import { enhance } from '$app/forms';
  import { Heading, Button, Alert, Link, Card } from '@zhuzh/ui';
  import type { UserSession } from '@zhuzh/api-client/auth/auth_messages';

  let { data } = $props();
  const sessions = $derived((data?.sessions ?? []) as UserSession[]);
  const error = $derived(data?.error as string | null | undefined);
  const revoked = $derived(data?.revoked ?? false);
  const revokedAll = $derived(data?.revokedAll ?? false);
  const userId = $derived(data?.userId as string);

  const errorMessages: Record<string, string> = {
    invalid: 'Invalid request.',
    revoke_failed: 'Failed to revoke session. Please try again.',
    revoke_all_failed: 'Failed to revoke sessions. Please try again.',
    server: 'Something went wrong. Please try again.',
  };
  const displayError = $derived(error ? (errorMessages[error] ?? 'Something went wrong.') : null);

  function formatDateTime(d: Date | undefined): string {
    if (!d) return '-';
    return new Date(d).toLocaleDateString(undefined, {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
    });
  }
</script>

<Heading level={1}>Sessions</Heading>
<p class="subtitle">User ID: {userId}</p>

{#if revoked}
  <Alert variant="info">Session revoked successfully.</Alert>
{/if}
{#if revokedAll}
  <Alert variant="info">All sessions have been revoked.</Alert>
{/if}
{#if displayError}
  <Alert variant="error">{displayError}</Alert>
{/if}

{#if sessions.length > 0}
  <div class="actions-bar">
    <form method="POST" action="?/revoke-all" use:enhance>
      <Button type="submit" variant="default" class="revoke-all-btn">Revoke all sessions</Button>
    </form>
  </div>

  <Card>
    <div class="table-wrap">
      <table class="data-table">
        <thead>
          <tr>
            <th>Device</th>
            <th>Login</th>
            <th>IP</th>
            <th>User Agent</th>
            <th>Created</th>
            <th>Last Active</th>
            <th>Expires</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each sessions as session (session.id)}
            <tr>
              <td>{session.deviceName || '-'}</td>
              <td>{session.loginMethod || '-'}</td>
              <td>{session.clientIp || '-'}</td>
              <td class="ua-cell">{session.userAgent || '-'}</td>
              <td>{formatDateTime(session.createdAt)}</td>
              <td>{formatDateTime(session.lastActiveAt)}</td>
              <td>{formatDateTime(session.expiresAt)}</td>
              <td>
                <form method="POST" action="?/revoke" use:enhance class="inline-form">
                  <input type="hidden" name="session_id" value={session.id} />
                  <Button type="submit" variant="default" class="revoke-btn">Revoke</Button>
                </form>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </Card>
{:else if !displayError}
  <Card><p class="muted">No active sessions found.</p></Card>
{/if}

<p><Link href="/users/{userId}">Back to user</Link></p>

<style>
  .subtitle {
    color: var(--color-text-muted);
    margin-bottom: var(--space-lg);
  }
  .actions-bar {
    margin-bottom: var(--space-md);
  }
  .table-wrap {
    overflow-x: auto;
  }
  .data-table {
    width: 100%;
    border-collapse: collapse;
  }
  .data-table th,
  .data-table td {
    padding: var(--space-sm) var(--space-md);
    text-align: left;
    border-bottom: 1px solid var(--color-border);
  }
  .ua-cell {
    max-width: 250px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .inline-form {
    display: inline;
  }
  .muted {
    color: var(--color-text-muted);
  }
</style>
