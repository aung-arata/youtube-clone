#!/bin/bash

# Start PHP-FPM in background
php-fpm -D

# Wait for database to be ready
sleep 5

# Initialize database tables
php -r "require '/var/www/admin-service/vendor/autoload.php'; 
       \$db = new App\Service\DatabaseService(); 
       \$db->initializeTables();" 2>/dev/null || true

# Start Nginx in foreground
nginx -g 'daemon off;'
