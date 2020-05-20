# setup

```bash
mkdir docker
cp -r ../laradock/nginx ./docker/nginx
cp -r ../laradock/php-fpm ./docker/php-fpm
cp ../laradock/docker-compose.yml docker-compose.yml
touch ./docker/php-fpm/startup.sh
```