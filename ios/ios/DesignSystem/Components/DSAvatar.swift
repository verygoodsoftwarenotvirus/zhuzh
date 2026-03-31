//
//  DSAvatar.swift
//  ios
//
//  Design System Avatar Component - User avatars with initials fallback
//

import SwiftUI

// MARK: - Avatar Size

enum DSAvatarSize {
  case xs
  case sm
  case md
  case lg
  case xl
  case xxl
  case custom(CGFloat)

  var dimension: CGFloat {
    switch self {
    case .xs:
      return DSTheme.AvatarSize.sm - 8  // 24
    case .sm:
      return DSTheme.AvatarSize.sm  // 32
    case .md:
      return DSTheme.AvatarSize.md  // 40
    case .lg:
      return DSTheme.AvatarSize.lg  // 50
    case .xl:
      return DSTheme.AvatarSize.xl  // 64
    case .xxl:
      return DSTheme.AvatarSize.xxl  // 80
    case .custom(let size):
      return size
    }
  }

  var font: Font {
    switch self {
    case .xs:
      return .system(size: 10, weight: .semibold)
    case .sm:
      return .system(size: 12, weight: .semibold)
    case .md:
      return .system(size: 14, weight: .semibold)
    case .lg:
      return .system(size: 18, weight: .semibold)
    case .xl:
      return .system(size: 22, weight: .semibold)
    case .xxl:
      return .system(size: 28, weight: .semibold)
    case .custom(let size):
      return .system(size: size * 0.35, weight: .semibold)
    }
  }
}

// MARK: - Avatar Style

enum DSAvatarStyle {
  case circle
  case rounded
  case square
}

// MARK: - DSAvatar Component

/// A user avatar with initials fallback and optional image.
///
/// Usage:
/// ```swift
/// // With name (generates initials automatically)
/// DSAvatar(name: "John Doe")
/// DSAvatar(name: "Jane Smith", size: .lg)
///
/// // With explicit initials
/// DSAvatar(initials: "JD")
///
/// // With image URL (when supported)
/// DSAvatar(name: "John Doe", imageURL: user.avatarURL)
///
/// // Custom color
/// DSAvatar(name: "John Doe", backgroundColor: .purple)
/// ```
struct DSAvatar: View {
  let initials: String
  let size: DSAvatarSize
  let style: DSAvatarStyle
  let backgroundColor: Color
  let foregroundColor: Color
  let imageURL: URL?

  init(
    name: String,
    size: DSAvatarSize = .md,
    style: DSAvatarStyle = .circle,
    backgroundColor: Color? = nil,
    foregroundColor: Color = DSTheme.Colors.textSecondary,
    imageURL: URL? = nil
  ) {
    self.initials = Self.generateInitials(from: name)
    self.size = size
    self.style = style
    self.backgroundColor = backgroundColor ?? Self.colorForName(name)
    self.foregroundColor = foregroundColor
    self.imageURL = imageURL
  }

  init(
    initials: String,
    size: DSAvatarSize = .md,
    style: DSAvatarStyle = .circle,
    backgroundColor: Color? = nil,
    foregroundColor: Color = DSTheme.Colors.textSecondary,
    imageURL: URL? = nil
  ) {
    self.initials = String(initials.prefix(2)).uppercased()
    self.size = size
    self.style = style
    self.backgroundColor = backgroundColor ?? .gray.opacity(0.3)
    self.foregroundColor = foregroundColor
    self.imageURL = imageURL
  }

  var body: some View {
    Group {
      if let url = imageURL {
        AsyncImage(url: url) { phase in
          switch phase {
          case .success(let image):
            image
              .resizable()
              .aspectRatio(contentMode: .fill)
          case .failure, .empty:
            initialsOverlay
          @unknown default:
            initialsOverlay
          }
        }
        .frame(width: size.dimension, height: size.dimension)
        .clipShape(avatarClipShape)
      } else {
        initialsOverlay
          .frame(width: size.dimension, height: size.dimension)
      }
    }
  }

  @ViewBuilder
  private var initialsOverlay: some View {
    Group {
      switch style {
      case .circle:
        Circle()
          .fill(backgroundColor)
          .overlay {
            Text(initials)
              .font(size.font)
              .foregroundColor(foregroundColor)
          }
      case .rounded:
        RoundedRectangle(cornerRadius: size.dimension * 0.2)
          .fill(backgroundColor)
          .overlay {
            Text(initials)
              .font(size.font)
              .foregroundColor(foregroundColor)
          }
      case .square:
        RoundedRectangle(cornerRadius: DSTheme.Radius.xs)
          .fill(backgroundColor)
          .overlay {
            Text(initials)
              .font(size.font)
              .foregroundColor(foregroundColor)
          }
      }
    }
  }

  private var avatarClipShape: AnyShape {
    switch style {
    case .circle:
      AnyShape(Circle())
    case .rounded:
      AnyShape(RoundedRectangle(cornerRadius: size.dimension * 0.2))
    case .square:
      AnyShape(RoundedRectangle(cornerRadius: DSTheme.Radius.xs))
    }
  }

  /// Generate initials from a name.
  static func generateInitials(from name: String) -> String {
    let words = name.split(separator: " ")
    if words.isEmpty {
      return "?"
    }
    if words.count == 1 {
      return String(words[0].prefix(1)).uppercased()
    }
    let firstInitial = words.first.map { String($0.prefix(1)) } ?? ""
    let lastInitial = words.last.map { String($0.prefix(1)) } ?? ""
    return (firstInitial + lastInitial).uppercased()
  }

  /// Generate a consistent color based on the name.
  static func colorForName(_ name: String) -> Color {
    let colors: [Color] = [
      .blue.opacity(0.3),
      .green.opacity(0.3),
      .orange.opacity(0.3),
      .purple.opacity(0.3),
      .pink.opacity(0.3),
      .teal.opacity(0.3),
      .indigo.opacity(0.3),
      .mint.opacity(0.3),
    ]
    let hash = abs(name.hashValue)
    return colors[hash % colors.count]
  }
}

// MARK: - Simpler Avatar View

/// A simpler avatar implementation that avoids complex generics
struct DSAvatarView: View {
  let initials: String
  let size: DSAvatarSize
  let backgroundColor: Color
  let foregroundColor: Color

  init(
    name: String,
    size: DSAvatarSize = .md,
    backgroundColor: Color? = nil,
    foregroundColor: Color = DSTheme.Colors.textSecondary
  ) {
    self.initials = DSAvatar.generateInitials(from: name)
    self.size = size
    self.backgroundColor = backgroundColor ?? DSAvatar.colorForName(name)
    self.foregroundColor = foregroundColor
  }

  init(
    initials: String,
    size: DSAvatarSize = .md,
    backgroundColor: Color = .gray.opacity(0.3),
    foregroundColor: Color = DSTheme.Colors.textSecondary
  ) {
    self.initials = String(initials.prefix(2)).uppercased()
    self.size = size
    self.backgroundColor = backgroundColor
    self.foregroundColor = foregroundColor
  }

  var body: some View {
    Circle()
      .fill(backgroundColor)
      .frame(width: size.dimension, height: size.dimension)
      .overlay {
        Text(initials)
          .font(size.font)
          .foregroundColor(foregroundColor)
      }
  }
}

// MARK: - Avatar Group

/// A group of overlapping avatars.
///
/// Usage:
/// ```swift
/// DSAvatarGroup(names: ["John", "Jane", "Bob"], maxVisible: 3)
/// ```
struct DSAvatarGroup: View {
  let names: [String]
  let size: DSAvatarSize
  let maxVisible: Int
  let overlapOffset: CGFloat

  init(
    names: [String],
    size: DSAvatarSize = .sm,
    maxVisible: Int = 4,
    overlapOffset: CGFloat? = nil
  ) {
    self.names = names
    self.size = size
    self.maxVisible = maxVisible
    self.overlapOffset = overlapOffset ?? (size.dimension * 0.6)
  }

  var body: some View {
    HStack(spacing: -overlapOffset + size.dimension) {
      ForEach(Array(visibleNames.enumerated()), id: \.offset) { index, name in
        DSAvatarView(name: name, size: size)
          .overlay(
            Circle()
              .stroke(Color(.systemBackground), lineWidth: 2)
          )
          .zIndex(Double(visibleNames.count - index))
      }

      if remainingCount > 0 {
        Circle()
          .fill(DSTheme.Colors.cardBackground)
          .frame(width: size.dimension, height: size.dimension)
          .overlay {
            Text("+\(remainingCount)")
              .font(size.font)
              .foregroundColor(DSTheme.Colors.textSecondary)
          }
          .overlay(
            Circle()
              .stroke(Color(.systemBackground), lineWidth: 2)
          )
      }
    }
  }

  private var visibleNames: [String] {
    Array(names.prefix(maxVisible))
  }

  private var remainingCount: Int {
    max(0, names.count - maxVisible)
  }
}

// MARK: - Preview

#Preview("Avatars") {
  ScrollView {
    VStack(spacing: DSTheme.Spacing.xl) {
      Group {
        Text("Sizes").font(.caption).foregroundColor(.secondary)
        HStack(spacing: DSTheme.Spacing.md) {
          DSAvatarView(name: "John Doe", size: .xs)
          DSAvatarView(name: "John Doe", size: .sm)
          DSAvatarView(name: "John Doe", size: .md)
          DSAvatarView(name: "John Doe", size: .lg)
          DSAvatarView(name: "John Doe", size: .xl)
        }
      }

      Divider()

      Group {
        Text("Different Names").font(.caption).foregroundColor(.secondary)
        HStack(spacing: DSTheme.Spacing.md) {
          DSAvatarView(name: "Alice")
          DSAvatarView(name: "Bob Smith")
          DSAvatarView(name: "Charlie Brown")
          DSAvatarView(name: "Diana Prince")
        }
      }

      Divider()

      Group {
        Text("Custom Initials").font(.caption).foregroundColor(.secondary)
        HStack(spacing: DSTheme.Spacing.md) {
          DSAvatarView(initials: "AB", size: .lg)
          DSAvatarView(initials: "XY", size: .lg, backgroundColor: .purple.opacity(0.3))
          DSAvatarView(initials: "?", size: .lg, backgroundColor: .red.opacity(0.3))
        }
      }

      Divider()

      Group {
        Text("Avatar Group").font(.caption).foregroundColor(.secondary)
        DSAvatarGroup(
          names: ["Alice", "Bob", "Charlie", "Diana", "Eve", "Frank"],
          size: .md,
          maxVisible: 4
        )
      }

      Group {
        Text("Small Avatar Group").font(.caption).foregroundColor(.secondary)
        DSAvatarGroup(
          names: ["Alice", "Bob", "Charlie"],
          size: .sm,
          maxVisible: 5
        )
      }
    }
    .padding()
  }
}
