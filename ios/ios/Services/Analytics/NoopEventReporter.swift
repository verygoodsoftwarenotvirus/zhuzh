//
//  NoopEventReporter.swift
//  ios
//
//  No-op implementation for tests and when analytics is disabled.
//

import Foundation

/// No-op EventReporter used when Segment write key is missing or for tests.
final class NoopEventReporter: EventReporter {
  func identify(userID: String, properties: [String: Any]) {}

  func track(event: String, properties: [String: Any]) {}

  func reset() {}
}
