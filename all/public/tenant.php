<?php
declare(strict_types=1);

header('X-Content-Type-Options: nosniff');
header('X-Frame-Options: SAMEORIGIN');
header('Referrer-Policy: no-referrer');
header('Permissions-Policy: geolocation=(), microphone=(), camera=()');

$tenant = defined('APP_TENANT') ? (string)APP_TENANT : 'checkcall';

require __DIR__ . '/index.php';
