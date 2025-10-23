# ai-tools-manager

## API docs
### /chat
#### POST /chat
Создает новый chat

Response: 
```json
{
  "chat_id": num
}
```
#### GET /chat/{id} (TODO: добавить пагинацию в будущем)
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
