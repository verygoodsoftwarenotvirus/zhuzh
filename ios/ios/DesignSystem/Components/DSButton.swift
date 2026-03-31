//
//  DSButton.swift
//  ios
//
//  Design System Button Component - Consistent button styling with variants
//

import SwiftUI

// MARK: - Button Style

enum DSButtonStyle {
  case primary
  case secondary
  case tertiary
  case destructive
  case destructiveGhost
  case ghost
  case outline

  var backgroundColor: Color {
    switch self {
    case .primary:
      return DSTheme.Colors.primary
    case .secondary:
      return DSTheme.Colors.secondary
    case .tertiary:
      return DSTheme.Colors.tertiary
    case .destructive:
      return DSTheme.Colors.error
    case .destructiveGhost, .ghost, .outline:
      return .clear
    }
  }

  var foregroundColor: Color {
    switch self {
    case .primary, .secondary, .tertiary, .destructive:
      return DSTheme.Colors.textOnPrimary
    case .destructiveGhost:
      return DSTheme.Colors.error
    case .ghost:
      return DSTheme.Colors.primary
    case .outline:
      return DSTheme.Colors.textPrimary
    }
  }

  var borderColor: Color? {
    switch self {
    case .outline:
      return DSTheme.Colors.border
    default:
      return nil
    }
  }

  var disabledBackgroundColor: Color {
    switch self {
    case .destructiveGhost, .ghost, .outline:
      return .clear
    default:
      return Color.gray.opacity(0.3)
    }
  }

  var disabledForegroundColor: Color {
    return Color.gray
  }
}

// MARK: - Button Size

enum DSButtonSize {
  case small
  case medium
  case large

  var verticalPadding: CGFloat {
    switch self {
    case .small:
      return DSTheme.Spacing.sm
    case .medium:
      return DSTheme.Spacing.md
    case .large:
      return DSTheme.Spacing.lg
    }
  }

  var horizontalPadding: CGFloat {
    switch self {
    case .small:
      return DSTheme.Spacing.md
    case .medium:
      return DSTheme.Spacing.lg
    case .large:
      return DSTheme.Spacing.xl
    }
  }

  var font: Font {
    switch self {
    case .small:
      return DSTheme.Typography.buttonSmall
    case .medium, .large:
      return DSTheme.Typography.button
    }
  }

  var iconSize: CGFloat {
    switch self {
    case .small:
      return DSTheme.IconSize.sm
    case .medium:
      return DSTheme.IconSize.md
    case .large:
      return DSTheme.IconSize.lg
    }
  }
}

// MARK: - DSButton Component

/// A styled button with consistent appearance and behavior.
///
/// Usage:
/// ```swift
/// DSButton("Submit") { handleSubmit() }
/// DSButton("Delete", style: .destructive) { handleDelete() }
/// DSButton("View Details", icon: "chevron.right", style: .ghost) { showDetails() }
///
/// // Async action
/// DSButton("Save", isLoading: isSaving) { await save() }
///
/// // Full width
/// DSButton("Continue", fullWidth: true) { next() }
/// ```
///
/// - Note: For sheet/toolbar Cancel and Save actions, use SwiftUI's `Button("Cancel")` /
///   `Button("Save")` in `cancellationAction` / `confirmationAction`; they receive system styling automatically.
struct DSButton: View {
  let title: String
  let icon: String?
  let iconPosition: IconPosition
  let style: DSButtonStyle
  let size: DSButtonSize
  let fullWidth: Bool
  let isLoading: Bool
  let isDisabled: Bool
  let action: () -> Void

  enum IconPosition {
    case leading
    case trailing
  }

  init(
    _ title: String,
    icon: String? = nil,
    iconPosition: IconPosition = .leading,
    style: DSButtonStyle = .primary,
    size: DSButtonSize = .medium,
    fullWidth: Bool = false,
    isLoading: Bool = false,
    isDisabled: Bool = false,
    action: @escaping () -> Void
  ) {
    self.title = title
    self.icon = icon
    self.iconPosition = iconPosition
    self.style = style
    self.size = size
    self.fullWidth = fullWidth
    self.isLoading = isLoading
    self.isDisabled = isDisabled
    self.action = action
  }

  var body: some View {
    Button(action: action) {
      HStack(spacing: DSTheme.Spacing.sm) {
        if isLoading {
          ProgressView()
            .progressViewStyle(CircularProgressViewStyle(tint: effectiveForegroundColor))
            .scaleEffect(0.8)
        } else if let icon = icon, iconPosition == .leading {
          Image(systemName: icon)
            .font(.system(size: size.iconSize))
        }

        Text(title)
          .font(size.font)

        if !isLoading, let icon = icon, iconPosition == .trailing {
          Image(systemName: icon)
            .font(.system(size: size.iconSize))
        }
      }
      .frame(maxWidth: fullWidth ? .infinity : nil)
      .padding(.vertical, size.verticalPadding)
      .padding(.horizontal, size.horizontalPadding)
      .foregroundColor(effectiveForegroundColor)
      .background(effectiveBackgroundColor)
      .cornerRadius(DSTheme.Radius.md)
      .overlay(
        RoundedRectangle(cornerRadius: DSTheme.Radius.md)
          .stroke(effectiveBorderColor, lineWidth: style.borderColor != nil ? 1 : 0)
      )
    }
    .disabled(isDisabled || isLoading)
  }

  private var effectiveBackgroundColor: Color {
    isDisabled ? style.disabledBackgroundColor : style.backgroundColor
  }

  private var effectiveForegroundColor: Color {
    isDisabled ? style.disabledForegroundColor : style.foregroundColor
  }

  private var effectiveBorderColor: Color {
    isDisabled ? Color.gray.opacity(0.3) : (style.borderColor ?? .clear)
  }
}

// MARK: - Navigation Link Button

/// A navigation link styled as a button.
///
/// Usage:
/// ```swift
/// DSButtonLink("View Recipes", icon: "book.closed", destination: RecipeListView())
/// DSButtonLink("Settings", style: .secondary, destination: SettingsView())
/// ```
struct DSButtonLink<Destination: View>: View {
  let title: String
  let icon: String?
  let iconPosition: DSButton.IconPosition
  let style: DSButtonStyle
  let size: DSButtonSize
  let fullWidth: Bool
  let destination: Destination

  init(
    _ title: String,
    icon: String? = nil,
    iconPosition: DSButton.IconPosition = .leading,
    style: DSButtonStyle = .primary,
    size: DSButtonSize = .medium,
    fullWidth: Bool = false,
    destination: Destination
  ) {
    self.title = title
    self.icon = icon
    self.iconPosition = iconPosition
    self.style = style
    self.size = size
    self.fullWidth = fullWidth
    self.destination = destination
  }

  var body: some View {
    NavigationLink(destination: destination) {
      HStack(spacing: DSTheme.Spacing.sm) {
        if let icon = icon, iconPosition == .leading {
          Image(systemName: icon)
            .font(.system(size: size.iconSize))
        }

        Text(title)
          .font(size.font)

        if let icon = icon, iconPosition == .trailing {
          Image(systemName: icon)
            .font(.system(size: size.iconSize))
        }
      }
      .frame(maxWidth: fullWidth ? .infinity : nil)
      .padding(.vertical, size.verticalPadding)
      .padding(.horizontal, size.horizontalPadding)
      .foregroundColor(style.foregroundColor)
      .background(style.backgroundColor)
      .cornerRadius(DSTheme.Radius.md)
      .overlay(
        RoundedRectangle(cornerRadius: DSTheme.Radius.md)
          .stroke(style.borderColor ?? .clear, lineWidth: style.borderColor != nil ? 1 : 0)
      )
    }
  }
}

// MARK: - Icon Button

/// A circular icon-only button.
///
/// Usage:
/// ```swift
/// DSIconButton("plus") { addItem() }
/// DSIconButton("trash", style: .destructive) { deleteItem() }
/// ```
struct DSIconButton: View {
  let icon: String
  let style: DSButtonStyle
  let size: DSButtonSize
  let isDisabled: Bool
  let action: () -> Void

  init(
    _ icon: String,
    style: DSButtonStyle = .ghost,
    size: DSButtonSize = .medium,
    isDisabled: Bool = false,
    action: @escaping () -> Void
  ) {
    self.icon = icon
    self.style = style
    self.size = size
    self.isDisabled = isDisabled
    self.action = action
  }

  var body: some View {
    Button(action: action) {
      Image(systemName: icon)
        .font(.system(size: size.iconSize))
        .foregroundColor(isDisabled ? style.disabledForegroundColor : style.foregroundColor)
        .frame(width: buttonSize, height: buttonSize)
        .background(isDisabled ? style.disabledBackgroundColor : style.backgroundColor)
        .clipShape(Circle())
    }
    .disabled(isDisabled)
  }

  private var buttonSize: CGFloat {
    switch size {
    case .small:
      return 32
    case .medium:
      return 40
    case .large:
      return 48
    }
  }
}

// MARK: - Preview

#Preview("Button Styles") {
  ScrollView {
    VStack(spacing: DSTheme.Spacing.lg) {
      Group {
        Text("Primary").font(.caption).foregroundColor(.secondary)
        DSButton("Primary Button", icon: "star.fill") { print("Tapped") }

        Text("Secondary").font(.caption).foregroundColor(.secondary)
        DSButton("Secondary Button", style: .secondary) { print("Tapped") }

        Text("Tertiary").font(.caption).foregroundColor(.secondary)
        DSButton("Tertiary Button", style: .tertiary) { print("Tapped") }

        Text("Destructive").font(.caption).foregroundColor(.secondary)
        DSButton("Delete", icon: "trash", style: .destructive) { print("Tapped") }

        Text("Destructive Ghost").font(.caption).foregroundColor(.secondary)
        DSButton("Sign Out", style: .destructiveGhost) { print("Tapped") }

        Text("Ghost").font(.caption).foregroundColor(.secondary)
        DSButton("Ghost Button", style: .ghost) { print("Tapped") }

        Text("Outline").font(.caption).foregroundColor(.secondary)
        DSButton("Outline Button", style: .outline) { print("Tapped") }
      }

      Divider()

      Group {
        Text("Full Width").font(.caption).foregroundColor(.secondary)
        DSButton("Full Width Button", fullWidth: true) { print("Tapped") }

        Text("Loading").font(.caption).foregroundColor(.secondary)
        DSButton("Saving...", isLoading: true) { print("Tapped") }

        Text("Disabled").font(.caption).foregroundColor(.secondary)
        DSButton("Disabled", isDisabled: true) { print("Tapped") }
      }

      Divider()

      Group {
        Text("Sizes").font(.caption).foregroundColor(.secondary)
        HStack {
          DSButton("Small", size: .small) { print("Tapped") }
          DSButton("Medium", size: .medium) { print("Tapped") }
          DSButton("Large", size: .large) { print("Tapped") }
        }
      }

      Divider()

      Group {
        Text("Icon Buttons").font(.caption).foregroundColor(.secondary)
        HStack(spacing: DSTheme.Spacing.md) {
          DSIconButton("plus", style: .primary) { print("Plus") }
          DSIconButton("heart.fill", style: .destructive) { print("Heart") }
          DSIconButton("gear") { print("Gear") }
          DSIconButton("trash", style: .ghost, isDisabled: true) { print("Trash") }
        }
      }
    }
    .padding()
  }
}
