# LOTest - Простое веб-приложение на Go

LOTest — это минимальный REST API сервер на Go для управления задачами (tasks) с поддержкой асинхронного логирования через канал.

## Функционал

- **Создание задач**: `POST /tasks`
- **Получение списка задач**: `GET /tasks` с возможностью фильтрации по статусу (`?status=""`)
- **Получение задачи по ID**: `GET /task/{id}`
- Асинхронное логирование всех действий через канал
- Грейсфул-шатдаун сервера при получении сигнала завершения

---

## Установка

1. Клонируйте репозиторий:

```bash
git clone https://github.com/SemenShakhray/lotest.git
cd lotest
```

2. Соберите приложение:

```bash
go build -o lotest ./cmd/main.go
```

3. Запуск:

``` bash
./lotest
```

## Примеры запросов

1. Создание задачи:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Task","description":"This is a test"}'
```

2. Получение всех задач:

```bash
curl -X GET http://localhost:8080/tasks
```

3. Получение задач с фильтром по статусу: 

```bash
curl -X GET http://localhost:8080/tasks?status=pending
```

4. Получение задачи по ID:

```bash
curl -X GET http://localhost:8080/task/<id>
```


