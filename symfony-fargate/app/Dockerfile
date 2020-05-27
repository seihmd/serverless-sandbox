FROM php:7.4.6-fpm-alpine3.11

ENV DEBIAN_FRONTEND noninteractive

# always run apt update when start and after add new source list, then clean up at end.
RUN set -xe; \
    apk update && \
    pecl channel-update pecl.php.net && \
    apk add \
      #
      #--------------------------------------------------------------------------
      # Mandatory Software's Installation
      #--------------------------------------------------------------------------
      #
      # Mandatory Software's such as ("mcrypt", "pdo_mysql", "libssl-dev", ....)
      # are installed on the base image 'laradock/php-fpm' image. If you want
      # to add more Software's or remove existing one, you need to edit the
      # base image (https://github.com/Laradock/php-fpm).
      #
      # next lines are here becase there is no auto build on dockerhub see https://github.com/laradock/laradock/pull/1903#issuecomment-463142846
      libzip-dev zip unzip && \
      docker-php-ext-configure zip && \
      # Install the zip extension
      docker-php-ext-install zip && \
      php -m | grep -q 'zip'

#RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/bin --filename=composer
COPY --from=composer /usr/bin/composer /usr/bin/composer
RUN composer global require hirak/prestissimo

WORKDIR /var/www/app
COPY . /var/www/app/

CMD composer install; php-fpm
EXPOSE 9000