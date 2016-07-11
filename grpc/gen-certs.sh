#!/bin/sh

# Small script to generate a CA and two signed certs (one for the server and one
# for the client).

set -ex

mkdir -p certs
cd certs

# Generate the CA key and self-cert.
openssl ecparam -name prime256v1 -out param
openssl req -x509 -nodes -days 365 -sha256 -subj '/CN=root' -newkey ec:param -keyout ca.key -out ca.pem

# Generate the server key and cert.
openssl req -new -nodes -sha256 -subj '/CN=server' -newkey ec:param -keyout server.key |
openssl x509 -req -CA ca.pem -CAkey ca.key -CAcreateserial -out server.pem

# Generate the client key and cert.
openssl req -new -nodes -sha256 -subj '/CN=client' -newkey ec:param -keyout client.key |
openssl x509 -req -CA ca.pem -CAkey ca.key -CAcreateserial -out client.pem

openssl verify -CAfile ca.pem server.pem client.pem
rm -f param