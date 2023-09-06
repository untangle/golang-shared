package zmqserver

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"github.com/stretchr/testify/assert"
	zreq "github.com/untangle/golang-shared/structs/protocolbuffers/ZMQRequest"
)

// MockProcesser is a mock implementation of the Processer interface for testing purposes.
type MockProcesser struct{}

func (mp *MockProcesser) Process(request *zreq.ZMQRequest) ([]byte, error) {
	if request != nil {
		data := []byte(request.Data)
		return data, nil
	}
	return nil, nil
}

func (mp *MockProcesser) ProcessError(processError string) ([]byte, error) {
	if len(processError) > 0 {
		return []byte(processError), nil
	}
	return nil, nil
}

func TestSocketServer(t *testing.T) {
	// Create a mock processer
	mockProcesser := &MockProcesser{}

	// Start the socket server in a goroutine
	go socketServer(mockProcesser)

	// Wait for some time to start socket server
	time.Sleep(1 * time.Second)

	// Create a ZMQ request and send it to the server
	reqSocket, _ := zmq.NewSocket(zmq.REQ)
	reqSocket.Connect("tcp://localhost:5555")
	defer reqSocket.Close()

	inputData := "ZMQ Request Data"
	request := zreq.ZMQRequest{
		Data: inputData,
	}
	zmqRequest, err := proto.Marshal(&request)
	if err != nil {
		logger.Warn("error in proto marshalling: %v \n", err.Error())
	}
	reqSocket.SendMessage(zmqRequest)

	// Wait for the server to process the request and receive the response.
	time.Sleep(1 * time.Second)

	// Capture the server's response
	response, err := reqSocket.RecvMessageBytes(zmq.DONTWAIT)
	if err != nil {
		t.Fatalf("Error receiving response from server: %v", err)
	}

	// Check the server's response
	expectedResponse := inputData
	actualResponse := string(response[0])
	assert.Equal(t, expectedResponse, actualResponse, "Server response should match expected response")

	// Stop the server gracefully
	close(isShutdown)
}

func TestProcessMessage(t *testing.T) {
	// Create a mock processer
	mockProcesser := &MockProcesser{}

	// Create a sample ZMQRequest
	inputData := "MockProcessMessage"
	request := &zreq.ZMQRequest{
		Data: inputData,
	}

	// Call the processMessage function
	reply, err := processMessage(mockProcesser, request)
	if err != nil {
		logger.Warn("Error in processMessage call : %v \n", err.Error())
	}

	// Check response of processMessage.
	expectedResponse := inputData
	actualResponse := string(reply)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestProcessErrorMessage(t *testing.T) {
	// Create a mock processer
	mockProcesser := &MockProcesser{}

	// Create a sample error message
	serverErr := "Sample Error"

	// Call the processErrorMessage function
	reply, err := processErrorMessage(mockProcesser, serverErr)
	if err != nil {
		logger.Warn("Error in pprocessErrorMessage call : %v \n", err.Error())
	}

	// Check response of processErrorMessage.
	expectedResponse := serverErr
	actualResponse := string(reply)
	assert.Equal(t, expectedResponse, actualResponse)
}
