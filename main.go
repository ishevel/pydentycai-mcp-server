// Package main implements a simple MCP server for PydanticAI agents.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "pydentycai-mcp-server"
	serverVersion = "0.1.0"
)

func main() {
	agentManager := NewAgentManager()

	// Создаем новый MCP-сервер
	mcpServer := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithLogging(),
		server.WithRecovery(),
		server.WithInstructions("Этот сервер управляет агентами PydanticAI."),
	)

	// Добавляем инструменты для управления агентами
	mcpServer.AddTools(
		server.ServerTool{
			Tool: mcp.NewTool("run_agent",
				mcp.WithDescription("Запускает агент PydanticAI с указанным ID и конфигурацией."),
				mcp.WithString("agent_id", mcp.Required(), mcp.Description("Уникальный идентификатор агента.")),
			),
			Handler: agentManager.RunAgent,
		},
		server.ServerTool{
			Tool: mcp.NewTool("stop_agent",
				mcp.WithDescription("Останавливает запущенный агент PydanticAI по его ID."),
				mcp.WithString("agent_id", mcp.Required(), mcp.Description("Уникальный идентификатор агента.")),
			),
			Handler: agentManager.StopAgent,
		},
		server.ServerTool{
			Tool: mcp.NewTool("get_agent_status",
				mcp.WithDescription("Возвращает текущий статус агента PydanticAI по его ID."),
				mcp.WithString("agent_id", mcp.Required(), mcp.Description("Уникальный идентификатор агента.")),
			),
			Handler: agentManager.GetAgentStatus,
		},
		server.ServerTool{
			Tool: mcp.NewTool("update_agent_config",
				mcp.WithDescription("Обновляет конфигурационный файл для агента PydanticAI."),
				mcp.WithString("agent_id", mcp.Required(), mcp.Description("Уникальный идентификатор агента.")),
				mcp.WithString("config_data", mcp.Required(), mcp.Description("Данные конфигурации в формате JSON (строка).")),
			),
			Handler: agentManager.UpdateAgentConfig,
		},
		server.ServerTool{
			Tool: mcp.NewTool("list_agents",
				mcp.WithDescription("Возвращает список всех зарегистрированных агентов и их статусы."),
			),
			Handler: agentManager.ListAgents,
		},
	)

	// Определяем порт для HTTP-сервера
	httpPort := ":8080" // Можно изменить на другой порт, если нужно

	// Создаем HTTP-сервер
	httpServer := server.NewStreamableHTTPServer(mcpServer)

	// Запускаем HTTP-сервер в отдельной горутине
	go func() {
		log.Printf("Запуск MCP-сервера %s (версия %s) через HTTP на порту %s...\n", serverName, serverVersion, httpPort)
		if err := httpServer.Start(httpPort); err != nil {
			log.Fatalf("Ошибка запуска HTTP-сервера: %v", err)
		}
	}()

	// Ожидаем сигнала завершения (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Получен сигнал завершения, остановка сервера...")
	log.Println("Сервер остановлен.")
}
