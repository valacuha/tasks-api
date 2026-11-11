# Tasks API (To-Do REST API на Go)

Микросервис для управления списком задач с хранением данных в оперативной памяти (in-memory) и обменом данными в формате JSON.

## Запуск проекта

Требования: **Go 1.20+**

```bash
# Клонирование и запуск сервера
go run ./cmd/server
```
Сервер запустится по адресу `http://localhost:8080`.

## Эндпоинты API

| Метод | Путь | Описание | Успешный статус | Коды ошибок |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/tasks` | Получить список всех задач | `200 OK` | `405` |
| `POST` | `/tasks` | Создать новую задачу | `201 Created` | `400`, `405`, `500` |
| `GET` | `/tasks/{id}` | Получить задачу по ID | `200 OK` | `400`, `404`, `405` |
| `PUT` | `/tasks/{id}` | Обновить задачу целиком | `200 OK` | `400`, `404`, `405`, `500` |
| `DELETE` | `/tasks/{id}` | Удалить задачу | `204 No Content` | `400`, `404`, `405`, `500` |
| `GET` | `/health` | Проверка работоспособности сервера | `200 OK` | `405` |

## Примеры запросов и ответов (curl)

### 1. Создание задачи (`POST /tasks`)
**Запрос:**
```bash
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Изучить Go", "done": false}'
```
**Ответ (`201 Created`):**
```json
{
  "id": 1,
  "title": "Изучить Go",
  "done": false,
  "created_at": "2026-11-11T10:00:00Z"
}
```

### 2. Получение списка задач (`GET /tasks`)
**Запрос:**
```bash
curl -i http://localhost:8080/tasks
```
**Ответ (`200 OK`):**
```json
[
  {
    "id": 1,
    "title": "Изучить Go",
    "done": false,
    "created_at": "2026-11-11T10:00:00Z"
  }
]
```

### 3. Обновление задачи (`PUT /tasks/{id}`)
**Запрос:**
```bash
curl -i -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Изучить Go и сдать проект", "done": true}'
```
**Ответ (`200 OK`):**
```json
{
  "id": 1,
  "title": "Изучить Go и сдать проект",
  "done": true,
  "created_at": "2026-11-11T10:00:00Z"
}
```

### 4. Удаление задачи (`DELETE /tasks/{id}`)
**Запрос:**
```bash
curl -i -X DELETE http://localhost:8080/tasks/1
```
**Ответ:** `204 No Content` (пустое тело).

### 5. Пример ответа с ошибкой (`400 Bad Request` / `404 Not Found`)
**Запрос:**
```bash
curl -i http://localhost:8080/tasks/999
```
**Ответ (`404 Not Found`):**
```json
{
  "error": "задача не найдена"
}
```
