//
//  DeepLinkHandlerTests.swift
//  iosTests
//
//  Unit tests for DeepLinkHandler URL parsing.
//

// swiftlint:disable force_unwrapping

import Foundation
@testable import ios
import Testing

struct DeepLinkHandlerTests {
  // MARK: - Initialization Tests

  @Test("DeepLinkHandler initializes with no pending destination")
  func testInitialization() {
    let handler = DeepLinkHandler()

    #expect(handler.pendingDestination == nil)
  }

  // MARK: - Accept Invitation URL Tests

  @Test("parseURL correctly parses accept_invitation URL with valid parameters")
  func testParseAcceptInvitationURL() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation?i=inv123&t=token456")!

    let result = handler.parseURL(url)

    #expect(result == .acceptInvitation(invitationID: "inv123", token: "token456"))
  }

  @Test("parseURL handles accept_invitation with URL-encoded parameters")
  func testParseAcceptInvitationURLEncoded() {
    let handler = DeepLinkHandler()
    // URL with special characters that would be encoded
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation?i=inv%2B123&t=token%3D456")!

    let result = handler.parseURL(url)

    #expect(result == .acceptInvitation(invitationID: "inv+123", token: "token=456"))
  }

  @Test("parseURL returns unknown for accept_invitation missing invitation ID")
  func testParseAcceptInvitationMissingID() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation?t=token456")!

    let result = handler.parseURL(url)

    #expect(result == .unknown)
  }

  @Test("parseURL returns unknown for accept_invitation missing token")
  func testParseAcceptInvitationMissingToken() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation?i=inv123")!

    let result = handler.parseURL(url)

    #expect(result == .unknown)
  }

  @Test("parseURL returns unknown for accept_invitation with no parameters")
  func testParseAcceptInvitationNoParams() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation")!

    let result = handler.parseURL(url)

    #expect(result == .unknown)
  }

  @Test("parseURL handles accept_invitation with empty parameter values")
  func testParseAcceptInvitationEmptyValues() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation?i=&t=")!

    let result = handler.parseURL(url)

    // Empty strings are still valid values, just not useful
    #expect(result == .acceptInvitation(invitationID: "", token: ""))
  }

  @Test("parseURL handles accept_invitation with extra parameters")
  func testParseAcceptInvitationExtraParams() {
    let handler = DeepLinkHandler()
    let url = URL(
      string: "https://www.zhuzh.dev/accept_invitation?i=inv123&t=token456&extra=ignored")!

    let result = handler.parseURL(url)

    #expect(result == .acceptInvitation(invitationID: "inv123", token: "token456"))
  }

  // MARK: - Reset Password URL Tests

  @Test("parseURL correctly parses reset_password URL")
  func testParseResetPasswordURL() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/reset_password?t=resettoken789")!

    let result = handler.parseURL(url)

    #expect(result == .resetPassword(token: "resettoken789"))
  }

  @Test("parseURL returns unknown for reset_password missing token")
  func testParseResetPasswordMissingToken() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/reset_password")!

    let result = handler.parseURL(url)

    #expect(result == .unknown)
  }

  @Test("parseURL returns unknown for reset_password with wrong parameter name")
  func testParseResetPasswordWrongParam() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/reset_password?token=resettoken789")!

    let result = handler.parseURL(url)

    #expect(result == .unknown)
  }

  // MARK: - Verify Email URL Tests

  @Test("parseURL correctly parses verify_email_address URL")
  func testParseVerifyEmailURL() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/verify_email_address?t=verifytoken123")!

    let result = handler.parseURL(url)

    #expect(result == .verifyEmail(token: "verifytoken123"))
  }

  @Test("parseURL returns unknown for verify_email_address missing token")
  func testParseVerifyEmailMissingToken() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/verify_email_address")!

    let result = handler.parseURL(url)

    #expect(result == .unknown)
  }

  // MARK: - Unknown Path Tests

  @Test("parseURL returns unknown for unrecognized path")
  func testParseUnknownPath() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/unknown_path?t=token")!

    let result = handler.parseURL(url)

    #expect(result == .unknown)
  }

  @Test("parseURL returns unknown for root path")
  func testParseRootPath() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/")!

    let result = handler.parseURL(url)

    #expect(result == .unknown)
  }

  @Test("parseURL returns unknown for empty path")
  func testParseEmptyPath() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev")!

    let result = handler.parseURL(url)

    #expect(result == .unknown)
  }

  @Test("parseURL handles path with trailing slash")
  func testParsePathWithTrailingSlash() {
    let handler = DeepLinkHandler()
    // Note: URL normalization may affect this
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation/?i=inv123&t=token456")!

    let result = handler.parseURL(url)

    // Trailing slash results in a different path, so it won't match
    #expect(result == .unknown)
  }

  // MARK: - Different Domains Tests

  @Test("parseURL works with development domain")
  func testParseDevelopmentDomain() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation?i=inv&t=tok")!

    let result = handler.parseURL(url)

    #expect(result == .acceptInvitation(invitationID: "inv", token: "tok"))
  }

  @Test("parseURL works with production domain")
  func testParseProductionDomain() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation?i=inv&t=tok")!

    let result = handler.parseURL(url)

    #expect(result == .acceptInvitation(invitationID: "inv", token: "tok"))
  }

  @Test("parseURL works with localhost")
  func testParseLocalhost() {
    let handler = DeepLinkHandler()
    let url = URL(string: "http://localhost:3000/accept_invitation?i=inv&t=tok")!

    let result = handler.parseURL(url)

    #expect(result == .acceptInvitation(invitationID: "inv", token: "tok"))
  }

  // MARK: - handleURL Tests

  @Test("handleURL sets pending destination for valid URL")
  func testHandleURLSetsPendingDestination() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation?i=inv123&t=token456")!

    handler.handleURL(url)

    #expect(handler.pendingDestination == .acceptInvitation(invitationID: "inv123", token: "token456"))
  }

  @Test("handleURL does not set pending destination for unknown URL")
  func testHandleURLIgnoresUnknownURL() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/unknown")!

    handler.handleURL(url)

    #expect(handler.pendingDestination == nil)
  }

  @Test("handleURL overwrites previous pending destination")
  func testHandleURLOverwritesPending() {
    let handler = DeepLinkHandler()
    let url1 = URL(string: "https://www.zhuzh.dev/accept_invitation?i=inv1&t=tok1")!
    let url2 = URL(string: "https://www.zhuzh.dev/accept_invitation?i=inv2&t=tok2")!

    handler.handleURL(url1)
    handler.handleURL(url2)

    #expect(handler.pendingDestination == .acceptInvitation(invitationID: "inv2", token: "tok2"))
  }

  // MARK: - clearPendingDestination Tests

  @Test("clearPendingDestination clears the pending destination")
  func testClearPendingDestination() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation?i=inv&t=tok")!

    handler.handleURL(url)
    #expect(handler.pendingDestination != nil)

    handler.clearPendingDestination()
    #expect(handler.pendingDestination == nil)
  }

  @Test("clearPendingDestination is safe to call when already nil")
  func testClearPendingDestinationWhenNil() {
    let handler = DeepLinkHandler()

    #expect(handler.pendingDestination == nil)

    // Should not crash
    handler.clearPendingDestination()

    #expect(handler.pendingDestination == nil)
  }

  // MARK: - DeepLinkDestination Equatable Tests

  @Test("DeepLinkDestination equality works for acceptInvitation")
  func testDestinationEqualityAcceptInvitation() {
    let dest1 = DeepLinkDestination.acceptInvitation(invitationID: "id1", token: "tok1")
    let dest2 = DeepLinkDestination.acceptInvitation(invitationID: "id1", token: "tok1")
    let dest3 = DeepLinkDestination.acceptInvitation(invitationID: "id2", token: "tok1")

    #expect(dest1 == dest2)
    #expect(dest1 != dest3)
  }

  @Test("DeepLinkDestination equality works for resetPassword")
  func testDestinationEqualityResetPassword() {
    let dest1 = DeepLinkDestination.resetPassword(token: "tok1")
    let dest2 = DeepLinkDestination.resetPassword(token: "tok1")
    let dest3 = DeepLinkDestination.resetPassword(token: "tok2")

    #expect(dest1 == dest2)
    #expect(dest1 != dest3)
  }

  @Test("DeepLinkDestination equality works for verifyEmail")
  func testDestinationEqualityVerifyEmail() {
    let dest1 = DeepLinkDestination.verifyEmail(token: "tok1")
    let dest2 = DeepLinkDestination.verifyEmail(token: "tok1")
    let dest3 = DeepLinkDestination.verifyEmail(token: "tok2")

    #expect(dest1 == dest2)
    #expect(dest1 != dest3)
  }

  @Test("DeepLinkDestination equality works for unknown")
  func testDestinationEqualityUnknown() {
    let dest1 = DeepLinkDestination.unknown
    let dest2 = DeepLinkDestination.unknown

    #expect(dest1 == dest2)
  }

  @Test("DeepLinkDestination different types are not equal")
  func testDestinationDifferentTypesNotEqual() {
    let invitation = DeepLinkDestination.acceptInvitation(invitationID: "id", token: "tok")
    let resetPassword = DeepLinkDestination.resetPassword(token: "tok")
    let verifyEmail = DeepLinkDestination.verifyEmail(token: "tok")
    let unknown = DeepLinkDestination.unknown

    #expect(invitation != resetPassword)
    #expect(invitation != verifyEmail)
    #expect(invitation != unknown)
    #expect(resetPassword != verifyEmail)
    #expect(resetPassword != unknown)
    #expect(verifyEmail != unknown)
  }

  // MARK: - Edge Cases

  @Test("parseURL handles URL with fragment")
  func testParseURLWithFragment() {
    let handler = DeepLinkHandler()
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation?i=inv&t=tok#section")!

    let result = handler.parseURL(url)

    // Fragment should be ignored, query params still parsed
    #expect(result == .acceptInvitation(invitationID: "inv", token: "tok"))
  }

  @Test("parseURL handles very long parameter values")
  func testParseURLLongValues() {
    let handler = DeepLinkHandler()
    let longID = String(repeating: "a", count: 1000)
    let longToken = String(repeating: "b", count: 1000)
    let url = URL(string: "https://www.zhuzh.dev/accept_invitation?i=\(longID)&t=\(longToken)")!

    let result = handler.parseURL(url)

    #expect(result == .acceptInvitation(invitationID: longID, token: longToken))
  }

  @Test("parseURL handles Unicode in parameter values")
  func testParseURLUnicodeValues() {
    let handler = DeepLinkHandler()
    // URL-encoded Unicode characters
    let url = URL(
      string:
        "https://www.zhuzh.dev/accept_invitation?i=%E4%B8%AD%E6%96%87&t=%F0%9F%8E%89")!

    let result = handler.parseURL(url)

    if case .acceptInvitation(let id, let token) = result {
      #expect(id == "中文")
      #expect(token == "🎉")
    } else {
      Issue.record("Expected acceptInvitation destination")
    }
  }

  #if DEBUG
    // MARK: - Debug Helper Tests

    @Test("simulateInvitationLink sets pending destination correctly")
    func testSimulateInvitationLink() {
      let handler = DeepLinkHandler()

      handler.simulateInvitationLink(invitationID: "test-id", token: "test-token")

      #expect(handler.pendingDestination == .acceptInvitation(invitationID: "test-id", token: "test-token"))
    }
  #endif
}

// swiftlint:enable force_unwrapping
