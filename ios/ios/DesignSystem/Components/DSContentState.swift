//
//  DSContentState.swift
//  ios
//
//  Design System Content State Components - Loading, Error, and Content wrapper
//

import SwiftUI

// MARK: - DSLoadingView

/// A consistent loading indicator view.
///
/// Usage:
/// ```swift
/// DSLoadingView()
/// DSLoadingView("Loading...")
/// DSLoadingView("Please wait", size: .large)
/// ```
struct DSLoadingView: View {
  let message: String?
  let size: Size

  enum Size {
    case small
    case medium
    case large

    var progressScale: CGFloat {
      switch self {
      case .small:
        return 0.8
      case .medium:
        return 1.0
      case .large:
        return 1.3
      }
    }

    var textFont: Font {
      switch self {
      case .small:
        return DSTheme.Typography.caption
      case .medium:
        return DSTheme.Typography.body
      case .large:
        return DSTheme.Typography.label
      }
    }

    var spacing: CGFloat {
      switch self {
      case .small:
        return DSTheme.Spacing.sm
      case .medium:
        return DSTheme.Spacing.md
      case .large:
        return DSTheme.Spacing.lg
      }
    }
  }

  init(_ message: String? = "Loading...", size: Size = .medium) {
    self.message = message
    self.size = size
  }

  var body: some View {
    VStack(spacing: size.spacing) {
      ProgressView()
        .scaleEffect(size.progressScale)

      if let message = message {
        Text(message)
          .font(size.textFont)
          .foregroundColor(DSTheme.Colors.textSecondary)
      }
    }
    .frame(maxWidth: .infinity, maxHeight: .infinity)
  }
}

// MARK: - DSErrorView

/// A consistent error display view with optional retry action.
///
/// Usage:
/// ```swift
/// DSErrorView("Something went wrong")
/// DSErrorView("Failed to load", onRetry: { await reload() })
/// DSErrorView("Network error", icon: "wifi.slash", onRetry: { await retry() })
/// ```
struct DSErrorView: View {
  let message: String
  let title: String
  let icon: String
  let iconColor: Color
  let retryTitle: String
  let onRetry: (() async -> Void)?
  let showEnvironmentSelector: Bool

  #if DEBUG
    @State private var showEnvironmentPicker = false
    @State private var selectedEnvironment: AppEnvironment = APIConfiguration.currentEnvironment
  #endif

  init(
    _ message: String,
    title: String = "Error",
    icon: String = "exclamationmark.triangle",
    iconColor: Color = DSTheme.Colors.warning,
    retryTitle: String = "Retry",
    onRetry: (() async -> Void)? = nil,
    showEnvironmentSelector: Bool = false
  ) {
    self.message = message
    self.title = title
    self.icon = icon
    self.iconColor = iconColor
    self.retryTitle = retryTitle
    self.onRetry = onRetry
    self.showEnvironmentSelector = showEnvironmentSelector
  }

  var body: some View {
    VStack(spacing: DSTheme.Spacing.lg) {
      Image(systemName: icon)
        .font(.system(size: DSTheme.IconSize.xxl))
        .foregroundColor(iconColor)

      VStack(spacing: DSTheme.Spacing.sm) {
        Text(title)
          .font(DSTheme.Typography.title3)
          .foregroundColor(DSTheme.Colors.textPrimary)

        Text(message)
          .font(DSTheme.Typography.body)
          .foregroundColor(DSTheme.Colors.textSecondary)
          .multilineTextAlignment(.center)
          .padding(.horizontal, DSTheme.Spacing.xl)
      }

      if let onRetry = onRetry {
        DSButton(retryTitle, icon: "arrow.clockwise", style: .primary) {
          Task {
            await onRetry()
          }
        }
      }

      #if DEBUG
        if showEnvironmentSelector {
          environmentButton
        }
      #endif
    }
    .frame(maxWidth: .infinity, maxHeight: .infinity)
    #if DEBUG
      .sheet(isPresented: $showEnvironmentPicker) {
        EnvironmentPickerSheet(
          selectedEnvironment: $selectedEnvironment,
          onDismiss: { showEnvironmentPicker = false }
        )
        .presentationDetents([.medium])
      }
    #endif
  }

  #if DEBUG
    private var environmentButton: some View {
      Button {
        showEnvironmentPicker = true
      } label: {
        HStack(spacing: DSTheme.Spacing.xs) {
          Image(systemName: selectedEnvironment.iconName)
            .font(.system(size: 12))
          Text(selectedEnvironment.displayName)
            .font(DSTheme.Typography.caption)
          Image(systemName: "chevron.up")
            .font(.system(size: 10, weight: .semibold))
        }
        .foregroundColor(DSTheme.Colors.textSecondary)
        .padding(.horizontal, DSTheme.Spacing.md)
        .padding(.vertical, DSTheme.Spacing.sm)
        .background(DSTheme.Colors.cardBackground)
        .cornerRadius(DSTheme.Radius.full)
      }
      .padding(.top, DSTheme.Spacing.sm)
    }
  #endif
}

// MARK: - DSContentState

/// A wrapper that handles loading, error, and content states.
///
/// Usage:
/// ```swift
/// DSContentState(
///   isLoading: viewModel.isLoading,
///   error: viewModel.errorMessage,
///   onRetry: { await viewModel.loadData() }
/// ) {
///   ScrollView {
///     // Your content here
///   }
/// }
///
/// // With custom loading message
/// DSContentState(
///   isLoading: viewModel.isLoading,
///   loadingMessage: "Fetching notifications...",
///   error: viewModel.errorMessage,
///   onRetry: { await viewModel.loadData() }
/// ) {
///   NotificationList(items: viewModel.items)
/// }
/// ```
struct DSContentState<Content: View>: View {
  let isLoading: Bool
  let loadingMessage: String?
  let error: String?
  let errorTitle: String
  let errorIcon: String
  let errorIconColor: Color
  let onRetry: (() async -> Void)?
  let showEnvironmentSelector: Bool
  @ViewBuilder let content: () -> Content

  init(
    isLoading: Bool,
    loadingMessage: String? = "Loading...",
    error: String? = nil,
    errorTitle: String = "Error",
    errorIcon: String = "exclamationmark.triangle",
    errorIconColor: Color = DSTheme.Colors.warning,
    onRetry: (() async -> Void)? = nil,
    showEnvironmentSelector: Bool = false,
    @ViewBuilder content: @escaping () -> Content
  ) {
    self.isLoading = isLoading
    self.loadingMessage = loadingMessage
    self.error = error
    self.errorTitle = errorTitle
    self.errorIcon = errorIcon
    self.errorIconColor = errorIconColor
    self.onRetry = onRetry
    self.showEnvironmentSelector = showEnvironmentSelector
    self.content = content
  }

  var body: some View {
    Group {
      if isLoading {
        DSLoadingView(loadingMessage)
      } else if let error = error {
        DSErrorView(
          error,
          title: errorTitle,
          icon: errorIcon,
          iconColor: errorIconColor,
          onRetry: onRetry,
          showEnvironmentSelector: showEnvironmentSelector
        )
      } else {
        content()
      }
    }
  }
}

// MARK: - Initializing State View

/// A view for the initializing/setup state before the view model is ready.
///
/// Usage:
/// ```swift
/// if let viewModel = viewModel {
///   // ... normal content
/// } else {
///   DSInitializingView()
/// }
/// ```
struct DSInitializingView: View {
  let message: String

  init(_ message: String = "Initializing...") {
    self.message = message
  }

  var body: some View {
    DSLoadingView(message)
  }
}

// MARK: - View Extension for Optional ViewModel Pattern

extension View {
  /// Wraps content in a DSContentState, handling the common optional viewModel pattern.
  ///
  /// Usage:
  /// ```swift
  /// someContent
  ///   .dsContentState(
  ///     viewModel: viewModel,
  ///     isLoading: \.isLoading,
  ///     error: \.errorMessage,
  ///     onRetry: { await $0.loadData() }
  ///   )
  /// ```
  @ViewBuilder
  func dsHandleState<VM>(
    viewModel: VM?,
    isLoading: Bool,
    error: String?,
    initializingMessage: String = "Initializing...",
    loadingMessage: String = "Loading...",
    onRetry: (() async -> Void)? = nil
  ) -> some View {
    if viewModel != nil {
      DSContentState(
        isLoading: isLoading,
        loadingMessage: loadingMessage,
        error: error,
        onRetry: onRetry
      ) {
        self
      }
    } else {
      DSInitializingView(initializingMessage)
    }
  }
}

// MARK: - Preview

#Preview("Content States") {
  TabView {
    // Loading
    DSLoadingView("Loading...")
      .tabItem { Label("Loading", systemImage: "hourglass") }

    // Error without retry
    DSErrorView("Unable to connect to the server. Please check your internet connection.")
      .tabItem { Label("Error", systemImage: "exclamationmark.triangle") }

    // Error with retry
    DSErrorView("Failed to load data", onRetry: { print("Retry tapped") })
      .tabItem { Label("Retry", systemImage: "arrow.clockwise") }

    // Content State wrapper - Loading
    DSContentState(isLoading: true, error: nil) {
      Text("Content")
    }
    .tabItem { Label("Wrapper", systemImage: "square.stack") }

    // Initializing
    DSInitializingView()
      .tabItem { Label("Init", systemImage: "gear") }
  }
}
