using Polly;
using Polly.Extensions.Http;

namespace BillingService.Api.Middleware
{
    public static class HttpPolicies
    {
        // Retry Policy
        public static IAsyncPolicy<HttpResponseMessage> GetRetryPolicy()
        {
            return HttpPolicyExtensions
                .HandleTransientHttpError()
                .WaitAndRetryAsync(
                    3,
                    retryAttempt => TimeSpan.FromSeconds(Math.Pow(2, retryAttempt))
                );
        }

        // Circuit Breaker Policy
        public static IAsyncPolicy<HttpResponseMessage> GetCircuitBreakerPolicy()
        {
            return HttpPolicyExtensions
                .HandleTransientHttpError()
                .CircuitBreakerAsync(
                    handledEventsAllowedBeforeBreaking: 3,
                    durationOfBreak: TimeSpan.FromSeconds(15)
                );
        }

        // Timeout Policy (opcional)
        public static IAsyncPolicy<HttpResponseMessage> GetTimeoutPolicy(int seconds = 10)
        {
            return Policy.TimeoutAsync<HttpResponseMessage>(TimeSpan.FromSeconds(seconds));
        }
    }
}
