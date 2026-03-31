//
//  MediaBucket.swift
//  ios
//

import Foundation

/// Bucket names for media uploads. Use the appropriate bucket for each use case:
/// - avatars: User profile photos
/// - recipes: Recipe images
/// - meals: Meal images
/// - custom: For any other bucket (pass raw string)
public enum MediaBucket: Sendable {
  case avatars
  case recipes
  case meals
  case custom(String)

  public var rawValue: String {
    switch self {
    case .avatars: return "avatars"
    case .recipes: return "recipes"
    case .meals: return "meals"
    case .custom(let name): return name
    }
  }
}
