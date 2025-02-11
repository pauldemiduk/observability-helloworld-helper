using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using System;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddLogging();  // Ensure logging is available



// ✅ Register Swagger services BEFORE building the app
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

// ✅ Ensure Swagger is used correctly
if (app.Environment.IsDevelopment() || app.Environment.IsProduction())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

// Function to retrieve cloud metadata from environment variables or default values
string GetCloudMetadata(string variable, string defaultValue) => 
    Environment.GetEnvironmentVariable(variable) ?? defaultValue;

// Get runtime language from environment variables
string runtimeLanguage = GetCloudMetadata("RUNTIME_LANGUAGE", "C#");

// Route: "/hello"
app.MapGet("/hello", async (HttpContext context) =>
{
    var provider = GetCloudMetadata("CLOUD_PROVIDER", "docker");
    var region = GetCloudMetadata("CLOUD_REGION", "local");
    var zone = GetCloudMetadata("CLOUD_ZONE", "local");

    string response = $"Hello, World! (Language: {runtimeLanguage}, Provider: {provider}, Region: {region}, Zone: {zone})";

    // Log request with standardized format
    var logger = context.RequestServices.GetRequiredService<ILogger<Program>>();
    logger.LogInformation("👋 Received request at /hello | Language: {Language}, Provider: {Provider}, Region: {Region}, Zone: {Zone}", 
                          runtimeLanguage, provider, region, zone);

    await context.Response.WriteAsync(response);
});

// Get host and port from environment variables (default: 8082)
string host = GetCloudMetadata("SERVICE_HOST", "0.0.0.0");
string port = GetCloudMetadata("SERVICE_PORT", "8082");

// Construct correct URL format
string url = $"http://{host}:{port}";

var logger = app.Services.GetRequiredService<ILogger<Program>>();
logger.LogInformation("🚀 Starting Hello World {Language} service on {Url} | Provider: {Provider}, Region: {Region}, Zone: {Zone}",
                      runtimeLanguage, url, GetCloudMetadata("CLOUD_PROVIDER", "docker"), 
                      GetCloudMetadata("CLOUD_REGION", "local"), 
                      GetCloudMetadata("CLOUD_ZONE", "local"));

app.Run(url);  // Start the application
