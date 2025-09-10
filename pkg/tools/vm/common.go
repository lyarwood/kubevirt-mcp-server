package vm

import (
	"github.com/mark3labs/mcp-go/mcp"
)

func newToolResultErr(err error) (*mcp.CallToolResult, error) {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: err.Error(),
			},
		},
	}, err
}