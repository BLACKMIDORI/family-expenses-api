# Family Expenses (API)

This project is the API of the [Family Expenses](https://github.com/BLACKMIDORI/family-expenses-mobile) application

# Setting up

## Environment Variables
### Postgres
`DB_URL`

`DB_NAME`

`DB_USER`

`DB_PASSWORD`

### Google Identity API
`GOOGLE_ALLOWED_ISSUERS=https://accounts.google.com`_comma-separated_

`GOOGLE_ALLOWED_AUDIENCES`_comma-separated_ : 
The OAuth Client Id of clients(web, mobile)

### JWT Settings
`JWT_ISSUER`: Issuer can be any string(usually is a url)

`JWT_AUDIENCES`_comma-separated_: Allowed clients for generating JWT access tokens

`JWT_SIGNING_KEY_PEM`: RSA private key in PEM format

`JWT_SIGNING_PUBLIC_KEY_PEM`: Equivalent RSA public key in PEM format
