FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS base
WORKDIR /app
EXPOSE 8080

FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

COPY ["src/Estoque.API/Estoque.API.csproj", "src/Estoque.API/"]
COPY ["src/Estoque.Application/Estoque.Application.csproj", "src/Estoque.Application/"]
COPY ["src/Estoque.Domain/Estoque.Domain.csproj", "src/Estoque.Domain/"]
COPY ["src/Estoque.Infrastructure/Estoque.Infrastructure.csproj", "src/Estoque.Infrastructure/"]

RUN dotnet restore "src/Estoque.API/Estoque.API.csproj"

COPY . .

RUN dotnet build "src/Estoque.API/Estoque.API.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "src/Estoque.API/Estoque.API.csproj" -c Release -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .

ENV ASPNETCORE_URLS=http://+:8080
ENV ASPNETCORE_ENVIRONMENT=Production
ENV ASPNETCORE_HTTP_PORTS=""

ENTRYPOINT ["dotnet", "Estoque.API.dll"]
