<?php
session_start();
session_destroy();
// Przekierowanie folder wyżej do głównego index.php
header('Location: ../index.php'); 
exit;
?>