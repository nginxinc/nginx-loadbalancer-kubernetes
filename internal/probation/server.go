/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package probation

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (

	// Ok is the message returned when a check passes.
	Ok = "OK"

	// ServiceNotAvailable is the message returned when a check fails.
	ServiceNotAvailable = "Service Not Available"

	// ListenPort is the port on which the health server will listen.
	ListenPort = 51031
)

// HealthServer is a server that spins up endpoints for the various k8s health checks.
type HealthServer struct {

	// The underlying HTTP server.
	httpServer *http.Server

	// Support for the "livez" endpoint.
	LiveCheck LiveCheck

	// Support for the "readyz" endpoint.
	ReadyCheck ReadyCheck

	// Support for the "startupz" endpoint.
	StartupCheck StartupCheck
}

// NewHealthServer creates a new HealthServer.
func NewHealthServer() *HealthServer {
	return &HealthServer{
		LiveCheck:    LiveCheck{},
		ReadyCheck:   ReadyCheck{},
		StartupCheck: StartupCheck{},
	}
}

// Start spins up the health server.
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

// Stop shuts down the health server.
func (hs *HealthServer) Stop() {
	if err := hs.httpServer.Close(); err != nil {
		logrus.Errorf("unable to stop probe listener on %s: %v", hs.httpServer.Addr, err)
	}
}

// HandleLive is the handler for the "livez" endpoint.
func (hs *HealthServer) HandleLive(writer http.ResponseWriter, request *http.Request) {
	hs.handleProbe(writer, request, &hs.LiveCheck)
}

// HandleReady is the handler for the "readyz" endpoint.
func (hs *HealthServer) HandleReady(writer http.ResponseWriter, request *http.Request) {
	hs.handleProbe(writer, request, &hs.ReadyCheck)
}

// HandleStartup is the handler for the "startupz" endpoint.
func (hs *HealthServer) HandleStartup(writer http.ResponseWriter, request *http.Request) {
	hs.handleProbe(writer, request, &hs.StartupCheck)
}

// handleProbe handles calling the appropriate Check method and writes the result to the client.
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
