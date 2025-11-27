| Component     | Why It Matters                                                                  |
| ------------- | ------------------------------------------------------------------------------- |
| **MX Record** | Does this domain accept email? Without MX, email cannot be delivered.           |
| **SPF**       | Does domain define a list of allowed mail senders? Protection against spoofing. |
| **DMARC**     | Does domain enforce email authenticity and reporting? Prevents impersonation.   |

Concurrency → Goroutines
Parallelism → When goroutines run on multiple CPU cores
Multithreading → Managed by Go runtime internally (you do NOT create threads yourself)

goroutines:very lightweight threads
channels:for communication
sync.WaitGroup to wait for all workers to finish

chan T means: a normal channel
<-chan T means: receive-only channel
chan<- T means: send-only channel

Why use receive-only & send-only channels?
This is good engineering practice.
It prevents accidental misuse
It guarantees workers behave correctly
It avoids race conditions

This ensures:
Workers only consume jobs
Workers only produce results

Go Uses an M:N Scheduler 

    M = OS threads

    N = Goroutines

Go maps many goroutines (N) onto few OS threads (M).
e.g 10000 goroutines → 4 OS threads → 4 CPU cores
Go scheduler is very fast
Goroutines are extremely lightweight
OS threads are heavy


Go Scheduler Responsibilities

    Pause goroutines waiting for I/O

    Move ready goroutines to runnable queue

    Assign queues to OS threads

    Run goroutines on available CPU cores

This is why Go runs 10,000 tasks easily while other languages struggle.

How Worker Pools Scale
lets say workerCount := 50
This means:

    50 goroutines actively processing domains

    All running in parallel (if you have enough cores)

Why it scales well?
DNS lookup is I/O bound
While one worker waits for DNS
→ other workers can run
 Go runtime automatically balances them

Even with 1000 workers:

    No thread explosion

    No memory blow

    No performance issues

This is because:

    Goroutines cost almost nothing

    Go uses M:N scheduling
