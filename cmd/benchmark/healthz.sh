ab -n 100000 -c 100 http://127.0.0.1:8080/health
wrk -t100 -c400 -d10s http://127.0.0.1:8080/health
