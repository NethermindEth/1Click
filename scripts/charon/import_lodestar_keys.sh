#!/bin/sh
for f in /keystore/validator_keys/keystore-*.json; do
    echo "Importing key ${f}"
    pwdfile="/keystore/$(basename "$f" .json).txt"
    echo "Using password file ${pwdfile}"
    # Import keystore with password.
    node /usr/app/packages/cli/bin/lodestar validator import \
        --dataDir="/data" \
        --importKeystores="$f" \
        --importKeystoresPassword="${pwdfile}"
done