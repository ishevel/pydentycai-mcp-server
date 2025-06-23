# PydanticAI MCP Server in Go

This project is a lightweight MCP (Model Context Protocol) server written in Go, designed for remote management of PydanticAI agents. It allows starting, stopping, getting status, and updating the configuration of PydanticAI agents via a standardized MCP protocol.

## Overview

The server provides a set of tools that enable interaction with PydanticAI agents, managing their lifecycle and configuration. Agent configurations are stored in text files (JSON) that are updated by MCP commands.

## MCP Server Tools

The MCP server provides the following tools for managing PydanticAI agents:

*   **`run_agent`**
    *   **Description:** Starts a PydanticAI agent with a specified ID.
    *   **Input Parameters:** `agent_id` (string, required) - Unique identifier for the agent to be launched.

*   **`stop_agent`**
    *   **Description:** Stops a running PydanticAI agent by its ID.
    *   **Input Parameters:** `agent_id` (string, required) - Unique identifier for the agent to be stopped.

*   **`update_agent_config`**
    *   **Description:** Updates the configuration file for a PydanticAI agent.
    *   **Input Parameters:** `agent_id` (string, required) - Unique identifier for the agent whose configuration needs to be updated; `config_data` (string, required) - Configuration data in JSON format (string).

*   **`get_agent_status`**
    *   **Description:** Returns the current status of a running PydanticAI agent by its ID.
    *   **Input Parameters:** `agent_id` (string, required) - Unique identifier for the agent whose status is to be retrieved.

*   **`list_agents`**
    *   **Description:** Returns a list of all registered agents and their statuses.
    *   **Input Parameters:** None.

## Connecting to MCP Clients (e.g., Cline for VS Code)

To connect our MCP server to a client like Cline in VS Code, you need to add the corresponding configuration to your `cline_mcp_settings.json` file (typically located at `C:/Users/vm-user/AppData/Roaming/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json`).

Since our server uses the Stdio transport (standard input/output), the configuration will look as follows:

```json
"pydentycai-mcp-server": {
  "autoApprove": [
    "run_agent",
    "stop_agent",
    "get_agent_status",
    "update_agent_config",
    "list_agents"
  ],
  "timeout": 60,
  "type": "stdio",
  "command": "./pydentycai-mcp-server.exe",
  "args": []
}
```

**Configuration Explanation:**
*   `"pydentycai-mcp-server"`: A unique name that will be used to identify our server in Cline.
*   `"autoApprove"`: A list of tools whose calls will be automatically approved by Cline. It is recommended to include all tools here to avoid constant confirmation prompts.
*   `"timeout"`: Maximum time in seconds to wait for a response from the server.
*   `"type": "stdio"`: Indicates that the server uses standard input/output for communication.
*   `"command": "./pydentycai-mcp-server.exe"`: The path to our Go server executable. It is assumed that the file is in the project's root directory.
*   `"args": []`: Additional command-line arguments for launching the server (none in this case).

After adding this configuration to `cline_mcp_settings.json` and restarting VS Code (or performing the relevant action in Cline to refresh servers), you will be able to call our server's tools directly from Cline.

## Running the Server

You can run the server using the following command in the project's root directory:

```bash
./pydentycai-mcp-server.exe
```

## Development and Testing

For development and testing, you can use standard Go commands:

*   `go mod tidy`: For managing dependencies.
*   `go build -o pydentycai-mcp-server.exe`: For building the executable.
