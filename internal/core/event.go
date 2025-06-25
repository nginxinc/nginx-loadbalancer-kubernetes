/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package core

import v1 "k8s.io/api/core/v1"

type EventType int

// Event types
const (

	// Created Represents the event type when a service is created
	Created EventType = iota

	// Updated Represents the event type when a service is updated
	Updated

	// Deleted Represents the event type when a service is deleted
	Deleted
)

// Event represents a service event
type Event struct {

	// Type represents the event type, one of the constant values defined above.
	Type EventType

	// Service represents the service object in its current state
	Service *v1.Service

	// PreviousService represents the service object in its previous state
	PreviousService *v1.Service

	// NodeIps represents the list of node IPs in the Cluster. This is populated by the Watcher when an event is created.
	// The Node IPs are needed by the BorderClient.
	NodeIps []string
}

// NewEvent factory method to create a new Event
func NewEvent(eventType EventType, service *v1.Service, previousService *v1.Service, nodeIps []string) Event {
	return Event{
		Type:            eventType,
		Service:         service,
		PreviousService: previousService,
		NodeIps:         nodeIps,
	}
}
