# Benchmark Result

API: GET /token/validate?token=*token*

## WRK

### Docker

**Information**

- Threads: 50
- Connections: 400
- Durations: 1m

| Thread Stats | Avg | Stdev | Max | +/- Stdev |
|---|---|---|---|---|
|    Latency  |  31.20ms |  19.77ms | 231.13ms |  74.74% |
|    Req/Sec |  270.84  |   93.93  |   1.67k  |  74.53% |

**1150981** requests in **1.00m**, **363.33MB** read

- Requests/sec: **19151.08**

- Transfer/sec:  **6.05MB**

### Local

**Information**

- Threads: 50
- Connections: 400
- Durations: 1m

| Thread Stats | Avg | Stdev | Max | +/- Stdev |
|---|---|---|---|---|
| Latency   | 21.15ms  |  8.47ms | 159.70ms |  77.33% |
|    Req/Sec  | 384.62  |  111.37  |  1.43k  |  73.47% |

**810606** requests in **1.00m**, **255.88MB** read

- Requests/sec: **13490.05**

- Transfer/sec:  **4.26MB**


## Apache Benchmark

### Docker

```
Document Path:          /token/validate?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Ijg2MjIzYzExLWRiMjYtNGVkOC1iOTBlLTI0MmIwOGY4ZjVmYyIsImlhdCI6IjIwMjMtMTAtMjlUMTU6MTk6MTIuNDM4MzU2NTk4KzA3OjAwIiwiZXhwIjoiMjAyMy0xMC0zMFQxNToxOToxMi40MzgzNTY2NTQrMDc6MDAiLCJkYXRhIjp7InBob25lIjoiKzg0OTE4OTE5MzE0IiwidXNlcklkIjoiYTM0NTA3NzctOWUxOC00OTVmLTk0NWMtY2QzNjYzMjAzZjk1In19.I8cCtOQBhGvjyohLPoMjcoq_1TieTUGU8jrb3hvuXLM
Document Length:        26 bytes
```

```
Concurrency Level:      100
Time taken for tests:   1.245 seconds
Complete requests:      5000
Failed requests:        0
Total transferred:      1655000 bytes
HTML transferred:       130000 bytes
Requests per second:    4016.60 [#/sec] (mean)
Time per request:       24.897 [ms] (mean)
Time per request:       0.249 [ms] (mean, across all concurrent requests)
Transfer rate:          1298.33 [Kbytes/sec] received
```


```
Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   0.8      1       7
Processing:     6   23   6.8     23      57
Waiting:        3   22   6.6     22      57
Total:          8   25   6.8     24      59
```

```
Percentage of the requests served within a certain time (ms)
    50%     24
    66%     26
    75%     28
    80%     29
    90%     33
    95%     38
    98%     42
    99%     48
    100%    59 (longest request)
```
    

### Local

```
Document Path:          /token/validate?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Ijg2MjIzYzExLWRiMjYtNGVkOC1iOTBlLTI0MmIwOGY4ZjVmYyIsImlhdCI6IjIwMjMtMTAtMjlUMTU6MTk6MTIuNDM4MzU2NTk4KzA3OjAwIiwiZXhwIjoiMjAyMy0xMC0zMFQxNToxOToxMi40MzgzNTY2NTQrMDc6MDAiLCJkYXRhIjp7InBob25lIjoiKzg0OTE4OTE5MzE0IiwidXNlcklkIjoiYTM0NTA3NzctOWUxOC00OTVmLTk0NWMtY2QzNjYzMjAzZjk1In19.I8cCtOQBhGvjyohLPoMjcoq_1TieTUGU8jrb3hvuXLM
Document Length:        26 bytes
```

```
Concurrency Level:      100
Time taken for tests:   0.242 seconds
Complete requests:      5000
Failed requests:        0
Total transferred:      1655000 bytes
HTML transferred:       130000 bytes
Requests per second:    20680.04 [#/sec] (mean)
Time per request:       4.836 [ms] (mean)
Time per request:       0.048 [ms] (mean, across all concurrent requests)
Transfer rate:          6684.66 [Kbytes/sec] received
```


```
Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   0.6      1       4
Processing:     1    4   1.8      4      11
Waiting:        1    3   1.6      3      11
Total:          1    5   1.7      4      12
```

```
Percentage of the requests served within a certain time (ms)
  50%      4
  66%      5
  75%      6
  80%      6
  90%      7
  95%      8
  98%      9
  99%     10
 100%     12 (longest request)
```
    
