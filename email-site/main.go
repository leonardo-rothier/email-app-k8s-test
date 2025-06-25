package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BackendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type ProxyEmailRequest struct {
	Sender  string `json:"sender" binding:"required"`
	To      string `json:"to" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Service Tester</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }

        .container {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            width: 100%;
            max-width: 800px;
            padding: 40px;
            backdrop-filter: blur(10px);
        }

        h1 {
            text-align: center;
            color: #333;
            margin-bottom: 30px;
            font-size: 2.5em;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }

        .tabs {
            display: flex;
            margin-bottom: 30px;
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);
        }

        .tab {
            flex: 1;
            padding: 15px;
            background: #f0f0f0;
            border: none;
            cursor: pointer;
            font-size: 16px;
            transition: all 0.3s ease;
            font-weight: 600;
        }

        .tab:hover {
            background: #e0e0e0;
        }

        .tab.active {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }

        .content {
            display: none;
            animation: fadeIn 0.5s ease;
        }

        .content.active {
            display: block;
        }

        @keyframes fadeIn {
            from {
                opacity: 0;
                transform: translateY(10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        .form-group {
            margin-bottom: 25px;
        }

        label {
            display: block;
            margin-bottom: 8px;
            color: #555;
            font-weight: 600;
            font-size: 14px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        input, textarea, select {
            width: 100%;
            padding: 15px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 16px;
            transition: all 0.3s ease;
            font-family: inherit;
            background-color: #fff;
        }

        input:focus, textarea:focus, select:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }

        textarea {
            resize: vertical;
            min-height: 150px;
        }

        button {
            width: 100%;
            padding: 15px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 10px;
            font-size: 18px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            text-transform: uppercase;
            letter-spacing: 1px;
            position: relative;
            overflow: hidden;
        }
        
        button:disabled {
            cursor: not-allowed;
            opacity: 0.7;
        }

        button:hover:not(:disabled) {
            transform: translateY(-2px);
            box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
        }

        button:active:not(:disabled) {
            transform: translateY(0);
        }

        .ip-info {
            background: #f8f9fa;
            border-radius: 10px;
            padding: 30px;
            text-align: center;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);
        }

        .ip-item {
            margin: 20px 0;
            padding: 20px;
            background: white;
            border-radius: 10px;
            box-shadow: 0 3px 10px rgba(0, 0, 0, 0.1);
        }

        .ip-label {
            font-size: 14px;
            color: #666;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 5px;
        }

        .ip-value {
            font-size: 24px;
            font-weight: 700;
            color: #333;
            font-family: 'Courier New', monospace;
        }

        .loading {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid #f3f3f3;
            border-top: 3px solid #667eea;
            border-radius: 50%;
            animation: spin 1s linear infinite;
            margin-left: 10px;
            vertical-align: middle;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }

        .message {
            padding: 15px;
            border-radius: 10px;
            margin-top: 20px;
            text-align: center;
            font-weight: 600;
            animation: slideIn 0.5s ease;
        }

        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateY(-10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        .success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }

        .error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }

        .icon {
            font-size: 24px;
            margin-right: 10px;
            vertical-align: middle;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📧 Email Service Tester</h1>
        
        <div class="tabs">
            <button class="tab active" onclick="showTab('email', event)">Enviar Email</button>
            <button class="tab" onclick="showTab('ip', event)">Informações de IP</button>
        </div>

        <div id="email" class="content active">
            <form id="emailForm">
                <div class="form-group">
                    <label for="sender">Remetente:</label>
                    <select id="sender" name="sender" required>
                        <option value="compras">Compras</option>
                        <option value="financeiro">Financeiro</option>
                        <option value="controle">Controle</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="to">Para (Email):</label>
                    <input type="email" id="to" name="to" required placeholder="destinatario@example.com">
                </div>
                <div class="form-group">
                    <label for="subject">Assunto:</label>
                    <input type="text" id="subject" name="subject" required placeholder="Assunto do email">
                </div>
                <div class="form-group">
                    <label for="body">Mensagem:</label>
                    <textarea id="body" name="body" required placeholder="Digite sua mensagem aqui..."></textarea>
                </div>
                <button type="submit">
                    <span id="sendButtonText">Enviar Email</span>
                </button>
            </form>
            <div id="emailMessage"></div>
        </div>

        <div id="ip" class="content">
            <button onclick="getIPInfo(event)">
                <span id="ipButtonText">Obter Informações do Email Service</span>
            </button>
            <div class="ip-info" id="ipInfo" style="display: none;">
                <div class="ip-item">
                    <div class="ip-label">Pod Name</div>
                    <div class="ip-value" id="podName">-</div>
                </div>
                <div class="ip-item">
                    <div class="ip-label">Server IP</div>
                    <div class="ip-value" id="serverIP">-</div>
                </div>
                <div class="ip-item">
                    <div class="ip-label">Outbound IP</div>
                    <div class="ip-value" id="outboundIP">-</div>
                </div>
                <div class="ip-item">
                    <div class="ip-label">Client IP</div>
                    <div class="ip-value" id="clientIP">-</div>
                </div>
            </div>
            <div id="ipMessage"></div>
        </div>
    </div>

    <script>
        function showTab(tabName, event) {
            const tabs = document.querySelectorAll('.tab');
            const contents = document.querySelectorAll('.content');
            
            tabs.forEach(tab => tab.classList.remove('active'));
            contents.forEach(content => content.classList.remove('active'));
            
            event.target.classList.add('active');
            document.getElementById(tabName).classList.add('active');
        }

        document.getElementById('emailForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const button = e.target.querySelector('button');
            const buttonText = document.getElementById('sendButtonText');
            const messageDiv = document.getElementById('emailMessage');
            
            const emailData = {
                sender: document.getElementById('sender').value,
                to: document.getElementById('to').value,
                subject: document.getElementById('subject').value,
                body: document.getElementById('body').value
            };

            buttonText.innerHTML = 'Enviando<span class="loading"></span>';
            button.disabled = true;
            messageDiv.innerHTML = '';

            try {
                const response = await fetch('/api/send-email', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(emailData)
                });

                const data = await response.json();

                if (response.ok) {
                    messageDiv.innerHTML = '<div class="message success"><span class="icon">✓</span>' + (data.message || 'Email enviado com sucesso!') + '</div>';
                    document.getElementById('to').value = '';
                    document.getElementById('subject').value = '';
                    document.getElementById('body').value = '';
                } else {
                    messageDiv.innerHTML = '<div class="message error"><span class="icon">✗</span>Erro ao enviar email: ' + (data.error || 'Erro desconhecido') + '</div>';
                }
            } catch (error) {
                messageDiv.innerHTML = '<div class="message error"><span class="icon">✗</span>Erro de conexão: ' + error.message + '</div>';
            } finally {
                buttonText.textContent = 'Enviar Email';
                button.disabled = false;
            }
        });

        async function getIPInfo(event) {
            const button = event.target.closest('button');
            const buttonText = document.getElementById('ipButtonText');
            const ipInfo = document.getElementById('ipInfo');
            const messageDiv = document.getElementById('ipMessage');
            
            buttonText.innerHTML = 'Obtendo informações<span class="loading"></span>';
            button.disabled = true;
            messageDiv.innerHTML = '';
            ipInfo.style.display = 'none';

            try {
                const response = await fetch('/api/get-ip');
                const data = await response.json();

                if (response.ok) {
                    document.getElementById('podName').textContent = data.pod_name || 'N/A';
                    document.getElementById('serverIP').textContent = data.server_ip || 'N/A';
                    document.getElementById('outboundIP').textContent = data.outbound_ip || 'N/A';
                    document.getElementById('clientIP').textContent = data.client_ip || 'N/A';
                    ipInfo.style.display = 'block';
                } else {
                    messageDiv.innerHTML = '<div class="message error"><span class="icon">✗</span>' + (data.error || 'Erro ao obter informações de IP') + '</div>';
                }
            } catch (error) {
                messageDiv.innerHTML = '<div class="message error"><span class="icon">✗</span>Erro de conexão: ' + error.message + '</div>';
            } finally {
                buttonText.textContent = 'Obter Informações do Email Service';
                button.disabled = false;
            }
        }
    </script>
</body>
</html>
`

// Config just to test pods
func createHTTPClient() *http.Client {
	transport := &http.Transport{
		DisableKeepAlives:   true,
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 0,
		IdleConnTimeout:     1 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}

func main() {
	r := gin.Default()

	// Serve the HTML page
	r.GET("/", func(c *gin.Context) {
		tmpl := template.Must(template.New("index").Parse(htmlTemplate))
		c.Header("Content-Type", "text/html")
		tmpl.Execute(c.Writer, nil)
	})

	r.POST("/api/send-email", func(c *gin.Context) {
		var proxyReq ProxyEmailRequest
		if err := c.ShouldBindJSON(&proxyReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}

		allowedSenders := map[string]bool{
			"compras":    true,
			"financeiro": true,
			"controle":   true,
		}
		if !allowedSenders[proxyReq.Sender] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender specified."})
			return
		}

		backendURL := fmt.Sprintf("http://email-service/send-email-%s", proxyReq.Sender)

		backendReq := BackendEmailRequest{
			To:      proxyReq.To,
			Subject: proxyReq.Subject,
			Body:    proxyReq.Body,
		}
		jsonData, _ := json.Marshal(backendReq)

		// Create a new HTTP client for each request (just for test purpose)
		client := createHTTPClient()

		// Forward the request to the dynamically determined backend endpoint.
		resp, err := client.Post(backendURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to connect to email service: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		// Read the response from the backend and forward it back to the browser.
		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		c.JSON(resp.StatusCode, result)
	})

	// Proxy endpoint for getting IP info
	r.GET("/api/get-ip", func(c *gin.Context) {
		client := createHTTPClient()

		resp, err := client.Get("http://email-service/get-ip")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to connect to email service: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		// Read the response
		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)

		c.JSON(resp.StatusCode, result)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.Run("0.0.0.0:80")
}
