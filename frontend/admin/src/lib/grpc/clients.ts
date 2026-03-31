/**
 * Admin gRPC clients: reads env, creates clients from @zhuzh/api-client createAdminGrpcClients.
 */

import { env } from '$env/dynamic/private';
import { createAdminGrpcClients } from '@zhuzh/api-client';

const clients = createAdminGrpcClients({
  serverUrl: env.GRPC_API_SERVER_URL ?? 'localhost:50051',
  insecure: env.DEVELOPING_LOCALLY === 'true',
});

export const authMetadata = clients.authMetadata;
export const adminLoginForToken = clients.adminLoginForToken;
export const beginPasskeyAuthentication = clients.beginPasskeyAuthentication;
export const finishPasskeyAuthentication = clients.finishPasskeyAuthentication;
export const getUser = clients.getUser;
export const getUsers = clients.getUsers;
export const getAccount = clients.getAccount;
export const getAccounts = clients.getAccounts;
export const searchForUsers = clients.searchForUsers;
export const getUsersForAccount = clients.getUsersForAccount;
export const getAccountsForUser = clients.getAccountsForUser;
export const adminUpdateUserStatus = clients.adminUpdateUserStatus;
export const adminSetPasswordChangeRequired = clients.adminSetPasswordChangeRequired;
export const updateUserDetails = clients.updateUserDetails;
export const updateAccount = clients.updateAccount;
export const getOAuth2Clients = clients.getOAuth2Clients;
export const getOAuth2Client = clients.getOAuth2Client;
export const createOAuth2Client = clients.createOAuth2Client;
export const archiveOAuth2Client = clients.archiveOAuth2Client;
export const getProducts = clients.getProducts;
export const getProduct = clients.getProduct;
export const createProduct = clients.createProduct;
export const updateProduct = clients.updateProduct;
export const archiveProduct = clients.archiveProduct;
export const getSubscription = clients.getSubscription;
export const createSubscription = clients.createSubscription;
export const updateSubscription = clients.updateSubscription;
export const archiveSubscription = clients.archiveSubscription;
export const getSubscriptionsForAccount = clients.getSubscriptionsForAccount;
export const getServiceSettings = clients.getServiceSettings;
export const searchForServiceSettings = clients.searchForServiceSettings;
export const getServiceSetting = clients.getServiceSetting;
export const createServiceSetting = clients.createServiceSetting;
export const archiveServiceSetting = clients.archiveServiceSetting;
export const getWaitlists = clients.getWaitlists;
export const getWaitlist = clients.getWaitlist;
export const getWaitlistSignupsForWaitlist = clients.getWaitlistSignupsForWaitlist;
export const getIssueReports = clients.getIssueReports;
export const getIssueReport = clients.getIssueReport;
export const updateIssueReport = clients.updateIssueReport;
export const archiveIssueReport = clients.archiveIssueReport;
export const testQueueMessage = clients.testQueueMessage;
export const getAuditLogEntriesForUser = clients.getAuditLogEntriesForUser;
export const getAuditLogEntriesForAccount = clients.getAuditLogEntriesForAccount;
export const trackEvent = clients.trackEvent;
export const trackAnonymousEvent = clients.trackAnonymousEvent;
export const adminListSessionsForUser = clients.adminListSessionsForUser;
export const adminRevokeUserSession = clients.adminRevokeUserSession;
export const adminRevokeAllUserSessions = clients.adminRevokeAllUserSessions;
export const revokeCurrentSession = clients.revokeCurrentSession;
