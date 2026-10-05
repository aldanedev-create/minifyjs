# MinifyJS benchmark results
Generated: 2026-10-05 03:10 UTC

## Compression

| Fixture | Input | minifyjs | esbuild | terser | uglify-js | rjsmin |
|---|---|---|---|---|---|---|
| `application.js` | 2391 B | 1144 B (52.2%) | 1105 B (53.8%) | 1172 B (51.0%) | 1173 B (50.9%) | 1460 B (38.9%) |
| `large.js` | 355525 B | 130474 B (63.3%) | 128475 B (63.9%) | 128475 B (63.9%) | 124458 B (65.0%) | 244570 B (31.2%) |
| `medium.js` | 44612 B | 18185 B (59.2%) | 17885 B (59.9%) | 17885 B (59.9%) | 17866 B (60.0%) | 25911 B (41.9%) |
| `small.js` | 895 B | 500 B (44.1%) | 496 B (44.6%) | 540 B (39.7%) | 529 B (40.9%) | 582 B (35.0%) |
| `tiny.js` | 32 B | 25 B (21.9%) | 26 B (18.8%) | 26 B (18.8%) | 24 B (25.0%) | 27 B (15.6%) |

## Throughput (medium.js)

fixture: medium.js (44612 B)
iterations: 10

tool            median ms       MB/s
------------------------------------
minifyjs            13.88        3.2
esbuild             10.45        4.3
terser             526.18        0.1
uglify-js          682.39        0.1
rjsmin               0.32      139.5

## Startup time

iterations: 20

tool            median ms
--------------------------
minifyjs             6.71
esbuild              6.43
terser             138.47
uglify-js          133.38
rjsmin              39.97

## Peak memory (large.js)

fixture: large.js (355525 B)

tool             peak RSS
--------------------------
minifyjs          23.3 MB
esbuild           25.0 MB
terser           192.4 MB
uglify-js        277.0 MB

## Methodology

See `bench/README.md`. Each measurement is the median of multiple runs in isolated subprocesses.
