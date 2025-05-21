# GoLoadBalancer
Understanding the fundamental knowledge of load balancer. How it operates?

## Load balancer setup

## Listener

## Target group

## Health check
Context from the AWS.
https://docs.aws.amazon.com/elasticloadbalancing/latest/network/target-group-health-checks.html
```
Network Load Balancers use active and passive health checks to determine whether a target is available to handle requests. By default, each load balancer node routes requests only to the healthy targets in its Availability Zone. If you enable cross-zone load balancing, each load balancer node routes requests to the healthy targets in all enabled Availability Zones.
```
Health check status are:
Initial, Healthy, Unhealthy, Draining, Unhealthy.draining, unavailable.

