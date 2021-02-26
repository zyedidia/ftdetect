# ftdetect

ftdetect is a library for detecting the filetype of source code files (what
programming language the file is written in). It primarily uses the file
extension and file name to determine the filetype, but also may use the first
line of the file (the header) for additional information (for example,
`#!/bin/bash` on the first line is a good indication of a shell file, even if
there is no extension).

The library is optimized for very fast detection in the common case and
supports saving the detection data structure to a file for very fast
loading/startup time.

A default set of detectors for 135 different languages is provided as an
embedded dataset.

The `./cmd/detect` directory contains an example tool which uses the default
detectors to perform file detection on the first argument to the CLI
application.
