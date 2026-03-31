//
//  EventReporter.swift
//  ios
//
//  Event capture interface. All analytics must go through this protocol.
//  No external analytics SDK imports in this file.
//

import Foundation

/// Protocol defining the event reporting interface.
/// Allows flexible implementations (Segment, Noop, etc.) with rigid exclusive use across the app.
protocol EventReporter: AnyObject {
  /// Associate the current identity with a user (e.g. after login).
  /// Maps to AddUser in the Go backend.
  func identify(userID: String, properties: [String: Any])

  /// Record an event. Uses current identity from identify.
  /// Maps to EventOccurred in the Go backend.
  func track(event: String, properties: [String: Any])

  /// Clear identity (e.g. on logout).
  func reset()
}
