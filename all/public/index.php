<?php
declare(strict_types=1);

// tenant приходит из tenant.php
$tenant = $tenant ?? 'checkcall';

$path = parse_url($_SERVER['REQUEST_URI'] ?? '/', PHP_URL_PATH) ?: '/';

if (str_starts_with($path, '/dist/')) {
    $file = __DIR__ . '/../web/dist' . substr($path, 5);
    if (is_file($file)) {
        $ext = pathinfo($file, PATHINFO_EXTENSION);
        $map = [
            'js' => 'application/javascript; charset=utf-8',
            'css' => 'text/css; charset=utf-8',
            'html' => 'text/html; charset=utf-8',
            'png' => 'image/png',
            'jpg' => 'image/jpeg',
            'svg' => 'image/svg+xml',
            'json' => 'application/json; charset=utf-8',
        ];
        header('Content-Type: ' . ($map[$ext] ?? 'application/octet-stream'));
        readfile($file);
        exit;
    }
    http_response_code(404);
    echo 'Not Found';
    exit;
}

header('Content-Type: text/html; charset=utf-8');

$title = $tenant === 'pulse' ? 'Pulse Poker' : 'CheckCall';
$apiBase = getenv('API_BASE_URL') ?: '/api';
?>
<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title><?= htmlspecialchars($title, ENT_QUOTES, 'UTF-8') ?></title>
  <link rel="stylesheet" href="/dist/app.css" />
  <script src="https://telegram.org/js/telegram-web-app.js"></script>

  <script>
    window.__APP__ = {
      tenant: <?= json_encode($tenant, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES) ?>,
      apiBase: <?= json_encode($apiBase, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES) ?>,
    };
  </script>
</head>
<body>
  <div id="app">Loading…</div>
  <script src="/dist/app.js"></script>
</body>
</html>
