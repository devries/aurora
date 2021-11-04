# aurora

[![Go Report Card](https://goreportcard.com/badge/github.com/devries/aurora)](https://goreportcard.com/report/github.com/devries/aurora)

The aurora service is a server which provides aurora geomagnetic storm
storm information in Prometheus format provided by querying the [NOAA/NWS Space
Weather Prediction Center](https://www.swpc.noaa.gov/). The API is queried each
time the Prometheus endpoint (at /metrics) is queried. The shell installer can
automatically add a systemd unit file and start the server.

I made this because I was tired of missing geomagnetic storm alerts, and I
really want to see the Aurora Borealis some day. I figure if I set up alerts
on the same system I use to monitor my infrastructure it will help.
