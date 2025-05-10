package services

import (
	"bytes"    // Для работы с буферами
	"context"  // Контексты
	"errors"   // Обработка ошибок
	"fmt"      // Форматирование строк
	"io"       // Ввод-вывод
	"net/http" // HTTP-клиент
	"time"     // Работа с временем

	"github.com/Misha-Glazunov/bank-api/internal/config" // Конфигурация
	"github.com/beevik/etree"                            // Парсинг XML
	"github.com/sirupsen/logrus"                         // Логирование
)

type centralBankServiceImpl struct {
	client *http.Client
	config *config.CentralCBConfig
	logger *logrus.Logger
}

func NewCentralBankService(cfg *config.Config, logger *logrus.Logger) CentralBankService {
	return &centralBankServiceImpl{
		client: &http.Client{Timeout: cfg.CentralCB.Timeout},
		config: &cfg.CentralCB,
		logger: logger,
	}
}

func buildSOAPRequest(cfg *config.CentralCBConfig) string {
	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
        <soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
            <soap12:Body>
                <KeyRate xmlns="http://web.cbr.ru/">
                    <fromDate>%s</fromDate>
                    <ToDate>%s</ToDate>
                </KeyRate>
            </soap12:Body>
        </soap12:Envelope>`, fromDate, toDate)
}

func sendRequest(soapRequest string, cfg *config.CentralCBConfig) ([]byte, error) {
	client := &http.Client{Timeout: cfg.Timeout}
	req, err := http.NewRequest(
		"POST",
		cfg.WSDLURL,
		bytes.NewBuffer([]byte(soapRequest)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://web.cbr.ru/KeyRate")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func parseXMLResponse(rawBody []byte) (float64, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(rawBody); err != nil {
		return 0, fmt.Errorf("ошибка парсинга XML: %v", err)
	}

	krElements := doc.FindElements("//diffgram/KeyRate/KR")
	if len(krElements) == 0 {
		return 0, errors.New("данные по ставке не найдены")
	}

	rateElement := krElements[0].FindElement("./Rate")
	if rateElement == nil {
		return 0, errors.New("тег Rate отсутствует")
	}

	rateStr := rateElement.Text()
	var rate float64
	if _, err := fmt.Sscanf(rateStr, "%f", &rate); err != nil {
		return 0, fmt.Errorf("ошибка конвертации ставки: %v", err)
	}

	return rate, nil
}

func (s *centralBankServiceImpl) GetCurrentRate(ctx context.Context) (float64, error) {
	// Исправленный вызов с передачей конфигурации
	soapRequest := buildSOAPRequest(s.config)          // <-- Добавлен s.config
	rawBody, err := sendRequest(soapRequest, s.config) // <-- Добавлен s.config

	if err != nil {
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}

	rate, err := parseXMLResponse(rawBody)
	if err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return rate + 5.0, nil
}
