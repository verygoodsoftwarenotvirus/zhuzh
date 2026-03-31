//
//  AppTheme.swift
//  ios
//
//  Design System Theme - Central configuration for colors, spacing, typography, and radii
//

import SwiftUI

// MARK: - Design System Theme

/// Central theme configuration for the app's design system.
/// Use these tokens throughout the app for consistent styling.
enum DSTheme {

  // MARK: - Colors

  enum Colors {
    // Brand colors
    static let primary = Color.blue
    static let secondary = Color.purple
    static let tertiary = Color.orange

    // Semantic colors
    static let success = Color.green
    static let warning = Color.orange
    static let error = Color.red
    static let info = Color.blue

    // Background colors
    static let cardBackground = Color(.systemGray6)
    static let cardBackgroundSelected = Color.blue.opacity(0.1)
    static let cardBackgroundElevated = Color(.systemGray5)

    // Text colors
    static let textPrimary = Color.primary
    static let textSecondary = Color.secondary
    static let textTertiary = Color(.tertiaryLabel)
    static let textOnPrimary = Color.white

    // Border colors
    static let border = Color(.separator)
    static let borderSelected = Color.blue
    static let borderSubtle = Color(.systemGray4)

    // Status colors (for badges, indicators)
    static let statusPending = Color.orange
    static let statusAccepted = Color.green
    static let statusRejected = Color.red
    static let statusCancelled = Color.gray
    static let statusUnknown = Color.secondary
  }

  // MARK: - Spacing

  enum Spacing {
    static let xxs: CGFloat = 2
    static let xs: CGFloat = 4
    static let sm: CGFloat = 8
    static let md: CGFloat = 12
    static let lg: CGFloat = 16
    static let xl: CGFloat = 24
    static let xxl: CGFloat = 32
    static let xxxl: CGFloat = 48
  }

  // MARK: - Corner Radius

  enum Radius {
    static let xs: CGFloat = 4
    static let sm: CGFloat = 8
    static let md: CGFloat = 10
    static let lg: CGFloat = 16
    static let xl: CGFloat = 20
    static let full: CGFloat = 9999  // For circular elements
  }

  // MARK: - Typography

  enum Typography {
    // Titles
    static let largeTitle = Font.largeTitle.weight(.bold)
    static let title1 = Font.title.weight(.bold)
    static let title2 = Font.title2.weight(.bold)
    static let title3 = Font.title3.weight(.semibold)

    // Body text
    static let bodyLarge = Font.body
    static let body = Font.subheadline
    static let bodySmall = Font.footnote

    // Labels
    static let label = Font.subheadline.weight(.medium)
    static let labelSmall = Font.caption.weight(.medium)
    static let caption = Font.caption
    static let captionSmall = Font.caption2

    // Special
    static let button = Font.body.weight(.semibold)
    static let buttonSmall = Font.subheadline.weight(.semibold)
  }

  // MARK: - Shadows

  enum Shadows {
    static let sm = ShadowStyle(color: .black.opacity(0.08), radius: 4, x: 0, y: 2)
    static let md = ShadowStyle(color: .black.opacity(0.1), radius: 8, x: 0, y: 4)
    static let lg = ShadowStyle(color: .black.opacity(0.12), radius: 16, x: 0, y: 8)
  }

  // MARK: - Animation

  enum Animation {
    static let fast = SwiftUI.Animation.easeInOut(duration: 0.15)
    static let normal = SwiftUI.Animation.easeInOut(duration: 0.25)
    static let slow = SwiftUI.Animation.easeInOut(duration: 0.4)
    static let spring = SwiftUI.Animation.spring(response: 0.3, dampingFraction: 0.7)
  }

  // MARK: - Icon Sizes

  enum IconSize {
    static let xs: CGFloat = 12
    static let sm: CGFloat = 16
    static let md: CGFloat = 20
    static let lg: CGFloat = 24
    static let xl: CGFloat = 32
    static let xxl: CGFloat = 48
  }

  // MARK: - Avatar Sizes

  enum AvatarSize {
    static let sm: CGFloat = 32
    static let md: CGFloat = 40
    static let lg: CGFloat = 50
    static let xl: CGFloat = 64
    static let xxl: CGFloat = 80
  }
}

// MARK: - Shadow Style Helper

struct ShadowStyle {
  let color: Color
  let radius: CGFloat
  let x: CGFloat
  let y: CGFloat
}

// MARK: - View Extension for Shadows

extension View {
  func dsShadow(_ style: ShadowStyle) -> some View {
    self.shadow(color: style.color, radius: style.radius, x: style.x, y: style.y)
  }
}
