# checksum-calc

[![Build Status](https://travis-ci.org/northbright/checksum-calc.svg?branch=master)](https://travis-ci.org/northbright/checksum-calc)
[![Go Report Card](https://goreportcard.com/badge/github.com/northbright/checksum-calc)](https://goreportcard.com/report/github.com/northbright/checksum-calc)

checksum-calc is a program which calculates file MD5, SHA-1 checksums. It's written in [Golang](http://golang.org).

#### Usage

    -f string
            File to calculate MD5 / SHA-1 hash. Ex: -f='my-cd.iso'
    Usage:
    checksum-calc -f=<file>
    Ex: checksum-calc -f='my-cd.iso'

#### Packaging

```
root@5e0a846c6cf3:/data# fpm -s dir -t rpm -n "static-checksum-utils" -v 1.0 ./bin                                                                                                             
Created package {:path=>"static-checksum-utils-1.0-1.x86_64.rpm"}
```

#### Screenshot

![Screenshot on Windows](images/screenshot-1.png)

#### License
* [MIT License](./LICENSE) 
