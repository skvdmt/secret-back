# Secret

### Translations
[Русский](./README_ru.md)

### Description
An application with a `/text` endpoint that requires Bearer token authorization.
Upon receiving the token, it verifies its validity by connecting via gRPC to the [access token validation](https://github.com/skvdmt/valid-access-back) application.
If the token is valid, it returns the secret.

### Download
```sh
git clone https://github.com/skvdmt/secret-back
```

## Links:
- [Docker image](https://hub.docker.com/r/skvdmt/secret-back-v1) — Docker image on docker hub
- [API](https://secret.skvdmt.ru/api/secret/v1/text) — Example /text
- [API Documentation](https://secret.skvdmt.ru/api/secret/v1/doc) — On Swagger-UI
- [Dmitry Skidanov](https://skvdmt.ru) — Author
