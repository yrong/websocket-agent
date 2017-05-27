// Copyright 2013 Joe Walnes and the websocketd team.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libwebsocketd

import (
	"time"
)

type Endpoint interface {
	StartReading()
	Terminate()
	Output() chan []byte
	ErrorOutput() chan []byte
	Send([]byte) bool
}

const(
	in = iota
	out
	error_out
)

type CommandInfo struct {
	Timestamp    int64 `json:"timestamp"`
	Message string `json:"message"`
	Type int `json:"type"`
	Pid  int `json:"pid"`
	SourceIp string `json:"src_ip"`
	Hostname string `json:"hostname"`
}


func PipeEndpoints(e1, e2 Endpoint, wsh *WebsocketdHandler) {
	e1.StartReading()
	e2.StartReading()

	defer e1.Terminate()
	defer e2.Terminate()

	for {
		select {
		case msg, ok := <-e1.Output():
			if !ok {
				return
			}
			pipe2OtherEndPointAndIndexToES(e1,e2,wsh,msg,out)
		case msg, ok := <-e1.ErrorOutput():
			if !ok {
				return
			}
			pipe2OtherEndPointAndIndexToES(e1,e2,wsh,msg,error_out)
		case msg, ok := <-e2.Output():
			if !ok  {
				return
			}
			pipe2OtherEndPointAndIndexToES(e1, e2,wsh, msg, in)
		}
	}
}

func pipe2OtherEndPointAndIndexToES(e1,e2 Endpoint, wsh *WebsocketdHandler,msg []byte,mtype int) {
	var command CommandInfo
	if process_endpoint, ok := e1.(*ProcessEndpoint); ok {
		command.Message = string(msg)
		command.Pid = process_endpoint.process.cmd.Process.Pid
		command.SourceIp = wsh.RemoteInfo.Host
		command.Timestamp = time.Now().Unix()
		command.Type = mtype
		command.Hostname = wsh.server.Config.HostName
	}
	if mtype == out||mtype == error_out{
		if ws_endpoint, ok := e2.(*WebSocketEndpoint); ok {
			ws_endpoint.SendJson(command)
		}
	}else if mtype == in {
		e1.Send(msg)
	}
	if wsh.server.Config.Log2ES {
		wsh.server.es_handler.index(command)
	}
}

