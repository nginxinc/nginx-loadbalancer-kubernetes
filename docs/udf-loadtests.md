## WRK load tests from Ubuntu Jumphost
## to Nginx LB server
## and direct to each k8s nodeport
## using WRK in a container

### 10.1.1.4 is the Nginx LB Server's IP addr

docker run --rm williamyeh/wrk -t4 -c50 -d2m -H 'Host: cafe.example.com' --timeout 2s https://10.1.1.4/coffee
Running 2m test @ https://10.1.1.4/coffee
  4 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    19.73ms   11.26ms 172.76ms   81.04%
    Req/Sec   626.50    103.68     1.03k    75.60%
  299460 requests in 2.00m, 481.54MB read
Requests/sec:   2493.52
Transfer/sec:      4.01MB

## To knode1

ubuntu@k8-jumphost:~$ docker run --rm williamyeh/wrk -t4 -c50 -d2m -H 'Host: cafe.example.com' --timeout 2s https://10.1.1.8:31269/coffee
Running 2m test @ https://10.1.1.8:31269/coffee
  4 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    17.87ms   10.63ms 151.45ms   80.16%
    Req/Sec   698.98    113.22     1.05k    75.67%
  334080 requests in 2.00m, 537.22MB read
Requests/sec:   2782.35
Transfer/sec:      4.47MB

## t0 knode2

ubuntu@k8-jumphost:~$ docker run --rm williamyeh/wrk -t4 -c50 -d2m -H 'Host: cafe.example.com' --timeout 2s https://10.1.1.10:31269/coffee
Running 2m test @ https://10.1.1.10:31269/coffee
  4 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    17.62ms   10.01ms 170.99ms   80.32%
    Req/Sec   703.96    115.07     1.09k    74.17%
  336484 requests in 2.00m, 541.41MB read
Requests/sec:   2801.89
Transfer/sec:      4.51MB

