package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	serverproto "protocol"
	"syscall"
	"time"

	glog "github.com/golang/glog"
	protobuf "github.com/golang/protobuf/proto"
)

type ServerApp struct {
	ServerAddr  *net.UDPAddr
	ServerConn  *net.UDPConn
	serverState *ServerState
}

func NewServerApp(listenAddress string) (*ServerApp, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", listenAddress)
	if err != nil {
		return nil, fmt.Errorf("ResolveUDPAddr at %s failed: %v", listenAddress, err)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("ListenUDP failed: %v", err)
	}

	return &ServerApp{ServerAddr: serverAddr, ServerConn: serverConn,
		serverState: newServerState()}, nil
}

func (app *ServerApp) HandleUpdateStateRequest(
	req *serverproto.UpdateStateRequest) {
	app.serverState.UpdateState(req.GetGroupIndex(), req.GetFactionId(),
		req.GetUnderProtection())
	glog.V(1).Infof("UpdateState, group: %d, faction: %d, state: %q",
		req.GetGroupIndex(), req.GetFactionId(), req.GetUnderProtection())
}

func (app *ServerApp) HandlePickRandomFaction(
	req *serverproto.PickRandomFactionRequest) *serverproto.PickRandomFactionResponse {
	randomFactionID := app.serverState.PickRandomFaction(req.GetGroupIndex(),
		req.GetSelfFactionId())
	glog.V(1).Infof("PickRandomFaction, group: %d, self: %d, result: %d",
		req.GetGroupIndex(), req.GetSelfFactionId(), randomFactionID)
	return &serverproto.PickRandomFactionResponse{
		FactionId: protobuf.Uint32(randomFactionID)}
}

func (app *ServerApp) HandleMessage(message serverproto.Message) ([]byte, error) {
	var cmd int32 = int32(message.GetCmd())
	if cmd == int32(serverproto.MessageCommand_UpdateStateRequestCmd) {
		var req serverproto.UpdateStateRequest
		err := protobuf.Unmarshal(message.GetPayload(), &req)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal UpdateStateRequest failed: %v", err)
		}
		app.HandleUpdateStateRequest(&req)
		return nil, nil
	} else if cmd == int32(serverproto.MessageCommand_PickRandomFactionRequestCmd) {
		var req serverproto.PickRandomFactionRequest
		err := protobuf.Unmarshal(message.GetPayload(), &req)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal PickRandomFactionRequest failed: %v", err)
		}

		resp := app.HandlePickRandomFaction(&req)
		encoded, err := protobuf.Marshal(resp)
		if err != nil {
			return nil, fmt.Errorf("Marshal PickRandomFactionResponse failed: %v", err)
		}

		respCommand := uint32(serverproto.MessageCommand_PickRandomFactionResponseCmd)
		reply := serverproto.Message{
			Cmd:     protobuf.Uint32(respCommand),
			Ctx:     message.Ctx,
			Payload: encoded}

		encoded, err = protobuf.Marshal(&reply)
		if err != nil {
			return nil, fmt.Errorf("Marshal Message failed: %v", err)
		}
		return encoded, nil
	} else {
		return nil, fmt.Errorf("Unknown command: %d", cmd)
	}
}

func (app *ServerApp) Run() {
	defer glog.Flush()
	defer app.ServerConn.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	quit := false

	const kBufSize int = 1024
	buf := make([]byte, kBufSize)

	waitDuration, _ := time.ParseDuration("100us")

Exit:
	for !quit {
		select {
		case s := <-c:
			glog.Infof("Receive signal %v, exit...", s)
			glog.Flush()
			quit = true
			break Exit
		default:
			deadline := time.Now().Add(waitDuration)
			app.ServerConn.SetReadDeadline(deadline)
			n, addr, err := app.ServerConn.ReadFromUDP(buf)
			if err != nil {
				nerr, ok := err.(net.Error)
				if !ok || !nerr.Timeout() {
					glog.Errorf("ReadFromUDP failed: %v", err)
				}
				continue
			}
			glog.Infof("Receive %d bytes from %v", n, addr)

			var message serverproto.Message
			err = protobuf.Unmarshal(buf[:n], &message)
			if err != nil {
				glog.Errorf("Unmarshal failed: %v", err)
				continue
			}
			bytes, err := app.HandleMessage(message)
			if err != nil {
				glog.Errorf("HandleMessage failed: %v", err)
				continue
			}

			if bytes != nil {
				n, err := app.ServerConn.WriteToUDP(bytes, addr)
				if err != nil {
					glog.Errorf("WriteToUDP (addr %v) failed: %v", err)
					continue
				}
				glog.Infof("Write %d bytes to %v", n, addr)
			}
		}
	}
}
