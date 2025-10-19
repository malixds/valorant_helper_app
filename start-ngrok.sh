#!/bin/bash

# Скрипт для запуска приложения с ngrok

echo "🚀 Запуск Valorant Bot с ngrok..."

# Проверяем, установлен ли ngrok
if ! command -v ngrok &> /dev/null; then
    echo "❌ ngrok не установлен. Установите его:"
    echo "   brew install ngrok/ngrok/ngrok"
    echo "   или скачайте с https://ngrok.com/"
    exit 1
fi

# Проверяем, запущен ли PostgreSQL
if ! pg_isready -h localhost -p 5432 &> /dev/null; then
    echo "❌ PostgreSQL не запущен. Запустите его:"
    echo "   brew services start postgresql"
    exit 1
fi

echo "📡 Запуск ngrok туннеля на порту 8080..."
ngrok http 8080 --log=stdout > ngrok.log 2>&1 &
NGROK_PID=$!

# Ждем, пока ngrok запустится
sleep 3

# Получаем ngrok URL
NGROK_URL=$(curl -s http://localhost:4040/api/tunnels | grep -o '"public_url":"[^"]*' | grep -o 'https://[^"]*' | head -1)

if [ -z "$NGROK_URL" ]; then
    echo "❌ Не удалось получить ngrok URL"
    kill $NGROK_PID
    exit 1
fi

echo "✅ Ngrok URL: $NGROK_URL"

# Обновляем конфигурацию
echo "🔧 Обновление конфигурации..."
sed -i.bak "s|# NGROK_URL=.*|NGROK_URL=$NGROK_URL|" config.env

echo "🤖 Запуск бота..."
go run main.go &
BOT_PID=$!

echo "✅ Бот запущен!"
echo "📱 Ngrok URL: $NGROK_URL"
echo "🌐 Веб-приложение: $NGROK_URL"
echo "🔗 Webhook: $NGROK_URL/webhook"
echo ""
echo "Для остановки нажмите Ctrl+C"

# Функция для очистки при выходе
cleanup() {
    echo "🛑 Остановка приложения..."
    kill $BOT_PID 2>/dev/null
    kill $NGROK_PID 2>/dev/null
    # Восстанавливаем конфигурацию
    mv config.env.bak config.env 2>/dev/null
    exit 0
}

trap cleanup SIGINT SIGTERM

# Ждем завершения
wait
