[🇷🇺 Русский](#русский) | [🇬🇧 English](#english)

---

# gen-secure-password-bot

## Русский

### Описание

**gen-secure-password-bot** — это Telegram-бот для генерации надёжных паролей с возможностью гибкой настройки параметров. Бот поддерживает два языка: русский и английский.

---

### Используемый стек и библиотеки

- **Go 1.24+**
- **Redis** — хранение пользовательских сессий
- **Docker** и **docker-compose** — для контейнеризации и запуска
- **go-telegram-bot-api** — работа с Telegram Bot API
- **go-i18n** — поддержка локализации
- **godotenv** — загрузка переменных окружения
- **go-redis** — клиент для Redis

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

2. Создайте файл .env

```bash
BOT_TOKEN=ваш_токен_бота
REDIS_PASSWORD=ваш_пароль_redis
```

3. Запустите через Docker Compose

```bash
docker-compose up --build
```

4. Добавьте бота в Telegram и начните диалог

Ссылка на оригинальный бот
https://t.me/GenSecurePasswordBot

---

# gen_secure_password_bot

## English

### Description

**gen-secure-password-bot** is a Telegram bot for generating strong passwords with flexible customization options. The bot supports two languages: Russian and English.

---

### Stack and Libraries

- Go 1.24+
- Redis — user session storage
- Docker and docker-compose — for containerization and running
- go-telegram-bot-api — Telegram Bot API integration
- go-i18n — localization support
- godotenv — environment variable loading
- go-redis — Redis client

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

### How to Run

1. Clone the repository

```bash
git clone https://github.com/obumax/gen-secure-password-bot.git
cd gen-secure-password-bot
```

2. Create a .env file

```bash
BOT_TOKEN=ваш_токен_бота
REDIS_PASSWORD=ваш_пароль_redis
```

3. Start with Docker Compose

```bash
docker-compose up --build
```

4. Add the bot in Telegram and start a conversation

Link to the original bot
https://t.me/GenSecurePasswordBot

---

### Лицензия / License

MIT License

Copyright (c) 2025 Maksim Obukhov https://github.com/obumax

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice