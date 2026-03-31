//
//  DesignSystemTests.swift
//  iosTests
//
//  Tests for the Design System components and theme tokens
//

import Foundation
import SwiftUI
@testable import ios
import Testing

// MARK: - Theme Token Tests

struct ThemeSpacingTests {
  @Test("Spacing tokens have correct relative sizes")
  func testSpacingRelativeSizes() {
    #expect(DSTheme.Spacing.xxs < DSTheme.Spacing.xs)
    #expect(DSTheme.Spacing.xs < DSTheme.Spacing.sm)
    #expect(DSTheme.Spacing.sm < DSTheme.Spacing.md)
    #expect(DSTheme.Spacing.md < DSTheme.Spacing.lg)
    #expect(DSTheme.Spacing.lg < DSTheme.Spacing.xl)
    #expect(DSTheme.Spacing.xl < DSTheme.Spacing.xxl)
  }

  @Test("Spacing tokens are positive values")
  func testSpacingPositiveValues() {
    #expect(DSTheme.Spacing.xxs > 0)
    #expect(DSTheme.Spacing.xs > 0)
    #expect(DSTheme.Spacing.sm > 0)
    #expect(DSTheme.Spacing.md > 0)
    #expect(DSTheme.Spacing.lg > 0)
    #expect(DSTheme.Spacing.xl > 0)
    #expect(DSTheme.Spacing.xxl > 0)
  }
}

struct ThemeRadiusTests {
  @Test("Radius tokens have correct relative sizes")
  func testRadiusRelativeSizes() {
    #expect(DSTheme.Radius.xs < DSTheme.Radius.sm)
    #expect(DSTheme.Radius.sm < DSTheme.Radius.md)
    #expect(DSTheme.Radius.md < DSTheme.Radius.lg)
    #expect(DSTheme.Radius.lg < DSTheme.Radius.xl)
  }

  @Test("Full radius is large enough for pill shapes")
  func testFullRadius() {
    #expect(DSTheme.Radius.full >= 100)
  }
}

struct ThemeIconSizeTests {
  @Test("Icon sizes have correct relative sizes")
  func testIconSizeRelativeSizes() {
    #expect(DSTheme.IconSize.sm < DSTheme.IconSize.md)
    #expect(DSTheme.IconSize.md < DSTheme.IconSize.lg)
    #expect(DSTheme.IconSize.lg < DSTheme.IconSize.xl)
  }
}

struct ThemeAvatarSizeTests {
  @Test("Avatar sizes have correct relative sizes")
  func testAvatarSizeRelativeSizes() {
    #expect(DSTheme.AvatarSize.sm < DSTheme.AvatarSize.md)
    #expect(DSTheme.AvatarSize.md < DSTheme.AvatarSize.lg)
    #expect(DSTheme.AvatarSize.lg < DSTheme.AvatarSize.xl)
    #expect(DSTheme.AvatarSize.xl < DSTheme.AvatarSize.xxl)
  }
}

// MARK: - Button Style Tests

struct DSButtonStyleTests {
  @Test("Primary style has primary background color")
  func testPrimaryStyleBackground() {
    let style = DSButtonStyle.primary
    #expect(style.backgroundColor == DSTheme.Colors.primary)
  }

  @Test("Destructive style has error background color")
  func testDestructiveStyleBackground() {
    let style = DSButtonStyle.destructive
    #expect(style.backgroundColor == DSTheme.Colors.error)
  }

  @Test("Ghost style has clear background")
  func testGhostStyleBackground() {
    let style = DSButtonStyle.ghost
    #expect(style.backgroundColor == .clear)
  }

  @Test("Outline style has clear background")
  func testOutlineStyleBackground() {
    let style = DSButtonStyle.outline
    #expect(style.backgroundColor == .clear)
  }

  @Test("Outline style has border color")
  func testOutlineStyleHasBorder() {
    let style = DSButtonStyle.outline
    #expect(style.borderColor != nil)
  }

  @Test("Primary style has no border")
  func testPrimaryStyleNoBorder() {
    let style = DSButtonStyle.primary
    #expect(style.borderColor == nil)
  }
}

struct DSButtonSizeTests {
  @Test("Button sizes have increasing padding")
  func testButtonSizePadding() {
    #expect(DSButtonSize.small.verticalPadding < DSButtonSize.medium.verticalPadding)
    #expect(DSButtonSize.medium.verticalPadding < DSButtonSize.large.verticalPadding)
  }

  @Test("Button sizes have increasing horizontal padding")
  func testButtonSizeHorizontalPadding() {
    #expect(DSButtonSize.small.horizontalPadding < DSButtonSize.medium.horizontalPadding)
    #expect(DSButtonSize.medium.horizontalPadding < DSButtonSize.large.horizontalPadding)
  }
}

// MARK: - Status Type Tests

struct DSStatusTypeTests {
  @Test("Status type pending has correct label")
  func testPendingLabel() {
    #expect(DSStatusType.pending.label == "Pending")
  }

  @Test("Status type success has correct label")
  func testSuccessLabel() {
    #expect(DSStatusType.success.label == "Success")
  }

  @Test("Status type error has correct label")
  func testErrorLabel() {
    #expect(DSStatusType.error.label == "Error")
  }

  @Test("Status type from string returns pending for 'pending'")
  func testFromStringPending() {
    #expect(DSStatusType.from("pending").label == "Pending")
    #expect(DSStatusType.from("Pending").label == "Pending")
  }

  @Test("Status type from string returns accepted for 'accepted'")
  func testFromStringAccepted() {
    #expect(DSStatusType.from("accepted").label == "Accepted")
  }

  @Test("Status type from string returns cancelled for 'cancelled'")
  func testFromStringCancelled() {
    #expect(DSStatusType.from("cancelled").label == "Cancelled")
    #expect(DSStatusType.from("canceled").label == "Cancelled")
  }

  @Test("Status types have distinct colors")
  func testStatusTypesHaveDistinctColors() {
    // These should have different colors for visual distinction
    let pendingColor = DSStatusType.pending.color
    let successColor = DSStatusType.success.color
    let errorColor = DSStatusType.error.color

    #expect(pendingColor != successColor)
    #expect(successColor != errorColor)
  }
}

// MARK: - Avatar Size Tests

struct DSAvatarSizeTests {
  @Test("Avatar size dimensions are correct")
  func testAvatarSizeDimensions() {
    #expect(DSAvatarSize.xs.dimension < DSAvatarSize.sm.dimension)
    #expect(DSAvatarSize.sm.dimension < DSAvatarSize.md.dimension)
    #expect(DSAvatarSize.md.dimension < DSAvatarSize.lg.dimension)
    #expect(DSAvatarSize.lg.dimension < DSAvatarSize.xl.dimension)
    #expect(DSAvatarSize.xl.dimension < DSAvatarSize.xxl.dimension)
  }

  @Test("Custom avatar size uses provided value")
  func testCustomAvatarSize() {
    let customSize: CGFloat = 100
    #expect(DSAvatarSize.custom(customSize).dimension == customSize)
  }
}

// MARK: - Card Style Tests

struct DSCardStyleTests {
  @Test("All card styles are distinct")
  func testCardStylesAreDistinct() {
    let styles: [DSCardStyle] = [.standard, .elevated, .outlined, .selected, .interactive]
    let uniqueCount = Set(styles.map { "\($0)" }).count
    #expect(uniqueCount == styles.count)
  }
}

// MARK: - Empty State Size Tests

struct DSEmptyStateSizeTests {
  @Test("Empty state sizes have increasing icon sizes")
  func testEmptyStateSizes() {
    #expect(DSEmptyStateSize.compact.iconSize < DSEmptyStateSize.standard.iconSize)
    #expect(DSEmptyStateSize.standard.iconSize < DSEmptyStateSize.large.iconSize)
  }

  @Test("Empty state sizes have increasing spacing")
  func testEmptyStateSpacing() {
    #expect(DSEmptyStateSize.compact.spacing < DSEmptyStateSize.standard.spacing)
    #expect(DSEmptyStateSize.standard.spacing < DSEmptyStateSize.large.spacing)
  }
}

// MARK: - Component Instantiation Tests

struct ComponentInstantiationTests {
  @Test("DSButton can be instantiated with default parameters")
  @MainActor
  func testDSButtonInstantiation() {
    _ = DSButton("Test") {}
  }

  @Test("DSButton can be instantiated with all style variants")
  @MainActor
  func testDSButtonStyles() {
    let styles: [DSButtonStyle] = [.primary, .secondary, .tertiary, .destructive, .ghost, .outline]
    for style in styles {
      _ = DSButton("Test", style: style) {}
    }
  }

  @Test("DSCard can be instantiated with all style variants")
  @MainActor
  func testDSCardStyles() {
    let styles: [DSCardStyle] = [.standard, .elevated, .outlined, .selected, .interactive]
    for style in styles {
      _ = DSCard(style: style) {
        Text("Content")
      }
    }
  }

  @Test("DSTextField can be instantiated with all type variants")
  @MainActor
  func testDSTextFieldTypes() {
    let binding = Binding<String>(get: { "" }, set: { _ in })
    let types: [DSTextFieldType] = [.text, .username, .email, .password, .phone, .number, .multiline]
    for type in types {
      _ = DSTextField("Label", text: binding, type: type)
    }
  }

  @Test("DSStatusBadge can be instantiated with status types")
  @MainActor
  func testDSStatusBadgeTypes() {
    let types: [DSStatusType] = [.pending, .success, .error, .warning, .info, .cancelled]
    for type in types {
      _ = DSStatusBadge(type)
    }
  }

  @Test("DSEmptyState can be instantiated with all size variants")
  @MainActor
  func testDSEmptyStateSizes() {
    let sizes: [DSEmptyStateSize] = [.compact, .standard, .large]
    for size in sizes {
      _ = DSEmptyState(
        icon: "star",
        title: "Title",
        message: "Message",
        size: size
      )
    }
  }

  @Test("DSAvatarView can be instantiated with name")
  @MainActor
  func testDSAvatarViewWithName() {
    _ = DSAvatarView(name: "John Doe", size: .md)
  }

  @Test("DSSection can be instantiated with title and content")
  @MainActor
  func testDSSectionInstantiation() {
    _ = DSSection("Title") {
      Text("Content")
    }
  }

  @Test("DSLoadingView can be instantiated")
  @MainActor
  func testDSLoadingViewInstantiation() {
    _ = DSLoadingView("Loading...")
  }

  @Test("DSInitializingView can be instantiated")
  @MainActor
  func testDSInitializingViewInstantiation() {
    _ = DSInitializingView()
  }

  @Test("DSIconButton can be instantiated")
  @MainActor
  func testDSIconButtonInstantiation() {
    _ = DSIconButton("star.fill") {}
  }

  @Test("DSStatusDot can be instantiated with all status types")
  @MainActor
  func testDSStatusDotTypes() {
    let types: [DSStatusType] = [.pending, .success, .error, .warning, .info, .cancelled]
    for type in types {
      _ = DSStatusDot(type)
    }
  }

  @Test("DSListRow can be instantiated")
  @MainActor
  func testDSListRowInstantiation() {
    _ = DSListRow(title: "Title")
  }
}
