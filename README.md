# Описание проекта 

    В дипломном проекте я написал веб-сервер, который реализует функциональность планировщика задач.
    Готовый планировщик задач позволяет вести календарь дел.
    Это задание на проверку и закрепление навыков по написанию веб-сервера, работе с REST API и базами данных. 

    Планировщик хранит задачи; каждая из них содержит дату дедлайна и заголовок с комментарием. Задачи могут повторяться по заданному правилу: например, ежегодно, через какое-то количество дней, в определённые дни месяца или недели. Если отметить такую задачу как выполненную, она переносится на следующую дату в соответствии с правилом. Обычные задачи при выполнении будут просто удаляться. 

    API содержит следующие операции:
    добавить задачу;
    получить список задач;
    удалить задачу;
    получить параметры задачи;
    изменить параметры задачи;
    отметить задачу как выполненную.

# Информация о выполненных заданиях

    Выполнены абсолютно все задания (со звездочкой в том числе)

# Структура проекта

    `main.go` (главный файл запуска приложения - ранее был в директории cmd, но так не проходили тесты гитхаба)

    `handlers (все обработчики запросов здесь):`
    - auth.go с обработчиком проверки аутентификации по паролю для входа в систему SignInHandler.
    - DeleteTask.go с обработчиком запроса на удаление задачи DeleteTaskHandler.
    - DoneTask.go с обработчиком DoneTaskHandler, который отмечает задачу выполненной и переносит, либо удаляет её.
    - GetTask.go с обработчиком GetTaskHandler, который нужен для получения полей задачи перед её редактированием.
    - HandlerSwitcherAndErrors.go , содержащий код обработки ошибок writeError и SwitcherHandler - переключатель между хендлерами PostTaskHandler, GetTaskHandler, PutTaskHandler, а также DeleteTaskHandler.
    - NextDate.go с обработчиком GetNextDateHandler для получения следующей даты.
    - PostTask.go с обработчиком для добавления задачи PostTaskHandler.
    - PutTask.go с обработчиком PutTaskHandler для передачи обновлённых полей задачи после её редактирования
    - TaskList.go с обработчиком GetTaskListHandler для получения списка задач, а также функцией isDate, проверяющей, является ли строка датой в нужном формате.

    `middleware:`
    - одноимённый с хендлером файл auth.go с функцией AuthMiddleware для проверки аутентификации пользователя через JWT-токен во время обработок запросов через SwitcherHandler, DoneTaskHandler и GetTaskListHandler.

    `models:`
    - файл models.go, который хранит в себе структуру планируемых задач планировщика.

    `rules:`
    - файл rules.go, в котором прописаны все вариации правил повторения задач. 
    Содержит в себе функции NextDate (не путать с GetNexDate из одноимённого обработчика NextDate), которая вычисляет следующую дату для задачи согласно заданным правилам повторения. Также содержит вспомогательные функции parseDays для парсинга дней месяца, parseMonths для парсинга месяцев, getLastDayOfMonth для получения последнего дня месяца, calculateNewDate для расчёта новой даты, getEarliestDate для получения ближайшей даты.

    `scheduler:`
    - файл scheduler.go, в котором осуществляется проверка существования базы данных, открытие или её создание (с обращением к файлу scheduler.sql)

    `tests`
    В этой директории находятся тесты для проверки API, которое реализовано в веб-сервере.

    `web`
    Эта директория содержит файлы фронтенда.

    `.env` 
    Файл, содержащий переменные окружения TODO_PASSWORD и TODO_PORT.

    `app` 
    Файл для успешного запуска приложения.

    `Dockerfile` 
    Файл, содержащий инструкцию для рабочей сборки контейнера.

    Обновлённые файлы `go.mod` и `go.sum` для корректной работы взаимосвязей приложения.

    База данных `scheduler.db`, хранящая в себе задачи планировщика.

    Файл `scheduler.sql`, содержащий в себе инструкцию для корректной постройки базы данных.

# Правильно запускать сервер:
go build -o app ./main.go
./app

# Команды для запуска тестов:
Шаг 1 - go test -run ^TestApp$ ./tests
Шаг 2 - go test -run ^TestDB$ ./tests
Шаг 3 - go test -run ^TestNextDate$ ./tests 
Шаг 4 - go test -run ^TestAddTask$ ./tests
Шаг 5 - go test -run ^TestTasks$ ./tests
Шаг 6.1 - go test -run ^TestTask$ ./tests
Шаг 6.2 - go test -run ^TestEditTask$ ./tests
Шаг 7.1 - go test -run ^TestDone$ ./tests
Шаг 7.2 - go test -run ^TestDelTask$ ./tests
Все тесты целиком - go test ./tests

# Примеры запуска тестов без кэша:
go test -count=1 -run ^TestTasks$ ./tests
go test -count=1 ./tests

# Получение куки через терминал:
curl -c cookies.txt -d '{"password":"12345"}' -H "Content-Type: application/json" -X POST http://localhost:7540/api/signin
TOKEN=$(grep 'token' cookies.txt | awk '{print $NF}')
echo $TOKEN

# Сборка образа:
docker buildx build -t todo-scheduler .

# Рабочая версия запуска контейнера:
docker run -d \
    -p 7540:7540 \
    -v ~\Desktop\go_final_project-main:/host \
    --name todo-app \
    todo-scheduler

# Также можно настроить параметры через переменные окружения при запуске:
docker run -d \
    -p 7540:7540 \
    -v ~\Desktop\go_final_project-main:/host \
    -e TODO_PORT=7540 \
    -e TODO_DBFILE=/host/scheduler.db \
    -e TODO_PASSWORD=12345 \
    --name todo-app \
    todo-scheduler

# Настройки тестов:
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = `` (разместить внутри свой)

# Переменные окружения в файле .env:
TODO_PASSWORD=12345
TODO_PORT=7540

# Страница с диалогом для ввода пароля 
http://localhost:7540/login.html
