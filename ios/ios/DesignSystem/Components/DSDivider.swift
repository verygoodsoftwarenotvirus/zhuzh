//
//  DSDivider.swift
//  ios
//
//  Design System Divider - Consistent divider styling with standard and soft variants
//

import SwiftUI

// MARK: - DSDivider Style

enum DSDividerStyle {
  /// System divider for standard separation
  case standard
  /// Softer divider with reduced opacity and optional horizontal inset
  case soft
}

// MARK: - DSDivider

/// A theme-consistent divider with standard and soft variants.
///
/// Usage:
/// ```swift
/// DSDivider()
/// DSDivider(style: .soft)
/// DSDivider(style: .soft, horizontalInset: DSTheme.Spacing.xl * 2)
/// ```
struct DSDivider: View {
  var style: DSDividerStyle = .standard
  var horizontalInset: CGFloat = 0

  var body: some View {
    switch style {
    case .standard:
      Divider()
        .padding(.leading, horizontalInset)
    case .soft:
      Rectangle()
        .fill(DSTheme.Colors.border.opacity(0.5))
        .frame(height: 1)
        .frame(maxWidth: .infinity)
        .padding(.horizontal, horizontalInset > 0 ? horizontalInset : DSTheme.Spacing.xl * 2)
    }
  }
}

// MARK: - DSSoftDivider

/// A soft divider with vertical spacing, for use between content sections.
/// Replaces the softSeparator pattern (Rectangle with opacity, horizontal padding, vertical spacing).
///
/// Usage:
/// ```swift
/// DSSoftDivider()
/// DSSoftDivider(verticalSpacing: DSTheme.Spacing.lg)
/// ```
struct DSSoftDivider: View {
  var verticalSpacing: CGFloat = DSTheme.Spacing.md
  var horizontalInset: CGFloat = DSTheme.Spacing.xl * 2

  var body: some View {
    VStack(spacing: 0) {
      Spacer()
        .frame(height: verticalSpacing)
      DSDivider(style: .soft, horizontalInset: horizontalInset)
      Spacer()
        .frame(height: verticalSpacing)
    }
  }
}

// MARK: - Preview

#Preview("DSDivider") {
  VStack(spacing: DSTheme.Spacing.lg) {
    Text("Above")
    DSDivider()
    Text("Below")

    Text("Above soft")
    DSDivider(style: .soft)
    Text("Below soft")

    Text("Above DSSoftDivider")
    DSSoftDivider()
    Text("Below DSSoftDivider")
  }
  .padding()
}
