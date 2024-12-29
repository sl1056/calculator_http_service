# Веб-сервис калькулятора

Простой веб-сервис для вычисления арифметических выражений, предоставленных пользователем. Сервис предоставляет API-эндпоинт, на который можно отправить выражение и получить результат в формате JSON.

## Возможности

- Поддерживает базовые арифметические операции: `+`, `-`, `*`, `/`.
- Работает с выражениями, содержащими скобки и десятичные числа.
- Проверяет входные данные и возвращает соответствующие коды ошибок HTTP для некорректных выражений или ошибок сервера.

## Документация API

### Эндпоинт

`POST /api/v1/calculate`

### Тело запроса

```json
{
    "expression": "арифметическое выражение"
}
```

### Ответ
## Успешный результат (HTTP 200)
``` json
{
    "result": "вычисленный результат"
}
```
## Некорректное выражение (HTTP 422)
``` json
{
    "error": "Expression is not valid"
}
```
## Внутренняя ошибка сервера (HTTP 500)
``` json
{
    "error": "Internal server error"
}
```

### Веб-сервис калькулятора
Запрос:
```json
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
Успешный ответ:
```json
{
    "result": "6.000000"
}
```

Запрос:
```json
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2-"
}'
```
Ответ:
```json
{
    "error":"Expression is not valid"
}
```

```json
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2a"
}'
```
Ответ:
```json
{
    "error":"Expression is not valid"
}
```

## Установка и запуск
# Клонируйте репозиторий:

```json
git clone https://github.com/sl1056/calculator_http_service.git
cd calculator_http_service
```


# Запустите сервис:

```json
go run calc_service.go
```

# Сервис будет доступен по адресу http://localhost:8080.

## Тестирование сервиса
# Вы можете протестировать сервис с помощью curl. Примеры корректных выражений:
``` json
"2+2*2"

"3.5/(2-0.5)"

```

