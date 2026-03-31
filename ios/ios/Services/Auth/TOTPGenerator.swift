//
//  TOTPGenerator.swift
//  ios
//
//  Created by Auto on 12/8/25.
//

import CryptoKit
import Foundation

/// Simple TOTP generator for development purposes
/// Generates TOTP codes using RFC 6238 (Time-based One-Time Password)
struct TOTPGenerator {
  /// Generate a TOTP code from a base32-encoded secret
  /// - Parameters:
  ///   - secret: Base32-encoded TOTP secret
  ///   - timeStep: Time step in seconds (default: 30)
  ///   - digits: Number of digits in the code (default: 6)
  /// - Returns: TOTP code as a string, or nil if generation fails
  static func generate(secret: String, timeStep: Int = 30, digits: Int = 6) -> String? {
    // Decode base32 secret
    guard let secretData = base32Decode(secret) else {
      print("❌ Failed to decode base32 secret")
      return nil
    }

    // Calculate time counter
    let time = Int64(Date().timeIntervalSince1970)
    let counter = time / Int64(timeStep)

    // Generate HMAC-SHA1
    let key = SymmetricKey(data: secretData)
    // Convert counter to 8-byte big-endian data
    var counterBytes = counter.bigEndian
    let counterData = withUnsafeBytes(of: &counterBytes) { Data($0) }

    let hmac = HMAC<Insecure.SHA1>.authenticationCode(for: counterData, using: key)
    let hmacData = Data(hmac)

    // Dynamic truncation
    let offset = Int(hmacData[hmacData.count - 1] & 0x0F)
    let binary =
      ((Int(hmacData[offset]) & 0x7F) << 24 | (Int(hmacData[offset + 1]) & 0xFF) << 16
        | (Int(hmacData[offset + 2]) & 0xFF) << 8 | (Int(hmacData[offset + 3]) & 0xFF))

    let otp = binary % Int(pow(10, Double(digits)))
    return String(format: "%0\(digits)d", otp)
  }

  /// Decode base32 string to Data
  private static func base32Decode(_ string: String) -> Data? {
    let base32Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
    let uppercaseString = string.uppercased().replacingOccurrences(of: " ", with: "")

    var bits = 0
    var value = 0
    var data = Data()

    for char in uppercaseString {
      guard let index = base32Alphabet.firstIndex(of: char) else {
        continue
      }

      value = (value << 5) | base32Alphabet.distance(from: base32Alphabet.startIndex, to: index)
      bits += 5

      if bits >= 8 {
        data.append(UInt8((value >> (bits - 8)) & 0xFF))
        bits -= 8
      }
    }

    return data.isEmpty ? nil : data
  }
}
