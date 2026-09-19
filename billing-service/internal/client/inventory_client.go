package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/souzannatha/Korp_Teste_Natha/billing-service/internal/model"
)

var (
	ErrInventoryUnavailable     = errors.New("inventory service unavailable")
	ErrInventoryTimeout         = errors.New("inventory service timed out")
	ErrInvalidInventoryResponse = errors.New("invalid response from inventory service")
	ErrProductNotFound          = errors.New("product not found")
	ErrInsufficientBalance      = errors.New("insufficient stock balance")
	ErrDeductionConflict        = errors.New("invoice already has a different stock deduction")
	ErrInventoryValidation      = errors.New("inventory rejected the request")
)

type InventoryClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewInventoryClient(baseURL string) (InventoryClient, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") ||
		parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return InventoryClient{}, errors.New("INVENTORY_URL must be a valid HTTP URL")
	}
	return InventoryClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				DialContext:         (&net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
		},
	}, nil
}

func (ic *InventoryClient) Close() {
	ic.httpClient.CloseIdleConnections()
}

func (ic *InventoryClient) GetProductByCode(ctx context.Context, code string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, ic.baseURL+"/product/"+url.PathEscape(code), nil)
	if err != nil {
		return err
	}
	response, err := ic.httpClient.Do(request)
	if err != nil {
		return inventoryConnectionError(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return inventoryResponseError(response)
	}
	var product struct {
		ID          int    `json:"id"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Balance     *int   `json:"balance"`
	}
	err = decodeResponse(response.Body, &product)
	if err != nil || product.ID <= 0 || product.Code != code ||
		product.Description == "" || product.Balance == nil || *product.Balance < 0 {
		return ErrInvalidInventoryResponse
	}
	return nil
}

func (ic *InventoryClient) DeductStock(ctx context.Context, deduction model.StockDeduction) error {
	content, err := json.Marshal(deduction)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, ic.baseURL+"/stock-deductions", bytes.NewReader(content))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := ic.httpClient.Do(request)
	if err != nil {
		return inventoryConnectionError(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return inventoryResponseError(response)
	}
	var result model.Response
	if err := decodeResponse(response.Body, &result); err != nil || result.Message != "stock deducted successfully" {
		return ErrInvalidInventoryResponse
	}
	return nil
}

func decodeResponse(body io.Reader, destination any) error {
	decoder := json.NewDecoder(io.LimitReader(body, 1024*1024))
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return ErrInvalidInventoryResponse
	}
	return nil
}

func inventoryConnectionError(err error) error {
	var networkErr net.Error
	if errors.Is(err, context.DeadlineExceeded) || errors.As(err, &networkErr) && networkErr.Timeout() {
		return ErrInventoryTimeout
	}
	return ErrInventoryUnavailable
}

func inventoryResponseError(response *http.Response) error {
	var result model.Response
	if err := decodeResponse(response.Body, &result); err != nil || result.Message == "" {
		return ErrInvalidInventoryResponse
	}
	switch response.StatusCode {
	case http.StatusBadRequest:
		return ErrInventoryValidation
	case http.StatusNotFound:
		return ErrProductNotFound
	case http.StatusConflict:
		if result.Message == ErrInsufficientBalance.Error() {
			return ErrInsufficientBalance
		}
		if result.Message == ErrDeductionConflict.Error() {
			return ErrDeductionConflict
		}
		return ErrInvalidInventoryResponse
	case http.StatusInternalServerError, http.StatusServiceUnavailable:
		return ErrInventoryUnavailable
	case http.StatusGatewayTimeout:
		return ErrInventoryTimeout
	default:
		return ErrInvalidInventoryResponse
	}
}
