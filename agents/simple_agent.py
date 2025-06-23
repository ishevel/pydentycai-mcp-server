# agents/simple_agent.py
import argparse
import json
import logging
from pydantic import BaseModel
from pydantic_settings import BaseSettings

# Настройка логирования
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class AgentConfig(BaseSettings):
    """
    Модель конфигурации агента PydanticAI.
    Использует Pydantic-Settings для загрузки из файла и переменных окружения.
    """
    name: str = "DefaultAgent"
    version: str = "1.0"
    log_level: str = "INFO"
    parameters: dict = {}

    class Config:
        env_file = ".env" # Опционально: для загрузки из .env файла
        env_file_encoding = "utf-8"

def load_config(config_path: str) -> AgentConfig:
    """
    Загружает конфигурацию агента из указанного JSON-файла.
    """
    try:
        with open(config_path, 'r', encoding='utf-8') as f:
            config_data = json.load(f)
        config = AgentConfig(**config_data)
        logging.info(f"Конфигурация успешно загружена из {config_path}")
        return config
    except FileNotFoundError:
        logging.error(f"Файл конфигурации не найден: {config_path}")
        raise
    except json.JSONDecodeError:
        logging.error(f"Ошибка декодирования JSON в файле: {config_path}")
        raise
    except Exception as e:
        logging.error(f"Неизвестная ошибка при загрузке конфигурации: {e}")
        raise

def main():
    parser = argparse.ArgumentParser(description="PydanticAI Simple Agent")
    parser.add_argument("--config", type=str, required=True,
                        help="Путь к JSON-файлу конфигурации агента.")
    args = parser.parse_args()

    try:
        config = load_config(args.config)
        logging.getLogger().setLevel(getattr(logging, config.log_level.upper(), logging.INFO))

        logging.info(f"Агент '{config.name}' (версия {config.version}) запущен.")
        logging.info(f"Параметры конфигурации: {config.parameters}")

        # Здесь может быть основная логика агента
        # Например, имитация работы
        import time
        for i in range(3):
            logging.info(f"Агент '{config.name}' выполняет задачу {i+1}...")
            time.sleep(config.parameters.get("interval_seconds", 1))
        logging.info(f"Агент '{config.name}' завершил работу.")

    except Exception as e:
        logging.critical(f"Критическая ошибка при запуске агента: {e}")
        exit(1)

if __name__ == "__main__":
    main()
