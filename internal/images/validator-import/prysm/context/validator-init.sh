#!/bin/sh
WALLET_DIR="/data/wallet"
rm -rf $WALLET_DIR
mkdir $WALLET_DIR

/validator wallet create --accept-terms-of-use --wallet-password-file=/keystore/keystore_password.txt --keymanager-kind=direct --wallet-dir="$WALLET_DIR"

tmpkeys="/tmpkeys"
mkdir -p ${tmpkeys}

for f in /keystore/validator_keys/keystore-*.json; do
    echo "Importing key ${f}"
    FILENAME=`echo ${key} | sed 's/.json//g'`
    cp "${f}" "${tmpkeys}"

    /validator accounts import \
        --accept-terms-of-use=true \
        --wallet-dir="$WALLET_DIR" \
        --keys-dir="${tmpkeys}" \
        --account-password-file="/keystore/${FILENAME}.txt" \
        --wallet-password-file=/keystore/keystore-password.txt

    filename="$(basename ${f})"
    rm "${tmpkeys}/${filename}"
done

rm -r ${tmpkeys}

echo "Imported all keys"