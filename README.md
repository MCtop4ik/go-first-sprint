### Описание
Проект go-first-sprint предназначен для вычисления простейших операций.
Он поддерживает сложение, вычитание, умножение и деление, а также скобки для расставления приоритетов для операций.

### Запуск
1. Убедитесь, что вы не питонист, и у вас скачан Гоша
2. ```git clone https://github.com/MCtop4ik/go-first-sprint.git```
3. Перейдите в папку проекта и запустите через ```go run cmd/main.go```

### Проект запускается на ```localhost:8080```

### Эндпоинт

```
POST /api/v1/calculate
```

### Заголовки

- `Content-Type: application/json`

### Тело запроса

Пример:

```json
{
  "expression": "2+2*2"
}
```

### Ответы

1. **Успешный запрос**

   **Статус-код:** `200 OK`  
   **Пример ответа:**

   ```json
   {
     "result": "6"
   }
   ```

2. **Ошибка обработки выражения**

   **Статус-код:** `422 Unprocessable Entity`  
   **Пример ответа:**

   ```json
   {
     "error": "Expression is not valid"
   }
   ```

3. **Неподдерживаемый метод**

   **Статус-код:** `405 Method Not Allowed`  
   **Пример ответа:**

   ```json
   {
     "error": "Method not allowed"
   }
   ```

4. **Некорректное тело запроса**

   **Статус-код:** `400 Bad Request`  
   **Пример ответа:**

   ```json
   {
     "error": "Invalid request body
   }
   ```

---
Я проверял на Postman
```
{
    "expression": "(2+2)*3"
}
```
```
{
    "result": 12
}
```
```
{
    "expression": "(2+2++3)"
}
```
```
{
    "result": 7
}
```
```
{
    "expression": "(2+2)*3/0"
}
```
```
{
    "error": "Expression is not valid"
}
```

```
{
    "expression": "(2+2)*3/0
}
```
```
{
    "error": "Invalid request body"
}
```
```
{
    "expression": "(2+2*3"
}
```
```
{
    "error": "Expression is not valid"
}
```
```
{
    "expression": "(2+2*3)a"
}
```
```
{
    "error": "Expression is not valid"
}
```
