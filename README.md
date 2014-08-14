bonza-go
========

A simple reverse proxy written in Go (as a learning exercise)

## Usage


    bonza [filename |  port] [-quiet] prefix1=>resource1 [prefix2=>resource2]...]


If invoked without arguments, reads configuration from a .bonza file in the same directory.
If invoked with a single argument that is a filename, reads configuration from the named file.
Otherwise reads arguments from the command line.

## File Format

Delimited lines, the first line is the port to bind on with the remaining lines consist of proxy mapping
expressions.  The # character may be used as a comment character and will result in everything after it being ignored.

e.g :

    #port
    8080
    #map /g to google
    /g=http://wwww.google.com

# Proxy Mapping Expressions

Expressions of the form uri_prefix=>proxy_url

Route requests with the given uri prefix to the proxied resource url.

The matching uri prefix of the request is removed and the remainder of the request appended to the proxy_url before 
invoking the request.

e.g. given a mapping expresssion /google=>http://google.net/ then a request to /google/blah?q=foo will result in a
request to http://google.net/blah?q=foo

If more than one uri_prefix matches the incoming request, then the longest such match is chosen.

Logging:
The system logs requests in the following format : {timestamp} {source_ip} {method} {uri} -> {proxied_url} > {response_code} in {request_time} ms
Logging may be disabled by passing -quiet on the command line or in the file.