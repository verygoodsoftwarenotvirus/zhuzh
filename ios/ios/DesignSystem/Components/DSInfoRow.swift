//
//  DSInfoRow.swift
//  ios
//
//  Design System Info Row - Card-style row with icon, text, and optional chevron
//

import SwiftUI

/// A card-style row with an icon in a colored circle, text, and optional chevron.
/// Use for navigation links (e.g. task summary, grocery list) inside NavigationLink or Button.
///
/// Usage:
/// ```swift
/// NavigationLink(destination: TaskListView(...)) {
///   DSInfoRow(icon: "checklist", text: "3 tasks remaining", color: .orange)
/// }
/// .buttonStyle(.plain)
///
/// DSInfoRow(icon: "cart.fill", text: "Grocery List (2 needed)", color: .blue, showChevron: false)
/// ```
struct DSInfoRow: View {
  let icon: String
  let text: String
  let color: Color
  var showChevron: Bool = true
  var strikethrough: Bool = false

  var body: some View {
    HStack(spacing: DSTheme.Spacing.md) {
      ZStack {
        Circle()
          .fill(color.opacity(0.15))
          .frame(width: 36, height: 36)

        Image(systemName: icon)
          .font(.system(size: 15, weight: .medium))
          .foregroundColor(color)
      }

      Text(text)
        .font(DSTheme.Typography.label)
        .foregroundColor(DSTheme.Colors.textPrimary)
        .strikethrough(strikethrough)

      Spacer()

      if showChevron {
        Image(systemName: "chevron.right")
          .font(.system(size: 13, weight: .semibold))
          .foregroundColor(DSTheme.Colors.textTertiary)
      }
    }
    .padding(DSTheme.Spacing.md)
    .background(DSTheme.Colors.cardBackground)
    .cornerRadius(DSTheme.Radius.lg)
    .overlay(
      RoundedRectangle(cornerRadius: DSTheme.Radius.lg)
        .stroke(DSTheme.Colors.border, lineWidth: 1)
    )
  }
}

// MARK: - Preview

#Preview("DSInfoRow") {
  VStack(spacing: DSTheme.Spacing.md) {
    DSInfoRow(icon: "checklist", text: "3 tasks remaining", color: DSTheme.Colors.warning)
    DSInfoRow(
      icon: "cart.fill", text: "Grocery List (2 ingredients needed)", color: DSTheme.Colors.primary)
    DSInfoRow(
      icon: "checkmark.circle", text: "All ingredients acquired", color: DSTheme.Colors.success)
    DSInfoRow(
      icon: "info.circle", text: "Without chevron", color: DSTheme.Colors.info, showChevron: false)
  }
  .padding()
}
