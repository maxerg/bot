# Архитектура (черновик)

## Слои (Go)
- cmd/ — точки входа (api)
- internal/domain — сущности/правила
- internal/app — use-cases
- internal/repository — интерфейсы/реализации хранилищ
- internal/transport/http — handlers + middleware
- internal/telegram — validate initData + Telegram API
- internal/crypto — AES-GCM, ключи, rotation

## Multi-tenant
tenant определяется доменом (через index.php виртуалхоста) и пробрасывается во фронт.
Дальше tenant будет идти в Go API через заголовок/параметр (добавим позже).

## Важно про “скрыть пакеты”
Шифрование в приложении скрывает payload от прокси/логов и случайного просмотра,
но **пользователь со своим браузером** всё равно сможет восстановить данные (ключ в клиенте).
