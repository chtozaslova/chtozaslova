# ChtoZaSlova
## Go library

This is the ChtoZaSlova Go library, command-line tool, and API server.

It converts between latitude/longitude coordinates and what3words addresses.

ChtoZaSlova is free and open source software published under the Creative
Commons CC0 license. Please copy and distribute widely.

## Using the Go library

```
    import "chtozaslova/chtozaslova"

    ...

    words, err := chtozaslova.LatLon2Words(37.234332, -115.806657)
    fmt.Println(words)
    lat, lon, err := chtozaslova.Words2LatLon("joyful.nail.harmonica")
    fmt.Printf("%.6f %.6f\n", lat, lon)
```

See `example/main.go` for an example.

`LatLon2Words` and `Words2LatLon` are the only functions considered part of the API, but you are welcome to use
the other exported functions if you find them useful.

## Using the command-line tool

The tool in `cmd/chtozaslova/` is a command-line interface for converting between latitude/longitude coordinates and what3words addresses.

Help text is as follows:

```
usage: chtozaslova [-i] [INPUT...]

Options:
  -i    read from stdin

If not -i, you should specify at least one INPUT in the form lat,lon or words. It is perfectly acceptable to mix lat,lon and word inputs in a single invocation.

Examples:
    chtozaslova -i < words.txt
    chtozaslova joyful.nail.harmonica
    chtozaslova 37.234332,-115.806657

Output is of the form "INPUT[tab]OUTPUT" with one line per INPUT, where OUTPUT will be either the converted INPUT, or an error message.
```

## Using the API server

The tool in `cmd/chtozaslovad/` is an HTTP server that provides a HTTP API for converting between latitude/longitude coordinates and what3words addresses.
It listens on port 8081.

Parameters are passed in the HTTP query string and responses are formatted as JSON.

Make requests like:

```
$ curl http://localhost:8081/api/convert-to-coordinates?words=joyful.nail.harmonica
{"words":"joyful.nail.harmonica","language":"en","coordinates":{"lat":37.234328,"lon":-115.806657}}
$ curl http://localhost:8081/api/convert-to-3wa?coordinates=37.234328,-115.806657
{"words":"joyful.nail.harmonica","language":"en","coordinates":{"lat":37.234328,"lon":-115.806657}}
```

In the event of an error, the `error` key will exist in the returned JSON hash, and the value of this key will be an error message:

```
$ curl http://localhost:8081/api/convert-to-coordinates?words=errormessage
{"error":"Error decoding words"}
```

## Development

The files of the library are in `*.go` in the root directory.

For the command line application, see `cmd/chtozaslova/main.go`.

For the API server, see `cmd/chtozaslovad/main.go`.
