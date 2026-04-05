//
//  MediaBucket.swift
//  ios
//

import Foundation

/// Bucket names for media uploads. Use the appropriate bucket for each use case:
/// - avatars: User profile photos
/// - custom: For any other bucket (pass raw string)
public enum MediaBucket: Sendable {
  case avatars
  case custom(String)

  public var rawValue: String {
    switch self {
    case .avatars: return "avatars"
    case .custom(let name): return name
    }
  }
}
