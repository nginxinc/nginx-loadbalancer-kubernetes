## NginxLB API testing:

List upstreams in block "nginx-lb-https":

curl http://10.1.1.4:9000/api/8/stream/upstreams/nginx-lb-https |jq

Add upstream with JSON:

curl -X POST "http://10.1.1.4:9000/api/8/stream/upstreams/nginx-lb-https/servers/" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"server\": \"10.1.1.8:31269\" }"

curl -X POST "http://10.1.1.4:9000/api/8/stream/upstreams/nginx-lb-https/servers/" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"server\": \"10.1.1.10:31269\" }"

Disable upstream #2 ( id required ):

curl -X PATCH -d "{ \"down\": true }" -s 'http://10.1.1.4:9000/api/8/stream/upstreams/nginx-lb-https/servers/2' -H "accept: application/json" -H "Content-Type: application/json"

Enable upstream #2 ( id required ):

curl -X PATCH -d "{ \"down\": false }" -s 'http://10.1.1.4:9000/api/8/stream/upstreams/nginx-lb-https/servers/2' -H "accept: application/json" -H "Content-Type: application/json"

Delete upstream server with id=0( id required ):

curl -X DELETE -s 'http://10.1.1.4:9000/api/8/stream/upstreams/nginx-lb-https/servers/0' -H "accept: application/json" -H "Content-Type: application/json"
