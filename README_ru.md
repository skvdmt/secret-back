# Секрет

### Translations
[En](./README.md)

### Описание
Приложение с эндтрипойнтом `/text` требующим токен доступа через Bearer авторизацию.
После получения токена проверяет его валидность соединяясь по gRPC с приложением [проверки валидности токена доступа](https://github.com/skvdmt/valid-access-back).
Если токен действителен, возвращается секрет.

### Скачать
```sh
git clone https://github.com/skvdmt/secret-back
```

## Ссылки:
- [Docker image](https://hub.docker.com/r/skvdmt/secret-back-v1) — Образ на docker hub.
- [Пример /text](https://secret.skvdmt.ru/api/secret/v1/text)
- [Swagger-UI Документация](https://secret.skvdmt.ru/api/secret/v1/doc)
- [Автор](https://skvdmt.ru) — Dmitry Skidanov.
