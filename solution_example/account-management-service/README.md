# Account management service

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) [![forthebadge](http://forthebadge.com/images/badges/built-with-love.svg)](http://forthebadge.com)

Микросервис для работы с балансами пользователей, резервации средств на специальном счёте и последующем признании выручки компании или возврата денег.
А также для получения данных для сводного отчёта по каждой услуге

Используемые технологии:
- PostgreSQL (в качестве хранилища данных)
- Docker (для запуска сервиса)
- Swagger (для документации API)
- Echo (веб фреймворк)
- golang-migrate/migrate (для миграций БД)
- pgx (драйвер для работы с PostgreSQL)
- golang/mock, testify (для тестирования)

Сервис был написан с Clean Architecture, что позволяет легко расширять функционал сервиса и тестировать его.
Также был реализован Graceful Shutdown для корректного завершения работы сервиса

# Getting Started

Для запуска сервиса с интеграцией с Google Drive необходимо предварительно:
- Зарегистрировать приложение в Google Cloud Platform: [Документация](https://developers.google.com/workspace/guides/create-project)
- Создать сервисный аккаунт и его секретный ключ: [Документация](https://developers.google.com/workspace/guides/create-credentials)
- Добавить секретный ключ в директорию secrets
- Добавить .env файл в директорию с проектом и заполнить его данными из .env.example,
указав `GOOGLE_DRIVE_JSON_FILE_PATH=secrets/your_credentials_file.json`
- Опционально, настроить `congig/config.yaml` под себя

Для запуска сервиса без интеграции с Google Drive достаточно заполнить .env файл,
оставив переменную `GOOGLE_DRIVE_JSON_FILE_PATH` пустой

# Usage

Запустить сервис можно с помощью команды `make compose-up`

Документацию после завпуска сервиса можно посмотреть по адресу `http://localhost:8080/swagger/index.html`
с портом 8080 по умолчанию

Для запуска тестов необходимо выполнить команду `make test`, для запуска тестов с покрытием `make cover` и `make cover-html` для получения отчёта в html формате

Для запуска линтера необходимо выполнить команду `make linter-golangci`

## Examples

Некоторые примеры запросов
- [Регистрация](#sign-up)
- [Аутентификация](#sign-in)
- [Пополнение счёта](#accounts-deposit)
- [Резервирование средств](#reservations-create)
- [Признание выручки](#reservations-revenue)
- [Возврат средств](#reservations-refund)
- [Получение истории операций пользователя](#operations-history)
- [Сводный отчёт по услугам с экспортом в Google Drive](#operations-report-link)
- [Сводный отчёт по услугам в формате csv файла](#operations-report-file)

### Регистрация <a name="sign-up"></a>

Регистрация сервиса:
```curl
curl --location --request POST 'http://localhost:8080/auth/sign-up' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"dannromm",
    "password":"Qwerty123!"
}'
```
Пример ответа:
```json
{
  "id": 1
}
```

### Аутентификация <a name="sign-in"></a>

Аутентификация сервиса для получения токена доступа:
```curl
curl --location --request POST 'http://localhost:8080/auth/sign-in' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"dannromm",
    "password":"Qwerty123!"
}'
```
Пример ответа:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY2MTUxMjEsImlhdCI6MTY2NjYwNzkyMSwiVXNlcklkIjoxfQ.c4jMWdmyXePtjTo_qrN6m9n-LQtHk_Q99OuzcpriYs4"
}
```

### Пополнение счёта <a name="accounts-deposit"></a>

Пополнение счёта пользователя на определённую сумму:
```curl
curl --location --request POST 'http://localhost:8080/api/v1/accounts/deposit' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY2MTUxMjEsImlhdCI6MTY2NjYwNzkyMSwiVXNlcklkIjoxfQ.c4jMWdmyXePtjTo_qrN6m9n-LQtHk_Q99OuzcpriYs4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": 1,
    "amount": 100
}'
```
Пример ответа:
```json
{
  "message": "success"
}
```

### Резервирование средств <a name="reservations-create"></a>

Резервирование средств по указанной услуге и номеру заказа:
```curl
curl --location --request POST 'http://localhost:8080/api/v1/reservations/create' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY2MTUxMjEsImlhdCI6MTY2NjYwNzkyMSwiVXNlcklkIjoxfQ.c4jMWdmyXePtjTo_qrN6m9n-LQtHk_Q99OuzcpriYs4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "account_id": 1,
    "product_id": 1,
    "order_id": 15,
    "amount": 10
}'
```
Пример ответа, с указанием id резервирования:
```json
{
  "id": 1
}
```

### Признание выручки <a name="reservations-revenue"></a>

Признание выручки по указанному резервированию:
```curl
curl --location --request POST 'http://localhost:8080/api/v1/reservations/revenue' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY2MTUxMjEsImlhdCI6MTY2NjYwNzkyMSwiVXNlcklkIjoxfQ.c4jMWdmyXePtjTo_qrN6m9n-LQtHk_Q99OuzcpriYs4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "account_id": 1,
    "product_id": 1,
    "order_id": 15,
    "amount": 10
}'
```
Пример ответа:
```json
{
  "message": "success"
}
```

### Возврат средств <a name="reservations-refund"></a>

В случае отказа от услуги можно вернуть средства на счёт пользователя:
```curl
curl --location --request POST 'http://localhost:8080/api/v1/reservations/refund' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY2MTUxMjEsImlhdCI6MTY2NjYwNzkyMSwiVXNlcklkIjoxfQ.c4jMWdmyXePtjTo_qrN6m9n-LQtHk_Q99OuzcpriYs4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "order_id": 15
}'
```
Пример ответа:
```json
{
  "message": "success"
}
```

### Получение истории операций пользователя <a name="accounts-history"></a>

Используется пагинация, по умолчанию возвращается последние 10 записей отсортированные по дате создания.

Для сортировке по сумме необходимо передать параметр `sort_type` со значением `amount`,
также можно явно указать сортировку по дате создания, передав значение параметра `date`

Для получения следующей страницы с данными необходимо передать параметр `offset` со значением `10` (по умолчанию 0)

Также можно указать количество записей на странице, передав параметр `limit` (максимальное значение 10, по умолчанию 10)
```curl
curl --location --request GET 'http://localhost:8080/api/v1/operations/history' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY2MTUxMjEsImlhdCI6MTY2NjYwNzkyMSwiVXNlcklkIjoxfQ.c4jMWdmyXePtjTo_qrN6m9n-LQtHk_Q99OuzcpriYs4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "account_id": 1
}'
```
Пример ответа:
```json
{
  "operations": [
    {
      "amount": 10,
      "operation": "refund",
      "time": "2022-10-24T11:06:06.896409Z",
      "product": "some product",
      "order": 15
    },
    {
      "amount": 10,
      "operation": "reservation",
      "time": "2022-10-24T11:06:02.431726Z",
      "product": "some product",
      "order": 15
    }
  ]
}
```

### Сводный отчёт по услугам с экспортом в Google Drive <a name="operations-report-link"></a>

Сервис формирует отчёт в разрезе каждой услуги, затем загружает его в Google Drive и возвращает ссылку на файл с открытым доступом на чтение:
```curl
curl --location --request GET 'http://localhost:8080/api/v1/operations/report-link' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY2MTUxMjEsImlhdCI6MTY2NjYwNzkyMSwiVXNlcklkIjoxfQ.c4jMWdmyXePtjTo_qrN6m9n-LQtHk_Q99OuzcpriYs4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "month": 10,
    "year": 2022
}'
```
Пример ответа:
```json
{
    "link": "https://drive.google.com/file/d/1rl91RS9n5l5kO9BDpHVQxpxYBegMzQC6/view?usp=sharing"
}
```

### Сводный отчёт по услугам в формате csv файла <a name="operations-report-file"></a>

Сервис формирует отчёт и возвращает его в виде csv файла:
```curl
curl --location --request GET 'http://localhost:8080/api/v1/operations/report-file' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjY2MTUxMjEsImlhdCI6MTY2NjYwNzkyMSwiVXNlcklkIjoxfQ.c4jMWdmyXePtjTo_qrN6m9n-LQtHk_Q99OuzcpriYs4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "month": 10,
    "year": 2022
}'
```
Пример ответа:
```csv
some product,30
```

# Decisions <a name="decisions"></a>

В ходе разработки был сомнения по тем или иным вопросам, которые были решены следующим образом:

1. При создании счёта стоит ли указывать id/uuid аккаунта в параметрах,
чтобы сервис мог использовать свои собственные id/uuid для идентификации и не хранить внешние id сервиса управления балансом.
> Решил, что не стоит, т.к. это увеличит время на разработку. Но возможно в будущем стоит добавить эту возможность
2. Как реализовать резервирование денег?
> Сначала была идея сделать отдельную сущность под отдельный счёт пользователя,
но потом решил, что достаточно хранить все резервирования в одной таблице и при необходимости делать по ней поиск
3. Как составлять отчёт?
> В задании указано, что нужно вернуть ссылку на отчёт. Была идея развернуть ftp сервер и хранить отчёты in-memory.
Всё-таки решил, что интеграция с Google Drive это интереснее, но в качестве альтернативы оставил возможность получить csv файл через http.
К тому же, архитектура позволяет в будущем легко переехать на внутреннее решение
4. Какой использовать тип пагинации?
> Каждый способ имеет свои преимущества и недостатки. Limit-Offset теряет в скорости работы и консистентности данных,
так как если между получениями смежных страниц, будут добавлены данные, то это приведёт к дублированию и потере записи. 
От использования курсора пришлось отказаться, так как на одну дату может быть множество операций,
тогда затруднительно получить отличные от первой страницы, нужно было бы увеличивать точность даты курсора, что усложнило бы разработку
