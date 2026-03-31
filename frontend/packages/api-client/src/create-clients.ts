/**
 * gRPC client factory. Creates client singletons and promisified API from config.
 * Per-request auth is passed via call-level metadata, not per-client credentials.
 */

import * as grpc from '@grpc/grpc-js';
import { AuthServiceClient } from './auth/auth_service.js';
import type {
  LoginForTokenRequest,
  LoginForTokenResponse,
  GetSelfRequest,
  GetSelfResponse,
  GetActiveAccountRequest,
  GetActiveAccountResponse,
  RequestPasswordResetTokenRequest,
  RequestPasswordResetTokenResponse,
  RedeemPasswordResetTokenRequest,
  RedeemPasswordResetTokenResponse,
  VerifyEmailAddressRequest,
  VerifyEmailAddressResponse,
  ListPasskeysRequest,
  ListPasskeysResponse,
  ArchivePasskeyRequest,
  ArchivePasskeyResponse,
  BeginPasskeyRegistrationRequest,
  BeginPasskeyRegistrationResponse,
  FinishPasskeyRegistrationRequest,
  FinishPasskeyRegistrationResponse,
  BeginPasskeyAuthenticationRequest,
  BeginPasskeyAuthenticationResponse,
  FinishPasskeyAuthenticationRequest,
  ListActiveSessionsRequest,
  ListActiveSessionsResponse,
  RevokeSessionRequest,
  RevokeSessionResponse,
  RevokeAllOtherSessionsRequest,
  RevokeAllOtherSessionsResponse,
  RevokeCurrentSessionRequest,
  RevokeCurrentSessionResponse,
  ExchangeTokenRequest,
  ExchangeTokenResponse,
} from './auth/auth_service_types.js';
import { IdentityServiceClient } from './identity/identity_service.js';
import type {
  GetAccountsForUserRequest,
  GetAccountsForUserResponse,
  GetSentAccountInvitationsRequest,
  GetSentAccountInvitationsResponse,
  CreateAccountInvitationRequest,
  CreateAccountInvitationResponse,
  CancelAccountInvitationRequest,
  CancelAccountInvitationResponse,
  UpdateAccountMemberPermissionsRequest,
  UpdateAccountMemberPermissionsResponse,
  UpdateAccountRequest,
  UpdateAccountResponse,
  UpdateUserUsernameRequest,
  UpdateUserUsernameResponse,
  UpdateUserDetailsRequest,
  UpdateUserDetailsResponse,
  UploadUserAvatarResponse,
} from './identity/identity_service_types.js';
import { UploadRequest, UploadMetadata } from './uploaded_media/uploaded_media_messages.js';
import { SettingsServiceClient } from './settings/settings_service.js';
import type {
  GetServiceSettingsRequest,
  GetServiceSettingsResponse,
  GetServiceSettingConfigurationsForUserRequest,
  GetServiceSettingConfigurationsForUserResponse,
  CreateServiceSettingConfigurationRequest,
  CreateServiceSettingConfigurationResponse,
  UpdateServiceSettingConfigurationRequest,
  UpdateServiceSettingConfigurationResponse,
} from './settings/settings_service_types.js';
import { AnalyticsServiceClient } from './analytics/analytics_service.js';
import type { Metadata } from '@grpc/grpc-js';

export interface GrpcClientConfig {
  serverUrl: string;
  insecure?: boolean;
}

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

const UPLOAD_AVATAR_CHUNK_SIZE = 64 * 1024; // 64 KB, match iOS

export function createGrpcClients(config: GrpcClientConfig) {
  const credentials = config.insecure ? grpc.credentials.createInsecure() : grpc.credentials.createSsl();
  const serverUrl = config.serverUrl;

  let authClient: AuthServiceClient | null = null;
  let identityClient: IdentityServiceClient | null = null;
  let settingsClient: SettingsServiceClient | null = null;
  let analyticsClient: AnalyticsServiceClient | null = null;

  function getAuthClient(): AuthServiceClient {
    if (!authClient) authClient = new AuthServiceClient(serverUrl, credentials);
    return authClient;
  }
  function getIdentityClient(): IdentityServiceClient {
    if (!identityClient) identityClient = new IdentityServiceClient(serverUrl, credentials);
    return identityClient;
  }
  function getSettingsClient(): SettingsServiceClient {
    if (!settingsClient) settingsClient = new SettingsServiceClient(serverUrl, credentials);
    return settingsClient;
  }
  function getAnalyticsClient(): AnalyticsServiceClient {
    if (!analyticsClient) analyticsClient = new AnalyticsServiceClient(serverUrl, credentials);
    return analyticsClient;
  }

  function authMetadata(oauth2AccessToken: string): Metadata {
    const m = new grpc.Metadata();
    m.add('authorization', `Bearer ${oauth2AccessToken}`);
    return m;
  }
  const emptyMetadata = new grpc.Metadata();

  return {
    authMetadata,

    loginForToken: (request: LoginForTokenRequest): Promise<LoginForTokenResponse> =>
      promisifyUnary<LoginForTokenRequest, LoginForTokenResponse>(getAuthClient().loginForToken.bind(getAuthClient()))(
        request,
        emptyMetadata,
      ),

    requestPasswordResetToken: (
      request: RequestPasswordResetTokenRequest,
    ): Promise<RequestPasswordResetTokenResponse> =>
      promisifyUnary<RequestPasswordResetTokenRequest, RequestPasswordResetTokenResponse>(
        getAuthClient().requestPasswordResetToken.bind(getAuthClient()),
      )(request, emptyMetadata),

    redeemPasswordResetToken: (request: RedeemPasswordResetTokenRequest): Promise<RedeemPasswordResetTokenResponse> =>
      promisifyUnary<RedeemPasswordResetTokenRequest, RedeemPasswordResetTokenResponse>(
        getAuthClient().redeemPasswordResetToken.bind(getAuthClient()),
      )(request, emptyMetadata),

    verifyEmailAddress: (request: VerifyEmailAddressRequest): Promise<VerifyEmailAddressResponse> =>
      promisifyUnary<VerifyEmailAddressRequest, VerifyEmailAddressResponse>(
        getAuthClient().verifyEmailAddress.bind(getAuthClient()),
      )(request, emptyMetadata),

    beginPasskeyRegistration: (oauth2Token: string): Promise<BeginPasskeyRegistrationResponse> =>
      promisifyUnary<BeginPasskeyRegistrationRequest, BeginPasskeyRegistrationResponse>(
        getAuthClient().beginPasskeyRegistration.bind(getAuthClient()),
      )({}, authMetadata(oauth2Token)),

    finishPasskeyRegistration: (
      oauth2Token: string,
      request: FinishPasskeyRegistrationRequest,
    ): Promise<FinishPasskeyRegistrationResponse> =>
      promisifyUnary<FinishPasskeyRegistrationRequest, FinishPasskeyRegistrationResponse>(
        getAuthClient().finishPasskeyRegistration.bind(getAuthClient()),
      )(request, authMetadata(oauth2Token)),

    beginPasskeyAuthentication: (
      request: BeginPasskeyAuthenticationRequest,
    ): Promise<BeginPasskeyAuthenticationResponse> =>
      promisifyUnary<BeginPasskeyAuthenticationRequest, BeginPasskeyAuthenticationResponse>(
        getAuthClient().beginPasskeyAuthentication.bind(getAuthClient()),
      )(request, emptyMetadata),

    finishPasskeyAuthentication: (request: FinishPasskeyAuthenticationRequest): Promise<LoginForTokenResponse> =>
      promisifyUnary<FinishPasskeyAuthenticationRequest, LoginForTokenResponse>(
        getAuthClient().finishPasskeyAuthentication.bind(getAuthClient()),
      )(request, emptyMetadata),

    getSelf: (oauth2Token: string): Promise<GetSelfResponse> =>
      promisifyUnary<GetSelfRequest, GetSelfResponse>(getAuthClient().getSelf.bind(getAuthClient()))(
        {},
        authMetadata(oauth2Token),
      ),

    getActiveAccount: (oauth2Token: string): Promise<GetActiveAccountResponse> =>
      promisifyUnary<GetActiveAccountRequest, GetActiveAccountResponse>(
        getAuthClient().getActiveAccount.bind(getAuthClient()),
      )({}, authMetadata(oauth2Token)),

    listPasskeys: (oauth2Token: string): Promise<ListPasskeysResponse> =>
      promisifyUnary<ListPasskeysRequest, ListPasskeysResponse>(getAuthClient().listPasskeys.bind(getAuthClient()))(
        {},
        authMetadata(oauth2Token),
      ),

    archivePasskey: (oauth2Token: string, request: ArchivePasskeyRequest): Promise<ArchivePasskeyResponse> =>
      promisifyUnary<ArchivePasskeyRequest, ArchivePasskeyResponse>(
        getAuthClient().archivePasskey.bind(getAuthClient()),
      )(request, authMetadata(oauth2Token)),

    listActiveSessions: (oauth2Token: string): Promise<ListActiveSessionsResponse> =>
      promisifyUnary<ListActiveSessionsRequest, ListActiveSessionsResponse>(
        getAuthClient().listActiveSessions.bind(getAuthClient()),
      )({ filter: undefined }, authMetadata(oauth2Token)),

    revokeSession: (oauth2Token: string, request: RevokeSessionRequest): Promise<RevokeSessionResponse> =>
      promisifyUnary<RevokeSessionRequest, RevokeSessionResponse>(getAuthClient().revokeSession.bind(getAuthClient()))(
        request,
        authMetadata(oauth2Token),
      ),

    revokeAllOtherSessions: (oauth2Token: string): Promise<RevokeAllOtherSessionsResponse> =>
      promisifyUnary<RevokeAllOtherSessionsRequest, RevokeAllOtherSessionsResponse>(
        getAuthClient().revokeAllOtherSessions.bind(getAuthClient()),
      )({}, authMetadata(oauth2Token)),

    revokeCurrentSession: (token: string): Promise<RevokeCurrentSessionResponse> =>
      promisifyUnary<RevokeCurrentSessionRequest, RevokeCurrentSessionResponse>(
        getAuthClient().revokeCurrentSession.bind(getAuthClient()),
      )({}, authMetadata(token)),

    exchangeToken: (refreshToken: string, desiredAccountId?: string): Promise<ExchangeTokenResponse> =>
      promisifyUnary<ExchangeTokenRequest, ExchangeTokenResponse>(getAuthClient().exchangeToken.bind(getAuthClient()))(
        { refreshToken, desiredAccountId: desiredAccountId ?? '' },
        emptyMetadata,
      ),

    getAccountsForUser: (
      oauth2Token: string,
      request: GetAccountsForUserRequest,
    ): Promise<GetAccountsForUserResponse> =>
      promisifyUnary<GetAccountsForUserRequest, GetAccountsForUserResponse>(
        getIdentityClient().getAccountsForUser.bind(getIdentityClient()),
      )(request, authMetadata(oauth2Token)),

    getSentAccountInvitations: (
      oauth2Token: string,
      request: GetSentAccountInvitationsRequest,
    ): Promise<GetSentAccountInvitationsResponse> =>
      promisifyUnary<GetSentAccountInvitationsRequest, GetSentAccountInvitationsResponse>(
        getIdentityClient().getSentAccountInvitations.bind(getIdentityClient()),
      )(request, authMetadata(oauth2Token)),

    createAccountInvitation: (
      oauth2Token: string,
      request: CreateAccountInvitationRequest,
    ): Promise<CreateAccountInvitationResponse> =>
      promisifyUnary<CreateAccountInvitationRequest, CreateAccountInvitationResponse>(
        getIdentityClient().createAccountInvitation.bind(getIdentityClient()),
      )(request, authMetadata(oauth2Token)),

    cancelAccountInvitation: (
      oauth2Token: string,
      request: CancelAccountInvitationRequest,
    ): Promise<CancelAccountInvitationResponse> =>
      promisifyUnary<CancelAccountInvitationRequest, CancelAccountInvitationResponse>(
        getIdentityClient().cancelAccountInvitation.bind(getIdentityClient()),
      )(request, authMetadata(oauth2Token)),

    updateAccountMemberPermissions: (
      oauth2Token: string,
      request: UpdateAccountMemberPermissionsRequest,
    ): Promise<UpdateAccountMemberPermissionsResponse> =>
      promisifyUnary<UpdateAccountMemberPermissionsRequest, UpdateAccountMemberPermissionsResponse>(
        getIdentityClient().updateAccountMemberPermissions.bind(getIdentityClient()),
      )(request, authMetadata(oauth2Token)),

    updateAccount: (oauth2Token: string, request: UpdateAccountRequest): Promise<UpdateAccountResponse> =>
      promisifyUnary<UpdateAccountRequest, UpdateAccountResponse>(
        getIdentityClient().updateAccount.bind(getIdentityClient()),
      )(request, authMetadata(oauth2Token)),

    updateUserUsername: (
      oauth2Token: string,
      request: UpdateUserUsernameRequest,
    ): Promise<UpdateUserUsernameResponse> =>
      promisifyUnary<UpdateUserUsernameRequest, UpdateUserUsernameResponse>(
        getIdentityClient().updateUserUsername.bind(getIdentityClient()),
      )(request, authMetadata(oauth2Token)),

    updateUserDetails: (oauth2Token: string, request: UpdateUserDetailsRequest): Promise<UpdateUserDetailsResponse> =>
      promisifyUnary<UpdateUserDetailsRequest, UpdateUserDetailsResponse>(
        getIdentityClient().updateUserDetails.bind(getIdentityClient()),
      )(request, authMetadata(oauth2Token)),

    uploadUserAvatar: async (
      oauth2Token: string,
      fileBuffer: Buffer,
      filename: string,
      contentType: string,
    ): Promise<UploadUserAvatarResponse> => {
      const client = getIdentityClient();
      const meta = authMetadata(oauth2Token);
      return new Promise((resolve, reject) => {
        const stream = client.uploadUserAvatar(meta, (err, response) => {
          if (err) reject(err);
          else if (response) resolve(response);
          else reject(new Error('No response'));
        });
        const metadataReq = UploadRequest.create({
          metadata: UploadMetadata.create({
            bucket: 'avatars',
            objectName: filename,
            contentType,
          }),
        });
        stream.write(metadataReq, (writeErr: unknown) => {
          if (writeErr) {
            const err = writeErr instanceof Error ? writeErr : new Error(String(writeErr));
            stream.destroy(err);
            reject(err);
            return;
          }
          let offset = 0;
          const writeNext = () => {
            if (offset >= fileBuffer.length) {
              stream.end();
              return;
            }
            const end = Math.min(offset + UPLOAD_AVATAR_CHUNK_SIZE, fileBuffer.length);
            const chunk = fileBuffer.subarray(offset, end);
            offset = end;
            stream.write(UploadRequest.create({ chunk: new Uint8Array(chunk) }), (chunkErr: unknown) => {
              if (chunkErr) {
                const err = chunkErr instanceof Error ? chunkErr : new Error(String(chunkErr));
                stream.destroy(err);
                reject(err);
                return;
              }
              writeNext();
            });
          };
          writeNext();
        });
      });
    },

    getServiceSettings: (
      oauth2Token: string,
      request: GetServiceSettingsRequest,
    ): Promise<GetServiceSettingsResponse> =>
      promisifyUnary<GetServiceSettingsRequest, GetServiceSettingsResponse>(
        getSettingsClient().getServiceSettings.bind(getSettingsClient()),
      )(request, authMetadata(oauth2Token)),

    getServiceSettingConfigurationsForUser: (
      oauth2Token: string,
      request: GetServiceSettingConfigurationsForUserRequest,
    ): Promise<GetServiceSettingConfigurationsForUserResponse> =>
      promisifyUnary<GetServiceSettingConfigurationsForUserRequest, GetServiceSettingConfigurationsForUserResponse>(
        getSettingsClient().getServiceSettingConfigurationsForUser.bind(getSettingsClient()),
      )(request, authMetadata(oauth2Token)),

    createServiceSettingConfiguration: (
      oauth2Token: string,
      request: CreateServiceSettingConfigurationRequest,
    ): Promise<CreateServiceSettingConfigurationResponse> =>
      promisifyUnary<CreateServiceSettingConfigurationRequest, CreateServiceSettingConfigurationResponse>(
        getSettingsClient().createServiceSettingConfiguration.bind(getSettingsClient()),
      )(request, authMetadata(oauth2Token)),

    updateServiceSettingConfiguration: (
      oauth2Token: string,
      request: UpdateServiceSettingConfigurationRequest,
    ): Promise<UpdateServiceSettingConfigurationResponse> =>
      promisifyUnary<UpdateServiceSettingConfigurationRequest, UpdateServiceSettingConfigurationResponse>(
        getSettingsClient().updateServiceSettingConfiguration.bind(getSettingsClient()),
      )(request, authMetadata(oauth2Token)),

    trackEvent: async (event: string, properties: Record<string, string> = {}): Promise<void> => {
      await promisifyUnary(getAnalyticsClient().trackEvent.bind(getAnalyticsClient()))(
        { source: 'web', event, properties },
        emptyMetadata,
      );
    },

    trackAnonymousEvent: async (
      event: string,
      anonymousId: string,
      properties: Record<string, string> = {},
    ): Promise<void> => {
      await promisifyUnary(getAnalyticsClient().trackAnonymousEvent.bind(getAnalyticsClient()))(
        { source: 'web', event, anonymousId, properties },
        emptyMetadata,
      );
    },
  };
}
