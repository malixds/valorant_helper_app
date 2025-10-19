# Valorant Teams Telegram Bot

Telegram бот с веб-приложением для управления командами в Valorant.

## Функциональность

- Создание и управление командами
- Вступление в команды
- Просмотр участников команды
- Telegram Mini App интерфейс

## Технологии

- **Backend**: Go, Gin, GORM, PostgreSQL
- **Frontend**: HTML, CSS, JavaScript, Telegram Web App API
- **Bot**: Telegram Bot API

## Установка и запуск

### 1. Установка зависимостей

```bash
go mod tidy
```

### 2. Настройка PostgreSQL

Создайте базу данных:

```sql
CREATE DATABASE valorant_db;
CREATE USER valorant_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE valorant_db TO valorant_user;
```

### 3. Настройка конфигурации

Отредактируйте `config.env` файл:

```
TELEGRAM_BOT_TOKEN=your_bot_token_here
WEBHOOK_URL=http://localhost:8080
PORT=8080

# PostgreSQL Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=valorant_user
DB_PASSWORD=your_password_here
DB_NAME=valorant_db
DB_SSLMODE=disable
```

### 4. Получение токена бота

1. Найдите @BotFather в Telegram
2. Создайте нового бота командой `/newbot`
3. Скопируйте полученный токен в конфигурацию

### 5. Запуск приложения

#### Локальная разработка (режим polling)

```bash
go run main.go
```

В этом режиме бот будет опрашивать Telegram API для получения новых сообщений. Webhook не нужен.

#### Продакшн (режим webhook)

1. Разверните приложение на сервере с публичным IP
2. Обновите `WEBHOOK_URL` на ваш домен
3. Запустите приложение

## API Endpoints

### Пользователи
- `GET /api/users/:telegram_id` - Получить пользователя
- `POST /api/users` - Создать пользователя
- `PUT /api/users/:telegram_id` - Обновить пользователя

### Команды
- `GET /api/teams` - Получить все команды
- `GET /api/teams/:id` - Получить команду по ID
- `POST /api/teams` - Создать команду
- `POST /api/teams/:team_id/join/:telegram_id` - Вступить в команду
- `POST /api/teams/leave/:telegram_id` - Покинуть команду

## Структура проекта

```
valorant_app/
├── main.go                 # Основной файл приложения
├── config/
│   └── config.go          # Конфигурация
├── database/
│   └── database.go        # Подключение к БД
├── models/
│   ├── user.go           # Модель пользователя
│   └── team.go           # Модель команды
├── handlers/
│   ├── user.go           # Обработчики пользователей
│   └── team.go           # Обработчики команд
├── bot/
│   └── bot.go            # Telegram бот
├── web/
│   ├── templates/
│   │   └── index.html    # HTML шаблон
│   └── static/
│       ├── css/
│       │   └── style.css # Стили
│       └── js/
│           └── app.js    # JavaScript логика
└── config.env            # Конфигурация
```

## Режимы работы

### Polling режим (для разработки)
- Бот опрашивает Telegram API каждые несколько секунд
- Не требует публичного URL
- Подходит для локальной разработки
- Автоматически включается при `WEBHOOK_URL=http://localhost:8080`

### Webhook режим (для продакшна)
- Telegram отправляет обновления на ваш сервер
- Требует публичный URL
- Более эффективен для продакшна
- Включается при указании реального домена в `WEBHOOK_URL`

## Развертывание

### Локальная разработка

1. Установите PostgreSQL
2. Создайте базу данных
3. Настройте `config.env`
4. Запустите: `go run main.go`

### Продакшн

1. Разверните приложение на сервере
2. Настройте SSL сертификат
3. Обновите `WEBHOOK_URL` на ваш домен
4. Запустите приложение

## Использование

1. Найдите вашего бота в Telegram
2. Нажмите "Start" или используйте кнопку меню
3. Откроется веб-приложение
4. Создайте команду или вступите в существующую
5. Управляйте участниками команды

## Разработка

Для разработки рекомендуется:

1. Использовать polling режим (localhost)
2. Настроить hot reload для фронтенда
3. Использовать логирование для отладки

## Лицензия

MIT License
