//
//  ContentView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import RevenueCat
import RevenueCatUI
import SwiftUI

struct ContentView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(DeepLinkHandler.self) private var deepLinkHandler
  @Environment(EventReporterService.self) private var eventReporterService
  @State private var showLogin: Bool = true

  // Invitation data from deep link
  @State private var pendingInvitationID: String = ""
  @State private var pendingInvitationToken: String = ""

  // Sheet for logged-in users who tap an invite link
  @State private var showAcceptInvitationSheet: Bool = false

  // Launch offering for paywall (avoids "default" offering which has no paywall)
  @State private var launchOffering: Offering?

  var body: some View {
    Group {
      if authManager.isAuthenticated && !authManager.oauth2AccessToken.isEmpty {
        if RevenueCatConfiguration.isConfigured, let offering = launchOffering {
          HomeView()
            .presentPaywallIfNeeded(
              requiredEntitlementIdentifier: EntitlementService.proEntitlementID,
              offering: offering
            )
        } else if RevenueCatConfiguration.isConfigured {
          HomeView()
            .task { launchOffering = await SubscriptionService.launchOffering() }
        } else {
          HomeView()
        }
      } else if showLogin {
        LoginView(
          showRegister: Binding(
            get: { false },
            set: { _ in showLogin = false }
          )
        )
      } else {
        RegisterView(
          showLogin: Binding(
            get: { true },
            set: { _ in showLogin = true }
          ),
          invitationID: pendingInvitationID,
          invitationToken: pendingInvitationToken
        )
      }
    }
    .onChange(of: deepLinkHandler.pendingDestination) { _, newDestination in
      handleDeepLink(newDestination)
    }
    .onAppear {
      // Handle any pending deep link on appear
      handleDeepLink(deepLinkHandler.pendingDestination)
    }
    .sheet(
      isPresented: $showAcceptInvitationSheet,
      onDismiss: { clearPendingInvitation() },
      content: {
        AcceptInvitationSheet(
          invitationID: pendingInvitationID,
          invitationToken: pendingInvitationToken,
          onAccepted: {}
        )
        .environment(authManager)
        .environment(eventReporterService)
      }
    )
  }

  private func handleDeepLink(_ destination: DeepLinkDestination?) {
    guard let destination = destination else { return }

    switch destination {
    case .acceptInvitation(let invitationID, let token):
      eventReporterService.reporter.track(
        event: "invitation_deep_link_opened",
        properties: ["invitation_id": invitationID])
      pendingInvitationID = invitationID
      pendingInvitationToken = token
      deepLinkHandler.clearPendingDestination()

      if authManager.isAuthenticated && !authManager.oauth2AccessToken.isEmpty {
        // Logged in: show accept-invitation sheet
        showAcceptInvitationSheet = true
      } else {
        // Not logged in: navigate to registration
        showLogin = false
      }

    case .resetPassword(let token):
      // swiftlint:disable:next todo
      // TODO: Navigate to password reset flow
      print("Password reset token: \(token)")
      deepLinkHandler.clearPendingDestination()

    case .verifyEmail(let token):
      // swiftlint:disable:next todo
      // TODO: Handle email verification
      print("Email verification token: \(token)")
      deepLinkHandler.clearPendingDestination()

    case .unknown:
      break
    }
  }

  private func clearPendingInvitation() {
    pendingInvitationID = ""
    pendingInvitationToken = ""
  }
}

#Preview {
  ContentView()
    .environment(EventReporterService())
    .environment(AuthenticationManager())
    .environment(DeepLinkHandler())
}
