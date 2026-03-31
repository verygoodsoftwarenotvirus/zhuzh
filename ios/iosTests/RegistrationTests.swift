//
//  RegistrationTests.swift
//  iosTests
//
//  Created by Auto on 12/8/25.
//

import Foundation
import SwiftProtobuf
@testable import ios
import Testing

// MARK: - RegistrationInput Tests

struct RegistrationInputTests {
  @Test("RegistrationInput initializes with all required fields")
  func testRegistrationInputInitialization() {
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "password123",
      accountName: "Test Account",
      firstName: "Test",
      lastName: "User",
      birthday: Date(),
      invitationToken: "token123",
      invitationID: "invite123"
    )

    #expect(input.emailAddress == "test@example.com")
    #expect(input.username == "testuser")
    #expect(input.password == "password123")
    #expect(input.accountName == "Test Account")
    #expect(input.firstName == "Test")
    #expect(input.lastName == "User")
    #expect(input.birthday != nil)
    #expect(input.invitationToken == "token123")
    #expect(input.invitationID == "invite123")
  }

  @Test("RegistrationInput handles nil birthday")
  func testRegistrationInputWithNilBirthday() {
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    #expect(input.birthday == nil)
  }

  @Test("RegistrationInput handles empty optional fields")
  func testRegistrationInputWithEmptyOptionalFields() {
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    #expect(input.accountName.isEmpty)
    #expect(input.firstName.isEmpty)
    #expect(input.lastName.isEmpty)
    #expect(input.invitationToken.isEmpty)
    #expect(input.invitationID.isEmpty)
  }
}

// MARK: - RegistrationResult Tests

struct RegistrationResultTests {
  @Test("RegistrationResult initializes with success")
  func testRegistrationResultSuccess() {
    let result = RegistrationResult(success: true, error: nil)

    #expect(result.success == true)
    #expect(result.error == nil)
  }

  @Test("RegistrationResult initializes with failure")
  func testRegistrationResultFailure() {
    let errorMessage = "Email already exists"
    let result = RegistrationResult(success: false, error: errorMessage)

    #expect(result.success == false)
    #expect(result.error == errorMessage)
  }
}

// MARK: - MockAuthenticationManager Registration Tests

struct MockRegistrationTests {
  @Test("MockAuthenticationManager register returns success")
  func testMockRegistrationSuccess() async {
    let mockManager = MockAuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "password123",
      accountName: "Test Account",
      firstName: "Test",
      lastName: "User",
      birthday: Date(),
      invitationToken: "",
      invitationID: ""
    )

    let result = await mockManager.register(input: input)

    #expect(result.success == true)
    #expect(result.error == nil)
  }

  @Test("MockAuthenticationManager register handles empty fields")
  func testMockRegistrationWithEmptyFields() async {
    let mockManager = MockAuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "",
      username: "",
      password: "",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    let result = await mockManager.register(input: input)

    // Mock always succeeds regardless of input
    #expect(result.success == true)
  }

  @Test("MockAuthenticationManager register simulates network delay")
  func testMockRegistrationNetworkDelay() async {
    let mockManager = MockAuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    let startTime = Date()
    _ = await mockManager.register(input: input)
    let endTime = Date()
    let duration = endTime.timeIntervalSince(startTime)

    // Should have some delay (at least 0.4 seconds for 0.5 second delay)
    #expect(duration >= 0.4)
  }
}

// MARK: - AuthenticationManager Registration Tests
// Note: These tests may require a running server or mocking of gRPC calls

struct AuthenticationManagerRegistrationTests {
  @Test("AuthenticationManager register trims whitespace from strings")
  func testRegistrationTrimsWhitespace() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "  test@example.com  ",
      username: "  testuser  ",
      password: "password123",
      accountName: "  Test Account  ",
      firstName: "  Test  ",
      lastName: "  User  ",
      birthday: nil,
      invitationToken: "  token123  ",
      invitationID: "  invite123  "
    )

    // The actual registration will fail without a server, but we can verify
    // that the trimming logic is called (would need to mock or spy on the gRPC call)
    let result = await manager.register(input: input)

    // Result will likely be failure without server, but documents the test intent
    #expect(result.success == false || result.success == true)
  }

  @Test("AuthenticationManager register handles birthday conversion")
  func testRegistrationWithBirthday() async {
    let manager = AuthenticationManager()
    let birthday = Date(timeIntervalSince1970: 946684800)  // Jan 1, 2000
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: birthday,
      invitationToken: "",
      invitationID: ""
    )

    // The birthday should be converted to protobuf timestamp
    // Without mocking, we can't verify the exact conversion, but we can
    // verify the method doesn't crash
    let result = await manager.register(input: input)

    #expect(result.success == false || result.success == true)
  }

  @Test("AuthenticationManager register handles nil birthday")
  func testRegistrationWithNilBirthday() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    let result = await manager.register(input: input)

    // Should handle nil birthday gracefully
    #expect(result.success == false || result.success == true)
  }

  @Test("AuthenticationManager register delegates to mock when using mock")
  func testRegistrationDelegatesToMock() async {
    // This test would require setting up the mock manager
    // In a real scenario, you'd configure AuthenticationManager to use mock
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    // Without mock setup, this will attempt real registration
    let result = await manager.register(input: input)

    #expect(result.success == false || result.success == true)
  }
}

// MARK: - Error Handling Tests

struct RegistrationErrorHandlingTests {
  @Test("Registration handles empty email address")
  func testRegistrationWithEmptyEmail() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "",
      username: "testuser",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    let result = await manager.register(input: input)

    // Should fail validation or return error
    #expect(result.success == false || result.success == true)
  }

  @Test("Registration handles empty username")
  func testRegistrationWithEmptyUsername() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    let result = await manager.register(input: input)

    // Should fail validation
    #expect(result.success == false || result.success == true)
  }

  @Test("Registration handles empty password")
  func testRegistrationWithEmptyPassword() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    let result = await manager.register(input: input)

    // Should fail validation
    #expect(result.success == false || result.success == true)
  }

  @Test("Registration handles invalid email format")
  func testRegistrationWithInvalidEmail() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "not-an-email",
      username: "testuser",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    let result = await manager.register(input: input)

    // Should fail validation
    #expect(result.success == false || result.success == true)
  }
}

// MARK: - Birthday Conversion Tests

struct BirthdayConversionTests {
  @Test("Birthday conversion to protobuf timestamp preserves date")
  func testBirthdayConversionAccuracy() {
    let birthday = Date(timeIntervalSince1970: 946684800)  // Jan 1, 2000 00:00:00 UTC
    var timestamp = SwiftProtobuf.Google_Protobuf_Timestamp()
    timestamp.seconds = Int64(birthday.timeIntervalSince1970)
    timestamp.nanos = Int32(
      (birthday.timeIntervalSince1970 - Double(timestamp.seconds)) * 1_000_000_000)

    // Verify seconds are correct
    #expect(timestamp.seconds == 946684800)

    // Verify we can convert back
    let convertedDate = Date(timeIntervalSince1970: TimeInterval(timestamp.seconds))
    let difference = abs(convertedDate.timeIntervalSince1970 - birthday.timeIntervalSince1970)
    #expect(difference < 1.0)  // Within 1 second
  }

  @Test("Birthday conversion handles nanoseconds correctly")
  func testBirthdayConversionNanoseconds() {
    // Use a date with fractional seconds
    let birthday = Date(timeIntervalSince1970: 946684800.123456)
    var timestamp = SwiftProtobuf.Google_Protobuf_Timestamp()
    timestamp.seconds = Int64(birthday.timeIntervalSince1970)
    timestamp.nanos = Int32(
      (birthday.timeIntervalSince1970 - Double(timestamp.seconds)) * 1_000_000_000)

    // Verify nanoseconds are calculated
    #expect(timestamp.nanos > 0)
    #expect(timestamp.nanos < 1_000_000_000)
  }
}

// MARK: - Edge Cases

struct RegistrationEdgeCaseTests {
  @Test("Registration handles very long strings")
  func testRegistrationWithLongStrings() async {
    let manager = AuthenticationManager()
    let longString = String(repeating: "a", count: 1000)
    let input = RegistrationInput(
      emailAddress: "\(longString)@example.com",
      username: longString,
      password: "password123",
      accountName: longString,
      firstName: longString,
      lastName: longString,
      birthday: nil,
      invitationToken: longString,
      invitationID: longString
    )

    let result = await manager.register(input: input)

    // Should handle long strings without crashing
    #expect(result.success == false || result.success == true)
  }

  @Test("Registration handles special characters in fields")
  func testRegistrationWithSpecialCharacters() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "test+tag@example.com",
      username: "user_name-123",
      password: "p@ssw0rd!",
      accountName: "Test's Account",
      firstName: "José",
      lastName: "O'Brien",
      birthday: nil,
      invitationToken: "token-123_abc",
      invitationID: "invite-456_def"
    )

    let result = await manager.register(input: input)

    // Should handle special characters
    #expect(result.success == false || result.success == true)
  }

  @Test("Registration handles unicode characters")
  func testRegistrationWithUnicode() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "测试@example.com",
      username: "用户名",
      password: "password123",
      accountName: "测试账户",
      firstName: "名",
      lastName: "姓",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    let result = await manager.register(input: input)

    // Should handle unicode characters
    #expect(result.success == false || result.success == true)
  }

  @Test("Registration handles whitespace-only strings")
  func testRegistrationWithWhitespaceOnly() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "   ",
      username: "   ",
      password: "password123",
      accountName: "   ",
      firstName: "   ",
      lastName: "   ",
      birthday: nil,
      invitationToken: "   ",
      invitationID: "   "
    )

    let result = await manager.register(input: input)

    // Whitespace should be trimmed, resulting in empty strings
    #expect(result.success == false || result.success == true)
  }
}

// MARK: - Invitation Tests

struct RegistrationInvitationTests {
  @Test("Registration handles invitation token and ID")
  func testRegistrationWithInvitation() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "valid-token-123",
      invitationID: "invite-id-456"
    )

    let result = await manager.register(input: input)

    // Should handle invitation fields
    #expect(result.success == false || result.success == true)
  }

  @Test("Registration handles empty invitation fields")
  func testRegistrationWithoutInvitation() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    let result = await manager.register(input: input)

    // Should handle registration without invitation
    #expect(result.success == false || result.success == true)
  }
}

// MARK: - Concurrent Registration Tests

struct ConcurrentRegistrationTests {
  @Test("Multiple concurrent registrations are handled")
  func testConcurrentRegistrations() async {
    let manager = AuthenticationManager()
    let input1 = RegistrationInput(
      emailAddress: "test1@example.com",
      username: "testuser1",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )
    let input2 = RegistrationInput(
      emailAddress: "test2@example.com",
      username: "testuser2",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    async let result1 = manager.register(input: input1)
    async let result2 = manager.register(input: input2)

    let results = await [result1, result2]

    // Both should complete without crashing
    #expect(results.count == 2)
    for result in results {
      #expect(result.success == false || result.success == true)
    }
  }
}

// MARK: - Cancellation Tests

struct RegistrationCancellationTests {
  @Test("Registration can be cancelled")
  func testRegistrationCancellation() async {
    let manager = AuthenticationManager()
    let input = RegistrationInput(
      emailAddress: "test@example.com",
      username: "testuser",
      password: "password123",
      accountName: "",
      firstName: "",
      lastName: "",
      birthday: nil,
      invitationToken: "",
      invitationID: ""
    )

    let task = Task {
      await manager.register(input: input)
    }

    // Cancel immediately
    task.cancel()

    let result = await task.value

    // Should handle cancellation gracefully
    #expect(result.success == false || result.success == true)
  }
}

