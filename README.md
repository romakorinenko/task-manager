![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/romakorinenko/task-manager/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/romakorinenko/task-manager)](https://goreportcard.com/report/github.com/romakorinenko/task-manager)
[![Coverage Status](https://coveralls.io/repos/github/romakorinenko/task-manager/badge.svg?branch=create_task_manager)](https://coveralls.io/github/romakorinenko/task-manager?branch=create_task_manager)

# Программа для управления задачами

### Описание:
Task Manager - это многопользовательский веб-сервис для ведения рабочих задач. В настоящий момент существует функционал 
для двух ролей пользователей: USER и ADMIN. Вход осуществляется по логину и паролю.

## USER может:
- видеть список задач заасайненых на себя;
- создавать новые задачи для себя;
- редактировать свои задачи;
- удалять свои задачи.

Таким образом, пользователь может организовать свою работу. Логин и пароль для тестового юзера: user:user

## ADMIN может:
- использовать весь функционал, доступный для USER роли;
- также видеть список всех задач, созданных всеми пользователями, и редактировать и удалять все задачи, 
а также создавать задачи на любого пользователя;
- просматривать prometheus-метрики приложения по пути `http://localhost:8080/metrics`;

Для администраторов существует админка в виде swagger, доступной по пути `http://localhost:8080/swagger/index.html`.
Админка предоставляет дополнительный функционал:
- создание и редактирование пользователей;
- блокировка пользователей (формальная, функционал приложения в данный момент доступен и заблокированным пользователям)
- получение списка задач по статусу или приоритету;
- получение списка всех пользователей системы.

Логин и пароль для тестового администратора: admin:admin

### Развертывание
Развертывание сервиса должно осуществляется с docker compose.

`docker-compose up -d`

### Тестирование
Написаны юнит тесты на core логику приложения:
`go test -race -count 100 -v -tags=unit ./...`

Написаны интеграционные тесты:
`go test -v -tags=integration ./...`

### TODO list
- создание функционала для заблокированных пользователей 
- разобраться с http методами, преимущественно в настоящий момент используются GET и POST запросы, поскольку html
без использования js не позволяет выполнять запросы с другими методами
- дополнить swagger описаниями передаваемых параметров
- добавить графану с дэшбордами
