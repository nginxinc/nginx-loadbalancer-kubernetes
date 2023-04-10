/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package core

import v1 "k8s.io/api/core/v1"

type EventType int

const (
	Created EventType = iota
	Updated
	Deleted
)

type Event struct {
	Type            EventType
	Service         *v1.Service
	PreviousService *v1.Service
	NodeIps         []string
}

func NewEvent(eventType EventType, service *v1.Service, previousService *v1.Service, nodeIps []string) Event {
	return Event{
		Type:            eventType,
		Service:         service,
		PreviousService: previousService,
		NodeIps:         nodeIps,
	}
}
