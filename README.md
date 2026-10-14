# Nintendo Badge Arcade Server

Community server implementation for the Nintendo 3DS application **Nintendo Badge Arcade**. This repository contains server code; it does not contain the game, Nintendo source code, badge artwork, or other Nintendo assets.

Nintendo ended the online functionality of Badge Arcade in April 2024. This project is an independent preservation and research effort and is not affiliated with or endorsed by Nintendo.

## What is included

- `authentication/`: a NEX authentication service. It derives the NEX password from the PID and configured HMAC secret. It initializes MongoDB for `pretendo_account.nexaccounts`; a PID lookup helper exists, but the current authentication flow does not call it.
- `secure/`: a NEX secure service with Badge Arcade handlers for maintenance status, DataStore operations, and Shop calls. DataStore objects are transferred through an S3-compatible object store; metadata and session data use PostgreSQL and MongoDB.
- `secure/third_party/nex-protocols-go/`: the local NEX protocol module referenced by `secure/go.mod` with a `replace` directive. Its upstream license and README are retained in that directory.

The source is based on the [PretendoNetwork Badge Arcade authentication server](https://github.com/PretendoNetwork/nintendo-badge-arcade-authentication) and [PretendoNetwork Badge Arcade secure server](https://github.com/PretendoNetwork/nintendo-badge-arcade-secure). Existing notices and the AGPL-3.0 license texts are retained. This repository also contains changes from the Protarium deployment snapshot; see [NOTICE.md](NOTICE.md).

## Build

Go 1.22 is used by the included Dockerfiles. Build each service from its module directory:

```sh
cd authentication
go test ./...
go build ./...

cd ../secure
go test ./...
go build ./...
```

The repository currently has no automated unit tests; `go test ./...` checks that packages compile.

## Configuration

Copy the relevant `.env.example` to `.env` in the service directory and fill in values for infrastructure you control. Never commit `.env` files or real database/object-store credentials.

The authentication listener binds UDP port `59400`. It advertises the secure service location using `SECURE_SERVER_LOCATION` and `SECURE_SERVER_PORT`. `PN_NEX_PASSWORD_SECRET` must be hexadecimal and decode to at least 32 bytes.

The secure listener binds UDP port `59401`. It requires MongoDB, PostgreSQL, and an S3-compatible endpoint. The `PN_NBA_CONFIG_S3_*` variables select the object-store bucket and key prefix used for Badge Arcade DataStore objects.

Both services use `KERBEROS_PASSWORD`. The authentication service initializes MongoDB through `MONGO_URI`. The secure service reads `MONGO_URI` and `DATABASE_URI` for its databases.

### Current implementation limits

These details are visible in the current source and should not be mistaken for complete game-service behavior:

- `GetRivToken` returns an empty string.
- `PostPlayLog` acknowledges the request but does not store its data.
- Badge Arcade `GetMetaByOwnerID` handles DataStore type `100` only.
- The secure service clears its `badge_arcade.sessions` collection during startup.
- Other NEX methods may return `NotImplemented` through the protocol library.

Review these behaviors and the database side effects before running the service against data you care about.

## License

The server code is distributed under the GNU Affero General Public License v3.0 (AGPL-3.0). See [LICENSE](LICENSE). The vendored protocol module retains its own license file. Third-party Go dependencies remain under their respective licenses.

## Official service status

Nintendo's [Badge Arcade service notice](https://www.nintendo.com/fr-fr/Assistance/Nintendo-3DS-et-Nintendo-2DS/Utilisation/Mise-a-jour-concernant-les-joueurs-de-Nintendo-Badge-Arcade-2455286.html) says that all Badge Arcade online functions ended in April 2024. Nintendo also notes that badges already placed on the HOME Menu remain available, while the application can no longer be used after service closure.
