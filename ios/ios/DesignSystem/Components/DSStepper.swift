//
//  DSStepper.swift
//  ios
//
//  Design System Stepper - Increment/decrement buttons for scale and numeric inputs
//

import SwiftUI

// MARK: - DSStepperButtons

/// A reusable increment/decrement control with optional center content.
/// Use for scale inputs, meal assignment, and other numeric steppers.
///
/// Usage:
/// ```swift
/// // Just the buttons
/// DSStepperButtons(onDecrement: { scale -= 0.25 }, onIncrement: { scale += 0.25 })
///
/// // With center value display
/// DSStepperButtons(onDecrement: { ... }, onIncrement: { ... }) {
///   Text("1.0")
///     .font(.subheadline)
///     .frame(minWidth: 36, alignment: .center)
/// }
/// ```
struct DSStepperButtons<Center: View>: View {
  let onDecrement: () -> Void
  let onIncrement: () -> Void
  @ViewBuilder let center: () -> Center

  init(
    onDecrement: @escaping () -> Void,
    onIncrement: @escaping () -> Void,
    @ViewBuilder center: @escaping () -> Center
  ) {
    self.onDecrement = onDecrement
    self.onIncrement = onIncrement
    self.center = center
  }

  var body: some View {
    HStack(spacing: DSTheme.Spacing.sm) {
      Button(action: onDecrement) {
        Image(systemName: "minus.circle")
          .font(.system(size: DSTheme.IconSize.lg))
          .foregroundColor(DSTheme.Colors.textPrimary)
      }
      .buttonStyle(.plain)

      center()

      Button(action: onIncrement) {
        Image(systemName: "plus.circle")
          .font(.system(size: DSTheme.IconSize.lg))
          .foregroundColor(DSTheme.Colors.textPrimary)
      }
      .buttonStyle(.plain)
    }
    .padding(.horizontal, DSTheme.Spacing.sm)
    .padding(.vertical, DSTheme.Spacing.xs)
    .background(DSTheme.Colors.cardBackgroundElevated)
    .cornerRadius(DSTheme.Radius.sm)
  }
}

// MARK: - Convenience initializer without center content

extension DSStepperButtons where Center == EmptyView {
  init(onDecrement: @escaping () -> Void, onIncrement: @escaping () -> Void) {
    self.onDecrement = onDecrement
    self.onIncrement = onIncrement
    self.center = { EmptyView() }
  }
}

// MARK: - Preview

#Preview("DSStepperButtons") {
  VStack(spacing: DSTheme.Spacing.lg) {
    DSStepperButtons(onDecrement: {}, onIncrement: {})

    DSStepperButtons(onDecrement: {}, onIncrement: {}) {
      Text("1.0")
        .font(DSTheme.Typography.body)
        .frame(minWidth: 36, alignment: .center)
    }
  }
  .padding()
}
