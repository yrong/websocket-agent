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
	Send([]byte) bool
}

const(
	in = iota
	out
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
		case msgOne, ok := <-e1.Output():
			if !ok || !e2.Send(msgOne) {
				return
			}
			if(wsh.server.Config.Log2ES){
				indexToEs(e1,e2,wsh,string(msgOne),out)
			}
		case msgTwo, ok := <-e2.Output():
			if !ok || !e1.Send(msgTwo) {
				return
			}
			if(wsh.server.Config.Log2ES) {
				indexToEs(e1, e2, wsh, string(msgTwo), in)
			}
		}
	}
}

func indexToEs(e1, e2 Endpoint, wsh *WebsocketdHandler,msg string,mtype int) {
	if process_endpoint, ok := e1.(*ProcessEndpoint); ok {
		var command CommandInfo
		command.Message = msg
		command.Pid = process_endpoint.process.cmd.Process.Pid
		command.SourceIp = wsh.RemoteInfo.Host
		command.Timestamp = time.Now().Unix()
		command.Type = mtype
		command.Hostname = wsh.server.Config.HostName
		wsh.server.es_handler.index(command)
	}
}
