using System;
using System.Diagnostics;
using System.Net;
using System.Text.Json;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Hosting;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

string serviceName = "hello-world-csharp";
string runtimeLanguage = "C#";
string cloudProvider = "docker";
string region = "local";
string zone = "local";
string ipAddress = Dns.GetHostAddresses(Dns.GetHostName())[0].ToString();
int servicePid = Environment.ProcessId;
string K8sContainerID = "unknown";

void Log(string level, string operation, string message, long? duration = null, string outcome = "")
{
    var logEntry = new
    {
        timestamp = DateTime.UtcNow.ToString("o"),
        log_level = level,
        service_name = serviceName,
        service_operation = operation,
        runtime_language = runtimeLanguage,
        cloud_provider = cloudProvider,
        region = region,
        zone = zone,
        ip_address = ipAddress,
        service_pid = servicePid,
        container_id = containerId,
        call_duration = duration,
        call_outcome = outcome,
        message = message
    };
    Console.WriteLine(JsonSerializer.Serialize(logEntry));
}

Log("INFO", "start", "Service started successfully");

app.MapGet("/hello", async (HttpContext context) =>
{
    var stopwatch = Stopwatch.StartNew();
    Log("INFO", "hello", "Operation started");
    await context.Response.WriteAsync($"Hello, World! (Language: {runtimeLanguage}, Provider: {cloudProvider}, Region: {region}, Zone: {zone})");
    stopwatch.Stop();
    Log("INFO", "hello", "Operation completed", stopwatch.ElapsedMilliseconds, "pass");
});

app.MapPost("/inbound", async (HttpContext context) =>
{
    var stopwatch = Stopwatch.StartNew();
    Log("INFO", "inbound", "Operation started");
    await context.Response.WriteAsync(JsonSerializer.Serialize(new { status = "inbound request processed" }));
    stopwatch.Stop();
    Log("INFO", "inbound", "Operation completed", stopwatch.ElapsedMilliseconds, "approved");
});

app.MapPost("/outbound", async (HttpContext context) =>
{
    var stopwatch = Stopwatch.StartNew();
    Log("INFO", "outbound", "Operation started");
    await context.Response.WriteAsync(JsonSerializer.Serialize(new { status = "outbound request processed" }));
    stopwatch.Stop();
    Log("INFO", "outbound", "Operation completed", stopwatch.ElapsedMilliseconds, "declined");
});

app.MapGet("/status", async (HttpContext context) =>
{
    Log("INFO", "status", "Health check operation started");
    await context.Response.WriteAsync(JsonSerializer.Serialize(new { status = "service running" }));
    Log("INFO", "status", "Health check operation completed", 0, "pass");
});

app.Run("http://0.0.0.0:9092");

