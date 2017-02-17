FROM heroku/cedar

RUN cd /tmp && git clone https://github.com/heroku/heroku-buildpack-{{ .framework }}
{{ if .environment }}
ENV {{ range $key, $value := .environment }}{{ $key }}={{ $value }} {{ end }}
{{ end }}
ENV HOME=/app
WORKDIR /app

COPY .cache/build/buildpack /var/build/cache/buildpack
COPY . /app
RUN chmod +x /app/entrypoint.sh

RUN /tmp/heroku-buildpack-{{ .framework }}/bin/compile /app /var/build/cache/buildpack

ENTRYPOINT ["/app/entrypoint.sh"]
