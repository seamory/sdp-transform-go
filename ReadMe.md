# SDP Transform

This project is an implementation of SDP transform based on the golang language, and the methods and interfaces have
been made as consistent as possible with the SDP transform implementation. Except for the processing of some data types,
the number type of the source project is all processed by the string of golang, so in the specific business, please
manually transform the corresponding data types.

A simple parser and writer of SDP. Defines internal grammar based on RFC4566 - SDP, RFC5245 - ICE, and many more.

For simplicity it will force values that are integers to integers and leave everything else as strings when parsing. The
module should be simple to extend or build upon, and is constructed rigorously.

# Usage - Parser

please reference: [SDP Transform Usage - Parser](https://github.com/clux/sdp-transform?tab=readme-ov-file#usage---parser)

# Usage - Writer

Please reference: [SDP Transform Usage - Writer](https://github.com/clux/sdp-transform?tab=readme-ov-file#usage---writer)

# Usage - Custom grammar

Please reference: [SDP Transform Usage - Custom grammar](https://github.com/clux/sdp-transform?tab=readme-ov-file#usage---custom-grammar)
