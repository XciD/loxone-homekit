# loxone-homekit
[![version](https://img.shields.io/badge/status-beta-orange.svg)](https://github.com/XciD/loxone-homekit)
[![Build Status](https://travis-ci.org/XciD/loxone-homekit.svg?branch=master)](https://travis-ci.org/XciD/loxone-homekit)
[![Go Report Card](https://goreportcard.com/badge/github.com/XciD/loxone-homekit)](https://goreportcard.com/report/github.com/XciD/loxone-homekit)
[![codecov](https://codecov.io/gh/XciD/loxone-homekit/branch/master/graph/badge.svg)](https://codecov.io/gh/XciD/loxone-homekit)
[![Pulls](https://img.shields.io/docker/pulls/xcid/loxone-homekit.svg)](https://hub.docker.com/r/xcid/loxone-homekit)
[![Layers](https://shields.beevelop.com/docker/image/layers/xcid/loxone-homekit/latest.svg)](https://hub.docker.com/r/xcid/loxone-homekit)
[![Size](https://shields.beevelop.com/docker/image/image-size/xcid/loxone-homekit/latest.svg)](https://hub.docker.com/r/xcid/loxone-homekit)


Loxone Homekit Integration in go

Work in progress


```
docker build -t homekit
docker run -v config.yaml:/config.yaml -v /config:/config --net=host -it homekit
```
