//
//  View+DesignSystem.swift
//  ios
//
//  Design System View Extensions - Convenience modifiers for common patterns
//

import SwiftUI

// MARK: - Conditional Modifiers

extension View {
  /// Apply a modifier conditionally.
  @ViewBuilder
  func `if`<Content: View>(_ condition: Bool, transform: (Self) -> Content) -> some View {
    if condition {
      transform(self)
    } else {
      self
    }
  }

  /// Apply a modifier conditionally with an else clause.
  @ViewBuilder
  func `if`<TrueContent: View, FalseContent: View>(
    _ condition: Bool,
    then trueTransform: (Self) -> TrueContent,
    else falseTransform: (Self) -> FalseContent
  ) -> some View {
    if condition {
      trueTransform(self)
    } else {
      falseTransform(self)
    }
  }
}

// MARK: - Loading Overlay

extension View {
  /// Add a loading overlay to any view.
  func dsLoading(_ isLoading: Bool, message: String? = nil) -> some View {
    self.overlay {
      if isLoading {
        ZStack {
          Color(.systemBackground).opacity(0.8)
          DSLoadingView(message)
        }
      }
    }
    .disabled(isLoading)
  }
}

// MARK: - Semantic Styling

extension View {
  /// Apply primary action styling.
  func dsPrimaryStyle() -> some View {
    self
      .font(DSTheme.Typography.button)
      .foregroundColor(DSTheme.Colors.textOnPrimary)
      .padding(.vertical, DSTheme.Spacing.md)
      .padding(.horizontal, DSTheme.Spacing.lg)
      .background(DSTheme.Colors.primary)
      .cornerRadius(DSTheme.Radius.md)
  }

  /// Apply secondary styling.
  func dsSecondaryStyle() -> some View {
    self
      .font(DSTheme.Typography.button)
      .foregroundColor(DSTheme.Colors.primary)
      .padding(.vertical, DSTheme.Spacing.md)
      .padding(.horizontal, DSTheme.Spacing.lg)
      .background(DSTheme.Colors.primary.opacity(0.1))
      .cornerRadius(DSTheme.Radius.md)
  }
}

// MARK: - Full Width

extension View {
  /// Make a view full width.
  func dsFullWidth(alignment: Alignment = .center) -> some View {
    self.frame(maxWidth: .infinity, alignment: alignment)
  }
}

// MARK: - Standard Padding

extension View {
  /// Apply standard screen padding.
  func dsScreenPadding() -> some View {
    self.padding(DSTheme.Spacing.lg)
  }

  /// Apply standard section spacing.
  func dsSectionSpacing() -> some View {
    self.padding(.vertical, DSTheme.Spacing.xl)
  }
}

// MARK: - Divider with Spacing

extension View {
  /// Add a divider below the view with standard spacing.
  func dsDivided(spacing: CGFloat = DSTheme.Spacing.md) -> some View {
    VStack(spacing: spacing) {
      self
      Divider()
    }
  }
}

// MARK: - Hide/Show

extension View {
  /// Conditionally hide the view.
  @ViewBuilder
  func dsHidden(_ hidden: Bool) -> some View {
    if hidden {
      self.hidden()
    } else {
      self
    }
  }
}

// MARK: - Redacted Loading

extension View {
  /// Apply redacted placeholder styling when loading.
  func dsRedactedWhenLoading(_ isLoading: Bool) -> some View {
    self.redacted(reason: isLoading ? .placeholder : [])
  }
}

// MARK: - Error Border

extension View {
  /// Apply error border when there's an error.
  func dsErrorBorder(_ hasError: Bool) -> some View {
    self.overlay(
      RoundedRectangle(cornerRadius: DSTheme.Radius.sm)
        .stroke(hasError ? DSTheme.Colors.error : .clear, lineWidth: 1)
    )
  }
}

// MARK: - Navigation Bar Styling

extension View {
  /// Apply standard navigation bar styling.
  func dsNavigationBar(title: String, displayMode: NavigationBarItem.TitleDisplayMode = .large)
    -> some View
  {
    self
      .navigationTitle(title)
      .navigationBarTitleDisplayMode(displayMode)
  }
}

// MARK: - Async Button Action

extension View {
  /// Wrap an async action for use in buttons.
  func dsAsyncAction(_ action: @escaping () async -> Void) -> () -> Void {
    return {
      Task {
        await action()
      }
    }
  }
}
