package main

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/issuereports"
	grpcconverters "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/converters"
	issuereportsgrpc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/generated/services/issue_reports"
	issuereportsconverters "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/issuereports/grpc/converters"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// boolResult is a shared result type for operations that don't return a domain object.
type boolResult struct {
	Success bool
}

var issueReportSchema = map[string]any{
	"ID":               stringField("The ID of the issue report"),
	"IssueType":        stringField("The type of issue"),
	"Details":          stringField("Detailed description of the issue"),
	"RelevantTable":    stringField("The database table the issue is related to, if any"),
	"RelevantRecordID": stringField("The ID of the record the issue is related to, if any"),
	"CreatedByUser":    stringField("The ID of the user who created the report"),
	"BelongsToAccount": stringField("The ID of the account this report belongs to"),
	"CreatedAt":        timestampField("When the report was created"),
	"LastUpdatedAt":    timestampField("When the report was last updated"),
	"ArchivedAt":       timestampField("When the report was archived"),
}

var getIssueReportTool = &mcp.Tool{
	Name:        "GetIssueReport",
	Description: "Get an issue report by its ID",
	InputSchema: schemaObject(map[string]any{
		"IssueReportID": stringField("The ID of the issue report to get"),
	}),
	OutputSchema: schemaObject(issueReportSchema),
}

type GetIssueReportInvocation struct {
	IssueReportID string `jsonschema:"description=The issue report ID"`
}

func (h *mcpToolManager) GetIssueReport() mcp.ToolHandlerFor[*GetIssueReportInvocation, *issuereports.IssueReport] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *GetIssueReportInvocation) (*mcp.CallToolResult, *issuereports.IssueReport, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		result, err := c.GetIssueReport(ctx, &issuereportsgrpc.GetIssueReportRequest{
			IssueReportId: x.IssueReportID,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, issuereportsconverters.ConvertGRPCIssueReportToIssueReport(result.Result), nil
	}
}

var getIssueReportsTool = &mcp.Tool{
	Name:        "GetIssueReports",
	Description: "Get issue reports with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(issueReportSchema)),
	}),
}

type (
	GetIssueReportsInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetIssueReportsResult struct {
		Results []*issuereports.IssueReport
	}
)

func (h *mcpToolManager) GetIssueReports() mcp.ToolHandlerFor[*GetIssueReportsInvocation, *GetIssueReportsResult] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *GetIssueReportsInvocation) (*mcp.CallToolResult, *GetIssueReportsResult, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		results, err := c.GetIssueReports(ctx, &issuereportsgrpc.GetIssueReportsRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetIssueReportsResult{}
		for _, r := range results.Results {
			out.Results = append(out.Results, issuereportsconverters.ConvertGRPCIssueReportToIssueReport(r))
		}
		return nil, out, nil
	}
}

var getIssueReportsForAccountTool = &mcp.Tool{
	Name:        "GetIssueReportsForAccount",
	Description: "Get issue reports belonging to a specific account",
	InputSchema: schemaObject(map[string]any{
		"AccountID": stringField("The account ID"),
		"Filter":    queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(issueReportSchema)),
	}),
}

type (
	GetIssueReportsForAccountInvocation struct {
		Filter    *filtering.QueryFilter
		AccountID string `jsonschema:"description=The account ID"`
	}

	GetIssueReportsForAccountResult struct {
		Results []*issuereports.IssueReport
	}
)

func (h *mcpToolManager) GetIssueReportsForAccount() mcp.ToolHandlerFor[*GetIssueReportsForAccountInvocation, *GetIssueReportsForAccountResult] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *GetIssueReportsForAccountInvocation) (*mcp.CallToolResult, *GetIssueReportsForAccountResult, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		results, err := c.GetIssueReportsForAccount(ctx, &issuereportsgrpc.GetIssueReportsForAccountRequest{
			AccountId: x.AccountID,
			Filter:    grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetIssueReportsForAccountResult{}
		for _, r := range results.Results {
			out.Results = append(out.Results, issuereportsconverters.ConvertGRPCIssueReportToIssueReport(r))
		}
		return nil, out, nil
	}
}

var createIssueReportTool = &mcp.Tool{
	Name:        "CreateIssueReport",
	Description: "Create a new issue report",
	InputSchema: schemaObject(map[string]any{
		"IssueType":        stringField("The type of issue"),
		"Details":          stringField("Detailed description of the issue"),
		"RelevantTable":    stringField("The database table the issue is related to, if any"),
		"RelevantRecordID": stringField("The ID of the record the issue is related to, if any"),
	}),
	OutputSchema: schemaObject(issueReportSchema),
}

type CreateIssueReportInvocation struct {
	*issuereports.IssueReportCreationRequestInput
}

func (h *mcpToolManager) CreateIssueReport() mcp.ToolHandlerFor[*CreateIssueReportInvocation, *issuereports.IssueReport] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *CreateIssueReportInvocation) (*mcp.CallToolResult, *issuereports.IssueReport, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		result, err := c.CreateIssueReport(ctx, &issuereportsgrpc.CreateIssueReportRequest{
			Input: issuereportsconverters.ConvertIssueReportCreationRequestInputToGRPCIssueReportCreationRequestInput(x.IssueReportCreationRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, issuereportsconverters.ConvertGRPCIssueReportToIssueReport(result.Created), nil
	}
}

var updateIssueReportTool = &mcp.Tool{
	Name:        "UpdateIssueReport",
	Description: "Update an existing issue report",
	InputSchema: schemaObject(map[string]any{
		"IssueReportID":    stringField("The ID of the issue report to update"),
		"IssueType":        stringField("New issue type (leave empty to skip)"),
		"Details":          stringField("New details (leave empty to skip)"),
		"RelevantTable":    stringField("New relevant table (leave empty to skip)"),
		"RelevantRecordID": stringField("New relevant record ID (leave empty to skip)"),
	}),
	OutputSchema: schemaObject(issueReportSchema),
}

type UpdateIssueReportInvocation struct {
	IssueReportID    string `jsonschema:"required,description=The ID of the issue report to update"`
	IssueType        string `jsonschema:"description=New issue type"`
	Details          string `jsonschema:"description=New details"`
	RelevantTable    string `jsonschema:"description=New relevant table"`
	RelevantRecordID string `jsonschema:"description=New relevant record ID"`
}

func (h *mcpToolManager) UpdateIssueReport() mcp.ToolHandlerFor[*UpdateIssueReportInvocation, *issuereports.IssueReport] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *UpdateIssueReportInvocation) (*mcp.CallToolResult, *issuereports.IssueReport, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		input := &issuereportsgrpc.IssueReportUpdateRequestInput{}
		if x.IssueType != "" {
			input.IssueType = &x.IssueType
		}
		if x.Details != "" {
			input.Details = &x.Details
		}
		if x.RelevantTable != "" {
			input.RelevantTable = &x.RelevantTable
		}
		if x.RelevantRecordID != "" {
			input.RelevantRecordId = &x.RelevantRecordID
		}

		result, err := c.UpdateIssueReport(ctx, &issuereportsgrpc.UpdateIssueReportRequest{
			IssueReportId: x.IssueReportID,
			Input:         input,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, issuereportsconverters.ConvertGRPCIssueReportToIssueReport(result.Updated), nil
	}
}

var archiveIssueReportTool = &mcp.Tool{
	Name:        "ArchiveIssueReport",
	Description: "Archive (soft-delete) an issue report",
	InputSchema: schemaObject(map[string]any{
		"IssueReportID": stringField("The ID of the issue report to archive"),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Success": boolField("Whether the archive was successful"),
	}),
}

type ArchiveIssueReportInvocation struct {
	IssueReportID string `jsonschema:"required,description=The ID of the issue report to archive"`
}

func (h *mcpToolManager) ArchiveIssueReport() mcp.ToolHandlerFor[*ArchiveIssueReportInvocation, *boolResult] {
	return func(ctx context.Context, req *mcp.CallToolRequest, x *ArchiveIssueReportInvocation) (*mcp.CallToolResult, *boolResult, error) {
		c, err := h.clientFromRequest(req)
		if err != nil {
			return nil, nil, err
		}

		_, err = c.ArchiveIssueReport(ctx, &issuereportsgrpc.ArchiveIssueReportRequest{
			IssueReportId: x.IssueReportID,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, &boolResult{Success: true}, nil
	}
}
