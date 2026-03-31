//
//  UserProfileViewModel.swift
//  ios
//

import Foundation
import GRPCCore
import GRPCNIOTransportHTTP2
import SwiftProtobuf
import SwiftUI

@Observable
@MainActor
class UserProfileViewModel {
  var user: Identity_User?
  var isLoading = false
  var errorMessage: String?
  var didSucceed = false

  var username: String = ""
  var firstName: String = ""
  var lastName: String = ""
  var birthday = Date()

  var hasTwoFactor: Bool {
    user?.hasTwoFactorSecretVerifiedAt == true
  }

  var usernameHasChanged: Bool {
    guard let user = user else { return false }
    return user.username != username
  }

  var detailsHasChanged: Bool {
    guard let user = user else { return false }
    let currentBirthday = user.hasBirthday ? dateFromTimestamp(user.birthday) : nil
    return user.firstName != firstName
      || user.lastName != lastName
      || (currentBirthday?.timeIntervalSince1970 ?? 0) != birthday.timeIntervalSince1970
  }

  private let authManager: AuthenticationManager

  init(authManager: AuthenticationManager) {
    self.authManager = authManager
  }

  func loadUser() async {
    isLoading = true
    errorMessage = nil

    do {
      let response = try await fetchUser()
      if response.hasResult {
        self.user = response.result
        initializeFormFields(from: response.result)
      }
    } catch {
      await authManager.invalidateCredentialsIfSessionError(error)
      errorMessage = ErrorDisplayFormatter.format(error, context: "load user").message
    }

    isLoading = false
  }

  func updateUsername() async -> Bool {
    guard usernameHasChanged, !username.isEmpty else { return false }

    return await performUpdate {
      let (clientManager, metadata) = try await getClientManagerAndMetadata()
      var request = Identity_UpdateUserUsernameRequest()
      request.newUsername = username

      _ = try await clientManager.client.identity.updateUserUsername(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )
      await loadUser()
    }
  }

  func updateUserDetails(currentPassword: String, totpToken: String = "") async -> Bool {
    guard detailsHasChanged else { return false }
    guard !currentPassword.isEmpty else {
      errorMessage = "Password is required"
      return false
    }

    return await performUpdate {
      let (clientManager, metadata) = try await getClientManagerAndMetadata()
      var input = Identity_UserDetailsUpdateRequestInput()
      input.firstName = firstName
      input.lastName = lastName
      input.birthday = timestampFromDate(birthday)
      input.currentPassword = currentPassword
      input.totpToken = totpToken

      var request = Identity_UpdateUserDetailsRequest()
      request.input = input

      _ = try await clientManager.client.identity.updateUserDetails(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )
      await loadUser()
    }
  }

  private func fetchUser() async throws -> Auth_GetSelfResponse {
    let (clientManager, metadata) = try await getClientManagerAndMetadata()
    return try await clientManager.client.auth.getSelf(
      Auth_GetSelfRequest(),
      metadata: metadata,
      options: clientManager.defaultCallOptions
    )
  }

  private func initializeFormFields(from user: Identity_User) {
    username = user.username
    firstName = user.firstName
    lastName = user.lastName
    birthday = user.hasBirthday ? dateFromTimestamp(user.birthday) : Date()
  }

  private func dateFromTimestamp(_ timestamp: SwiftProtobuf.Google_Protobuf_Timestamp) -> Date {
    let seconds = TimeInterval(timestamp.seconds)
    let nanos = TimeInterval(timestamp.nanos) / 1_000_000_000.0
    return Date(timeIntervalSince1970: seconds + nanos)
  }

  private func timestampFromDate(_ date: Date) -> SwiftProtobuf.Google_Protobuf_Timestamp {
    var timestamp = SwiftProtobuf.Google_Protobuf_Timestamp()
    timestamp.seconds = Int64(date.timeIntervalSince1970)
    timestamp.nanos = 0
    return timestamp
  }

  private func performUpdate(operation: () async throws -> Void) async -> Bool {
    isLoading = true
    errorMessage = nil

    do {
      try await operation()
      didSucceed = true
      isLoading = false
      return true
    } catch {
      await authManager.invalidateCredentialsIfSessionError(error)
      errorMessage = ErrorDisplayFormatter.format(error, context: "update").message
      isLoading = false
      return false
    }
  }

  private func getClientManagerAndMetadata() async throws -> (
    ClientManager<HTTP2ClientTransport.TransportServices>, GRPCCore.Metadata
  ) {
    guard let clientManager = try? authManager.getClientManager() else {
      throw NSError(
        domain: "UserProfileViewModel", code: 1,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get client manager"])
    }

    guard let oauth2Token = await authManager.getOAuth2AccessToken() else {
      throw NSError(
        domain: "UserProfileViewModel", code: 2,
        userInfo: [NSLocalizedDescriptionKey: "Failed to get OAuth2 access token"])
    }

    let metadata = clientManager.authenticatedMetadata(accessToken: oauth2Token)
    return (clientManager, metadata)
  }
}
