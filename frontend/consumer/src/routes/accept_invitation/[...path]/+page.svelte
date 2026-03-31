<script lang="ts">
  import { browser } from '$app/environment';
  import { PageContainer } from '@zhuzh/ui';
  import { COMPANY_NAME, APP_STORE_URL } from '@zhuzh/branding';

  let isIOS = $state(false);
  if (browser) {
    const ua = navigator.userAgent.toLowerCase();
    isIOS = ua.includes('iphone') || ua.includes('ipad');
  }
</script>

<PageContainer>
  <h1>Accept Invitation</h1>
  {#if isIOS}
    <p>
      To accept this invitation, open the {COMPANY_NAME} app. If you don't have the app, download it from the App Store.
    </p>
    <a href={APP_STORE_URL} class="app-link">Open in App Store</a>
  {:else}
    <p>
      This invitation link is for the {COMPANY_NAME} mobile app. Open this link on your iPhone or iPad to accept the invitation.
    </p>
    <a href={APP_STORE_URL} class="app-link">Get the App</a>
  {/if}
</PageContainer>

<style>
  .app-link {
    display: inline-block;
    margin-top: var(--space-md);
    padding: var(--space-sm) var(--space-lg);
    background: var(--color-text);
    color: var(--color-surface);
    text-decoration: none;
    border-radius: var(--radius-sm);
    font-weight: var(--font-weight-medium);
  }
  .app-link:hover {
    opacity: 0.9;
  }
</style>
