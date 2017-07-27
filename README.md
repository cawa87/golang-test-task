# golang-test-task
Тестовая задача для Golang разработчика


## Задача

- Форкнуть репозиторий
- Написать сервис на Golang, который бы принимал POST-запрос с json-массивом url-ов в теле и возвращал ответ с JSON-массивом вида:

```
{
  "type": "array",
  "items": {
    "type": "object",
    "required": ["url", "meta"],
    "properties": {
      "url": {
        "type": "string",
        "format": "uri",
        "description": "uri from input list",
      },
      "meta": {
        "type": "object",
        "required": ["status"],
        "properties": {
          "status": {
              "type": "integer",
              "description": "Response status of this uri"
          },
          "content-type": {
               "type": "string",
               "description": "In case of 2XX response status, value of mime-type part of Content-Type header (if exists)"
          },
          "content-length": {
                "type": "integer",
                "minimum": 0,
                "description": "In case of 2XX response status, length of response body (be careful, response could be chunked)."
          }
        }
      },
      "elements": {
        "type": "array",
        "description": "In case of 2XX response status, \"text\/html\" content type and non-zero content length, list of HTML-tags, occured.",
        "items": {
          "type": "object",
          "required": ["tag-name", "count"],
          "properties": {
            "tag-name": {"type": "string"},
            "count": {
              "type": "integer",
              "minimum": 1,
              "description": "Number of times, the current tag occures in response"
            }
          }
        }
      }
    }
  }
}
```

- Завернуть сервис в docker-контейнер.

## Результат

Сборка бинарника

`CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags "-s" -a -installsuffix cgo -o service`

Сборка образа

`docker build -t service .`

Запуск

`docker run -it -p 8080:8080 service`

Проверка

`curl -X POST -H 'Content-Type:application/json' --data '["http://mail.ru", "http://ya.ru", "google.com"]' http://localhost:8080 | jq`