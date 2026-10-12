<section class="container py-5">
  <nav aria-label="Okruszki">
    <ol class="breadcrumb small">
      <li class="breadcrumb-item"><a href="index.php">Strona główna</a></li>
      <li class="breadcrumb-item active" aria-current="page"><?= htmlspecialchars($current['title'], ENT_QUOTES, 'UTF-8') ?></li>
    </ol>
  </nav>
  <div class="placeholder-card p-4 p-md-5 mt-3">
    <?php if ($page === '404'): ?>
      <span class="badge text-bg-secondary mb-3">404</span>
      <h1 class="display-6 fw-semibold">Nie znaleziono strony</h1>
      <p class="lead text-secondary">Sprawdź adres lub wróć na stronę główną.</p>
    <?php else: ?>
      <span class="badge text-bg-secondary mb-3">Etap w przygotowaniu</span>
      <h1 class="display-6 fw-semibold"><?= htmlspecialchars($current['title'], ENT_QUOTES, 'UTF-8') ?></h1>
      <p class="lead text-secondary"><?= htmlspecialchars($current['description'], ENT_QUOTES, 'UTF-8') ?></p>
      <p>Ta podstrona jest już dostępna w nawigacji. Jej właściwa funkcjonalność zostanie dodana w jednym z kolejnych etapów projektu.</p>
    <?php endif; ?>
    <a href="index.php" class="btn btn-primary mt-3">Wróć na stronę główną</a>
  </div>
</section>
