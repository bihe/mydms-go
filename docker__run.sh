#!/bin/sh
docker run -it -p 127.0.0.1:3000:3000 -v $PWD:/opt/mydms/etc mydms
