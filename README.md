# RestrictedPassportService

How to use:
1. `cd src`
2. `go build`
3. Update passports DataBase: `RestrictedPassportService update`
4. Start API:  `RestrictedPassportService api`

## Api usage
`curl http://localhost/passportApi/check?passport=PASSPORT_NUMBER`

will returns `{"passportBanned":false,"success":true}`
