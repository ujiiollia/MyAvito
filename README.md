# Тестовое задание для стажёра Backend
## Сервис баннеров
Реализован сервис, который позволяет показывать пользователям баннеры, в зависимости от требуемой фичи и тега пользователя, а также управлять баннерами и связанными с ними тегами и фичами.
## Общие вводные
**Баннер** — это документ, описывающий какой-либо элемент пользовательского интерфейса. Технически баннер представляет собой  JSON-документ неопределенной структуры. 
**Тег** — это сущность для обозначения группы пользователей; представляет собой число (ID тега). 
**Фича** — это домен или функциональность; представляет собой число (ID фичи).  
1. Один баннер может быть связан только с одной фичей и несколькими тегами
2. При этом один тег, как и одна фича, могут принадлежать разным баннерам одновременно
3. Фича и тег однозначно определяют баннер

Так как баннеры являются для пользователя вспомогательным функционалом, допускается, если пользователь в течение короткого срока будет получать устаревшую информацию.  При этом существует часть пользователей (порядка 10%), которым обязательно получать самую актуальную информацию. Для таких пользователей нужно предусмотреть механизм получения информации напрямую из БД.

## Запуск сревиса
Для поднятия докер контейнера с БД Postgresql выполнить команду:
```shell
make up
```
Для зпуска сервиса выполнить команду:
```shell
make run
```

## Условия
1. Использован этот [API](https://github.com/avito-tech/backend-trainee-assignment-2024/blob/main/api.yaml)
2. Тегов и фичей небольшое количество (до 1000), RPS — 1k, SLI времени ответа — 50 мс, SLI успешности ответа — 99.99%
3. Для авторизации доступов должны используется 2 вида токенов: пользовательский и админский.  Получение баннера может происходить с помощью пользовательского или админского токена, а все остальные действия могут выполняться только с помощью админского токена.  
> выбрана работа с Bearer токенами
4. Если при получении баннера передан флаг use_last_revision, необходимо отдавать самую актуальную информацию.  В ином случае допускается передача информации, которая была актуальна 5 минут назад.
> Если флаг use_last_revision = true, то будет получена новейшая информация из БД. Иначе будет взято сначение из кеша (Если оно там есть). 
5. Баннеры могут быть временно выключены. Если баннер выключен, то обычные пользователи не должны его получать, при этом админы должны иметь к нему доступ.
6. Добавлен метод удаления баннеров по фиче или тегу.

## Стек
- **Язык сервиса:** GO. 
- **База данных:** PostgreSQL
- **Деплой зависимостей и самого сервиса** использован Docker и Docker Compose.
