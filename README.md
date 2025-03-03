# Simple Cloudflare DDNS Updater

Минималистичный DDNS клиент для Cloudflare на Go

## Features
- Одна бинарная сборка
- Автозапуск через systemd
- Проверка IP каждые 10 минут
- Логирование в systemd journal
- Простая JSON-конфигурация

## Требования
- Linux (тестировано на CentOS 7+)
- Go 1.15+ (только для сборки из исходников)
- API токен Cloudflare с правами DNS:Edit

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
sudo mkdir -p /etc/cloudflare-ddns
sudo cp cloudflare-ddns.conf.example /etc/cloudflare-ddns.conf
sudo vi /etc/cloudflare-ddns.conf
sudo chmod 600 /etc/cloudflare-ddns.conf

sudo mv cloudflare-ddns /usr/local/bin/
sudo mkdir -p /var/lib/cloudflare-ddns
sudo cp cloudflare-ddns.service /etc/systemd/system/
sudo systemctl daemon-reload
```

# проверка работы сервиса
```bash
# Статус сервиса
sudo systemctl status cloudflare-ddns

# Логи в реальном времени
journalctl -u cloudflare-ddns -f

# Проверка текущего IP
cat /var/lib/cloudflare-ddns/current_ip
```

# запустите сервис
```bash
sudo systemctl enable --now cloudflare-ddns
```