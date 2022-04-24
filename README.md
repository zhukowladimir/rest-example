# Разработка REST веб-сервиса
> Приложение писалось в качестве домашнего задания из курса "Сервис-ориентированные архитектуры" ПМИ ВШЭ.

## Постановка задачи

**Цель:** на языке высокого уровня (Java, C#, Python и др. – на выбор обучающегося) реализовать REST веб-сервис и клиент для него, обеспечивающий сбор и предоставление клиентам возможностей работы с ресурсами игры «SOA-мафия».  

**Задача:** 
1. Реализовать REST-сервис, который предоставляет возможность добавления, просмотра, редактирования и удаления следующей информации по профилю игрока: никнейм, аватар, пол, email. Должна быть обеспечена возможность получения профиля как отдельного игрока, так и перечня игроков.
2. Реализовать сбор и представление посредством REST-сервиса статистики по игрокам и проведенным сессиям игр. Статистика по игроку должна генерироваться в виде PDF-документа, содержащего информацию: профиль игрока, количество сессий, количество побед, количество поражений, общее время в игре. Статистика по игроку должна генерироваться по асинхронному запросу, возвращающему URL, по которому в дальнейшем будет доступен PDF-документ со сгенерированной статистикой. Генерация статистики должна быть реализована на основе паттерна «Очередь заданий».
3. *Организовать регистрацию, авторизацию и разграничение прав пользователей к редактированию профиля игрока в REST-сервисе*
4. *Организовать единый механизм регистрации и авторизации пользователей с использованием JWT как для основного приложения, так и для REST-сервиса*

## Запуск приложения 
### Docker-compose
> Необходимо наличие `docker-compose`

```
docker-compose build
docker-compose up
```

## Работа с приложением 

### Как отправлять запросы

Например, через библиотеку ![httpie](https://httpie.io/):
```
# Получить информацию по пользователю с идентификатором c32240c6-4f4c-4100-a401-826f70faeb4a
http GET http://127.0.0.1:8080/players/c32240c6-4f4c-4100-a401-826f70faeb4a/stats

# Добавить пользователя username
http POST http://127.0.0.1:8080/players Username='username' Sex='male' Email='username@gmail.com' Avatar='https://i.imgur.com/n0XluaE.jpg'

# Получить файл user_stats.pdf со статистикой пользователя с идентификатором c32240c6-4f4c-4100-a401-826f70faeb4a
http --print=b GET http://127.0.0.1:8080/players/c32240c6-4f4c-4100-a401-826f70faeb4a/stats/466466649.pdf > user_stats.pdf
```

## API
### POST /players
Добавить нового игрока
#### Example Input
```
{
  "Username": "username",
  "Avatar": "https://i.imgur.com/n0XluaE.jpg",
  "Sex": "male",
  "Email": "username@gmail.com"
}
```
#### Example Response
```
{
    "Avatar": "https://i.imgur.com/07ddZW7.jpeg",
    "Email": "username@gmail.com",
    "ID": "a633a2c6-8734-4202-95ec-bdd12061611b",
    "Sex": "male",
    "Username": "username"
}
```

### GET /players
Получить информацию обо всех игроках в базе данных
#### Example Input
```
```
#### Example Response
```
[
    {
        "Avatar": "https://i.imgur.com/n0XluaE.jpg",
        "Email": "ch4d@gmail.com",
        "ID": "c32240c6-4f4c-4100-a401-826f70faeb4a",
        "Sex": "male",
        "Username": "gigachad"
    },
    {
        "Avatar": "https://i.imgur.com/07ddZW7.jpeg",
        "Email": "username@gmail.com",
        "ID": "a633a2c6-8734-4202-95ec-bdd12061611b",
        "Sex": "male",
        "Username": "username"
    }
]
```

### GET /players/{id}
Получить информацию об игроке с данным id
#### Example Input
```
```
#### Example Response
```
{
    "Avatar": "https://i.imgur.com/n0XluaE.jpg",
    "Email": "ch4d@gmail.com",
    "ID": "c32240c6-4f4c-4100-a401-826f70faeb4a",
    "Sex": "male",
    "Username": "gigachad"
}
```

### PUT /players/{id}
Обновить информацию об игроке с данным id
#### Example Input
```
{
  "Sex": "god"
}
```
#### Example Response
```
{
    "Avatar": "https://i.imgur.com/n0XluaE.jpg",
    "Email": "ch4d@gmail.com",
    "ID": "c32240c6-4f4c-4100-a401-826f70faeb4a",
    "Sex": "god",
    "Username": "gigachad"
}
```

### DELETE /players/{id}
Удалить игрока с данным id
#### Example Input
```
```
#### Example Response
```
```

### PUT /players/{id}/stats 
Обработать информацию о прошедешей сессии у игрока с данным id и обновить его статистику
#### Example Input
```
  "IsWin": "true"
  "Duration": "13"
```
#### Example Response
```
```

### GET /players/{id}/stats
Получить статистику игрока с данным id
#### Example Input
```
```
#### Example Response
```
127.0.0.1:8080/players/c32240c6-4f4c-4100-a401-826f70faeb4a/stats/995420084.pdf
```

### GET /players/{id}/stats/{filename}
Получить filename.pdf файл с отчетом по стастике игрока с данным id
#### Example Input
```
```
#### Example Response
```
Bianry data
```
