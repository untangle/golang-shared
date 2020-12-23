package restdZmqServer

import (
	"errors"
	"sync"
	"syscall"
	"time"

	zmq "github.com/pebbe/zmq4"
	"github.com/untangle/golang-shared/services/logger"
	"github.com/TiffanyKalin-untangle/fake-packetd/services/dispatch"
	prep "github.com/untangle/golang-shared/structs/protocolbuffers/PacketdReply"
	zreq "github.com/untangle/golang-shared/structs/protocolbuffers/ZMQRequest"
	"google.golang.org/protobuf/proto"
	spb "google.golang.org/protobuf/types/known/structpb"
)

var isShutdown = make(chan struct{})
var wg sync.WaitGroup

type Processer interface {
	process(request *zreq.ZMQRequest) (processedReply []byte, processErr error) 
}

func Startup(proc Processer) {
	logger.Info("Starting zmq service...\n")
	socketServer()

}

func Shutdown() {
	close(isShutdown)
	wg.Wait()
}

func socketServer() {
	zmqSocket, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		logger.Warn("Failed to create zmq socket...", err)
	}

	zmqSocket.Bind("tcp://*:5555")
	wg.Add(1)
	go func(waitgroup *sync.WaitGroup, socket *zmq.Socket) {
		defer socket.Close()
		defer waitgroup.Done()
		tick := time.Tick(500 * time.Millisecond)
		for {
			select {
			case <-isShutdown:
				logger.Info("Shutdown is seen\n")
				return
			case <-tick:
				logger.Info("Listening for requests\n")
				requestRaw, err := socket.RecvMessageBytes(zmq.DONTWAIT)
				if err != nil {
					if zmq.AsErrno(err) != zmq.AsErrno(syscall.EAGAIN) {
						logger.Warn("Error on receive ", err, "\n")
					}
					continue
				}

				// Process message
				request := &zreq.ZMQRequest{}
				if err := proto.Unmarshal(requestRaw[0], request); err != nil {
					logger.Warn("Error on unmasharlling ", err, "\n")
					continue
				}
				logger.Info("Received ", request, "\n")

				reply, err := processMessage(request)
				if err != nil {
					logger.Warn("Error on processing reply: ", err, "\n")
					continue
				}

				socket.SendMessage(reply)
				logger.Info("Sent ", reply, "\n")
			}
		} 
	}(&wg, zmqSocket)
}

func processMessage(request *zreq.ZMQRequest) (processedReply []byte, processErr error) {
	function := request.Function
	reply := &prep.PacketdReply{}

	if function == "GetConntrackTable" {
		conntrackTable := dispatch.GetConntrackTable()
		for _, v := range conntrackTable {
			conntrackStruct, err := spb.NewStruct(v)

			if err != nil {
				return nil, errors.New("Error getting conntrack table: " + err.Error())
			}

			reply.Conntracks = append(reply.Conntracks, conntrackStruct)
		}
		
		
	}

	encodedReply, err := proto.Marshal(reply)
	if err != nil {
		return nil, errors.New("Error encoding reply: " + err.Error())
	}

	return encodedReply, nil
}