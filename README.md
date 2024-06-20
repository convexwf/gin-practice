# Go-Gin mTLS Template

This is a template for Go-Gin server with mTLS authentication.

## Generate mSSL Certificates

```shell
# Create directories for certificates
mkdir -p tls

# Generate CA
$ openssl req -x509 -newkey rsa:4096 -keyout tls/ca.key -out tls/ca.crt -days 365
Enter PEM pass phrase:
Verifying - Enter PEM pass phrase:
...
Country Name (2 letter code) [AU]:CN
State or Province Name (full name) [Some-State]:guangdong
Locality Name (eg, city) []:guangzhou
Organization Name (eg, company) [Internet Widgits Pty Ltd]:my.inc
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:
Email Address []:
```

You can decide the above information by yourself.

```shell
# Edit the openssl.cnf file
$ vim tls/openssl.cnf
# Add the following content
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
prompt = no

[req_distinguished_name]
countryName = CN
stateOrProvinceName = guangdong
localityName = guangzhou
organizationName = my.inc
commonName =

[v3_req]
subjectAltName = @alt_names

[alt_names]
DNS.1 = mtls-server
```

Then we can generate the server and client certificates with the following commands.

```shell
# Generate server certificate
$ openssl req -newkey rsa:2048 -nodes -keyout tls/server.key -out tls/server.csr -subj "/CN=mtls-server" --config tls/openssl.cnf
# Sign the server certificate by CA
$ openssl x509 -req -in tls/server.csr -out tls/server.crt -CA tls/ca.crt -CAkey tls/ca.key -CAcreateserial -days 365 -extensions v3_req -extfile tls/openssl.cnf
Certificate request self-signature ok
subject=CN = mtls-server
Enter pass phrase for tls/ca.key:
```

The password is the one you set when generating the CA certificate.

```shell
# Generate client certificate
$ openssl req -newkey rsa:2048 -nodes -keyout tls/client.key -out tls/client.csr -subj "/CN=mtls-client" -config tls/openssl.cnf
# Sign the client certificate by CA
$ openssl x509 -req -in tls/client.csr -out tls/client.crt -CA tls/ca.crt -CAkey tls/ca.key -CAcreateserial -days 365 -extensions v3_req -extfile tls/openssl.cnf
Certificate request self-signature ok
subject=CN = mtls-client
Enter pass phrase for tls/ca.key:
```

Now you can check the directory `tls` to see the generated certificates.

```shell
$ tree tls
tls
├── ca.crt
├── ca.key
├── ca.srl
├── client.crt
├── client.csr
├── client.key
├── openssl.cnf
├── server.crt
├── server.csr
└── server.key
```

## deploy client and server

```shell
$ docker-compose up -d
# docker-compose up -d --build
```
