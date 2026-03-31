//
//  SubscriptionService.swift
//  ios
//
//  Subscription management: customer info, offerings, purchases.
//  Product identifiers (monthly, yearly) must match App Store Connect and RevenueCat dashboard.
//

import Foundation
import RevenueCat

enum SubscriptionService {
  /// Product identifiers as configured in App Store Connect and RevenueCat.
  enum ProductID {
    static let monthly = "monthly"
    static let yearly = "yearly"
  }

  /// Fetches current customer info. Returns nil when RevenueCat is not configured or on error.
  static func customerInfo() async -> CustomerInfo? {
    guard RevenueCatConfiguration.isConfigured else { return nil }
    do {
      return try await Purchases.shared.customerInfo()
    } catch {
      print("⚠️ SubscriptionService: Failed to fetch customer info: \(error)")
      return nil
    }
  }

  /// Offering identifier for the launch paywall. Must match RevenueCat dashboard.
  static let launchOfferingID = "launch"

  /// Fetches available offerings. Returns nil when RevenueCat is not configured or on error.
  static func offerings() async -> Offerings? {
    guard RevenueCatConfiguration.isConfigured else { return nil }
    do {
      return try await Purchases.shared.offerings()
    } catch {
      print("⚠️ SubscriptionService: Failed to fetch offerings: \(error)")
      return nil
    }
  }

  /// Fetches the "launch" offering. Returns nil when not configured or on error.
  static func launchOffering() async -> Offering? {
    guard let offerings = await offerings() else { return nil }
    return offerings[launchOfferingID]
  }

  /// Purchases a package. Returns (customerInfo, true) on success, (nil, false) on failure.
  static func purchase(package: Package) async -> (CustomerInfo?, Bool) {
    guard RevenueCatConfiguration.isConfigured else {
      return (nil, false)
    }
    do {
      let result = try await Purchases.shared.purchase(package: package)
      if !result.userCancelled {
        return (result.customerInfo, true)
      }
      return (nil, false)
    } catch {
      print("⚠️ SubscriptionService: Purchase failed: \(error)")
      return (nil, false)
    }
  }

  /// Restores previous purchases. Returns customerInfo on success, nil on failure.
  static func restorePurchases() async -> CustomerInfo? {
    guard RevenueCatConfiguration.isConfigured else { return nil }
    do {
      return try await Purchases.shared.restorePurchases()
    } catch {
      print("⚠️ SubscriptionService: Restore failed: \(error)")
      return nil
    }
  }
}
