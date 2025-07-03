/*
🔧 Detailed Problem Statement
You are tasked with designing and implementing a load balancer system that distributes incoming network traffic across a pool of backend servers. The system should efficiently route client requests to backend servers using configurable load balancing algorithms. It must also handle dynamic server registration and removal, maintain real-time metrics, and operate reliably under concurrent load.

🔹 Request Routing
Incoming client requests must be accepted by the load balancer and forwarded to one of the available backend servers.

The load balancer must support multiple routing strategies including:

Round Robin: Sequentially distribute traffic across all servers.

Least Connections: Route to the server with the fewest active connections.

Consistent Hashing: Route requests based on a hash of request attributes (e.g., client IP), ensuring consistent routing for similar clients.

🔹 Backend Server Pool
A list of backend servers must be maintained in memory.

Each server must be able to receive and process forwarded requests.

Servers should expose metadata like current connection count or server health.

The system must support adding or removing backend servers dynamically at runtime.

🔹 Concurrency and Synchronization
The system must handle multiple client requests concurrently and route them in a thread-safe manner.

Server selection logic must be safely shared across multiple goroutines.

Synchronization primitives (e.g., mutexes, channels) must be used to ensure safe concurrent access.

🔹 Health Checks and Failover
The load balancer should periodically perform health checks on backend servers.

Unhealthy servers should be temporarily removed from the active pool.

Recovered servers should be automatically re-added.

🔹 Metrics and Monitoring (Optional, Advanced)
Track metrics such as:

# Number of requests handled per server

# Average response time

# Number of active connections

Expose these metrics via a status endpoint (e.g., /metrics).

🔹 Graceful Shutdown
On shutdown, the load balancer should stop accepting new requests.

Existing in-flight requests should be completed before exiting.

Resources should be cleaned up properly.
*/
package main

func main()  {
	
}