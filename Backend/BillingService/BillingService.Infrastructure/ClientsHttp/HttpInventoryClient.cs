using System.Net.Http.Json;
using BillingService.Application.ContractsHttp;
using BillingService.Application.DTOs;
using Microsoft.Extensions.Logging;
using Polly.CircuitBreaker;


namespace BillingService.Infrastructure.ClientsHttp
{
    public class HttpInventoryClient : IHttpInventoryClient
    {
        private readonly HttpClient _http;
        private readonly ILogger<HttpInventoryClient> _logger;
        private DateTime _lastHealthCheck = DateTime.MinValue;
        private bool _lastHealthStatus = false;

        public HttpInventoryClient(HttpClient http, ILogger<HttpInventoryClient> logger)
        {
            _http = http;
            _logger = logger;
        }

        public async Task<bool> IsInventoryServiceAvailableAsync()
        {
            // reaproveita o health por 5 segundos
            if ((DateTime.Now - _lastHealthCheck).TotalSeconds < 5)
                return _lastHealthStatus;

            try
            {
                // Realiza a chamada ao endpoint de health check
                var response = await _http.GetAsync("/health");

                // Atualiza o timestamp do último health check
                _lastHealthCheck = DateTime.Now;

                // Se a resposta não for bem-sucedida, considera o serviço indisponível
                if (!response.IsSuccessStatusCode)
                {
                    _lastHealthStatus = false;
                    return false;
                }

                //Evitar erro por conta de resposta vazia
                var content = await response.Content.ReadAsStringAsync();
                // Considera o serviço disponível se o conteúdo não for vazio
                _lastHealthStatus = !string.IsNullOrWhiteSpace(content);

                // Retorna o status do último health check
                return _lastHealthStatus;
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Erro ao verificar a disponibilidade do InventoryService.");
                return false;
            }
        }

        public async Task<ItemDto?> GetProductByCodeAsync(string productCode)
        {
            if (!await IsInventoryServiceAvailableAsync())
                throw new InvalidOperationException("InventoryService está offline.");

            try
            {
                // Consulta o produto pelo código
                var response = await _http.GetAsync($"/products/{productCode}");

                // Se o produto não for encontrado, retorna null
                if (response.StatusCode == System.Net.HttpStatusCode.NotFound)
                    return null;

                // Lança uma exceção para outros erros HTTP
                response.EnsureSuccessStatusCode();

                // Desserializa o conteúdo da resposta para ItemDto
                var item = await response.Content.ReadFromJsonAsync<ItemDto>();
                if (item == null)
                    throw new Exception($"Produto {productCode} retornou vazio do InventoryService.");

                return item;
            }
            catch (BrokenCircuitException)
            {
                // Circuit breaker aberto
                throw new InvalidOperationException(
                    "InventoryService está instável (circuit breaker aberto). Tente novamente em instantes.");
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Erro ao consultar o produto {ProductCode} no estoque.", productCode);
                throw;
            }
        }

        public async Task AdjustStockAsync(string productCode, int quantity, string operationKey)
        {
            if (!await IsInventoryServiceAvailableAsync())
                throw new InvalidOperationException("InventoryService está offline.");

            try
            {
                if (quantity == 0)
                    throw new ArgumentOutOfRangeException(nameof(quantity), "Quantity must be non-zero.");

                var movementType = quantity > 0 ? "add" : "remove";

                // Envia a requisição para ajustar o estoque
                var response = await _http.PostAsJsonAsync($"/stock/{movementType}", new
                {
                    productCode,
                    quantity = Math.Abs(quantity),
                    operationKey
                });

                // Verifica se a resposta indica sucesso
                if (!response.IsSuccessStatusCode)
                {
                    var error = await response.Content.ReadAsStringAsync();
                    throw new InvalidOperationException($"InventoryService error: {error}");
                }
            }
            catch (BrokenCircuitException)
            {
                // Circuit breaker aberto
                throw new InvalidOperationException(
                    "InventoryService está instável (circuit breaker aberto). Tente novamente em instantes.");
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, $"Erro ao movimentar o produto:{productCode} no estoque.");
                throw;
            }
        }
    }
}
