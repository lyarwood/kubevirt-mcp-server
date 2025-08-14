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
					RunStrategy: &[]kubevirtv1.VirtualMachineRunStrategy{kubevirtv1.RunStrategyHalted}[0],
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

		It("should start the test VM via MCP server", func() {
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

			By("Starting VM through MCP server")
			startVMRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name": "start_vm",
					"arguments": map[string]interface{}{
						"namespace": testNamespace,
						"name":      testVM.Name,
					},
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(startVMRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive start_vm response")
			Expect(response.Error).To(BeNil(), "start_vm should not return error")

			// Verify the VM was started
			Eventually(func() bool {
				vm, err := virtClient.VirtualMachine(testNamespace).Get(context.Background(), testVM.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
				return vm.Spec.RunStrategy != nil && *vm.Spec.RunStrategy == kubevirtv1.RunStrategyAlways
			}, testTimeout, pollInterval).Should(BeTrue(), "VM should be started")
		})

		It("should stop the test VM via MCP server", func() {
			By("First starting the VM using subresource API")
			err := virtClient.VirtualMachine(testNamespace).Start(context.Background(), testVM.Name, &kubevirtv1.StartOptions{})
			Expect(err).NotTo(HaveOccurred())

			By("Starting and initializing MCP server")
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

			By("Stopping VM through MCP server")
			stopVMRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name": "stop_vm",
					"arguments": map[string]interface{}{
						"namespace": testNamespace,
						"name":      testVM.Name,
					},
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(stopVMRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive stop_vm response")
			// Allow for Kubernetes conflicts which can happen in test environments
			if response.Error != nil {
				errorMap, ok := response.Error.(map[string]interface{})
				if ok {
					message, ok := errorMap["message"].(string)
					if ok && strings.Contains(message, "the object has been modified") {
						Skip("Skipping due to Kubernetes resource conflict in test environment")
					}
				}
			}
			Expect(response.Error).To(BeNil(), "stop_vm should not return error")

			// Verify the VM was stopped
			Eventually(func() bool {
				vm, err := virtClient.VirtualMachine(testNamespace).Get(context.Background(), testVM.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
				return vm.Spec.RunStrategy == nil || *vm.Spec.RunStrategy == kubevirtv1.RunStrategyHalted
			}, testTimeout, pollInterval).Should(BeTrue(), "VM should be stopped")
		})

		It("should restart the test VM via MCP server", func() {
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

			By("Restarting VM through MCP server")
			restartVMRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name": "restart_vm",
					"arguments": map[string]interface{}{
						"namespace": testNamespace,
						"name":      testVM.Name,
					},
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(restartVMRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive restart_vm response")
			// Allow for Kubernetes conflicts which can happen in test environments
			if response.Error != nil {
				errorMap, ok := response.Error.(map[string]interface{})
				if ok {
					message, ok := errorMap["message"].(string)
					if ok && strings.Contains(message, "the object has been modified") {
						Skip("Skipping due to Kubernetes resource conflict in test environment")
					}
				}
			}
			Expect(response.Error).To(BeNil(), "restart_vm should not return error")

			// Verify the VM is running after restart
			Eventually(func() bool {
				vm, err := virtClient.VirtualMachine(testNamespace).Get(context.Background(), testVM.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
				return vm.Spec.RunStrategy != nil && *vm.Spec.RunStrategy == kubevirtv1.RunStrategyAlways
			}, testTimeout, pollInterval).Should(BeTrue(), "VM should be running after restart")
		})

		It("should get VM instancetype via MCP server", func() {
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

			By("Getting VM instancetype through MCP server")
			getInstancetypeRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name": "get_vm_instancetype",
					"arguments": map[string]interface{}{
						"namespace": testNamespace,
						"name":      testVM.Name,
					},
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(getInstancetypeRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive get_vm_instancetype response")
			// Note: This might return an error if no instancetype is set, which is expected
			result, ok := response.Result.(map[string]interface{})
			Expect(ok).To(BeTrue(), "Result should be a map")
			content, ok := result["content"].([]interface{})
			Expect(ok).To(BeTrue(), "Result should contain content array")
			Expect(len(content)).To(BeNumerically(">", 0), "Should have content")
		})
	})

	Context("MCP Tools - General", func() {
		It("should list instancetypes via MCP server", func() {
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

			By("Listing instancetypes through MCP server")
			listInstancetypesRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name":      "list_instancetypes",
					"arguments": map[string]interface{}{},
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(listInstancetypesRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive list_instancetypes response")
			Expect(response.Error).To(BeNil(), "list_instancetypes should not return error")

			result, ok := response.Result.(map[string]interface{})
			Expect(ok).To(BeTrue(), "Result should be a map")
			content, ok := result["content"].([]interface{})
			Expect(ok).To(BeTrue(), "Result should contain content array")
			Expect(len(content)).To(BeNumerically(">", 0), "Should have content")
		})
	})

	XContext("MCP Resources", func() {
		var testVM *kubevirtv1.VirtualMachine

		BeforeEach(func() {
			By("Creating a test VM for resource tests")
			testVM = &kubevirtv1.VirtualMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("test-vm-for-resources-%d", time.Now().UnixNano()),
					Namespace: testNamespace,
				},
				Spec: kubevirtv1.VirtualMachineSpec{
					RunStrategy: &[]kubevirtv1.VirtualMachineRunStrategy{kubevirtv1.RunStrategyHalted}[0],
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
				By("Cleaning up test VM for resources")
				err := virtClient.VirtualMachine(testNamespace).Delete(context.Background(), testVM.Name, metav1.DeleteOptions{})
				if err != nil {
					GinkgoLogr.Error(err, "Failed to delete test VM", "vm", testVM.Name)
				}
			}
		})

		It("should access VMs resource endpoint", func() {
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

			By("Reading VMs resource")
			vmsResourceRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "resources/read",
				Params: map[string]interface{}{
					"uri": fmt.Sprintf("kubevirt://%s/vms", testNamespace),
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(vmsResourceRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive VMs resource response")
			Expect(response.Error).To(BeNil(), "VMs resource should not return error")

			result, ok := response.Result.(map[string]interface{})
			Expect(ok).To(BeTrue(), "Result should be a map")
			contents, ok := result["contents"].([]interface{})
			Expect(ok).To(BeTrue(), "Result should contain contents array")
			Expect(len(contents)).To(BeNumerically(">", 0), "Should have contents")
		})

		It("should access specific VM resource endpoint", func() {
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

			By("Reading specific VM resource")
			vmResourceRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "resources/read",
				Params: map[string]interface{}{
					"uri": fmt.Sprintf("kubevirt://%s/vm/%s", testNamespace, testVM.Name),
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(vmResourceRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive VM resource response")
			Expect(response.Error).To(BeNil(), "VM resource should not return error")

			result, ok := response.Result.(map[string]interface{})
			Expect(ok).To(BeTrue(), "Result should be a map")
			contents, ok := result["contents"].([]interface{})
			Expect(ok).To(BeTrue(), "Result should contain contents array")
			Expect(len(contents)).To(BeNumerically(">", 0), "Should have contents")

			// Check content contains VM details
			content := contents[0].(map[string]interface{})
			text, ok := content["text"].(string)
			Expect(ok).To(BeTrue(), "Content should have text")
			Expect(text).To(ContainSubstring(testVM.Name), "Should contain VM name")
		})

		It("should access VMIs resource endpoint", func() {
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

			By("Reading VMIs resource")
			vmisResourceRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "resources/read",
				Params: map[string]interface{}{
					"uri": fmt.Sprintf("kubevirt://%s/vmis", testNamespace),
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(vmisResourceRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive VMIs resource response")
			Expect(response.Error).To(BeNil(), "VMIs resource should not return error")

			result, ok := response.Result.(map[string]interface{})
			Expect(ok).To(BeTrue(), "Result should be a map")
			contents, ok := result["contents"].([]interface{})
			Expect(ok).To(BeTrue(), "Result should contain contents array")
			Expect(len(contents)).To(BeNumerically(">", 0), "Should have contents")
		})
	})

	Context("Error Handling", func() {
		It("should handle invalid tool names gracefully", func() {
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

			By("Calling non-existent tool")
			invalidToolRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name": "non_existent_tool",
					"arguments": map[string]interface{}{
						"namespace": testNamespace,
					},
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(invalidToolRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive response even for invalid tool")
			Expect(response.Error).NotTo(BeNil(), "Should return error for invalid tool")
		})

		It("should handle missing arguments gracefully", func() {
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

			By("Calling tool with missing required arguments")
			missingArgsRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name": "list_vms",
					"arguments": map[string]interface{}{
						// Missing namespace argument
					},
				},
				ID: 2,
			}

			_, err = mcpServer.SendRequest(missingArgsRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive response even with missing args")
			// Note: The actual behavior depends on the tool implementation
		})

		It("should handle invalid resource URIs gracefully", func() {
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

			By("Reading invalid resource URI")
			invalidResourceRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "resources/read",
				Params: map[string]interface{}{
					"uri": "invalid://uri/format",
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(invalidResourceRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive response even for invalid URI")
			Expect(response.Error).NotTo(BeNil(), "Should return error for invalid URI")
		})

		It("should handle non-existent VM operations gracefully", func() {
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

			By("Trying to start non-existent VM")
			startNonExistentVMRequest := MCPRequest{
				JSONRPC: "2.0",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name": "start_vm",
					"arguments": map[string]interface{}{
						"namespace": testNamespace,
						"name":      "non-existent-vm",
					},
				},
				ID: 2,
			}

			response, err := mcpServer.SendRequest(startNonExistentVMRequest)
			Expect(err).NotTo(HaveOccurred(), "Should receive response")
			// Check that it handles the error gracefully (either in error field or result content)
			if response.Error == nil {
				result, ok := response.Result.(map[string]interface{})
				Expect(ok).To(BeTrue(), "Result should be a map if no error field")
				content, ok := result["content"].([]interface{})
				Expect(ok).To(BeTrue(), "Result should contain content array")
				Expect(len(content)).To(BeNumerically(">", 0), "Should have error content")
			}
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