# TeleGospel бот

Небольшой бот для индивидуального прочтения христианского лекционария (плана чтения Библии).
Из провайдеров доступен только Bible Gateway.
Список переводов можно расширять свободно — достаточно вписать их в соответствующий `translations.json` с соблюдением формата. Вы можете выбрать любой вариант, лишь бы он был поддерживаем Bible Gateway.
Из планов доступен только Revised Common Lectionary. Пока нет необходимости в ином лекционарии — расширения не планируется.

## Инструкции по запуску на индивидуальной машине

Убедитесь, что у вас установлен компилятор Go (версии, совместимой с 1.20.2) и sqlite3.
Перед запуском создайте файл .env в корневой папке и введите туда пару `TG=%ваш_токен_телеграм_бота%`. ([О том, как его получить](https://core.telegram.org/bots/features#botfather))

Далее, находясь в корневой папке, исполните любую из следующих команд:
Для запуска:
```
go run .
```

Для компиляции исполняемого файла:
```
go build .
```

Учтите, что во время первого запуска бота будет произведена настройка базы данных sqlite.
