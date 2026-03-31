//
//  DeepLinkHandler.swift
//  ios
//
//  Created for Universal Links support.
//

import Foundation

/// Represents a parsed deep link destination
enum DeepLinkDestination: Equatable {
  case acceptInvitation(invitationID: String, token: String)
  case resetPassword(token: String)
  case verifyEmail(token: String)
  case unknown
}

/// Handles parsing and routing of deep links (Universal Links)
@Observable
class DeepLinkHandler {
  /// The current pending deep link destination, if any
  var pendingDestination: DeepLinkDestination?

  /// Parses a Universal Link URL and returns the appropriate destination
  /// - Parameter url: The incoming URL from a Universal Link
  /// - Returns: The parsed destination, or .unknown if not recognized
  func parseURL(_ url: URL) -> DeepLinkDestination {
    guard let components = URLComponents(url: url, resolvingAgainstBaseURL: true) else {
      return .unknown
    }

    let path = components.path
    let queryItems = components.queryItems ?? []

    switch path {
    case "/accept_invitation":
      // Format: /accept_invitation?i={invitationID}&t={token}
      guard let invitationID = queryItems.first(where: { $0.name == "i" })?.value,
        let token = queryItems.first(where: { $0.name == "t" })?.value
      else {
        return .unknown
      }
      return .acceptInvitation(invitationID: invitationID, token: token)

    case "/reset_password":
      // Format: /reset_password?t={token}
      guard let token = queryItems.first(where: { $0.name == "t" })?.value else {
        return .unknown
      }
      return .resetPassword(token: token)

    case "/verify_email_address":
      // Format: /verify_email_address?t={token}
      guard let token = queryItems.first(where: { $0.name == "t" })?.value else {
        return .unknown
      }
      return .verifyEmail(token: token)

    default:
      return .unknown
    }
  }

  /// Handle an incoming URL and set the pending destination
  /// - Parameter url: The incoming URL
  func handleURL(_ url: URL) {
    let destination = parseURL(url)
    if destination != .unknown {
      pendingDestination = destination
    }
  }

  /// Clear the pending destination after it's been handled
  func clearPendingDestination() {
    pendingDestination = nil
  }

  #if DEBUG
    /// For testing: manually trigger an invitation deep link
    func simulateInvitationLink(invitationID: String, token: String) {
      pendingDestination = .acceptInvitation(invitationID: invitationID, token: token)
    }
  #endif
}
