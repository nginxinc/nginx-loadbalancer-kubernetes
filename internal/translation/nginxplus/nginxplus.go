// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package nginxplus

import (
	"errors"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
	nginxClient "github.com/nginxinc/nginx-plus-go-client/client"
	v1 "k8s.io/api/networking/v1"
)

func Translate(event *core.Event) (*core.Event, error) {
	var err error

	switch event.Type {
	case core.Created:
		event.NginxUpstream, err = translateCreated(event.Ingress)
		if err != nil {
			return event, err
		}

		return event, nil

	case core.Updated:
		event.NginxUpstream, err = translateUpdated(event.Ingress, event.PreviousIngress)
		if err != nil {
			return event, err
		}

		return event, nil

	case core.Deleted:
		event.NginxUpstream, err = translateDeleted(event.Ingress)
		if err != nil {
			return event, err
		}

		return event, nil
	}

	return nil, errors.New("unknown event type")
}

func translateCreated(ingress *v1.Ingress) (*nginxClient.UpstreamServer, error) {
	return nil, nil
}

func translateUpdated(ingress, previousIngress *v1.Ingress) (*nginxClient.UpstreamServer, error) {
	return nil, nil
}

func translateDeleted(ingress *v1.Ingress) (*nginxClient.UpstreamServer, error) {
	return nil, nil
}
