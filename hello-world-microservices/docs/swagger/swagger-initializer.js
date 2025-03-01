window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "http://localhost:8080/openapi/golang.yaml",  // ✅ Auto-load your OpenAPI spec
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
};
