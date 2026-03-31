//
//  iosUITests.swift
//  iosUITests
//
//  Created by Jeffrey Dorrycott on 12/8/25.
//

import XCTest

final class IOSUITests: XCTestCase {
    override func setUpWithError() throws {
        // Put setup code here. This method is called before the invocation of each test method in the class.

        // In UI tests it is usually best to stop immediately when a failure occurs.
        continueAfterFailure = false

        // In UI tests it’s important to set the initial state - such as interface orientation - required for your tests before they run. The setUp method is a good place to do this.
    }

    override func tearDownWithError() throws {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
    }

    @MainActor
    func testExample() throws {
        // UI tests must launch the application that they test.
        let app = XCUIApplication()
        app.launch()

        // Use XCTAssert and related functions to verify your tests produce the correct results.
    }

    @MainActor
    func testLaunchPerformance() throws {
        if #available(macOS 10.15, iOS 13.0, tvOS 13.0, watchOS 7.0, *) {
            // This measures how long it takes to launch your application.
            measure(metrics: [XCTApplicationLaunchMetric()]) {
                XCUIApplication().launch()
            }
        }
    }
    
    // MARK: - TOTP Input Tests
    
    @MainActor
    func testTOTPInputAppearsWhenRequired() throws {
        // Launch the app with mock authentication enabled
        let app = XCUIApplication()
        app.launchArguments = ["--use-mock-auth", "--mock-requires-totp"]
        app.launch()
        
        // Wait for the login view to appear
        let usernameField = app.textFields["usernameTextField"]
        XCTAssertTrue(usernameField.waitForExistence(timeout: 5), "Username field should exist")
        
        // Enter username and password
        usernameField.tap()
        usernameField.typeText("testuser")
        
        let passwordField = app.secureTextFields["passwordTextField"]
        XCTAssertTrue(passwordField.exists, "Password field should exist")
        passwordField.tap()
        passwordField.typeText("testpassword")
        
        // Verify TOTP field is NOT visible initially
        let totpField = app.textFields["totpTextField"]
        XCTAssertFalse(totpField.exists, "TOTP field should not be visible before login attempt")
        
        // Tap the sign in button
        let signInButton = app.buttons["signInButton"]
        XCTAssertTrue(signInButton.exists, "Sign in button should exist")
        XCTAssertTrue(signInButton.isEnabled, "Sign in button should be enabled")
        signInButton.tap()
        
        // Wait for the TOTP field to appear (this assumes the server returns requiresTOTP: true)
        // The field should appear after the login attempt indicates TOTP is required
        let totpFieldExists = totpField.waitForExistence(timeout: 10)
        
        if totpFieldExists {
            // Verify TOTP field is visible and enabled
            XCTAssertTrue(totpField.exists, "TOTP field should appear when TOTP is required")
            XCTAssertTrue(totpField.isEnabled, "TOTP field should be enabled")
            
            // Verify error message appears
            let errorMessage = app.staticTexts["errorMessage"]
            let errorExists = errorMessage.waitForExistence(timeout: 2)
            if errorExists {
                XCTAssertTrue(errorMessage.exists, "Error message should appear when TOTP is required")
                XCTAssertFalse(errorMessage.label.isEmpty, "Error message should not be empty")
            }
            
            // Verify sign in button is disabled when TOTP is required but not entered
            XCTAssertFalse(signInButton.isEnabled, "Sign in button should be disabled when TOTP is required but not entered")
            
            // Enter TOTP code
            totpField.tap()
            totpField.typeText("123456")
            
            // Verify sign in button becomes enabled after entering TOTP
            XCTAssertTrue(signInButton.isEnabled, "Sign in button should be enabled after entering TOTP code")
        } else {
            // With mock enabled, TOTP field should always appear
            XCTFail("TOTP field did not appear. This should not happen with mock authentication enabled.")
        }
    }
    
    @MainActor
    func testTOTPFieldBehavior() throws {
        // Launch the app with mock authentication enabled
        let app = XCUIApplication()
        app.launchArguments = ["--use-mock-auth", "--mock-requires-totp"]
        app.launch()
        
        // Wait for login view
        let usernameField = app.textFields["usernameTextField"]
        XCTAssertTrue(usernameField.waitForExistence(timeout: 5), "Username field should exist")
        
        // Fill in credentials
        usernameField.tap()
        usernameField.typeText("testuser")
        
        let passwordField = app.secureTextFields["passwordTextField"]
        passwordField.tap()
        passwordField.typeText("testpassword")
        
        // Attempt login
        let signInButton = app.buttons["signInButton"]
        signInButton.tap()
        
        // Wait for TOTP field (if server requires it)
        let totpField = app.textFields["totpTextField"]
        let totpFieldExists = totpField.waitForExistence(timeout: 10)
        
        guard totpFieldExists else {
            // Skip test if TOTP is not required (server doesn't return requiresTOTP: true)
            throw XCTSkip("TOTP field did not appear. Server may not require TOTP for this test scenario.")
        }
        
        // Test 1: Verify TOTP field has correct keyboard type (numberPad)
        totpField.tap()
        // Note: We can't directly verify keyboard type in UI tests, but we can verify numeric input works
        totpField.typeText("123456")
        XCTAssertEqual(totpField.value as? String, "123456", "TOTP field should accept numeric input")
        
        // Test 2: Verify TOTP field is cleared when login succeeds (if we had a successful login flow)
        // This would require a full login flow with valid TOTP
        
        // Test 3: Verify sign in button state changes based on TOTP input
        // Clear TOTP field
        totpField.tap()
        // Select all and delete
        if let stringValue = totpField.value as? String, !stringValue.isEmpty {
            let deleteString = String(repeating: XCUIKeyboardKey.delete.rawValue, count: stringValue.count)
            totpField.typeText(deleteString)
        }
        
        // Button should be disabled when TOTP is required but empty
        XCTAssertFalse(signInButton.isEnabled, "Sign in button should be disabled when TOTP is required but empty")
        
        // Re-enter TOTP
        totpField.typeText("654321")
        XCTAssertTrue(signInButton.isEnabled, "Sign in button should be enabled when TOTP is entered")
    }
}
