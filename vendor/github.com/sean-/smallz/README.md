# `smallz`
Naive implementation of `pigz` or `gzip` written in Go.

`smallz` is a parallel compression utility that includes a rate-limiting feature
or bandwidth throttle.  Instead of having to compress and send the output
through `pv(1)`, it's now possible to include everything both `pigz(1)` and
`pv(1)` in the same utility.  This is especially useful for database
administrators who need to not be CPU bound, but also need to be mindful of
either disk IO or network bandwidth.  Additionally, there is an uncompressed
mode so you can pass uncompressed data to the IO throttle.

See @klauspost's excellent
writeup,
[Re-balancing Deflate Compression Levels](https://blog.klauspost.com/rebalancing-deflate-compression-levels/),
for additional considerations regarding the use of `gzip`.

As ~90min airplane-ware goes, this should be a useful utility.  Please test.
There are a number of rough edges in terms of usability that could be improved.
See TODOs or FIXME notes for items to clean up (these need to be moved to Github
issues).  "Patches welcome!"  `smallz` is written in Go and inherently easy to
deploy as a "Go-gettable" standalone program (i.e. no dependencies on
`libz(3)`).

## Installation

```
% export GOPATH=`go env GOPATH`
% export PATH=$GOPATH/bin:$PATH
% go get -u github.com/sean-/smallz
% which smallz
```

## Usage

```
% smallz -h
smallz is an parallel compression utility with throttling

Usage:
  smallz <FILE> [flags]

Examples:
  $ echo 'Hello World' | smallz -c -9 - > stdout.gz
  $ smallz -dc stdout.gz
  Hello World
  $ echo 'Hello World' | time smallz -i=1B -c -0 -
  $ smallz -dc stdout.gz
  Hello World

Flags:
  -b, --block-size string   Specify the block-size to use (default "1MiB")
  -C, --compress            Compress the input (default true)
  -0, --compress-0          Skip compressing the input
  -1, --compress-1          Compress the input using the specified compression level
  -2, --compress-2          Compress the input using the specified compression level
  -3, --compress-3          Compress the input using the specified compression level
  -4, --compress-4          Compress the input using the specified compression level
  -5, --compress-5          Compress the input using the specified compression level
  -6, --compress-6          Compress the input using the specified compression level
  -7, --compress-7          Compress the input using the specified compression level
  -8, --compress-8          Compress the input using the specified compression level
  -9, --compress-9          Compress the input using the specified compression level
  -d, --decompress          Decompress the input
      --enable-agent        Enable the gops(1) agent interface
      --enable-pprof        Enable the pprof endpoint interface
  -h, --help                help for smallz
  -i, --io-limit string     Specify the rate at which IO should be ingested, 0B to disable (default "0B")
  -F, --log-format string   Specify the log format ("auto", "zerolog", or "human") (default "auto")
  -l, --log-level string    Log level (default "INFO")
  -t, --num-threads int     Specify the output file to use (default 4)
  -o, --output string       Specify the output file to use
      --pprof-port uint16   Specify the pprof port (default 4243)
  -c, --stdout              This option specifies that output will go to the standard output stream, leaving files intact. (default true)
      --use-color           Use ASCII colors (default true)
```


## Benchmarks

In the interest of preserving the practice of pseudo-benchmarking, the following are three different tests:

1) Using `gzip` for compression.
2) Using `pigz`
3) Using `smallz`

### `gzip` Test

```
  PID USERNAME   PRI NICE   SIZE    RES STATE   C   TIME    WCPU COMMAND
56419 seanc       85    0 57524K 37656K CPU2    2   0:06  74.22% smallz{small
56419 seanc       84    0 57524K 37656K RUN     1   0:06  70.79% smallz{small
56419 seanc       81    0 57524K 37656K CPU0    0   0:04  61.35% smallz{small
56419 seanc       52    0 57524K 37656K uwait   0   0:05  47.99% smallz{small
56419 seanc       52    0 57524K 37656K uwait   1   0:04  33.94% smallz{small
56419 seanc       79    0 57524K 37656K CPU3    3   0:04  25.52% smallz{small

% find /usr/local/bin -type f -print0 | xargs -0 cat | gzip -9 -c - > /dev/null

Time spent in user mode   (CPU seconds) : 50.669s
Time spent in kernel mode (CPU seconds) : 1.020s
Total time                              : 0:50.82s
CPU utilization (percentage)            : 101.6%
Times the process was swapped           : 0
Times of major page faults              : 0
Times of minor page faults              : 959
Time spent in user mode   (CPU seconds) : 50.575s
Time spent in kernel mode (CPU seconds) : 0.989s
Total time                              : 0:50.69s
CPU utilization (percentage)            : 101.6%
Times the process was swapped           : 0
Times of major page faults              : 0
Times of minor page faults              : 957
% find /usr/local/bin -type f -print0 | xargs -0 cat | gzip -9 -c - > gzip.gz
% du -a gzip.gz
200209	gzip.gz
% gzcat gzip.gz | md5 
3d382c184569f3ed0a1dea9f402907ad
% cat ./gzip.gz | /usr/bin/time gzcat -d - | md5
        2.92 real         1.79 user         0.16 sys
3d382c184569f3ed0a1dea9f402907ad
```

### `pigz` Test

```
% find /usr/local/bin -type f -print0 | xargs -0 cat | pigz -9 -c - > /dev/null

Time spent in user mode   (CPU seconds) : 62.714s
Time spent in kernel mode (CPU seconds) : 1.509s
Total time                              : 0:17.23s
CPU utilization (percentage)            : 372.6%
Times the process was swapped           : 0
Times of major page faults              : 7
Times of minor page faults              : 2118
Time spent in user mode   (CPU seconds) : 62.598s
Time spent in kernel mode (CPU seconds) : 1.718s
Total time                              : 0:17.12s
CPU utilization (percentage)            : 375.5%
Times the process was swapped           : 0
Times of major page faults              : 0
Times of minor page faults              : 2119
% find /usr/local/bin -type f -print0 | xargs -0 cat | pigz -9 -c - > pigz.gz
% du -a pigz.gz
200105	pigz.gz
% gzcat pigz.gz | md5 
3d382c184569f3ed0a1dea9f402907ad
% cat ./pigz.gz | /usr/bin/time pigz -d - | md5
        2.13 real         2.42 user         0.28 sys
3d382c184569f3ed0a1dea9f402907ad
```

### `smallz` Test

```
% find /usr/local/bin -type f -print0 | xargs -0 cat | ./smallz -9 -c - > /dev/null

Time spent in user mode   (CPU seconds) : 49.070s
Time spent in kernel mode (CPU seconds) : 1.447s
Total time                              : 0:14.86s
CPU utilization (percentage)            : 339.9%
Times the process was swapped           : 0
Times of major page faults              : 2
Times of minor page faults              : 9902
Time spent in user mode   (CPU seconds) : 49.377s
Time spent in kernel mode (CPU seconds) : 1.385s
Total time                              : 0:14.97s
CPU utilization (percentage)            : 339.0%
Times the process was swapped           : 0
Times of major page faults              : 2
Times of minor page faults              : 9377
% find /usr/local/bin -type f -print0 | xargs -0 cat | pigz -9 -c - > smallz.gz
% du -a smallz.gz
205897	smallz.gz
% gzcat smallz.gz |md5
3d382c184569f3ed0a1dea9f402907ad
% cat ./smallz.gz | /usr/bin/time pigz -d - | md5
        1.69 real         1.90 user         0.26 sys
3d382c184569f3ed0a1dea9f402907ad
% du -a /etc/services 
81	/etc/services

% echo Test various rate-limits
Test various rate-limits
% cat /etc/services | /usr/bin/time smallz -i=4KiB -c -0 - > /dev/null 
       86.21 real         0.20 user         0.24 sys
% cat /etc/services | /usr/bin/time smallz -i=16KiB -c -0 - > /dev/null
       20.97 real         0.05 user         0.05 sys
% cat /etc/services | /usr/bin/time smallz -i=32KiB -c -0 - > /dev/null
       10.42 real         0.03 user         0.03 sys
% cat /etc/services | /usr/bin/time smallz -i=64KiB -c -0 - > /dev/null
        5.12 real         0.00 user         0.03 sys
% cat /etc/services | /usr/bin/time smallz -i=128KiB -c -0 - > /dev/null
        2.52 real         0.00 user         0.01 sys
% cat /etc/services | /usr/bin/time smallz -i=256KiB -c -0 - > /dev/null
        1.29 real         0.01 user         0.01 sys
% cat /etc/services | /usr/bin/time smallz -i=1MiB -c -0 - > /dev/null
        0.32 real         0.00 user         0.00 sys
% cat /etc/services | /usr/bin/time smallz -i=1MiB -b=25KiB -t=4 -c -0 - > /dev/null
        0.31 real         0.00 user         0.00 sys
% cat /etc/services | /usr/bin/time smallz -i=1TiB -b=16KiB -t=4 -c -0 - > /dev/null
        0.00 real         0.01 user         0.00 sys
% cat /etc/services | /usr/bin/time smallz -c -0 - > /dev/null
        0.00 real         0.00 user         0.00 sys

% echo cat-stone benchmarking
cat-stone benchmarking
% doas find /usr/local -type f -print0 | xargs -0 doas cat | /usr/bin/time cat > /dev/null
       18.23 real         0.58 user         3.10 sys

Time spent in user mode   (CPU seconds) : 0.873s
Time spent in kernel mode (CPU seconds) : 21.334s
Total time                              : 0:18.23s
CPU utilization (percentage)            : 121.7%
Times the process was swapped           : 0
Times of major page faults              : 0
Times of minor page faults              : 11508
% doas find /usr/local -type f -print0 | xargs -0 doas cat | /usr/bin/time cat > /dev/null
       18.20 real         0.52 user         3.19 sys
% doas find /usr/local -type f -print0 | xargs -0 doas cat | /usr/bin/time smallz -0 -c - > /dev/null
       18.16 real         1.82 user         3.11 sys

Time spent in user mode   (CPU seconds) : 2.164s
Time spent in kernel mode (CPU seconds) : 21.268s
Total time                              : 0:18.17s
CPU utilization (percentage)            : 128.8%
Times the process was swapped           : 0
Times of major page faults              : 0
Times of minor page faults              : 11980
% doas find /usr/local -type f -print0 | xargs -0 doas cat | /usr/bin/time smallz -0 -c - > /dev/null
       18.72 real         1.57 user         3.41 sys
% doas find /usr/local -type f -print0 | xargs -0 doas cat | /usr/bin/time smallz -0 -c - > /dev/null
       18.13 real         1.55 user         3.39 sys
% doas find /usr/local -type f -print0 | xargs -0 doas cat | /usr/bin/time smallz -0 -c - > /dev/null
       18.32 real         1.64 user         3.32 sys
% doas find /usr/local -type f -print0 | xargs -0 doas cat | /usr/bin/time smallz -0 -i=1TiB -c - > /dev/null
       18.30 real         1.91 user         3.36 sys
% echo Good enough
Good enough
```
