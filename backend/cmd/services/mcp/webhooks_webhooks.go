package main

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/webhooks"
	grpcconverters "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/converters"
	webhooksgrpc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/generated/services/webhooks"
	webhooksconverters "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/webhooks/grpc/converters"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var webhookTriggerConfigSchema = map[string]any{
	"ID":               stringField("The ID of the trigger config"),
	"BelongsToWebhook": stringField("The ID of the webhook this config belongs to"),
	"TriggerEventID":   stringField("The ID of the trigger event"),
	"CreatedAt":        timestampField("When the trigger config was created"),
	"ArchivedAt":       timestampField("When the trigger config was archived"),
}

var webhookSchema = map[string]any{
	"ID":               stringField("The ID of the webhook"),
	"Name":             stringField("The webhook name"),
	"URL":              stringField("The webhook URL"),
	"Method":           stringField("The HTTP method (GET, POST, PUT, PATCH, DELETE)"),
	"ContentType":      stringField("The content type (application/json or application/xml)"),
	"BelongsToAccount": stringField("The ID of the account this webhook belongs to"),
	"CreatedByUser":    stringField("The ID of the user who created this webhook"),
	"TriggerConfigs":   arrayType(schemaObject(webhookTriggerConfigSchema)),
	"CreatedAt":        timestampField("When the webhook was created"),
	"LastUpdatedAt":    timestampField("When the webhook was last updated"),
	"ArchivedAt":       timestampField("When the webhook was archived"),
}

var webhookTriggerEventSchema = map[string]any{
	"ID":            stringField("The ID of the trigger event"),
	"Name":          stringField("The trigger event name"),
	"Description":   stringField("The trigger event description"),
	"CreatedAt":     timestampField("When the trigger event was created"),
	"LastUpdatedAt": timestampField("When the trigger event was last updated"),
	"ArchivedAt":    timestampField("When the trigger event was archived"),
}

var getWebhookTool = &mcp.Tool{
	Name:        "GetWebhook",
	Description: "Get a webhook by its ID",
	InputSchema: schemaObject(map[string]any{
		"WebhookID": stringField("The ID of the webhook to get"),
	}),
	OutputSchema: schemaObject(webhookSchema),
}

type GetWebhookInvocation struct {
	WebhookID string `jsonschema:"description=The webhook ID"`
}

func (h *mcpToolManager) GetWebhook() mcp.ToolHandlerFor[*GetWebhookInvocation, *webhooks.Webhook] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *GetWebhookInvocation) (*mcp.CallToolResult, *webhooks.Webhook, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		result, err := c.GetWebhook(ctx, &webhooksgrpc.GetWebhookRequest{
			WebhookId: x.WebhookID,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, webhooksconverters.ConvertGRPCWebhookToWebhook(result.Result), nil
	}
}

var getWebhooksTool = &mcp.Tool{
	Name:        "GetWebhooks",
	Description: "Get webhooks with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(webhookSchema)),
	}),
}

type (
	GetWebhooksInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetWebhooksResult struct {
		Results []*webhooks.Webhook
	}
)

func (h *mcpToolManager) GetWebhooks() mcp.ToolHandlerFor[*GetWebhooksInvocation, *GetWebhooksResult] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *GetWebhooksInvocation) (*mcp.CallToolResult, *GetWebhooksResult, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		results, err := c.GetWebhooks(ctx, &webhooksgrpc.GetWebhooksRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetWebhooksResult{}
		for _, w := range results.Results {
			out.Results = append(out.Results, webhooksconverters.ConvertGRPCWebhookToWebhook(w))
		}
		return nil, out, nil
	}
}

var createWebhookTool = &mcp.Tool{
	Name:        "CreateWebhook",
	Description: "Create a new webhook",
	InputSchema: schemaObject(map[string]any{
		"Name":        stringField("The webhook name"),
		"URL":         stringField("The webhook URL"),
		"Method":      stringField("The HTTP method (GET, POST, PUT, PATCH, DELETE)"),
		"ContentType": stringField("The content type (application/json or application/xml)"),
		"Events": arrayType(objectType(map[string]any{
			"ID":          stringField("Optional: ID of an existing trigger event to reference"),
			"Name":        stringField("Name for a new trigger event (used if ID is not provided)"),
			"Description": stringField("Description for a new trigger event"),
		})),
	}),
	OutputSchema: schemaObject(webhookSchema),
}

type CreateWebhookInvocation struct {
	*webhooks.WebhookCreationRequestInput
}

func (h *mcpToolManager) CreateWebhook() mcp.ToolHandlerFor[*CreateWebhookInvocation, *webhooks.Webhook] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *CreateWebhookInvocation) (*mcp.CallToolResult, *webhooks.Webhook, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		result, err := c.CreateWebhook(ctx, &webhooksgrpc.CreateWebhookRequest{
			Input: webhooksconverters.ConvertWebhookCreationRequestInputToGRPCWebhookCreationRequestInput(x.WebhookCreationRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, webhooksconverters.ConvertGRPCWebhookToWebhook(result.Created), nil
	}
}

var archiveWebhookTool = &mcp.Tool{
	Name:        "ArchiveWebhook",
	Description: "Archive (soft-delete) a webhook",
	InputSchema: schemaObject(map[string]any{
		"WebhookID": stringField("The ID of the webhook to archive"),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Success": boolField("Whether the archive was successful"),
	}),
}

type ArchiveWebhookInvocation struct {
	WebhookID string `jsonschema:"required,description=The ID of the webhook to archive"`
}

func (h *mcpToolManager) ArchiveWebhook() mcp.ToolHandlerFor[*ArchiveWebhookInvocation, *boolResult] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *ArchiveWebhookInvocation) (*mcp.CallToolResult, *boolResult, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		_, err = c.ArchiveWebhook(ctx, &webhooksgrpc.ArchiveWebhookRequest{
			WebhookId: x.WebhookID,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, &boolResult{Success: true}, nil
	}
}

var getWebhookTriggerEventsTool = &mcp.Tool{
	Name:        "GetWebhookTriggerEvents",
	Description: "Get webhook trigger events (the catalog of available event types)",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(webhookTriggerEventSchema)),
	}),
}

type (
	GetWebhookTriggerEventsInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetWebhookTriggerEventsResult struct {
		Results []*webhooks.WebhookTriggerEvent
	}
)

func (h *mcpToolManager) GetWebhookTriggerEvents() mcp.ToolHandlerFor[*GetWebhookTriggerEventsInvocation, *GetWebhookTriggerEventsResult] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *GetWebhookTriggerEventsInvocation) (*mcp.CallToolResult, *GetWebhookTriggerEventsResult, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		results, err := c.GetWebhookTriggerEvents(ctx, &webhooksgrpc.GetWebhookTriggerEventsRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetWebhookTriggerEventsResult{}
		for _, e := range results.Results {
			out.Results = append(out.Results, &webhooks.WebhookTriggerEvent{
				CreatedAt:     grpcconverters.ConvertPBTimestampToTime(e.CreatedAt),
				LastUpdatedAt: grpcconverters.ConvertPBTimestampToTimePointer(e.LastUpdatedAt),
				ArchivedAt:    grpcconverters.ConvertPBTimestampToTimePointer(e.ArchivedAt),
				ID:            e.Id,
				Name:          e.Name,
				Description:   e.Description,
			})
		}
		return nil, out, nil
	}
}

var createWebhookTriggerEventTool = &mcp.Tool{
	Name:        "CreateWebhookTriggerEvent",
	Description: "Create a new webhook trigger event in the catalog",
	InputSchema: schemaObject(map[string]any{
		"Name":        stringField("The trigger event name"),
		"Description": stringField("The trigger event description"),
	}),
	OutputSchema: schemaObject(webhookTriggerEventSchema),
}

type CreateWebhookTriggerEventInvocation struct {
	Name        string `jsonschema:"description=The trigger event name"`
	Description string `jsonschema:"description=The trigger event description"`
}

func (h *mcpToolManager) CreateWebhookTriggerEvent() mcp.ToolHandlerFor[*CreateWebhookTriggerEventInvocation, *webhooks.WebhookTriggerEvent] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *CreateWebhookTriggerEventInvocation) (*mcp.CallToolResult, *webhooks.WebhookTriggerEvent, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		result, err := c.CreateWebhookTriggerEvent(ctx, &webhooksgrpc.CreateWebhookTriggerEventRequest{
			Input: &webhooksgrpc.WebhookTriggerEventCreationRequestInput{
				Name:        x.Name,
				Description: x.Description,
			},
		})
		if err != nil {
			return nil, nil, err
		}
		e := result.Created
		return nil, &webhooks.WebhookTriggerEvent{
			CreatedAt:     grpcconverters.ConvertPBTimestampToTime(e.CreatedAt),
			LastUpdatedAt: grpcconverters.ConvertPBTimestampToTimePointer(e.LastUpdatedAt),
			ArchivedAt:    grpcconverters.ConvertPBTimestampToTimePointer(e.ArchivedAt),
			ID:            e.Id,
			Name:          e.Name,
			Description:   e.Description,
		}, nil
	}
}

var addWebhookTriggerConfigTool = &mcp.Tool{
	Name:        "AddWebhookTriggerConfig",
	Description: "Add a trigger event configuration to a webhook",
	InputSchema: schemaObject(map[string]any{
		"WebhookID":      stringField("The ID of the webhook"),
		"TriggerEventID": stringField("The ID of the trigger event to subscribe to"),
	}),
	OutputSchema: schemaObject(webhookTriggerConfigSchema),
}

type AddWebhookTriggerConfigInvocation struct {
	WebhookID      string `jsonschema:"required,description=The ID of the webhook"`
	TriggerEventID string `jsonschema:"required,description=The ID of the trigger event"`
}

func (h *mcpToolManager) AddWebhookTriggerConfig() mcp.ToolHandlerFor[*AddWebhookTriggerConfigInvocation, *webhooks.WebhookTriggerConfig] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *AddWebhookTriggerConfigInvocation) (*mcp.CallToolResult, *webhooks.WebhookTriggerConfig, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		result, err := c.AddWebhookTriggerConfig(ctx, &webhooksgrpc.AddWebhookTriggerConfigRequest{
			WebhookId: x.WebhookID,
			Input: &webhooksgrpc.WebhookTriggerConfigCreationRequestInput{
				BelongsToWebhook: x.WebhookID,
				TriggerEventId:   x.TriggerEventID,
			},
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, webhooksconverters.ConvertGRPCWebhookTriggerConfigToWebhookTriggerConfig(result.Created), nil
	}
}

var archiveWebhookTriggerConfigTool = &mcp.Tool{
	Name:        "ArchiveWebhookTriggerConfig",
	Description: "Archive (remove) a trigger event configuration from a webhook",
	InputSchema: schemaObject(map[string]any{
		"WebhookID":              stringField("The ID of the webhook"),
		"WebhookTriggerConfigID": stringField("The ID of the trigger config to archive"),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Success": boolField("Whether the archive was successful"),
	}),
}

type ArchiveWebhookTriggerConfigInvocation struct {
	WebhookID              string `jsonschema:"required,description=The ID of the webhook"`
	WebhookTriggerConfigID string `jsonschema:"required,description=The ID of the trigger config to archive"`
}

func (h *mcpToolManager) ArchiveWebhookTriggerConfig() mcp.ToolHandlerFor[*ArchiveWebhookTriggerConfigInvocation, *boolResult] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *ArchiveWebhookTriggerConfigInvocation) (*mcp.CallToolResult, *boolResult, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		_, err = c.ArchiveWebhookTriggerConfig(ctx, &webhooksgrpc.ArchiveWebhookTriggerConfigRequest{
			WebhookId:              x.WebhookID,
			WebhookTriggerConfigId: x.WebhookTriggerConfigID,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, &boolResult{Success: true}, nil
	}
}
