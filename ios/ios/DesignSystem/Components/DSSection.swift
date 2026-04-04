//
//  DSSection.swift
//  ios
//
//  Design System Section Component - Section with header and optional subtitle
//

import SwiftUI

// MARK: - Section Style

enum DSSectionStyle {
  case standard
  case compact
  case card
}

// MARK: - DSSection Component

/// A section container with a styled header and optional subtitle.
///
/// Usage:
/// ```swift
/// DSSection("Members") {
///   ForEach(members) { MemberRow(member: $0) }
/// }
///
/// DSSection("Household", subtitle: "People in your household") {
///   ForEach(members) { MemberRow(member: $0) }
/// }
///
/// DSSection("Tasks", style: .compact) {
///   TaskList()
/// }
/// ```
struct DSSection<Content: View>: View {
  let title: String
  let subtitle: String?
  let style: DSSectionStyle
  let spacing: CGFloat
  let showDivider: Bool
  @ViewBuilder let content: () -> Content

  init(
    _ title: String,
    subtitle: String? = nil,
    style: DSSectionStyle = .standard,
    spacing: CGFloat? = nil,
    showDivider: Bool = false,
    @ViewBuilder content: @escaping () -> Content
  ) {
    self.title = title
    self.subtitle = subtitle
    self.style = style
    self.spacing = spacing ?? (style == .compact ? DSTheme.Spacing.sm : DSTheme.Spacing.md)
    self.showDivider = showDivider
    self.content = content
  }

  var body: some View {
    VStack(alignment: .leading, spacing: spacing) {
      // Header
      DSSectionHeader(title: title, subtitle: subtitle, style: headerStyle)

      if showDivider {
        Divider()
      }

      // Content
      if style == .card {
        DSCard {
          VStack(alignment: .leading, spacing: spacing) {
            content()
          }
          .frame(maxWidth: .infinity, alignment: .leading)
        }
      } else {
        content()
      }
    }
  }

  private var headerStyle: DSSectionHeaderStyle {
    switch style {
    case .compact:
      return .compact
    default:
      return .standard
    }
  }
}

// MARK: - Section Header Style

enum DSSectionHeaderStyle {
  case standard
  case compact
  case large
  /// Caption-sized label for drawer sections and small headers
  case label
}

// MARK: - DSSectionHeader Component

/// A standalone section header that can be used independently.
///
/// Usage:
/// ```swift
/// DSSectionHeader(title: "Members")
/// DSSectionHeader(title: "Details", subtitle: "Additional information")
/// DSSectionHeader(title: "Upcoming Events", systemImage: "calendar.badge.clock")
/// DSSectionHeader(title: "Upcoming", systemImage: "calendar") {
///   Text("3 plans").font(.caption).foregroundColor(.secondary)
/// }
/// ```
struct DSSectionHeader<Trailing: View>: View {
  let title: String
  let subtitle: String?
  let systemImage: String?
  let style: DSSectionHeaderStyle
  let titleColor: Color
  @ViewBuilder let trailing: () -> Trailing

  init(
    title: String,
    subtitle: String? = nil,
    systemImage: String? = nil,
    style: DSSectionHeaderStyle = .standard,
    titleColor: Color? = nil,
    @ViewBuilder trailing: @escaping () -> Trailing
  ) {
    self.title = title
    self.subtitle = subtitle
    self.systemImage = systemImage
    self.style = style
    self.titleColor =
      titleColor ?? (style == .label ? DSTheme.Colors.textSecondary : DSTheme.Colors.textPrimary)
    self.trailing = trailing
  }

  var body: some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.xs) {
      headerContent

      if let subtitle = subtitle, !subtitle.isEmpty {
        Text(subtitle)
          .font(subtitleFont)
          .foregroundColor(DSTheme.Colors.textSecondary)
      }
    }
    .padding(.horizontal, horizontalPadding)
  }

  private var titleFont: Font {
    switch style {
    case .large:
      return DSTheme.Typography.title1
    case .standard:
      return DSTheme.Typography.title2
    case .compact:
      return DSTheme.Typography.title3
    case .label:
      return DSTheme.Typography.caption.weight(.semibold)
    }
  }

  private var subtitleFont: Font {
    switch style {
    case .large:
      return DSTheme.Typography.body
    case .standard, .compact, .label:
      return DSTheme.Typography.bodySmall
    }
  }

  private var horizontalPadding: CGFloat {
    switch style {
    case .compact:
      return DSTheme.Spacing.xxs
    case .label:
      return DSTheme.Spacing.md
    default:
      return DSTheme.Spacing.xs
    }
  }

  @ViewBuilder
  private var headerContent: some View {
    if let systemImage = systemImage {
      HStack {
        Label(title, systemImage: systemImage)
          .font(titleFont)
          .foregroundColor(titleColor)

        Spacer()

        trailing()
      }
    } else {
      HStack {
        Text(title)
          .font(titleFont)
          .foregroundColor(titleColor)

        Spacer()

        trailing()
      }
    }
  }
}

// MARK: - DSSectionHeader without trailing

extension DSSectionHeader where Trailing == EmptyView {
  init(
    title: String,
    subtitle: String? = nil,
    systemImage: String? = nil,
    style: DSSectionHeaderStyle = .standard,
    titleColor: Color? = nil
  ) {
    self.title = title
    self.subtitle = subtitle
    self.systemImage = systemImage
    self.style = style
    self.titleColor =
      titleColor ?? (style == .label ? DSTheme.Colors.textSecondary : DSTheme.Colors.textPrimary)
    self.trailing = { EmptyView() }
  }
}

// MARK: - Rule-Flanked Header

/// A section header with horizontal rules flanking the title.
///
/// Usage:
/// ```swift
/// DSRuleFlankedHeader(title: "Up Next", color: .orange)
/// DSRuleFlankedHeader(title: "Have", color: .green)
/// ```
struct DSRuleFlankedHeader: View {
  let title: String
  let color: Color
  var strikethrough: Bool = false

  init(title: String, color: Color, strikethrough: Bool = false) {
    self.title = title
    self.color = color
    self.strikethrough = strikethrough
  }

  var body: some View {
    HStack(spacing: 12) {
      Rectangle()
        .fill(Color.secondary.opacity(0.5))
        .frame(height: 1)
        .frame(maxWidth: .infinity)
      Text(title)
        .font(.headline)
        .fontWeight(.semibold)
        .foregroundColor(color)
        .strikethrough(strikethrough)
      Rectangle()
        .fill(Color.secondary.opacity(0.5))
        .frame(height: 1)
        .frame(maxWidth: .infinity)
    }
    .padding(.vertical, 8)
    .frame(maxWidth: .infinity)
    .background(Color(uiColor: .systemBackground))
  }
}

// MARK: - Empty Section Content

/// A helper view for displaying empty state within a section.
struct DSSectionEmptyContent: View {
  let message: String
  let icon: String?

  init(_ message: String, icon: String? = nil) {
    self.message = message
    self.icon = icon
  }

  var body: some View {
    HStack(spacing: DSTheme.Spacing.sm) {
      if let icon = icon {
        Image(systemName: icon)
          .font(.subheadline)
          .foregroundColor(DSTheme.Colors.textSecondary)
      }

      Text(message)
        .font(DSTheme.Typography.body)
        .foregroundColor(DSTheme.Colors.textSecondary)
    }
    .padding(.vertical, DSTheme.Spacing.sm)
    .padding(.horizontal, DSTheme.Spacing.xs)
  }
}

// MARK: - Preview

#Preview("Section Styles") {
  ScrollView {
    VStack(spacing: DSTheme.Spacing.xl) {
      DSSection("Standard Section") {
        Text("Content goes here")
          .dsCard()
      }

      DSSection("With Subtitle", subtitle: "This is additional context") {
        Text("More content")
          .dsCard()
      }

      DSSection("Compact Style", style: .compact) {
        VStack(spacing: DSTheme.Spacing.sm) {
          Text("Item 1")
          Text("Item 2")
          Text("Item 3")
        }
        .dsCard()
      }

      DSSection("Card Style", subtitle: "Content wrapped in a card", style: .card) {
        VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
          Text("Item 1")
          Divider()
          Text("Item 2")
          Divider()
          Text("Item 3")
        }
      }

      DSSection("Empty Section") {
        DSSectionEmptyContent("No items available", icon: "tray")
      }
    }
    .padding()
  }
}
