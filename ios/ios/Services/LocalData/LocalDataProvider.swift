//
//  LocalDataProvider.swift
//  ios
//
//  Provides offline data from a bundled JSON seed file.
//

import Foundation

// MARK: - LocalDataProvider

@MainActor
class LocalDataProvider {
  static let shared = LocalDataProvider()

  private var isLoaded = false

  func loadIfNeeded() {
    guard !isLoaded else { return }
    isLoaded = true
  }
}
