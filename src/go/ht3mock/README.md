# ht3mock

Mockup to simulate HTTP workload from a client to a server
using profiles for

- number of requests
- size of requests and responses
- max. number of requests in parallel
- client and server processing time

Load profiles are applied in parallel.
Background loads may also be applied, they run until all non-background are finished.

The test is used to compare HTTP/1.1 and HTTP/3 performances in various conditions.
