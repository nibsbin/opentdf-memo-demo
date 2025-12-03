// OpenTDF MCP Server
// Run as a standalone Model Context Protocol server for OpenTDF operations
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/opentdf/platform/protocol/go/policy/attributes"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/opentdf/platform/sdk"
)

// Tool input/output schemas
// EncryptToolInput defines the input for the encrypt tool
type EncryptToolInput struct {
	Input        string   `json:"input,omitempty" jsonschema:"Path to plaintext file to encrypt (mutually exclusive with data)"`
	Data         string   `json:"data,omitempty" jsonschema:"Literal data to encrypt (mutually exclusive with input)"`
	Attributes   []string `json:"attributes" jsonschema:"Data attributes (FQNs) to apply during encryption"`
	Output       string   `json:"output,omitempty" jsonschema:"Output file path (optional returns base64 if not specified)"`
	ClientID     string   `json:"clientId,omitempty" jsonschema:"OAuth client ID for OpenTDF platform authentication"`
	ClientSecret string   `json:"clientSecret,omitempty" jsonschema:"OAuth client secret for OpenTDF platform authentication"`
}

type EncryptToolOutput struct {
	Success    bool   `json:"success"`
	OutputFile string `json:"outputFile"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
}

// DecryptToolInput defines the input for the decrypt tool
type DecryptToolInput struct {
	Input        string `json:"input" jsonschema:"Path to encrypted file or base64 encoded data"`
	Output       string `json:"output,omitempty" jsonschema:"Output file path (optional returns plaintext if not specified)"`
	ClientID     string `json:"clientId,omitempty" jsonschema:"OAuth client ID for OpenTDF platform authentication"`
	ClientSecret string `json:"clientSecret,omitempty" jsonschema:"OAuth client secret for OpenTDF platform authentication"`
}

type DecryptToolOutput struct {
	Success      bool   `json:"success"`
	DecryptedData string `json:"decryptedData,omitempty"`
	Error        string `json:"error,omitempty"`
}

type ListAttributesToolInput struct {
	Namespace    string `json:"namespace,omitempty" jsonschema:"Filter by namespace (e.g. https://example.com)"`
	Verbose      bool   `json:"verbose,omitempty" jsonschema:"Show detailed attribute information"`
	ClientID     string `json:"clientId,omitempty" jsonschema:"OAuth client ID for OpenTDF platform authentication"`
	ClientSecret string `json:"clientSecret,omitempty" jsonschema:"OAuth client secret for OpenTDF platform authentication"`
}

type ListAttributesToolOutput struct {
	Success    bool                   `json:"success"`
	Attributes []AttributeInfo        `json:"attributes,omitempty"`
	Error      string                 `json:"error,omitempty"`
}

type AttributeInfo struct {
	Namespace string   `json:"namespace"`
	Name      string   `json:"name"`
	Values    []string `json:"values,omitempty"`
	FQN       string   `json:"fqn"`
}

// JWT Claims structure for agent authentication
type JWTClaims struct {
	Sub         string   `json:"sub"`         // Subject (agent ID)
	Iss         string   `json:"iss"`         // Issuer
	Aud         string   `json:"aud"`         // Audience
	Iat         int64    `json:"iat"`         // Issued at
	Exp         int64    `json:"exp"`         // Expiration
	AgentName   string   `json:"agent_name"`  // Agent display name
	Permissions []string `json:"permissions"` // Granted permissions
}

// parseJWT parses a JWT token and extracts claims (mock implementation - no signature verification)
func parseJWT(token string) (*JWTClaims, error) {
	if token == "" {
		return nil, fmt.Errorf("no JWT token provided")
	}

	// Split the JWT into parts
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	// Decode the payload (second part)
	payload := parts[1]

	// Add padding if needed for base64 decoding
	if l := len(payload) % 4; l > 0 {
		payload += strings.Repeat("=", 4-l)
	}

	// Decode base64
	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	// Parse JSON
	var claims JWTClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse JWT claims: %w", err)
	}

	return &claims, nil
}

// validateJWT performs basic validation on JWT claims (mock implementation)
func validateJWT(claims *JWTClaims) error {
	if claims == nil {
		return fmt.Errorf("no claims provided")
	}

	// Check audience
	if claims.Aud != "opentdf-mcp" {
		return fmt.Errorf("invalid audience: expected 'opentdf-mcp', got '%s'", claims.Aud)
	}

	// Check expiration
	now := time.Now().Unix()
	if claims.Exp > 0 && claims.Exp < now {
		return fmt.Errorf("token expired")
	}

	// Check issued at
	if claims.Iat > 0 && claims.Iat > now {
		return fmt.Errorf("token used before issued")
	}

	return nil
}

// logAgentAuthentication logs the agent authentication details
func logAgentAuthentication() {
	token := getAgentJWT()
	if token == "" {
		log.Println("WARNING: No agent JWT token found. Running without agent authentication.")
		log.Println("In production, this would require OAuth-issued JWT after user consent.")
		return
	}

	claims, err := parseJWT(token)
	if err != nil {
		log.Printf("WARNING: Failed to parse agent JWT: %v\n", err)
		return
	}

	if err := validateJWT(claims); err != nil {
		log.Printf("WARNING: JWT validation failed: %v\n", err)
		return
	}

	log.Println("=== Agent Authentication ===")
	log.Printf("Agent ID: %s\n", claims.Sub)
	log.Printf("Agent Name: %s\n", claims.AgentName)
	log.Printf("Issuer: %s\n", claims.Iss)
	log.Printf("Permissions: %v\n", claims.Permissions)
	log.Printf("Expires: %s\n", time.Unix(claims.Exp, 0).Format(time.RFC3339))
	log.Println("NOTE: This is a MOCK JWT for demo purposes.")
	log.Println("In production, JWT would be issued by OAuth after user consent.")
	log.Println("===========================")
}

// getSDKClientMCP creates an authenticated OpenTDF SDK client for MCP
// If clientID and clientSecret are provided, they take precedence over environment variables
func getSDKClientMCP(clientID, clientSecret string) (*sdk.SDK, error) {
	platformEndpoint := getPlatformEndpoint()

	// Use provided credentials if available, otherwise fall back to config
	if clientID == "" {
		clientID = getClientID()
	}
	if clientSecret == "" {
		clientSecret = getClientSecret()
	}

	var opts []sdk.Option
	if clientID != "" && clientSecret != "" {
		opts = append(opts, sdk.WithClientCredentials(clientID, clientSecret, nil))
	} else {
		opts = append(opts, sdk.WithInsecurePlaintextConn())
	}

	client, err := sdk.New(platformEndpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create SDK client: %w", err)
	}

	return client, nil
}

// MCPEncrypt encrypts data with the given attributes
func MCPEncrypt(ctx context.Context, req *mcp.CallToolRequest, input EncryptToolInput) (*mcp.CallToolResult, EncryptToolOutput, error) {
	client, err := getSDKClientMCP(input.ClientID, input.ClientSecret)
	if err != nil {
		return nil, EncryptToolOutput{Success: false, Error: err.Error()}, nil
	}
	defer client.Close()

	// Validate input: either Input or Data must be provided, but not both
	if input.Input != "" && input.Data != "" {
		return nil, EncryptToolOutput{Success: false, Error: "cannot specify both 'input' and 'data' parameters"}, nil
	}
	if input.Input == "" && input.Data == "" {
		return nil, EncryptToolOutput{Success: false, Error: "must specify either 'input' (file path) or 'data' (literal data)"}, nil
	}

	// Get the data to encrypt
	var dataToEncrypt string
	if input.Input != "" {
		// Read file contents
		fileData, err := os.ReadFile(input.Input)
		if err != nil {
			return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to read input file: %v", err)}, nil
		}
		dataToEncrypt = string(fileData)
	} else {
		// Use literal data
		dataToEncrypt = input.Data
	}

	// Always use nanoTDF format
	outputFile := input.Output
	if outputFile == "" {
		outputFile = "encrypted.ntdf"
	}

	// Encrypt
	reader := strings.NewReader(dataToEncrypt)
	file, err := os.Create(outputFile)
	if err != nil {
		return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to create output file: %v", err)}, nil
	}
	defer file.Close()

	baseKasURL := getPlatformEndpoint()
	if !strings.HasPrefix(baseKasURL, "http://") && !strings.HasPrefix(baseKasURL, "https://") {
		baseKasURL = "http://" + baseKasURL
	}

	// Create nanoTDF config
	nanoConfig, err := client.NewNanoTDFConfig()
	if err != nil {
		return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to create nanoTDF config: %v", err)}, nil
	}

	if len(input.Attributes) > 0 {
		if err := nanoConfig.SetAttributes(input.Attributes); err != nil {
			return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to set attributes: %v", err)}, nil
		}
	}

	nanoConfig.EnableECDSAPolicyBinding()

	if err := nanoConfig.SetKasURL(baseKasURL + "/kas"); err != nil {
		return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to set KAS URL: %v", err)}, nil
	}

	if _, err := client.CreateNanoTDF(file, reader, *nanoConfig); err != nil {
		return nil, EncryptToolOutput{Success: false, Error: fmt.Sprintf("failed to encrypt: %v", err)}, nil
	}

	msg := fmt.Sprintf("Successfully encrypted data to %s", outputFile)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg},
		},
	}, EncryptToolOutput{Success: true, OutputFile: outputFile, Message: msg}, nil
}

// MCPDecrypt decrypts a TDF or nanoTDF file
func MCPDecrypt(ctx context.Context, req *mcp.CallToolRequest, input DecryptToolInput) (*mcp.CallToolResult, DecryptToolOutput, error) {
	client, err := getSDKClientMCP(input.ClientID, input.ClientSecret)
	if err != nil {
		return nil, DecryptToolOutput{Success: false, Error: err.Error()}, nil
	}
	defer client.Close()

	file, err := os.Open(input.Input)
	if err != nil {
		return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to open input file: %v", err)}, nil
	}
	defer file.Close()

	// Detect format by reading magic bytes
	var magic [3]byte
	var isNano bool
	n, err := io.ReadFull(file, magic[:])
	switch {
	case err != nil:
		return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to read magic bytes: %v", err)}, nil
	case n < 3:
		return nil, DecryptToolOutput{Success: false, Error: "file too small; no magic number found"}, nil
	case bytes.HasPrefix(magic[:], []byte("L1L")):
		isNano = true
	default:
		isNano = false
	}

	// Reset file position
	if _, err := file.Seek(0, 0); err != nil {
		return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to seek to beginning: %v", err)}, nil
	}

	var output bytes.Buffer

	if isNano {
		if _, err := client.ReadNanoTDF(&output, file); err != nil {
			return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to decrypt nanoTDF: %v", err)}, nil
		}
	} else {
		tdfReader, err := client.LoadTDF(file)
		if err != nil {
			return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to load TDF: %v", err)}, nil
		}
		if _, err := io.Copy(&output, tdfReader); err != nil && err != io.EOF {
			return nil, DecryptToolOutput{Success: false, Error: fmt.Sprintf("failed to decrypt TDF: %v", err)}, nil
		}
	}

	decryptedData := output.String()
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Successfully decrypted:\n%s", decryptedData)},
		},
	}, DecryptToolOutput{Success: true, DecryptedData: decryptedData}, nil
}

// MCPListAttributes lists available attributes
func MCPListAttributes(ctx context.Context, req *mcp.CallToolRequest, input ListAttributesToolInput) (*mcp.CallToolResult, ListAttributesToolOutput, error) {
	client, err := getSDKClientMCP(input.ClientID, input.ClientSecret)
	if err != nil {
		return nil, ListAttributesToolOutput{Success: false, Error: err.Error()}, nil
	}
	defer client.Close()

	// List namespaces
	var nsuris []string
	if input.Namespace == "" {
		listResp, err := client.Namespaces.ListNamespaces(ctx, &namespaces.ListNamespacesRequest{})
		if err != nil {
			return nil, ListAttributesToolOutput{Success: false, Error: fmt.Sprintf("failed to list namespaces: %v", err)}, nil
		}
		for _, n := range listResp.GetNamespaces() {
			nsuris = append(nsuris, n.GetFqn())
		}
	} else {
		nsuris = []string{input.Namespace}
	}

	var attributesList []AttributeInfo
	var textOutput strings.Builder

	for _, ns := range nsuris {
		// Parse namespace to get host
		parts := strings.Split(ns, "//")
		var host string
		if len(parts) > 1 {
			host = strings.Split(parts[1], "/")[0]
		} else {
			host = ns
		}

		lsr, err := client.Attributes.ListAttributes(ctx, &attributes.ListAttributesRequest{
			Namespace: host,
		})
		if err != nil {
			continue
		}

		textOutput.WriteString(fmt.Sprintf("Namespace: %s\n", ns))
		for _, a := range lsr.GetAttributes() {
			values := []string{}
			if input.Verbose {
				for _, v := range a.GetValues() {
					values = append(values, v.GetValue())
				}
			}

			attrInfo := AttributeInfo{
				Namespace: ns,
				Name:      a.GetName(),
				Values:    values,
				FQN:       a.GetFqn(),
			}
			attributesList = append(attributesList, attrInfo)

			if input.Verbose {
				textOutput.WriteString(fmt.Sprintf("  Attribute: %s\n", a.GetFqn()))
				if len(values) > 0 {
					textOutput.WriteString(fmt.Sprintf("  Values: %s\n", strings.Join(values, ", ")))
				}
			} else {
				textOutput.WriteString(fmt.Sprintf("  %s\n", a.GetFqn()))
			}
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: textOutput.String()},
		},
	}, ListAttributesToolOutput{Success: true, Attributes: attributesList}, nil
}

func runMCPServer() error {
	// Log agent authentication details
	logAgentAuthentication()

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "opentdf-mcp",
		Version: "1.0.0",
	}, &mcp.ServerOptions{
		InitializedHandler: func(ctx context.Context, req *mcp.InitializedRequest) {
			log.Println("MCP server initialized")
		},
	})

	// Add encrypt tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "encrypt",
		Description: "Encrypt data using OpenTDF with the specified attributes. Creates a nanoTDF file (.ntdf). Specify either 'input' (file path) or 'data' (literal text).",
	}, MCPEncrypt)

	// Add decrypt tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "decrypt",
		Description: "Decrypt a TDF or nanoTDF file and return the plaintext data. Automatically detects the format.",
	}, MCPDecrypt)

	// Add list attributes tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_attributes",
		Description: "List available data attributes from the OpenTDF platform. Use verbose mode to see attribute values. Filter by namespace if needed.",
	}, MCPListAttributes)

	// Run server over stdio
	log.Println("Starting OpenTDF MCP server on stdio...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

func main() {
	if err := runMCPServer(); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}
