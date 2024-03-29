/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

/*
Package application includes support for applying updates to the Border servers.

"Border Servers" are servers that are exposed to the outside world and direct traffic into the cluster.

The BorderClient module defines an interface that can be implemented to support other Border Server types.
To add a Border Server client...
1. Create a module that implements the BorderClient interface;
2. Add a new constant in application_constants.go that acts as a key for selecting the client;
3. Update the NewBorderClient factory method in border_client.go that returns the client;

At this time the only supported Border Servers are NGINX Plus servers.

The two Border Server clients for NGINX Plus are:
- NginxHTTPBorderClient: updates NGINX Plus servers using HTTP Upstream methods on the NGINX Plus API.
- NginxStreamBorderClient: updates NGINX Plus servers using Stream Upstream methods on the NGINX Plus API.

Both of these implementations use the NGINX Plus client module to communicate with the NGINX Plus server.

Selection of the appropriate client is based on the Annotations present on the Service definition, e.g.:

	annotations:
	  nginxinc.io/nlk-<upstream name>: <value>

where <upstream name> is the name of the upstream in the NGINX Plus configuration
and <value> is one of the constants in application_constants.go.
*/

package application
