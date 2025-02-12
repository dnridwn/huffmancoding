# HuffmanCoding

An implementation of the Huffman algorithm for data compression and decompression using Go.
This project utilizes a priority queue data structure to build a Huffman tree, enabling relatively efficient data compression.

## Features
- **Data Compression**
The `Compress` function transforms input data into a compressed bit representation based on the frequency of each byte.
- **Data Decompression**
The `Decompress` function restores the original data from the compressed output using the constructed Huffman tree.
- **Efficient Data Structure**
Implements a priority queue to repeatedly select the node with the lowest frequency, in line with Huffman Coding principles.
*Implementation details can be found in [huffmancoding.go](huffmancoding.go)*

## Installation
To add **HuffmanCoding** to your project, run the following command:
```bash
go get github.com/dnridwn/huffmancoding
```

## Usage
This project provides compression and decompression functions that can be directly used in your Go programs. Here’s an example:
```
package main

import (
    "fmt"
    "log"

    "github.com/dnridwn/huffmancoding"
)

func main() {
    // Data to compress
    data := []byte("This is a sample data to be compressed using Huffman Coding.")

    // Compress the data
    compressed, totalBitCount, tree, err := huffmancoding.Compress(data)
    if err != nil {
        log.Fatalf("Compression failed: %v", err)
    }
    fmt.Printf("Compressed Data: %v\nTotal bits: %d\n", compressed, totalBitCount)

    // Decompress the data
    decompressed, err := huffmancoding.Decompress(compressed, totalBitCount, tree)
    if err != nil {
        log.Fatalf("Decompression failed: %v", err)
    }
    fmt.Printf("Original Data: %s\n", decompressed)
}
```

**Note:**
In the `Decompress` function, the length parameter refers to the total bit count produced during compression. Make sure to pass this value correctly.

## Code Structure
- **huffmancoding.go**
Contains the main implementation of the compression and decompression algorithms, including:
    - Building a frequency map for bytes.
    - Implementing the priority queue and constructing the Huffman tree.
    - `The Compress and Decompress functions.`

- **huffmancoding_test.go**
Contains unit tests to validate the functionality of both compression and decompression.

## Contribution
Contributions are welcome! Feel free to fork this repository and submit a pull request with improvements, features, or fixes.

## License
This project is licensed under the [MIT License](LICENSE).
