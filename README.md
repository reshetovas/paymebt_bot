1. Установка шлюза для webhook
  ssh -R 80:localhost:8080 serveo.net
=> получаем адрес шлюза, который доступен из интернета
Адрес шлюза проксирует на local host

2. Передача адреса в ТГ
  curl --location 'https://api.telegram.org/bot{token}/setWebhook' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'url=https://56b6bea995243304beba6f6e1ac0fe85.serveo.net'

3. 3. Запуск приложения, используя либо telebot.v3
  go get gopkg.in/telebot.v3
