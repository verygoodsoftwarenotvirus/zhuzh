//
//  Branding.swift
//  ios
//

import Foundation

/// Core branding constants - change these to rebrand the project.
enum Branding {
  static let companyName = "Zhuzh"
  static let companyNameSlug = "zhuzh"
  static let publicDomain = "\(companyNameSlug).com"
  static let proEntitlementName = "\(companyName) Pro"

  // Keychain key prefix
  static let keychainPrefix = "com.\(companyNameSlug)"

  // URLs
  static let productionAPIDomain = "http-api.\(publicDomain)"
  static let productionWebDomain = "www.\(publicDomain)"
  static let productionGRPCDomain = "api.\(publicDomain)"
  static let productionMediaDomain = "media.\(publicDomain)"

  static let productionAPIURL = "https://\(productionAPIDomain)"
  static let productionWebURL = "https://\(productionWebDomain)"
  static let productionMediaURL = "https://\(productionMediaDomain)"
}
