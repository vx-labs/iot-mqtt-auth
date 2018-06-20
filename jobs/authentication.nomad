job "authentication" {
  datacenters = ["dc1"]
  type = "service"
  update {
    max_parallel = 1
    min_healthy_time = "10s"
    healthy_deadline = "3m"
    health_check = "checks"
    auto_revert = true
    canary = 0
  }
  group "app" {
    count = 1
    restart {
      attempts = 10
      interval = "5m"
      delay = "15s"
      mode = "delay"
    }
    ephemeral_disk {
      size = 300
    }
    task "store" {
      driver = "docker"
      env {
        CONSUL_HTTP_ADDR="172.17.0.1:8500"
      }
      config {
        image = "jbonachera/authentication:local"
        port_map {
          AuthenticationService = 7994
          health = 9000
        }
      }
      resources {
        cpu    = 500
        memory = 128
        network {
          mbits = 10
          port "AuthenticationService" {}
          port "health" {}
        }
      }
      service {
        name = "AuthenticationService"
        port = "AuthenticationService"
        tags=  ["leader"]
        check {
          type     = "http"
          path     = "/health"
          port     = "health"
          interval = "5s"
          timeout  = "2s"
        }
      }
    }
  }
}

