import Foundation
import GRPCCore
import GRPCNIOTransportHTTP2

@Observable
@MainActor
class SessionsViewModel {
  var sessions: [Auth_UserSession] = []
  var isLoading = false
  var errorMessage: String?

  private let authManager: AuthenticationManager

  init(authManager: AuthenticationManager) {
    self.authManager = authManager
  }

  func loadSessions() async {
    isLoading = true
    errorMessage = nil

    do {
      let (clientManager, metadata) = try getClientManagerAndJwtMetadata()
      let response = try await clientManager.client.auth.listActiveSessions(
        Auth_ListActiveSessionsRequest(),
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )
      sessions = response.sessions
    } catch {
      await authManager.invalidateCredentialsIfSessionError(error)
      errorMessage = "Failed to load sessions"
      print("❌ Error loading sessions: \(error)")
    }

    isLoading = false
  }

  func revokeSession(sessionID: String) async -> Bool {
    isLoading = true
    errorMessage = nil

    do {
      let (clientManager, metadata) = try getClientManagerAndJwtMetadata()
      var request = Auth_RevokeSessionRequest()
      request.sessionID = sessionID
      _ = try await clientManager.client.auth.revokeSession(
        request,
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )
      sessions.removeAll { $0.id == sessionID }
      isLoading = false
      return true
    } catch {
      await authManager.invalidateCredentialsIfSessionError(error)
      errorMessage = "Failed to revoke session"
      print("❌ Error revoking session: \(error)")
      isLoading = false
      return false
    }
  }

  func revokeAllOtherSessions() async -> Bool {
    isLoading = true
    errorMessage = nil

    do {
      let (clientManager, metadata) = try getClientManagerAndJwtMetadata()
      _ = try await clientManager.client.auth.revokeAllOtherSessions(
        Auth_RevokeAllOtherSessionsRequest(),
        metadata: metadata,
        options: clientManager.defaultCallOptions
      )
      sessions.removeAll { !$0.isCurrent }
      isLoading = false
      return true
    } catch {
      await authManager.invalidateCredentialsIfSessionError(error)
      errorMessage = "Failed to revoke sessions"
      print("❌ Error revoking sessions: \(error)")
      isLoading = false
      return false
    }
  }

  // Use JWT (not OAuth2) for session management so isCurrent is set correctly
  private func getClientManagerAndJwtMetadata() throws -> (
    ClientManager<HTTP2ClientTransport.TransportServices>, GRPCCore.Metadata
  ) {
    let clientManager = try authManager.getClientManager()
    let metadata = clientManager.authenticatedMetadata(accessToken: authManager.accessToken)
    return (clientManager, metadata)
  }
}
