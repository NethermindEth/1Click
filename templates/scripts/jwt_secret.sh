#!/bin/sh
openssl rand -hex 32 | tr -d "\n" > {{.Path}}/jwtsecret