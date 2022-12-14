# Introduction

# Key Concepts

## Service Account

A Service Account is an account used to identify services (non-humans) to Google.
A Service Account has a list of Service Account Keys, which can be used to authenticate to Google.

## Service Account Keys

A Service Account Key is a public/private keypair generated by Google.  Google retains the public
key, while the customer is given the private key.  The private key can be used to [sign JWTs and
authenticate Service Accounts to Google](https://developers.google.com/identity/protocols/OAuth2ServiceAccount#authorizingrequests).
