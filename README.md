This is a tool to check the x509 certificates served up from PostgreSQL.

Inspiration, and at present most of the code, comes from
https://github.com/chr4/pg-check-cert. I wanted to be able to show the certs
that were returned, and also do full cert verification.

## Usage
```
% ./check_pg_cert postgresql://example.com:5432
0: Subject: CN=example.com
   Issuer: CN=InCommon RSA Server CA,OU=InCommon,O=Internet2,L=Ann Arbor,ST=MI,C=US
   NotAfter: 2020-04-19 23:59:59 +0000 UTC
1: Subject: CN=InCommon RSA Server CA,OU=InCommon,O=Internet2,L=Ann Arbor,ST=MI,C=US
   Issuer: CN=USERTrust RSA Certification Authority,O=The USERTRUST Network,L=Jersey City,ST=New Jersey,C=US
   NotAfter: 2024-10-05 23:59:59 +0000 UTC
2: Subject: CN=USERTrust RSA Certification Authority,O=The USERTRUST Network,L=Jersey City,ST=New Jersey,C=US
   Issuer: CN=AddTrust External CA Root,OU=AddTrust External TTP Network,O=AddTrust AB,C=SE
   NotAfter: 2020-05-30 10:48:38 +0000 UTC
```
