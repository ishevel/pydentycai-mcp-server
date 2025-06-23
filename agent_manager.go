// Package main provides functions for managing PydanticAI agents.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// JSONResult представляет собой результат, который может быть преобразован в JSON.
type JSONResult struct {
	Data map[string]any `json:"data"`
}

// ToMCPResult преобразует JSONResult в mcp.Result.
func (r JSONResult) ToMCPResult() mcp.Result {
	return mcp.Result{
		Meta: map[string]any{
			"data": r.Data,
		},
	}
}

const (
	configsDir = "configs/agents"
)

// AgentProcess представляет запущенный процесс агента PydanticAI.
type AgentProcess struct {
	PID        int
	Cmd        *exec.Cmd
	Status     string // "running", "stopped", "error"
	ConfigPath string
	StartTime  time.Time
	Error      string
}

// AgentManager управляет жизненным циклом агентов PydanticAI.
type AgentManager struct {
	agents map[string]*AgentProcess
	mu     sync.RWMutex // Мьютекс для защиты доступа к карте агентов
}

// NewAgentManager создает новый экземпляр AgentManager.
func NewAgentManager() *AgentManager {
	return &AgentManager{
		agents: make(map[string]*AgentProcess),
	}
}

// runAgent запускает агент PydanticAI.
func (am *AgentManager) RunAgent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments.(map[string]any)
	agentID := args["agent_id"].(string)
	configPath := filepath.Join(configsDir, fmt.Sprintf("%s.json", agentID))

	am.mu.Lock()
	defer am.mu.Unlock()

	if _, ok := am.agents[agentID]; ok && am.agents[agentID].Status == "running" {
		return nil, fmt.Errorf("агент '%s' уже запущен", agentID)
	}

	// Проверяем существование файла конфигурации
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл конфигурации для агента '%s' не найден по пути: %s", agentID, configPath)
	}

	// Путь к исполняемому файлу Python
	pythonExecutable := "python" // Или "python3" в зависимости от системы

	// Путь к скрипту агента
	agentScriptPath := filepath.Join("agents", "simple_agent.py") // Предполагаем, что simple_agent.py - это наш агент

	cmd := exec.CommandContext(ctx, pythonExecutable, agentScriptPath, "--config", configPath)

	// Перенаправляем stdout и stderr агента в логи Go-сервера
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("не удалось запустить агент '%s': %w", agentID, err)
	}

	agentProcess := &AgentProcess{
		PID:        cmd.Process.Pid,
		Cmd:        cmd,
		Status:     "running",
		ConfigPath: configPath,
		StartTime:  time.Now(),
	}
	am.agents[agentID] = agentProcess

	log.Printf("Агент '%s' запущен с PID %d. Конфигурация: %s\n", agentID, agentProcess.PID, configPath)

	// Запускаем горутины для чтения вывода агента
	go am.readPipe(stdoutPipe, fmt.Sprintf("[Агент %s - stdout]", agentID))
	go am.readPipe(stderrPipe, fmt.Sprintf("[Агент %s - stderr]", agentID))

	// Запускаем горутину для ожидания завершения процесса
	go func() {
		err := cmd.Wait()
		am.mu.Lock()
		defer am.mu.Unlock()
		if agentProcess, ok := am.agents[agentID]; ok {
			agentProcess.Status = "stopped"
			if err != nil {
				agentProcess.Status = "error"
				agentProcess.Error = err.Error()
				log.Printf("Агент '%s' (PID %d) завершился с ошибкой: %v\n", agentID, agentProcess.PID, err)
			} else {
				log.Printf("Агент '%s' (PID %d) успешно завершился.\n", agentID, agentProcess.PID)
			}
		}
	}()

	return &mcp.CallToolResult{
		Result: JSONResult{
			Data: map[string]any{
				"agent_id": agentID,
				"pid":      agentProcess.PID,
				"status":   agentProcess.Status,
				"message":  fmt.Sprintf("Агент '%s' запущен с PID %d.", agentID, agentProcess.PID),
			},
		}.ToMCPResult(),
	}, nil
}

// readPipe читает данные из pipe и логирует их.
func (am *AgentManager) readPipe(pipe io.ReadCloser, prefix string) {
	buf := make([]byte, 1024)
	for {
		n, err := pipe.Read(buf)
		if n > 0 {
			log.Printf("%s %s", prefix, string(buf[:n]))
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("Ошибка чтения из pipe %s: %v\n", prefix, err)
			}
			break
		}
	}
}

// stopAgent останавливает запущенный агент PydanticAI.
func (am *AgentManager) StopAgent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments.(map[string]any)
	agentID := args["agent_id"].(string)

	am.mu.Lock()
	defer am.mu.Unlock()

	agentProcess, ok := am.agents[agentID]
	if !ok || agentProcess.Status != "running" {
		return nil, fmt.Errorf("агент '%s' не запущен или не найден", agentID)
	}

	if err := agentProcess.Cmd.Process.Kill(); err != nil {
		agentProcess.Status = "error"
		agentProcess.Error = err.Error()
		return nil, fmt.Errorf("не удалось остановить агент '%s' (PID %d): %w", agentID, agentProcess.PID, err)
	}

	agentProcess.Status = "stopped"
	log.Printf("Агент '%s' (PID %d) остановлен.\n", agentID, agentProcess.PID)

	return &mcp.CallToolResult{
		Result: JSONResult{
			Data: map[string]any{
				"agent_id": agentID,
				"pid":      agentProcess.PID,
				"status":   agentProcess.Status,
				"message":  fmt.Sprintf("Агент '%s' (PID %d) остановлен.", agentID, agentProcess.PID),
			},
		}.ToMCPResult(),
	}, nil
}

// getAgentStatus возвращает статус агента PydanticAI.
func (am *AgentManager) GetAgentStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments.(map[string]any)
	agentID := args["agent_id"].(string)

	am.mu.RLock()
	defer am.mu.RUnlock()

	agentProcess, ok := am.agents[agentID]
	if !ok {
		return nil, fmt.Errorf("агент '%s' не найден", agentID)
	}

	return &mcp.CallToolResult{
		Result: JSONResult{
			Data: map[string]any{
				"agent_id":    agentID,
				"pid":         agentProcess.PID,
				"status":      agentProcess.Status,
				"config_path": agentProcess.ConfigPath,
				"start_time":  agentProcess.StartTime.Format(time.RFC3339),
				"error":       agentProcess.Error,
			},
		}.ToMCPResult(),
	}, nil
}

// updateAgentConfig обновляет конфигурацию агента PydanticAI.
func (am *AgentManager) UpdateAgentConfig(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments.(map[string]any)
	agentID := args["agent_id"].(string)
	configData, ok := args["config_data"].(string)
	if !ok {
		return nil, fmt.Errorf("отсутствует или некорректный параметр 'config_data'")
	}

	configPath := filepath.Join(configsDir, fmt.Sprintf("%s.json", agentID))

	// Создаем директорию, если она не существует
	if err := os.MkdirAll(configsDir, 0755); err != nil {
		return nil, fmt.Errorf("не удалось создать директорию для конфигураций: %w", err)
	}

	// Записываем данные конфигурации в файл
	err := os.WriteFile(configPath, []byte(configData), 0644)
	if err != nil {
		return nil, fmt.Errorf("не удалось записать конфигурацию в файл '%s': %w", configPath, err)
	}

	log.Printf("Конфигурация для агента '%s' успешно обновлена в файле: %s\n", agentID, configPath)

	return &mcp.CallToolResult{
		Result: JSONResult{
			Data: map[string]any{
				"agent_id":    agentID,
				"config_path": configPath,
				"message":     fmt.Sprintf("Конфигурация для агента '%s' успешно обновлена.", agentID),
			},
		}.ToMCPResult(),
	}, nil
}

// listAgents возвращает список всех зарегистрированных агентов и их статусы.
func (am *AgentManager) ListAgents(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	agentsList := make([]map[string]any, 0, len(am.agents))
	for id, agent := range am.agents {
		agentsList = append(agentsList, map[string]any{
			"agent_id":    id,
			"pid":         agent.PID,
			"status":      agent.Status,
			"config_path": agent.ConfigPath,
			"start_time":  agent.StartTime.Format(time.RFC3339),
			"error":       agent.Error,
		})
	}

	return &mcp.CallToolResult{
		Result: JSONResult{
			Data: map[string]any{
				"agents": agentsList,
			},
		}.ToMCPResult(),
	}, nil
}
