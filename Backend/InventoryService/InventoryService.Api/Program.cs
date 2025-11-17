using InventoryService.Api.Endpoints;
using InventoryService.Application.Services;
using InventoryService.Domain.Repositories;
using InventoryService.Infrastructure.Data;
using InventoryService.Infrastructure.Repositories;
using Microsoft.EntityFrameworkCore;
using Microsoft.OpenApi.Models;
using InventoryService.Api.Middleware;

var builder = WebApplication.CreateBuilder(args);

// Db
var cs = builder.Configuration.GetConnectionString("Default")
 ?? throw new InvalidOperationException("Connection string não encontrada.");

// registra o DbContext com o provider do PostgreSQL
builder.Services.AddDbContext<InventoryDbContext>(opt =>
    opt.UseNpgsql(cs, npgsql => npgsql.EnableRetryOnFailure())
);

builder.Services.AddScoped<IStockMovementRepository, StockMovementRepository>();
builder.Services.AddScoped<IProductRepository, ProductRepository>();
builder.Services.AddScoped<ProductService>();
builder.Services.AddScoped<StockService>();

//Configurar o CORS
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowFrontend", policy =>
    {
        policy
            .WithOrigins(
                "https://localhost:4200",  // Frontend Angular
                "https://localhost:5002"
            )
            .AllowAnyHeader()
            .AllowAnyMethod()
            .AllowAnyOrigin();
    });
});

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(options =>
{
    options.EnableAnnotations();
    options.SwaggerDoc("v1", new OpenApiInfo
    {
        Title = "Inventory API v1",
        Version = "v1"
    });
});

builder.Services.AddHealthChecks();

var app = builder.Build();

//Habilitar CORS globalmente
app.UseCors("AllowFrontend");

app.UseMiddleware<ErrorHandlingMiddleware>();

app.UseSwagger();
app.UseSwaggerUI(options =>
{
    options.SwaggerEndpoint("/swagger/v1/swagger.json", "Inventory API v1");
});

// Map endpoints
app.MapProductEndpoints();
app.MapStockEndpoints();
app.MapHealthEndpoints();

app.UseHttpsRedirection();
app.Run();