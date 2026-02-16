# Telegram WebApp Bot (multi-tenant)

Каркас под Telegram Web App + админку в браузере.

## Точки входа
- /var/www/html/bot/checkcall.ru/index.php -> tenant checkcall
- /var/www/html/bot/pulse-poker.ru/index.php -> tenant pulse

## Идея
- Frontend: HTML/CSS/JS (Telegram WebApp + обычный браузер для админки)
- Backend: Go (основной, Clean Architecture)
- PHP: тонкий вход под виртуалхосты

## Шифрование пакетов (план)
AES-256-GCM: браузер WebCrypto <-> Go crypto/aes.
