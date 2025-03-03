# Simple Cloudflare DDNS Updater

Минималистичный DDNS клиент для Cloudflare на Go

## Features
- Одна бинарная сборка
- Автозапуск через systemd
- Проверка IP каждые 10 минут
- Логирование в systemd journal

## Установка
sudo yum install -y golang
git clone https://github.com/yourusername/cloudflare-ddns
cd cloudflare-ddns

# Сборка
go build -o cloudflare-ddns

# Установка
sudo mv cloudflare-ddns /usr/local/bin/
sudo mkdir -p /var/lib/cloudflare-ddns
sudo systemctl daemon-reload
sudo systemctl enable --now cloudflare-ddns

# проверка работы сервиса
sudo systemctl status cloudflare-ddns
journalctl -u cloudflare-ddns -f
