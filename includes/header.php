<?php require_once 'auth.php'; ?>
<!doctype html>
<html lang="pl">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="description" content="Spółdzielnia Mieszkaniowa Mirmiłowo — informacje dla mieszkańców i aktualności osiedla.">
  <title><?= htmlspecialchars($current['title'], ENT_QUOTES, 'UTF-8') ?> — Spółdzielnia Mieszkaniowa „Mirmiłowo”</title>
  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.8/dist/css/bootstrap.min.css" rel="stylesheet">
  <link href="assets/css/style.css" rel="stylesheet">
</head>
<body>
<header>
  <nav class="navbar navbar-expand-lg bg-white border-bottom" aria-label="Nawigacja główna">
    <div class="container">
      <a class="navbar-brand fw-bold" href="index.php">Mirmiłowo <span class="brand-subtitle d-block fw-normal">Spółdzielnia Mieszkaniowa</span></a>
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#menu" aria-controls="menu" aria-expanded="false" aria-label="Rozwiń menu">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="menu">
        <ul class="navbar-nav ms-auto gap-lg-1">
          <?php
          $groups = [];
          foreach ($routes as $slug => $data) {
              if ($slug === 'start') { continue; }
              $group = $data['group'];
              if ($group === null) {
                  $groups[] = ['kind' => 'single', 'slug' => $slug, 'data' => $data];
              } else {
                  if (!isset($groupIndexes[$group])) {
                      $groupIndexes[$group] = count($groups);
                      $groups[] = ['kind' => 'group', 'name' => $group, 'children' => []];
                  }
                  $groups[$groupIndexes[$group]]['children'][$slug] = $data;
              }
          }
          foreach ($groups as $entry):
          ?>
            <?php if ($entry['kind'] === 'single'): ?>
              <li class="nav-item">
                <a class="nav-link <?= $page === $entry['slug'] ? 'active' : '' ?>" href="index.php?page=<?= urlencode($entry['slug']) ?>" <?= $page === $entry['slug'] ? 'aria-current="page"' : '' ?>><?= htmlspecialchars($entry['data']['title'], ENT_QUOTES, 'UTF-8') ?></a>
              </li>
            <?php else: ?>
              <li class="nav-item dropdown">
                <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false"><?= htmlspecialchars($entry['name'], ENT_QUOTES, 'UTF-8') ?></a>
                <ul class="dropdown-menu">
                  <?php foreach ($entry['children'] as $slug => $data): ?>
                    <li><a class="dropdown-item <?= $page === $slug ? 'active' : '' ?>" href="index.php?page=<?= urlencode($slug) ?>"><?= htmlspecialchars($data['title'], ENT_QUOTES, 'UTF-8') ?></a></li>
                  <?php endforeach; ?>
                </ul>
              </li>
            <?php endif; ?>
          <?php endforeach; ?>
          
          <!-- TUTAJ DODANO PRZYCISK LOGOWANIA/WYLOGOWANIA -->
          <?php if (isLoggedIn()): ?>
              <li class="nav-item ms-lg-3">
                  <a class="nav-link text-danger fw-bold" href="src/logout.php">Wyloguj się</a>
              </li>
          <?php else: ?>
              <li class="nav-item ms-lg-3">
                  <!-- Zmiana w linijce poniżej: -->
                  <a class="nav-link text-primary fw-bold" href="index.php?page=login">Zaloguj</a>
              </li>
          <?php endif; ?>
          <!-- KONIEC PRZYCISKU LOGOWANIA -->

        </ul>
      </div>
    </div>
  </nav>
</header>
<main id="tresc">