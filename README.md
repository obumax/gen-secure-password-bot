[🇷🇺 Русский](#русский) | [🇬🇧 English](#english)

---

# gen-secure-password-bot

## Русский

### Описание

**gen-secure-password-bot** — это Telegram-бот для генерации надёжных паролей с возможностью гибкой настройки параметров. Бот поддерживает два языка: русский и английский.

---

### Используемый стек и библиотеки

- **Go 1.24.2**
- **Redis** — хранение пользовательских сессий
- **Docker** и **docker-compose** — для контейнеризации и запуска
- **go-telegram-bot-api** — работа с Telegram Bot API
- **go-i18n** — поддержка локализации
- **godotenv** — загрузка переменных окружения
- **go-redis** — клиент для Redis
- **GitHub Actions** - CI/CD для автоматической сборки, тестирования и публикации Docker-образа при пуше новых тегов в репозиторий
- **Codecov** - автоматический сбор и анализ покрытия тестами в CI/CD для контроля качества кода
---

### Особенности генерации пароля

- Длина пароля: от 4 до 35 символов
- Гарантируется наличие хотя бы одного символа из каждой выбранной категории
- Категории символов:
  - Прописные буквы (A-Z)
  - Строчные буквы (a-z)
  - Цифры (0-9)
  - Специальные символы (!@#№$;%^:&?*()-_=+[]{}<>.,/|`~)
- Возможность исключить похожие символы (il1O0)
- Перемешивание символов для повышения стойкости

---

### Возможности настройки

- Быстрая генерация пароля по умолчанию
- Кастомизация:
  - Выбор длины пароля
  - Включение/отключение категорий символов
  - Исключение похожих символов
- Генерация нового пароля с теми же параметрами одной кнопкой
- Сброс и начало заново

---

### Локализация

- Поддерживаются русский и английский языки
- Язык выбирается при первом запуске или командой `/lang`
- Все сообщения и кнопки локализованы

---

### Инструкция по запуску

#### 1. Клонируйте репозиторий

```bash
git clone https://github.com/obumax/gen-secure-password-bot.git
cd gen-secure-password-bot
```

#### 2. Установите и запустите Redis на сервере

```bash
sudo apt update
sudo apt install redis-server -y
sudo systemctl enable redis-server
sudo systemctl start redis-server
```

#### 3. Создайте файл .env

```bash
BOT_TOKEN=ваш_токен_бота
REDIS_PASSWORD=ваш_пароль_redis # если у вас включен пароль
```

#### 4. Настройте переменные окружения для подключения к Redis (на хост-машине)

В Docker Compose укажите подключение к Redis на сервере, а не к контейнеру через host.docker.internal (или 172.17.0.1 для Linux):

```bash
environment:
  - BOT_TOKEN=${BOT_TOKEN}
  - REDIS_URL=redis://:${REDIS_PASSWORD}@host.docker.internal:6379/0
```

или

```bash
environment:
  - BOT_TOKEN=${BOT_TOKEN}
  - REDIS_URL=redis://:${REDIS_PASSWORD}@172.17.0.1:6379/0
```

#### 5. Запустите через Docker Compose

```bash
docker-compose up --build
```

#### 6. Добавьте бота в Telegram и начните диалог

Ссылка на оригинальный бот
https://t.me/GenSecurePasswordBot

### Авто-деплой через GitHub Actions

Проект настроен на автоматическую сборку и загрузку Docker-образа при пуше нового тега в репозиторий.
С помощью GitHub Actions образ публикуется на DockerHub.
Добавьте теги:

```bash
git tag v1.0.0 && git push --tags
```

Подождите ~1-2 минуты — на DockerHub появится свежий релиз.
Файл workflow содержит всю логику CI/CD (сборка, тесты, пуш).

### Деплой через Docker на сервере

#### 1. Установите Docker на вашем сервере:

```bash
curl -fsSL https://get.docker.com | sh
```

#### 2. Установите/убедитесь, что Redis запущен на сервере:

```bash
sudo apt install redis-server -y
sudo systemctl enable redis-server
sudo systemctl start redis-server
```

#### 3. Запустите контейнер бота, указав подключение к Redis на сервере:

```bash
docker run -d \
  --name gen-secure-password-bot \
  -e BOT_TOKEN=ваш_токен_бота \
  -e REDIS_URL=redis://:ваш_пароль_redis@host.docker.internal:6379/0 \
  obumax/gen-secure-password-bot:latest
```

Для Linux используйте 172.17.0.1 вместо host.docker.internal:

```bash
docker run -d \
  --name gen-secure-password-bot \
  -e BOT_TOKEN=ваш_токен_бота \
  -e REDIS_URL=redis://:ваш_пароль_redis@host.docker.internal:6379/0 \
  obumax/gen-secure-password-bot:latest
```

#### 4. Используйте Watchtower для автообновления контейнера при выходе новой версии:

```bash
docker run -d \
  --name watchtower \
  --restart always \
  -v /var/run/docker.sock:/var/run/docker.sock \
  containrrr/watchtower \
  gen-secure-password-bot --interval 60
```

#### Контейнер работает 24/7, автозапуск обеспечен средствами Docker и Watchtower.

---

# gen_secure_password_bot

## English

### Description

**gen-secure-password-bot** is a Telegram bot for generating strong passwords with flexible customization options. The bot supports two languages: Russian and English.

---

### Stack and Libraries

- **Go 1.24.2**
- **Redis** — user session storage
- **Docker** and **docker-compose** — for containerization and running
- **go-telegram-bot-api** — Telegram Bot API integration
- **go-i18n** — localization support
- **godotenv** — environment variable loading
- **go-redis** — Redis client
- **GitHub Actions** — CI/CD for automatic build, testing, and DockerHub publishing on new tags
- **Codecov** — Automated code coverage collection and analysis in CI/CD to monitor code quality

---

### Password Generation Features

- Password length: from 4 to 35 characters
- At least one character from each selected category is guaranteed
- Character categories:
    - Uppercase letters (A-Z)
    - Lowercase letters (a-z)
    - Digits (0-9)
    - Special symbols (!@#№$;%^:&?*()-_=+[]{}<>.,/|`~)

Option to exclude similar characters (il1O0)
Characters are shuffled for extra security

### Customization Options

- Quick password generation with default settings
- Customization:
    - Choose password length
    - Enable/disable character categories
    - Exclude similar characters
    - Generate a new password with the same settings in one click
    - Reset and start over

---

### Localization

- Russian and English supported
- Language is selected on first launch /start or via /lang command
- All messages and buttons are localized

---

###  Quick Start Guide

#### 1. Clone the repository

```bash
git clone https://github.com/obumax/gen-secure-password-bot.git
cd gen-secure-password-bot
```

#### 2. Install and Start Redis on your server

```bash
sudo apt update
sudo apt install redis-server -y
sudo systemctl enable redis-server
sudo systemctl start redis-server
```

#### 3. Create a .env file

```bash
BOT_TOKEN=your_bot_token
REDIS_PASSWORD=your_redis_password # if enabled
```

#### 4. Configure environment variables for Redis connection

In your docker-compose.yml, make sure to connect to the host Redis (not a container):

```bash
environment:
  - BOT_TOKEN=${BOT_TOKEN}
  - REDIS_URL=redis://:${REDIS_PASSWORD}@host.docker.internal:6379/0
```

Or on Linux, use Docker’s gateway address:

```bash
environment:
  - BOT_TOKEN=${BOT_TOKEN}
  - REDIS_URL=redis://:${REDIS_PASSWORD}@172.17.0.1:6379/0
```

#### 5. Start the bot with Docker Compose

```bash
docker-compose up --build
```

#### 6. Add the bot in Telegram and start chatting

Link to the original bot
https://t.me/GenSecurePasswordBot

### Auto-deploy via GitHub Actions

When a new tag is pushed, Docker image is built and pushed to DockerHub automatically. Just run:

```bash
git tag v1.0.0 && git push --tags
```

Please wait 1-2 minutes for the latest release to appear on DockerHub.
The workflow file contains all the CI/CD logic (build, tests, push).

### Deploy and auto-update on server

#### 1. Install Docker on your server:

```bash
curl -fsSL https://get.docker.com | sh
```

#### 2. Make sure Redis is running on the server

As shown above.

#### 3. Run the bot container, connecting to the host Redis:

```bash
docker run -d \
  --name gen-secure-password-bot \
  -e BOT_TOKEN=your_telegram_bot_token \
  -e REDIS_URL=redis://:your_redis_password@host.docker.internal:6379/0 \
  obumax/gen-secure-password-bot:latest
```

On Linux, use 172.17.0.1 instead of host.docker.internal:

```bash
docker run -d \
  --name gen-secure-password-bot \
  -e BOT_TOKEN=your_telegram_bot_token \
  -e REDIS_URL=redis://:your_redis_password@172.17.0.1:6379/0 \
  obumax/gen-secure-password-bot:latest
```

#### 4. For auto-updates use Watchtower (Optional):

```bash
docker run -d \
  --name watchtower \
  --restart always \
  -v /var/run/docker.sock:/var/run/docker.sock \
  containrrr/watchtower \
  gen-secure-password-bot --interval 60
```

#### Container works 24/7, with automatic restart and updates via Docker & Watchtower tools.

---

## Лицензия / License

MIT License

Copyright (c) 2025 Maksim Obukhov https://github.com/obumax

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice