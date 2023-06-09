#!/bin/bash

#This is for CICD workflow that has the key pulled from the secrets repo and store as an env var
#Generates server keys and signs them with the CA
if [[ -z ${CA_SUB} ]]; then
  echo "No CA Subject ENV var, making one for testing"
  CA_SUB="/C=US/O=theOTW/OU=Engineering/CN=k8srelay"
fi
echo "Using Subject $CA_SUB"
if [[ -z ${CA_KEY} ]]; then
  echo "No CA KEY ENV var, making one for testing"
  openssl genrsa -out myCA.key 2048
  openssl req -x509 -new -nodes -key myCA.key -sha256 -days 1825 -out myCA.pem -subj $CA_SUB
else
  echo "Using CA from Secrets"
  echo "$CA_CERT" | base64 -d > myCA.pem
  echo "$CA_KEY" | base64 -d > myCA.key
fi
echo "Generating Private Key"
openssl genrsa -out out/k8srelay.key 2048
echo "Generating Private CSR"
openssl req -new -key out/k8srelay.key -out out/k8srelay.csr  -subj $CA_SUB
echo "Generating Private Signing"
openssl x509 -req -in out/k8srelay.csr -CA myCA.pem -CAkey myCA.key -CAcreateserial -out out/k8srelay.crt -days 825 -sha256 -extfile x509.config
cp myCA.pem out/myCA.pem
#rm  myCA.key
