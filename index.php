<?php
// Jedyny publiczny punkt wejscia aplikacji.
$routes = require __DIR__ . '/config/routes.php';
$page = $_GET['page'] ?? 'start';
if (!is_string($page) || !array_key_exists($page, $routes)) {
    http_response_code(404);
    $page = '404';
    $current = [
        'title' => 'Nie znaleziono strony',
        'group' => null,
        'description' => 'Podany adres nie istnieje.',
    ];
} else {
    $current = $routes[$page];
}

require __DIR__ . '/includes/header.php';
if ($page === 'start') {
    require __DIR__ . '/pages/home.php';
} else {
    require __DIR__ . '/pages/placeholder.php';
}
require __DIR__ . '/includes/footer.php';
