# catch2-compare

Generate benchmark result comparison report from catch2 XML reports.

## Usage

```sh
$ /path/to/your/benchmark --reporter XML::out=result-a.xml

# You probably improved the performance of your code.

$ /path/to/your/benchmark --reporter XML::out=result-b.xml

# Compare two reports.
$ catch2-compare result-a.xml result-b.xml
...

# You can compare reports in the directory.
# Note that it doesn't recursively traverse directories.
$ catch2-compare reports-target reports-source
...
```

## Sample

You can have same one by running:
```sh
$ catch2-compare sample/target sample/source
```

```diff
::: [test 1] test case A
+++ /path/to/test-A.cpp:42
+ benchmark a                       4ms          2ms     50.00%
- benchmark b                       3ms          6ms   -100.00%
  benchmark c                        7s            -          -
  benchmark d                         -         1m7s          -
+ benchmark e                      10µs          2µs     80.00%
  benchmark f                         -       7h0m0s          -
≈ benchmark g                     999ns          1µs     -0.10%

::: [test 1] test case C
+++ /path/to/test-C.cpp:69
  benchmark a                         -         13ns          -

::: [test 3] test case A
+++ /path/to/test-A.cpp:12
  benchmark a                         -          11s          -
  Very long long benchm...            -      1h0m13s          -

```
