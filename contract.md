# Спецификация
## GET /task/{id} - получает таску по айди
### response:
    
```
200 - запрос выполнен успешно 
    body:{
        id: number,
        title: string, 
        text: string,
        author: string,
        urgent: bool
    }
404  - таска не найдена
429 - слишком много запросов с единицу времени 
500 - На сервере произошла ошибка, в результате которой он не может успешно обработать запрос.
```
**requestBody - отсутствует**
---

## GET /task - получает список всех таск
### response:

```
200 - запрос выполнен корректно
    body:{
        array: Task[]
    }
429 - слишком много запросов с единицу времени 
500 - На сервере произошла ошибка, в результате которой он не может успешно обработать запрос.
```
**requestBody - отсутствует**
---

## POST /task - создает таску по форме
### response:
   
```
201 - запрос успешно обработан, новая таска создана
    body:{
        answer: Task
    }
429 - слишком много запросов с единицу времени 
500 - На сервере произошла ошибка, в результате которой он не может успешно обработать запрос.
507 - запрос не выполнен, потому что серверу не удалось сохранить данные
```

### request:

```
body:{
    title: string, 
    text: string,
    author: string,
    urgent: bool
}
```
---


## PATCH /task/{id} - обновляет значение полей таски
### response:
   
```
200 - запрос успешно обработан, данные таски с указанным айди обновлены
400 - в теле данные указаны в неправильном формате или синтаксисе
404 - таска не нашлась по данному айдишнику
429 - слишком много запросов с единицу времени 
500 - На сервере произошла ошибка, в результате которой он не может успешно обработать запрос.
507 - запрос не выполнен, потому что серверу не удалось сохранить данные 
```

### request:
   
```
body:{
    title: abobe?:string, 
    text: abobe?:string,
    author: abobe?:string,
    urgent: abobe?:bool
}
```
---

## PUT /task/{id} -  удаляет текущую таску, создает новую с тем же айдишником и с полями из тела запроса
### response:
   
```
200 - запрос успешно обработан, данные таска с указанным айди обновлена
201 - тк по заданому айди не было таски, она была создана, причем успешно
400 - в теле данные указаны в неправильном формате или синтаксисе
429 - слишком много запросов с единицу времени 
500 - На сервере произошла ошибка, в результате которой он не может успешно обработать запрос.
507 - запрос не выполнен, потому что серверу не удалось сохранить данные 
```
### request:
   
```
body:{
    title: string, 
    text: string,
    author: string,
    urgent: bool
}
```
---
## DELETE /task/{id} - удаляет таску по айдишнику
### response:
```   
200 - запрос успешно обработан, таска удалена
404 - таска по данному айдишнику не найдена или не существует 
500 - На сервере произошла ошибка, в результате которой он не может успешно обработать запрос.
```
**requestBody - отсутствует** 

## Class:
```
Task:{
    id: number,
    title: string, 
    text: string,
    author: string,
    urgent: bool
}
```