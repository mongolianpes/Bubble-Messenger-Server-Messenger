# Bubble Messenger
## Назначение
Микросервис для обмена сообщениями между пользователями в проекте [Bubble](https://github.com/mongolianpes/Bubble-Messenger-Server) <br>
Через сообщения пользователей также передаються специальные сообщения инициализации аудио диалога

## Функционал
- Обмен сообщениями между пользователями

## RPC запросы
```
service MessengerService {
  rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);
  rpc CheckChats(CheckChatsRequest) returns (CheckChatsResponse);
  rpc DelMessage(DelMessageRequest) returns (DelMessageResponse);
}
```
Подробнее в файле messenger/proto/messenger.proto

## Переменные окружения
Для корректной работы необходимы переменные окружения:<br>
DB_HOST - адрес БД<br>
DB_PORT - порт БД<br>
DB_USER - пользователи БД<br>
DB_PASSWORD - пароль пользователя БД<br>
DB_NAME - название БД

## Порты
Для корректной работы необходимо открыть порт 8086 для взаимодействия сервисов через gRPC
