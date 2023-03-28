/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package probation

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	Ok                  = "OK"
	ServiceNotAvailable = "Service Not Available"
	ListenPort          = 51031
)

type HealthServer struct {
	httpServer   *http.Server
	LiveCheck    LiveCheck
	ReadyCheck   ReadyCheck
	StartupCheck StartupCheck
}

func NewHealthServer() *HealthServer {
	return &HealthServer{
		LiveCheck:    LiveCheck{},
		ReadyCheck:   ReadyCheck{},
		StartupCheck: StartupCheck{},
	}
}

func (hs *HealthServer) Start() {
	logrus.Debugf("Starting probe listener on port %d", ListenPort)

	address := fmt.Sprintf(":%d", ListenPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/livez", hs.HandleLive)
	mux.HandleFunc("/readyz", hs.HandleReady)
	mux.HandleFunc("/startupz", hs.HandleStartup)
	hs.httpServer = &http.Server{Addr: address, Handler: mux}

	go func() {
		if err := hs.httpServer.ListenAndServe(); err != nil {
			logrus.Errorf("unable to start probe listener on %s: %v", hs.httpServer.Addr, err)
		}
	}()

	logrus.Info("Started probe listener on", hs.httpServer.Addr)
}

func (hs *HealthServer) Stop() {
	if err := hs.httpServer.Close(); err != nil {
		logrus.Errorf("unable to stop probe listener on %s: %v", hs.httpServer.Addr, err)
	}
}

func (hs *HealthServer) HandleLive(writer http.ResponseWriter, request *http.Request) {
	hs.handleProbe(writer, request, &hs.LiveCheck)
}

func (hs *HealthServer) HandleReady(writer http.ResponseWriter, request *http.Request) {
	hs.handleProbe(writer, request, &hs.ReadyCheck)
}

func (hs *HealthServer) HandleStartup(writer http.ResponseWriter, request *http.Request) {
	hs.handleProbe(writer, request, &hs.StartupCheck)
}

func (hs *HealthServer) handleProbe(writer http.ResponseWriter, _ *http.Request, check Check) {
	if check.Check() {
		writer.WriteHeader(http.StatusOK)

		if _, err := fmt.Fprint(writer, Ok); err != nil {
			logrus.Error(err)
		}

	} else {
		writer.WriteHeader(http.StatusServiceUnavailable)

		if _, err := fmt.Fprint(writer, ServiceNotAvailable); err != nil {
			logrus.Error(err)
		}
	}
}
