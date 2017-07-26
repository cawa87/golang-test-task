# golang-test-task
Тестовая задача для Golang разработчика


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

## Сборка и запуск

1. Установить [gb](https://getgb.io/)
1. Восстановить зависимости: `gb vendor restore`
1. Собрать проект: `gb build`
1. Запустить тесты: `gb test -v`
1. Собрать docker-контейнер: `docker built -t linkfetcher .`

## Замечания

1. По-хорошему, контейнер должен быть на базе Alpine. Для этого нужно, чтобы сборка и прогон тестов тоже проходили в докере. Не хочу копировать сюда свою баше-лапшу для этого, поскольку по-хорошему этим все равно должна заниматься билд-сервер.
    * Тут же возникает проблема того, что сейчас в контекст докера попадает все исходники, бинарники и проч-проч-проч.
1. Внешние данные пользователей нельзя без фильтрации передавать в этот сервис: для этого должа быть фильтрация частных IP-диапазонов. 
1. Cтоит группировать ссылки по серверам и/или использовать ограниченное количество воркеров для загрузки ссылок. Сделал так, т.к. детали реализации зависят от ньюансов реального использования, которого здесь нет.
1. В описании задачи не было ничего об обработке ошибок. Что, например, делать, если не получилось разрешить доменное имя? Для этих целей я добавил поле "error" в Meta и взял смелость подменять status на 500 в этом случае.
