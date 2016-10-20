package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	serverproto "protocol"

	"github.com/golang/glog"
	protobuf "github.com/golang/protobuf/proto"
)

type ServerApp struct {
	m              sync.RWMutex
	wg             sync.WaitGroup
	serverState    *ServerState
	httpServer     *http.Server
	udpConn        *net.UDPConn
	exitChan       chan struct{}
	managerAddress string
}

func NewServerApp(udpAddr string, mgrAddr string) (*ServerApp, error) {
	httpServer := &http.Server{
		Addr: udpAddr,
	}

	udpServerAddr, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("ResolveUDPAddr failed: %v", err)
	}

	udpServerConn, err := net.ListenUDP("udp", udpServerAddr)
	if err != nil {
		return nil, fmt.Errorf("ListenUDP failed: %v", err)
	}

	return &ServerApp{
		serverState:    newServerState(),
		httpServer:     httpServer,
		udpConn:        udpServerConn,
		exitChan:       make(chan struct{}),
		managerAddress: mgrAddr,
	}, nil
}

func (app *ServerApp) Run() {
	defer glog.Flush()
	defer app.udpConn.Close()
	http.HandleFunc("/", app.handleQueryAll)
	http.HandleFunc("/update", app.handleUpdateServerState)
	go http.ListenAndServe(app.managerAddress, nil)

	app.wg.Add(1)
	go app.listenAndServeUDP()
	app.waitForSignal()
	close(app.exitChan)
	app.wg.Wait()
}

func (app *ServerApp) waitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	s := <-c
	glog.Infof("Receive signal %v, exit...", s)
}

func (app *ServerApp) listenAndServeUDP() {
	defer app.wg.Done()
	quit := false
	for !quit {
		select {
		case <-app.exitChan:
			quit = true
		default:
			app.udpConn.SetReadDeadline(time.Now().Add(time.Millisecond * 10))
			buf := make([]byte, 1024)
			n, addr, err := app.udpConn.ReadFromUDP(buf)
			if err != nil {
				nerr, ok := err.(net.Error)
				if !ok {
					glog.Errorf("Convert to net.Error failed")
				} else if !nerr.Timeout() {
					glog.Errorf("ReadFromUDP failed: %v", err)
				} else {
					// timeout
				}
				continue
			}

			var req serverproto.QueryUserGroup
			err = protobuf.Unmarshal(buf[:n], &req)
			if err != nil {
				glog.Errorf("Unmarshal QueryUserGroup failed: %v", err)
				continue
			}

			var resp serverproto.QueryUserGroupReply
			app.handleQueryUser(&req, &resp)

			reply, err := protobuf.Marshal(&resp)
			if err != nil {
				glog.Fatal(err)
			}

			n, err = app.udpConn.WriteToUDP(reply, addr)
			if err != nil {
				glog.Errorf("WriteToUDP failed: %v", err)
				continue
			}
		}
	}
}

type errorReason struct {
	Reason string
}

func (app *ServerApp) errorOutput(w http.ResponseWriter, err error) {
	templateContent := "<h1>{{.Reason}}</h1>"
	tmpl := template.Must(template.New("ErrorOutput").Parse(templateContent))
	reason := errorReason{Reason: err.Error()}
	tmpl.Execute(w, reason)
}

func (app *ServerApp) handleQueryUser(req *serverproto.QueryUserGroup,
	resp *serverproto.QueryUserGroupReply) {
	uin := req.GetUin()
	group := app.serverState.GetUserGroup(uin)
	resp.Uin = protobuf.Uint32(uin)
	resp.UserGroup = protobuf.Uint32(group)
	resp.SalesFrom = protobuf.Uint32(uint32(app.serverState.salesBegin.Unix()))
	resp.SalesTo = protobuf.Uint32(uint32(app.serverState.salesEnd.Unix()))
}

func (app *ServerApp) handleUpdateServerState(w http.ResponseWriter,
	r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		app.errorOutput(w, fmt.Errorf("ParseMediaType failed: %v", err))
		return
	}
	if !strings.HasPrefix(mediaType, "multipart/") {
		app.errorOutput(w, fmt.Errorf("Invalid mediaType: %s", mediaType))
		return
	}
	contentLength, err := strconv.ParseUint(r.Header.Get("Content-Length"),
		10, 32)
	if err != nil {
		app.errorOutput(w, fmt.Errorf("Parse Content-Length failed: %v", err))
		return
	}
	allContent := make([]byte, 0, contentLength)
	mr := multipart.NewReader(r.Body, params["boundary"])
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Fatal(err)
		}
		data, err := ioutil.ReadAll(p)
		if err != nil {
			glog.Fatal(err)
		}
		allContent = append(allContent, data...)
	}
	scanner := bufio.NewScanner(bytes.NewReader(allContent))
	lineNo := 0
	var begin time.Time
	var end time.Time
	userGroups := make(map[uint32]uint32)
	const timeLayout = "2006-01-02 15:04:05"
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		glog.Fatal(err)
	}
	for scanner.Scan() {
		lineNo += 1
		line := strings.TrimSpace(scanner.Text())

		if lineNo == 1 {
			begin, err = time.ParseInLocation(timeLayout, line, location)
			if err != nil {
				app.errorOutput(w, fmt.Errorf("Parse sales begin(line:1) failed: %v",
					err))
				return
			}
		} else if lineNo == 2 {
			end, err = time.ParseInLocation(timeLayout, line, location)
			if err != nil {
				app.errorOutput(w, fmt.Errorf("Parse sales end(line:2) failed: %v",
					err))
				return
			}
		} else {
			fields := strings.Fields(line)
			if len(fields) != 2 {
				app.errorOutput(w, fmt.Errorf("Invalid content at line: %d", lineNo))
				return
			}
			uin, err := strconv.ParseUint(fields[0], 10, 32)
			if err != nil {
				app.errorOutput(w, fmt.Errorf("Invalid uin at line: %d", lineNo))
				return
			}
			group, err := strconv.ParseUint(fields[1], 10, 32)
			if err != nil {
				app.errorOutput(w, fmt.Errorf("Invalid group at line: %d", lineNo))
				return
			}
			userGroups[uint32(uin)] = uint32(group)
		}
	}
	app.serverState.m.Lock()
	defer app.serverState.m.Unlock()

	app.serverState.salesBegin = begin
	app.serverState.salesEnd = end
	app.serverState.userGroups = userGroups
	glog.Infof("ServerState updated, begin: %v, end: %v, size: %d",
		begin, end, len(userGroups))
	w.Write([]byte("<h1>Update succeed!</h1>"))
}

func (app *ServerApp) handleQueryAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	app.m.RLock()
	defer app.m.RUnlock()

	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "<h1>Begin: %v</h1>", app.serverState.salesBegin)
	fmt.Fprintf(&buffer, "<h1>End: %v</h1>", app.serverState.salesEnd)
	fmt.Fprintf(&buffer, "<h1>User Num: %d</h1>", len(app.serverState.userGroups))
	for k, v := range app.serverState.userGroups {
		fmt.Fprintf(&buffer, "%d: %d<br />", k, v)
	}
	w.Write(buffer.Bytes())
}
