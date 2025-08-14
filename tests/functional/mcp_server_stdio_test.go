package functional_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

// MCPRequest represents a JSON-RPC request to the MCP server
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      int         `json:"id"`
}

// MCPResponse represents a JSON-RPC response from the MCP server
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	ID      int         `json:"id"`
}

// MCPServer represents a running MCP server process
type MCPServer struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	reader *bufio.Scanner
}

var _ = Describe("MCP Server Stdio Functional Tests", func() {
	var (
		mcpServer     *MCPServer
		testNamespace string
	)

	BeforeEach(func() {
		// Create test namespace with nanoseconds for uniqueness
		testNamespace = fmt.Sprintf("mcp-test-%d", time.Now().UnixNano())
		By(fmt.Sprintf("Creating test namespace: %s", testNamespace))
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: testNamespace,
			},
		}
		_, err := virtClient.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		// Stop MCP server if running
		if mcpServer != nil {
			mcpServer.Stop()
			mcpServer = nil
		}

		// Clean up test namespace
		if testNamespace != "" {
			By(fmt.Sprintf("Cleaning up test namespace: %s", testNamespace))
			err := virtClient.CoreV1().Namespaces().Delete(context.Background(), testNamespace, metav1.DeleteOptions{})
			if err != nil {
				GinkgoLogr.Error(err, "Failed to delete test namespace", "namespace", testNamespace)
			}
			// Wait a moment for namespace deletion to complete
			time.Sleep(100 * time.Millisecond)
		}
	})

	Context("MCP Server Communication", func() {
		It("should start MCP server and initialize successfully", func() {
			By("Starting MCP server")
			var err error
			mcpServer, err = StartMCPServer()
			Expect(err).NotTo(HaveOccurred(), "Should be able to start MCP server")
			Expect(mcpServer).NotTo(BeNil())

			By("Sending initialize request")
			initRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "initialize",
				Params: map[string]interface{}{
					"capabilities": map[string]interface{}{},
					"clientInfo": map[string]interface{}{
						"name":    "test-client",
						"version": "1.0.0",
					},
				},
				ID: 1,
			}

			response, err := mcpServer.SendRequest(initRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive initialize response")
			Expect(response.Error).To(BeNil(), "Initialize should not return error")
			Expect(response.Result).NotTo(BeNil(), "Initialize should return result")
		})

		It("should list tools successfully", func() {
			By("Starting and initializing MCP server")
			var err error
			mcpServer, err = StartMCPServer()
			Expect(err).NotTo(HaveOccurred())

			// Initialize first
			initRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "initialize",
				Params: map[string]interface{}{
					"capabilities": map[string]interface{}{},
					"clientInfo": map[string]interface{}{
						"name":    "test-client",
						"version": "1.0.0",
					},
				},
				ID: 1,
			}
			_, err = mcpServer.SendRequest(initRequest)
			Expect(err).NotTo(HaveOccurred())

			By("Requesting tools list")
			toolsRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "tools/list",
				Params:  map[string]interface{}{},
				ID:      2,
			}

			response, err := mcpServer.SendRequest(toolsRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive tools list response")
			Expect(response.Error).To(BeNil(), "Tools list should not return error")
			Expect(response.Result).NotTo(BeNil(), "Tools list should return result")

			// Verify some expected tools are present
			result, ok := response.Result.(map[string]interface{})
			Expect(ok).To(BeTrue(), "Result should be a map")
			tools, ok := result["tools"].([]interface{})
			Expect(ok).To(BeTrue(), "Result should contain tools array")
			Expect(len(tools)).To(BeNumerically(">", 0), "Should have at least one tool")

			// Check for expected tool names
			toolNames := make([]string, 0, len(tools))
			for _, tool := range tools {
				toolMap, ok := tool.(map[string]interface{})
				Expect(ok).To(BeTrue(), "Each tool should be a map")
				name, ok := toolMap["name"].(string)
				Expect(ok).To(BeTrue(), "Each tool should have a name")
				toolNames = append(toolNames, name)
			}
			Expect(toolNames).To(ContainElement("list_vms"), "Should contain list_vms tool")
			Expect(toolNames).To(ContainElement("start_vm"), "Should contain start_vm tool")
			Expect(toolNames).To(ContainElement("stop_vm"), "Should contain stop_vm tool")
		})

		It("should call list_vms tool successfully", func() {
			By("Starting and initializing MCP server")
			var err error
			mcpServer, err = StartMCPServer()
			Expect(err).NotTo(HaveOccurred())

			// Initialize first
			initRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "initialize",
				Params: map[string]interface{}{
					"capabilities": map[string]interface{}{},
					"clientInfo": map[string]interface{}{
						"name":    "test-client",
						"version": "1.0.0",
					},
				},
				ID: 1,
			}
			_, err = mcpServer.SendRequest(initRequest)
			Expect(err).NotTo(HaveOccurred())

			By("Calling list_vms tool")
			listVMsRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name": "list_vms",
					"arguments": map[string]interface{}{
						"namespace": testNamespace,
					},
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(listVMsRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive list_vms response")
			Expect(response.Error).To(BeNil(), "list_vms should not return error")
			Expect(response.Result).NotTo(BeNil(), "list_vms should return result")

			// Verify the result structure
			result, ok := response.Result.(map[string]interface{})
			Expect(ok).To(BeTrue(), "Result should be a map")
			content, ok := result["content"].([]interface{})
			Expect(ok).To(BeTrue(), "Result should contain content array")
			
			// Should be empty initially (no VMs in test namespace)
			Expect(len(content)).To(Equal(1), "Should have one content item (text response)")
			contentItem, ok := content[0].(map[string]interface{})
			Expect(ok).To(BeTrue(), "Content item should be a map")
			text, ok := contentItem["text"].(string)
			Expect(ok).To(BeTrue(), "Content should have text")
			Expect(strings.TrimSpace(text)).To(Equal(""), "Should return empty string when no VMs found")
		})
	})

	Context("KubeVirt Integration with VMs", func() {
		var testVM *kubevirtv1.VirtualMachine

		BeforeEach(func() {
			By("Creating a test VM for MCP server to discover")
			testVM = &kubevirtv1.VirtualMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("test-vm-for-mcp-%d", time.Now().UnixNano()),
					Namespace: testNamespace,
				},
				Spec: kubevirtv1.VirtualMachineSpec{
					Running: &[]bool{false}[0],
					Template: &kubevirtv1.VirtualMachineInstanceTemplateSpec{
						Spec: kubevirtv1.VirtualMachineInstanceSpec{
							Domain: kubevirtv1.DomainSpec{
								Resources: kubevirtv1.ResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceMemory: resource.MustParse("128Mi"),
									},
								},
								Devices: kubevirtv1.Devices{
									Disks: []kubevirtv1.Disk{
										{
											Name: "containerdisk",
											DiskDevice: kubevirtv1.DiskDevice{
												Disk: &kubevirtv1.DiskTarget{
													Bus: "virtio",
												},
											},
										},
									},
								},
							},
							Volumes: []kubevirtv1.Volume{
								{
									Name: "containerdisk",
									VolumeSource: kubevirtv1.VolumeSource{
										ContainerDisk: &kubevirtv1.ContainerDiskSource{
											Image: "quay.io/kubevirt/cirros-container-disk-demo",
										},
									},
								},
							},
						},
					},
				},
			}

			createdVM, err := virtClient.VirtualMachine(testNamespace).Create(context.Background(), testVM, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
			testVM = createdVM
		})

		AfterEach(func() {
			if testVM != nil {
				By("Cleaning up test VM")
				err := virtClient.VirtualMachine(testNamespace).Delete(context.Background(), testVM.Name, metav1.DeleteOptions{})
				if err != nil {
					GinkgoLogr.Error(err, "Failed to delete test VM", "vm", testVM.Name)
				}
			}
		})

		It("should discover the test VM via MCP server", func() {
			By("Starting and initializing MCP server")
			var err error
			mcpServer, err = StartMCPServer()
			Expect(err).NotTo(HaveOccurred())

			// Initialize first
			initRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "initialize",
				Params: map[string]interface{}{
					"capabilities": map[string]interface{}{},
					"clientInfo": map[string]interface{}{
						"name":    "test-client",
						"version": "1.0.0",
					},
				},
				ID: 1,
			}
			_, err = mcpServer.SendRequest(initRequest)
			Expect(err).NotTo(HaveOccurred())

			By("Listing VMs through MCP server")
			listVMsRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name": "list_vms",
					"arguments": map[string]interface{}{
						"namespace": testNamespace,
					},
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(listVMsRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive list_vms response")
			Expect(response.Error).To(BeNil(), "list_vms should not return error")

			// Verify the VM is found
			result, ok := response.Result.(map[string]interface{})
			Expect(ok).To(BeTrue(), "Result should be a map")
			content, ok := result["content"].([]interface{})
			Expect(ok).To(BeTrue(), "Result should contain content array")
			
			Expect(len(content)).To(Equal(1), "Should have one content item")
			contentItem, ok := content[0].(map[string]interface{})
			Expect(ok).To(BeTrue(), "Content item should be a map")
			text, ok := contentItem["text"].(string)
			Expect(ok).To(BeTrue(), "Content should have text")
			Expect(text).To(ContainSubstring(testVM.Name), "Should contain the test VM name")
		})
	})
})

// StartMCPServer starts a new MCP server process and returns a handle to communicate with it
func StartMCPServer() (*MCPServer, error) {
	// Get the project root directory
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		return nil, fmt.Errorf("failed to get project root: %w", err)
	}

	// Get kubeconfig from kubevirtci
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = filepath.Join(projectRoot, "_kubevirtci/_ci-configs/k8s-1.33/.kubeconfig")
	}

	// Start the MCP server process
	binaryPath := filepath.Join(projectRoot, "kubevirt-mcp-server")
	cmd := exec.Command(binaryPath)
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", kubeconfigPath))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start MCP server: %w", err)
	}

	return &MCPServer{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		reader: bufio.NewScanner(stdout),
	}, nil
}

// SendRequest sends a JSON-RPC request to the MCP server and waits for a response
func (s *MCPServer) SendRequest(request MCPRequest) (*MCPResponse, error) {
	// Serialize the request
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send the request
	_, err = s.stdin.Write(append(requestBytes, '\n'))
	if err != nil {
		return nil, fmt.Errorf("failed to write request: %w", err)
	}

	// Read the response
	if !s.reader.Scan() {
		if err := s.reader.Err(); err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}
		return nil, fmt.Errorf("no response received")
	}

	responseText := s.reader.Text()
	if strings.TrimSpace(responseText) == "" {
		return nil, fmt.Errorf("empty response received")
	}

	// Parse the response
	var response MCPResponse
	if err := json.Unmarshal([]byte(responseText), &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// Stop terminates the MCP server process
func (s *MCPServer) Stop() error {
	if s.stdin != nil {
		s.stdin.Close()
	}
	if s.stdout != nil {
		s.stdout.Close()
	}
	if s.stderr != nil {
		s.stderr.Close()
	}
	if s.cmd != nil && s.cmd.Process != nil {
		return s.cmd.Process.Kill()
	}
	return nil
}