using Estoque.Application.Interfaces;
using Estoque.Application.UseCases.Movimentacoes;
using Estoque.Application.UseCases.Produtos;
using Estoque.Infrastructure.Persistence;
using Estoque.Infrastructure.Persistence.Repositories;
using Microsoft.EntityFrameworkCore;
using Serilog;
using Estoque.API.HealthChecks;
using Estoque.API.Middleware;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;

var builder = WebApplication.CreateBuilder(args);

// Serilog
Log.Logger = new LoggerConfiguration()
    .ReadFrom.Configuration(builder.Configuration)
    .Enrich.FromLogContext()
    .WriteTo.Console()
    .CreateLogger();

builder.Host.UseSerilog();

// DbContext
var connectionString = builder.Configuration.GetConnectionString("SqlServer")
    ?? throw new InvalidOperationException("Connection string 'SqlServer' não configurada.");

builder.Services.AddDbContext<EstoqueDbContext>(options =>
    options.UseSqlServer(connectionString, sql =>
        sql.EnableRetryOnFailure(maxRetryCount: 3)));

// Repositórios e UoW
builder.Services.AddScoped<IProdutoRepository, ProdutoRepository>();
builder.Services.AddScoped<IMovimentacaoRepository, MovimentacaoRepository>();
builder.Services.AddScoped<IUnitOfWork, UnitOfWork>();

// Use cases
builder.Services.AddScoped<CriarProdutoUseCase>();
builder.Services.AddScoped<ObterProdutoUseCase>();
builder.Services.AddScoped<ListarProdutosUseCase>();
builder.Services.AddScoped<ListarProdutosPaginadoUseCase>();
builder.Services.AddScoped<DeduzirLoteUseCase>();
builder.Services.AddScoped<EstornarLoteUseCase>();
builder.Services.AddScoped<AdicionarEntradaUseCase>();

// Health checks
builder.Services.AddHealthChecks()
    .AddDbContextCheck<EstoqueDbContext>("sqlserver", tags: ["db", "ready"]);

builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new() { Title = "Estoque API", Version = "v1" });
    c.EnableAnnotations();
    c.DocumentFilter<Estoque.API.Infrastructure.HealthCheckSchemaFilter>();
});

builder.Services.AddCors(o => o.AddPolicy("AllowAll", p =>
    p.AllowAnyOrigin().AllowAnyMethod().AllowAnyHeader()));

var app = builder.Build();

// Migrations automáticas na inicialização
using (var scope = app.Services.CreateScope())
{
    var db = scope.ServiceProvider.GetRequiredService<EstoqueDbContext>();
    await db.Database.EnsureCreatedAsync();
    Log.Information("Migrations aplicadas com sucesso.");
}

app.UseMiddleware<ExceptionMiddleware>();
app.UseSerilogRequestLogging();
app.UseCors("AllowAll");

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

// Health check endpoint com resposta JSON detalhada
app.MapHealthChecks("/health", new()
{
    ResponseWriter = HealthCheckResponseWriter.WriteAsync
});

app.MapControllers();

Log.Information("Serviço de Estoque iniciado na porta {Port}", builder.Configuration["ASPNETCORE_URLS"]);

app.Run();
