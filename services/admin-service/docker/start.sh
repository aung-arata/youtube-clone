#!/bin/sh

# Start PHP-FPM in background
php-fpm -D

# Initialize database tables (in background, wait for DB to be ready)
(
    sleep 5
    php -r "require '/var/www/admin-service/vendor/autoload.php'; 
           \$db = new App\Service\DatabaseService(); 
           \$db->initializeTables();" 2>/dev/null || true
) &

# Start Nginx in foreground
nginx -g 'daemon off;'
