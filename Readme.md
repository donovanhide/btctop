btctop
======

Top for Bitcoin exchanges! Uses the [bitcoincharts.com API](http://bitcoincharts.com/about/markets-api/) for data and displays the data in sortable and filterable format in a terminal window. Inspired by [goxtool](http://prof7bit.github.io/goxtool/)!

Latest Binaries
---------------

Available on the [release page](https://github.com/donovanhide/btctop/releases).


Installation from source
------------------------

[Install Go](http://golang.org/doc/install) and then (for OS X and linux):

```bash
export GOPATH=/path/to/gopath
go get github.com/donovanhide/btctop
./bin/btctop
```

TODO
----

* Dynamic window resizing
* Graphing of last 24 hours of trades per exchange
* Reconcile figures...