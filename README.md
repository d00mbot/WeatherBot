# Subscription bot

Subsciption bot - is a telegram bot that provides daily weather forecast for its subscribers.
The weather forecast is carried out due to the geolocation provided by the user by subscribing to the bot.

[OpenWeatherMap.org](https://openweathermap.org/api) is used as a data provider.

[MongoDB](https://www.mongodb.com/) is used to store data.

[Docker](https://www.docker.com/) is required to run the application.

### To run the application:
The application is run from docker containers, one for the bot and another for the local mongo database.
In docker-compose.yml set the required environment variables:
  >  `- WEATHER_TOKEN=YOUR_WEATHER_API_TOKEN`

  >  `- TELEGRAM_TOKEN=YOUR_TELEGRAM_BOT_TOKEN`

  >  `- MONGO_URI=YOUR_MONGO_URI (if empty - default value: "mongodb://localhost:27017/")`

  >  `- LOG_LEVEL=PRODUCTION (if empty - default value: "DEBUG")`

  >  `- RESPONSES_CONFIG_PATH=/config (if empty - default value for local usage: "../config")`
    
Run the `docker-compose up` command to build and run containers. Bot container runs on port `8080` and mongoDB runs on port `27100`.

Go version: `1.19`
