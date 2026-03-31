//
//  DSStatusBadge.swift
//  ios
//
//  Design System Status Badge Component - Status indicators with color coding
//

import SwiftUI

// MARK: - Status Type

enum DSStatusType {
  case pending
  case accepted
  case rejected
  case cancelled
  case success
  case warning
  case error
  case info
  case custom(String, Color)

  var label: String {
    switch self {
    case .pending:
      return "Pending"
    case .accepted:
      return "Accepted"
    case .rejected:
      return "Rejected"
    case .cancelled:
      return "Cancelled"
    case .success:
      return "Success"
    case .warning:
      return "Warning"
    case .error:
      return "Error"
    case .info:
      return "Info"
    case .custom(let label, _):
      return label
    }
  }

  var color: Color {
    switch self {
    case .pending:
      return DSTheme.Colors.statusPending
    case .accepted, .success:
      return DSTheme.Colors.statusAccepted
    case .rejected, .error:
      return DSTheme.Colors.statusRejected
    case .cancelled:
      return DSTheme.Colors.statusCancelled
    case .warning:
      return DSTheme.Colors.warning
    case .info:
      return DSTheme.Colors.info
    case .custom(_, let color):
      return color
    }
  }

  /// Create a status type from a raw string value (e.g., from API)
  static func from(_ rawValue: String) -> DSStatusType {
    switch rawValue.lowercased() {
    case "pending":
      return .pending
    case "accepted":
      return .accepted
    case "rejected":
      return .rejected
    case "cancelled", "canceled":
      return .cancelled
    case "success":
      return .success
    case "warning":
      return .warning
    case "error":
      return .error
    case "info":
      return .info
    default:
      return .custom(rawValue.isEmpty ? "Unknown" : rawValue.capitalized, .secondary)
    }
  }
}

// MARK: - Badge Style

enum DSBadgeStyle {
  case dot  // Colored dot + text
  case pill  // Colored background pill
  case outline  // Outlined pill
  case minimal  // Just text with color
}

// MARK: - Badge Size

enum DSBadgeSize {
  case small
  case medium
  case large

  var dotSize: CGFloat {
    switch self {
    case .small:
      return 6
    case .medium:
      return 8
    case .large:
      return 10
    }
  }

  var font: Font {
    switch self {
    case .small:
      return DSTheme.Typography.captionSmall
    case .medium:
      return DSTheme.Typography.caption
    case .large:
      return DSTheme.Typography.labelSmall
    }
  }

  var padding: EdgeInsets {
    switch self {
    case .small:
      return EdgeInsets(
        top: DSTheme.Spacing.xxs,
        leading: DSTheme.Spacing.sm,
        bottom: DSTheme.Spacing.xxs,
        trailing: DSTheme.Spacing.sm
      )
    case .medium:
      return EdgeInsets(
        top: DSTheme.Spacing.xs,
        leading: DSTheme.Spacing.sm,
        bottom: DSTheme.Spacing.xs,
        trailing: DSTheme.Spacing.sm
      )
    case .large:
      return EdgeInsets(
        top: DSTheme.Spacing.sm,
        leading: DSTheme.Spacing.md,
        bottom: DSTheme.Spacing.sm,
        trailing: DSTheme.Spacing.md
      )
    }
  }
}

// MARK: - DSStatusBadge Component

/// A status indicator badge with consistent styling.
///
/// Usage:
/// ```swift
/// DSStatusBadge(.pending)
/// DSStatusBadge(.accepted, style: .pill)
/// DSStatusBadge(.custom("Active", .purple), size: .large)
///
/// // From raw string
/// DSStatusBadge(status: invitation.status)
/// ```
struct DSStatusBadge: View {
  let status: DSStatusType
  let style: DSBadgeStyle
  let size: DSBadgeSize

  init(
    _ status: DSStatusType,
    style: DSBadgeStyle = .dot,
    size: DSBadgeSize = .medium
  ) {
    self.status = status
    self.style = style
    self.size = size
  }

  /// Convenience initializer for raw string status values
  init(
    status: String,
    style: DSBadgeStyle = .dot,
    size: DSBadgeSize = .medium
  ) {
    self.status = DSStatusType.from(status)
    self.style = style
    self.size = size
  }

  var body: some View {
    Group {
      switch style {
      case .dot:
        dotBadge
      case .pill:
        pillBadge
      case .outline:
        outlineBadge
      case .minimal:
        minimalBadge
      }
    }
  }

  private var dotBadge: some View {
    HStack(spacing: DSTheme.Spacing.sm) {
      Circle()
        .fill(status.color)
        .frame(width: size.dotSize, height: size.dotSize)

      Text(status.label)
        .font(size.font)
        .foregroundColor(DSTheme.Colors.textSecondary)
    }
  }

  private var pillBadge: some View {
    Text(status.label)
      .font(size.font)
      .foregroundColor(.white)
      .padding(size.padding)
      .background(status.color)
      .cornerRadius(DSTheme.Radius.full)
  }

  private var outlineBadge: some View {
    Text(status.label)
      .font(size.font)
      .foregroundColor(status.color)
      .padding(size.padding)
      .background(status.color.opacity(0.1))
      .cornerRadius(DSTheme.Radius.full)
      .overlay(
        RoundedRectangle(cornerRadius: DSTheme.Radius.full)
          .stroke(status.color, lineWidth: 1)
      )
  }

  private var minimalBadge: some View {
    Text(status.label)
      .font(size.font)
      .foregroundColor(status.color)
  }
}

// MARK: - Status Dot Only

/// A simple colored status dot without text.
///
/// Usage:
/// ```swift
/// DSStatusDot(.pending)
/// DSStatusDot(.accepted, size: .large)
/// ```
struct DSStatusDot: View {
  let status: DSStatusType
  let size: DSBadgeSize

  init(_ status: DSStatusType, size: DSBadgeSize = .medium) {
    self.status = status
    self.size = size
  }

  init(status: String, size: DSBadgeSize = .medium) {
    self.status = DSStatusType.from(status)
    self.size = size
  }

  var body: some View {
    Circle()
      .fill(status.color)
      .frame(width: size.dotSize, height: size.dotSize)
  }
}

// MARK: - Preview

#Preview("Status Badges") {
  ScrollView {
    VStack(spacing: DSTheme.Spacing.xl) {
      Group {
        Text("Dot Style (default)").font(.caption).foregroundColor(.secondary)
        HStack(spacing: DSTheme.Spacing.lg) {
          DSStatusBadge(.pending)
          DSStatusBadge(.accepted)
          DSStatusBadge(.rejected)
          DSStatusBadge(.cancelled)
        }
      }

      Group {
        Text("Pill Style").font(.caption).foregroundColor(.secondary)
        HStack(spacing: DSTheme.Spacing.sm) {
          DSStatusBadge(.pending, style: .pill)
          DSStatusBadge(.accepted, style: .pill)
          DSStatusBadge(.rejected, style: .pill)
        }
      }

      Group {
        Text("Outline Style").font(.caption).foregroundColor(.secondary)
        HStack(spacing: DSTheme.Spacing.sm) {
          DSStatusBadge(.pending, style: .outline)
          DSStatusBadge(.accepted, style: .outline)
          DSStatusBadge(.rejected, style: .outline)
        }
      }

      Group {
        Text("Minimal Style").font(.caption).foregroundColor(.secondary)
        HStack(spacing: DSTheme.Spacing.lg) {
          DSStatusBadge(.pending, style: .minimal)
          DSStatusBadge(.accepted, style: .minimal)
          DSStatusBadge(.rejected, style: .minimal)
        }
      }

      Divider()

      Group {
        Text("Sizes").font(.caption).foregroundColor(.secondary)
        VStack(spacing: DSTheme.Spacing.md) {
          DSStatusBadge(.pending, style: .pill, size: .small)
          DSStatusBadge(.pending, style: .pill, size: .medium)
          DSStatusBadge(.pending, style: .pill, size: .large)
        }
      }

      Divider()

      Group {
        Text("From Raw String").font(.caption).foregroundColor(.secondary)
        VStack(spacing: DSTheme.Spacing.sm) {
          DSStatusBadge(status: "pending")
          DSStatusBadge(status: "ACCEPTED")
          DSStatusBadge(status: "cancelled")
          DSStatusBadge(status: "unknown_status")
        }
      }

      Divider()

      Group {
        Text("Custom").font(.caption).foregroundColor(.secondary)
        HStack(spacing: DSTheme.Spacing.sm) {
          DSStatusBadge(.custom("Active", .purple), style: .pill)
          DSStatusBadge(.custom("Draft", .cyan), style: .outline)
          DSStatusBadge(.info)
          DSStatusBadge(.warning)
        }
      }

      Divider()

      Group {
        Text("Status Dots Only").font(.caption).foregroundColor(.secondary)
        HStack(spacing: DSTheme.Spacing.lg) {
          DSStatusDot(.pending)
          DSStatusDot(.accepted)
          DSStatusDot(.rejected)
          DSStatusDot(.cancelled)
        }
      }
    }
    .padding()
  }
}
