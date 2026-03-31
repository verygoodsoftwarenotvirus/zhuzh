//
//  DSEmptyState.swift
//  ios
//
//  Design System Empty State Component - Consistent empty state display
//

import SwiftUI

// MARK: - Empty State Size

enum DSEmptyStateSize {
  case compact
  case standard
  case large

  var iconSize: CGFloat {
    switch self {
    case .compact:
      return DSTheme.IconSize.xl
    case .standard:
      return DSTheme.IconSize.xxl
    case .large:
      return 64
    }
  }

  var titleFont: Font {
    switch self {
    case .compact:
      return DSTheme.Typography.label
    case .standard:
      return DSTheme.Typography.title3
    case .large:
      return DSTheme.Typography.title2
    }
  }

  var messageFont: Font {
    switch self {
    case .compact:
      return DSTheme.Typography.caption
    case .standard, .large:
      return DSTheme.Typography.body
    }
  }

  var spacing: CGFloat {
    switch self {
    case .compact:
      return DSTheme.Spacing.sm
    case .standard:
      return DSTheme.Spacing.md
    case .large:
      return DSTheme.Spacing.lg
    }
  }
}

// MARK: - DSEmptyState Component

/// A consistent empty state view with icon, title, message, and optional action.
///
/// Usage:
/// ```swift
/// DSEmptyState(
///   icon: "tray",
///   title: "No Items",
///   message: "You haven't added any items yet."
/// )
///
/// DSEmptyState(
///   icon: "magnifyingglass",
///   title: "No Results",
///   message: "Try adjusting your search.",
///   actionTitle: "Clear Search",
///   action: { clearSearch() }
/// )
///
/// // Compact for inline use
/// DSEmptyState(
///   icon: "doc",
///   message: "No documents",
///   size: .compact
/// )
/// ```
struct DSEmptyState: View {
  let icon: String?
  let title: String?
  let message: String
  let actionTitle: String?
  let size: DSEmptyStateSize
  let action: (() -> Void)?

  init(
    icon: String? = nil,
    title: String? = nil,
    message: String,
    actionTitle: String? = nil,
    size: DSEmptyStateSize = .standard,
    action: (() -> Void)? = nil
  ) {
    self.icon = icon
    self.title = title
    self.message = message
    self.actionTitle = actionTitle
    self.size = size
    self.action = action
  }

  var body: some View {
    VStack(spacing: size.spacing) {
      if let icon = icon {
        Image(systemName: icon)
          .font(.system(size: size.iconSize))
          .foregroundColor(DSTheme.Colors.textTertiary)
      }

      VStack(spacing: DSTheme.Spacing.xs) {
        if let title = title {
          Text(title)
            .font(size.titleFont)
            .foregroundColor(DSTheme.Colors.textPrimary)
        }

        Text(message)
          .font(size.messageFont)
          .foregroundColor(DSTheme.Colors.textSecondary)
          .multilineTextAlignment(.center)
      }

      if let actionTitle = actionTitle, let action = action {
        DSButton(actionTitle, style: .primary, size: .small, action: action)
          .padding(.top, DSTheme.Spacing.sm)
      }
    }
    .padding(size == .compact ? DSTheme.Spacing.md : DSTheme.Spacing.xl)
    .frame(maxWidth: size == .compact ? nil : .infinity)
  }
}

// MARK: - Preview

#Preview("Empty States") {
  ScrollView {
    VStack(spacing: DSTheme.Spacing.xxl) {
      DSEmptyState(
        icon: "tray",
        title: "No Items",
        message: "You haven't added any items yet."
      )

      Divider()

      DSEmptyState(
        icon: "magnifyingglass",
        title: "No Results Found",
        message: "Try adjusting your search or filters.",
        actionTitle: "Clear Filters"
      ) {
        print("Clear filters")
      }

      Divider()

      DSEmptyState(
        icon: "person.2",
        title: "No Members",
        message: "Invite someone to join your household.",
        actionTitle: "Send Invite",
        size: .large
      ) {
        print("Send invite")
      }

      Divider()

      Text("Compact (inline)").font(.caption).foregroundColor(.secondary)
      DSEmptyState(
        icon: "doc",
        message: "No documents available",
        size: .compact
      )
    }
    .padding()
  }
}
