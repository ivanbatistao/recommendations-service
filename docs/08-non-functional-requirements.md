# Non-Functional Requirements

## Availability
The system must continue responding even if events temporarily stop arriving.

## Scalability
It must be possible to increase event-processing capacity without modifying the API.

## Low Latency
`Target: P95 < 200 ms`

## Throughput
- Pending

## Observability
We want to be able to answer questions such as:
- How many events are arriving?
- How many are failing?
- What is the throughput?
- What is the P99?
- Which worker is taking the longest?