//
//  DSListRow.swift
//  ios
//
//  Design System List Row Component - Consistent list row styling
//

import SwiftUI

// MARK: - List Row Style

enum DSListRowStyle {
  case standard
  case card
  case plain
}

// MARK: - DSListRow Component

/// A consistent list row with optional leading/trailing content and chevron.
///
/// Usage:
/// ```swift
/// DSListRow(title: "Settings")
/// DSListRow(title: "Profile", subtitle: "Edit your profile")
/// DSListRow(title: "Notifications", icon: "bell")
/// DSListRow(title: "Account", icon: "person", showChevron: true)
///
/// // As navigation link
/// DSListRow(title: "Settings", destination: SettingsView())
///
/// // With custom trailing
/// DSListRow(title: "Theme") {
///   Toggle("", isOn: $isDarkMode)
/// }
/// ```
struct DSListRow<Trailing: View>: View {
  let title: String
  let subtitle: String?
  let icon: String?
  let iconColor: Color
  let style: DSListRowStyle
  let showChevron: Bool
  let trailing: (() -> Trailing)?
  let action: (() -> Void)?

  init(
    title: String,
    subtitle: String? = nil,
    icon: String? = nil,
    iconColor: Color = DSTheme.Colors.textSecondary,
    style: DSListRowStyle = .standard,
    showChevron: Bool = false,
    action: (() -> Void)? = nil,
    @ViewBuilder trailing: @escaping () -> Trailing
  ) {
    self.title = title
    self.subtitle = subtitle
    self.icon = icon
    self.iconColor = iconColor
    self.style = style
    self.showChevron = showChevron
    self.trailing = trailing
    self.action = action
  }

  var body: some View {
    Group {
      if let action = action {
        Button(action: action) {
          rowContent
        }
        .buttonStyle(.plain)
      } else {
        rowContent
      }
    }
    .applyStyle(style)
  }

  private var rowContent: some View {
    HStack(spacing: DSTheme.Spacing.md) {
      // Leading icon
      if let icon = icon {
        Image(systemName: icon)
          .font(.system(size: DSTheme.IconSize.md))
          .foregroundColor(iconColor)
          .frame(width: DSTheme.IconSize.lg)
      }

      // Title and subtitle
      VStack(alignment: .leading, spacing: DSTheme.Spacing.xxs) {
        Text(title)
          .font(DSTheme.Typography.label)
          .foregroundColor(DSTheme.Colors.textPrimary)

        if let subtitle = subtitle {
          Text(subtitle)
            .font(DSTheme.Typography.caption)
            .foregroundColor(DSTheme.Colors.textSecondary)
        }
      }

      Spacer()

      // Trailing content
      if let trailing = trailing {
        trailing()
      }

      // Chevron
      if showChevron {
        Image(systemName: "chevron.right")
          .font(.system(size: DSTheme.IconSize.sm))
          .foregroundColor(DSTheme.Colors.textTertiary)
      }
    }
    .contentShape(Rectangle())
  }
}

// MARK: - Convenience initializer without trailing

extension DSListRow where Trailing == EmptyView {
  init(
    title: String,
    subtitle: String? = nil,
    icon: String? = nil,
    iconColor: Color = DSTheme.Colors.textSecondary,
    style: DSListRowStyle = .standard,
    showChevron: Bool = false,
    action: (() -> Void)? = nil
  ) {
    self.title = title
    self.subtitle = subtitle
    self.icon = icon
    self.iconColor = iconColor
    self.style = style
    self.showChevron = showChevron
    self.trailing = nil
    self.action = action
  }
}

// MARK: - Navigation Link Row

/// A list row that navigates to a destination.
struct DSListRowLink<Destination: View>: View {
  let title: String
  let subtitle: String?
  let icon: String?
  let iconColor: Color
  let style: DSListRowStyle
  let destination: Destination

  init(
    title: String,
    subtitle: String? = nil,
    icon: String? = nil,
    iconColor: Color = DSTheme.Colors.textSecondary,
    style: DSListRowStyle = .standard,
    destination: Destination
  ) {
    self.title = title
    self.subtitle = subtitle
    self.icon = icon
    self.iconColor = iconColor
    self.style = style
    self.destination = destination
  }

  var body: some View {
    NavigationLink(destination: destination) {
      HStack(spacing: DSTheme.Spacing.md) {
        if let icon = icon {
          Image(systemName: icon)
            .font(.system(size: DSTheme.IconSize.md))
            .foregroundColor(iconColor)
            .frame(width: DSTheme.IconSize.lg)
        }

        VStack(alignment: .leading, spacing: DSTheme.Spacing.xxs) {
          Text(title)
            .font(DSTheme.Typography.label)
            .foregroundColor(DSTheme.Colors.textPrimary)

          if let subtitle = subtitle {
            Text(subtitle)
              .font(DSTheme.Typography.caption)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }
        }

        Spacer()

        Image(systemName: "chevron.right")
          .font(.system(size: DSTheme.IconSize.sm))
          .foregroundColor(DSTheme.Colors.textTertiary)
      }
    }
    .applyStyle(style)
  }
}

// MARK: - Style Modifier

extension View {
  @ViewBuilder
  fileprivate func applyStyle(_ style: DSListRowStyle) -> some View {
    switch style {
    case .standard:
      self
        .padding(.vertical, DSTheme.Spacing.md)
        .padding(.horizontal, DSTheme.Spacing.sm)
    case .card:
      self
        .padding(DSTheme.Spacing.lg)
        .background(DSTheme.Colors.cardBackground)
        .cornerRadius(DSTheme.Radius.md)
    case .plain:
      self
        .padding(.vertical, DSTheme.Spacing.sm)
    }
  }
}

// MARK: - Divider Row

/// A simple divider for use in lists.
struct DSListDivider: View {
  let inset: CGFloat

  init(inset: CGFloat = 0) {
    self.inset = inset
  }

  var body: some View {
    Divider()
      .padding(.leading, inset)
  }
}

// MARK: - Preview

#Preview("List Rows") {
  NavigationStack {
    ScrollView {
      VStack(spacing: 0) {
        Text("Standard Style").font(.caption).foregroundColor(.secondary).frame(
          maxWidth: .infinity, alignment: .leading
        ).padding()

        VStack(spacing: 0) {
          DSListRow(title: "Simple Row")
          DSListDivider()
          DSListRow(title: "With Subtitle", subtitle: "This is a subtitle")
          DSListDivider()
          DSListRow(title: "With Icon", icon: "star")
          DSListDivider()
          DSListRow(title: "With Chevron", icon: "gear", showChevron: true)
          DSListDivider()
          DSListRow(title: "Tappable", icon: "hand.tap", showChevron: true) {
            print("Tapped")
          }
        }
        .padding(.horizontal)

        Text("Card Style").font(.caption).foregroundColor(.secondary).frame(
          maxWidth: .infinity, alignment: .leading
        ).padding()

        VStack(spacing: DSTheme.Spacing.sm) {
          DSListRow(title: "Card Row", icon: "folder", style: .card, showChevron: true)
          DSListRow(
            title: "With Subtitle", subtitle: "Additional info", icon: "doc", style: .card,
            showChevron: true)
        }
        .padding(.horizontal)

        Text("With Trailing Content").font(.caption).foregroundColor(.secondary).frame(
          maxWidth: .infinity, alignment: .leading
        ).padding()

        VStack(spacing: 0) {
          DSListRow(title: "Toggle Option", icon: "moon") {
            Toggle("", isOn: .constant(true))
              .labelsHidden()
          }
          DSListDivider()
          DSListRow(title: "Count", icon: "bell") {
            Text("5")
              .font(DSTheme.Typography.caption)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }
        }
        .padding(.horizontal)

        Text("Navigation Link").font(.caption).foregroundColor(.secondary).frame(
          maxWidth: .infinity, alignment: .leading
        ).padding()

        VStack(spacing: DSTheme.Spacing.sm) {
          DSListRowLink(
            title: "Go to Settings",
            icon: "gear",
            style: .card,
            destination: Text("Settings View")
          )
        }
        .padding(.horizontal)
      }
    }
    .navigationTitle("List Rows")
  }
}
