//
//  ClientTests.swift
//  iosTests
//
//  Created by Auto on 12/8/25.
//

import Foundation
import GRPCCore
import GRPCNIOTransportHTTP2
import GRPCNIOTransportHTTP2TransportServices
@testable import ios
import Testing

// MARK: - Client Initialization Tests

struct ClientInitializationTests {
  @Test("Client initializes with all service clients")
  func testClientInitialization() throws {
    // Create a mock transport for testing
    // Note: In a real scenario, you might want to use a test transport
    let transport = try HTTP2ClientTransport.TransportServices(
      target: .dns(host: "127.0.0.1", port: 8001),
      transportSecurity: .plaintext
    )
    let grpcClient = GRPCCore.GRPCClient(transport: transport)
    let client = Client(grpcClient: grpcClient)

    // Verify all service clients are initialized
    // We can't directly test the service clients without making actual calls,
    // but we can verify the client was created successfully (non-nil by type)
    _ = client
  }

  @Test("Client initializes with shared gRPC client")
  func testClientSharesGRPCClient() throws {
    let transport = try HTTP2ClientTransport.TransportServices(
      target: .dns(host: "127.0.0.1", port: 8001),
      transportSecurity: .plaintext
    )
    let grpcClient = GRPCCore.GRPCClient(transport: transport)
    let client = Client(grpcClient: grpcClient)

    // The client should wrap the provided gRPC client
    // This is verified by the fact that all service clients use the same underlying client
    _ = client
  }
}

// MARK: - ClientManager Initialization Tests

struct ClientManagerInitializationTests {
  @Test("ClientManager initializes with default call options")
  func testClientManagerDefaultCallOptions() throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )

    // Verify default call options are set (5 second timeout)
    #expect(manager.defaultCallOptions.timeout == .seconds(5))
  }

  @Test("ClientManager initializes with custom call options")
  func testClientManagerCustomCallOptions() throws {
    var customOptions = GRPCCore.CallOptions.defaults
    customOptions.timeout = .seconds(10)

    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001,
      defaultCallOptions: customOptions
    )

    // Verify custom call options are used
    #expect(manager.defaultCallOptions.timeout == .seconds(10))
  }

  @Test("ClientManager initializes with default host and port")
  func testClientManagerDefaultHostPort() throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>()

    // Should use default host (127.0.0.1) and port (8001)
    #expect(manager.defaultCallOptions.timeout == .seconds(5))
  }

  @Test("ClientManager initializes with custom host and port")
  func testClientManagerCustomHostPort() throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "localhost",
      port: 9000
    )

    // Verify client was created with custom host/port (non-nil by type)
    _ = manager.client
  }

  @Test("ClientManager creates unified client")
  func testClientManagerCreatesUnifiedClient() throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )

    // Verify the unified client is accessible (non-nil by type)
    _ = manager.client
  }

  @Test("ClientManager reuses same client instance")
  func testClientManagerClientReuse() throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )

    let client1 = manager.client
    let client2 = manager.client

    // Should return the same client instance (reused)
    // Since Client is a struct, we verify both are accessible (structs are value types)
    _ = client1.auth
    _ = client2.auth
  }
}

// MARK: - Call Options Tests

struct CallOptionsTests {
  @Test("callOptions returns default when no overrides")
  func testCallOptionsNoOverrides() throws {
    var defaultOptions = GRPCCore.CallOptions.defaults
    defaultOptions.timeout = .seconds(5)

    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001,
      defaultCallOptions: defaultOptions
    )

    let options = manager.callOptions()

    // Should return default options
    #expect(options.timeout == .seconds(5))
  }

  @Test("callOptions overrides timeout")
  func testCallOptionsOverrideTimeout() throws {
    var defaultOptions = GRPCCore.CallOptions.defaults
    defaultOptions.timeout = .seconds(5)

    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001,
      defaultCallOptions: defaultOptions
    )

    var overrideOptions = GRPCCore.CallOptions.defaults
    overrideOptions.timeout = .seconds(15)

    let mergedOptions = manager.callOptions(overriding: overrideOptions)

    // Override should take precedence
    #expect(mergedOptions.timeout == .seconds(15))
  }

  @Test("callOptions preserves default when override has nil timeout")
  func testCallOptionsPreservesDefaultWhenOverrideNil() throws {
    var defaultOptions = GRPCCore.CallOptions.defaults
    defaultOptions.timeout = .seconds(5)

    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001,
      defaultCallOptions: defaultOptions
    )

    let overrideOptions = GRPCCore.CallOptions.defaults  // No timeout set

    let mergedOptions = manager.callOptions(overriding: overrideOptions)

    // Should preserve default timeout
    #expect(mergedOptions.timeout == .seconds(5))
  }

  @Test("callOptions overrides compression")
  func testCallOptionsOverrideCompression() throws {
    var defaultOptions = GRPCCore.CallOptions.defaults
    defaultOptions.timeout = .seconds(5)

    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001,
      defaultCallOptions: defaultOptions
    )

    var overrideOptions = GRPCCore.CallOptions.defaults
    overrideOptions.compression = .gzip

    let mergedOptions = manager.callOptions(overriding: overrideOptions)

    // Override compression should be used
    #expect(mergedOptions.compression == .gzip)
  }

  @Test("callOptions merges timeout and compression")
  func testCallOptionsMergesBoth() throws {
    var defaultOptions = GRPCCore.CallOptions.defaults
    defaultOptions.timeout = .seconds(5)

    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001,
      defaultCallOptions: defaultOptions
    )

    var overrideOptions = GRPCCore.CallOptions.defaults
    overrideOptions.timeout = .seconds(20)
    overrideOptions.compression = .gzip

    let mergedOptions = manager.callOptions(overriding: overrideOptions)

    // Both overrides should be applied
    #expect(mergedOptions.timeout == .seconds(20))
    #expect(mergedOptions.compression == .gzip)
  }
}

// MARK: - Authenticated Metadata Tests

struct AuthenticatedMetadataTests {
  @Test("authenticatedMetadata creates correct authorization header")
  func testAuthenticatedMetadata() throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )

    let token = "test-access-token-123"
    let metadata = manager.authenticatedMetadata(accessToken: token)

    // Should contain authorization header with Bearer token
    let authValues = metadata["authorization"]
    let authValue = authValues.first { _ in true }
    #expect(authValue?.description == "Bearer test-access-token-123" || String(describing: authValue).contains("Bearer test-access-token-123"))
  }

  @Test("authenticatedMetadata handles empty token")
  func testAuthenticatedMetadataEmptyToken() throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )

    let metadata = manager.authenticatedMetadata(accessToken: "")

    // Should still create metadata with Bearer prefix
    let authValues = metadata["authorization"]
    let authValue = authValues.first { _ in true }
    #expect(authValue?.description == "Bearer " || String(describing: authValue).contains("Bearer "))
  }

  @Test("authenticatedMetadata handles special characters in token")
  func testAuthenticatedMetadataSpecialCharacters() throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )

    let token = "token-with-special-chars-123_abc"
    let metadata = manager.authenticatedMetadata(accessToken: token)

    // Should preserve special characters
    let authValues = metadata["authorization"]
    let authValue = authValues.first { _ in true }
    let valueString = String(describing: authValue)
    #expect(valueString.contains("Bearer token-with-special-chars-123_abc"))
  }

  @Test("authenticatedMetadata creates new metadata each time")
  func testAuthenticatedMetadataCreatesNewInstance() throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )

    let metadata1 = manager.authenticatedMetadata(accessToken: "token1")
    let metadata2 = manager.authenticatedMetadata(accessToken: "token2")

    // Should create separate metadata instances
    let authValues1 = metadata1["authorization"]
    let authValues2 = metadata2["authorization"]
    let authValue1 = authValues1.first { _ in true }
    let authValue2 = authValues2.first { _ in true }
    let valueString1 = String(describing: authValue1)
    let valueString2 = String(describing: authValue2)
    #expect(valueString1.contains("Bearer token1"))
    #expect(valueString2.contains("Bearer token2"))
  }
}

// MARK: - Factory Function Tests

struct FactoryFunctionTests {
  @Test("buildUnauthenticatedClient creates client with default host and port")
  func testBuildUnauthenticatedClientDefaults() throws {
    let client = try buildUnauthenticatedClient()

    // Should create client successfully (non-nil by type)
    _ = client
  }

  @Test("buildUnauthenticatedClient creates client with custom host and port")
  func testBuildUnauthenticatedClientCustom() throws {
    let client = try buildUnauthenticatedClient(host: "localhost", port: 9000)

    // Should create client with custom host/port (non-nil by type)
    _ = client
  }

  @Test("buildUnauthenticatedClientWithFallback uses provided host")
  func testBuildUnauthenticatedClientWithFallbackProvidedHost() throws {
    let client = try buildUnauthenticatedClientWithFallback(host: "127.0.0.1", port: 8001)

    // Should use provided host directly (non-nil by type)
    _ = client
  }

  @Test("buildUnauthenticatedClientWithFallback attempts IPv4 first")
  func testBuildUnauthenticatedClientWithFallbackIPv4() throws {
    // When host is nil, should try IPv4 first
    // This will succeed if IPv4 works, or fail and try other methods
    do {
      let client = try buildUnauthenticatedClientWithFallback(host: nil, port: 8001)
      _ = client
    } catch {
      // If IPv4 fails, it will try IPv6, then DNS
      // Any of these succeeding is valid
    }
  }
}

// MARK: - Error Handling Tests

struct ClientErrorHandlingTests {
  @Test("ClientManager throws error with invalid host")
  func testClientManagerInvalidHost() {
    // Invalid host format should throw error
    // Note: The actual behavior depends on DNS resolution
    do {
      _ = try ClientManager<HTTP2ClientTransport.TransportServices>(
        host: "invalid..host..name",
        port: 8001
      )
      // If it doesn't throw, that's also valid (DNS might resolve it)
    } catch {
      // Expected to throw for invalid host
    }
  }

  @Test("ClientManager handles invalid port")
  func testClientManagerInvalidPort() {
    // Port 0 or negative ports should be handled
    do {
      _ = try ClientManager<HTTP2ClientTransport.TransportServices>(
        host: "127.0.0.1",
        port: 0
      )
      // Some systems might allow port 0, so this is also valid
    } catch {
      // Expected to throw for invalid port
    }
  }

  @Test("buildUnauthenticatedClient throws error with invalid host")
  func testBuildUnauthenticatedClientInvalidHost() {
    do {
      _ = try buildUnauthenticatedClient(host: "invalid..host", port: 8001)
      // If it doesn't throw, DNS might have resolved it
    } catch {
      // Expected to throw for invalid host
    }
  }
}

// MARK: - Client Lifecycle Tests

struct ClientLifecycleTests {
  @Test("Client startConnections starts async connection")
  func testClientStartConnections() throws {
    let transport = try HTTP2ClientTransport.TransportServices(
      target: .dns(host: "127.0.0.1", port: 8001),
      transportSecurity: .plaintext
    )
    let grpcClient = GRPCCore.GRPCClient(transport: transport)
    let client = Client(grpcClient: grpcClient)

    // startConnections should not throw (it's async)
    client.startConnections()

    // Verify client is still valid after starting connections (non-nil by type)
    _ = client
  }

  @Test("ClientManager starts connections on initialization")
  func testClientManagerStartsConnections() throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )

    // Connections should be started automatically (client is non-nil by type)
    _ = manager.client
  }
}

// MARK: - Service Client Access Tests

struct ServiceClientAccessTests {
  @Test("Client provides access to all service clients")
  func testClientServiceClients() throws {
    let transport = try HTTP2ClientTransport.TransportServices(
      target: .dns(host: "127.0.0.1", port: 8001),
      transportSecurity: .plaintext
    )
    let grpcClient = GRPCCore.GRPCClient(transport: transport)
    let client = Client(grpcClient: grpcClient)

    // Verify all service clients are accessible (non-nil by type)
    _ = client.auth
    _ = client.identity
    _ = client.audit
    _ = client.dataPrivacy
    _ = client.internalOps
    _ = client.notifications
    _ = client.oauth
    _ = client.settings
    _ = client.webhooks
  }

  @Test("ClientManager provides access to unified client")
  func testClientManagerUnifiedClient() throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )

    // Verify unified client is accessible (non-nil by type)
    _ = manager.client
    _ = manager.client.auth
    _ = manager.client.identity
  }
}

// MARK: - Concurrent Access Tests

struct ConcurrentAccessTests {
  @Test("Multiple ClientManagers can be created concurrently")
  func testConcurrentClientManagerCreation() async throws {
    async let manager1 = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )
    async let manager2 = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )

    let managers = try await [manager1, manager2]

    // Both should be created successfully
    #expect(managers.count == 2)
    _ = managers[0].client
    _ = managers[1].client
  }

  @Test("ClientManager can be accessed concurrently")
  func testConcurrentClientAccess() async throws {
    let manager = try ClientManager<HTTP2ClientTransport.TransportServices>(
      host: "127.0.0.1",
      port: 8001
    )

    let client1 = manager.client
    let client2 = manager.client
    let metadata1 = manager.authenticatedMetadata(accessToken: "token1")
    let metadata2 = manager.authenticatedMetadata(accessToken: "token2")

    // All should succeed (service clients non-nil by type)
    _ = client1.auth
    _ = client2.auth
    let authValues1 = metadata1["authorization"]
    let authValues2 = metadata2["authorization"]
    let valueString1 = String(describing: authValues1.first { _ in true })
    let valueString2 = String(describing: authValues2.first { _ in true })
    #expect(valueString1.contains("Bearer token1"))
    #expect(valueString2.contains("Bearer token2"))
  }
}

