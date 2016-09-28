package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/websocket"
	"golang.org/x/sys/unix"
)

type Client struct {
	ws          *websocket.Conn
	tickChan    <-chan time.Time
	messageChan chan *Message
	errChan     chan error
	doneChan    chan bool
}

func (client *Client) sendReply(result int, msg string) {
	var replyMsg ReplyMessage
	replyMsg.Msg = msg
	replyMsg.Result = result
	if err := websocket.JSON.Send(client.ws, replyMsg); err != nil {
		client.errChan <- err
	}
}

func newClient(ws *websocket.Conn) (client *Client) {
	tickChan := time.NewTicker(time.Second * 1).C
	messageChan := make(chan *Message)
	errChan := make(chan error)
	doneChan := make(chan bool)
	return &Client{ws, tickChan, messageChan, errChan, doneChan}
}

type Message struct {
	Result   int    `json:"result, omitempty"`
	Op       string `json:"op"`
	Year     int    `json:"y"`
	Month    int    `json:"mo"`
	Day      int    `json:"d"`
	Hour     int    `json:"h"`
	Minute   int    `json:"m"`
	Second   int    `json:"s"`
	LockedBy string `json:"lockedby"`
}

type ReplyMessage struct {
	Result int    `json:"result"`
	Msg    string `json:"msg"`
}

func newMessage(t time.Time) (msg Message) {
	msg.Year = t.Year()
	msg.Month = int(t.Month())
	msg.Day = t.Day()
	msg.Hour = t.Hour()
	msg.Minute = t.Minute()
	msg.Second = t.Second()
	msg.LockedBy = lockedBy
	return msg
}

func timeMgrHandler(ws *websocket.Conn) {
	client := newClient(ws)
	go readFromClient(client)
Exit:
	for {
		select {
		case now := <-client.tickChan:
			updateClientTime(client, now)
		case err := <-client.errChan:
			{
				glog.Error(err)
				break Exit
			}
		case <-client.doneChan:
			{
				break Exit
			}
		case msg := <-client.messageChan:
			handlerOperation(client, msg)
		}
	}
}

func updateClientTime(client *Client, now time.Time) {
	msg := newMessage(now)
	if err := websocket.JSON.Send(client.ws, msg); err != nil {
		client.errChan <- err
	}
}

func readFromClient(client *Client) {
	for {
		var msg Message
		err := websocket.JSON.Receive(client.ws, &msg)
		if err == io.EOF {
			client.doneChan <- true
			break
		} else if err != nil {
			client.errChan <- err
			break
		} else {
			client.messageChan <- &msg
		}
	}
}

func handlerOperation(client *Client, msg *Message) {
	switch msg.Op {
	case "modify":
		modifySystemTime(client, msg)
	case "unlock":
		unlockSystemTime(client, msg)
	case "restart":
		restartApache(client, msg)
	default:
		glog.Errorf("Invalid msg.Op: %s\n", msg.Op)
	}
}

func modifySystemTime(client *Client, msg *Message) {
	if lockedBy != "" && lockedBy != msg.LockedBy {
		client.sendReply(-2, "已被其他人锁定！")
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Date(msg.Year, time.Month(msg.Month), msg.Day, msg.Hour,
		msg.Minute, msg.Second, 0, loc)
	tv := unix.NsecToTimeval(t.UnixNano())
	err := unix.Settimeofday(&tv)
	if err != nil {
		client.sendReply(-1, fmt.Sprintf("修改时间失败：%s", err.Error()))
	} else {
		client.sendReply(0, "修改成功！")
		lockedBy = msg.LockedBy
	}
}

func unlockSystemTime(client *Client, msg *Message) {
	if lockedBy != "" && lockedBy != msg.LockedBy {
		client.sendReply(-2, "已被其他人锁定！")
	} else {
		lockedBy = ""
		client.sendReply(0, "解锁成功！")
	}
}

func restartApache(client *Client, msg *Message) {
	cmd := exec.Command("/usr/local/apr-prefork/bin/apachectl", "restart")
	var err error
	defer func() {
		if err != nil {
			client.sendReply(-1, fmt.Sprintf("重启失败：%s", err.Error()))
		}
	}()
	err = cmd.Start()
	if err != nil {
		return
	}
	err = cmd.Wait()
	if err != nil {
		return
	}
	client.sendReply(0, fmt.Sprintf("重启成功！"))
}

var lockedBy string

var host = flag.String("host", "0.0.0.0", "Listening host")
var port = flag.Int("port", 8123, "Listening port")

func main() {
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)
	glog.Infof("Listening address: %v", addr)
	http.Handle("/", websocket.Handler(timeMgrHandler))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
