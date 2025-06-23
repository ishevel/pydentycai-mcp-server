# Процедура установки PydanticAI для Windows

Эта инструкция описывает шаги по установке Python и PydanticAI в операционной системе Windows.

## 1. Установка Python

Если у вас еще не установлен Python, следуйте этим шагам:

1.  **Скачайте установщик Python:**
    Перейдите на официальный сайт Python: [https://www.python.org/downloads/windows/](https://www.python.org/downloads/windows/)
    Скачайте последнюю стабильную версию Python 3 (рекомендуется 3.8 или выше). Выберите "Windows installer (64-bit)" для большинства современных систем.

2.  **Запустите установщик:**
    Найдите скачанный файл `.exe` и запустите его.

3.  **Важные опции установки:**
    *   **Обязательно установите флажок "Add Python X.Y to PATH"** (где X.Y — версия Python) в самом низу окна установщика. Это позволит запускать Python из командной строки.
    *   Выберите "Install Now" для стандартной установки или "Customize installation" для выбора компонентов и пути установки. Для большинства пользователей "Install Now" будет достаточно.

4.  **Завершите установку:**
    Дождитесь завершения установки. После успешной установки вы увидите сообщение "Setup was successful".

5.  **Проверьте установку Python:**
    Откройте командную строку (нажмите `Win + R`, введите `cmd` и нажмите `Enter`).
    Введите следующие команды и убедитесь, что они выводят версии Python и pip:
    ```cmd
    python --version
    pip --version
    ```
    Если вы видите ошибки, убедитесь, что Python был добавлен в PATH. Возможно, потребуется перезапустить командную строку или даже компьютер.

## 2. Создание виртуального окружения (рекомендуется)

Использование виртуального окружения помогает изолировать зависимости проекта.

1.  **Перейдите в директорию вашего проекта:**
    Откройте командную строку и перейдите в корневую директорию вашего проекта MCP-сервера (например, `cd C:\Users\vm-user\Documents\GitHub\pydentycai-mcp-server`).

2.  **Создайте виртуальное окружение:**
    ```cmd
    python -m venv venv
    ```
    Это создаст папку `venv` в вашей текущей директории.

3.  **Активируйте виртуальное окружение:**
    ```cmd
    .\venv\Scripts\activate
    ```
    Вы увидите `(venv)` перед вашей командной строкой, что указывает на активное виртуальное окружение.

## 3. Установка PydanticAI и зависимостей

После активации виртуального окружения установите необходимые библиотеки.

1.  **Установите Pydantic и Pydantic-Settings:**
    ```cmd
    pip install pydantic pydantic-settings
    ```
    PydanticAI, вероятно, будет использовать эти библиотеки для определения и загрузки конфигурации.

2.  **Установите PydanticAI (если доступно как отдельный пакет):**
    Если PydanticAI является отдельным пакетом, который можно установить через pip, используйте:
    ```cmd
    pip install pydanticai
    ```
    *Примечание: На момент написания PydanticAI может быть не публичным пакетом PyPI. В этом случае вам потребуется получить исходный код PydanticAI и установить его локально или использовать его напрямую из исходников в вашем проекте.*

    **Если PydanticAI не является пакетом PyPI:**
    Вам нужно будет разместить файлы PydanticAI в поддиректории вашего проекта (например, `pydantic_agents/`) и запускать их напрямую. Убедитесь, что все необходимые зависимости Python для PydanticAI установлены в вашем виртуальном окружении.

## 4. Проверка установки

После установки вы можете создать простой Python-скрипт для проверки Pydantic и Pydantic-Settings.

1.  **Создайте файл `test_pydantic.py`:**
    ```python
    # test_pydantic.py
    from pydantic import BaseModel
    from pydantic_settings import BaseSettings

    class AgentConfig(BaseSettings):
        name: str = "DefaultAgent"
        version: str = "1.0"
        api_key: str

    class User(BaseModel):
        id: int
        name: str = "John Doe"

    try:
        config = AgentConfig(api_key="your_secret_key")
        print(f"Pydantic-Settings Config: {config.model_dump_json(indent=2)}")

        user = User(id=123)
        print(f"Pydantic Model: {user.model_dump_json(indent=2)}")
        print("Pydantic and Pydantic-Settings installed successfully!")
    except Exception as e:
        print(f"Error during Pydantic test: {e}")

    ```

2.  **Запустите скрипт:**
    Убедитесь, что ваше виртуальное окружение активно.
    ```cmd
    python test_pydantic.py
    ```
    Вы должны увидеть вывод JSON-конфигурации и сообщения об успешной установке.

Теперь ваша среда Windows готова для работы с агентами PydanticAI, которые будут управляться MCP-сервером на Go.
