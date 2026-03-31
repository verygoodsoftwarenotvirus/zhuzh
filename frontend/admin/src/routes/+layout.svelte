<script lang="ts">
  import '../app.css';
  import { Link, PageContainer } from '@zhuzh/ui';

  let { children } = $props();

  let mobileMenuOpen = $state(false);
</script>

<div class="layout-root">
  <header class="layout-header">
    <div class="container">
      <div><Link href="/">Admin</Link></div>
      <nav class="layout-links" class:open={mobileMenuOpen}>
        <Link href="/">Dashboard</Link>
        <Link href="/users">Users</Link>
        <Link href="/accounts">Accounts</Link>
        <Link href="/oauth2_clients">OAuth2 Clients</Link>
        <Link href="/waitlists">Waitlists</Link>
        <Link href="/issue_reports">Issue Reports</Link>
        <Link href="/settings">Settings</Link>
        <Link href="/products">Products</Link>
        <Link href="/subscriptions">Subscriptions</Link>
        <Link href="/queue_test">Queue Test</Link>
        <Link href="/analytics_test">Analytics Test</Link>
        <Link href="/logout">Logout</Link>
      </nav>
      <button
        type="button"
        class="mobile-menu-btn"
        aria-label="Toggle menu"
        onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
      >
        <span></span>
        <span></span>
        <span></span>
      </button>
    </div>
  </header>

  <main class="layout-main">
    <PageContainer>
      {@render children()}
    </PageContainer>
  </main>
</div>

<style>
  .layout-root {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }
  .layout-header {
    background: var(--color-surface);
    border-bottom: 1px solid var(--color-border);
    padding: var(--space-md);
  }
  .container {
    max-width: 64rem;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
  }
  .layout-links {
    display: flex;
    gap: var(--space-md);
    align-items: center;
  }
  .mobile-menu-btn {
    display: none;
    flex-direction: column;
    gap: 4px;
    background: none;
    border: none;
    cursor: pointer;
    padding: var(--space-sm);
  }
  .mobile-menu-btn span {
    width: 1.25rem;
    height: 2px;
    background: var(--color-text);
  }
  @media (max-width: 768px) {
    .mobile-menu-btn {
      display: flex;
    }
    .layout-links {
      display: none;
      width: 100%;
      flex-direction: column;
      align-items: flex-start;
      padding-top: var(--space-md);
    }
    .layout-links.open {
      display: flex;
    }
  }
  .layout-main {
    flex: 1;
    padding: var(--space-lg);
  }
</style>
