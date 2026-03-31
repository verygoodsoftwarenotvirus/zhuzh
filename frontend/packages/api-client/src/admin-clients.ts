/**
 * Admin gRPC client factory. Uses AdminLoginForToken (no OAuth exchange) and exposes
 * all admin-required methods with token-first signature.
 */

import * as grpc from '@grpc/grpc-js';
import type { Metadata } from '@grpc/grpc-js';
import { AuthServiceClient } from './auth/auth_service.js';
import type {
  AdminLoginForTokenRequest,
  AdminListSessionsForUserRequest,
  AdminRevokeUserSessionRequest,
  AdminRevokeAllUserSessionsRequest,
} from './auth/auth_service_types.js';
import type {
  LoginForTokenResponse,
  ListActiveSessionsResponse,
  RevokeSessionResponse,
  RevokeAllOtherSessionsResponse,
  RevokeCurrentSessionRequest,
  RevokeCurrentSessionResponse,
} from './auth/auth_service_types.js';
import { IdentityServiceClient } from './identity/identity_service.js';
import { OAuthServiceClient } from './oauth/oauth_service.js';
import { PaymentsServiceClient } from './payments/payments_service.js';
import { SettingsServiceClient } from './settings/settings_service.js';
import { WaitlistsServiceClient } from './waitlists/waitlists_service.js';
import { IssueReportsServiceClient } from './issue_reports/issue_reports_service.js';
import { InternalOperationsClient } from './internal_ops/internal_ops_service.js';
import { AuditServiceClient } from './audit/audit_service.js';
import { AnalyticsServiceClient } from './analytics/analytics_service.js';
import type { GrpcClientConfig } from './create-clients.js';

function promisifyUnary<TRequest, TResponse>(
  call: (
    req: TRequest,
    metadata: Metadata,
    callback: (err: grpc.ServiceError | null, res: TResponse) => void,
  ) => grpc.ClientUnaryCall,
): (req: TRequest, metadata: Metadata) => Promise<TResponse> {
  return (req, metadata) =>
    new Promise((resolve, reject) => {
      call(req, metadata, (err, res) => {
        if (err) reject(err);
        else if (res) resolve(res);
        else reject(new Error('No response'));
      });
    });
}

export function createAdminGrpcClients(config: GrpcClientConfig) {
  const credentials = config.insecure ? grpc.credentials.createInsecure() : grpc.credentials.createSsl();
  const serverUrl = config.serverUrl;

  let authClient: AuthServiceClient | null = null;
  let identityClient: IdentityServiceClient | null = null;
  let oauthClient: OAuthServiceClient | null = null;
  let paymentsClient: PaymentsServiceClient | null = null;
  let settingsClient: SettingsServiceClient | null = null;
  let waitlistsClient: WaitlistsServiceClient | null = null;
  let issueReportsClient: IssueReportsServiceClient | null = null;
  let internalOpsClient: InternalOperationsClient | null = null;
  let auditClient: AuditServiceClient | null = null;
  let analyticsClient: AnalyticsServiceClient | null = null;

  const get = {
    auth: () => {
      if (!authClient) authClient = new AuthServiceClient(serverUrl, credentials);
      return authClient;
    },
    identity: () => {
      if (!identityClient) identityClient = new IdentityServiceClient(serverUrl, credentials);
      return identityClient;
    },
    oauth: () => {
      if (!oauthClient) oauthClient = new OAuthServiceClient(serverUrl, credentials);
      return oauthClient;
    },
    payments: () => {
      if (!paymentsClient) paymentsClient = new PaymentsServiceClient(serverUrl, credentials);
      return paymentsClient;
    },
    settings: () => {
      if (!settingsClient) settingsClient = new SettingsServiceClient(serverUrl, credentials);
      return settingsClient;
    },
    waitlists: () => {
      if (!waitlistsClient) waitlistsClient = new WaitlistsServiceClient(serverUrl, credentials);
      return waitlistsClient;
    },
    issueReports: () => {
      if (!issueReportsClient) issueReportsClient = new IssueReportsServiceClient(serverUrl, credentials);
      return issueReportsClient;
    },
    internalOps: () => {
      if (!internalOpsClient) internalOpsClient = new InternalOperationsClient(serverUrl, credentials);
      return internalOpsClient;
    },
    audit: () => {
      if (!auditClient) auditClient = new AuditServiceClient(serverUrl, credentials);
      return auditClient;
    },
    analytics: () => {
      if (!analyticsClient) analyticsClient = new AnalyticsServiceClient(serverUrl, credentials);
      return analyticsClient;
    },
  };

  function authMetadata(token: string): Metadata {
    const m = new grpc.Metadata();
    m.add('authorization', `Bearer ${token}`);
    return m;
  }
  const emptyMetadata = new grpc.Metadata();

  return {
    authMetadata,

    adminLoginForToken: (request: AdminLoginForTokenRequest): Promise<LoginForTokenResponse> =>
      promisifyUnary<AdminLoginForTokenRequest, LoginForTokenResponse>(get.auth().adminLoginForToken.bind(get.auth()))(
        request,
        emptyMetadata,
      ),

    beginPasskeyAuthentication: (request: { username?: string }) =>
      promisifyUnary(get.auth().beginPasskeyAuthentication.bind(get.auth()))(
        { username: request.username ?? '' },
        emptyMetadata,
      ),
    finishPasskeyAuthentication: (request: { challenge: string; username: string; assertionResponse: Uint8Array }) =>
      promisifyUnary(get.auth().finishPasskeyAuthentication.bind(get.auth()))(request, emptyMetadata),

    adminListSessionsForUser: (
      token: string,
      request: AdminListSessionsForUserRequest,
    ): Promise<ListActiveSessionsResponse> =>
      promisifyUnary<AdminListSessionsForUserRequest, ListActiveSessionsResponse>(
        get.auth().adminListSessionsForUser.bind(get.auth()),
      )(request, authMetadata(token)),

    adminRevokeUserSession: (token: string, request: AdminRevokeUserSessionRequest): Promise<RevokeSessionResponse> =>
      promisifyUnary<AdminRevokeUserSessionRequest, RevokeSessionResponse>(
        get.auth().adminRevokeUserSession.bind(get.auth()),
      )(request, authMetadata(token)),

    adminRevokeAllUserSessions: (
      token: string,
      request: AdminRevokeAllUserSessionsRequest,
    ): Promise<RevokeAllOtherSessionsResponse> =>
      promisifyUnary<AdminRevokeAllUserSessionsRequest, RevokeAllOtherSessionsResponse>(
        get.auth().adminRevokeAllUserSessions.bind(get.auth()),
      )(request, authMetadata(token)),

    revokeCurrentSession: (token: string): Promise<RevokeCurrentSessionResponse> =>
      promisifyUnary<RevokeCurrentSessionRequest, RevokeCurrentSessionResponse>(
        get.auth().revokeCurrentSession.bind(get.auth()),
      )({}, authMetadata(token)),

    // Identity – request types are intentionally loose; callers pass proto-shaped objects
    getUser: (token: string, request: { userId: string }) =>
      promisifyUnary(get.identity().getUser.bind(get.identity()))(request as any, authMetadata(token)),
    getUsers: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.identity().getUsers.bind(get.identity()))(
        request as unknown as Parameters<IdentityServiceClient['getUsers']>[0],
        authMetadata(token),
      ),
    getAccount: (token: string, request: { accountId: string }) =>
      promisifyUnary(get.identity().getAccount.bind(get.identity()))(request as any, authMetadata(token)),
    getAccounts: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.identity().getAccounts.bind(get.identity()))(
        request as unknown as Parameters<IdentityServiceClient['getAccounts']>[0],
        authMetadata(token),
      ),
    searchForUsers: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.identity().searchForUsers.bind(get.identity()))(
        request as unknown as Parameters<IdentityServiceClient['searchForUsers']>[0],
        authMetadata(token),
      ),
    getUsersForAccount: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.identity().getUsersForAccount.bind(get.identity()))(
        request as unknown as Parameters<IdentityServiceClient['getUsersForAccount']>[0],
        authMetadata(token),
      ),
    getAccountsForUser: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.identity().getAccountsForUser.bind(get.identity()))(
        request as unknown as Parameters<IdentityServiceClient['getAccountsForUser']>[0],
        authMetadata(token),
      ),
    adminUpdateUserStatus: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.identity().adminUpdateUserStatus.bind(get.identity()))(
        request as unknown as Parameters<IdentityServiceClient['adminUpdateUserStatus']>[0],
        authMetadata(token),
      ),
    adminSetPasswordChangeRequired: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.identity().adminSetPasswordChangeRequired.bind(get.identity()))(
        request as unknown as Parameters<IdentityServiceClient['adminSetPasswordChangeRequired']>[0],
        authMetadata(token),
      ),
    updateUserDetails: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.identity().updateUserDetails.bind(get.identity()))(
        request as unknown as Parameters<IdentityServiceClient['updateUserDetails']>[0],
        authMetadata(token),
      ),
    updateAccount: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.identity().updateAccount.bind(get.identity()))(
        request as unknown as Parameters<IdentityServiceClient['updateAccount']>[0],
        authMetadata(token),
      ),

    // OAuth
    getOAuth2Clients: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.oauth().getOAuth2Clients.bind(get.oauth()))(
        request as unknown as Parameters<OAuthServiceClient['getOAuth2Clients']>[0],
        authMetadata(token),
      ),
    getOAuth2Client: (token: string, request: { oauth2ClientId: string }) =>
      promisifyUnary(get.oauth().getOAuth2Client.bind(get.oauth()))(request as any, authMetadata(token)),
    createOAuth2Client: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.oauth().createOAuth2Client.bind(get.oauth()))(
        request as unknown as Parameters<OAuthServiceClient['createOAuth2Client']>[0],
        authMetadata(token),
      ),
    archiveOAuth2Client: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.oauth().archiveOAuth2Client.bind(get.oauth()))(
        request as unknown as Parameters<OAuthServiceClient['archiveOAuth2Client']>[0],
        authMetadata(token),
      ),

    // Payments
    getProducts: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.payments().getProducts.bind(get.payments()))(
        request as unknown as Parameters<PaymentsServiceClient['getProducts']>[0],
        authMetadata(token),
      ),
    getProduct: (token: string, request: { productId: string }) =>
      promisifyUnary(get.payments().getProduct.bind(get.payments()))(request as any, authMetadata(token)),
    createProduct: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.payments().createProduct.bind(get.payments()))(
        request as unknown as Parameters<PaymentsServiceClient['createProduct']>[0],
        authMetadata(token),
      ),
    updateProduct: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.payments().updateProduct.bind(get.payments()))(
        request as unknown as Parameters<PaymentsServiceClient['updateProduct']>[0],
        authMetadata(token),
      ),
    archiveProduct: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.payments().archiveProduct.bind(get.payments()))(
        request as unknown as Parameters<PaymentsServiceClient['archiveProduct']>[0],
        authMetadata(token),
      ),
    getSubscription: (token: string, request: { subscriptionId: string }) =>
      promisifyUnary(get.payments().getSubscription.bind(get.payments()))(request as any, authMetadata(token)),
    createSubscription: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.payments().createSubscription.bind(get.payments()))(
        request as unknown as Parameters<PaymentsServiceClient['createSubscription']>[0],
        authMetadata(token),
      ),
    updateSubscription: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.payments().updateSubscription.bind(get.payments()))(
        request as unknown as Parameters<PaymentsServiceClient['updateSubscription']>[0],
        authMetadata(token),
      ),
    archiveSubscription: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.payments().archiveSubscription.bind(get.payments()))(
        request as unknown as Parameters<PaymentsServiceClient['archiveSubscription']>[0],
        authMetadata(token),
      ),
    getSubscriptionsForAccount: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.payments().getSubscriptionsForAccount.bind(get.payments()))(
        request as unknown as Parameters<PaymentsServiceClient['getSubscriptionsForAccount']>[0],
        authMetadata(token),
      ),

    // Settings
    getServiceSettings: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.settings().getServiceSettings.bind(get.settings()))(
        request as unknown as Parameters<SettingsServiceClient['getServiceSettings']>[0],
        authMetadata(token),
      ),
    searchForServiceSettings: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.settings().searchForServiceSettings.bind(get.settings()))(
        request as unknown as Parameters<SettingsServiceClient['searchForServiceSettings']>[0],
        authMetadata(token),
      ),
    getServiceSetting: (token: string, request: { serviceSettingId: string }) =>
      promisifyUnary(get.settings().getServiceSetting.bind(get.settings()))(request as any, authMetadata(token)),
    createServiceSetting: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.settings().createServiceSetting.bind(get.settings()))(
        request as unknown as Parameters<SettingsServiceClient['createServiceSetting']>[0],
        authMetadata(token),
      ),
    archiveServiceSetting: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.settings().archiveServiceSetting.bind(get.settings()))(
        request as unknown as Parameters<SettingsServiceClient['archiveServiceSetting']>[0],
        authMetadata(token),
      ),

    // Waitlists
    getWaitlists: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.waitlists().getWaitlists.bind(get.waitlists()))(
        request as unknown as Parameters<WaitlistsServiceClient['getWaitlists']>[0],
        authMetadata(token),
      ),
    getWaitlist: (token: string, request: { waitlistId: string }) =>
      promisifyUnary(get.waitlists().getWaitlist.bind(get.waitlists()))(request as any, authMetadata(token)),
    getWaitlistSignupsForWaitlist: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.waitlists().getWaitlistSignupsForWaitlist.bind(get.waitlists()))(
        request as unknown as Parameters<WaitlistsServiceClient['getWaitlistSignupsForWaitlist']>[0],
        authMetadata(token),
      ),

    // Issue Reports
    getIssueReports: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.issueReports().getIssueReports.bind(get.issueReports()))(
        request as unknown as Parameters<IssueReportsServiceClient['getIssueReports']>[0],
        authMetadata(token),
      ),
    getIssueReport: (token: string, request: { issueReportId: string }) =>
      promisifyUnary(get.issueReports().getIssueReport.bind(get.issueReports()))(request as any, authMetadata(token)),
    updateIssueReport: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.issueReports().updateIssueReport.bind(get.issueReports()))(
        request as unknown as Parameters<IssueReportsServiceClient['updateIssueReport']>[0],
        authMetadata(token),
      ),
    archiveIssueReport: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.issueReports().archiveIssueReport.bind(get.issueReports()))(
        request as unknown as Parameters<IssueReportsServiceClient['archiveIssueReport']>[0],
        authMetadata(token),
      ),

    // Internal Ops
    testQueueMessage: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.internalOps().testQueueMessage.bind(get.internalOps()))(
        request as unknown as Parameters<InternalOperationsClient['testQueueMessage']>[0],
        authMetadata(token),
      ),

    // Audit
    getAuditLogEntriesForUser: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.audit().getAuditLogEntriesForUser.bind(get.audit()))(
        request as unknown as Parameters<AuditServiceClient['getAuditLogEntriesForUser']>[0],
        authMetadata(token),
      ),
    getAuditLogEntriesForAccount: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.audit().getAuditLogEntriesForAccount.bind(get.audit()))(
        request as unknown as Parameters<AuditServiceClient['getAuditLogEntriesForAccount']>[0],
        authMetadata(token),
      ),

    // Analytics
    trackEvent: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.analytics().trackEvent.bind(get.analytics()))(
        request as unknown as Parameters<AnalyticsServiceClient['trackEvent']>[0],
        authMetadata(token),
      ),
    trackAnonymousEvent: (token: string, request: Record<string, unknown>) =>
      promisifyUnary(get.analytics().trackAnonymousEvent.bind(get.analytics()))(
        request as unknown as Parameters<AnalyticsServiceClient['trackAnonymousEvent']>[0],
        authMetadata(token),
      ),

  };
}
