# Connecting to the Zhuzh MCP Server

## Prerequisites

- An **admin** account on the Zhuzh instance
- The backend API service running (HTTP on `:8000`, gRPC on `:8001` for local dev)

## Quick Start (Local Development)

### 1. Start the MCP server

```bash
make mcp
```

This runs the MCP server on port `8888` with hot-reload (via `air`), proxied on port `9999`.

### 2. Test with the MCP Inspector

```bash
make mcp_inspector
```

This launches the `@modelcontextprotocol/inspector` UI, which connects to `http://localhost:8888`. You can browse available tools and invoke them interactively.

## Connecting an MCP Client

### Claude Desktop / Claude Code

Add the server to your MCP client configuration. For a deployed server:

```json
{
  "mcpServers": {
    "zhuzh": {
      "type": "streamable-http",
      "url": "https://mcp.zhuzh.dev/mcp"
    }
  }
}
```

For local development:

```json
{
  "mcpServers": {
    "zhuzh": {
      "type": "streamable-http",
      "url": "http://localhost:8888/mcp"
    }
  }
}
```

When the client connects, it will initiate the OAuth2 flow automatically.

## Authentication Flow

The MCP server implements OAuth2 authorization code with PKCE. Compliant MCP clients handle this transparently, but here's what happens:

### 1. Discovery

The client fetches server metadata:

- `GET /.well-known/oauth-protected-resource` (RFC 9728)
- `GET /.well-known/oauth-authorization-server` (RFC 8414)

### 2. Dynamic Client Registration

The client registers itself (RFC 7591):

```
POST /register
Content-Type: application/json

{
  "redirect_uris": ["http://localhost:..."],
  "client_name": "My MCP Client"
}
```

Returns a `client_id` and `client_secret`.

### 3. Authorization (Login)

The client opens a browser to:

```
GET /authorize?response_type=code&client_id=...&redirect_uri=...&code_challenge=...&code_challenge_method=S256&state=...
```

A login form is rendered. Enter your:

| Field | Description |
|-------|-------------|
| **Username** | Your admin username |
| **Password** | Your admin password |
| **TOTP Code** | Your 2FA code (if enabled) |

On success, the server redirects back to the client with an authorization code.

**Only admin accounts can authenticate.** Regular user accounts will be rejected.

### 4. Token Exchange

The client exchanges the authorization code for tokens:

```
POST /token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&code=...&code_verifier=...&client_id=...&redirect_uri=...
```

Returns:

```json
{
  "access_token": "...",
  "token_type": "Bearer",
  "refresh_token": "...",
  "expires_in": 86400
}
```

### 5. Authenticated MCP Requests

All tool calls go to `POST /mcp` with the bearer token:

```
Authorization: Bearer <access_token>
```

### Token Lifetimes

| Token | Lifetime |
|-------|----------|
| Authorization code | 5 minutes |
| Access token | 24 hours |
| Refresh token | 7 days |

Refresh tokens are single-use. Each refresh returns a new access/refresh token pair.

## Transport Modes

The server supports three transport modes (set via `--transport` flag):

| Transport | Description | Use Case |
|-----------|-------------|----------|
| `http` (default) | Stateless streamable HTTP | Production, API gateways |
| `sse` | Server-Sent Events | Long-lived streaming connections |
| `stdio` | Standard I/O | CLI tools, local piping |

## Troubleshooting

**"Access denied. Admin credentials required."**
Only admin accounts can log in. Ensure you're using admin credentials.

**Token expired / "unknown token"**
Access tokens last 24 hours. Your MCP client should automatically refresh using the refresh token. If refresh also fails (after 7 days), you'll need to re-authenticate.

**Connection refused on localhost:8888**
Make sure the MCP server is running (`make mcp`) and the backend API is running on ports 8000/8001.

**PKCE verification failed**
The client must use S256 code challenge method. This is handled automatically by compliant MCP clients.
