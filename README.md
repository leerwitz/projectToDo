# Project To do
# Содержание:
- [Краткое описание](#краткое-описание)
- [Установка локально](#установка-локально)
- [Создание и запуск Docker контейнера](#создание-и-запуск-docker-контейнера)
- [Использование](#использование) 

# Краткое описание 
Проект представляет собой сайт на котором можно оставлять какие то задачи как напоминание. 
Сама задача разделена на несколько частей:
- Заголовок
- Текст
- Автор
- Срочное оно или не срочное
- ID
Задание можно создать, удалить и редактировать. Также реализован поиск 
по заголовку и по ID (прим. поиск по ID всегда следует начинать с)
Сервер реализован по Restful api, для задач реализован CRUD.

# Установка локально

На устройстве должен быть установлен Docker версии 
## Установка Docker:
### На Linux(ubuntu)
#### Шаг 1: Удалите старые версии Docker (если они установлены)
    ```bash
    sudo apt-get remove docker docker-engine docker.io containerd runc
    ```

#### Шаг 2: Обновите индекс пакетов и установите необходимые пакеты для установки через HTTPS
    ```bash
    sudo apt-get update
    sudo apt-get install \
        ca-certificates \
        curl \
        gnupg \
        lsb-release
    ```

#### Шаг 3: Добавьте официальный GPG ключ Docker
    ```bash
    echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    ```


#### Шаг 4: Добавьте репозиторий Docker в источники APT
    ```bash
    echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    ```


#### Шаг 5: Обновите индекс пакетов и установите Docker Engine
    ```bash
    sudo apt-get update
    sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    ```


#### Шаг 6: Проверьте, что Docker установлен корректно
    ```bash
    sudo docker run hello-world
    ```
Если вы увидите сообщение "Hello from Docker!", значит Docker установлен и работает правильно.
---

### На Windows 

#### Шаг 1: Проверьте системные требования

Перед установкой Docker убедитесь, что ваша система соответствует следующим требованиям:
- Windows 10 64-бит: Pro, Enterprise, или Education (Build 16299 или выше) или Windows 11.
- Включённая функция Hyper-V и контейнеры Windows.

#### Шаг 2: Скачайте Docker Desktop

1. Перейдите на официальный сайт Docker и скачайте установочный файл Docker Desktop: [Docker Desktop для Windows](https://desktop.docker.com/win/stable/Docker%20Desktop%20Installer.exe).

#### Шаг 3: Установите Docker Desktop

1. Запустите скачанный установочный файл `Docker Desktop Installer.exe`.
2. Следуйте инструкциям установщика. По умолчанию будут установлены все необходимые компоненты.

#### Шаг 4: Настройте Docker Desktop

1. После завершения установки запустите Docker Desktop.
2. При первом запуске Docker Desktop предложит вам включить необходимые функции Windows, такие как Hyper-V и контейнеры. Нажмите "Ok", чтобы разрешить изменения.
3. Перезагрузите компьютер, если это потребуется.

#### Шаг 5: Проверьте установку Docker

1. После перезагрузки запустите Docker Desktop.
2. Откройте командную строку (CMD) или PowerShell и выполните команду:
   ```bash
   docker --version
3. Выполните тестовый запуск контейнера:
    ```bash
    docker run hello-world
    ```
Если вы увидите сообщение "Hello from Docker!", значит Docker установлен и работает правильно.
---
### На MacOS

#### Шаг 1: Проверьте системные требования

Перед установкой Docker убедитесь, что ваша система соответствует следующим требованиям:
- macOS версии 10.15 или выше.
- Аппаратная виртуализация должна быть включена в настройках BIOS.

#### Шаг 2: Скачайте Docker Desktop

1. Перейдите на официальный сайт Docker и скачайте установочный файл Docker Desktop для Mac: [Docker Desktop для Mac](https://desktop.docker.com/mac/stable/Docker.dmg).

#### Шаг 3: Установите Docker Desktop

1. Откройте скачанный файл `Docker.dmg`.
2. Перетащите Docker в папку "Applications".

#### Шаг 4: Запустите Docker Desktop

1. Откройте Docker из папки "Applications".
2. При первом запуске Docker Desktop может запросить доступ к системным настройкам для установки необходимых компонентов. Нажмите "Ok" и введите свой пароль администратора.

#### Шаг 5: Настройте Docker Desktop

1. Следуйте инструкциям на экране для завершения настройки Docker Desktop.
2. После завершения настройки вы увидите значок Docker в строке меню.

#### Шаг 6: Проверьте установку Docker

1. Откройте терминал (Terminal).
2. Введите команду:
   ```bash
   docker --version
3. Выполните тестовый запуск контейнера:
    ```bash
    docker run hello-world
    ```
Если вы увидите сообщение "Hello from Docker!", значит Docker установлен и работает правильно.

## Создание и запуск Docker контейнера

Зайдите в корневую директорию проекта(projectToDo) и введите в терминал: 
    ```bash
    sudo docker-compose up -d --build
    ```
Чтобы удалить контейнер введите:
    ```bash
    sudo docker-compose down
    ```
Чтобы посмотреть логи контейнера используйте:
    ```bash
    sudo docker-compose logs
    ```

## Использование
После сборки контейнер зайдите на http://localhost в браузере,
теперь вы можете начать пользоваться приложением.