//
//  KeychainManager.swift
//  ios
//

import Foundation
import Security

/// Thin wrapper around the iOS Keychain Services API.
/// All items are stored as generic passwords scoped to the app's bundle identifier.
enum KeychainManager {

  // MARK: - Keys

  enum Key: String {
    case accessToken
    case refreshToken
    case oauth2AccessToken
    case oauth2RefreshToken
    case oauth2TokenExpiresAt
    case username
    case userID
    case accountID

    var keychainKey: String {
      "\(Branding.keychainPrefix).\(rawValue)"
    }
  }

  // MARK: - Public API

  static func save(_ value: String, for key: Key) {
    guard let data = value.data(using: .utf8) else { return }
    save(data: data, for: key)
  }

  static func loadString(for key: Key) -> String? {
    guard let data = loadData(for: key) else { return nil }
    return String(data: data, encoding: .utf8)
  }

  static func delete(key: Key) {
    let query: [String: Any] = [
      kSecClass as String: kSecClassGenericPassword,
      kSecAttrAccount as String: key.keychainKey,
    ]
    SecItemDelete(query as CFDictionary)
  }

  static func deleteAll() {
    for key in [
      Key.accessToken, .refreshToken,
      .oauth2AccessToken, .oauth2RefreshToken, .oauth2TokenExpiresAt,
      .username, .userID, .accountID,
    ] {
      delete(key: key)
    }
  }

  // MARK: - Helpers

  private static func save(data: Data, for key: Key) {
    // Delete any existing item first so the add doesn't fail with errSecDuplicateItem
    delete(key: key)

    let query: [String: Any] = [
      kSecClass as String: kSecClassGenericPassword,
      kSecAttrAccount as String: key.keychainKey,
      kSecValueData as String: data,
      kSecAttrAccessible as String: kSecAttrAccessibleAfterFirstUnlock,
    ]
    SecItemAdd(query as CFDictionary, nil)
  }

  private static func loadData(for key: Key) -> Data? {
    let query: [String: Any] = [
      kSecClass as String: kSecClassGenericPassword,
      kSecAttrAccount as String: key.keychainKey,
      kSecReturnData as String: true,
      kSecMatchLimit as String: kSecMatchLimitOne,
    ]

    var result: AnyObject?
    let status = SecItemCopyMatching(query as CFDictionary, &result)
    guard status == errSecSuccess else { return nil }
    return result as? Data
  }
}
