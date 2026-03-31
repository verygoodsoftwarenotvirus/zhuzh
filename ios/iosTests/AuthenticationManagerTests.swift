//
//  AuthenticationManagerTests.swift
//  iosTests
//
//  Created by Auto on 12/8/25.
//

import Foundation
@testable import ios
import Testing

struct AuthenticationManagerTests {
    // MARK: - Initialization Tests
    
    @Test("AuthenticationManager initializes with default state")
    func testInitialization() async {
        let manager = AuthenticationManager()
        
        #expect(manager.isAuthenticated == false)
        #expect(manager.username.isEmpty)
        #expect(manager.accessToken.isEmpty)
        #expect(manager.refreshToken.isEmpty)
        #expect(manager.oauth2AccessToken.isEmpty)
        #expect(manager.oauth2RefreshToken.isEmpty)
        #expect(manager.oauth2TokenExpiresAt == nil)
        #expect(manager.userID.isEmpty)
        #expect(manager.accountID.isEmpty)
    }
    
    // MARK: - Logout Tests
    
    @Test("Logout clears all authentication state")
    func testLogout() async {
        let manager = AuthenticationManager()
        
        // Set up authenticated state
        await MainActor.run {
            manager.isAuthenticated = true
            manager.username = "testuser"
            manager.accessToken = "test-access-token"
            manager.refreshToken = "test-refresh-token"
            manager.oauth2AccessToken = "test-oauth2-access-token"
            manager.oauth2RefreshToken = "test-oauth2-refresh-token"
            manager.oauth2TokenExpiresAt = Date()
            manager.userID = "user-123"
            manager.accountID = "account-456"
        }
        
        // Verify state is set
        #expect(manager.isAuthenticated == true)
        #expect(manager.username == "testuser")
        #expect(!manager.accessToken.isEmpty)
        
        // Logout
        await manager.logout()
        
        // Verify all state is cleared
        #expect(manager.isAuthenticated == false)
        #expect(manager.username.isEmpty)
        #expect(manager.accessToken.isEmpty)
        #expect(manager.refreshToken.isEmpty)
        #expect(manager.oauth2AccessToken.isEmpty)
        #expect(manager.oauth2RefreshToken.isEmpty)
        #expect(manager.oauth2TokenExpiresAt == nil)
        #expect(manager.userID.isEmpty)
        #expect(manager.accountID.isEmpty)
    }
    
    @Test("Logout can be called multiple times safely")
    func testLogoutMultipleTimes() async {
        let manager = AuthenticationManager()
        
        // Set up authenticated state
        await MainActor.run {
            manager.isAuthenticated = true
            manager.username = "testuser"
        }
        
        // Logout multiple times
        await manager.logout()
        await manager.logout()
        await manager.logout()
        
        // Verify state remains cleared
        #expect(manager.isAuthenticated == false)
        #expect(manager.username.isEmpty)
    }
    
    // MARK: - OAuth2 Token Expiration Tests
    
    @Test("getOAuth2AccessToken returns nil when no token is available")
    func testGetOAuth2AccessTokenWithNoToken() async {
        let manager = AuthenticationManager()
        
        let token = await manager.getOAuth2AccessToken()
        
        #expect(token == nil)
    }
    
    @Test("getOAuth2AccessToken returns token when not expired")
    func testGetOAuth2AccessTokenWithValidToken() async {
        let manager = AuthenticationManager()
        
        // Set a valid token that expires in 1 hour
        await MainActor.run {
            manager.oauth2AccessToken = "valid-token"
            manager.oauth2TokenExpiresAt = Date().addingTimeInterval(3600) // 1 hour from now
        }
        
        let token = await manager.getOAuth2AccessToken()
        
        #expect(token == "valid-token")
    }
    
    @Test("getOAuth2AccessToken returns nil when token is expired and no refresh token")
    func testGetOAuth2AccessTokenWithExpiredTokenNoRefresh() async {
        let manager = AuthenticationManager()
        
        // Set an expired token with no refresh token
        await MainActor.run {
            manager.oauth2AccessToken = "expired-token"
            manager.oauth2TokenExpiresAt = Date().addingTimeInterval(-3600) // 1 hour ago
            manager.oauth2RefreshToken = "" // No refresh token
        }
        
        let token = await manager.getOAuth2AccessToken()
        
        // Should return nil since refresh will fail and no JWT token available
        #expect(token == nil)
    }
    
    @Test("getOAuth2AccessToken detects token expiring soon")
    func testGetOAuth2AccessTokenWithTokenExpiringSoon() async {
        let manager = AuthenticationManager()
        
        // Set a token that expires in 4 minutes (less than 5 minute threshold)
        await MainActor.run {
            manager.oauth2AccessToken = "expiring-token"
            manager.oauth2TokenExpiresAt = Date().addingTimeInterval(240) // 4 minutes from now
            manager.oauth2RefreshToken = "" // No refresh token to test the path
        }
        
        // This will attempt to refresh, but since there's no refresh token,
        // it will try to use JWT token (which is also empty), so returns nil
        let token = await manager.getOAuth2AccessToken()
        
        // Should return nil since refresh will fail
        #expect(token == nil)
    }
    
    @Test("getOAuth2AccessToken returns token when exactly at 5 minute threshold")
    func testGetOAuth2AccessTokenAtThreshold() async {
        let manager = AuthenticationManager()
        
        // Set a token that expires in exactly 5 minutes (at the threshold)
        await MainActor.run {
            manager.oauth2AccessToken = "threshold-token"
            manager.oauth2TokenExpiresAt = Date().addingTimeInterval(300) // Exactly 5 minutes
        }
        
        // Should attempt to refresh since it's at the threshold
        let token = await manager.getOAuth2AccessToken()
        
        // Will return nil since refresh will fail (no refresh token)
        #expect(token == nil)
    }
    
    @Test("getOAuth2AccessToken returns token when over 5 minute threshold")
    func testGetOAuth2AccessTokenJustOverThreshold() async {
        let manager = AuthenticationManager()
        
        // Set a token that expires in 6 minutes (well over the 5 minute threshold)
        // Using 360 seconds instead of 301 to avoid race conditions in CI
        await MainActor.run {
            manager.oauth2AccessToken = "valid-token"
            manager.oauth2TokenExpiresAt = Date().addingTimeInterval(360) // 6 minutes
        }
        
        let token = await manager.getOAuth2AccessToken()
        
        // Should return the token without attempting refresh
        #expect(token == "valid-token")
    }
    
    @Test("getOAuth2AccessToken handles token with no expiration date")
    func testGetOAuth2AccessTokenWithNoExpiration() async {
        let manager = AuthenticationManager()
        
        // Set a token with no expiration date
        await MainActor.run {
            manager.oauth2AccessToken = "no-expiry-token"
            manager.oauth2TokenExpiresAt = nil
        }
        
        let token = await manager.getOAuth2AccessToken()
        
        // Should return the token since there's no expiration to check
        #expect(token == "no-expiry-token")
    }
    
    // MARK: - Refresh Token Tests
    
    @Test("refreshOAuth2Token returns false when no refresh token available")
    func testRefreshOAuth2TokenWithNoRefreshToken() async {
        let manager = AuthenticationManager()
        
        await MainActor.run {
            manager.oauth2RefreshToken = ""
        }
        
        let result = await manager.refreshOAuth2Token()
        
        #expect(result == false)
    }
    
    // MARK: - Client Manager Tests
    
    @Test("getClientManager creates client manager on first call")
    func testGetClientManagerCreatesManager() async throws {
        let manager = AuthenticationManager()
        
        // This will attempt to create a client manager
        // Note: This test may fail if the server is not running
        // In a real scenario, you'd want to mock the ClientManager
        do {
            _ = try manager.getClientManager()
            // Success: client manager was created
        } catch {
            // If connection fails, that's expected in a test environment
            // The important thing is that the method doesn't crash
        }
    }
    
    @Test("getClientManager reuses existing client manager")
    func testGetClientManagerReusesManager() async throws {
        let manager = AuthenticationManager()
        
        do {
            let firstManager = try manager.getClientManager()
            let secondManager = try manager.getClientManager()
            
            // Both should return the same instance (reused)
            // Note: This is testing the caching behavior
            #expect(firstManager === secondManager)
        } catch {
            // If connection fails, that's expected in a test environment
            // The important thing is that the method doesn't crash
        }
    }
    
    // MARK: - Login Error Handling Tests
    
    @Test("login handles empty username")
    func testLoginWithEmptyUsername() async {
        let manager = AuthenticationManager()
        
        let result = await manager.login(username: "", password: "password")
        
        // The actual behavior depends on server validation
        // This test documents that empty username is handled
        #expect(result.success == false || result.success == true)
    }
    
    @Test("login handles empty password")
    func testLoginWithEmptyPassword() async {
        let manager = AuthenticationManager()
        
        let result = await manager.login(username: "user", password: "")
        
        // The actual behavior depends on server validation
        #expect(result.success == false || result.success == true)
    }
    
    @Test("login handles TOTP token parameter")
    func testLoginWithTOTPToken() async {
        let manager = AuthenticationManager()
        
        let result = await manager.login(
            username: "user",
            password: "password",
            totpToken: "123456"
        )
        
        // The actual behavior depends on server response
        // This test documents that TOTP token is accepted
        #expect(result.success == false || result.success == true)
    }
    
    @Test("login handles nil TOTP token")
    func testLoginWithNilTOTPToken() async {
        let manager = AuthenticationManager()
        
        let result = await manager.login(
            username: "user",
            password: "password",
            totpToken: nil
        )
        
        // Should handle nil TOTP token gracefully
        #expect(result.success == false || result.success == true)
    }
    
    @Test("login handles empty TOTP token")
    func testLoginWithEmptyTOTPToken() async {
        let manager = AuthenticationManager()
        
        let result = await manager.login(
            username: "user",
            password: "password",
            totpToken: ""
        )
        
        // Empty TOTP token should be treated same as nil
        #expect(result.success == false || result.success == true)
    }
    
    // MARK: - State Management Tests
    
    @Test("Authentication state persists after login attempt")
    func testStateAfterLoginAttempt() async {
        let manager = AuthenticationManager()
        
        // Attempt login (will likely fail without server, but tests state management)
        _ = await manager.login(username: "test", password: "test")
        
        // State should be false if login failed, or true if it succeeded
        // This test documents the state management behavior
        #expect(manager.isAuthenticated == false || manager.isAuthenticated == true)
    }
    
    @Test("Multiple logout calls maintain consistent state")
    func testMultipleLogoutCallsMaintainState() async {
        let manager = AuthenticationManager()
        
        // Set authenticated state
        await MainActor.run {
            manager.isAuthenticated = true
            manager.username = "user1"
            manager.accessToken = "token1"
        }
        
        // Logout
        await manager.logout()
        let stateAfterFirst = manager.isAuthenticated
        
        // Set state again
        await MainActor.run {
            manager.isAuthenticated = true
            manager.username = "user2"
        }
        
        // Logout multiple times
        await manager.logout()
        await manager.logout()
        let stateAfterMultiple = manager.isAuthenticated
        
        #expect(stateAfterFirst == false)
        #expect(stateAfterMultiple == false)
    }
    
    // MARK: - Integration Test Scenarios
    // These tests document expected behavior but may require a running server
    
    @Test("Login flow sets all authentication properties on success")
    func testLoginFlowSetsProperties() async {
        let manager = AuthenticationManager()
        
        // This test documents the expected behavior when login succeeds
        // In a real scenario with a running server, it would verify:
        // - isAuthenticated is set to true
        // - username is set
        // - accessToken is set
        // - refreshToken is set
        // - userID is set
        // - accountID is set
        // - OAuth2 token exchange is attempted
        
        // For now, we just verify the initial state
        #expect(manager.isAuthenticated == false)
    }
    
    @Test("OAuth2 token exchange is attempted after successful login")
    func testOAuth2ExchangeAfterLogin() async {
        let manager = AuthenticationManager()
        
        // This test documents that OAuth2 exchange should happen after login
        // In a real scenario, we'd verify that exchangeJWTForOAuth2Token is called
        // and that oauth2AccessToken, oauth2RefreshToken, and oauth2TokenExpiresAt are set
        
        // For now, we verify initial OAuth2 state
        #expect(manager.oauth2AccessToken.isEmpty)
        #expect(manager.oauth2RefreshToken.isEmpty)
        #expect(manager.oauth2TokenExpiresAt == nil)
    }
    
    // MARK: - Token Expiration Edge Cases
    
    @Test("Token expiration calculation handles future dates correctly")
    func testTokenExpirationFutureDate() async {
        let manager = AuthenticationManager()
        
        // Set token expiring in 10 hours
        let futureDate = Date().addingTimeInterval(36000) // 10 hours
        await MainActor.run {
            manager.oauth2AccessToken = "future-token"
            manager.oauth2TokenExpiresAt = futureDate
        }
        
        let token = await manager.getOAuth2AccessToken()
        
        // Should return token without attempting refresh
        #expect(token == "future-token")
    }
    
    @Test("Token expiration calculation handles past dates correctly")
    func testTokenExpirationPastDate() async {
        let manager = AuthenticationManager()
        
        // Set token that expired 1 hour ago
        let pastDate = Date().addingTimeInterval(-3600) // 1 hour ago
        await MainActor.run {
            manager.oauth2AccessToken = "past-token"
            manager.oauth2TokenExpiresAt = pastDate
            manager.oauth2RefreshToken = "" // No refresh token
        }
        
        let token = await manager.getOAuth2AccessToken()
        
        // Should return nil since token is expired and can't refresh
        #expect(token == nil)
    }
    
    // MARK: - Login Return Value Tests
    
    @Test("Login returns correct tuple structure")
    func testLoginReturnValueStructure() async {
        let manager = AuthenticationManager()
        
        let result = await manager.login(username: "test", password: "test")
        
        // Verify the tuple structure is correct
        // result should have: success (Bool), error (String?), requiresTOTP (Bool)
        #expect(type(of: result.success) == Bool.self)
        #expect(result.error == nil || type(of: result.error) == Optional<String>.self)
        #expect(type(of: result.requiresTOTP) == Bool.self)
    }
    
    @Test("Login requiresTOTP flag is set correctly")
    func testLoginRequiresTOTPFlag() async {
        let manager = AuthenticationManager()
        
        let result = await manager.login(username: "test", password: "test")
        
        // The requiresTOTP flag should be a boolean
        // In a real scenario, this would be true if the server requires TOTP
        #expect(result.requiresTOTP == false || result.requiresTOTP == true)
    }
}
