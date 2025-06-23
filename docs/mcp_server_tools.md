# Инструменты MCP-сервера для PydanticAI

Этот документ описывает набор инструментов (tools), которые MCP-сервер на Go будет предоставлять для взаимодействия с агентами PydanticAI. Эти инструменты позволят удаленно управлять жизненным циклом агентов и их конфигурацией, используя сериализованные данные.

## Обзор инструментов

MCP-сервер будет предоставлять следующие основные инструменты:

1.  **`run_agent`**: Запускает новый экземпляр агента PydanticAI.
2.  **`stop_agent`**: Останавливает запущенный экземпляр агента PydanticAI.
3.  **`update_agent_config`**: Обновляет конфигурацию существующего агента PydanticAI.
4.  **`get_agent_status`**: Получает текущий статус запущенного агента.

## Детали инструментов

### 1. `run_agent`

**Описание:** Запускает новый процесс агента PydanticAI с указанной конфигурацией. Сервер создаст или обновит файл конфигурации и передаст его путь агенту.

**Входные параметры (`arguments`):**

```json
{
  "agent_id": {
    "type": "string",
    "description": "Уникальный идентификатор для запускаемого агента. Используется для управления его конфигурацией и процессом."
  },
  "config_data": {
    "type": "object",
    "description": "Объект JSON, представляющий полную конфигурацию агента PydanticAI. Этот объект будет сериализован в файл."
  },
  "agent_script_path": {
    "type": "string",
    "description": "Относительный путь к Python-скрипту агента PydanticAI, который должен быть запущен. Например, 'agents/my_agent.py'."
  }
}
```

**Пример использования:**

```xml
<use_mcp_tool>
<server_name>pydanticai-mcp-server</server_name>
<tool_name>run_agent</tool_name>
<arguments>
{
  "agent_id": "my_first_agent",
  "config_data": {
    "name": "AgentAlpha",
    "version": "1.0",
    "log_level": "INFO",
    "parameters": {
      "threshold": 0.8,
      "interval_seconds": 60
    }
  },
  "agent_script_path": "agents/simple_agent.py"
}
</arguments>
</use_mcp_tool>
```

### 2. `stop_agent`

**Описание:** Останавливает запущенный процесс агента PydanticAI по его идентификатору.

**Входные параметры (`arguments`):**

```json
{
  "agent_id": {
    "type": "string",
    "description": "Уникальный идентификатор агента, который нужно остановить."
  }
}
```

**Пример использования:**

```xml
<use_mcp_tool>
<server_name>pydanticai-mcp-server</server_name>
<tool_name>stop_agent</tool_name>
<arguments>
{
  "agent_id": "my_first_agent"
}
</arguments>
</use_mcp_tool>
```

### 3. `update_agent_config`

**Описание:** Обновляет конфигурационный файл для указанного агента. Это не перезапускает агента автоматически; для применения изменений может потребоваться ручной перезапуск агента с помощью `stop_agent` и `run_agent`.

**Входные параметры (`arguments`):**

```json
{
  "agent_id": {
    "type": "string",
    "description": "Уникальный идентификатор агента, конфигурацию которого нужно обновить."
  },
  "config_data": {
    "type": "object",
    "description": "Объект JSON, представляющий новую полную конфигурацию агента PydanticAI. Этот объект полностью перезапишет существующий файл конфигурации."
  }
}
```

**Пример использования:**

```xml
<use_mcp_tool>
<server_name>pydanticai-mcp-server</server_name>
<tool_name>update_agent_config</tool_name>
<arguments>
{
  "agent_id": "my_first_agent",
  "config_data": {
    "name": "AgentAlpha",
    "version": "1.1",
    "log_level": "DEBUG",
    "parameters": {
      "threshold": 0.9,
      "interval_seconds": 30
    }
  }
}
</arguments>
</use_mcp_tool>
```

### 4. `get_agent_status`

**Описание:** Получает текущий статус запущенного агента, включая его PID и, возможно, последние строки логов или другую метаинформацию.

**Входные параметры (`arguments`):**

```json
{
  "agent_id": {
    "type": "string",
    "description": "Уникальный идентификатор агента, статус которого нужно получить."
  }
}
```

**Пример использования:**

```xml
<use_mcp_tool>
<server_name>pydanticai-mcp-server</server_name>
<tool_name>get_agent_status</tool_name>
<arguments>
{
  "agent_id": "my_first_agent"
}
</arguments>
</use_mcp_tool>
```

## Примечания по реализации

*   **Сериализация/Десериализация:** MCP-сервер будет использовать стандартную библиотеку Go `encoding/json` для сериализации `config_data` в файлы и десериализации при необходимости.
*   **Управление файлами:** Конфигурационные файлы будут храниться в выделенной директории (например, `configs/agents/`) внутри проекта, с именами, соответствующими `agent_id` (например, `configs/agents/my_first_agent.json`).
*   **Обработка ошибок:** Все инструменты должны возвращать информативные сообщения об ошибках в случае сбоев (например, агент не найден, ошибка записи файла, ошибка запуска процесса).
*   **PydanticAI:** Агенты PydanticAI должны быть спроектированы так, чтобы читать свою конфигурацию из файла, путь к которому будет передан им как аргумент командной строки при запуске. Pydantic будет использоваться внутри агентов для валидации этой конфигурации.
