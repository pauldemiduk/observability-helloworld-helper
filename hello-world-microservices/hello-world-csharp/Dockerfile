# Use official .NET 8 SDK as the build environment
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy only the project file first (improves caching)
COPY *.csproj ./

# Restore dependencies
RUN dotnet restore

# Copy the rest of the source code
COPY . .

# Publish the application
RUN dotnet publish -c Release -o /app/out
# RUN dotnet publish -c Release -o /app/build-out

# Use a lightweight .NET runtime for running the app
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime

# Set the working directory inside the runtime container
WORKDIR /app

# Copy the built application from the build stage
COPY --from=build /app/out .

# Set default environment variables
ENV SERVICE_HOST="0.0.0.0"
ENV SERVICE_PORT="9092"
ENV CLOUD_PROVIDER="docker"
ENV CLOUD_REGION="local"
ENV CLOUD_ZONE="local"
ENV RUNTIME_LANGUAGE="C#"

# Expose the port dynamically
EXPOSE ${SERVICE_PORT}

# Command to run the application
CMD ["dotnet", "hello-world-csharp.dll"]
# CMD ["dotnet", "/app/build-out/hello-world-csharp.dll"]
