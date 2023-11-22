#!/usr/bin/env zsh

echo generate CA Root
openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -out ca.crt -keyout ca.key -subj "/C=SO/ST=Earth/L=MyLocation/O=MyOrganiz/OU=MyOrgUnit/CN=localhost"
echo generate Server Certs
openssl genrsa -out server.key 2048
echo generate server Cert Signing request
openssl req -new -key server.key -days 3650 -out server.csr -subj "/C=SO/ST=Earth/L=MyLocation/O=MyOrganiz/OU=MyOrgUnit/CN=localhost"
echo sign with CA Root
openssl x509  -req -in server.csr -extfile <(printf "subjectAltName=DNS:localhost") -CA ca.crt -CAkey ca.key -days 3650 -sha256 -CAcreateserial -out server.crt
echo generate Client Certs
openssl genrsa -out client.key 2048
echo generate client Cert Signing request
openssl req -new -key client.key -days 3650 -out client.csr -subj "/C=SO/ST=Earth/L=MyLocation/O=$O/OU=$OU/CN=localhost"
echo sign with CA Root
openssl x509  -req -in client.csr -extfile <(printf "subjectAltName=DNS:localhost") -CA ca.crt -CAkey ca.key -out client.crt -days 3650 -sha256 -CAcreateserial