Запущен в Yandex Cloud - http://62.84.113.234:8080/swagger/index.html

Запуск в docker:
- выполнить команду `docker build . -t kafkatest:latest`
- выполнить команду `docker-compose up`

Запуск локально:
`go run cmd/kafka-test/main.go cmd/kafka-test/flags.go`

Запуск с .env файлом (по умолчанию используется ./.env в основной папке): `go run cmd/kafka-test/main.go cmd/kafka-test/flags.go --config-path path/to/.env`

Информация по флагам: `go run cmd/kafka-test/main.go cmd/kafka-test/flags.go --help`

Переменные окружения:

| Название  | Описание | Default |
| ------------- | ------------- | ------------- |
| MIGRATIONS_FOLDER | Папка с миграциями для бд | ./migrations |
| PORT | Порт на котором запустится сервер | 8080 |
| HOST_ADDR | Хост адрес сервера, к которому будет обращаться сваггер для получения документации | localhost |
| DB_PROVIDER | Провайдер базы данных (поддерживается только postgres) | postgres |
| DB_URL | Url к базе данных | localhost:5432 |
| DB_USER | Username для доступа к бд | admin |
| DB_PASSWORD | Пароль для доступа к бд | admin |
| DB_DBNAME | Название базы данных | test |
| KAFKA_MESSAGES_TOPIC | Название топика для обработки сообщений | messages |
| KAFKA_ERR_MESSAGES_TOPIC | Название топика для ошибок обработки сообщений | err-messages |
| DB_APPEND_MIGRATIONS | Нужно ли попробовать применить миграции | true |
| DB_FORCE_RECREATE | Нужно ли пересоздать бд | false |
| KAFKA_BROKERS | Url к брокерам кафки через,запятую | localhost:9092 |
| LOGS_PATH | путь к файлу с логами | ./logs.json |
