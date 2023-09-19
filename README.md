# Nordlayer

This is a test script to limit bandwitch to any IP using [go-tc](https://github.com/florianl/go-tc) library.

## Build command

```
go build -o limit main
```

## Run commands

To set limit to default `80.249.99.148` IP with default `100Kbit` bandwidth through default interface `wlp5s0` run: 
```
sudo ./limit create
```

To reset limit run:
```
sudo ./limit delete
```

To get parameters info:
```
sudo ./limit help
```