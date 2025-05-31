# GoLoadBalancer
Understanding the fundamental knowledge of load balancer. Layer 7 Load balancer for now(HTTP Levels), but I would like to make it works for the Layer 4, which is TCP/IP layer.

## Load balancer setup
* Load balancer configuration setup requirements.

## Listener

## Target group

## Health check
* Load balancer should check whether connected servers are available and healthy.
Context from the AWS.
https://docs.aws.amazon.com/elasticloadbalancing/latest/network/target-group-health-checks.html
```
Network Load Balancers use active and passive health checks to determine whether a target is available to handle requests. By default, each load balancer node routes requests only to the healthy targets in its Availability Zone. If you enable cross-zone load balancing, each load balancer node routes requests to the healthy targets in all enabled Availability Zones.
```
Health check status are:
Initial, Healthy, Unhealthy, Draining, Unhealthy.draining, unavailable.



## Algorithms
###  Round Robin(Implemented)
* Distribute incoming client requests across multiple servers in a sequential, cyclical manner.
* Each incoming request is forwarded to the next server, and it loops back to the first server when it hits the end of the server list.
### Weighted Round Robin(In Progress)
* hmmm - researching...
### Least connections
* New connections are connected to the server with the fewest active connections.

## Reference
* https://medium.com/@leonardo5621_66451/building-a-load-balancer-in-go-1c68131dc0ef
* https://www.geeksforgeeks.org/how-do-you-load-balance-tcp-traffic/#importance-of-load-balancing-tcp-traffic
