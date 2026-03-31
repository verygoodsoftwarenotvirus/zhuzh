//
//  StepTimerNotificationService.swift
//  ios
//
//  Schedules local notifications when recipe step timers reach their minimum elapsed time,
//  so the user is notified even when the app is backgrounded or the device is locked.
//

import Foundation
import UserNotifications

/// Schedules and cancels local notifications for recipe step timers.
/// Notifications fire when the timer reaches its minimum elapsed time (user can proceed).
@available(macOS 15.0, iOS 18.0, watchOS 11.0, tvOS 18.0, visionOS 2.0, *)
enum StepTimerNotificationService {
  private static let identifierPrefix = "step_timer:"

  /// Schedules a local notification to fire when the step timer reaches its minimum elapsed time.
  /// - Parameters:
  ///   - stepKey: Unique key for the step (e.g. "recipeID:stepID"), used as notification identifier
  ///   - recipeName: Recipe name for the notification body
  ///   - stepName: Step name for the notification body
  ///   - minSeconds: Seconds until the notification should fire (minimum timer duration)
  static func scheduleNotification(
    stepKey: String,
    recipeName: String,
    stepName: String,
    minSeconds: UInt32
  ) {
    let identifier = identifierPrefix + stepKey
    let content = UNMutableNotificationContent()
    content.title = "Timer ready"
    content.body = "\(stepName) — \(recipeName)"
    content.sound = .default
    content.categoryIdentifier = "STEP_TIMER"

    let trigger = UNTimeIntervalNotificationTrigger(
      timeInterval: TimeInterval(minSeconds),
      repeats: false
    )
    let request = UNNotificationRequest(identifier: identifier, content: content, trigger: trigger)

    UNUserNotificationCenter.current().add(request) { error in
      if let error {
        print("⚠️ Failed to schedule step timer notification: \(error.localizedDescription)")
      }
    }
  }

  /// Cancels any pending notification for the given step timer.
  static func cancelNotification(stepKey: String) {
    let identifier = identifierPrefix + stepKey
    UNUserNotificationCenter.current().removePendingNotificationRequests(withIdentifiers: [
      identifier
    ])
  }
}
