# matthewrdale.com [![Build Status](https://travis-ci.org/matthewdale/matthewrdale.com.svg?branch=master)](https://travis-ci.org/matthewdale/matthewrdale.com) [![GoDoc](https://godoc.org/github.com/matthewdale/matthewrdale.com?status.svg)](https://godoc.org/github.com/matthewdale/matthewrdale.com)
Source code for my personal homepage, matthewrdale.com

## Deploying via Docker
On build machine:
```bash
docker-compose build
docker save -o matthewrdalecom_web.tar matthewrdalecom_web:latest
scp matthewrdalecom_web.tar <remote>
scp docker-compose.yml <remote>
```

On run machine:
```bash
docker load -i matthewrdalecom_web.tar
docker-compose down --remove-orphans
docker-compose up -d
```
