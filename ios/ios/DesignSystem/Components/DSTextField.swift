//
//  DSTextField.swift
//  ios
//
//  Design System TextField Component - Consistent text input styling
//

import SwiftUI

// MARK: - Text Field Type

enum DSTextFieldType {
  case text
  case username
  case email
  case password
  case phone
  case number
  case multiline

  var keyboardType: UIKeyboardType {
    switch self {
    case .email:
      return .emailAddress
    case .phone:
      return .phonePad
    case .number:
      return .decimalPad
    default:
      return .default
    }
  }

  var textContentType: UITextContentType? {
    switch self {
    case .email:
      return .emailAddress
    case .username:
      return .username
    case .password:
      return .password
    case .phone:
      return .telephoneNumber
    default:
      return nil
    }
  }

  var autocapitalization: TextInputAutocapitalization {
    switch self {
    case .email, .password, .username:
      return .never
    default:
      return .sentences
    }
  }

  var disableAutocorrection: Bool {
    switch self {
    case .email, .password, .username:
      return true
    default:
      return false
    }
  }
}

// MARK: - Text Field Style

enum DSTextFieldStyle {
  case standard
  case outlined
  case filled

  var backgroundColor: Color {
    switch self {
    case .standard:
      return .clear
    case .outlined:
      return .clear
    case .filled:
      return DSTheme.Colors.cardBackground
    }
  }

  var borderColor: Color {
    switch self {
    case .standard:
      return .clear
    case .outlined:
      return DSTheme.Colors.border
    case .filled:
      return .clear
    }
  }
}

// MARK: - DSTextField Component

/// A styled text field with consistent appearance and behavior.
///
/// Usage:
/// ```swift
/// DSTextField("Username", text: $username, type: .username)
/// DSTextField("Email", text: $email, type: .email)
/// DSTextField("Password", text: $password, type: .password)
/// DSTextField("Phone", text: $phone, type: .phone)
/// DSTextField("Notes", text: $notes, type: .multiline)
///
/// // With label
/// DSTextField("Email Address", text: $email, type: .email, label: "Contact Email")
///
/// // With error
/// DSTextField("Email", text: $email, type: .email, error: emailError)
///
/// // Disabled
/// DSTextField("Email", text: $email, isDisabled: true)
/// ```
struct DSTextField: View {
  let placeholder: String
  @Binding var text: String
  let type: DSTextFieldType
  let style: DSTextFieldStyle
  let label: String?
  let helperText: String?
  let error: String?
  let icon: String?
  let isDisabled: Bool
  let lineLimit: ClosedRange<Int>?

  @FocusState private var isFocused: Bool

  init(
    _ placeholder: String,
    text: Binding<String>,
    type: DSTextFieldType = .text,
    style: DSTextFieldStyle = .standard,
    label: String? = nil,
    helperText: String? = nil,
    error: String? = nil,
    icon: String? = nil,
    isDisabled: Bool = false,
    lineLimit: ClosedRange<Int>? = nil
  ) {
    self.placeholder = placeholder
    self._text = text
    self.type = type
    self.style = style
    self.label = label
    self.helperText = helperText
    self.error = error
    self.icon = icon
    self.isDisabled = isDisabled
    self.lineLimit = lineLimit ?? (type == .multiline ? 3...6 : 1...1)
  }

  var body: some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.xs) {
      // Label
      if let label = label {
        Text(label)
          .font(DSTheme.Typography.labelSmall)
          .foregroundColor(DSTheme.Colors.textSecondary)
      }

      // Input field
      HStack(spacing: DSTheme.Spacing.sm) {
        if let icon = icon {
          Image(systemName: icon)
            .font(.system(size: DSTheme.IconSize.sm))
            .foregroundColor(DSTheme.Colors.textSecondary)
        }

        inputField
      }
      .padding(.horizontal, DSTheme.Spacing.md)
      .padding(.vertical, DSTheme.Spacing.md)
      .background(effectiveBackgroundColor)
      .cornerRadius(DSTheme.Radius.sm)
      .overlay(
        RoundedRectangle(cornerRadius: DSTheme.Radius.sm)
          .stroke(effectiveBorderColor, lineWidth: 1)
      )

      // Helper text or error
      if let error = error, !error.isEmpty {
        Text(error)
          .font(DSTheme.Typography.caption)
          .foregroundColor(DSTheme.Colors.error)
      } else if let helperText = helperText {
        Text(helperText)
          .font(DSTheme.Typography.caption)
          .foregroundColor(DSTheme.Colors.textTertiary)
      }
    }
  }

  @ViewBuilder
  private var inputField: some View {
    Group {
      switch type {
      case .password:
        SecureField(placeholder, text: $text)
      case .multiline:
        TextField(placeholder, text: $text, axis: .vertical)
          .lineLimit(lineLimit ?? 3...6)
      default:
        TextField(placeholder, text: $text)
      }
    }
    .font(DSTheme.Typography.body)
    .foregroundColor(isDisabled ? DSTheme.Colors.textTertiary : DSTheme.Colors.textPrimary)
    .keyboardType(type.keyboardType)
    .textInputAutocapitalization(type.autocapitalization)
    .autocorrectionDisabled(type.disableAutocorrection)
    .disabled(isDisabled)
    .focused($isFocused)
  }

  private var effectiveBackgroundColor: Color {
    if isDisabled {
      return DSTheme.Colors.cardBackground.opacity(0.5)
    }
    return style.backgroundColor
  }

  private var effectiveBorderColor: Color {
    if error != nil && !error!.isEmpty {
      return DSTheme.Colors.error
    }
    if isFocused {
      return DSTheme.Colors.primary
    }
    if isDisabled {
      return DSTheme.Colors.border.opacity(0.5)
    }
    switch style {
    case .standard:
      return DSTheme.Colors.border
    case .outlined:
      return DSTheme.Colors.border
    case .filled:
      return .clear
    }
  }
}

// MARK: - Labeled Field Group

/// A group of text fields with a shared label.
///
/// Usage:
/// ```swift
/// DSFieldGroup("Address") {
///   DSTextField("Street", text: $street)
///   DSTextField("City", text: $city)
///   DSTextField("Zip", text: $zip, type: .number)
/// }
/// ```
struct DSFieldGroup<Content: View>: View {
  let label: String
  let spacing: CGFloat
  @ViewBuilder let content: () -> Content

  init(
    _ label: String,
    spacing: CGFloat = DSTheme.Spacing.md,
    @ViewBuilder content: @escaping () -> Content
  ) {
    self.label = label
    self.spacing = spacing
    self.content = content
  }

  var body: some View {
    VStack(alignment: .leading, spacing: DSTheme.Spacing.sm) {
      Text(label)
        .font(DSTheme.Typography.label)
        .foregroundColor(DSTheme.Colors.textSecondary)

      VStack(spacing: spacing) {
        content()
      }
    }
  }
}

// MARK: - Preview

#Preview("Text Field Styles") {
  ScrollView {
    VStack(spacing: DSTheme.Spacing.xl) {
      Group {
        Text("Standard").font(.caption).foregroundColor(.secondary)
        DSTextField("Email", text: .constant(""), type: .email)

        Text("With Label").font(.caption).foregroundColor(.secondary)
        DSTextField("Enter password", text: .constant(""), type: .password, label: "Password")

        Text("With Icon").font(.caption).foregroundColor(.secondary)
        DSTextField("Search", text: .constant(""), icon: "magnifyingglass")

        Text("With Error").font(.caption).foregroundColor(.secondary)
        DSTextField(
          "Email", text: .constant("invalid"), type: .email, error: "Please enter a valid email")

        Text("With Helper").font(.caption).foregroundColor(.secondary)
        DSTextField(
          "Username", text: .constant(""), type: .username, helperText: "Choose a unique username")

        Text("Disabled").font(.caption).foregroundColor(.secondary)
        DSTextField("Email", text: .constant("disabled@example.com"), isDisabled: true)
      }

      Divider()

      Group {
        Text("Multiline").font(.caption).foregroundColor(.secondary)
        DSTextField("Enter notes...", text: .constant(""), type: .multiline)

        Text("Phone").font(.caption).foregroundColor(.secondary)
        DSTextField("Phone number", text: .constant(""), type: .phone)

        Text("Number").font(.caption).foregroundColor(.secondary)
        DSTextField("Amount", text: .constant(""), type: .number)
      }

      Divider()

      DSFieldGroup("Contact Information") {
        DSTextField("Email", text: .constant(""), type: .email)
        DSTextField("Phone", text: .constant(""), type: .phone)
      }
    }
    .padding()
  }
}
