package main

import (
	"context"
	"time"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	grpcconverters "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/converters"
	identitygrpc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/generated/services/identity"
	identityconverters "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/identity/grpc/converters"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var userSchema = map[string]any{
	"ID":            stringField("The ID of the user"),
	"Username":      stringField("The user's username"),
	"EmailAddress":  stringField("The user's email address"),
	"FirstName":     stringField("The user's first name"),
	"LastName":      stringField("The user's last name"),
	"AccountStatus": stringField("The user's account status (good, unverified, banned, terminated)"),
	"CreatedAt":     timestampField("When the user was created"),
	"LastUpdatedAt": timestampField("When the user was last updated"),
	"ArchivedAt":    timestampField("When the user was archived"),
}

var accountSchema = map[string]any{
	"ID":            stringField("The ID of the account"),
	"Name":          stringField("The account name"),
	"BillingStatus": stringField("The account's billing status"),
	"BelongsToUser": stringField("The ID of the user who owns this account"),
	"ContactPhone":  stringField("The account's contact phone number"),
	"City":          stringField("The account's city"),
	"State":         stringField("The account's state"),
	"ZipCode":       stringField("The account's zip code"),
	"Country":       stringField("The account's country"),
	"AddressLine1":  stringField("The account's address line 1"),
	"AddressLine2":  stringField("The account's address line 2"),
	"CreatedAt":     timestampField("When the account was created"),
	"LastUpdatedAt": timestampField("When the account was last updated"),
	"ArchivedAt":    timestampField("When the account was archived"),
}

var getUserTool = &mcp.Tool{
	Name:        "GetUser",
	Description: "Get a user by their ID",
	InputSchema: schemaObject(map[string]any{
		"UserID": stringField("The ID of the user to get"),
	}),
	OutputSchema: schemaObject(userSchema),
}

type GetUserInvocation struct {
	UserID string `jsonschema:"description=The user ID"`
}

func (h *mcpToolManager) GetUser() mcp.ToolHandlerFor[*GetUserInvocation, *identity.User] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *GetUserInvocation) (*mcp.CallToolResult, *identity.User, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		result, err := c.GetUser(ctx, &identitygrpc.GetUserRequest{
			UserId: x.UserID,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, identityconverters.ConvertGRPCUserToUser(result.Result), nil
	}
}

var getUsersTool = &mcp.Tool{
	Name:        "GetUsers",
	Description: "Get users with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(userSchema)),
	}),
}

type (
	GetUsersInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetUsersResult struct {
		Results []*identity.User
	}
)

func (h *mcpToolManager) GetUsers() mcp.ToolHandlerFor[*GetUsersInvocation, *GetUsersResult] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *GetUsersInvocation) (*mcp.CallToolResult, *GetUsersResult, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		results, err := c.GetUsers(ctx, &identitygrpc.GetUsersRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetUsersResult{}
		for _, u := range results.Results {
			out.Results = append(out.Results, identityconverters.ConvertGRPCUserToUser(u))
		}
		return nil, out, nil
	}
}

var searchForUsersTool = &mcp.Tool{
	Name:        "SearchForUsers",
	Description: "Search for users by query string",
	InputSchema: schemaObject(map[string]any{
		"Query":            stringField("The search query string"),
		"UseSearchService": boolField("Whether to use the search service (if false, uses database search)"),
		"Filter":           queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(userSchema)),
	}),
}

type (
	SearchForUsersInvocation struct {
		Filter           *filtering.QueryFilter
		Query            string
		UseSearchService bool
	}

	SearchForUsersResult struct {
		Results []*identity.User
	}
)

func (h *mcpToolManager) SearchForUsers() mcp.ToolHandlerFor[*SearchForUsersInvocation, *SearchForUsersResult] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *SearchForUsersInvocation) (*mcp.CallToolResult, *SearchForUsersResult, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		results, err := c.SearchForUsers(ctx, &identitygrpc.SearchForUsersRequest{
			Query:            x.Query,
			UseSearchService: x.UseSearchService,
			Filter:           grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &SearchForUsersResult{}
		for _, u := range results.Results {
			out.Results = append(out.Results, identityconverters.ConvertGRPCUserToUser(u))
		}
		return nil, out, nil
	}
}

var getAccountTool = &mcp.Tool{
	Name:        "GetAccount",
	Description: "Get an account by its ID",
	InputSchema: schemaObject(map[string]any{
		"AccountID": stringField("The ID of the account to get"),
	}),
	OutputSchema: schemaObject(accountSchema),
}

type GetAccountInvocation struct {
	AccountID string `jsonschema:"description=The account ID"`
}

func (h *mcpToolManager) GetAccount() mcp.ToolHandlerFor[*GetAccountInvocation, *identity.Account] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *GetAccountInvocation) (*mcp.CallToolResult, *identity.Account, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		result, err := c.GetAccount(ctx, &identitygrpc.GetAccountRequest{
			AccountId: x.AccountID,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, identityconverters.ConvertGRPCAccountToAccount(result.Result), nil
	}
}

var getAccountsForUserTool = &mcp.Tool{
	Name:        "GetAccountsForUser",
	Description: "Get accounts belonging to a specific user",
	InputSchema: schemaObject(map[string]any{
		"UserID": stringField("The user ID"),
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(accountSchema)),
	}),
}

type (
	GetAccountsForUserInvocation struct {
		Filter *filtering.QueryFilter
		UserID string `jsonschema:"description=The user ID"`
	}

	GetAccountsForUserResult struct {
		Results []*identity.Account
	}
)

func (h *mcpToolManager) GetAccountsForUser() mcp.ToolHandlerFor[*GetAccountsForUserInvocation, *GetAccountsForUserResult] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *GetAccountsForUserInvocation) (*mcp.CallToolResult, *GetAccountsForUserResult, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		results, err := c.GetAccountsForUser(ctx, &identitygrpc.GetAccountsForUserRequest{
			UserId: x.UserID,
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetAccountsForUserResult{}
		for _, a := range results.Results {
			out.Results = append(out.Results, identityconverters.ConvertGRPCAccountToAccount(a))
		}
		return nil, out, nil
	}
}

var updateUserDetailsTool = &mcp.Tool{
	Name:        "UpdateUserDetails",
	Description: "Update a user's details (first name, last name, birthday)",
	InputSchema: schemaObject(map[string]any{
		"FirstName":       stringField("The user's new first name"),
		"LastName":        stringField("The user's new last name"),
		"Birthday":        timestampField("The user's birthday (ISO 8601)"),
		"CurrentPassword": stringField("The user's current password for verification"),
		"TOTPToken":       stringField("The user's current TOTP token for verification"),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Success": boolField("Whether the update was successful"),
	}),
}

type UpdateUserDetailsInvocation struct {
	Birthday        time.Time `jsonschema:"description=The user's birthday"`
	FirstName       string    `jsonschema:"description=The user's new first name"`
	LastName        string    `jsonschema:"description=The user's new last name"`
	CurrentPassword string    `jsonschema:"description=The user's current password"`
	TOTPToken       string    `jsonschema:"description=The user's current TOTP token"`
}

func (h *mcpToolManager) UpdateUserDetails() mcp.ToolHandlerFor[*UpdateUserDetailsInvocation, *boolResult] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *UpdateUserDetailsInvocation) (*mcp.CallToolResult, *boolResult, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		_, err = c.UpdateUserDetails(ctx, &identitygrpc.UpdateUserDetailsRequest{
			Input: &identitygrpc.UserDetailsUpdateRequestInput{
				FirstName:       x.FirstName,
				LastName:        x.LastName,
				Birthday:        grpcconverters.ConvertTimeToPBTimestamp(x.Birthday),
				CurrentPassword: x.CurrentPassword,
				TotpToken:       x.TOTPToken,
			},
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, &boolResult{Success: true}, nil
	}
}
