using BillingService.Api.Endpoints;
using BillingService.Api.Middleware;
using BillingService.Application.ContractsHttp;
using BillingService.Application.Interfaces;
using BillingService.Application.Services;
using BillingService.Domain.Interfaces;
using BillingService.Infrastructure.ClientsHttp;
using BillingService.Infrastructure.Data;
using BillingService.Infrastructure.Pdf;
using BillingService.Infrastructure.Repositories;
using Microsoft.EntityFrameworkCore;
using Microsoft.OpenApi.Models;
using QuestPDF.Infrastructure;

// Configure a licença gratuita Community
QuestPDF.Settings.License = LicenseType.Community;

var builder = WebApplication.CreateBuilder(args);

// HTTP Client for InventoryService + Circuit Breaker + Retry + Timeout Policies for HttpClient
builder.Services.AddHttpClient<IHttpInventoryClient, HttpInventoryClient>(client =>
{
    // base address should point to InventoryService (adjust host/port as needed)
    client.BaseAddress = new Uri(builder.Configuration["InventoryService:BaseUrl"] ?? "http://inventoryservice:5001/");
    client.Timeout = TimeSpan.FromSeconds(30);
})
.AddPolicyHandler(HttpPolicies.GetRetryPolicy())
.AddPolicyHandler(HttpPolicies.GetCircuitBreakerPolicy())
.AddPolicyHandler(HttpPolicies.GetTimeoutPolicy());

// Configuration
var configuration = builder.Configuration;
var connectionString = configuration.GetConnectionString("Default")
    ?? throw new InvalidOperationException("Connection string 'Default' not found.");

// DbContext
builder.Services.AddDbContext<BillingDbContext>(opts =>
    opts.UseSqlServer(connectionString, sqlServerOptions =>
        sqlServerOptions.EnableRetryOnFailure()));

// Repositories & Services
builder.Services.AddScoped<IInvoiceRepository, InvoiceRepository>();
builder.Services.AddScoped<InvoiceService>();

builder.Services.AddScoped<IInvoicePrinter, QuestPdfInvoicePrinter>();

//Configurar o CORS
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowFrontend", policy =>
    {
        policy
            .WithOrigins(
                "https://localhost:4200"  // Frontend Angular
            )
            .AllowAnyHeader()
            .AllowAnyMethod()
            .AllowAnyOrigin();
    });
});

// Middleware & Swagger
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new OpenApiInfo { Title = "Billing API", Version = "v1" });
});

builder.Services.AddHealthChecks();

var app = builder.Build();

//Habilitar CORS globalmente
app.UseCors("AllowFrontend");

// Global error handling middleware (register before endpoints)
app.UseMiddleware<ErrorHandlingMiddleware>();

app.UseSwagger();
app.UseSwaggerUI(c =>
{
    c.SwaggerEndpoint("/swagger/v1/swagger.json", "Billing API v1");
});

// Map endpoints
app.MapInvoiceEndpoints();
app.MapHealthEndpoints();

app.UseHttpsRedirection();
app.Run();