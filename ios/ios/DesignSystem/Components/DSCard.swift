//
//  DSCard.swift
//  ios
//
//  Design System Card Component - Consistent card styling
//

import SwiftUI

// MARK: - Card Style

enum DSCardStyle {
  case standard
  case elevated
  case outlined
  case selected
  case interactive

  var backgroundColor: Color {
    switch self {
    case .standard:
      return DSTheme.Colors.cardBackground
    case .elevated:
      return DSTheme.Colors.cardBackgroundElevated
    case .outlined:
      return Color(.systemBackground)
    case .selected:
      return DSTheme.Colors.cardBackgroundSelected
    case .interactive:
      return DSTheme.Colors.cardBackground
    }
  }

  var borderColor: Color? {
    switch self {
    case .outlined:
      return DSTheme.Colors.border
    case .selected:
      return DSTheme.Colors.borderSelected
    default:
      return nil
    }
  }

  var borderWidth: CGFloat {
    switch self {
    case .outlined:
      return 1
    case .selected:
      return 2
    default:
      return 0
    }
  }

  var shadow: ShadowStyle? {
    switch self {
    case .elevated:
      return DSTheme.Shadows.sm
    default:
      return nil
    }
  }
}

// MARK: - DSCard Component

/// A styled card container with consistent padding, background, and corner radius.
///
/// Usage:
/// ```swift
/// DSCard {
///   HStack {
///     Text("Card content")
///     Spacer()
///   }
/// }
///
/// DSCard(style: .selected) {
///   Text("Selected card")
/// }
///
/// DSCard(style: .interactive, action: { handleTap() }) {
///   Text("Tappable card")
/// }
/// ```
struct DSCard<Content: View>: View {
  let style: DSCardStyle
  let padding: CGFloat
  let cornerRadius: CGFloat
  let action: (() -> Void)?
  @ViewBuilder let content: () -> Content

  init(
    style: DSCardStyle = .standard,
    padding: CGFloat = DSTheme.Spacing.lg,
    cornerRadius: CGFloat = DSTheme.Radius.md,
    action: (() -> Void)? = nil,
    @ViewBuilder content: @escaping () -> Content
  ) {
    self.style = style
    self.padding = padding
    self.cornerRadius = cornerRadius
    self.action = action
    self.content = content
  }

  var body: some View {
    Group {
      if let action = action {
        Button(action: action) {
          cardContent
        }
        .buttonStyle(.plain)
      } else {
        cardContent
      }
    }
  }

  private var cardContent: some View {
    content()
      .padding(padding)
      .background(style.backgroundColor)
      .cornerRadius(cornerRadius)
      .overlay(
        RoundedRectangle(cornerRadius: cornerRadius)
          .stroke(style.borderColor ?? .clear, lineWidth: style.borderWidth)
      )
      .ifLet(style.shadow) { view, shadow in
        view.dsShadow(shadow)
      }
      .contentShape(Rectangle())
  }
}

// MARK: - View Modifier

struct DSCardModifier: ViewModifier {
  let style: DSCardStyle
  let padding: CGFloat
  let cornerRadius: CGFloat

  func body(content: Content) -> some View {
    content
      .padding(padding)
      .background(style.backgroundColor)
      .cornerRadius(cornerRadius)
      .overlay(
        RoundedRectangle(cornerRadius: cornerRadius)
          .stroke(style.borderColor ?? .clear, lineWidth: style.borderWidth)
      )
      .ifLet(style.shadow) { view, shadow in
        view.dsShadow(shadow)
      }
  }
}

// MARK: - View Extension

extension View {
  /// Apply card styling to any view.
  ///
  /// Usage:
  /// ```swift
  /// HStack {
  ///   Text("Content")
  ///   Spacer()
  /// }
  /// .dsCard()
  ///
  /// VStack { ... }
  /// .dsCard(style: .elevated)
  /// ```
  func dsCard(
    style: DSCardStyle = .standard,
    padding: CGFloat = DSTheme.Spacing.lg,
    cornerRadius: CGFloat = DSTheme.Radius.md
  ) -> some View {
    self.modifier(DSCardModifier(style: style, padding: padding, cornerRadius: cornerRadius))
  }
}

// MARK: - Helper Extension

extension View {
  /// Conditionally applies a transformation if the value is non-nil.
  @ViewBuilder
  func ifLet<T, Result: View>(_ value: T?, transform: (Self, T) -> Result) -> some View {
    if let value = value {
      transform(self, value)
    } else {
      self
    }
  }
}

// MARK: - Preview

#Preview("Card Styles") {
  ScrollView {
    VStack(spacing: DSTheme.Spacing.lg) {
      DSCard {
        HStack {
          Text("Standard Card")
          Spacer()
        }
      }

      DSCard(style: .elevated) {
        HStack {
          Text("Elevated Card")
          Spacer()
        }
      }

      DSCard(style: .outlined) {
        HStack {
          Text("Outlined Card")
          Spacer()
        }
      }

      DSCard(style: .selected) {
        HStack {
          Text("Selected Card")
          Spacer()
        }
      }

      DSCard(style: .interactive, action: { print("Tapped") }) {
        HStack {
          Text("Interactive Card - Tap me")
          Spacer()
          Image(systemName: "chevron.right")
            .foregroundColor(.secondary)
        }
      }

      // Using modifier
      HStack {
        Text("Using .dsCard() modifier")
        Spacer()
      }
      .dsCard(style: .elevated)
    }
    .padding()
  }
}
