FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS base
WORKDIR /app
EXPOSE 80
EXPOSE 443

FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /src
COPY ["Kodiiing/Kodiiing.csproj", "Kodiiing/"]
RUN dotnet restore "Kodiiing/Kodiiing.csproj"
COPY . .
WORKDIR "/src/Kodiiing"
RUN dotnet build "Kodiiing.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "Kodiiing.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Kodiiing.dll"]
