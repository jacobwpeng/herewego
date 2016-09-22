package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	serverproto "protocol"

	"github.com/golang/glog"
	protobuf "github.com/golang/protobuf/proto"
)

type ServerApp struct {
	serverAddr  *net.UDPAddr
	serverConn  *net.UDPConn
	serverState *ServerState
}

func NewServerApp(listeningAddress string) (*ServerApp, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", listeningAddress)
	if err != nil {
		return nil, fmt.Errorf("ResolveUDPAddr at %s failed: %v",
			listeningAddress, err)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("ListenUDP failed: %v", err)
	}

	return &ServerApp{serverAddr: serverAddr, serverConn: serverConn,
		serverState: newServerState()}, nil
}

func (app *ServerApp) HandleUpdateRegion(req *serverproto.UpdateRegionRequest) {
	glog.V(1).Info(req)
	app.serverState.UpdateUser(req.GetUin(), req.GetLevel(), req.GetRegion())
}

func (app *ServerApp) HandleRegionStatus(req *serverproto.RegionStatusRequest,
	resp *serverproto.RegionStatusResponse) {
	glog.V(1).Info(req)
	status := app.serverState.GetRegionStatus()
	resp.Status = make([]*serverproto.RegionStatus, 0, len(status))
	for _, v := range status {
		protoStatus := serverproto.RegionStatus{
			Region: protobuf.Uint32(v.region),
			Num:    protobuf.Uint32(v.userNum),
		}
		resp.Status = append(resp.Status, &protoStatus)
	}
	glog.V(1).Info(resp)
}

func (app *ServerApp) HandlePickFromRegion(req *serverproto.PickFromRegionRequest,
	resp *serverproto.PickFromRegionResponse) {
	glog.V(1).Info(req)
	uins := app.serverState.PickUser(req.GetSelfUin(),
		req.GetSelfLevel(), req.GetExpectRegion())
	resp.Uin = append(resp.Uin, uins...)
	glog.V(1).Info(resp)
}

func (app *ServerApp) HandleMessage(
	message *serverproto.Message) ([]byte, error) {
	cmd := serverproto.MessageCommand(message.GetCmd())
	if cmd == serverproto.MessageCommand_UpdateRegionRequestCmd {
		var req serverproto.UpdateRegionRequest
		err := protobuf.Unmarshal(message.GetPayload(), &req)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal UpdateRegionRequest failed: %v", err)
		}
		app.HandleUpdateRegion(&req)
		return nil, nil
	} else if cmd == serverproto.MessageCommand_RegionStatusRequestCmd {
		var req serverproto.RegionStatusRequest
		err := protobuf.Unmarshal(message.GetPayload(), &req)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal RegionStatusRequest failed: %v", err)
		}
		var resp serverproto.RegionStatusResponse
		app.HandleRegionStatus(&req, &resp)

		encoded, err := protobuf.Marshal(&resp)
		if err != nil {
			return nil, fmt.Errorf("Marshal RegionStatusResponse failed: %v", err)
		}

		respCmd := uint32(serverproto.MessageCommand_RegionStatusResponseCmd)
		reply := serverproto.Message{
			Cmd:     protobuf.Uint32(respCmd),
			Ctx:     message.Ctx,
			Payload: encoded,
		}

		encoded, err = protobuf.Marshal(&reply)
		if err != nil {
			return nil, fmt.Errorf("Marshal Message failed: %v", err)
		}
		return encoded, nil
	} else if cmd == serverproto.MessageCommand_PickFromRegionRequestCmd {
		var req serverproto.PickFromRegionRequest
		err := protobuf.Unmarshal(message.GetPayload(), &req)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal PickFromRegionRequest failed: %v", err)
		}
		var resp serverproto.PickFromRegionResponse
		app.HandlePickFromRegion(&req, &resp)

		encoded, err := protobuf.Marshal(&resp)
		if err != nil {
			return nil, fmt.Errorf("Marshal PickFromRegionResponse failed: %v", err)
		}

		respCmd := uint32(serverproto.MessageCommand_PickFromRegionResponseCmd)
		reply := serverproto.Message{
			Cmd:     protobuf.Uint32(respCmd),
			Ctx:     message.Ctx,
			Payload: encoded,
		}
		encoded, err = protobuf.Marshal(&reply)
		if err != nil {
			return nil, fmt.Errorf("Marshal Message failed: %v", err)
		}
		return encoded, nil
	}
	return nil, fmt.Errorf("Invalid cmd: %d", cmd)
}

func (app *ServerApp) Run() {
	defer glog.Flush()
	defer app.serverConn.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	quit := false

	const kBufSize int = 1024
	buf := make([]byte, kBufSize)

Exit:
	for !quit {
		select {
		case s := <-c:
			glog.Infof("Receive signal %v, exit...", s)
			quit = true
			break Exit
		default:
			app.serverConn.SetReadDeadline(time.Now().Add(time.Millisecond * 10))
			n, addr, err := app.serverConn.ReadFromUDP(buf)
			if err != nil {
				nerr, ok := err.(net.Error)
				if !ok || !nerr.Timeout() {
					glog.Errorf("ReadFromUDP failed: %v", err)
				}
				continue
			}
			glog.V(2).Infof("Receive %d bytes from %v", n, addr)

			var message serverproto.Message
			err = protobuf.Unmarshal(buf[:n], &message)
			if err != nil {
				glog.Errorf("Unmarshal wrapper message failed: %v", err)
				continue
			}

			bytes, err := app.HandleMessage(&message)
			if err != nil {
				glog.Errorf("HandleMessage failed: %v", err)
				continue
			}

			if bytes != nil {
				n, err := app.serverConn.WriteToUDP(bytes, addr)
				if err != nil {
					glog.Errorf("WriteToUDP (addr %v) failed: %v", addr, err)
					continue
				}
				glog.V(2).Infof("Write %d bytes to %v", n, addr)
			}
		}
	}
}
