//
//  HomeView.swift
//  ios
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import Combine
import SwiftUI

struct HomeView: View {
  @Environment(AuthenticationManager.self) private var authManager
  @Environment(EventReporterService.self) private var eventReporterService

  var body: some View {
    NavigationStack {
      VStack(spacing: 0) {
        HStack(alignment: .center, spacing: DSTheme.Spacing.md) {
          Text("\(greeting), \(authManager.username)!")
            .font(DSTheme.Typography.title1)
            .foregroundColor(DSTheme.Colors.textPrimary)
            .frame(maxWidth: .infinity, alignment: .leading)
        }
        .padding(.horizontal, DSTheme.Spacing.lg)
        .padding(.vertical, 14)

        Spacer()
      }
      .navigationTitle("")
      .navigationBarTitleDisplayMode(.inline)
      .toolbar(.hidden, for: .navigationBar)
    }
  }

  private var greeting: String {
    let hour = Calendar.current.component(.hour, from: Date())
    switch hour {
    case 0..<12:
      return "Good morning"
    case 12..<17:
      return "Good afternoon"
    default:
      return "Good evening"
    }
  }
}

#Preview {
  let authManager = AuthenticationManager()
  authManager.isAuthenticated = true
  authManager.username = "John Doe"
  authManager.userID = "user123"
  authManager.accountID = "account123"

  return HomeView()
    .environment(authManager)
}
