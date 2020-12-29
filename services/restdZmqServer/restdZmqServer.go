package restdzmqserver

import (
	"sync"
	"syscall"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/untangle/golang-shared/services/logger"
	zreq "github.com/untangle/golang-shared/structs/protocolbuffers/ZMQRequest"
	"google.golang.org/protobuf/proto"
)

var isShutdown = make(chan struct{})
var wg sync.WaitGroup

// Processer is function for processing server functions
type Processer interface {
	Process(request *zreq.ZMQRequest) (processedReply []byte, processErr error) 
}

// Startup the server function 
func Startup(processer Processer) {
	logger.Info("Starting zmq service...\n")
	socketServer(processer)
}

// Shutdown the server function 
func Shutdown() {
	close(isShutdown)
	wg.Wait()
}

/* Main server funcion, creates socket and runs goroutine to keep server open */
func socketServer(processer Processer) {
	zmqSocket, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		logger.Warn("Failed to create zmq socket...", err)
	}

	zmqSocket.Bind("tcp://*:5555")
	wg.Add(1)
	go func(waitgroup *sync.WaitGroup, socket *zmq.Socket, proc Processer) {
		defer socket.Close()
		defer waitgroup.Done()
		tick := time.Tick(500 * time.Millisecond)
		for {
			select {
			case <-isShutdown:
				logger.Info("Shutdown is seen\n")
				return
			case <-tick:
				logger.Debug("Listening for requests\n")
				var serverErr error
				serverErr = nil
				var reply []byte
				var replyErr error
				requestRaw, err := socket.RecvMessageBytes(zmq.DONTWAIT)
				if err != nil {
					if zmq.AsErrno(err) != zmq.AsErrno(syscall.EAGAIN) {
						serverErr = "Error on receive " + err.Error()
					}
				} else {
					request := &zreq.ZMQRequest{}
					err := proto.Unmarshal(requestRaw[0], request)
					if err != nil {
						serverErr = "Error on unmasharlling " + err.Error()
					} else {
						logger.Debug("Received ", request, "\n")

						reply, replyErr = processMessage(proc, request)
						if replyErr != nil {
							serverErr = "Error on processing reply: " + err.Error()
						}
					}
				}

				if serverErr != nil {
					errReply := &prep.PacketdReply{ServerError: serverErr}
					reply, replyErr = proto.Marshal(errReply)
					if replyErr != nil {
						socket.SendMessage("Bad Error")
						continue
					}
				}

				socket.SendMessage(reply)
				logger.Info("Sent ", reply, "\n")
			}
		} 
	}(&wg, zmqSocket, processer)
}

/* helper function that calls the interface's process function */
func processMessage(proc Processer, request *zreq.ZMQRequest) (processedReply []byte, processErr error) {
	return proc.Process(request)
}