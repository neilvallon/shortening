# Shortening
A bijective base<sub>64</sub> encoder for creating short URLs.


## Padding significant
Most encoders use decimal bases to convert characters.
This can waste space by refusing to use the zero-value in the most significant position.

For example, in base<sub>10</sub> the numbers **1**, **01**, **00001**,
are all equivalent despite being padded with leading 0's.

This is why non-bijective encoders will skip values such as **A** and **AAA**
in favor of possibly longer strings.

By using these values we can recover 64<sup>(n-1)</sup> IDs
of length *n*.

| Length | Radix<sub>64</sub> |   Shortening  | difference |
|--------|--------------------|---------------|------------|
|    1   |                 63 |            64 |          1 |
|    2   |              4,032 |         4,096 |         64 |
|    3   |            258,048 |       262,144 |       4096 |
|    4   |         16,515,072 |    16,777,216 |    262,144 |
|    5   |      1,056,964,608 | 1,073,741,824 | 16,777,216 |


## Performance
|  Benchmark   |  ns/op |
|--------------|--------|
|  Encode64    |   46.2 |
|  Decode64    |   21.6 |
|  Encode32    |   46.4 |
|  Decode32    |   23.7 |

* go1.12.4 darwin/amd64
* Intel Xeon X5675 - 3.06 GHz

## References
Wikipedia: [Bijective numeration](https://en.wikipedia.org/wiki/Bijective_numeration)
