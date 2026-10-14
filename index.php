<?php
ob_start(); // Zatrzymuje wyświetlanie HTML, dopóki skrypty nie zakończą pracy
require_once __DIR__ . '/includes/auth.php';

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

// ZMIENIONA LOGIKA:
// Sprawdzamy, czy plik o nazwie z paska adresu istnieje w folderze pages
$target_file = __DIR__ . '/pages/' . $page . '.php';

if ($page === 'start') {
    require __DIR__ . '/pages/home.php';
} elseif (file_exists($target_file)) {
    // Jeśli utworzyłeś już plik (np. login.php), system go załaduje
    require $target_file;
} else {
    // Jeśli pliku jeszcze nie ma, ładujemy zaślepkę
    require __DIR__ . '/pages/placeholder.php';
}

require __DIR__ . '/includes/footer.php';