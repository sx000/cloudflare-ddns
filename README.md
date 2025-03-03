# Simple Cloudflare DDNS Updater

Минималистичный DDNS клиент для Cloudflare на Go

## Features
- Одна бинарная сборка
- Автозапуск через systemd
- Проверка IP каждые 10 минут
- Логирование в systemd journal

## Установка
```bash
sudo yum install -y golang
git clone https://github.com/sx000/cloudflare-ddns
cd cloudflare-ddns
```

# Сборка
```bash
go build -o cloudflare-ddns
```

# Установка
```bash
sudo mv cloudflare-ddns /usr/local/bin/
sudo mkdir -p /var/lib/cloudflare-ddns
sudo systemctl daemon-reload
sudo systemctl enable --now cloudflare-ddns
```

# проверка работы сервиса
```bash
sudo systemctl status cloudflare-ddns
journalctl -u cloudflare-ddns -f
```