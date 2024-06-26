#!/bin/bash
mysql -u $MYSQL_USER -p"$MYSQL_PASSWORD" --force < sql/schema/readingCopilotSchema.sql

mkdir -p "tls"
cd "tls"
if [ -f "cert.pem" ] && [ -f "key.pem" ]; then
    echo "Certificate files already exist. No need to regenerate."
    exit 0
fi
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost