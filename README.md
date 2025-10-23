# ai-tools-manager

## API docs
### /chat
#### POST /chat
Создает новый chat

Response: 
```json
{
  "chat_id": 0
}
```

#### GET /chat/list
Получает список всех чатов

Response: 
```json
{
  "chats": [
    {
      "chat_id": 0
    }
  ]
}
```


### /message
#### GET /message/{chat_id} (TODO: добавить пагинацию в будущем)
Получает список всех сообщений чата 

Response:
```json
{
  "messages": [
    {
      "role": "assistant",
      "content": "example text"
    }
  ]
}
```

#### POST /message/{chat_id)
Отправлет сообщениие в чат


### /lib
#### GET /lib/list
Получение всех установленных библиотек

Response: 
```json
[
  {
    "name": "example name",
    "description": "example description"
    "git_url": "http://example-link.com"
  }
]
```
#### POST /lib
Cкачивает библиотеку (сохрвняет в бд) и создает knative контейнер с service

Request:
```json
{
  "git_url": "http://example-link.com"
}
```

#### DELETE /lib/{lib_url}
Удвляет библиотеку, контейнер, service

