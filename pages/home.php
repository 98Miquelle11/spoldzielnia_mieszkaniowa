<section class="hero py-5">
  <div class="container py-lg-4">
    <p class="text-uppercase fw-semibold mb-2 text-muted small letter-spacing">Jedno osiedle. Wspólna przestrzeń.</p>
    <h1 class="display-5 fw-bold">Spółdzielnia Mieszkaniowa „Mirmiłowo”</h1>
    <p class="lead mt-3 mb-4 hero-lead">Informacje o życiu osiedla, terminach, inwestycjach i sprawach mieszkańców w jednym miejscu.</p>
    <a href="index.php?page=kalendarz" class="btn btn-primary btn-lg">Sprawdź kalendarz</a>
    <a href="index.php?page=e-mieszkaniec" class="btn btn-outline-primary btn-lg ms-2">E-mieszkaniec</a>
  </div>
</section>
<section class="container py-5" aria-labelledby="sekcje-tytul">
  <h2 class="h3 mb-2" id="sekcje-tytul">Najważniejsze sekcje</h2>
  <p class="text-secondary mb-4">Dostępna jest nawigacja i strony informacyjne. Funkcje będą wdrażane w kolejnych etapach.</p>
  <div class="row g-4">
    <?php
    $cards = [
        ['title' => 'Ogłoszenia', 'desc' => 'Wydarzenia, przetargi i oferty wynajmu.', 'slug' => 'imprezy'],
        ['title' => 'Kalendarz', 'desc' => 'Terminy, spotkania i sprawy pilne.', 'slug' => 'kalendarz'],
        ['title' => 'Inwestycje', 'desc' => 'Zmiany i prace na terenie osiedla.', 'slug' => 'inwestycje'],
        ['title' => 'Strefa mieszkańca', 'desc' => 'Dokumenty, regulaminy i czat.', 'slug' => 'druki'],
        ['title' => 'E-mieszkaniec', 'desc' => 'Sprawy i wnioski online — moduł planowany.', 'slug' => 'e-mieszkaniec'],
        ['title' => 'Kontakt', 'desc' => 'Jak skontaktować się ze spółdzielnią.', 'slug' => 'kontakt'],
    ];
    foreach ($cards as $card):
    ?>
      <div class="col-md-6 col-lg-4">
        <article class="card h-100 border-0 shadow-sm">
          <div class="card-body p-4">
            <h3 class="h5 card-title"><?= htmlspecialchars($card['title'], ENT_QUOTES, 'UTF-8') ?></h3>
            <p class="card-text text-secondary"><?= htmlspecialchars($card['desc'], ENT_QUOTES, 'UTF-8') ?></p>
            <a href="index.php?page=<?= urlencode($card['slug']) ?>" class="stretched-link link-primary">Przejdź do sekcji <span aria-hidden="true">→</span></a>
          </div>
        </article>
      </div>
    <?php endforeach; ?>
  </div>
</section>
