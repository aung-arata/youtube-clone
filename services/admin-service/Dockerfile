# Symfony Admin Service Dockerfile

FROM php:8.2-cli AS base

# Install system dependencies
RUN apt-get update && apt-get install -y \
    git \
    unzip \
    libpq-dev \
    libzip-dev \
    && docker-php-ext-install pdo pdo_pgsql zip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

WORKDIR /var/www/admin-service

# Copy composer files first for better caching
COPY composer.json composer.lock* ./

# Install dependencies (will be minimal in production)
RUN composer install --no-dev --no-scripts --no-autoloader --prefer-dist || composer install --no-dev --no-scripts --prefer-dist

# Copy application code
COPY . .

# Optimize autoloader and run scripts
RUN composer dump-autoload --no-dev --classmap-authoritative

# Create var directory and set permissions
RUN mkdir -p var/cache var/log \
    && chmod -R 777 var

# Set up Symfony runtime
ENV APP_ENV=prod
ENV APP_DEBUG=0

EXPOSE 8085

# Use Symfony CLI server (or php built-in server)
CMD ["php", "-S", "0.0.0.0:8085", "-t", "public"]
