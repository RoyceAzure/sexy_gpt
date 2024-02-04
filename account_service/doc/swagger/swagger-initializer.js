window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    url: "sexy_gpt.swagger.json",
    dom_id: '#swagger-ui',
    deepLinking: true,
    requestInterceptor: (request) => {
      if (request.headers.Authorization) {
          request.headers.Authorization = `Bearer ${request.headers.Authorization}`;
      }
      return request;
  },
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
