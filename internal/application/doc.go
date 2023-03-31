/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

/*
Package application includes support for applying updates to the Border servers.

"Border Servers" are the servers that are exposed to the outside world and direct traffic into the cluster.
At this time the only supported Border Servers are NGINX Plus servers. The BorderClient module defines
an interface that can be implemented to support other Border Server types.

- HttpBorderClient: updates NGINX Plus servers using HTTP Upstream methods on the NGINX Plus API.
- TcpBorderClient: updates NGINX Plus servers using Stream Upstream methods on the NGINX Plus API.

Selection of the appropriate client is based on the Annotations present on the NodePort Service definition.
*/

package application
