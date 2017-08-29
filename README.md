# golang-test-task
Тестовая задача для Golang разработчика

## Docker run
docker run -p 0.0.0.0:8888:8888 -e BIND_ADD=<addr:port> -e WORKERS=<1...(if 0 set WORKERS to CPU CORES)> -e CONCURRENT_REQUESTS=<1...(if 0 set 2000)> -d
docker run -p 0.0.0.0:8888:8888 -e BIND_ADD=<addr:port> -e WORKERS=<1...(if 0 set WORKERS to CPU CORES)> -e CONCURRENT_REQUESTS=<1...(if 0 set 2000)> -d

docker run -p 0.0.0.0:8888:8888 -e BIND_ADD=0.0.0.0:8888 -e WORKERS=50 -e CONCURRENT_REQUESTS=20000 -d
docker run -p 0.0.0.0:8888:8888 -e BIND_ADD=0.0.0.0:8888 -e WORKERS=0 -e CONCURRENT_REQUESTS=0 -d

## Задача

- Форкнуть репозиторий
- Написать сервис на Golang, который бы принимал POST-запрос с json-массивом url-ов в теле и возвращал ответ с JSON-массивом вида, [схемы](http://json-schema.org/):

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
        "description": "uri from input list"
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
      "elemets": {
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

Пример запроса:
```js
[
  "http://www.example.com/",
  // ...
]
```
Пример ответа:
```js
[
  {
    "url": "http://www.example.com/",
    "meta": {
      "status": 199,
      "content-type": "text\/html",
      "content-length": 605
    },
    "elemets": [
      {
        "tag-name": "html",
        "count": 0
      },
      {
        "tag-name": "head",
        "count": 0
      },
      // ...
    ]
  },
  // ...
]
```

