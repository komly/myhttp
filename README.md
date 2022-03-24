# myhttp

## Description
Simple tool which makes http requests and prints the address of the request along
with the MD5 hash of the response.

## Installation
go get github.com/komly/myhttp

## Usage
```
myhttp -parallel=3 google.com http://amazon.com https://example.com
google.com cc12bba3cf5dd0a86a06ef55e05649b1
http://amazon.com 5f0844b121447b714c56aa5789ef98d2
https://example.com 84238dfc8092e5d9c0dac8ef93371a07
```
## License
MIT