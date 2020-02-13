# GoLinkFinder

A minimal JS endpoint extractor

# Why?

To extract endpoints in both HTML source and embedded javascript files. Useful for bug hunters, red teamers, infosec ninjas.

# Version

1.0.0-alpha

# Usage?

```[-d|--domain] is required
usage: goLinkFinder [-h|--help] -d|--domain "<value>" [-o|--out "<value>"]
                    GoLinkFinder
Arguments:

  -h  --help    Print help information
  -d  --domain  Input a URL.
  -o  --out     File name :  (e.g : output.txt)
```

# How?

best used with grep

```
GoLinkFinder -d https://github.com | grep api
```

Output :

```
 "https://api.github.com/_private/browser/stats"
 "https://api.github.com/_private/browser/errors"
```

you can easily pipe out its with your other tools.

# Watch

[![asciicast](https://asciinema.org/a/HSM3Po0HC8s03XtXw3kw2UuHa.svg)](https://asciinema.org/a/HSM3Po0HC8s03XtXw3kw2UuHa)

# Requirements

Go >= 1.13

# Installation

```
Git clone https://github.com/0xsha/GoLinkFinder.git
cd GoLinkFinder
go build GoLinkfinder.go
```

# Feature request or found an issue?

Please write a patch to fix it and then pull a request.

# References

Python implementation:
https://github.com/GerbenJavado/LinkFinder
