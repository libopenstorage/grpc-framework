window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    urls: [{ name: "apis/example/apiv1/example.swagger.json", url: "../apis/example/apiv1/example.swagger.json" },{ name: "apis/hello/apiv1/hello.swagger.json", url: "../apis/hello/apiv1/hello.swagger.json" }],
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

  //</editor-fold>
};
