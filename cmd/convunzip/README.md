# convunzip: unzip file in unicode file name which zipped in local encoding

Install dependency package, iconv-go:

    $ go get github.com/djimenez/iconv-go

Usage:

    $ convunzip <-from encoding_of_zipfile> <-to encoding_of_env> zip_file ...
