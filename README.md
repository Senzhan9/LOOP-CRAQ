generate用法：
./craq-generate --key_count=10 --read_proportion=0 --value_length=4 --distribution=zipf --total_operations=1000 output.txt

benchtest用法：
./craq-bench-test -start 2025-06-16T14:51:00 -file output.txt -ops 1000 -concurrency 80

注意：设置start时间时，要先看系统的时间，命令为date。系统时间和实际时间很可能不相同。
慢慢涨-concurrency

