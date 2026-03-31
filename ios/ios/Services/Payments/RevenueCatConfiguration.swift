//
//  RevenueCatConfiguration.swift
//  ios
//
//  RevenueCat API key from Info.plist. Injected via Secrets.xcconfig at build time.
//

import Foundation

enum RevenueCatConfiguration {
  private static let placeholderAPIKey = "your-revenuecat-api-key-here"

  /// RevenueCat API key from Info.plist. Injected via Secrets.xcconfig at build time.
  static var revenueCatAPIKey: String {
    Bundle.main.infoDictionary?["RevenueCatAPIKey"] as? String ?? ""
  }

  /// Whether RevenueCat is properly configured (non-empty, non-placeholder key).
  static var isConfigured: Bool {
    let key = revenueCatAPIKey.trimmingCharacters(in: .whitespaces)
    return !key.isEmpty && key != placeholderAPIKey
  }
}
