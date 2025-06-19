package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	"email-service/metrics"
)

type EmailRequest struct {
	To         string `json:"to" binding:"required"`
	Subject    string `json:"subject" binding:"required"`
	Body       string `json:"body" binding:"required"`
	Filename   string `json:"filename,omitempty"`
	Attachment string `json:"attachment,omitempty"`
}

type SenderConfig struct {
	GmailUser     string
	GmailPassword string
}

var (
	senderConfigs map[string]SenderConfig
	smtpHost      = "smtp.gmail.com"
	stmpPort      = "587"
)

func init() {
	senderConfigs = make(map[string]SenderConfig)
	senders := []string{"compras", "financeiro", "controle"}

	log.Println("Loading sender configurations...")

	for _, senderName := range senders {
		userEnvKey := fmt.Sprintf("SENDER_%s_USER", strings.ToUpper(senderName))
		passEnvKey := fmt.Sprintf("SENDER_%s_PASS", strings.ToUpper(senderName))

		user := os.Getenv(userEnvKey)
		pass := os.Getenv(passEnvKey)

		if user == "" || pass == "" {
			log.Fatalf("Environment variables %s and %s must be set for sender '%s'", userEnvKey, passEnvKey, senderName)
		}

		senderConfigs[senderName] = SenderConfig{
			GmailUser:     user,
			GmailPassword: pass,
		}
		log.Printf("Loaded configuration for sender: %s", senderName)
	}
}

func sendEmailHtmlFormat(config SenderConfig, to, subject, body, filename, attachment string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", config.GmailUser)
	m.SetHeader("To", to)
	//m.SetAddressHeader("Cc", "") #caso tenha algum
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if filename != "" && attachment != "" {
		decodedAttachment, err := base64.StdEncoding.DecodeString(attachment)
		if err != nil {
			return fmt.Errorf("failed to decode attachment: %w", err)
		}

		m.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(decodedAttachment)
			return err
		}))
	}

	// send email
	port, _ := strconv.Atoi(stmpPort)
	d := gomail.NewDialer(smtpHost, port, config.GmailUser, config.GmailPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Factory patterns for handler
func createEmailHandler(senderName string) gin.HandlerFunc {
	config, ok := senderConfigs[senderName]

	if !ok {
		log.Fatalf("No Configuration found for sender '%s'", senderName)
	}

	return func(c *gin.Context) {
		var request EmailRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := sendEmailHtmlFormat(config, request.To, request.Subject, request.Body, request.Filename, request.Attachment)

		if err != nil {
			metrics.EmailErrors.WithLabelValues(senderName).Inc()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		metrics.EmailsSent.WithLabelValues(senderName).Inc()
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Email From '%s' sent successfully", senderName)})
	}
}

func getIpHandler(c *gin.Context) {
	cmd := exec.Command("hostname", "-i")
	output, err := cmd.Output()

	hostname, _ := os.Hostname()

	var hostnameIP string
	if err == nil {
		hostnameIP = strings.TrimSpace(string(output))
	}

	var dialIP string
	if conn, err := net.Dial("udp", "8.8.8.8"); err == nil {
		defer conn.Close()
		dialIP = conn.LocalAddr().(*net.UDPAddr).IP.String()
	}

	c.JSON(http.StatusOK, gin.H{
		"pod_name":    hostname,
		"server_ip":   hostnameIP,
		"outbound_ip": dialIP,
		"client_ip":   c.ClientIP(),
	})
}

func main() {
	router := gin.Default()

	p := ginprometheus.NewPrometheus("email_service")
	p.Use(router)

	trustedProxies := []string{
		"192.168.1.0/24",
	}

	err := router.SetTrustedProxies(trustedProxies)
	if err != nil {
		log.Fatalf("error on creating trusted proxies: %v", err)
	}

	router.POST("/send-email-compras", createEmailHandler("compras"))
	router.POST("/send-email-financeiro", createEmailHandler("financeiro"))
	router.POST("/send-email-controle", createEmailHandler("controle"))

	router.GET("/get-ip", getIpHandler)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(router.Run("0.0.0.0:" + port))
}
