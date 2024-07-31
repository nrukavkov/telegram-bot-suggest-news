# Telegram Bot Suggest News

This project is a Telegram bot that allows users to suggest news for a channel. The bot is built with Go using the `telebot` library and includes rate limiting functionality to prevent spam. 

## Features

- Users can suggest news to a Telegram channel.
- Rate limiting to prevent spam (configurable).
- Dockerized for easy deployment.

## Prerequisites

- Go 1.17 or later
- Docker
- Docker Compose
- A Telegram bot token (You can create a bot and get the token from [BotFather](https://t.me/botfather))
- Admin ID (Your Telegram user ID)

## Variables

- `ADMIN_ID` - your personal ID in Telegram
- `RATE_LIMIT` - limit of message per time interval. For example, `5:60` allows users suggest 5 news per 1 minute (60s)
- `TELEGRAM_BOT_TOKEN` - telegram bot token

## Hosting

You can run your news bot without buying some hosting. The trick is using github actions. For achive it follow these steps...

### 1. Create a new personal bot 

Open https://t.me/botfather and create a new bot. Select a name and save bot token. Remember, one bot - one app - one channel. If you need this bot for 2 channels or more - follow steps from 1st till end for each telegram channels. Save token and open your Telegram Bot, than press `Start`. You will get nothing from already created bot, but it is ok, because you didn't run your application for catch requests.

### 2. Fork this repo

It will create your own copy of code and you might set and use your personal secrets and variables. 

### 3. Set var TELEGRAM_BOT_TOKEN in repo

Follow settings/secrets/actions and add `New repository secret` and set `TELEGRAM_BOT_TOKEN` with token from step 1.

### 4. Build your bot app

Follow `Actions` page and enable actions. Open action `Build and Push Image` and press `Run workflow` for building your own instance of Telegram Bot. Wait until success finished.

### 5. First run

Follow `Actions` page and open action `Hosting` and press `Run workflow` for first run.

### 6. Get your ADMIN_ID

`ADMIN_ID` is your personal ID in telegram. But do not worry it is easy to know. Open bot in Telegram and write command `/getid`. You will get an ADMIN_ID. 

### 7. Final setup and run

Follow settings/variables/actions and add `New repository variable`. 

- `ADMIN_ID` - from previous step
- `RATE_LIMIT` - recommended value is `5:60`. It allows 1 user send not more then 5 messages in 60 seconds. But it is up to you and you can set `0:0` for unlimited suggested news.

Now we need to restart our application. Follow `Actions` page and and press `Run workflow` again. It will cancel current run application and create a newone with all needed params. 

### 8. Test

Open bot in telegram and send some message. You will receive a message back. So you can make post in Telegram Channel with link to your bot and PIN it!

## Local Development

### 1. Clone the repository

```bash
git clone https://github.com/nrukavkov/telegram-bot-suggest-news.git
cd telegram-bot-suggest-news
```

### 2. Create a .env.local file
Create a .env.local file in the root of your project with the following content:

```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
ADMIN_ID=your_admin_id
RATE_LIMIT=5:60
```

### 3. Build and Run the Bot
```bash
go build -o telegram-bot-suggest-news
./telegram-bot-suggest-news

