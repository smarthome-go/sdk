## Changelog for v0.19.1

- **BUGFIX**: The auth mode `CookieToken` was broken and is now fixed
- The cookies were neither received nor appended which caused a `invalid credentials` error on any connection function after the inital login
