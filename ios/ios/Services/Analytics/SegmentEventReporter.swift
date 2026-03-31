//
//  SegmentEventReporter.swift
//  ios
//
//  Segment-backed EventReporter using Analytics-Swift.
//  This is the ONLY file that should import Segment.
//

import Foundation
import Segment

/// Segment-backed EventReporter. Reports events to Segment via Analytics-Swift.
final class SegmentEventReporter: EventReporter {
  private let analytics: Analytics

  init(writeKey: String) {
    let configuration = Configuration(writeKey: writeKey)
      .setTrackedApplicationLifecycleEvents(.all)
      .flushInterval(10)
    analytics = Analytics(configuration: configuration)
  }

  func identify(userID: String, properties: [String: Any]) {
    let traits = Self.encodeToTraits(properties)
    analytics.identify(userId: userID, traits: traits)
  }

  func track(event: String, properties: [String: Any]) {
    let props = Self.encodeToProperties(properties)
    analytics.track(name: event, properties: props)
  }

  func reset() {
    analytics.reset()
  }

  /// Encodes [String: Any] to a Codable type for Segment's identify traits.
  private static func encodeToTraits(_ dict: [String: Any]) -> [String: String] {
    dict.compactMapValues { value in
      if let s = value as? String { return s }
      if let n = value as? Int { return String(n) }
      if let b = value as? Bool { return b ? "true" : "false" }
      return nil
    }
  }

  /// Encodes [String: Any] to a Codable type for Segment's track properties.
  private static func encodeToProperties(_ dict: [String: Any]) -> [String: String]? {
    let converted = encodeToTraits(dict)
    return converted.isEmpty ? nil : converted
  }
}
