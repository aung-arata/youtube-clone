# Multi-stage Dockerfile for Admin Service

FROM php:8.2-fpm-alpine AS base

# Install system dependencies
RUN apk add --no-cache \
    nginx \
    postgresql-dev \
    curl \
    && docker-php-ext-install pdo pdo_pgsql

# Configure PHP-FPM
RUN sed -i 's/listen = 127.0.0.1:9000/listen = 9000/g' /usr/local/etc/php-fpm.d/www.conf

# Install Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

WORKDIR /var/www/admin-service

# Copy application files
COPY . .

# Install dependencies (if any)
RUN composer install --no-dev --optimize-autoloader --no-interaction 2>/dev/null || true

# Set permissions
RUN chown -R www-data:www-data /var/www/admin-service/var \
    && chmod -R 775 /var/www/admin-service/var

# Configure Nginx
COPY docker/nginx.conf /etc/nginx/http.d/default.conf

# Create startup script
COPY docker/start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 8085

CMD ["/start.sh"]
